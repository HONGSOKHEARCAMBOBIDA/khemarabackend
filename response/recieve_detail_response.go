package response

type RecieveDetailResponse struct {
	ID          int     `json:"recieve_detail_id" gorm:"column:recieve_detail_id"`
	RecieveID   int     `json:"receive_id" gorm:"column:receive_id"`
	Principle   float64 `json:"principal" gorm:"column:principal"`
	Rate        float64 `json:"rate" gorm:"column:rate"`
	Income      float64 `json:"income" gorm:"column:income"`
	PayrollDate string  `json:"payroll_date" gorm:"column:payroll_date"`
	PayrollType string  `json:"payroll_type" gorm:"column:payroll_type"`
}
