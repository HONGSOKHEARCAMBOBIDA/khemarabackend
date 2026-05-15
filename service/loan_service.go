package service

import (
	"errors"
	"fmt"
	"math"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type LoanService interface {
	CreateLoan(userID int, input request.LoanRequest) error
	GetLoan(id int, filters map[string]string, pagination request.Pagination) ([]response.LoanResponse, *model.PaginationMetadata, error)
	DeleteLoan(id int) error
}

type loanservice struct {
	db *gorm.DB
}

func NewLoanService() LoanService {
	return &loanservice{
		db: config.DB,
	}
}

func (s *loanservice) CreateLoan(userID int, input request.LoanRequest) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var lastLoan model.Loan
	tx.Order("id DESC").First(&lastLoan)
	newCode := utils.GenerateLoanCode(lastLoan.ID)

	var setting model.Setting
	if err := tx.Where("`key` = ?", "LOANRATE").First(&setting).Error; err != nil {
		tx.Rollback()
		return err
	}

	loanRate, err := strconv.ParseFloat(setting.Value, 64)
	if err != nil {
		tx.Rollback()
		return err
	}

	var loanByEmployee model.Loan
	var numberLoan int
	if err := tx.Where("employee_id = ?", input.EmployeeID).Order("id DESC").First(&loanByEmployee).Error; err != nil {
		numberLoan = 1
	} else {
		numberLoan = loanByEmployee.NumberofLoan + 1
	}

	var user model.User
	if err := tx.Where("employee_id =?", input.EmployeeID).First(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	var pair model.CurrencyPair
	if err := tx.Where("base_currency_id =? AND target_currency_id =?", 2, input.CurrencyID).First(&pair).Error; err != nil {
		tx.Rollback()
		return err
	}

	var exchange model.ExchangeRate
	if err := tx.Where("pair_id =?", pair.ID).First(&exchange).Error; err != nil {
		tx.Rollback()
		return err
	}

	exchangetoUSD := math.Ceil(input.LoanAmount/exchange.Rate*100) / 100
	exchangetoSource := math.Round(exchangetoUSD*exchange.Rate*100) / 100

	totalInterest := math.Round(exchangetoSource*loanRate/100*100*float64(input.LoanDuration)) / 100
	totalPeriods := input.LoanDuration * 2
	principlePerPeriod := math.Round((exchangetoUSD/float64(totalPeriods))*exchange.Rate*100) / 100
	interestPerPeriod := math.Ceil(totalInterest/float64(totalPeriods)*100) / 100
	monthlyPaymentAmount := principlePerPeriod + interestPerPeriod

	now := time.Now().Format("2006-01-02")

	loanStartDate, err := time.Parse("2006-01-02", input.LoanStartDate)
	if err != nil {
		tx.Rollback()
		return err
	}

	loan := model.Loan{
		Code:                 newCode,
		BranchID:             user.BranchID,
		EmployeeID:           input.EmployeeID,
		LoanAmount:           input.LoanAmount,
		CurrencyID:           input.CurrencyID,
		ApproveDate:          now,
		LoanStartDate:        input.LoanStartDate,
		LoanEndDate:          nil,
		LoanRateAmount:       totalInterest,
		NumberofLoan:         numberLoan,
		ApproveBy:            userID,
		LoanPurpose:          input.LoanPurpose,
		MonthlyPaymentAmount: monthlyPaymentAmount,
		Status:               1,
		LoanDuration:         input.LoanDuration,
	}
	if err := tx.Create(&loan).Error; err != nil {
		tx.Rollback()
		return err
	}

	schedules := make([]model.Schedule, 0, totalPeriods)
	currentDate := loanStartDate

	for i := 1; i <= totalPeriods; i++ {
		var paymentDate time.Time

		if i%2 == 1 {
			paymentDate = time.Date(currentDate.Year(), currentDate.Month(), 15, 0, 0, 0, 0, time.UTC)
		} else {
			paymentDate = time.Date(currentDate.Year(), currentDate.Month(), 30, 0, 0, 0, 0, time.UTC)
			currentDate = currentDate.AddDate(0, 1, 0)
		}

		schedules = append(schedules, model.Schedule{
			LoanID:          loan.ID,
			PaymentDate:     paymentDate.Format("2006-01-02"),
			PaidDate:        nil,
			PrincipleAmount: principlePerPeriod,
			RateAmount:      interestPerPeriod,
			IncomeAmount:    principlePerPeriod + interestPerPeriod,
			PrinciplePaid:   nil,
			RatePaid:        nil,
			IncomePaid:      nil,
			ScheduleNumber:  i,
			Status:          1,
		})
	}

	if err := tx.Create(&schedules).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *loanservice) GetLoan(id int, filters map[string]string, pagination request.Pagination) ([]response.LoanResponse, *model.PaginationMetadata, error) {
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

	var loans []response.LoanResponse
	var totalCount int64

	offset := (pagination.Page - 1) * pagination.PageSize
	query := s.db.Table("loans l").
		Select(`
			l.id AS id,
			l.code AS code,
			b.id AS branch_id,
			b.name AS branch_name,
			e.id AS employee_id,
			e.name_kh AS employee_name,
			e.gender AS employee_gender,
			ep.dob AS employee_dob,
			u.contact AS employee_contact,
			e.code AS employee_code,
			l.loan_amount AS loan_amount,
			c.id AS currency_id,
			c.code AS currency_code,
			l.approve_date AS approve_date,
			l.loan_start_date AS loan_start_date,
			l.loan_end_date AS loan_end_date,
			l.number_of_loan AS number_of_loan,
			uu.id AS approve_by_id,
			eu.name_kh AS approve_by_name,
			l.loan_purpose AS loan_purpose,
			l.status AS loan_status,
			l.loan_duration AS loan_duration
		`).
		Joins("LEFT JOIN branches b ON b.id = l.branch_id").
		Joins("LEFT JOIN employees e ON e.id = l.employee_id").
		Joins("LEFT JOIN employee_profiles ep ON ep.employee_id = e.id").
		Joins("LEFT JOIN users u ON u.employee_id = e.id").
		Joins("LEFT JOIN currencies c ON c.id = l.currency_id").
		Joins("LEFT JOIN users uu ON uu.id = l.approve_by").
		Joins("LEFT JOIN employees eu ON eu.id = uu.employee_id")

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
				return []response.LoanResponse{}, &model.PaginationMetadata{
					Page:       pagination.Page,
					PageSize:   pagination.PageSize,
					TotalCount: 0,
					TotalPages: 0,
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
		case "search":
			like := "%" + value + "%"
			query = query.Where("(e.name_en LIKE ? OR e.name_kh LIKE ?)", like, like)
		case "employee_id":
			if role.Level >= 4 {
				query = query.Where("l.employee_id = ?", value)
			}
		case "branch_id":
			if role.Level >= 4 {
				query = query.Where("l.branch_id = ?", value)
			}
		case "status":
			query = query.Where("l.status = ?", value)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	if err := query.Offset(offset).Limit(pagination.PageSize).Scan(&loans).Error; err != nil {
		return nil, nil, err
	}

	for i := range loans {
		loans[i].EmployeeDob = helper.FormatDate(loans[i].EmployeeDob)
		loans[i].ApproveDate = helper.FormatDate(loans[i].ApproveDate)
		loans[i].LoanStartDate = helper.FormatDate(loans[i].LoanStartDate)
	}

	if len(loans) > 0 {
		loanIDs := make([]int, len(loans))
		for i, loan := range loans {
			loanIDs[i] = loan.ID
		}

		var allSchedules []response.ScheduleResponse
		if err := s.db.Table("schedules s").
			Select(`
				s.loan_id AS loan_id,
				s.id AS schedule_id,
				s.payment_date AS payment_date,
				s.paid_date AS paid_date,
				s.principle_amount AS principle_amount,
				s.rate_amount AS rate_amount,
				s.income_amount AS income_amount,
				s.principle_paid AS principle_paid,
				s.rate_paid AS rate_paid,
				s.income_paid AS income_paid,
				s.status AS status
			`).
			Where("s.loan_id IN ?", loanIDs).
			Scan(&allSchedules).Error; err != nil {
			return nil, nil, err
		}

		scheduleMap := make(map[int][]response.ScheduleResponse)
		for j := range allSchedules {
			allSchedules[j].PaidDate = helper.FormatDate(allSchedules[j].PaidDate)
			allSchedules[j].PaymentDate = helper.FormatDate(allSchedules[j].PaymentDate)
			loanID := allSchedules[j].LoanID
			scheduleMap[loanID] = append(scheduleMap[loanID], allSchedules[j])
		}

		// Assign grouped schedules back to each loan
		for i := range loans {
			loans[i].ScheduleResponse = scheduleMap[loans[i].ID]
			if loans[i].ScheduleResponse == nil {
				loans[i].ScheduleResponse = []response.ScheduleResponse{}
			}
		}
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PageSize)))
	metadata := &model.PaginationMetadata{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	return loans, metadata, nil
}

func (s *loanservice) DeleteLoan(id int) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var loan model.Loan

	if err := tx.First(&loan, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	today := time.Now().Format("2006-01-02")

	if loan.ApproveDate != today {
		tx.Rollback()
		return fmt.Errorf("can only delete loans approved today")
	}

	if err := tx.Where("loan_id =?", id).Delete(&model.Schedule{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.Loan{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
