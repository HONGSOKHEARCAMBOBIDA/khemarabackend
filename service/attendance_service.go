package service

import (
	"fmt"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AttendanceService interface {
	CheckIn(id int, input request.LocationRequest) error
	CheckOut(id int, input request.LocationRequest) error
	GetAttendance(filter map[string]string, pagination request.Pagination) ([]response.AttendanceResponse, *model.PaginationMetadata, error)
}

type attendanceservice struct {
	db *gorm.DB
}

func NewAttendanceService() AttendanceService {
	return &attendanceservice{
		db: config.DB,
	}
}

func (s *attendanceservice) CheckIn(id int, input request.LocationRequest) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	currentDate := now.Format("2006-01-02")

	dayOfWeek := (int(now.Weekday())+6)%7 + 1
	var user model.User
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	var branch model.Branch
	if err := tx.First(&branch, user.BranchID).Error; err != nil {
		tx.Rollback()
		return err
	}
	var shiftpattern model.ShiftPattern
	if err := tx.Where("employee_id = ? AND day_of_week_id = ?", user.EmployeeID, dayOfWeek).
		First(&shiftpattern).Error; err != nil {
		tx.Rollback()
		return err
	}

	if shiftpattern.Isdayoff {
		tx.Rollback()
		return nil
	}

	branchLat, err := strconv.ParseFloat(branch.Latitude, 64)
	if err != nil {
		tx.Rollback()
		return err
	}
	branchLng, err := strconv.ParseFloat(branch.Longitude, 64)
	if err != nil {
		tx.Rollback()
		return err
	}
	userLat, err := strconv.ParseFloat(input.Latitude, 64)
	if err != nil {
		tx.Rollback()
		return err
	}
	userLng, err := strconv.ParseFloat(input.Longitude, 64)
	if err != nil {
		tx.Rollback()
		return err
	}

	distance := utils.CalculateDistance(branchLat, branchLng, userLat, userLng)
	inzone := distance <= float64(input.BranchRadius)

	var attendancelog model.AttendanceLog
	err = tx.Where("employee_id = ? AND check_date = ?", user.EmployeeID, currentDate).
		Order("id DESC").
		First(&attendancelog).Error

	var session model.ShiftSession
	var shiftOrder int

	if err != nil {

		if err := tx.Where("shift_id = ?", shiftpattern.ShiftID).
			Order("shift_order ASC").
			First(&session).Error; err != nil {
			tx.Rollback()
			return err
		}
		shiftOrder = session.ShiftOrder
	} else {
		if attendancelog.StatusAttendanceLogID == 1 {
			tx.Rollback()
			return fmt.Errorf("អ្នកមិនទាន់ចុចចេញពីធ្វេីការទេ!")
		}
		shiftOrder = attendancelog.ShiftSessionOrder + 1

		if err := tx.Where("shift_id = ? AND shift_order = ?", shiftpattern.ShiftID, shiftOrder).
			First(&session).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	startTime, err := time.Parse("15:04:05", session.StartTime)
	if err != nil {
		return err
	}

	shiftStart := time.Date(
		now.Year(), now.Month(), now.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second(),
		0, now.Location(),
	)

	checkInEarly := 0
	if now.Before(shiftStart) {
		checkInEarly = 1
	}

	isLate := 0
	if now.After(shiftStart) {
		isLate = 1
	}

	checkInOnTime := 0
	if now.Equal(shiftStart) {
		checkInOnTime = 1
	}

	newLog := model.AttendanceLog{
		EmployeeID:            user.EmployeeID,
		CheckDate:             currentDate,
		Note:                  "",
		BranchID:              int(branch.ID),
		StatusAttendanceLogID: 1,
		ShiftSessionOrder:     shiftOrder,
	}

	if err := tx.Create(&newLog).Error; err != nil {
		tx.Rollback()
		return err
	}

	newRecord := model.AttendanceRecord{
		AttendanceLogID:  newLog.ID,
		ShiftSessionID:   session.ID,
		CheckTime:        now,
		CheckInEarly:     &checkInEarly,
		CheckInOnTime:    &checkInOnTime,
		IsLate:           &isLate,
		IsLeftEarly:      nil,
		CheckOutOnTime:   nil,
		CheckOutOverTime: nil,
		Latitude:         userLat,
		Logitude:         userLng,
		Note:             input.Note,
		Iszoone:          inzone,
		Type:             "IN",
	}

	if err := tx.Create(&newRecord).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err == nil {
		if err := tx.Model(&model.AttendanceLog{}).
			Where("id = ?", attendancelog.ID).
			Update("shift_session_order", shiftOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// var employee model.Employee
	// if err := tx.First(&employee, user.EmployeeID).Error; err != nil {
	// 	tx.Rollback()
	// }
	// workTime := fmt.Sprintf("%s - %s", session.StartTime, session.EndTime)
	// lateText := "⏰ ស្កែនទាន់ម៉ោង"
	// if isLate == 1 {
	// 	lateText = "🔴 ចូលធ្វេីការយឺត"
	// }
	// zoneText := "📍 ស្កែនក្នុងតំបន់ក្រុមហ៊ុន"
	// if !inzone {
	// 	zoneText = "⚠️ ស្កែនក្រៅតំបន់ក្រុមហ៊ុន"
	// }
	// message := fmt.Sprintf(
	// 	"🟢 <b>CHECK IN</b>\n\n"+
	// 		"👤 ឈ្មោះ: %s\n"+
	// 		"📲 ឈ្មោះអង់គ្លេស: %s\n"+
	// 		"ID: %s\n"+
	// 		"🏢 សាខា: %s\n"+
	// 		"🕒 ម៉ោងធ្វើការ: %s\n"+
	// 		"🕒 Check-in: %s\n"+
	// 		"%s\n"+
	// 		"%s\n"+
	// 		"មូលហេតុ: %s\n",

	// 	employee.NameKh,
	// 	employee.NameEn,
	// 	employee.Code,
	// 	branch.Name,
	// 	workTime,
	// 	now.Format("15:04:05"),
	// 	lateText,
	// 	zoneText,
	// 	input.Note,
	// )
	// go helper.SendTelegramMessage(message)
	return tx.Commit().Error
}

func (s *attendanceservice) CheckOut(id int, input request.LocationRequest) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	dayOfWeek := int(now.Weekday())
	dayOfWeek = (dayOfWeek+6)%7 + 1
	var user model.User
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	var branch model.Branch
	if err := tx.First(&branch, user.BranchID).Error; err != nil {
		tx.Rollback()
		return err
	}
	var shiftpattern model.ShiftPattern
	if err := tx.Where("employee_id = ? AND day_of_week_id = ?", user.EmployeeID, dayOfWeek).
		First(&shiftpattern).Error; err != nil {
		tx.Rollback()
		return err
	}

	if shiftpattern.Isdayoff {
		tx.Rollback()
		return nil
	}

	var attendancelog model.AttendanceLog
	err := tx.Where("employee_id =? AND status_attendance_log_id =?", user.EmployeeID, 1).Order("id DESC").First(&attendancelog).Error
	if err != nil {
		tx.Rollback()
		return err
	} else {
		var attendancerecore model.AttendanceRecordRes
		if err := tx.Where("attendance_log_id =?", attendancelog.ID).First(&attendancerecore).Error; err != nil {
			tx.Rollback()
			return err
		}
		var session model.ShiftSession
		if err := tx.Where("id =?", attendancerecore.ShiftSessionID).First(&session).Error; err != nil {
			tx.Rollback()
			return err
		}
		lat, err := strconv.ParseFloat(branch.Latitude, 64)
		if err != nil {
			tx.Rollback()
			return err
		}
		log, err := strconv.ParseFloat(branch.Longitude, 64)
		if err != nil {
			tx.Rollback()
			return err
		}
		latput, err := strconv.ParseFloat(input.Latitude, 64)
		if err != nil {
			tx.Rollback()
			return err
		}
		logput, err := strconv.ParseFloat(input.Longitude, 64)
		if err != nil {
			tx.Rollback()
			return err
		}
		distance := utils.CalculateDistance(lat, log, latput, logput)
		inzone := distance <= float64(input.BranchRadius)
		now := time.Now()

		endTime, err := time.Parse("15:04:05", session.EndTime)
		if err != nil {
			return err
		}

		currentDate := now.Format("2006-01-02")

		shiftEnd := time.Date(
			now.Year(), now.Month(), now.Day(),
			endTime.Hour(), endTime.Minute(), endTime.Second(),
			0, now.Location(),
		)

		isLeftEarly := 0
		if now.Before(shiftEnd) {
			isLeftEarly = 1
		}

		checkOutOverTime := 0
		if now.After(shiftEnd) {
			checkOutOverTime = 1
		}

		checkOutOnTime := 0
		grace := 1 * time.Minute

		if now.After(shiftEnd.Add(-grace)) && now.Before(shiftEnd.Add(grace)) {
			checkOutOnTime = 1
		}
		attendancelog.StatusAttendanceLogID = 2
		attendancelog.CheckDate = currentDate
		if err := tx.Save(&attendancelog).Error; err != nil {
			tx.Rollback()
			return err
		}
		newattendancerecord := model.AttendanceRecord{
			AttendanceLogID:  attendancelog.ID,
			ShiftSessionID:   session.ID,
			CheckTime:        now,
			CheckInEarly:     nil,
			CheckInOnTime:    nil,
			IsLate:           nil,
			IsLeftEarly:      &isLeftEarly,
			CheckOutOnTime:   &checkOutOnTime,
			CheckOutOverTime: &checkOutOverTime,
			Latitude:         latput,
			Logitude:         logput,
			Note:             input.Note,
			Iszoone:          inzone,
			Type:             "OUT",
		}
		if err := tx.Create(&newattendancerecord).Error; err != nil {
			tx.Rollback()
			return err
		}
		// var employee model.Employee
		// if err := tx.First(&employee, user.EmployeeID).Error; err != nil {
		// 	tx.Rollback()
		// }
		// workTime := fmt.Sprintf("%s - %s", session.StartTime, session.EndTime)

		// earlyText := "⏰ ស្កែនត្រូវម៉ោង"
		// if isLeftEarly == 1 {
		// 	earlyText = "🔴 ចេញមុនម៉ោងកំណត់"
		// }

		// zoneText := "📍 ស្កែនក្នុងតំបន់ក្រុមហ៊ុន"
		// if !inzone {
		// 	zoneText = "⚠️ ស្កែនក្រៅតំបន់ក្រុមហ៊ុន"
		// }

		// message := fmt.Sprintf(
		// 	"🟢 <b>CHECK OUT</b>\n\n"+
		// 		"👤 ឈ្មោះ: %s\n"+
		// 		"📲 ឈ្មោះអង់គ្លេស: %s\n"+
		// 		"ID: %s\n"+
		// 		"🏢 សាខា: %s\n"+
		// 		"🕒 ម៉ោងធ្វើការ: %s\n"+
		// 		"🕒 Check-out: %s\n"+
		// 		"%s\n"+
		// 		"%s\n",

		// 	employee.NameKh,
		// 	employee.NameEn,
		// 	employee.Code,
		// 	branch.Name,
		// 	workTime,
		// 	now.Format("15:04:05"),
		// 	earlyText,
		// 	zoneText,
		// )

		// go helper.SendTelegramMessage(message)
	}

	return tx.Commit().Error
}

func (s *attendanceservice) GetAttendance(filter map[string]string, pagination request.Pagination) ([]response.AttendanceResponse, *model.PaginationMetadata, error) {
	var attendances []response.AttendanceResponse
	var totalCount int64
	offset := (pagination.Page - 1) * pagination.PageSize

	query := s.db.Table("employees e").
		Select(`
			e.id AS employee_id,
			e.code AS employee_code,
			e.name_en AS employee_name_en,
			e.name_kh AS employee_name_kh,
			e.gender AS gender,
			p.id AS position_id,
			p.display_name AS position_display_name,
			d.id AS department_id,
			d.display_name AS department_display_name,
			o.id AS office_id,
			o.name AS office_name,
			ep.profile_image AS profile
		`).
		Joins("INNER JOIN attendance_logs alog ON alog.employee_id = e.id").
		Joins("LEFT JOIN positions p ON p.id = e.position_id").
		Joins("LEFT JOIN departments d ON d.id = p.department_id").
		Joins("LEFT JOIN offices o ON o.id = e.office_id").
		Joins("LEFT JOIN employee_profiles ep ON ep.employee_id = e.id").
		Group("e.id, e.code, e.name_en, e.name_kh, e.gender, p.id, p.display_name, d.id, d.display_name, o.id, o.name, ep.profile_image")

	checkDateFrom := filter["check_date_from"]
	checkDateTo := filter["check_date_to"]
	boolFilterMap := map[string]string{
		"check_in_early":     "check_in_early",
		"check_in_on_time":   "check_in_on_time",
		"is_late":            "is_late",
		"is_left_early":      "is_left_early",
		"check_out_on_time":  "check_out_on_time",
		"check_out_overtime": "check_out_overtime",
	}

	for key, value := range filter {
		if value == "" {
			continue
		}
		switch key {
		case "name":
			query = query.Where("e.name_kh LIKE ? OR e.name_en LIKE ?", "%"+value+"%", "%"+value+"%")
		case "department_id":
			query = query.Where("d.id = ?", value)
		case "employee_id":
			query = query.Where("e.id = ?", value)
		case "office_id":
			query = query.Where("e.office_id = ?", value)

		case "check_date_from":
			query = query.Where("alog.check_date >= ?", value)
		case "check_date_to":
			query = query.Where("alog.check_date <= ?", value)

		case "check_in_early", "check_in_on_time", "is_late",
			"is_left_early", "check_out_on_time", "check_out_overtime":
			col := boolFilterMap[key]
			query = query.Where(fmt.Sprintf(`
				EXISTS (
					SELECT 1 FROM attendance_records ar
					JOIN attendance_logs al2 ON al2.id = ar.attendance_log_id
					WHERE al2.employee_id = e.id
					AND ar.%s = ?
				)`, col), value)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	if err := query.Offset(offset).Limit(pagination.PageSize).Scan(&attendances).Error; err != nil {
		return nil, nil, err
	}

	for i := range attendances {
		var attendancelogs []response.AttendanceLogResponse

		logQuery := s.db.Table("attendance_logs al").
			Select(`
				al.id AS id,
				al.check_date AS check_date,
				b.id AS branch_id,
				b.name AS branch_name,
				st.id AS status_attendance_log_id,
				st.name AS status_attendance_log_name
			`).
			Joins("LEFT JOIN branches b ON b.id = al.branch_id").
			Joins("LEFT JOIN status_attendance_logs st ON st.id = al.status_attendance_log_id").
			Where("al.employee_id = ?", attendances[i].EmployeeID)

		if checkDateFrom != "" {
			logQuery = logQuery.Where("al.check_date >= ?", checkDateFrom)
		}
		if checkDateTo != "" {
			logQuery = logQuery.Where("al.check_date <= ?", checkDateTo)
		}

		for key, col := range boolFilterMap {
			if val := filter[key]; val != "" {
				logQuery = logQuery.Where(fmt.Sprintf(`
					EXISTS (
						SELECT 1 FROM attendance_records ar
						WHERE ar.attendance_log_id = al.id
						AND ar.%s = ?
					)`, col), val)
			}
		}

		if err := logQuery.Scan(&attendancelogs).Error; err != nil {
			return nil, nil, err
		}

		for j := range attendancelogs {
			var attendancerecords []response.AttendanceRecordResponse

			recordQuery := s.db.Table("attendance_records ar").
				Select(`
					ar.id AS id,
					s.id AS shift_session_id,
					s.start_time AS start_time,
					s.end_time AS end_time,
					s.shift_order AS shift_order,
					ar.check_time AS check_time,
					ar.check_in_early AS check_in_early,
					ar.check_in_on_time AS check_in_on_time,
					ar.is_late AS is_late,
					ar.is_left_early AS is_left_early,
					ar.check_out_on_time AS check_out_on_time,
					ar.check_out_overtime AS check_out_overtime,
					ar.latitude AS latitude,
					ar.longitude AS longitude,
					ar.note AS note,
					ar.iszonecheckin AS iszonecheckin,
					ar.type AS type
				`).
				Joins("LEFT JOIN shift_sessions s ON s.id = ar.shift_session_id").
				Where("ar.attendance_log_id = ?", attendancelogs[j].ID)

			for key, col := range boolFilterMap {
				if val := filter[key]; val != "" {
					recordQuery = recordQuery.Where(fmt.Sprintf("ar.%s = ?", col), val)
				}
			}

			if err := recordQuery.Scan(&attendancerecords).Error; err != nil {
				return nil, nil, err
			}

			attendancelogs[j].AttendanceRecordResponse = attendancerecords
		}

		attendances[i].AttendanceLogResponse = attendancelogs
	}

	paginationMeta := &model.PaginationMetadata{
		CurrentPage: pagination.Page,
		PageSize:    pagination.PageSize,
		TotalCount:  totalCount,
		TotalPages:  (int(totalCount) + pagination.PageSize - 1) / pagination.PageSize,
	}

	return attendances, paginationMeta, nil
}
