package response

type ShiftResponse struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	Isactive             bool                   `json:"is_active" gorm:"column:is_active"`
	BranchID             int                    `json:"branch_id"`
	BranchName           string                 `json:"branch_name"`
	ShiftSessionResponse []ShiftSessionResponse `json:"shiftsessionresponse" gorm:"-"`
}
