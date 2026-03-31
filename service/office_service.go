package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type OfficeService interface {
	GetAllOffice() ([]model.Office, error)
}

type officeservice struct {
	db *gorm.DB
}

func NewOfficeService() OfficeService {
	return &officeservice{
		db: config.DB,
	}
}

func (s *officeservice) GetAllOffice() ([]model.Office, error) {
	var offices []model.Office
	if err := s.db.Find(&offices).Error; err != nil {
		return nil, err
	}
	return offices, nil
}
