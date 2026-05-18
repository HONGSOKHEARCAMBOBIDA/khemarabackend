package response

type PayrollDrafResponse struct {
	EmployeeID     int     `json:"employee_id"`
	EmployeeName   string  `json:"employee_name"`
	BranchID       int     `json:"branch_id"`
	BranchName     string  `json:"branch_name"`
	SalaryID       int     `json:"salary_id"`
	BaseSalary     float32 `json:"base_salary"`
	DailyRate      float32 `json:"daily_rate"`
	HalfSalary     float32 `json:"half_salary" gorm:"column:half_salary"`
	Pensionfund    float32 `json:"pensionfund" gorm:"column:pensionfund"`
	TotalWorkDay   int     `json:"total_work_day" gorm:"column:total_work_day"`
	LoanDeduction  float32 `json:"loan_deduction" gorm:"column:loan_deduction"`
	TotalDeduction float32 `json:"total_deduction" gorm:"column:total_deduction"`
	NetSalary      float32 `json:"net_salary" gorm:"column:net_salary"`
	Comment        string  `json:"comment"`
	LoanID         int     `json:"loan_id"`
	CurrenySymbol  string  `json:"currency_symbol" gorm:"column:currency_symbol"`
}
