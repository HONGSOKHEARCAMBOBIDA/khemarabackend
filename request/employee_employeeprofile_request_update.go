package request

type EmployeeEmpoyeeProfileRequestUpdate struct {
	NameEn         string `form:"name_en"`
	NameKh         string `form:"name_kh"`
	NationalID     string `form:"national_id_number"`
	Gender         int    `form:"gender"`
	PositionID     int    `form:"position_id"`
	HireDate       string `form:"hire_date" gorm:"column:hire_date"`
	PromoteDate    string `form:"promote_date"`
	IsPromote      bool   `form:"is_promote"`
	EmployeeTypeID int    `form:"employee_type_id"`
	Isactive       bool   `form:"is_active"`
	OfficeID       int    `form:"office_id"`

	PositionLevelID         int    `form:"position_level_id"`
	DoB                     string `form:"dob"`
	VillageIdOfBirth        int    `form:"village_id_of_birth"`
	MaterialStatus          string `form:"material_status"`
	VillageIDCurrentAddress int    `form:"village_id_current_address"`
	FamilyPhone             string `form:"family_phone"`
	BankName                string `form:"bank_name"`
	BankAccountNumber       string `form:"bank_account_number"`
	Note                    string `form:"note"`
	ReportoID               int    `form:"report_to"`
	WifeName                string `form:"wife_name"`
	HusBanName              string `form:"husban_name"`
	SonNumber               int    `form:"son_number"`
	DaughterNumber          int    `form:"daughter_number"`
}
