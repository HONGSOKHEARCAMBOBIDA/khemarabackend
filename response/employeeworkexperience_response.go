package response

type EmployeeWorkExperienceResponse struct {
	ID             int    `json:"id"`
	CompanyName    string `json:"company_name"`
	PositionTitle  string `json:"position_title"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	JobDescription string `json:"job_description"`
}
