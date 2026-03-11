package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type CurrencyPairService interface {
	CreateCurrencyPair(input request.CurrencyPairRequestCreate) error
	GetCurrencypair() ([]response.CurrencyPairResponse, error)
	UpdateCurrencyPaire(id int, input request.CurrencyPairRequestUpdate) error
	ChangeStatusCurrencyPair(id int) error
}

type currencypairService struct {
	db *gorm.DB
}

func NewCurrencyPairService() CurrencyPairService {
	return &currencypairService{
		db: config.DB,
	}
}
func (cps *currencypairService) CreateCurrencyPair(input request.CurrencyPairRequestCreate) error {
	tx := cps.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	newcurrencypair := model.CurrencyPair{
		BaseCurrencyID:   input.BaseCurrencyID,
		TargetCurrencyID: input.TargetCurrencyID,
		IsActive:         true,
	}
	if err := tx.Create(&newcurrencypair).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (cps *currencypairService) GetCurrencypair() ([]response.CurrencyPairResponse, error) {
	var currencypair []response.CurrencyPairResponse

	db := cps.db.Table("currency_pairs").
		Select(`
			currency_pairs.id AS id,
			base.id AS base_currency_id,
			base.code AS base_currency_code,
			base.symbol AS base_currency_symbol,
			base.name AS base_currency_name,
			base.is_active AS base_currency_is_active,
			target.id AS target_currency_id,
			target.code AS target_currency_code,
			target.symbol AS target_currency_symbol,
			target.name AS target_currency_name,
			target.is_active AS target_currency_is_active,
			currency_pairs.is_active AS is_active
		`).
		Joins("INNER JOIN currencies AS base ON base.id = currency_pairs.base_currency_id").
		Joins("INNER JOIN currencies AS target ON target.id = currency_pairs.target_currency_id").
		Order("currency_pairs.id DESC")

	if err := db.Scan(&currencypair).Error; err != nil {
		return nil, err
	}

	return currencypair, nil
}

func (cps *currencypairService) UpdateCurrencyPaire(id int, input request.CurrencyPairRequestUpdate) error {
	result := cps.db.Model(&model.CurrencyPair{}).Where("id =?", id).Updates(input)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("currency pair not foundn or not update")
	}
	return nil
}

func (cps *currencypairService) ChangeStatusCurrencyPair(id int) error {
	result := cps.db.Model(&model.CurrencyPair{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("currency pair not found or not update")

	}
	return nil
}
