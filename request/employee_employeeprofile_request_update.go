package request

type EmployeeEmpoyeeProfileRequestUpdate struct {
	// employee
	NameEn         string `json:"name_en"`
	NameKh         string `json:"name_kh"`
	NationalID     string `json:"national_id_number" gorm:"column:national_id_number"`
	Gender         int    `json:"gender"`
	PositionID     int    `json:"position_id"`
	HireDate       string `json:"hire_date"`
	PromoteDate    string `json:"promote_date"`
	IsPromote      bool   `json:"is_promote"`
	EmployeeTypeID int    `json:"employee_type_id"`
	Isactive       bool   `json:"is_active" gorm:"column:is_active"`
	OfficeID       int    `json:"office_id"`
	// employeeprofile
	PositionLevelID         int    `json:"position_level_id"`
	DoB                     string `json:"dob" gorm:"column:dob"`
	VillageIdOfBirth        int    `json:"village_id_of_birth"`
	MaterialStatus          string `json:"material_status"`
	VillageIDCurrentAddress int    `json:"village_id_current_address"`
	FamilyPhone             string `json:"family_phone"`
	BankName                string `json:"bank_name"`
	BankAccountNumber       string `json:"bank_account_number"`
	Note                    string `json:"note"`
	ReportoID               int    `json:"report_to" gorm:"column:report_to"`
	WifeName                string `json:"wife_name"`
	HusBanName              string `json:"husban_name" gorm:"column:husban_name"`
	SonNumber               int    `json:"son_number"`
	DaughterNumber          int    `json:"daughter_number"`
}
