package response

type EmployeeProfileResponse struct {
	ID                  int    `json:"id"`
	PositionLevelID     int    `json:"position_level_id"`
	PositionLevelName   string `json:"position_level_name"`
	DoB                 string `json:"dob" gorm:"column:dob"`
	MaterialStatus      string `json:"material_status" gorm:"column:material_status"`
	ProvinceIDBirth     string `json:"province_id_birth"`
	ProvinceNameBirth   string `json:"province_name_birth"`
	DistrictIDBirth     int    `json:"district_id_birth"`
	DistrictNameBirth   string `json:"district_name_birth"`
	CommunceIDBirth     int    `json:"communce_id_birth"`
	CommunceNameBirth   string `json:"communce_name_birth"`
	VillageIdBirth      int    `json:"village_id_birth"`
	VillageNameBirth    string `json:"village_name_birth"`
	ProfileImage        string `json:"profile_image"`
	ProvinceIDCurrent   int    `json:"province_id_current"`
	ProvinceNameCurrent string `json:"province_name_current"`
	DistrictIDCurrent   int    `json:"distirct_id_current" gorm:"column:distirct_id_current"`
	DistrictNameCurrent string `json:"district_name_current"`
	CommunceIDCurrent   int    `json:"communce_id_current"`
	CommunceNameCurrent string `json:"communce_name_current"`
	VillageIdCurrent    int    `json:"village_id_current"`
	VillageNameCurrent  string `json:"village_name_current"`
	FamilyPhone         string `json:"family_phone"`
	BankName            string `json:"bank_name"`
	BankAccountNumber   string `json:"bank_account_number"`
	QrCodeBankAccount   string `json:"qr_code_bank_account"`
	Note                string `json:"note"`
	ReportoID           int    `json:"report_to" gorm:"column:report_to"`
	ReportoName         string `json:"report_to_name" gorm:"column:report_to_name"`
	WifiName            string `json:"wife_name" gorm:"column:wife_name"`
	HusBanName          string `json:"husban_name" gorm:"column:husban_name"`
	SonNumber           int    `json:"son_number" gorm:"column:son_number"`
	DaughterNumber      int    `json:"daughter_number" gorm:"column:daughter_number"`
}
