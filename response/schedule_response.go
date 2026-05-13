package response

type ScheduleResponse struct {
	ID              int    `json:"schedule_id"`
	PaymentDate     string `json:"payment_date"`
	PaidDate        string `json:"paid_date"`
	PrincipleAmount string `json:"principle_amount"`
	RateAmount      string `json:"rate_amount"`
	IncomeAmount    string `json:"income_amount"`
	PrinciplePaid   string `json:"principle_paid"`
	RatePaid        string `json:"rate_paid"`
	IncomePaid      string `json:"income_paid"`
	Status          int    `json:"staus"`
}
