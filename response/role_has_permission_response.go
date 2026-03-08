package response

type PermissionWithAssignedRole struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	GroupName   string `json:"group_name" gorm:"column:group_name"`
	Shortname   int    `json:"short_name" gorm:"column:short_name"`
	Assigned    bool   `json:"assigned"`
	RoleID      int    `json:"role_id"`
}
