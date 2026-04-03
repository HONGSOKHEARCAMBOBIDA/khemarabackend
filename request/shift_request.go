package request

type ShiftRequestCreate struct {
	Name     string `json:"name"`
	BranchID int    `json:"branch_id"`
}

type ShiftRequestUpdate struct {
	Name     *string `json:"name"`
	BranchID *int    `json:"branch_id"`
}
