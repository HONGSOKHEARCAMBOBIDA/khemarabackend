package model

type Shift struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	BranchID int    `json:"branch_id"`
	Isactive bool   `json:"is_active" gorm:"column:is_active"`
}
