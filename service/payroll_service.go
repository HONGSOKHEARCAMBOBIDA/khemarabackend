package service

import (
	"mysql/config"
	"mysql/request"

	"gorm.io/gorm"
)

type PayrollService interface {
	CreatePayroll(userID int, input request.PayrollRequestCreate) error
}

type payrollservice struct {
	db *gorm.DB
}

func NewPayrollService() PayrollService {
	return &payrollservice{
		db: config.DB,
	}
}

func (s *payrollservice) CreatePayroll(userID int, input request.PayrollRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
}
