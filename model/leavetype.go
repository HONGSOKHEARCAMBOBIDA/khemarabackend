package model

type LeaveType struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Deduct_Type_Id   int     `json:"deduct_type_id"`
	Deduct_Type_Code string  `json:"deduct_type_code"`
	Deduct_Type_Name string  `json:"deduct_type_name"`
	CurrencyID       int     `json:"currency_id"`
	CurrencyName     string  `json:"currency_name"`
	Amount           float64 `json:"amount"`
	Description      string  `json:"description"`
	Isactive         bool    `json:"is_active" gorm:"column:is_active"`
}
