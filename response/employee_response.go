package response

type EmployeeResponseDetail struct {
	UserID                  int                              `json:"user_id"`
	UserName                string                           `json:"username" gorm:"column:username"`
	Contact                 string                           `json:"contact"`
	BranchID                int                              `json:"branch_id"`
	BranchName              string                           `json:"branch_name"`
	RoleID                  int                              `json:"role_id"`
	RoleName                string                           `json:"role_name"`
	RoleDisplayName         string                           `json:"role_display_name"`
	UserActive              bool                             `json:"user_active"`
	ManageBranchID          int                              `json:"manage_branch_id"`
	ManageBranchName        string                           `json:"manage_branch_name"`
	Parts                   []UserPartResponse               `json:"parts" gorm:"-"`
	Branches                []UserBranchResponse             `json:"branches" gorm:"-"`
	EmployeeRespons         []EmployeeRespons                `json:"employees" gorm:"-"`
	EmployeeEducations      []EmployeeEducationRespons       `json:"employeeeducation" gorm:"-"`
	EmployeeProfies         []EmployeeProfileResponse        `json:"employeeprofies" gorm:"-"`
	EmployeeWorkExperiences []EmployeeWorkExperienceResponse `json:"employeeworkexperiences" gorm:"-"`
	Salarys                 []SalaryResponse                 `json:"salarys" gorm:"-"`
	ShiftPatterns           []ShiftPatternResponse           `json:"shift_patterns" gorm:"-"`
}

type EmployeeRespons struct {
	ID                 int    `json:"id"`
	NameEn             string `json:"name_en"`
	NameKh             string `json:"name_kh"`
	NationalIDNumber   string `json:"national_id_number"`
	Gender             int    `json:"gender"`
	PositionID         int    `json:"position_id"`
	PositionName       string `json:"position_name"`
	HireDate           string `json:"hire_date"`
	PromoteDate        string `json:"promote_date"`
	IsPromote          bool   `json:"is_promote"`
	EmployeeTypeID     int    `json:"employee_type_id"`
	EmployeeTypeName   string `json:"employee_type_name"`
	Isactive           bool   `json:"is_active" gorm:"column:is_active"`
	OfficeID           int    `json:"office_id"`
	OffinceName        string `json:"office_name" gorm:"column:office_name"`
	OffinceDisplayName string `json:"office_display_name" gorm:"column:office_display_name"`
	CreateBy           int    `json:"create_by"`
	CreateByName       string `json:"create_by_name"`
}
