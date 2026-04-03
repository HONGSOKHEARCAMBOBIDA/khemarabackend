package request

type ShiftSessionRequestCreate struct {
	SessionName string `json:"session_name"`
	ShiftID     int    `json:"shift_id"`
	ShiftOrder  int    `json:"shift_order"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type ShiftSessionRequestUpdate struct {
	SessionName *string `json:"session_name"`
	ShiftID     *int    `json:"shift_id"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
}
