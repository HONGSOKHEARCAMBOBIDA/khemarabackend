package model

type Permission struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	DisplayName string `json:"display_name" gorm:"column:display_name"`
	Group       string `json:"group_name" gorm:"column:group_name"`
	Short       int    `json:"short_name" gorm:"column:short_name"`
}
