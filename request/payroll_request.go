package request

type PayrollRequestCreate struct {
	EmployeeID     int     `form:"employee_id"`
	SalaryID       int     `form:"salary_id" gorm:"column:salary_id"`
	BranchID       int     `form:"branch_id"`
	PayRollTypeID  int     `form:"payroll_type_id" gorm:"column:payroll_type_id"`
	BasicSalary    float64 `form:"basic_salary" gorm:"column:basic_salary"`
	HalfSalary     float64 `form:"half_salary" gorm:"column:half_salary"`
	Pensionfund    float64 `form:"pension_fund" gorm:"column:pension_fund"`
	TotalWorkDay   int     `form:"total_work_day" gorm:"column:total_work_day"`
	PayrollDate    string  `form:"payroll_date" gorm:"column:payroll_date"`
	LoanDeduction  float64 `form:"loan_deduction" gorm:"column:loan_deduction"`
	Isbonus        bool    `form:"is_bonus" gorm:"column:is_bonus"`
	BonusType      int     `form:"bonus_type" gorm:"column:bonus_type"`
	BonusAmount    float64 `form:"bonus_amount" gorm:"column:bonus_amount"`
	TotalDeduction float64 `form:"total_deduction" gorm:"column:total_deduction"`
	NetSalary      float64 `form:"net_salary" gorm:"column:net_salary"`
	CurrencyID     int     `form:"currency_id" gorm:"column:currency_id"`
	Note           string  `form:"note" gorm:"column:note"`
	Comment        string  `form:"comment"`
	LoanID         int     `form:"loan_id"`
}
