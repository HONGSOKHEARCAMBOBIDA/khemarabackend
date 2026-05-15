package model

type Payroll struct {
	ID             int     `json:"id"`
	SalaryID       int     `json:"salary_id" gorm:"column:salary_id"`
	BranchID       int     `json:"branch_id" gorm:"column:branch_id"`
	PayRollTypeID  int     `json:"payroll_type_id" gorm:"column:payroll_type_id"`
	BasicSalary    float64 `json:"basic_salary" gorm:"column:basic_salary"`
	HalfSalary     float64 `json:"half_salary" gorm:"column:half_salary"`
	Pensionfund    float64 `json:"pension_fund" gorm:"column:pension_fund"`
	TotalWorkDay   int     `json:"total_work_day" gorm:"column:total_work_day"`
	PayrollDate    string  `json:"payroll_date" gorm:"column:payroll_date"`
	LoanDeduction  float64 `json:"loan_deduction" gorm:"column:loan_deduction"`
	Isbonus        bool    `json:"is_bonus" gorm:"column:is_bonus"`
	BonusType      int     `json:"bonus_type" gorm:"column:bonus_type"`
	BonusAmount    float64 `json:"bonus_amount" gorm:"column:bonus_amount"`
	TotalDeduction float64 `json:"total_deduction" gorm:"column:total_deduction"`
	NetSalary      float64 `json:"net_salary" gorm:"column:net_salary"`
	CurrencyID     int     `json:"currency_id" gorm:"column:currency_id"`
	StatusID       int     `json:"status_id" gorm:"column:status_id"`
	SubmittedBy    int     `json:"submitted_by" gorm:"column:submitted_by"`
	SubmittedDate  string  `json:"submitted_date" gorm:"column:submitted_date"`
	Note           string  `json:"note" gorm:"column:note"`
}
