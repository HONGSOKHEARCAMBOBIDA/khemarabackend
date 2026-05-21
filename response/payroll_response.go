package response

type PayrollResponse struct {
	ID                int    `json:"id"`
	EmployeeNameEn    string `json:"employee_name_en"`
	EmployeeNameKh    string `json:"employee_name_kh"`
	EmployeeGender    int    `json:"employee_gender"`
	PositionName      string `json:"position_name"`
	OfficeName        string `json:"office_name"`
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
	BonusTypeName     string `json:"bonus_type_name"`
	BonusAmount       string `json:"bonus_amount"`
	TotalDeduction    string `json:"total_deduction"`
	NetSalary         string `json:"net_salary"`
	CurrencyName      string `json:"currency_name"`
	CurrencyCode      string `json:"currency_code"`
	StatusName        string `json:"status_name"`
	Note              string `json:"note"`
	BranchName        string `json:"branch_name"`
	ShowApprovebutton bool   `json:"show_approve_button"`
}
