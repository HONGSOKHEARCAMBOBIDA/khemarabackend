package request

type CurrencyRequestCreate struct {
	Code   string `json:"code" gorm:"column:code"`
	Symbol string `json:"symbol" gorm:"column:symbol"`
	Name   string `json:"name" gorm:"column:name"`
}

type CurrencyRequestUpdate struct {
	Code   *string `json:"code" gorm:"column:code"`
	Symbol *string `json:"symbol" gorm:"column:symbol"`
	Name   *string `json:"name" gorm:"column:name"`
}
