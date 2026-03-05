package request

type RoleRequestCreate struct {
	Name        string `json:"name"`
	DisPlayName string `json:"display_name"`
}

type RoleRequestUpdate struct {
	Name        *string `json:"name"`
	DisPlayName *string `json:"display_name"`
}
