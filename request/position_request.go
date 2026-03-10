package request

type PositionRequestCreate struct {
	Name         string `json:"name" binding:"required"`
	DisplayName  string `json:"display_name" binding:"required"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

type PositionRequestUpdate struct {
	Name         *string `json:"name"`
	DisplayName  *string `json:"display_name"`
	DepartmentID *uint   `json:"department_id"`
}
