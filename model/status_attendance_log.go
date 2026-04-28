package model

type StatusAttendanceLog struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Isactive bool   `json:"is_active" gorm:"column:is_active"`
}
