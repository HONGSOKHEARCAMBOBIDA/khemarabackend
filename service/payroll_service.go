package service

import (
	"fmt"
	"math"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PayrollService interface {
	CreatePayroll(userID int, input []request.PayrollRequestCreate) error
	DeletePayroll(id int) error
	GetDraftPayroll(branch_id int, currency_id int, payroll_type int) ([]response.PayrollDrafResponse, error)
	GetPayroll(userID int, filters map[string]string) ([]response.PayrollResponse, error)
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
					Principle:  payprinciple / rateusd.Rate * rate.Rate,
					Rate:       payinterest / rateusd.Rate * rate.Rate,
					Income:     (payprinciple + payinterest) / rateusd.Rate * rate.Rate,
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

func (s *payrollservice) GetDraftPayroll(branchID int, currencyID int, payrollType int) ([]response.PayrollDrafResponse, error) {
	var payrollDraft []response.PayrollDrafResponse

	pensionfundExpr := "0 AS pensionfund"
	if payrollType == 2 {
		pensionfundExpr = `
			s.base_salary * stpensionfund.value
				/ COALESCE(er_to_usd.rate, 1)
				* COALESCE(er_from_usd.rate, 1) / 100   AS pensionfund`
	}

	query := s.db.Table("employees e").
		Select(`
			e.id AS employee_id,
			e.name_kh AS employee_name,
			b.id branch_id,
			b.name AS branch_name,
			s.id AS salary_id,
			st.value AS total_work_day,
			s.base_salary
				/ COALESCE(er_to_usd.rate, 1)
				* COALESCE(er_from_usd.rate, 1) AS base_salary,

			s.daily_rate
				/ COALESCE(er_to_usd.rate, 1)
				* COALESCE(er_from_usd.rate, 1) AS daily_rate,
			`+pensionfundExpr+`,
			l.id AS loan_id,
			c.symbol AS currency_symbol,

			(
				SELECT COALESCE(SUM(
					CASE
						WHEN COALESCE(ps.income_amount, 0) != COALESCE(ps.income_paid, 0)
						AND DATE(ps.payment_date) <= CURRENT_DATE
							THEN (COALESCE(ps.income_amount, 0) - COALESCE(ps.income_paid, 0)) / COALESCE(loanrate.rate, 1) * COALESCE(er_from_usd.rate, 1)

						ELSE 0
					END
				), 0)
				FROM schedules ps
				WHERE ps.loan_id = l.id
			) AS loan_deduction
		`).
		Joins("LEFT JOIN users u ON u.employee_id = e.id").
		Joins("LEFT JOIN branches b ON b.id = u.branch_id").
		Joins("LEFT JOIN salaries s ON s.employee_id = e.id AND s.is_active = 1").
		Joins("LEFT JOIN loans l ON l.employee_id = e.id AND l.status = 1").
		Joins("LEFT JOIN settings st ON st.key = 'WORKDAY'").
		Joins("LEFT JOIN currency_pairs loanpair ON loanpair.base_currency_id = 2 AND loanpair.target_currency_id = l.currency_id").
		Joins("LEFT JOIN exchange_rates loanrate ON loanrate.pair_id = loanpair.id").
		Joins("LEFT JOIN currency_pairs cp_to_usd ON cp_to_usd.base_currency_id = 2 AND cp_to_usd.target_currency_id = s.currency_id").
		Joins("LEFT JOIN exchange_rates er_to_usd ON er_to_usd.pair_id = cp_to_usd.id").
		Joins("LEFT JOIN currency_pairs cp_from_usd ON cp_from_usd.base_currency_id = 2 AND cp_from_usd.target_currency_id = ?", currencyID).
		Joins("LEFT JOIN exchange_rates er_from_usd ON er_from_usd.pair_id = cp_from_usd.id").
		Joins("LEFT JOIN currencies c ON c.id =?", currencyID).
		Where("u.branch_id = ?", branchID)
	if payrollType == 2 {
		query = query.Joins("LEFT JOIN settings stpensionfund ON stpensionfund.key = 'PENSIONFUND'")
	}

	if err := query.Scan(&payrollDraft).Error; err != nil {
		return nil, err
	}

	return payrollDraft, nil
}

func (s *payrollservice) GetPayroll(userID int, filters map[string]string) ([]response.PayrollResponse, error) {
	var payroll []response.PayrollResponse
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	var role model.Role
	if err := s.db.First(&role, user.RoleID).Error; err != nil {
		return nil, err
	}
	query := s.db.Table("payrolls p").
		Select(`
		p.id AS id,
		b.name AS branch_name,
		e.name_en AS employee_name_en,
		e.name_kh AS employee_name_kh,
		e.gender AS employee_gender,
		ps.display_name AS position_name,
		o.name AS office_name,
		ep.profile_image AS profile_image,
		ep.bank_name AS bank_name,
		ep.bank_account_number AS bank_account_number,
		ep.qr_code_bank_account AS qr_code_bank_account,
		p.basic_salary AS basic_salary,
		p.half_salary AS half_salary,
		p.pension_fund AS pensionfund,
		p.total_work_day AS total_work_day,
		p.payroll_date AS payroll_date,
		p.loan_deduction AS loan_deduction,
		p.is_bonus AS is_bonus,
		bo.name AS bonus_type_name,
		p.bonus_amount AS bonus_amount,
		p.total_deduction AS total_deduction,
		p.net_salary AS net_salary,
		c.name AS currency_name,
		c.code AS currency_code,
		prs.name AS status_name,
		p.note AS note
	`).
		Joins("LEFT JOIN branches b ON b.id = p.branch_id").
		Joins("LEFT JOIN salaries ON salaries.id = p.salary_id AND salaries.is_active = 1").
		Joins("LEFT JOIN employees e ON e.id = salaries.employee_id").
		Joins("LEFT JOIN positions ps ON ps.id = e.position_id").
		Joins("LEFT JOIN offices o ON o.id = e.office_id").
		Joins("LEFT JOIN employee_profiles ep ON ep.employee_id = e.id").
		Joins("LEFT JOIN bonus_types bo ON bo.id = p.bonus_type").
		Joins("LEFT JOIN currencies c ON c.id = p.currency_id").
		Joins("LEFT JOIN pay_roll_statuses prs ON prs.id = p.status_id").
		Joins("LEFT JOIN payroll_approvals pa ON pa.payroll_id = p.id")

	if role.Level < 4 {
		query = query.Where("e.id = ?", user.EmployeeID)
	} else {
		switch user.ManageBranch {
		case 1:
			query = query.Where("p.branch_id =?", user.BranchID)
		case 2:
			var branchIDs []int
			if err := s.db.Model(&model.UserBranch{}).
				Where("user_id =?", user.ID).
				Pluck("branch_id", &branchIDs).Error; err != nil {
				return nil, fmt.Errorf("failed to fetch user branches: %w", err)
			}
			if len(branchIDs) == 0 {
				return []response.PayrollResponse{}, nil
			}
			query = query.Where("p.branch_id IN ?", branchIDs)
		case 3:
		}
	}

	for key, value := range filters {
		if value != "" {
			switch key {
			case "branch_id":
				query = query.Where("p.branch_id =?", value)
			case "name":
				query = query.Where("e.name_kh LIKE ? OR e.name_en LIKE ?", "%"+value+"%", "%"+value+"%")
			case "position_id":
				query = query.Where("e.position_id =?", value)
			case "status_id":
				query = query.Where("p.status_id =?", value)
			case "office_id":
				query = query.Where("e.office_id =?", value)
			case "department_id":
				query = query.Where("p.department_id =?", value)

			}
		}
	}

	if err := query.Scan(&payroll).Error; err != nil {
		return nil, err
	}
	if role.Level < 4 {
		for i := range payroll {
			payroll[i].ApproveStep1 = false
			payroll[i].ApproveStep2 = false
			payroll[i].ApproveStep3 = false
			payroll[i].ApproveStep4 = false
		}
	} else {
		var workflows []model.ApproveWorkflow
		if err := s.db.Where("role_name = ?", role.Name).Find(&workflows).Error; err != nil || len(workflows) == 0 {
			for i := range payroll {
				payroll[i].ApproveStep1 = false
				payroll[i].ApproveStep2 = false
				payroll[i].ApproveStep3 = false
				payroll[i].ApproveStep4 = false
			}
		} else {
			// Which steps this role owns
			allowedSteps := make(map[int]bool)
			for _, w := range workflows {
				allowedSteps[w.StepOrder] = true
			}

			// Fetch all approval records for these payrolls
			payrollIDs := make([]int, len(payroll))
			for i, p := range payroll {
				payrollIDs[i] = p.ID
			}

			type ApprovalRow struct {
				PayrollID int    `gorm:"column:payroll_id"`
				StepOrder int    `gorm:"column:step_order"`
				Status    string `gorm:"column:status"`
			}
			var approvals []ApprovalRow
			s.db.Table("payroll_approvals").
				Select("payroll_id, step_order, status").
				Where("payroll_id IN ?", payrollIDs).
				Scan(&approvals)

			// Build map: payrollID -> stepOrder -> isApproved
			approvedMap := make(map[int]map[int]bool)
			for _, a := range approvals {
				if approvedMap[a.PayrollID] == nil {
					approvedMap[a.PayrollID] = make(map[int]bool)
				}
				if strings.EqualFold(a.Status, "APPROVED") {
					approvedMap[a.PayrollID][a.StepOrder] = true
				}
			}

			for i := range payroll {
				pid := payroll[i].ID
				steps := approvedMap[pid] // which steps are already approved for this payroll

				// Step 1: show if role owns it (no prerequisite)
				if allowedSteps[1] {
					payroll[i].ApproveStep1 = true
				}
				// Step 2: show if role owns it AND step 1 already approved
				if allowedSteps[2] {
					payroll[i].ApproveStep2 = steps[1]
				}
				// Step 3: show if role owns it AND step 2 already approved
				if allowedSteps[3] {
					payroll[i].ApproveStep3 = steps[2]
				}
				// Step 4: show if role owns it AND step 3 already approved
				if allowedSteps[4] {
					payroll[i].ApproveStep4 = steps[3]
				}
			}
		}
	}
	for i := range payroll {
		payroll[i].PayrollDate = helper.FormatDate(payroll[i].PayrollDate)
	}
	return payroll, nil

}
