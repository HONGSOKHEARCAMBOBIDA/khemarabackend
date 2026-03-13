package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type VillageService interface {
	GetVillage(id int) ([]model.Village, error)
}

type villageservice struct {
	db *gorm.DB
}

func NewVillageService() VillageService {
	return &villageservice{
		db: config.DB,
	}
}

func (s villageservice) GetVillage(id int) ([]model.Village, error) {
	var village []model.Village

	if err := s.db.Where("communce_id =?", id).Find(&village).Error; err != nil {
		return nil, err
	}
	return village, nil
}
