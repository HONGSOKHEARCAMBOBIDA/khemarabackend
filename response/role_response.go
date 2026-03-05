package response

type RoleRespons struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisPlayName string `json:"display_name"`
	IsActive    bool   `json:"is_active"`
}
