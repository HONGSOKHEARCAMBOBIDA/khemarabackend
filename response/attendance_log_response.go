package response

type AttendanceLogResponse struct {
	ID                       int                        `json:"id"`
	CheckDate                string                     `json:"check_date"`
	BranchID                 int                        `json:"branch_id"`
	BranchName               string                     `json:"branch_name"`
	StatusAttendanceLogID    int                        `json:"status_attendance_log_id"`
	StatusAttendanceLogName  string                     `json:"status_attendance_log_name"`
	AttendanceRecordResponse []AttendanceRecordResponse `json:"attendancerecordresponse" gorm:"-"`
}
