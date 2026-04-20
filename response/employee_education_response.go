package response

type EmployeeEducationRespons struct {
	ID                 int    `json:"id"`
	EducationLevelID   int    `json:"education_level_id"`
	EducationLevelName string `json:"education_level_name"`
	MajorField         string `json:"major_field_of_study" gorm:"column:major_field_of_study"`
	StartDate          string `json:"start_date"`
	EndDate            string `json:"end_date"`
	Note               string `json:"note"`
	Image              string `json:"image"`
	CreateBy           int    `json:"create_by"`
	CreateByName       string `json:"create_by_name"`
	UpdateBy           int    `json:"update_by_name"`
}
