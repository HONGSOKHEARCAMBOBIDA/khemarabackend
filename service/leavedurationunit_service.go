package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type LeaveDurationService interface {
	GetLeaveDuration() ([]model.LeaveDurationUnit, error)
}

type leavedurationservice struct {
	db *gorm.DB
}

func NewLeaveDurationService() LeaveDurationService {
	return &leavedurationservice{
		db: config.DB,
	}
}

func (s *leavedurationservice) GetLeaveDuration() ([]model.LeaveDurationUnit, error) {
	var leavedurationunit []model.LeaveDurationUnit
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	if err := tx.Find(&leavedurationunit).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return leavedurationunit, nil
}
