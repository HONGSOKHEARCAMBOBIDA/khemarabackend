package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type ExchangeRateService interface {
	GetExchangeRate() ([]response.ExchangeRateResponse, error)
	CreateExchangeRate(input request.ExchangeRateRequestCreate) error
	UpdateExchangeRate(id int, input request.ExchangeRateRequestUpdate) error
	ChangeStatusExchangeRate(id int) error
}

type exchangerateservice struct {
	db *gorm.DB
}

func NewExchangeRateService() ExchangeRateService {
	return &exchangerateservice{
		db: config.DB,
	}
}

func (s *exchangerateservice) GetExchangeRate() ([]response.ExchangeRateResponse, error) {
	var exchangerateresponse []response.ExchangeRateResponse
	db := config.DB.Table("exchange_rates").Select(`
	exchange_rates.id AS id,
	b.id AS base_currency_id,
	b.code AS base_currency_code,
	b.symbol AS base_currency_symbol,
	b.name AS base_currency_name,
	b.is_active AS base_currency_is_active,
	t.id AS target_currency_id,
	t.code AS target_currency_code,
	t.symbol AS target_currency_symbol,
	t.name AS target_currency_name,
	t.is_active AS target_currency_is_active,
	exchange_rates.rate AS rate,
	exchange_rates.is_active AS is_active,
	exchange_rates.is_edit AS is_edit,
	exchange_rates.pair_id AS pair_id
	
	`).
		Joins("INNER JOIN currency_pairs c ON c.id = exchange_rates.pair_id").
		Joins("INNER JOIN currencies b ON b.id = c.base_currency_id").
		Joins("INNER JOIN currencies t ON t.id = c.target_currency_id").
		Order("exchange_rates.id desc")
	if err := db.Scan(&exchangerateresponse).Error; err != nil {
		return nil, err
	}
	return exchangerateresponse, nil
}

func (s *exchangerateservice) CreateExchangeRate(input request.ExchangeRateRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if input.PairID == 0 {
		return errors.New("pair id is required")
	}
	if input.Rate == 0 {
		return errors.New("rate is required")
	}
	newexchangerate := model.ExchangeRate{
		PairID:   input.PairID,
		Rate:     input.Rate,
		Isactive: true,
	}
	if err := tx.Create(&newexchangerate).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *exchangerateservice) UpdateExchangeRate(id int, input request.ExchangeRateRequestUpdate) error {
	updates := map[string]interface{}{}
	if input.PairID != nil {
		updates["pair_id"] = *input.PairID
	}
	if input.Rate != nil {
		updates["rate"] = *input.Rate
	}
	updates["is_edit"] = true
	if len(updates) == 0 {
		return errors.New("no field to udpate")
	}
	result := s.db.Model(&model.ExchangeRate{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no data changed")
	}
	return nil
}

func (s *exchangerateservice) ChangeStatusExchangeRate(id int) error {
	result := s.db.Model(&model.ExchangeRate{}).Where("id =?", id).Update("is_active", gorm.Expr("NTO is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no data changed")
	}
	return nil
}
