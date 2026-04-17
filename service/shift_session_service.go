package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type ShiftSessionService interface {
	GetAllShiftSession() ([]response.ShiftSessionResponse, error)
	GetByShiftID(id int) ([]response.ShiftSessionResponse, error)
	CreateShiftSession(input request.ShiftSessionRequestCreate) error
	UpdateShiftSession(id int, input request.ShiftSessionRequestUpdate) error
	ChangeStatusShiftSession(id int) error
}

type shiftsessionservice struct {
	db *gorm.DB
}

func Newshiftsessionservice() ShiftSessionService {
	return &shiftsessionservice{
		db: config.DB,
	}
}

func (s *shiftsessionservice) GetAllShiftSession() ([]response.ShiftSessionResponse, error) {
	var shiftsessions []response.ShiftSessionResponse
	db := s.db.Table("shift_sessions ss").
		Select(`
	ss.id AS id,
	ss.session_name AS session_name,
	ss.shift_order AS shift_order,
	ss.start_time AS start_time,
	ss.end_time AS end_time,
	ss.is_active AS is_active,
	s.id AS shift_id,
	s.name AS shift_name,
	b.id AS branch_id,
	b.  name AS branch_name
	`).
		Joins("LEFT JOIN shifts s ON s.id = ss.shift_id").
		Joins("LEFT JOIN branches b ON b.id = s.branch_id")
	if err := db.Order("ss.id DESC").Scan(&shiftsessions).Error; err != nil {
		return nil, err
	}
	return shiftsessions, nil
}

func (s *shiftsessionservice) GetByShiftID(id int) ([]response.ShiftSessionResponse, error) {
	var shiftsessions []response.ShiftSessionResponse
	db := s.db.Table("shift_sessions ss").
		Select(`
	ss.id AS id,
	ss.session_name AS session_name,
	ss.shift_order AS shift_order,
	ss.start_time AS start_time,
	ss.end_time AS end_time,
	ss.is_active AS is_active,
	s.id AS shift_id,
	s.name AS shift_name,
	b.id AS branch_id,
	b.name AS branch_name
	`).
		Joins("LEFT JOIN shifts s ON s.id = ss.shift_id").
		Joins("LEFT JOIN branches b ON b.id = s.branch_id")
	if err := db.Order("ss.id DESC").Where("s.id = ?", id).Scan(&shiftsessions).Error; err != nil {
		return nil, err
	}
	return shiftsessions, nil
}

func (s *shiftsessionservice) CreateShiftSession(input request.ShiftSessionRequestCreate) error {
	if input.SessionName == "" {
		return errors.New("session name is required")
	}
	if input.ShiftID == 0 {
		return errors.New("shift id is required")
	}
	if input.StartTime == "" {
		return errors.New("start time is required")
	}
	if input.EndTime == "" {
		return errors.New("end time is required")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	var lastSession model.ShiftSession
	err := tx.Where("shift_id = ?", input.ShiftID).Order("shift_order DESC").First(&lastSession).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	newShiftSession := model.ShiftSession{
		SessionName: input.SessionName,
		ShiftID:     input.ShiftID,
		ShiftOrder:  lastSession.ShiftOrder + 1,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Isactive:    true,
	}
	if err := tx.Create(&newShiftSession).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (s *shiftsessionservice) UpdateShiftSession(id int, input request.ShiftSessionRequestUpdate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()
	updates := map[string]interface{}{}
	if input.SessionName != nil {
		updates["session_name"] = *input.SessionName
	}
	if input.ShiftID != nil {
		updates["shift_id"] = *input.ShiftID
		var lastSession model.ShiftSession
		err := tx.Where("shift_id = ?", input.ShiftID).Order("shift_order DESC").First(&lastSession).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		updates["shift_id"] = *input.ShiftID
		if lastSession.ShiftID != id {
			updates["shift_order"] = lastSession.ShiftOrder + 1
		}

	}
	if input.StartTime != nil {
		updates["start_time"] = *input.StartTime
	}
	if input.EndTime != nil {
		updates["end_time"] = *input.EndTime
	}
	result := s.db.Model(&model.ShiftSession{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("shift session not found or already up to date")
	}
	return nil
}

func (s *shiftsessionservice) ChangeStatusShiftSession(id int) error {
	result := s.db.Model(&model.ShiftSession{}).Where("id = ?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("shift not found")
	}
	return nil
}
