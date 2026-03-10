package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"

	"gorm.io/gorm"
)

type PositionLevelService interface {
	GetPositionLevel() ([]model.PositionLevel, error)
	CreatePositionLevel(input request.PositionLevelRequestCreate) error
	UpdatePositionLevel(id int, input request.PositionLevelRequestUpdate) error
	ChangeStatusPositionLevel(id int) error
}

type positionlevelservice struct {
	db *gorm.DB
}

func NewPositionLevelService() PositionLevelService {
	return &positionlevelservice{
		db: config.DB,
	}
}

func (s *positionlevelservice) GetPositionLevel() ([]model.PositionLevel, error) {
	var positionlevel []model.PositionLevel
	if err := s.db.Find(&positionlevel).Error; err != nil {
		return nil, err
	}
	return positionlevel, nil
}

func (s *positionlevelservice) CreatePositionLevel(input request.PositionLevelRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.DisplayName == "" {
		return errors.New("display name is required")
	}
	newpositionlevel := model.PositionLevel{
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Isactive:    true,
	}
	if err := s.db.Create(&newpositionlevel).Error; err != nil {
		return err
	}
	return nil
}

func (s *positionlevelservice) UpdatePositionLevel(id int, input request.PositionLevelRequestUpdate) error {
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.DisplayName != nil {
		updates["display_name"] = *input.DisplayName
	}
	if len(updates) == 0 {
		return errors.New("no field to update")
	}
	result := s.db.Model(&model.PositionLevel{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no data changed")
	}
	return nil
}

func (s *positionlevelservice) ChangeStatusPositionLevel(id int) error {
	result := s.db.Model(&model.PositionLevel{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("not data changed")
	}
	return nil
}
