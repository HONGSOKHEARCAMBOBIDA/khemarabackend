package response

type EmployeeResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Isactive bool   `json:"is_active"`
}
