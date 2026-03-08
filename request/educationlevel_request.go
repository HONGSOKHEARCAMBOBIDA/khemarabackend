package request

type EducationLevelRequestCreate struct {
	Name string `json:"name"`
}

type EducationLevelRequestUpdate struct {
	Name *string `json:"name"`
}
