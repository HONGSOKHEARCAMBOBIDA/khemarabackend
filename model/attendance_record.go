package model

import "time"

type AttendanceRecord struct {
	ID               int       `json:"id"`
	AttendanceLogID  int       `json:"attendance_log_id"`
	ShiftSessionID   int       `json:"shift_session_id"`
	CheckTime        time.Time `json:"check_time"`
	CheckInEarly     *int      `json:"check_in_early" gorm:"column:check_in_early"`
	CheckInOnTime    *int      `json:"check_in_on_time" gorm:"column:check_in_on_time"`
	IsLate           *int      `json:"is_late"`
	IsLeftEarly      *int      `json:"is_left_early"`
	CheckOutOnTime   *int      `json:"check_out_on_time" gorm:"column:check_out_on_time"`
	CheckOutOverTime *int      `json:"check_out_overtime" gorm:"column:check_out_overtime"`
	Latitude         float64   `json:"latitude"`
	Logitude         float64   `json:"longitude" gorm:"column:longitude"`
	Note             string    `json:"note"`
	Iszoone          bool      `json:"iszonecheckin" gorm:"column:iszonecheckin"`
	Type             string    `json:"type"`
}

type AttendanceRecordRes struct {
	ID              int     `json:"id"`
	AttendanceLogID int     `json:"attendance_log_id"`
	ShiftSessionID  int     `json:"shift_session_id"`
	CheckTime       string  `json:"check_time"`
	IsLate          *int    `json:"is_late"`
	IsLeftEarly     *int    `json:"is_left_early"`
	Latitude        float64 `json:"latitude"`
	Logitude        float64 `json:"longitude" gorm:"column:longitude"`
	Note            string  `json:"note"`
	Iszoone         bool    `json:"iszonecheckin" gorm:"column:iszonecheckin"`
}

func (AttendanceRecordRes) TableName() string {
	return "attendance_records"
}
