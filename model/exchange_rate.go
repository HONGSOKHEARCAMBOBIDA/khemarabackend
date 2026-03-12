package model

type ExchangeRate struct {
	ID       int     `json:"id" gorm:"primarykey"`
	PairID   int     `json:"pair_id" gorm:"column:pair_id"`
	Rate     float64 `json:"rate" gorm:"column:rate"`
	Isactive bool    `json:"is_active" gorm:"column:is_active"`
	IsEdit   bool    `json:"is_edit" gorm:"column:is_edit"`
}
