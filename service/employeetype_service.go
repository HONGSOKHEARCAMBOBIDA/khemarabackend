package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"

	"gorm.io/gorm"
)

type EmployeeTypeService interface {
	GetEmployeeType() ([]model.EmployeeType, error)
	CreateEmployeeType(input request.EmployeeTypeRequestCreate) error
	UpdateEmployeeType(id int, input request.EmployeeTypeRequestUpdate) error
	ChangeStatusEmployeeType(id int) error
}

type employeetypeservice struct {
	db *gorm.DB
}

func NewEmployeeTypeService() EmployeeTypeService {
	return &employeetypeservice{
		db: config.DB,
	}
}

func (s *employeetypeservice) GetEmployeeType() ([]model.EmployeeType, error) {
	var employeetype []model.EmployeeType
	if err := s.db.Find(&employeetype).Error; err != nil {
		return nil, err
	}
	return employeetype, nil
}

func (s *employeetypeservice) CreateEmployeeType(input request.EmployeeTypeRequestCreate) error {
	if input.Name == "" {
		return errors.New("name is required")
	}
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	newemployeetype := model.EmployeeType{
		Name:     input.Name,
		Isactive: true,
	}
	if err := tx.Create(&newemployeetype).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *employeetypeservice) UpdateEmployeeType(id int, input request.EmployeeTypeRequestUpdate) error {

	update := map[string]interface{}{}

	if input.Name != nil {
		update["name"] = *input.Name
	}

	if len(update) == 0 {
		return errors.New("no field to update")
	}

	result := s.db.Model(&model.EmployeeType{}).Where("id =?", id).Updates(update)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("employee type not found")
	}

	return nil

}

func (s *employeetypeservice) ChangeStatusEmployeeType(id int) error {
	result := s.db.Model(&model.EmployeeType{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("employee type not found")
	}
	return nil
}
