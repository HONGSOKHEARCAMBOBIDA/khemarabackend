package model

type Department struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Isactive    bool   `json:"is_active"`
}
