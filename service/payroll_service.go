package service

import (
	"fmt"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/utils"
	"time"

	"gorm.io/gorm"
)

type PayrollService interface {
	CreatePayroll(userID int, input []request.PayrollRequestCreate) error
}

type payrollservice struct {
	db *gorm.DB
}

func NewPayrollService() PayrollService {
	return &payrollservice{
		db: config.DB,
	}
}

func (s *payrollservice) CreatePayroll(userID int, input []request.PayrollRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if len(input) == 0 {
		tx.Rollback()
		return fmt.Errorf("no field")
	}

	now := time.Now().Format("2006-01-02")

	for _, inputs := range input {
		p := model.Payroll{
			SalaryID:       inputs.SalaryID,
			BranchID:       inputs.BranchID,
			PayRollTypeID:  inputs.PayRollTypeID,
			BasicSalary:    inputs.BasicSalary,
			HalfSalary:     inputs.HalfSalary,
			Pensionfund:    inputs.Pensionfund,
			TotalWorkDay:   inputs.TotalWorkDay,
			PayrollDate:    inputs.PayrollDate,
			LoanDeduction:  inputs.LoanDeduction,
			Isbonus:        inputs.Isbonus,
			BonusType:      inputs.BonusType,
			BonusAmount:    inputs.BonusAmount,
			TotalDeduction: inputs.TotalDeduction,
			NetSalary:      inputs.NetSalary,
			CurrencyID:     inputs.CurrencyID,
			StatusID:       1,
			SubmittedBy:    userID,
			SubmittedDate:  now,
			Note:           inputs.Note,
		}
		if err := tx.Create(&p).Error; err != nil {
			tx.Rollback()
			return err
		}

		pa := model.PayrollApproval{
			PayrollID:  p.ID,
			ApproveBy:  userID,
			Status:     "APPROVED",
			Comment:    "FIRST STEP",
			ActionDate: now,
			StepOrder:  1,
		}

		if err := tx.Create(&pa).Error; err != nil {
			tx.Rollback()
			return err
		}

		if inputs.Pensionfund > 0 {
			ps := model.Pensionfund{
				EmployeeID: inputs.EmployeeID,
				BranchID:   inputs.BranchID,
				Amount:     inputs.Pensionfund,
				CurrencyID: inputs.CurrencyID,
				Date:       now,
				PayrollID:  p.ID,
			}
			if err := tx.Create(&ps).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		if inputs.LoanDeduction > 0 {
			var loan model.Loan
			if err := tx.First(&loan, inputs.LoanID).Error; err != nil {
				continue
			}
			var pair model.CurrencyPair
			if err := tx.Where("base_currency_id =? AND target_currency_id =?", inputs.CurrencyID, loan.CurrencyID).First(&pair).Error; err != nil {
				tx.Rollback()
				return err
			}
			var rate model.ExchangeRate
			if err := tx.Where("pair_id =?", pair.ID).First(&rate).Error; err != nil {
				tx.Rollback()
				return err
			}
			recieve := model.Recieve{
				Code:         utils.GenerateRecieveCode(),
				BranchID:     inputs.BranchID,
				LoanID:       loan.ID,
				RecieveDate:  now,
				TotalRecieve: inputs.LoanDeduction,
				CurrencyID:   inputs.CurrencyID,
				Note:         inputs.Note,
				RecieveBy:    userID,
				PayrollID:    p.ID,
			}
			newremaining := inputs.TotalDeduction * rate.Rate

			for remaining > 0 {
				var schedule model.Schedule
				if err := tx.Where("loan_id =? AND status =?", loan.ID, 1).Order("schedule_number ASC").First(&schedule).Error; err != nil {
					tx.Rollback()
					return err
				}
				principlepaid, interestpaid := 0.0, 0.0
				if schedule.PrinciplePaid != nil {
					principlepaid = *schedule.PrinciplePaid
				}
				if schedule.RatePaid != nil {
					interestpaid = *schedule.IncomePaid
				}
			}
		}
	}
}
