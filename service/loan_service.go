package service

import (
	"fmt"
	"math"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type LoanService interface {
	CreateLoan(userID int, input request.LoanRequest) error
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

	// Generate loan code from last loan ID
	var lastLoan model.Loan
	tx.Order("id DESC").First(&lastLoan)
	newCode := fmt.Sprintf("LOAN-%05d", lastLoan.ID+1)

	// Get loan rate from settings
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

	// Count how many loans this employee has had
	var loanByEmployee model.Loan
	var numberLoan int
	if err := tx.Where("employee_id = ?", input.EmployeeID).Order("id DESC").First(&loanByEmployee).Error; err != nil {
		numberLoan = 1 // no previous loan
	} else {
		numberLoan = loanByEmployee.NumberofLoan + 1
	}

	// Calculations
	totalInterest := math.Round(input.LoanAmount*loanRate/100*100) / 100
	totalPeriods := input.LoanDuration * 2
	principlePerPeriod := math.Ceil(input.LoanAmount/float64(totalPeriods)*100) / 100
	interestPerPeriod := math.Ceil(totalInterest/float64(totalPeriods)*100) / 100
	monthlyPaymentAmount := math.Ceil((principlePerPeriod+interestPerPeriod)/100) * 100

	now := time.Now().Format("2006-01-02")

	loanStartDate, err := time.Parse("2006-01-02", input.LoanStartDate)
	if err != nil {
		tx.Rollback()
		return err
	}
	loanEndDate := loanStartDate.AddDate(0, input.LoanDuration, 0).Format("2006-01-02")

	// Create loan
	loan := model.Loan{
		Code:                 newCode,
		EmployeeID:           input.EmployeeID,
		LoanAmount:           input.LoanAmount,
		CurrencyID:           input.CurrencyID,
		ApproveDate:          now,
		LoanStartDate:        input.LoanStartDate,
		LoanEndDate:          loanEndDate,
		LoanRateAmount:       loanRate, // ✅ store rate, not total interest
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

	// Generate schedules (2 per month: 1st and 16th)
	schedules := make([]model.Schedule, 0, totalPeriods)
	remainingPrinciple := input.LoanAmount
	currentDate := loanStartDate

	for i := 1; i <= totalPeriods; i++ {
		var paymentDate time.Time
		if i%2 == 1 {
			// Odd → 1st of current month
			paymentDate = time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		} else {
			// Even → 16th of current month, then move to next month
			paymentDate = time.Date(currentDate.Year(), currentDate.Month(), 16, 0, 0, 0, 0, time.UTC)
			currentDate = currentDate.AddDate(0, 1, 0)
		}

		// Last period absorbs any rounding remainder
		periodPrinciple := principlePerPeriod
		if i == totalPeriods {
			periodPrinciple = math.Round(remainingPrinciple*100) / 100
		}
		remainingPrinciple -= periodPrinciple

		schedules = append(schedules, model.Schedule{
			LoanID:          loan.ID,
			PaymentDate:     paymentDate.Format("2006-01-02"),
			PrincipleAmount: periodPrinciple,
			RateAmount:      interestPerPeriod,
			IncomeAmount:    periodPrinciple + interestPerPeriod,
			PrinciplePaid:   0,
			RatePaid:        0,
			IncomePaid:      0,
			ScheduleNumber:  i,
			Status:          0,
		})
	}

	if err := tx.Create(&schedules).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
