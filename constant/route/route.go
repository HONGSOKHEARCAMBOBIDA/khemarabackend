package route

const (
	// Role
	ViewRole         = "viewrole"
	AddRole          = "addrole"
	EditRole         = "editrole/:id"
	ChangeStatusRole = "changestatusrole/:id"

	// RoleHasPermission
	AddPermissionTORole      = "add.permission.to.role"
	RemovePermissionFromRole = "remove.permission.from.role"
)
