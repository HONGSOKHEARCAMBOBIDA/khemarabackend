package request

type ShiftPatternRequestUpdate struct {
	EmployeeID int `json:"employee_id"`
	ShiftID    int `json:"shift_id"`
}
