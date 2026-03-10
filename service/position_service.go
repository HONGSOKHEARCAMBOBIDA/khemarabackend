package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type PositionService interface {
	GetAllPosition() ([]response.PositionResponse, error)
	GetByDepartmentID(id int) ([]response.PositionResponse, error)
	CreatePosition(input request.PositionRequestCreate) error
	UpdatePosition(id int, input request.PositionRequestUpdate) error
	ChangeStatusPosition(id int) error
}

type positionservice struct {
	db *gorm.DB
}

func NewPositionService() PositionService {
	return &positionservice{
		db: config.DB,
	}
}

func (s *positionservice) GetAllPosition() ([]response.PositionResponse, error) {
	var position []response.PositionResponse
	db := s.db.Table("positions p").
		Select(`
		p.id AS id,
		p.name AS name,
		p.display_name AS display_name,
		p.is_active AS is_active,
		d.id AS department_id,
		d.name AS department_name,
		d.display_name AS department_display_name
	`).
		Joins("LEFT JOIN departments d ON d.id = p.department_id")
	if err := db.Order("p.id DESC").Scan(&position).Error; err != nil {
		return nil, err
	}
	return position, nil
}

func (s *positionservice) GetByDepartmentID(id int) ([]response.PositionResponse, error) {
	var position []response.PositionResponse
	db := s.db.Table("positions p").
		Select(`
		p.id AS id,
		p.name AS name,
		p.display_name AS display_name,
		p.is_active AS is_active,
		d.id AS department_id,
		d.name AS department_name,
		d.display_name AS department_display_name
	`).
		Joins("LEFT JOIN departments d ON d.id = p.department_id").Where("d.id =?", id)
	if err := db.Order("p.id DESC").Scan(&position).Error; err != nil {
		return nil, err
	}
	return position, nil
}

func (s *positionservice) CreatePosition(input request.PositionRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.DisplayName == "" {
		return errors.New("display name is required")
	}
	if input.DepartmentID == 0 {
		return errors.New("department_id is required")
	}
	newposition := model.Position{
		Name:         input.Name,
		DisplayName:  input.DisplayName,
		DepartmentID: input.DepartmentID,
		Isactive:     true,
	}
	if err := tx.Create(&newposition).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *positionservice) UpdatePosition(id int, input request.PositionRequestUpdate) error {
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.DisplayName != nil {
		updates["display_name"] = *input.DisplayName
	}
	if input.DepartmentID != nil {
		updates["department_id"] = *input.DepartmentID
	}
	if len(updates) == 0 {
		return errors.New("no field to update")
	}
	result := s.db.Model(&model.Position{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no data update")
	}
	return nil
}

func (s *positionservice) ChangeStatusPosition(id int) error {
	result := s.db.Model(&model.Position{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no data update")
	}
	return nil
}
