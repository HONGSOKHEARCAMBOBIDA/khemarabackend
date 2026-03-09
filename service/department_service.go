package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"

	"gorm.io/gorm"
)

type DepartmentService interface {
	GetDepartment() ([]model.Department, error)
	CreateDepartment(input request.DepartmentRequestCreate) error
	UpdateDepartment(id int, input request.DepartmentRequestUpdate) error
	ChangeStatusDepartment(id int) error
}

type departmentservice struct {
	db *gorm.DB
}

func NewDepartmentService() DepartmentService {
	return &departmentservice{
		db: config.DB,
	}
}

func (s *departmentservice) GetDepartment() ([]model.Department, error) {
	var departments []model.Department
	if err := s.db.Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

func (s *departmentservice) CreateDepartment(input request.DepartmentRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.DisplayName == "" {
		return errors.New("display name is required")
	}
	newdepartment := model.Department{
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Isactive:    true,
	}
	if err := s.db.Create(&newdepartment).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *departmentservice) UpdateDepartment(id int, input request.DepartmentRequestUpdate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.DisplayName != nil {
		updates["display_name"] = *input.DisplayName
	}
	if len(updates) == 0 {
		return errors.New("no field to update")
	}
	result := s.db.Model(&model.Department{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("data no changed")
	}
	return nil

}

func (s *departmentservice) ChangeStatusDepartment(id int) error {
	result := s.db.Model(&model.Department{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
