package request

type EmployeeWorkExperienceRequestUpdate struct {
	CompanyName    string `form:"company_name"`
	PositionTitle  string `form:"position_title"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
	JobDescription string `form:"job_description"`
}

type EmployeeWorkExperienceRequestCreate struct {
	EmployeeID     int    `form:"employee_id"`
	CompanyName    string `form:"company_name"`
	PositionTitle  string `form:"position_title"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
	JobDescription string `form:"job_description"`
}
