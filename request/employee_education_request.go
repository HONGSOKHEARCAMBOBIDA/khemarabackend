package request

type EmployeeEducationRequestUpdate struct {
	EducationLevelID int    `form:"education_level_id"`
	MajorField       string `form:"major_field_of_study" gorm:"column:major_field_of_study"`
	StartDate        string `form:"start_date"`
	EndDate          string `form:"end_date"`
	Note             string `form:"note"`
}

type EmployeeEducationRequestCreate struct {
	EmployeeID       int    `form:"employee_id"`
	EducationLevelID int    `form:"education_level_id"`
	MajorField       string `form:"major_field_of_study" gorm:"column:major_field_of_study"`
	StartDate        string `form:"start_date"`
	EndDate          string `form:"end_date"`
	Note             string `form:"note"`
}
