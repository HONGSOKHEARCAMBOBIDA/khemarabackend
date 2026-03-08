package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type RoleHasPermissionService interface {
	CreateRoleHasPermission(input request.RoleHasPermissionRequestCreate) error
	DeleteRoleHasPermission(input request.RoleHasPermissionRequestDelete) error
	GetRoleHasPermission(id int) ([]response.PermissionWithAssignedRole, error)
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

func (s *rolehaspermissionservice) GetRoleHasPermission(id int) ([]response.PermissionWithAssignedRole, error) {

	var rolehaspermission []response.PermissionWithAssignedRole

	err := s.db.Table("permissions p").
		Select(`
			p.id AS id,
			p.name AS name,
			p.display_name AS display_name,
			p.group_name AS group_name,
			p.short_name AS short_name,
			CASE 
				WHEN rhp.permission_id IS NULL THEN false
				ELSE true
			END AS assigned
		`).
		Joins("LEFT JOIN role_has_permissions rhp ON rhp.permission_id = p.id AND rhp.role_id = ?", id).
		Scan(&rolehaspermission).Error

	if err != nil {
		return nil, err
	}

	return rolehaspermission, nil
}
