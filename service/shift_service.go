package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type ShiftService interface {
	GetAllShift() ([]response.ShiftResponse, error)
	GetByBranchID(id int) ([]response.ShiftResponse, error)
	CreateShift(input request.ShiftRequestCreate) error
	UpdateShift(id int, input request.ShiftRequestUpdate) error
	ChangeStatusShift(id int) error
}

type shiftservice struct {
	db *gorm.DB
}

func NewShiftService() ShiftService {
	return &shiftservice{
		db: config.DB,
	}
}

func (s *shiftservice) GetAllShift() ([]response.ShiftResponse, error) {
	var shifts []response.ShiftResponse
	db := s.db.Table("shifts s").
		Select(`
		s.id AS id,
		s.name AS name,
		s.is_active AS is_active,
		b.id AS branch_id,
		b.name AS branch_name
	`).Joins("LEFT JOIN branches b ON b.id = s.branch_id")
	if err := db.Order("s.id DESC").Scan(&shifts).Error; err != nil {
		return nil, err
	}
	return shifts, nil
}

func (s *shiftservice) GetByBranchID(id int) ([]response.ShiftResponse, error) {
	var shifts []response.ShiftResponse
	db := s.db.Table("shifts s").
		Select(`
		s.id AS id,
		s.name AS name,
		s.is_active AS is_active,
		b.id AS branch_id,
		b.name AS branch_name
	`).Joins("LEFT JOIN branches b ON b.id = s.branch_id").Where("s.branch_id =?", id)
	if err := db.Order("s.id DESC").Scan(&shifts).Error; err != nil {
		return nil, err
	}
	return shifts, nil
}

func (s *shiftservice) CreateShift(input request.ShiftRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.BranchID == 0 {
		return errors.New("branch id is required")
	}
	newshift := model.Shift{
		Name:     input.Name,
		BranchID: input.BranchID,
		Isactive: true,
	}
	if err := tx.Create(&newshift).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *shiftservice) UpdateShift(id int, input request.ShiftRequestUpdate) error {
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.BranchID != nil {
		updates["branch_id"] = *input.BranchID
	}
	if len(updates) == 0 {
		return errors.New("no field to update")
	}
	result := s.db.Model(&model.Shift{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("shift not found or already up to date")
	}
	return nil
}

func (s *shiftservice) ChangeStatusShift(id int) error {
	result := s.db.Model(&model.Shift{}).Where("id = ?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("shift not found")
	}
	return nil
}
