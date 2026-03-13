package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type ProvinceService interface {
	GetProvince() ([]model.Province, error)
}

type provinceservice struct {
	db *gorm.DB
}

func NewProvinceService() ProvinceService {
	return &provinceservice{
		db: config.DB,
	}
}

func (s *provinceservice) GetProvince() ([]model.Province, error) {
	var province []model.Province
	if err := s.db.Find(&province).Error; err != nil {
		return nil, err
	}
	return province, nil
}
