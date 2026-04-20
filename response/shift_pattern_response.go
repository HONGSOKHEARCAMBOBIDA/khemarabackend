package response

type ShiftPatternResponse struct {
	ID            int             `json:"id"`
	DayOfWeekID   int             `json:"day_of_week_id"`
	DayOfWeekName string          `json:"day_of_week_name"`
	IsdayOff      bool            `json:"is_dayoff" gorm:"column:is_dayoff"`
	ShiftID       int             `json:"shift_id"`
	ShiftResponse []ShiftResponse `json:"shift_response" gorm:"-"`
}
