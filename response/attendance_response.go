package response

type AttendanceResponse struct {
	AttendanceLogID       int                     `json:"attendance_log_id"`
	EmployeeID            int                     `json:"employee_id"`
	EmployeeCode          string                  `json:"employee_code"`
	EmployeeNameEn        string                  `json:"employee_name_en"`
	EmployeeNameKh        string                  `json:"employee_name_kh"`
	Gender                int                     `json:"gender"`
	PositionID            int                     `json:"position_id"`
	PositionDisplayName   string                  `json:"position_display_name"`
	DepartmentID          int                     `json:"department_id"`
	DepartmentDisplayName string                  `json:"department_display_name"`
	OfficeID              int                     `json:"office_id"`
	OfficeName            string                  `json:"office_name"`
	Profile               string                  `json:"profile"`
	AttendanceLogResponse []AttendanceLogResponse `json:"attendancelogresponse" gorm:"-"`
	BranchName            string                  `json:"branch_name" gorm:"column:branch_name"`
}

type AttendanceLogResponseV2 struct {
	ID                         int                          `json:"id"`
	EmployeeCode               string                       `json:"employee_code"`
	EmployeeName               string                       `json:"employee_name"`
	EmployeeNameEn             string                       `json:"employee_name_en"`
	PositionName               string                       `json:"position_name"`
	DepartmentName             string                       `json:"department_name"`
	CheckDate                  string                       `json:"check_date" gorm:"column:check_date"`
	BranchName                 string                       `json:"branch_name"`
	StatusName                 string                       `json:"status_name"`
	AttendanceRecordResponseV2 []AttendanceRecordResponseV2 `json:"attendance_record" gorm:"-"`
}

type AttendanceRecordResponseV2 struct {
	ID               int     `json:"id"`
	AttendanceLogID  int     `json:"attendance_log_id"`
	SessionName      string  `json:"session_name"`
	StartTime        string  `json:"start_time"`
	EndTime          string  `json:"end_time"`
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
