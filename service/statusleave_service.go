package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type StatusLeaveService interface {
	GetStatusLeave() ([]model.StatusLeave, error)
}

type statusleaveservice struct {
	db *gorm.DB
}

func NewStatusLeaveService() StatusLeaveService {
	return &statusleaveservice{
		db: config.DB,
	}
}

func (s *statusleaveservice) GetStatusLeave() ([]model.StatusLeave, error) {
	var statusleave []model.StatusLeave
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	if err := tx.Find(&statusleave).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return statusleave, nil
}
