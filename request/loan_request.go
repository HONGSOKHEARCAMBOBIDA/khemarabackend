package request

type LoanRequest struct {
	EmployeeID    int     `json:"employee_id" gorm:"column:employee_id"`
	LoanAmount    float64 `json:"loan_amount" gorm:"column:loan_amount"`
	CurrencyID    int     `json:"currency_id" gorm:"column:currency_id"`
	LoanStartDate string  `json:"loan_start_date" gorm:"column:loan_start_date"`
	LoanEndDate   string  `json:"loan_end_date" gorm:"column:loan_end_date"`
	LoanPurpose   string  `json:"loan_purpose" gorm:"column:loan_purpose"`
	LoanDuration  int     `json:"loan_duration"`
}
