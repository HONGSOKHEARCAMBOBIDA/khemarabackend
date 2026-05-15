package model

type RecieveDetail struct {
	ID         int     `json:"id"`
	LoanID     int     `json:"loan_id" gorm:"column:loan_id"`
	RecieveID  int     `json:"receive_id" gorm:"column:receive_id"`
	ScheduleID int     `json:"schedule_id" gorm:"column:schedule_id"`
	Principle  float64 `json:"principal" gorm:"column:principal"`
	Rate       float64 `json:"rate" gorm:"column:rate"`
	Income     float64 `json:"income" gorm:"column:income"`
}
