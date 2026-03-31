package model

type EmployeeWorkExperience struct {
	ID             int    `json:"id"`
	EmployeeID     int    `json:"employee_id"`
	CompanyName    string `json:"company_name"`
	PositionTitle  string `json:"position_title"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	JobDescription string `json:"job_description"`
	CreateBy       int    `json:"create_by"`
	UpdateBy       int    `json:"update_by"`
}
