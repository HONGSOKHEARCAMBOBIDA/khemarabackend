package model

type CurrencyPair struct {
	ID               int  `json:"id" gorm:"primaryKey"`
	BaseCurrencyID   int  `json:"base_currency_id" gorm:"column:base_currency_id"`
	TargetCurrencyID int  `json:"target_currency_id" gorm:"column:target_currency_id"`
	IsActive         bool `json:"is_active" gorm:"column:is_active"`
}
