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

	// Employee Type
	ViewEmployeeType         = "view.employee.type"
	AddEmployeeType          = "add.employee.type"
	UpdateEmployeeType       = "update.employee.type"
	ChangeStatusEmployeeType = "change.status.employee.type"

	// Education Level
	ViewEducationLevel         = "view.education.level"
	AddEducationLevel          = "add.education.level"
	UpdateEducationLevel       = "update.education.level"
	ChangeStatusEducationLevel = "change.status.education.level"

	// Branch
	ViewBranch         = "view.branch"
	AddBranch          = "add.branch"
	UpdateBranch       = "update.branch"
	ChangeStatusBranch = "change.status.branch"

	// Department
	ViewDepartment         = "view.department"
	AddDepartment          = "add.department"
	UpdateDepartment       = "update.department"
	ChangeStatusDepartment = "change.status.department"

	// Position
	ViewPosition         = "view.position"
	AddPosition          = "add.position"
	UpdatePosition       = "update.position"
	ChangeStatusPosition = "change.status.position"

	// PositionLevel
	ViewPositionLevel         = "view.position.level"
	AddPositionLevel          = "add.position.level"
	UpdatePositionLevel       = "update.position.level"
	ChangeStatusPositionLevel = "change.status.position.level"
)
