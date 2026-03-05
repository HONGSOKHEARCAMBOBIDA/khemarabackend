package model

type Role struct {
	ID          uint         `gorm:"primarykey" json:"id"`
	Name        string       `json:"name"`
	DisPlayName string       `json:"display_name" gorm:"column:display_name"`
	IsActive    bool         `json:"is_active"`
	Permissions []Permission `gorm:"many2many:role_has_permissions"`
}
