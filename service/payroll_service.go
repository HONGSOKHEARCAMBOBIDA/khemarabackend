package service

import (
	"fmt"
	"math"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/utils"
	"time"

	"gorm.io/gorm"
)

type PayrollService interface {
	CreatePayroll(userID int, input []request.PayrollRequestCreate) error
	DeletePayroll(id int) error
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
				tx.Rollback()
				return err
			}

			var pairtoUSD model.CurrencyPair
			if err := tx.Where("base_currency_id = ? AND target_currency_id = ?", 2, inputs.CurrencyID).First(&pairtoUSD).Error; err != nil {
				tx.Rollback()
				return err
			}

			var rate model.ExchangeRate
			if err := tx.Where("pair_id = ?", pairtoUSD.ID).First(&rate).Error; err != nil {
				tx.Rollback()
				return err
			}

			exchangetoUSD := inputs.LoanDeduction / rate.Rate

			var pairtoLoan model.CurrencyPair
			if err := tx.Where("base_currency_id =? AND target_currency_id =?", 2, loan.CurrencyID).First(&pairtoLoan).Error; err != nil {
				tx.Rollback()
				return err
			}

			var rateusd model.ExchangeRate
			if err := tx.Where("pair_id =?", pairtoLoan.ID).First(&rateusd).Error; err != nil {
				tx.Rollback()
				return err
			}

			exchangetoloancurrency := exchangetoUSD * rateusd.Rate

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
			if err := tx.Create(&recieve).Error; err != nil {
				tx.Rollback()
				return err
			}

			newremaining := exchangetoloancurrency

			for newremaining > 0 {
				var schedule model.Schedule
				if err := tx.Where("loan_id = ? AND status = ?", loan.ID, 1).Order("schedule_number ASC").First(&schedule).Error; err != nil {
					tx.Rollback()
					return err
				}

				principlepaid := 0.0
				interestpaid := 0.0
				if schedule.PrinciplePaid != nil {
					principlepaid = *schedule.PrinciplePaid
				}
				if schedule.RatePaid != nil {
					interestpaid = *schedule.RatePaid
				}

				principledue := schedule.PrincipleAmount - principlepaid
				interestdue := schedule.RateAmount - interestpaid

				payprinciple := math.Min(newremaining, principledue)
				newremaining -= payprinciple
				payinterest := math.Min(newremaining, interestdue)
				newremaining -= payinterest

				newprinciple := principlepaid + payprinciple
				newinterest := interestpaid + payinterest

				paidincome := newprinciple + newinterest

				schedule.PrinciplePaid = &newprinciple
				schedule.RatePaid = &newinterest
				schedule.IncomePaid = &paidincome

				totaldue := schedule.PrincipleAmount + schedule.RateAmount
				if paidincome >= totaldue {
					schedule.PaidDate = &now
					schedule.Status = 0
				}

				updates := map[string]interface{}{
					"principle_paid": schedule.PrinciplePaid,
					"rate_paid":      schedule.RatePaid,
					"income_paid":    schedule.IncomePaid,
					"status":         schedule.Status,
				}
				if schedule.PaidDate != nil {
					updates["paid_date"] = schedule.PaidDate
				}

				if err := tx.Model(&model.Schedule{}).Where("id = ?", schedule.ID).Updates(updates).Error; err != nil {
					tx.Rollback()
					return err
				}

				recievedetail := model.RecieveDetail{
					LoanID:     loan.ID,
					RecieveID:  recieve.ID,
					ScheduleID: schedule.ID,
					Principle:  payprinciple,
					Rate:       payinterest,
					Income:     payprinciple + payinterest,
				}
				if err := tx.Create(&recievedetail).Error; err != nil {
					tx.Rollback()
					return err
				}
			}

			var pendingcount int64
			if err := tx.Model(&model.Schedule{}).Where("loan_id = ? AND status = ?", loan.ID, 1).Count(&pendingcount).Error; err != nil {
				tx.Rollback()
				return err
			}
			if pendingcount == 0 {
				if err := tx.Model(&model.Loan{}).Where("id = ?", loan.ID).Updates(map[string]interface{}{
					"status":        0,
					"loan_end_date": now,
				}).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	return tx.Commit().Error
}

func (s *payrollservice) DeletePayroll(id int) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var payroll model.Payroll
	if err := tx.First(&payroll, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if payroll.LoanDeduction > 0 {
		var recieve model.Recieve
		if err := tx.Where("payroll_id = ?", payroll.ID).First(&recieve).Error; err != nil {
			tx.Rollback()
			return err
		}

		var loan model.Loan
		if err := tx.First(&loan, recieve.LoanID).Error; err != nil {
			tx.Rollback()
			return err
		}

		var pairtoUSD model.CurrencyPair
		if err := tx.Where("base_currency_id = ? AND target_currency_id = ?", 2, payroll.CurrencyID).First(&pairtoUSD).Error; err != nil {
			tx.Rollback()
			return err
		}
		var rate model.ExchangeRate
		if err := tx.Where("pair_id = ?", pairtoUSD.ID).First(&rate).Error; err != nil {
			tx.Rollback()
			return err
		}
		exchangetoUSD := payroll.LoanDeduction / rate.Rate

		var pairtoLoan model.CurrencyPair
		if err := tx.Where("base_currency_id = ? AND target_currency_id = ?", 2, loan.CurrencyID).First(&pairtoLoan).Error; err != nil {
			tx.Rollback()
			return err
		}
		var rateusd model.ExchangeRate
		if err := tx.Where("pair_id = ?", pairtoLoan.ID).First(&rateusd).Error; err != nil {
			tx.Rollback()
			return err
		}
		exchangetoloancurrency := exchangetoUSD * rateusd.Rate

		var recievedetails []model.RecieveDetail
		if err := tx.Where("receive_id = ?", recieve.ID).Find(&recievedetails).Error; err != nil {
			tx.Rollback()
			return err
		}

		remaining := exchangetoloancurrency
		for _, detail := range recievedetails {
			var schedule model.Schedule
			if err := tx.First(&schedule, detail.ScheduleID).Error; err != nil {
				tx.Rollback()
				return err
			}

			principlepaid := 0.0
			interestpaid := 0.0
			if schedule.PrinciplePaid != nil {
				principlepaid = *schedule.PrinciplePaid
			}
			if schedule.RatePaid != nil {
				interestpaid = *schedule.RatePaid
			}

			newprinciple := math.Max(principlepaid-detail.Principle, 0)
			newinterest := math.Max(interestpaid-detail.Rate, 0)
			newincome := newprinciple + newinterest
			remaining -= (detail.Principle + detail.Rate)

			updates := map[string]interface{}{
				"principle_paid": newprinciple,
				"rate_paid":      newinterest,
				"income_paid":    newincome,
				"status":         1,
				"paid_date":      nil,
			}
			if err := tx.Model(&model.Schedule{}).Where("id = ?", schedule.ID).Updates(updates).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		if err := tx.Model(&model.Loan{}).Where("id = ?", loan.ID).Updates(map[string]interface{}{
			"status":        1,
			"loan_end_date": nil,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Where("receive_id = ?", recieve.ID).Delete(&model.RecieveDetail{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Delete(&model.Recieve{}, recieve.ID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if payroll.Pensionfund > 0 {
		if err := tx.Where("payroll_id = ?", id).Delete(&model.Pensionfund{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("payroll_id = ?", id).Delete(&model.PayrollApproval{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.Payroll{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
