package model

type ShiftPattern struct {
	ID          int  `json:"id"`
	EmployeeID  int  `json:"employee_id"`
	DayOfWeekID int  `json:"day_of_week_id"`
	ShiftID     int  `json:"shift_id"`
	Isdayoff    bool `json:"is_day_off" gorm:"column:is_day_off"`
}
