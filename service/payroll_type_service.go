package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type PayRollTypeService interface {
	GetPayrollType() ([]model.PayRollType, error)
}

type payrolltypeservice struct {
	db *gorm.DB
}

func NewPayRollTypeService() PayRollTypeService {
	return &payrolltypeservice{
		db: config.DB,
	}
}

func (s payrolltypeservice) GetPayrollType() ([]model.PayRollType, error) {
	var payrolltype []model.PayRollType
	if err := s.db.Find(&payrolltype).Error; err != nil {
		return nil, err
	}
	return payrolltype, nil
}
