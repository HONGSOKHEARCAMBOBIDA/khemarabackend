package model

type ShiftSession struct {
	ID          int    `json:"id"`
	SessionName string `json:"session_name"`
	ShiftID     int    `json:"shift_id"`
	ShiftOrder  int    `json:"shift_order"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Isactive    bool   `json:"is_active"`
}
