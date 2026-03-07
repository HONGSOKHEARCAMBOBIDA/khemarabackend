package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"

	"gorm.io/gorm"
)

type RoleHasPermissionService interface {
	CreateRoleHasPermission(input request.RoleHasPermissionRequestCreate) error
	DeleteRoleHasPermission(input request.RoleHasPermissionRequestDelete) error
}

type rolehaspermissionservice struct {
	db *gorm.DB
}

func NewRoleHasPermissionService() RoleHasPermissionService {
	return &rolehaspermissionservice{
		db: config.DB,
	}
}

func (s *rolehaspermissionservice) CreateRoleHasPermission(input request.RoleHasPermissionRequestCreate) error {

	if len(input.PermissionIDs) == 0 {
		return errors.New("permission id cannot be empty")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {

		var roleHasPermissions []model.RoleHasPermission

		for _, permissionID := range input.PermissionIDs {
			roleHasPermissions = append(roleHasPermissions, model.RoleHasPermission{
				RoleID:       input.RoleID,
				PermissionID: uint(permissionID),
			})
		}

		if err := tx.Create(&roleHasPermissions).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *rolehaspermissionservice) DeleteRoleHasPermission(input request.RoleHasPermissionRequestDelete) error {

	if len(input.PermissionIDs) == 0 {
		return errors.New("permission id cannot be empty")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {

		err := tx.
			Where("role_id = ? AND permission_id IN ?", input.RoleID, input.PermissionIDs).
			Delete(&model.RoleHasPermission{}).
			Error

		if err != nil {
			return err
		}

		return nil
	})
}
