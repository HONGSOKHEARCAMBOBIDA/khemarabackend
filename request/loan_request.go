package request

type LoanRequest struct {
	EmployeeID    int     `form:"employee_id" gorm:"column:employee_id"`
	LoanAmount    float64 `form:"loan_amount" gorm:"column:loan_amount"`
	CurrencyID    int     `form:"currency_id" gorm:"column:currency_id"`
	LoanStartDate string  `form:"loan_start_date" gorm:"column:loan_start_date"`
	LoanEndDate   string  `form:"loan_end_date" gorm:"column:loan_end_date"`
	LoanPurpose   string  `form:"loan_purpose" gorm:"column:loan_purpose"`
	LoanDuration  int     `form:"loan_duration"`
}
