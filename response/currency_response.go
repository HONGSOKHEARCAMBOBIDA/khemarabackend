package response

type CurrencyResponse struct {
	ID       int    `json:"id"`
	Code     string `json:"code" gorm:"column:code"`
	Symbol   string `json:"symbol" gorm:"column:symbol"`
	Name     string `json:"name" gorm:"column:name"`
	Isactive bool   `json:"is_active" gorm:"column:is_active"`
}
