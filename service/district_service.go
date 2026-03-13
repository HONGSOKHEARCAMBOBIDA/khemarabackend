package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type DistrictService interface {
	GetDistrict(id int) ([]model.District, error)
}

type districtservice struct {
	db *gorm.DB
}

func NewDistrictService() DistrictService {
	return &districtservice{
		db: config.DB,
	}
}

func (s *districtservice) GetDistrict(id int) ([]model.District, error) {
	var district []model.District
	if err := s.db.Where("province_id =?", id).Find(&district).Error; err != nil {
		return nil, err
	}
	return district, nil
}
