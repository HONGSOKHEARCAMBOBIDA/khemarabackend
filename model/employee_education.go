package model

type EmployeeEducation struct {
	ID               int    `json:"id"`
	EmployeeID       int    `json:"employee_id"`
	EducationLevelID int    `json:"education_level_id"`
	MajorField       string `json:"major_field_of_study"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	Note             string `json:"note"`
	Image            string `json:"image"`
	CreateBy         int    `json:"create_by"`
	UpdateBy         int    `json:"update_by"`
}
