package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"

	"gorm.io/gorm"
)

type CurrencyService interface {
	GetCurrency() ([]model.Currency, error)
	CreateCurrency(input request.CurrencyRequestCreate) error
	UpdateCurrency(id int, intput request.CurrencyRequestUpdate) error
	ChangeStatusCurrency(id int) error
}

type currencyservice struct {
	db *gorm.DB
}

func NewCurrencyService() CurrencyService {
	return &currencyservice{
		db: config.DB,
	}
}

func (s *currencyservice) GetCurrency() ([]model.Currency, error) {
	var currency []model.Currency
	if err := s.db.Find(&currency).Error; err != nil {
		return nil, err
	}
	return currency, nil
}

func (s *currencyservice) CreateCurrency(input request.CurrencyRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if input.Code == "" {
		return errors.New("code is required")
	}
	if input.Symbol == "" {
		return errors.New("symbol is required")
	}
	if input.Name == "" {
		return errors.New("name is required")
	}
	newcurrency := model.Currency{
		Code:     input.Code,
		Symbol:   input.Symbol,
		Name:     input.Name,
		Isactive: true,
	}
	if err := s.db.Create(&newcurrency).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *currencyservice) UpdateCurrency(id int, intput request.CurrencyRequestUpdate) error {
	updates := map[string]interface{}{}
	if intput.Code != nil {
		updates["code"] = *intput.Code
	}
	if intput.Symbol != nil {
		updates["symbol"] = *intput.Symbol
	}
	if intput.Name != nil {
		updates["name"] = *intput.Name
	}
	if len(updates) == 0 {
		return errors.New("no field to update")
	}
	result := s.db.Model(&model.Currency{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no data chanage")
	}
	return nil
}

func (s *currencyservice) ChangeStatusCurrency(id int) error {
	result := s.db.Model(&model.Currency{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no da")
	}
	return nil
}
