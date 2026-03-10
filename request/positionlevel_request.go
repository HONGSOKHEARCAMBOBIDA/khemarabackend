package request

type PositionLevelRequestCreate struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type PositionLevelRequestUpdate struct {
	Name        *string `json:"name"`
	DisplayName *string `json:"display_name"`
}
