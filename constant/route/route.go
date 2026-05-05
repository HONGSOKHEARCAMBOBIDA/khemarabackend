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

	// Branch
	ViewBranch         = "view.branch"
	AddBranch          = "add.branch"
	UpdateBranch       = "update.branch/:id"
	ChangeStatusBranch = "change.status.branch/:id"

	// Department
	ViewDepartment         = "view.department"
	AddDepartment          = "add.department"
	UpdateDepartment       = "update.department/:id"
	ChangeStatusDepartment = "change.status.department/:id"

	// Position
	ViewPosition             = "view.position"
	ViewPositionByDepartment = "view.position.bydepartment/:id"
	AddPosition              = "add.position"
	UpdatePosition           = "update.position/:id"
	ChangeStatusPosition     = "change.status.position/:id"

	// PositionLevel
	ViewPositionLevel         = "view.position.level"
	AddPositionLevel          = "add.position.level"
	UpdatePositionLevel       = "update.position.level/:id"
	ChangeStatusPositionLevel = "change.status.position.level/:id"

	// Currency
	ViewCurrency         = "view.currency"
	AddCurrency          = "add.currency"
	UpdateCurrency       = "upate.currency/:id"
	ChangeStatusCurrency = "change.status.currency/:id"

	// CurrencyPair
	ViewCurrencyPair         = "view.currency.pair"
	AddCurrencyPair          = "add.currency.pair"
	UpdateCurrencyPair       = "update.currency_pair/:id"
	ChangeStatusCurrencyPair = "change.status.currency.pair/:id"

	// ExchangeRate
	ViewExchangeRate         = "view.exchange.rate"
	AddExchangeRate          = "add.exchange.rate"
	UpdateExchangeRate       = "update.exchange.rate/:id"
	ChangeStatusExchangeRate = "change.status.exchange.rate/:id"

	// ManageBranch
	ViewManageBranch = "view.manage.branch"

	// Address
	ViewProvince = "view.province"
	ViewDistrict = "view.district/:id"
	ViewCommunce = "view.communce/:id"
	ViewVillage  = "view.village/:id"

	// Dayofweek
	ViewDayofweek = "view.day.of.week"

	// Office
	ViewOffice = "view.office"

	// Shift
	ViewAllShift        = "view.all.shift"
	ViewShfitByBranchID = "view.shift.by.branch.id/:id"
	AddShift            = "add.shift"
	UpdateShift         = "update.shift/:id"
	ChangeStatusShift   = "change.status.shift/:id"

	// ShiftSession
	ViewAllShiftSession       = "view.all.shift.session"
	ViewShiftSessionByShiftID = "view.shift.session.by.shift.id/:id"
	AddShiftSession           = "add.shift.session"
	UpdateShiftSession        = "update.shift.session/:id"
	ChangeStatusShiftSession  = "change.status.shift.session/:id"
	ViewShiftSessionV2        = "view.shift.session"

	// User
	ViewUserByBranch = "view.user.by.branch/:id"
	AddUser          = "add.user"
	UpdateUser       = "update.user/:id"
	ChangeStatusUser = "change.status.user/:id"

	// Part
	ViewPart = "view.part"

	// Employee
	ViewEmployee         = "view.employee"
	EditEmployee         = "edit.employee/:id"
	EditEducation        = "edit.education/:id"
	CreateEducation      = "add.education"
	EditWorkExperience   = "edit.work.experience/:id"
	CreateWorkExperience = "add.work.experience"
	EditSalary           = "edit.salary/:id"
	CreateSalary         = "add.salary"
	ChangeShiftPattern   = "edit.shift.pattern/:id"
	ChangeShift          = "edit.shift"

	// StatusAttendanceLog
	ViewStatusAttendance = "view.status.attendance"

	// Attendance
	AddAttendance  = "check.in"
	ViewAttendance = "view.attendance"
	CheckOut       = "check.out"

	// LeaveType
	ViewLeaveType = "view.leave.type"
)
