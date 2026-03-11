package request

type CurrencyPairRequestCreate struct {
	BaseCurrencyID   int `json:"base_currency_id" gorm:"column:base_currency_id"`
	TargetCurrencyID int `json:"target_currency_id" gorm:"column:target_currency_id"`
}

type CurrencyPairRequestUpdate struct {
	BaseCurrencyID   int `json:"base_currency_id" gorm:"column:base_currency_id"`
	TargetCurrencyID int `json:"target_currency_id" gorm:"column:target_currency_id"`
}
