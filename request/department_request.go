package request

type DepartmentRequestCreate struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type DepartmentRequestUpdate struct {
	Name        *string `json:"name"`
	DisplayName *string `json:"display_name"`
}
