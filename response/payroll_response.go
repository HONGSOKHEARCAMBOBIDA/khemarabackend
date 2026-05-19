package response

type PayrollResponse struct {
	ID                int    `json:"id"`
	SalaryID          int    `json:"salary_id"`
	EmployeeID        int    `json:"employee_id"`
	EmployeeNameEn    string `json:"employee_name_en"`
	EmployeeNameKh    string `json:"employee_name_kh"`
	EmployeeGender    int    `json:"employee_gender"`
	PositionID        int    `json:"position_id"`
	PositionName      string `json:"position_name"`
	OfficeID          int    `json:"office_id"`
	OfficeName        string `json:"office_name"`
	ProfileImage      string `json:"profile_image"`
	BankName          string `json:"bank_name"`
	BankAccountNumber string `json:"bank_account_number"`
	QrCodeBankAccount string `json:"qr_code_bank_account"`
	BasicSalary       string `json:"basic_salary"`
	HalfSalary        string `json:"half_salary"`
	Pensionfund       string `json:"pensionfund"`
	TotalWorkDay      int    `json:"total_work_day"`
	PayrollDate       string `json:"payroll_date"`
	LoanDeduction     string `json:"loan_deduction"`
	Isbonus           bool   `json:"is_bonus"`
	BonusType         int    `json:"bonus_type_id"`
	BonusTypeName     string `json:"bonus_type_name"`
	BonusAmount       string `json:"bonus_amount"`
	TotalDeduction    string `json:"total_deduction"`
	NetSalary         string `json:"net_salary"`
	CurrencyID        int    `json:"currency_id"`
	CurrencyName      string `json:"currency_name"`
	CurrencyCode      string `json:"currency_code"`
	StatusID          int    `json:"status_id"`
	StatusName        string `json:"status_name"`
	Note              string `json:"note"`
}
