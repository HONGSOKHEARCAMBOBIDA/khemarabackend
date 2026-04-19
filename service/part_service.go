package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type PartService interface {
	GetPart() ([]model.Part, error)
}

type partservice struct {
	db *gorm.DB
}

func NewPartService() PartService {
	return &partservice{
		db: config.DB,
	}
}

func (s *partservice) GetPart() ([]model.Part, error) {
	var part []model.Part
	if err := s.db.Find(&part).Error; err != nil {
		return nil, err
	}
	return part, nil
}
