package model

type Loan struct {
	ID                   int     `json:"id"`
	Code                 string  `json:"code" gorm:"column:code"`
	BranchID             int     `json:"branch_id" gorm:"column:branch_id"`
	EmployeeID           int     `json:"employee_id" gorm:"column:employee_id"`
	LoanAmount           float64 `json:"loan_amount" gorm:"column:loan_amount"`
	CurrencyID           int     `json:"currency_id" gorm:"column:currency_id"`
	ApproveDate          string  `json:"approve_date" gorm:"column:approve_date"`
	LoanStartDate        string  `json:"loan_start_date" gorm:"column:loan_start_date"`
	LoanEndDate          *string `json:"loan_end_date" gorm:"column:loan_end_date"`
	LoanRateAmount       float64 `json:"loan_rate_amount" gorm:"column:loan_rate_amount"`
	NumberofLoan         int     `json:"number_of_loan" gorm:"column:number_of_loan"`
	ApproveBy            int     `json:"approve_by" gorm:"column:approve_by"`
	LoanPurpose          string  `json:"loan_purpose" gorm:"column:loan_purpose"`
	MonthlyPaymentAmount float64 `json:"monthly_payment_amount" gorm:"column:monthly_payment_amount"`
	Status               int     `json:"status" gorm:"column:status"`
	LoanDuration         int     `json:"loan_duration" gorm:"column:loan_duration"`
}
