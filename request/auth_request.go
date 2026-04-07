package request

type AuthRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	NameEn                  string  `form:"name_en" binding:"required"`
	NameKh                  string  `form:"name_kh" binding:"required"`
	NationalID              string  `form:"national_id_number" binding:"required"`
	Gender                  int     `form:"gender" binding:"required"`
	PositionID              int     `form:"position_id" binding:"required"`
	HireDate                string  `form:"hire_date"`
	PromoteDate             string  `form:"promote_date"`
	EmployeeTypeID          int     `form:"employee_type_id" binding:"required"`
	OfficeID                int     `form:"office_id" binding:"required"`
	Contact                 string  `form:"contact"`
	BranchID                int     `form:"branch_id" binding:"required"`
	RoleID                  int     `form:"role_id" binding:"required"`
	ManageBranch            int     `form:"manage_branch" binding:"required"`
	PositionLevelID         int     `form:"position_level_id"`
	DoB                     string  `form:"dob"`
	VillageIdOfBirth        int     `form:"village_id_of_birth"`
	MaterialStatus          string  `form:"material_status"`
	ProfileImage            string  `form:"profile_image"`
	VillageIDCurrentAddress int     `form:"village_id_current_address"`
	FamilyPhone             string  `form:"family_phone"`
	BankName                string  `form:"bank_name" binding:"required"`
	BankAccountNumber       string  `form:"bank_account_number" binding:"required"`
	QrCodeBankAccount       string  `form:"qr_code_bank_account" binding:"required"`
	Note                    string  `form:"note"`
	ReportoID               int     `form:"report_to"`
	WifeName                string  `form:"wife_name"`
	HusBanName              string  `form:"husban_name"`
	SonNumber               int     `form:"son_number"`
	DaughterNumber          int     `form:"daughter_number"`
	CompanyName             string  `form:"company_name"`
	PositionTitle           string  `form:"position_title"`
	StartDate               string  `form:"start_date"`
	EndDate                 string  `form:"end_date"`
	JobDescription          string  `form:"job_description"`
	EducationLevelID        int     `form:"education_level_id"`
	MajorField              string  `form:"major_field_of_study"`
	StartDateEducation      string  `form:"start_date_eud"`
	EndDateEducation        string  `form:"end_date_eud"`
	NoteEducation           string  `form:"note"`
	Image                   string  `form:"image"`
	PartIDs                 []int   `form:"part_ids"`
	BranchIDs               []int   `form:"branch_ids"`
	Dayofweeks              []int   `form:"day_of_weeks"`
	ShiftID                 int     `json:"shift_id"`
	Isdayoff                []bool  `json:"is_day_of"`
	BaseSalary              float64 `json:"base_salary" binding:"required"`
	WorkDay                 int     `json:"work_day" binding:"required"`
	DailyRate               float64 `json:"daily_rate" binding:"required"`
	CurrencyID              int     `json:"currency_id" binding:"required"`
}
