package model

type Schedule struct {
	ID              int      `json:"id"`
	LoanID          int      `json:"loan_id" gorm:"column:loan_id"`
	PaymentDate     string   `json:"payment_date" gorm:"column:payment_date"`
	PaidDate        *string  `json:"paid_date" gorm:"column:paid_date"`
	PrincipleAmount float64  `json:"principle_amount" gorm:"column:principle_amount"`
	RateAmount      float64  `json:"rate_amount" gorm:"column:rate_amount"`
	IncomeAmount    float64  `json:"income_amount" gorm:"column:income_amount"`
	PrinciplePaid   *float64 `json:"principle_paid" gorm:"column:principle_paid"`
	RatePaid        *float64 `json:"rate_paid" gorm:"column:rate_paid"`
	IncomePaid      *float64 `json:"income_paid" gorm:"column:income_paid"`
	ScheduleNumber  int      `json:"schedule_number" gorm:"column:schedule_number"`
	Status          int      `json:"status" gorm:"column:status"`
}
