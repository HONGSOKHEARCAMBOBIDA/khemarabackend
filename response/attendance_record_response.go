package response

type AttendanceRecordResponse struct {
	ID               int     `json:"id"`
	ShiftSessionID   int     `json:"shift_session_id"`
	StartTime        string  `json:"start_time"`
	EndTime          string  `json:"end_time"`
	ShiftOrder       int     `json:"shift_order"`
	CheckTime        string  `json:"check_time"`
	CheckInEarly     int     `json:"check_in_early" gorm:"column:check_in_early"`
	CheckInOnTime    int     `json:"check_in_on_time" gorm:"column:check_in_on_time"`
	IsLate           int     `json:"is_late"`
	IsLeftEarly      int     `json:"is_left_early"`
	CheckOutOnTime   int     `json:"check_out_on_time" gorm:"column:check_out_on_time"`
	CheckOutOverTime int     `json:"check_out_overtime" gorm:"column:check_out_overtime"`
	Latitude         float64 `json:"latitude"`
	Logitude         float64 `json:"longitude" gorm:"column:longitude"`
	Note             string  `json:"note"`
	Iszoone          bool    `json:"iszonecheckin" gorm:"column:iszonecheckin"`
	Type             string  `json:"type"`
}
