package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"time"

	"gorm.io/gorm"
)

type ShiftSessionService interface {
	GetAllShiftSession() ([]response.ShiftSessionResponse, error)
	GetByShiftID(id int) ([]response.ShiftSessionResponse, error)
	CreateShiftSession(input request.ShiftSessionRequestCreate) error
	UpdateShiftSession(id int, input request.ShiftSessionRequestUpdate) error
	ChangeStatusShiftSession(id int) error
	GetShiftSessionV2(id int) (response.ShiftSessionResponsev2, error)
}

type shiftsessionservice struct {
	db *gorm.DB
}

func Newshiftsessionservice() ShiftSessionService {
	return &shiftsessionservice{
		db: config.DB,
	}
}

func (s *shiftsessionservice) GetAllShiftSession() ([]response.ShiftSessionResponse, error) {
	var shiftsessions []response.ShiftSessionResponse
	db := s.db.Table("shift_sessions ss").
		Select(`
	ss.id AS id,
	ss.session_name AS session_name,
	ss.shift_order AS shift_order,
	ss.start_time AS start_time,
	ss.end_time AS end_time,
	ss.is_active AS is_active,
	s.id AS shift_id,
	s.name AS shift_name,
	b.id AS branch_id,
	b.  name AS branch_name
	`).
		Joins("LEFT JOIN shifts s ON s.id = ss.shift_id").
		Joins("LEFT JOIN branches b ON b.id = s.branch_id")
	if err := db.Order("ss.id DESC").Scan(&shiftsessions).Error; err != nil {
		return nil, err
	}
	return shiftsessions, nil
}

func (s *shiftsessionservice) GetByShiftID(id int) ([]response.ShiftSessionResponse, error) {
	var shiftsessions []response.ShiftSessionResponse
	db := s.db.Table("shift_sessions ss").
		Select(`
	ss.id AS id,
	ss.session_name AS session_name,
	ss.shift_order AS shift_order,
	ss.start_time AS start_time,
	ss.end_time AS end_time,
	ss.is_active AS is_active,
	s.id AS shift_id,
	s.name AS shift_name,
	b.id AS branch_id,
	b.name AS branch_name
	`).
		Joins("LEFT JOIN shifts s ON s.id = ss.shift_id").
		Joins("LEFT JOIN branches b ON b.id = s.branch_id")
	if err := db.Where("s.id = ?", id).Scan(&shiftsessions).Error; err != nil {
		return nil, err
	}
	return shiftsessions, nil
}

func (s *shiftsessionservice) CreateShiftSession(input request.ShiftSessionRequestCreate) error {
	if input.SessionName == "" {
		return errors.New("session name is required")
	}
	if input.ShiftID == 0 {
		return errors.New("shift id is required")
	}
	if input.StartTime == "" {
		return errors.New("start time is required")
	}
	if input.EndTime == "" {
		return errors.New("end time is required")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	var lastSession model.ShiftSession
	err := tx.Where("shift_id = ?", input.ShiftID).Order("shift_order DESC").First(&lastSession).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	newShiftSession := model.ShiftSession{
		SessionName: input.SessionName,
		ShiftID:     input.ShiftID,
		ShiftOrder:  lastSession.ShiftOrder + 1,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Isactive:    true,
	}
	if err := tx.Create(&newShiftSession).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (s *shiftsessionservice) UpdateShiftSession(id int, input request.ShiftSessionRequestUpdate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()
	updates := map[string]interface{}{}
	if input.SessionName != nil {
		updates["session_name"] = *input.SessionName
	}
	if input.ShiftID != nil {
		updates["shift_id"] = *input.ShiftID
		var lastSession model.ShiftSession
		err := tx.Where("shift_id = ?", input.ShiftID).Order("shift_order DESC").First(&lastSession).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		updates["shift_id"] = *input.ShiftID
		if lastSession.ShiftID != id {
			updates["shift_order"] = lastSession.ShiftOrder + 1
		}

	}
	if input.StartTime != nil {
		updates["start_time"] = *input.StartTime
	}
	if input.EndTime != nil {
		updates["end_time"] = *input.EndTime
	}
	result := s.db.Model(&model.ShiftSession{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("shift session not found or already up to date")
	}
	return nil
}

func (s *shiftsessionservice) ChangeStatusShiftSession(id int) error {
	result := s.db.Model(&model.ShiftSession{}).Where("id = ?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("shift not found")
	}
	return nil
}

func (s *shiftsessionservice) GetShiftSessionV2(id int) (response.ShiftSessionResponsev2, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return response.ShiftSessionResponsev2{}, tx.Error
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
		return response.ShiftSessionResponsev2{}, err
	}

	var shiftpattern model.ShiftPattern
	if err := tx.Where("employee_id = ? AND day_of_week_id = ?", user.EmployeeID, dayOfWeek).
		First(&shiftpattern).Error; err != nil {
		tx.Rollback()
		return response.ShiftSessionResponsev2{}, err
	}

	if shiftpattern.Isdayoff {
		tx.Commit()
		return response.ShiftSessionResponsev2{ShowCheckIn: false, ShowCheckOut: false}, nil
	}

	var attendancelog model.AttendanceLog
	attendanceErr := tx.Where("employee_id = ? AND check_date = ?", user.EmployeeID, currentDate).
		Order("id DESC").
		First(&attendancelog).Error

	var session model.ShiftSession
	var showCheckIn, showCheckOut bool

	if attendanceErr != nil {
		// No log today — show check-in for first session
		if err := tx.Where("shift_id = ?", shiftpattern.ShiftID).
			Order("shift_order ASC").
			First(&session).Error; err != nil {
			tx.Rollback()
			return response.ShiftSessionResponsev2{}, err
		}
		showCheckIn = true
		showCheckOut = false

	} else if attendancelog.StatusAttendanceLogID == 1 {
		// Checked in but not yet checked out
		if err := tx.Where("shift_id = ? AND shift_order = ?", shiftpattern.ShiftID, attendancelog.ShiftSessionOrder).
			First(&session).Error; err != nil {
			tx.Rollback()
			return response.ShiftSessionResponsev2{}, err
		}
		showCheckIn = false
		showCheckOut = true

	} else {
		// Last session completed — advance to next session
		nextOrder := attendancelog.ShiftSessionOrder + 1
		if err := tx.Where("shift_id = ? AND shift_order = ?", shiftpattern.ShiftID, nextOrder).
			First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Commit()
				return response.ShiftSessionResponsev2{ShowCheckIn: false, ShowCheckOut: false}, nil
			}
			tx.Rollback()
			return response.ShiftSessionResponsev2{}, err
		}
		showCheckIn = true
		showCheckOut = false
	}

	// ✅ Parse times AFTER session is populated from DB
	startTime, _ := time.Parse("15:04:05", session.StartTime)
	endTime, _ := time.Parse("15:04:05", session.EndTime)

	isLate := 0
	isLeftEarly := 0

	if showCheckIn {
		// Late if checking in after the session start time
		if now.Hour() > startTime.Hour() ||
			(now.Hour() == startTime.Hour() && now.Minute() > startTime.Minute()) {
			isLate = 1
		}
	}

	if showCheckOut {
		// Left early if checking out before the session end time
		shiftEnd := time.Date(now.Year(), now.Month(), now.Day(),
			endTime.Hour(), endTime.Minute(), endTime.Second(), 0, now.Location())
		if now.Before(shiftEnd) {
			isLeftEarly = 1
		}
	}

	if err := tx.Commit().Error; err != nil {
		return response.ShiftSessionResponsev2{}, err
	}

	return response.ShiftSessionResponsev2{
		StartTime:    session.StartTime,
		EndTime:      session.EndTime,
		ShowCheckIn:  showCheckIn,
		ShowCheckOut: showCheckOut,
		IsLate:       isLate,
		IsLeftEarly:  isLeftEarly,
	}, nil
}
