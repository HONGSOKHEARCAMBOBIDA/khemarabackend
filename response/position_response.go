package response

type PositionResponse struct {
	ID                    uint   `json:"id"`
	Name                  string `json:"name"`
	DisplayName           string `json:"display_name" gorm:"column:display_name"`
	DepartmentID          uint   `json:"department_id" gorm:"column:department_id"`
	DepartmentName        string `json:"department_name" gorm:"column:department_name"`
	DepartmentDisplayName string `json:"department_display_name"`
	Isactive              bool   `json:"is_active" gorm:"column:is_active"`
}
