package model

type LeaveDuration struct {
	ID             int     `json:"id"`
	LeaveID        int     `json:"leave_id" gorm:"column:leave_id"`
	DurationVlaue  float64 `json:"duration_value" gorm:"column:duration_value"`
	DurationUnitID int     `json:"duration_unit_id" gorm:"column:duration_unit_id"`
	StartTime      *string `json:"start_time" gorm:"column:start_time"`
	EndTime        *string `json:"end_time" gorm:"column:end_time"`
	Note           *string `json:"note" gorm:"column:note"`
}
