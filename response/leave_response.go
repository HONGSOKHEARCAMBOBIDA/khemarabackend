package response

type LeaveResponse struct {
	ID                   int     `json:"id"`
	EmployeeID           int     `json:"employee_id"`
	EmployeeNameEn       string  `json:"employee_name_en"`
	EmployeeNameKh       string  `json:"employee_name_kh"`
	EmployeeGender       int     `json:"employee_gender"`
	PositionID           int     `json:"position_id"`
	PositionName         string  `json:"position_name"`
	OfficeID             int     `json:"office_id"`
	OfficeName           string  `json:"office_name"`
	LeaveTypeID          int     `json:"leave_type_id"`
	LeaveTypeName        string  `json:"leave_type_name"`
	DeductTypeID         int     `json:"deduct_type_id"`
	DeductTypeCode       string  `json:"deduct_type_code"`
	DeductTypeName       string  `json:"deduct_type_name"`
	StartDate            string  `json:"start_date"`
	EndDate              string  `json:"end_date"`
	Desscription         string  `json:"description"`
	StatusLeaveID        int     `json:"status_leave_id"`
	StatusLeaveName      string  `json:"status_leave_name"`
	ApproveByID          int     `json:"approve_by_id"`
	ApproveByName        string  `json:"approve_by_name"`
	BranchID             int     `json:"branch_id"`
	BranchName           string  `json:"branch_name"`
	LeaveDurationID      int     `json:"leave_duration_id"`
	DurationValue        float64 `json:"duration_value"`
	DurationUnitID       int     `json:"duration_unit_id"`
	DurationUnitCode     string  `json:"duration_unit_code"`
	DurationUnitNameEn   string  `json:"duration_unit_name_en"`
	DurationUnitNameKh   string  `json:"duration_unit_name_kh"`
	DurationUnitToMinute int     `json:"duration_unit_tominute"`
}
