package model

type Pensionfund struct {
	ID         int     `json:"id"`
	EmployeeID int     `json:"employee_id" gorm:"column:employee_id"`
	BranchID   int     `json:"branch_id" gorm:"column:branch_id"`
	Amount     float64 `json:"amount" gorm:"column:amount"`
	CurrencyID int     `json:"currency_id" gorm:"column:currency_id"`
	Date       string  `json:"date" gorm:"column:date"`
	PayrollID  int     `json:"payroll_id" gorm:"column:payroll_id"`
}
