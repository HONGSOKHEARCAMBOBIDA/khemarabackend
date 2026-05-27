package service

import (
	"errors"
	"fmt"

	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type LeaveService interface {
	CreateLeave(id int, input request.LeaveCreate) error
	UpdateLeave(id int, input request.LeaveUpdate) error
	GetLeave(id int, filters map[string]string, pagination request.Pagination) ([]response.LeaveResponse, *model.PaginationMetadata, error)
	ApproveLeave(id int, input request.LeaveApproveRequest, userID int) error
}

type leaveservice struct {
	db *gorm.DB
}

func NewLeaveService() LeaveService {
	return &leaveservice{db: config.DB}
}

func (s *leaveservice) CreateLeave(id int, input request.LeaveCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user model.User
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with id %d not found", id)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	var emplog model.Employee
	if err := tx.First(&emplog, user.EmployeeID).Error; err != nil {
		tx.Rollback()
		return err
	}

	newLeave := model.Leave{
		EmployeeID:    user.EmployeeID,
		LeaveTypeID:   input.LeaveTypeID,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		BackDate:      input.BackDate,
		Description:   input.Description,
		StatusLeaveID: 1,
		ApproveByID:   input.ApproveByID,
		BranchID:      user.BranchID,
		CreateBy:      id,
	}

	if err := tx.Create(&newLeave).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create leave: %w", err)
	}

	leaveDurations := make([]model.LeaveDuration, 0, len(input.DurationVlaue))
	for i := range input.DurationVlaue {
		leaveDurations = append(leaveDurations, model.LeaveDuration{
			LeaveID:        newLeave.ID,
			DurationVlaue:  input.DurationVlaue[i],
			DurationUnitID: input.DurationUnitID[i],
			StartTime:      nil,
			EndTime:        nil,
			Note:           nil,
		})
	}

	if err := tx.Create(&leaveDurations).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create leave duration: %w", err)
	}

	var employyeeapprove model.Employee
	if err := tx.First(&employyeeapprove, input.ApproveByID).Error; err != nil {
		tx.Rollback()
		return err
	}

	var durationText string

	for i := range input.DurationVlaue {
		var durationunit model.LeaveDurationUnit
		if err := tx.First(&durationunit, input.DurationUnitID[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
		durationText += fmt.Sprintf(
			" %.0f %s",
			input.DurationVlaue[i],
			durationunit.NameKh,
		)
	}

	message := fmt.Sprintf(
		"សូមជម្រាបសួរលោកគ្រូ!\n"+
			"មានបុគ្គលិកឈ្មោះៈ %s\n"+
			"សុំអនុញ្ញាតច្បាប់ឈប់សម្រាក %s\n"+
			"ចាប់ពីថ្ងៃទី %s រហូតដល់ថ្ងៃទី %s នឹងចូលបម្រើការងារវិញនៅថ្ងៃទី %s\n"+
			"*មូលហេតុៈ %s\n"+
			"សូមអរគុណ",
		emplog.NameKh,
		durationText,
		input.StartDate,
		input.EndDate,
		input.BackDate,
		input.Description,
	)

	go helper.SendTelegramMessage(message, employyeeapprove.TelegramChatID)

	return tx.Commit().Error
}

func (s *leaveservice) UpdateLeave(id int, input request.LeaveUpdate) error {

	if len(input.DurationVlaue) != len(input.DurationUnitID) {
		return fmt.Errorf("duration_value and duration_unit_id must have the same length")
	}

	if len(input.DurationVlaue) == 0 {
		return fmt.Errorf("at least one duration is required")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var leave model.Leave
	if err := tx.Where("id = ? AND status_leave_id = ?", id, 1).First(&leave).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("leave not found or cannot be updated")
		}
		return fmt.Errorf("failed to fetch leave: %w", err)
	}

	if err := tx.Model(&leave).Updates(map[string]interface{}{
		"leave_type_id": input.LeaveTypeID,
		"start_date":    input.StartDate,
		"end_date":      input.EndDate,
		"back_date":     input.BackDate,
		"description":   input.Description,
		"approve_by_id": input.ApproveByID,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update leave: %w", err)
	}

	if err := tx.Where("leave_id = ?", id).Delete(&model.LeaveDuration{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete old leave durations: %w", err)
	}

	leaveDurations := make([]model.LeaveDuration, 0, len(input.DurationVlaue))
	for i := range input.DurationVlaue {
		leaveDurations = append(leaveDurations, model.LeaveDuration{
			LeaveID:        id,
			DurationVlaue:  input.DurationVlaue[i],
			DurationUnitID: input.DurationUnitID[i],
			StartTime:      nil,
			EndTime:        nil,
			Note:           nil,
		})
	}

	if err := tx.Create(&leaveDurations).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create new leave durations: %w", err)
	}

	return tx.Commit().Error
}

func (s *leaveservice) GetLeave(id int, filters map[string]string, pagination request.Pagination) ([]response.LeaveResponse, *model.PaginationMetadata, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	var role model.Role
	if err := s.db.First(&role, user.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("role with id %d not found", user.RoleID)
		}
		return nil, nil, fmt.Errorf("failed to fetch role: %w", err)
	}

	var leaves []response.LeaveResponse
	var totalCount int64
	offset := (pagination.Page - 1) * pagination.PageSize

	query := s.db.Table("leaves l").
		Select(`
			l.id                    AS id,
			u.contact 				AS employee_phone,
			e.id                    AS employee_id,
			e.code					AS employee_code,
			e.name_en               AS employee_name_en,
			e.name_kh               AS employee_name_kh,
			e.gender                AS employee_gender,
			p.id                    AS position_id,
			p.display_name          AS position_name,
			o.id                    AS office_id,
			o.name                  AS office_name,
			lt.id                   AS leave_type_id,
			lt.name                 AS leave_type_name,
			ddt.id                  AS deduct_type_id,
			ddt.code                AS deduct_type_code,
			ddt.name                AS deduct_type_name,
			l.start_date            AS start_date,
			l.end_date              AS end_date,
			l.back_date             AS back_date,
			l.description           AS description,
			stl.id                  AS status_leave_id,
			stl.name                AS status_leave_name,
			l.approve_by_id         AS approve_by_id,
			ep.name_kh              AS approve_by_name,
			b.id                    AS branch_id,
			b.name                  AS branch_name
		`).
		Joins("LEFT JOIN employees e ON e.id = l.employee_id").
		Joins("LEFT JOIN positions p ON p.id = e.position_id").
		Joins("LEFT JOIN offices o ON o.id = e.office_id").
		Joins("LEFT JOIN leave_types lt ON lt.id = l.leave_type_id").
		Joins("LEFT JOIN deduct_types ddt ON ddt.id = lt.deduct_type_id").
		Joins("LEFT JOIN status_leaves stl ON stl.id = l.status_leave_id").
		Joins("LEFT JOIN branches b ON b.id = l.branch_id").
		Joins("LEFT JOIN employees ep ON ep.id = l.approve_by_id").
		Joins("LEFT JOIN users u ON u.employee_id = e.id").
		Order("l.id DESC")

	if role.Level < 4 {

		query = query.Where("l.employee_id = ?", user.EmployeeID)

	} else {

		switch user.ManageBranch {
		case 1:
			query = query.Where("l.branch_id = ?", user.BranchID)

		case 2:

			var branchIDs []int
			if err := s.db.
				Model(&model.UserBranch{}).
				Where("user_id = ?", user.ID).
				Pluck("branch_id", &branchIDs).Error; err != nil {
				return nil, nil, fmt.Errorf("failed to fetch user branches: %w", err)
			}
			if len(branchIDs) == 0 {
				return []response.LeaveResponse{}, &model.PaginationMetadata{
					Page: pagination.Page, PageSize: pagination.PageSize,
					TotalCount: 0, TotalPages: 0,
				}, nil
			}
			query = query.Where("l.branch_id IN ?", branchIDs)

		case 3:

		}
	}

	for key, value := range filters {
		if value == "" {
			continue
		}
		switch key {
		case "employee_id":

			if role.Level >= 4 {
				query = query.Where("l.employee_id = ?", value)
			}

		case "branch_id":

			if role.Level >= 4 {
				query = query.Where("l.branch_id = ?", value)
			}

		case "office_id":

			query = query.Where("o.id = ?", value)

		case "status_leave_id":

			query = query.Where("l.status_leave_id = ?", value)

		case "leave_type_id":

			query = query.Where("l.leave_type_id = ?", value)

		case "start_date":

			query = query.Where("l.start_date >= ?", value)

		case "end_date":

			query = query.Where("l.end_date <= ?", value)

		case "search":

			like := "%" + value + "%"
			query = query.Where("(e.name_en LIKE ? OR e.name_kh LIKE ?)", like, like)

		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count leaves: %w", err)
	}

	if err := query.
		Limit(pagination.PageSize).
		Offset(offset).
		Scan(&leaves).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch leaves: %w", err)
	}

	for i := range leaves {
		var leaveDurations []response.LeaveDurationResponse
		if err := s.db.Table("leave_durations ld").
			Select(`
			ld.id AS id,
			ld.duration_value AS duration_value,
			ldu.id AS duration_unit_id,
			ldu.code AS duration_unit_code,
			ldu.name_en AS duration_unit_name_en,
			ldu.name_km AS duration_unit_name_kh
		`).
			Joins("LEFT JOIN leave_duration_units ldu ON ldu.id = ld.duration_unit_id").
			Where("ld.leave_id = ?", leaves[i].ID).Scan(&leaveDurations).Error; err != nil {
			return nil, nil, err
		}
		leaves[i].LeaveDurationResponse = leaveDurations
	}

	for i := range leaves {
		leaves[i].StartDate = helper.FormatDate(leaves[i].StartDate)
		leaves[i].EndDate = helper.FormatDate(leaves[i].EndDate)
		leaves[i].BackDate = helper.FormatDate(leaves[i].BackDate)
	}

	totalPages := int(totalCount) / pagination.PageSize
	if int(totalCount)%pagination.PageSize != 0 {
		totalPages++
	}

	metadata := &model.PaginationMetadata{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	return leaves, metadata, nil
}

func (s *leaveservice) ApproveLeave(
	id int,
	input request.LeaveApproveRequest,
	userID int,
) error {

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user model.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("user not found")
	}

	result := tx.Model(&model.Leave{}).
		Where("id = ? AND approve_by_id =?", id, user.EmployeeID).
		Updates(map[string]interface{}{
			"status_leave_id": input.StatusLeave,
		})

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("leave cannot approve")
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("អ្នកមិនអាចអនុម័តបានទេ")
	}

	var setting model.Setting

	if err := tx.Where("`key` = ?", "ATTENDANCEGROUPCHATID").First(&setting).Error; err != nil {
		tx.Rollback()
		return err
	}

	var leaveDurations []model.LeaveDuration

	if err := tx.Where("leave_id =?", id).Find(&leaveDurations).Error; err != nil {
		tx.Rollback()
		return err
	}

	var durationText string

	for i := range leaveDurations {
		var durationunit model.LeaveDurationUnit
		if err := tx.First(&durationunit, leaveDurations[i].DurationUnitID).Error; err != nil {
			tx.Rollback()
			return err
		}
		durationText += fmt.Sprintf(
			" %.0f %s",
			leaveDurations[i].DurationVlaue,
			durationunit.NameKh,
		)
	}

	var leave model.Leave
	if err := tx.First(&leave, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	var employee model.Employee

	if err := tx.First(&employee, leave.EmployeeID).Error; err != nil {
		tx.Rollback()
		return err
	}

	message := fmt.Sprintf(
		"សូមជម្រាបសួរលោកគ្រូ-អ្នកគ្រូ!\n"+
			"មានបុគ្គលិកឈ្មោះៈ %s\n"+
			"សុំអនុញ្ញាតច្បាប់ឈប់សម្រាក %s\n"+
			"ចាប់ពីថ្ងៃទី %s រហូតដល់ថ្ងៃទី %s នឹងចូលបម្រើការងារវិញនៅថ្ងៃទី %s\n"+
			"*មូលហេតុៈ %s\n"+
			"សូមអរគុណ",
		employee.NameKh,
		durationText,
		helper.FormatDate(leave.StartDate),
		helper.FormatDate(leave.EndDate),
		helper.FormatDate(leave.BackDate),
		leave.Description,
	)

	go helper.SendTelegramMessage(message, setting.Value)
	return tx.Commit().Error
}
