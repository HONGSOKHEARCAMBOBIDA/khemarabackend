package request

type ExchangeRateRequestCreate struct {
	PairID int     `json:"pair_id" gorm:"column:pair_id"`
	Rate   float64 `json:"rate" gorm:"column:rate"`
}

type ExchangeRateRequestUpdate struct {
	PairID *int     `json:"pair_id" gorm:"column:pair_id"`
	Rate   *float64 `json:"rate" gorm:"column:rate"`
}
