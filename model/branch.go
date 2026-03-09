package model

type Branch struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Radius    int    `json:"radius"`
	Isactive  bool   `json:"is_active" gorm:"column:is_active"`
}
