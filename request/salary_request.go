package request

type SalaryRequestUpdate struct {
	BaseSalary    float64 `form:"base_salary"`
	WorkDay       int     `form:"work_day"`
	DailyRate     float64 `form:"daily_rate"`
	EffectiveDate string  `form:"effective_date"`
	CurrencyID    int     `form:"currency_id"`
}

type SalaryRequestCreate struct {
	EmployeeID    int     `form:"employee_id"`
	BaseSalary    float64 `form:"base_salary"`
	WorkDay       int     `form:"work_day"`
	DailyRate     float64 `form:"daily_rate"`
	EffectiveDate string  `form:"effective_date"`
	CurrencyID    int     `form:"currency_id"`
}
