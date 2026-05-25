package response

type LeaveDurationResponse struct {
	ID                 int     `json:"id"`
	DurationValue      float64 `json:"duration_value"`
	DurationUnitID     int     `json:"duration_unit_id"`
	DurationUnitCode   string  `json:"duration_unit_code"`
	DurationUnitNameEn string  `json:"duration_unit_name_en"`
	DurationUnitNameKh string  `json:"duration_unit_name_kh"`
}
