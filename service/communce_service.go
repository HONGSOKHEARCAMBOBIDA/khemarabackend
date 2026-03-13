package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type CommunceService interface {
	GetCommunce(id int) ([]model.Communce, error)
}

type communceservice struct {
	db *gorm.DB
}

func NewCommunceService() CommunceService {
	return &communceservice{
		db: config.DB,
	}
}

func (s *communceservice) GetCommunce(id int) ([]model.Communce, error) {
	var communces []model.Communce
	if err := s.db.Where("district_id =?", id).Find(&communces).Error; err != nil {
		return nil, err
	}
	return communces, nil
}
