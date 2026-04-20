package response

type ShiftSessionResponse struct {
	ID          int    `json:"id"`
	SessionName string `json:"session_name"`
	// ShiftID     int    `json:"shift_id"`
	// ShiftName   string `json:"shift_name"`
	// BranchID    int    `json:"branch_id"`
	// BranchName  string `json:"branch_name"`
	ShiftOrder int    `json:"shift_order"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Isactive   bool   `json:"is_active" gorm:"column:is_active"`
}
