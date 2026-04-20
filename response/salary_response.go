package response

type SalaryResponse struct {
	ID             int    `json:"id"`
	BaseSalary     string `json:"base_salary"`
	WorkDay        int    `json:"work_day"`
	DailyRate      string `json:"daily_rate"`
	EffectiveDate  string `json:"effective_date"`
	ExpireDate     string `json:"expire_date"`
	CurrencyID     int    `json:"currency_id"`
	CurrencyCode   string `json:"currency_code"`
	CurrencySymbol string `json:"currency_symbol"`
	CurrencyName   string `json:"currency_name"`
	Isactive       bool   `json:"is_active" gorm:"column:is_active"`
}
