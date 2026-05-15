package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type PayrollStatusService interface {
	GetPayrollStatus() ([]model.PayRollStatus, error)
}

type payrollstatusservice struct {
	db *gorm.DB
}

func NewPayrollStatusService() PayrollStatusService {
	return &payrollstatusservice{
		db: config.DB,
	}
}

func (s *payrollstatusservice) GetPayrollStatus() ([]model.PayRollStatus, error) {
	var payrollstatus []model.PayRollStatus
	if err := s.db.Find(&payrollstatus).Error; err != nil {
		return nil, err
	}
	return payrollstatus, nil
}
