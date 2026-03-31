package model

type EmployeeProfile struct {
	ID                      int    `json:"id"`
	EmployeeID              int    `json:"employee_id"`
	PositionLevelID         int    `json:"position_level_id"`
	DoB                     string `json:"dob"`
	VillageIdOfBirth        int    `json:"village_id_of_birth"`
	MaterialStatus          string `json:"material_status"`
	ProfileImage            string `json:"profile_image"`
	VillageIDCurrentAddress int    `json:"village_id_current_address"`
	FamilyPhone             string `json:"family_phone"`
	BankName                string `json:"bank_name"`
	BankAccountNumber       string `json:"bank_account_number"`
	QrCodeBankAccount       string `json:"qr_code_bank_account"`
	Note                    string `json:"note"`
	ReportoID               int    `json:"report_to"`
	WifeName                string `json:"wife_name"`
	HusBanName              string `json:"husban_name"`
	SonNumber               int    `json:"son_number"`
	DaughterNumber          int    `json:"daughter_number"`
	CreateBy                int    `json:"create_by"`
	UpdateBy                int    `json:"update_by"`
}
