package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type StatusAttendanceLogService interface {
	GetStatusAttendanceLogService() ([]model.StatusAttendanceLog, error)
}

type statusattendancelogservice struct {
	db *gorm.DB
}

func NewStatusAttendanceLogService() StatusAttendanceLogService {
	return &statusattendancelogservice{
		db: config.DB,
	}
}

func (s *statusattendancelogservice) GetStatusAttendanceLogService() ([]model.StatusAttendanceLog, error) {
	var statusattendancelog []model.StatusAttendanceLog
	if err := s.db.Find(&statusattendancelog).Error; err != nil {
		return nil, err
	}
	return statusattendancelog, nil
}
