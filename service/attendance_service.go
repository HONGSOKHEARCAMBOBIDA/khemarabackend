package service

import (
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AttendanceService interface {
	CheckIn(Id int, input request.LocationRequest, branchId int) error
}

type attendanceservice struct {
	db *gorm.DB
}

func NewAttendanceService() AttendanceService {
	return &attendanceservice{
		db: config.DB,
	}
}

func (s *attendanceservice) CheckIn(id int, input request.LocationRequest, branchId int) error {
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

	var shiftpattern model.ShiftPattern
	if err := tx.Where("employee_id = ? AND day_of_week_id = ?", id, dayOfWeek).
		First(&shiftpattern).Error; err != nil {
		tx.Rollback()
		return err
	}

	if shiftpattern.Isdayoff {
		tx.Rollback()
		return nil
	}

	var branch model.Branch
	if err := tx.First(&branch, branchId).Error; err != nil {
		tx.Rollback()
		return err
	}

	var attendancelog model.AttendanceLog
	err := tx.Where("employee_id = ? AND status_attendance_log_id = ?", id, 1).
		Order("id DESC").
		First(&attendancelog).Error

	if err != nil {
		// First check-in → get first session
		var session model.ShiftSession
		if err := tx.Where("shift_id = ?", shiftpattern.ShiftID).
			Order("shift_order ASC").
			First(&session).Error; err != nil {
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
		inzone := distance <= float64(branch.Radius)
		currentDate := time.Now().Format("2006-01-02")
		now := time.Now()
		start_time, _ := time.Parse("15:04:05", session.StartTime)
		is_late := 0
		if now.Hour() > start_time.Hour() || (now.Hour() == start_time.Hour() && now.Minute() > starttime.Minute()) {
			is_late = 1
		}

		newlog := model.AttendanceLog{
			EmployeeID:            id,
			CheckDate:             currentDate,
			Note:                  "",
			BranchID:              branchId,
			StatusAttendanceLogID: 1,
			ShiftSessionOrder:     session.ShiftOrder,
		}

		if err := tx.Create(&newlog).Error; err != nil {
			tx.Rollback()
			return err
		}

		newlogrecore := model.AttendanceRecord{
			AttendanceLogID: newlog.ID,
			ShiftSessionID:  session.ID,
			CheckTime:       now,
			IsLate:          is_late,
			IsLeftEarly:     nil,
			Latitude:        lat,
			Logitude:        log,
			Note:            input.Note,
			Iszoone:         inzone,
		}

		if err := tx.Create(&newlogrecore).Error; err != nil {
			tx.Rollback()
			return err
		}

		// TODO: create attendance log here

	} else {
		// Next session
		nextOrder := attendancelog.ShiftSessionOrder + 1

		var session model.ShiftSession
		if err := tx.Where("shift_id = ? AND shift_order = ?", shiftpattern.ShiftID, nextOrder).
			First(&session).Error; err != nil {
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
		inzone := distance <= float64(branch.Radius)
		currentDate := time.Now().Format("2006-01-02")
		now := time.Now()
		start_time, _ := time.Parse("15:04:05", session.StartTime)
		is_late := 0
		if now.Hour() > start_time.Hour() || (now.Hour() == start_time.Hour() && now.Minute() > starttime.Minute()) {
			is_late = 1
		}
		newlog := model.AttendanceLog{
			EmployeeID:            id,
			CheckDate:             currentDate,
			Note:                  "",
			BranchID:              branchId,
			StatusAttendanceLogID: 1,
			ShiftSessionOrder:     session.ShiftOrder,
		}
		if err := tx.Create(&newlog).Error; err != nil {
			tx.Rollback()
			return err
		}
		newlogrecore := model.AttendanceRecord{
			AttendanceLogID: newlog.ID,
			ShiftSessionID:  session.ID,
			CheckTime:       now,
			IsLate:          is_late,
			IsLeftEarly:     nil,
			Latitude:        lat,
			Logitude:        log,
			Note:            input.Note,
			Iszoone:         inzone,
		}

		if err := tx.Create(&newlogrecore).Error; err != nil {
			tx.Rollback()
			return err
		}
		// TODO: create next session log
	}

	return tx.Commit().Error
}
