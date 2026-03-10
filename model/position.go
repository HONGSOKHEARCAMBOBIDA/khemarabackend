package model

type Position struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	DepartmentID uint   `json:"department_id"`
	Isactive     bool   `json:"is_active" gorm:"column:is_active"`
}
