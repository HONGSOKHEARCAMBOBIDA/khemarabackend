package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"

	"gorm.io/gorm"
)

type BranchService interface {
	GetBranch() ([]model.Branch, error)
	CreateBranch(input request.BranchRequestCreate) error
	UpdateBranch(id int, input request.BranchRequesUpdate) error
	ChangeStatusBranch(id int) error
}

type branchservice struct {
	db *gorm.DB
}

func NewBranchService() BranchService {
	return &branchservice{
		db: config.DB,
	}
}

func (s *branchservice) GetBranch() ([]model.Branch, error) {
	var branch []model.Branch
	if err := s.db.Find(&branch).Error; err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *branchservice) CreateBranch(input request.BranchRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	newbranch := model.Branch{
		Name:      input.Name,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Radius:    input.Radius,
		Isactive:  true,
	}
	if err := s.db.Create(&newbranch).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *branchservice) UpdateBranch(id int, input request.BranchRequesUpdate) error {
	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Latitude != nil {
		updates["latitude"] = *input.Latitude
	}
	if input.Longitude != nil {
		updates["longitude"] = *input.Longitude
	}
	if input.Radius != nil {
		updates["radius"] = *input.Radius
	}
	if len(updates) == 0 {
		return errors.New(" no field to update")
	}
	result := s.db.Model(&model.Branch{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no data change")
	}
	return nil
}

func (s *branchservice) ChangeStatusBranch(id int) error {
	result := s.db.Model(&model.Branch{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
