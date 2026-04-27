package model

import "time"

type Salary struct {
	ID            int        `json:"id"`
	EmployeeID    int        `json:"employee_id"`
	BaseSalary    float64    `json:"base_salary"`
	WorkDay       int        `json:"work_day"`
	DailyRate     float64    `json:"daily_rate"`
	EffectiveDate string     `json:"effective_date"`
	ExpireDate    *time.Time `json:"expire_date"`
	CurrencyID    int        `json:"currency_id"`
	Isactive      bool       `json:"is_active" gorm:"column:is_active"`
}
