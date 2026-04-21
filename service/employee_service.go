package service

import (
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type EmployeeService interface {
	GetEmployee(filters map[string]string, pagination request.Pagination) ([]response.EmployeeResponseDetail, *model.PaginationMetadata, error)
}

type employeeservice struct {
	db *gorm.DB
}

func NewEmployeeService() EmployeeService {
	return &employeeservice{
		db: config.DB,
	}
}

func (s *employeeservice) GetEmployee(filters map[string]string, pagination request.Pagination) ([]response.EmployeeResponseDetail, *model.PaginationMetadata, error) {
	var employees []response.EmployeeResponseDetail
	var totalCount int64

	offset := (pagination.Page - 1) * pagination.PageSize

	// pagination.Page is current page number
	// pagination.PageSize is number record per page
	// formula calculates how many records to skip before fetching data

	query := s.db.Table("users u").
		Select(`
			u.id AS user_id,
			u.username AS username,
			u.contact AS contact,
			b.id AS branch_id,
			b.name AS branch_name,
			r.id AS role_id,
			r.name AS role_name,
			r.display_name AS role_display_name,
			u.is_active AS user_active,
			mb.id AS manage_branch_id,
			mb.name AS manage_branch_name
		`).
		Joins("LEFT JOIN branches b ON b.id = u.branch_id").
		Joins("LEFT JOIN roles r ON r.id = u.role_id").
		Joins("LEFT JOIN manage_branches mb ON mb.id = u.manage_branch")

	for key, value := range filters {
		if value != "" {
			switch key {
			case "username":
				query = query.Where("u.username LIKE ?", "%"+value+"%")
			case "branch_id":
				query = query.Where("b.id = ?", value)
			case "role_id":
				query = query.Where("r.id = ?", value)
			case "is_active":
				query = query.Where("u.is_active = ?", value)
			}
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	if err := query.Offset(offset).Limit(pagination.PageSize).Find(&employees).Error; err != nil {
		return nil, nil, err
	}

	for i := range employees {
		var userParts []response.UserPartResponse
		if err := s.db.Table("user_parts up").
			Select(`
				up.id AS id,
				p.id AS part_id,
				p.name AS part_name,
				p.display_name AS part_display_name
			`).
			Joins("LEFT JOIN parts p ON p.id = up.part_id").
			Where("up.user_id = ?", employees[i].UserID).
			Scan(&userParts).Error; err != nil {
			return nil, nil, err
		}
		employees[i].Parts = userParts

		var userBranches []response.UserBranchResponse
		if err := s.db.Table("user_branches ub").
			Select(`
				ub.id AS id,
				ub.branch_id AS branch_id,
				b.name AS branch_name
			`).
			Joins("LEFT JOIN branches b ON b.id = ub.branch_id").
			Where("ub.user_id = ?", employees[i].UserID).
			Scan(&userBranches).Error; err != nil {
			return nil, nil, err
		}
		employees[i].Branches = userBranches

		var employeeDetails []response.EmployeeRespons
		if err := s.db.Table("employees e").
			Select(`
				e.id AS id,
				e.name_en AS name_en,
				e.name_kh AS name_kh,
				e.national_id_number AS national_id_number,
				e.gender AS gender,
				e.position_id AS position_id,
				p.display_name AS position_name,
				e.hire_date AS hire_date,
				e.promote_date AS promote_date,
				e.is_promote AS is_promote,
				e.employee_type_id AS employee_type_id,
				et.name AS employee_type_name,
				e.is_active AS is_active,
				e.office_id AS office_id,
				o.name AS office_name,
				o.display_name AS office_display_name,
				e.create_by AS create_by,
				cu.username AS create_by_name
			`).
			Joins("LEFT JOIN positions p ON p.id = e.position_id").
			Joins("LEFT JOIN employee_types et ON et.id = e.employee_type_id").
			Joins("LEFT JOIN offices o ON o.id = e.office_id").
			Joins("LEFT JOIN users cu ON cu.id = e.create_by").
			Where("e.id = (SELECT employee_id FROM users WHERE id = ?)", employees[i].UserID).
			Scan(&employeeDetails).Error; err != nil {
			return nil, nil, err
		}
		for i := range employeeDetails {
			employeeDetails[i].HireDate = helper.FormatDate(employeeDetails[i].HireDate)
			employeeDetails[i].PromoteDate = helper.FormatDate(employeeDetails[i].PromoteDate)
		}
		employees[i].EmployeeRespons = employeeDetails

		if len(employeeDetails) > 0 {
			employeeID := employeeDetails[0].ID

			// Get employee education
			var educations []response.EmployeeEducationRespons
			if err := s.db.Table("employee_educations ee").
				Select(`
					ee.id AS id,
					ee.education_level_id AS education_level_id,
					el.name AS education_level_name,
					ee.major_field_of_study AS major_field_of_study,
					ee.start_date AS start_date,
					ee.end_date AS end_date,
					ee.note AS note,
					ee.image AS image,
					ee.create_by AS create_by,
					cu.username AS create_by_name,
					ee.update_by AS update_by
				`).
				Joins("LEFT JOIN education_levels el ON el.id = ee.education_level_id").
				Joins("LEFT JOIN users cu ON cu.id = ee.create_by").
				Where("ee.employee_id = ?", employeeID).
				Scan(&educations).Error; err != nil {
				return nil, nil, err
			}
			for i := range educations {
				educations[i].StartDate = helper.FormatDate(educations[i].StartDate)
				educations[i].EndDate = helper.FormatDate(educations[i].EndDate)
			}
			employees[i].EmployeeEducations = educations

			// Get employee profile
			var profiles []response.EmployeeProfileResponse
			if err := s.db.Table("employee_profiles ep").
				Select(`
					ep.id AS id,
					ep.position_level_id AS position_level_id,
					pl.display_name AS position_level_name,
					ep.dob AS dob,
					pb.id AS province_id_birth,
					pb.name AS province_name_birth,
					db.id AS district_id_birth,
					db.name AS district_name_birth,
					cb.id AS communce_id_birth,
					cb.name AS communce_name_birth,
					vb.id AS village_id_birth,
					vb.name AS village_name_birth,
					ep.profile_image AS profile_image,
					pc.id AS province_id_current,
					pc.name AS province_name_current,
					dc.id AS distirct_id_current,
					dc.name AS district_name_current,
					cc.id AS communce_id_current,
					cc.name AS communce_name_current,
					vc.id AS village_id_current,
					vc.name AS village_name_current,
					ep.family_phone AS family_phone,
					ep.bank_name AS bank_name,
					ep.bank_account_number AS bank_account_number,
					ep.qr_code_bank_account AS qr_code_bank_account,
					ep.note AS note,
					ep.report_to AS report_to,
					ru.name_kh AS report_to_name,
					ep.wife_name AS wife_name,
					ep.husban_name AS husban_name,
					ep.son_number AS son_number,
					ep.daughter_number AS daughter_number,
					ep.material_status AS material_status
				`).
				Joins("LEFT JOIN position_levels pl ON pl.id = ep.position_level_id").
				Joins("LEFT JOIN villages vb ON vb.id = ep.village_id_of_birth").
				Joins("LEFT JOIN communces cb ON cb.id = vb.communce_id").
				Joins("LEFT JOIN districts db ON db.id = cb.district_id").
				Joins("LEFT JOIN provinces pb ON pb.id = db.province_id").
				Joins("LEFT JOIN villages vc ON vc.id = ep.village_id_current_address").
				Joins("LEFT JOIN communces cc ON cc.id = vc.communce_id").
				Joins("LEFT JOIN districts dc ON dc.id = cc.district_id").
				Joins("LEFT JOIN provinces pc ON pc.id = dc.province_id").
				Joins("LEFT JOIN employees ru ON ru.id = ep.report_to").
				Where("ep.employee_id = ?", employeeID).
				Scan(&profiles).Error; err != nil {
				return nil, nil, err
			}
			for i := range profiles {
				profiles[i].DoB = helper.FormatDate(profiles[i].DoB)
			}
			employees[i].EmployeeProfies = profiles

			// Get work experiences
			var workExperiences []response.EmployeeWorkExperienceResponse
			if err := s.db.Table("employee_work_experiences ewe").
				Select(`
					ewe.id AS id,
					ewe.company_name AS company_name,
					ewe.position_title AS position_title,
					ewe.start_date AS start_date,
					ewe.end_date AS end_date,
					ewe.job_description AS job_description
				`).
				Where("ewe.employee_id = ?", employeeID).
				Scan(&workExperiences).Error; err != nil {
				return nil, nil, err
			}
			for i := range workExperiences {
				workExperiences[i].StartDate = helper.FormatDate(workExperiences[i].StartDate)
				workExperiences[i].EndDate = helper.FormatDate(workExperiences[i].EndDate)
			}
			employees[i].EmployeeWorkExperiences = workExperiences

			// Get salaries
			var salaries []response.SalaryResponse
			if err := s.db.Table("salaries s").
				Select(`
					s.id AS id,
					s.base_salary AS base_salary,
					s.work_day AS work_day,
					s.daily_rate AS daily_rate,
					s.effective_date AS effective_date,
					s.expire_date AS expire_date,
					s.currency_id AS currency_id,
					c.code AS currency_code,
					c.symbol AS currency_symbol,
					c.name AS currency_name,
					s.is_active AS is_active
				`).
				Joins("LEFT JOIN currencies c ON c.id = s.currency_id").
				Where("s.employee_id = ?", employeeID).
				Scan(&salaries).Error; err != nil {
				return nil, nil, err
			}
			for i := range salaries {
				salaries[i].EffectiveDate = helper.FormatDate(salaries[i].EffectiveDate)
				salaries[i].ExpireDate = helper.FormatDate(salaries[i].ExpireDate)
			}
			employees[i].Salarys = salaries

			// Get shift patterns
			var shiftPatterns []response.ShiftPatternResponse
			if err := s.db.Table("shift_patterns esp").
				Select(`
					esp.id AS id,
					esp.day_of_week_id AS day_of_week_id,
					esp.shift_id AS shift_id,
					esp.is_day_off AS is_dayoff,
					dow.name AS day_of_week_name
				`).
				Joins("LEFT JOIN day_of_weeks dow ON dow.id = esp.day_of_week_id").
				Where("esp.employee_id = ?", employeeID).
				Scan(&shiftPatterns).Error; err != nil {
				return nil, nil, err
			}

			// Get shifts for each pattern
			for j := range shiftPatterns {
				var shifts []response.ShiftResponse
				if err := s.db.Table("shifts s").
					Select(`
						s.id AS id,
						s.name AS name,
						s.is_active AS is_active,
						s.branch_id AS branch_id,
						b.name AS branch_name
					`).
					Joins("LEFT JOIN branches b ON b.id = s.branch_id").
					Where("s.id = ?", shiftPatterns[j].ShiftID).
					Scan(&shifts).Error; err != nil {
					return nil, nil, err
				}

				// Get shift sessions for each shift
				for k := range shifts {
					var sessions []response.ShiftSessionResponse
					if err := s.db.Table("shift_sessions ss").
						Select(`
							ss.id AS id,
							ss.session_name AS session_name,
							ss.shift_order AS shift_order,
							ss.start_time AS start_time,
							ss.end_time AS end_time,
							ss.is_active AS is_active
						`).
						Where("ss.shift_id = ?", shifts[k].ID).
						Order("ss.shift_order ASC").
						Scan(&sessions).Error; err != nil {
						return nil, nil, err
					}
					shifts[k].ShiftSessionResponse = sessions
				}
				shiftPatterns[j].ShiftResponse = shifts
			}
			employees[i].ShiftPatterns = shiftPatterns
		}
	}

	// Create pagination metadata
	paginationMeta := &model.PaginationMetadata{
		CurrentPage: pagination.Page,
		PageSize:    pagination.PageSize,
		TotalCount:  totalCount,
		TotalPages:  (int(totalCount) + pagination.PageSize - 1) / pagination.PageSize,
	}

	return employees, paginationMeta, nil
}
