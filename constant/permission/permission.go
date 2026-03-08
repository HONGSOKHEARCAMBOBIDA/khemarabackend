package permission

const (
	// Role
	ViewRole         = "role.read"
	AddRole          = "role.create"
	EditRole         = "role.update"
	ChangeStatusRole = "role.change.status"

	// RoleHasPermission
	ViewRoleHasPermission    = "view.role.has.permission"
	AddPermissionTORole      = "add.permission.to.role"
	RemovePermissionFromRole = "remove.permission.from.role"
)
