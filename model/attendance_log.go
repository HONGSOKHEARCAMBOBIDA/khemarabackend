package model

type AttendanceLog struct {
	ID                    int    `json:"id"`
	EmployeeID            int    `json:"employee_id"`
	CheckDate             string `json:"check_date"`
	Note                  string `json:"note"`
	BranchID              int    `json:"branch_id"`
	StatusAttendanceLogID int    `json:"status_attendance_log_id"`
	ShiftSessionOrder     int    `json:"shift_session_order"`
}
