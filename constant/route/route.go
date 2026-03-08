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

	// Employee Type
	ViewEmployeeType         = "view.employee.type"
	AddEmployeeType          = "add.employee.type"
	UpdateEmployeeType       = "update.employee.type/:id"
	ChangeStatusEmployeeType = "change.status.employee.type/:id"

	// Education Level
	ViewEducationLevel         = "view.education.level"
	AddEducationLevel          = "add.education.level"
	UpdateEducationLevel       = "update.education.level/:id"
	ChangeStatusEducationLevel = "change.status.education.level/:id"
)
