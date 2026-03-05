package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"

	"gorm.io/gorm"
)

type RoleService interface {
	GetRole() ([]model.Role, error)
	CreateRole(input request.RoleRequestCreate) error
	UpdateRole(id int, input request.RoleRequestUpdate) error
	ChangeStatusRole(id int) error
}

type roleservice struct {
	db *gorm.DB
}

func NewRoleService() RoleService {
	return &roleservice{
		db: config.DB,
	}
}

func (s *roleservice) GetRole() ([]model.Role, error) {
	var roles []model.Role

	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *roleservice) CreateRole(input request.RoleRequestCreate) error {

	if input.Name == "" {
		return errors.New("name is required")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	role := model.Role{
		Name:        input.Name,
		DisPlayName: input.DisPlayName,
		IsActive:    true,
	}

	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *roleservice) UpdateRole(id int, input request.RoleRequestUpdate) error {

	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}

	if input.DisPlayName != nil {
		updates["display_name"] = *input.DisPlayName
	}

	if len(updates) == 0 {
		return errors.New("no field to update")
	}

	result := s.db.Model(&model.Role{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("role not found")
	}

	return nil
}

func (s *roleservice) ChangeStatusRole(id int) error {

	result := s.db.Model(&model.Role{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("role not found or status not changed")
	}

	return nil
}
