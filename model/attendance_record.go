package model

import "time"

type AttendanceRecord struct {
	ID              int       `json:"id"`
	AttendanceLogID int       `json:"attendance_log_id"`
	ShiftSessionID  int       `json:"shift_session_id"`
	CheckTime       time.Time `json:"check_time"`
	IsLate          int       `json:"is_late"`
	IsLeftEarly     *int      `json:"is_left_early"`
	Latitude        float64   `json:"latitude"`
	Logitude        float64   `json:"longitude" gorm:"column:longitude"`
	Note            string    `json:"note"`
	Iszoone         bool      `json:"iszonecheckin" gorm:"column:iszonecheckin"`
}
