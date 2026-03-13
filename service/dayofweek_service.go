package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type DayOfWeekService interface {
	GetDayOfWeek() ([]model.DayOfWeek, error)
}

type dayofweekservice struct {
	db *gorm.DB
}

func NewDayOfWeekService() DayOfWeekService {
	return &dayofweekservice{
		db: config.DB,
	}
}

func (s *dayofweekservice) GetDayOfWeek() ([]model.DayOfWeek, error) {
	var dayofweeks []model.DayOfWeek
	if err := s.db.Find(&dayofweeks).Error; err != nil {
		return nil, err
	}
	return dayofweeks, nil
}
