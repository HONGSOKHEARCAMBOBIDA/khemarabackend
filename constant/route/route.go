package route

const (
	// Role
	ViewRole         = "viewrole"
	AddRole          = "addrole"
	EditRole         = "editrole/:id"
	ChangeStatusRole = "changestatusrole/:id"

	// RoleHasPermission
	ViewRoleHasPermission    = "view.role.has.permission/:id"
	AddPermissionTORole      = "add.permission.to.role"
	RemovePermissionFromRole = "remove.permission.from.role"
)
