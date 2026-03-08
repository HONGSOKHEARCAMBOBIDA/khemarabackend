package request

type EmployeeTypeRequestCreate struct {
	Name string `json:"name"`
}

type EmployeeTypeRequestUpdate struct {
	Name *string `json:"name"`
}
