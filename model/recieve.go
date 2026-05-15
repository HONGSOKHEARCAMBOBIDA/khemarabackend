package model

type Recieve struct {
	ID           int     `json:"id"`
	Code         string  `json:"code" gorm:"column:code"`
	BranchID     int     `json:"branch_id" gorm:"column:branch_id"`
	LoanID       int     `json:"loan_id" gorm:"column:loan_id"`
	RecieveDate  string  `json:"receive_date" gorm:"column:receive_date"`
	TotalRecieve float64 `json:"total_receive" gorm:"column:total_receive"`
	CurrencyID   int     `json:"currency_id" gorm:"column:currency_id"`
	Note         string  `json:"note" gorm:"column:note"`
	RecieveBy    int     `json:"receive_by" gorm:"column:receive_by"`
	PayrollID    int     `json:"payroll_id" gorm:"column:payroll_id"`
}
