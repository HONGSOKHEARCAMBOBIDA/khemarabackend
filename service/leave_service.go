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
	GetLeave(id int, filters map[string]string, pagination request.Pagination) ([]response.LeaveResponse, *model.PaginationMetadata, error)
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

	newLeave := model.Leave{
		EmployeeID:    user.EmployeeID,
		LeaveTypeID:   input.LeaveTypeID,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
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

	newLeaveDuration := model.LeaveDuration{
		LeaveID:        newLeave.ID,
		DurationVlaue:  input.DurationVlaue,
		DurationUnitID: input.DurationUnitID,
		StartTime:      nil,
		EndTime:        nil,
		Note:           nil,
	}

	if err := tx.Create(&newLeaveDuration).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create leave duration: %w", err)
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
			e.id                    AS employee_id,
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
			l.description           AS description,
			stl.id                  AS status_leave_id,
			stl.name                AS status_leave_name,
			l.approve_by_id         AS approve_by_id,
			ep.name_kh              AS approve_by_name,
			b.id                    AS branch_id,
			b.name                  AS branch_name,
			ld.id                   AS leave_duration_id,
			ld.duration_value       AS duration_value,
			ldn.id                  AS duration_unit_id,
			ldn.code                AS duration_unit_code,
			ldn.name_en             AS duration_unit_name_en,
			ldn.name_km             AS duration_unit_name_kh,
			ldn.to_minutes * ld.duration_value         AS duration_unit_tominute
		`).
		Joins("LEFT JOIN employees e ON e.id = l.employee_id").
		Joins("LEFT JOIN positions p ON p.id = e.position_id").
		Joins("LEFT JOIN offices o ON o.id = e.office_id").
		Joins("LEFT JOIN leave_types lt ON lt.id = l.leave_type_id").
		Joins("LEFT JOIN deduct_types ddt ON ddt.id = lt.deduct_type_id").
		Joins("LEFT JOIN status_leaves stl ON stl.id = l.status_leave_id").
		Joins("LEFT JOIN branches b ON b.id = l.branch_id").
		Joins("LEFT JOIN leave_durations ld ON ld.leave_id = l.id").
		Joins("LEFT JOIN leave_duration_units ldn ON ldn.id = ld.duration_unit_id").
		Joins("LEFT JOIN employees ep ON ep.id = l.approve_by_id").
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
		leaves[i].StartDate = helper.FormatDate(leaves[i].StartDate)
		leaves[i].EndDate = helper.FormatDate(leaves[i].EndDate)
		hours, display := helper.FormatDuration(leaves[i].DurationUnitToMinute)
		leaves[i].DurationHour = hours
		leaves[i].DurationDisplay = display
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
