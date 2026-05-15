package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type BonusTypeService interface {
	GetBonusType() ([]model.BonusType, error)
}

type bonustypeservice struct {
	db *gorm.DB
}

func NewBonusTypeService() BonusTypeService {
	return &bonustypeservice{
		db: config.DB,
	}
}

func (s bonustypeservice) GetBonusType() ([]model.BonusType, error) {
	var bonustype []model.BonusType
	if err := s.db.Find(&bonustype).Error; err != nil {
		return nil, err
	}
	return bonustype, nil
}
