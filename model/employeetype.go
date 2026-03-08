package model

type EmployeeType struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Isactive bool   `json:"is_active" gorm:"column:is_active"`
}
