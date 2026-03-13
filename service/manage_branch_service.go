package service

import (
	"mysql/config"
	"mysql/model"

	"gorm.io/gorm"
)

type ManageBranchService interface {
	GetManageBranch() ([]model.ManageBranch, error)
}

type managebranchservice struct {
	db *gorm.DB
}

func NewManageBranchService() ManageBranchService {
	return &managebranchservice{
		db: config.DB,
	}
}
func (s *managebranchservice) GetManageBranch() ([]model.ManageBranch, error) {
	var managebranch []model.ManageBranch
	if err := s.db.Find(&managebranch).Error; err != nil {
		return nil, err
	}
	return managebranch, nil
}
