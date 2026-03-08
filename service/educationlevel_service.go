package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"

	"gorm.io/gorm"
)

type EducationLevelService interface {
	GetEducationLevel() ([]model.EducationLevel, error)
	CreateEducationLevel(input request.EducationLevelRequestCreate) error
	UpdateEducationLevel(id int, input request.EducationLevelRequestUpdate) error
	ChangeStatusEducationLevel(id int) error
}

type educationLevelService struct {
	db *gorm.DB
}

func NewEducationLevelService() EducationLevelService {
	return &educationLevelService{
		db: config.DB,
	}
}

func (s *educationLevelService) GetEducationLevel() ([]model.EducationLevel, error) {
	var educationlevel []model.EducationLevel
	if err := s.db.Find(&educationlevel).Error; err != nil {
		return nil, err
	}
	return educationlevel, nil
}

func (s *educationLevelService) CreateEducationLevel(input request.EducationLevelRequestCreate) error {
	if input.Name == "" {
		return errors.New("name is required")
	}
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	neweducation := model.EducationLevel{
		Name:     input.Name,
		Isactive: true,
	}
	if err := tx.Create(&neweducation).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *educationLevelService) UpdateEducationLevel(id int, input request.EducationLevelRequestUpdate) error {

	update := map[string]interface{}{}

	if input.Name != nil {
		update["name"] = *input.Name
	}

	if len(update) == 0 {
		return errors.New("no field to update")
	}

	result := s.db.Model(&model.EducationLevel{}).Where("id =?", id).Updates(update)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("education level not found")
	}

	return nil

}

func (s *educationLevelService) ChangeStatusEducationLevel(id int) error {
	result := s.db.Model(&model.EducationLevel{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("education level not found")
	}
	return nil
}
