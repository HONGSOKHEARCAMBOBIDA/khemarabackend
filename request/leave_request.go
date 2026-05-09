package request

type LeaveCreate struct {
	LeaveTypeID    int     `form:"leave_type_id"`
	StartDate      string  `form:"start_date"`
	EndDate        string  `form:"end_date"`
	Description    string  `form:"description"`
	ApproveByID    int     `form:"approve_by"`
	DurationVlaue  float64 `form:"duration_value"`
	DurationUnitID int     `form:"duration_unit_id"`
}
