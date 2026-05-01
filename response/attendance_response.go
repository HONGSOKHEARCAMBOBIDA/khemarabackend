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
}
