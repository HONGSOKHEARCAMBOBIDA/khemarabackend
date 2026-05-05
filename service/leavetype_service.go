package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type LeaveTypeService interface {
	GetLeaveType() ([]model.LeaveType, error)
}

type leavetypeservice struct {
	db *gorm.DB
}

func NewLeaveTypeService() LeaveTypeService {
	return &leavetypeservice{
		db: config.DB,
	}
}

func (s *leavetypeservice) GetLeaveType() ([]model.LeaveType, error) {
	var leavetype []model.LeaveType
	db := s.db.Table("leave_types lt").
		Select(`
	lt.id AS id,
	lt.name AS name,
	dt.id AS deduct_type_id,
	dt.code AS deduct_type_code,
	dt.name AS deduct_type_name,
	c.id AS currency_id,
	c.name AS currency_name,
	lt.amount AS amount,
	lt.description AS description,
	lt.is_active AS is_active
	`).
		Joins("LEFT JOIN deduct_types dt ON dt.id = lt.deduct_type_id").
		Joins("LEFT JOIN currencies c ON c.id = lt.currency_id")
	if err := db.Order("lt.id DESC").Scan(&leavetype).Error; err != nil {
		return nil, err
	}
	return leavetype, nil
}
