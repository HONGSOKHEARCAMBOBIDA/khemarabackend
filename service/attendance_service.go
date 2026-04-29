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
	CheckIn(input request.LocationRequest) error
	CheckOut(input request.LocationRequest) error
}

type attendanceservice struct {
	db *gorm.DB
}

func NewAttendanceService() AttendanceService {
	return &attendanceservice{
		db: config.DB,
	}
}

func (s *attendanceservice) CheckIn(input request.LocationRequest) error {
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
	if err := tx.Where("employee_id = ? AND day_of_week_id = ?", input.EmployeeID, dayOfWeek).
		First(&shiftpattern).Error; err != nil {
		tx.Rollback()
		return err
	}

	if shiftpattern.Isdayoff {
		tx.Rollback()
		return nil
	}

	var attendancelog model.AttendanceLog
	err := tx.Where("employee_id = ? AND check_date =?", input.EmployeeID, now.Format("2006-01-02")).
		Order("id DESC").
		First(&attendancelog).Error

	if err != nil {
		var session model.ShiftSession
		if err := tx.Where("shift_id = ?", shiftpattern.ShiftID).
			Order("shift_order ASC").
			First(&session).Error; err != nil {
			tx.Rollback()
			return err
		}
		lat, err := strconv.ParseFloat(input.BranchLatitude, 64)
		if err != nil {
			tx.Rollback()
			return err
		}
		log, err := strconv.ParseFloat(input.BranchLongitude, 64)
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
		currentDate := time.Now().Format("2006-01-02")
		now := time.Now()
		start_time, _ := time.Parse("15:04:05", session.StartTime)
		is_late := 0
		if now.Hour() > start_time.Hour() || (now.Hour() == start_time.Hour() && now.Minute() > start_time.Minute()) {
			is_late = 1
		}

		newlog := model.AttendanceLog{
			EmployeeID:            input.EmployeeID,
			CheckDate:             currentDate,
			Note:                  "",
			BranchID:              input.BranchID,
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
			IsLate:          &is_late,
			IsLeftEarly:     nil,
			Latitude:        lat,
			Logitude:        log,
			Note:            input.Note,
			Iszoone:         inzone,
			Type:            "IN",
		}

		if err := tx.Create(&newlogrecore).Error; err != nil {
			tx.Rollback()
			return err
		}

	} else {

		nextOrder := attendancelog.ShiftSessionOrder + 1

		var session model.ShiftSession
		if err := tx.Where("shift_id = ? AND shift_order = ?", shiftpattern.ShiftID, nextOrder).
			First(&session).Error; err != nil {
			tx.Rollback()
			return err
		}

		lat, err := strconv.ParseFloat(input.BranchLatitude, 64)
		if err != nil {
			tx.Rollback()
			return err
		}
		log, err := strconv.ParseFloat(input.BranchLongitude, 64)
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
		start_time, _ := time.Parse("15:04:05", session.StartTime)
		currentDate := time.Now().Format("2006-01-02")
		is_late := 0
		if now.Hour() > start_time.Hour() || (now.Hour() == start_time.Hour() && now.Minute() > start_time.Minute()) {
			is_late = 1
		}
		newattendancelog := model.AttendanceLog{
			EmployeeID:            input.EmployeeID,
			CheckDate:             currentDate,
			Note:                  "",
			BranchID:              input.BranchID,
			StatusAttendanceLogID: 1,
			ShiftSessionOrder:     nextOrder,
		}
		if err := tx.Create(&newattendancelog).Error; err != nil {
			tx.Rollback()
			return err
		}
		newlogrecore := model.AttendanceRecord{
			AttendanceLogID: newattendancelog.ID,
			ShiftSessionID:  session.ID,
			CheckTime:       now,
			IsLate:          &is_late,
			IsLeftEarly:     nil,
			Latitude:        lat,
			Logitude:        log,
			Note:            input.Note,
			Iszoone:         inzone,
			Type:            "IN",
		}

		if err := tx.Create(&newlogrecore).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&model.AttendanceLog{}).Where("id =?", attendancelog.ID).Update("shift_session_order", nextOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *attendanceservice) CheckOut(input request.LocationRequest) error {
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
	if err := tx.Where("employee_id = ? AND day_of_week_id = ?", input.EmployeeID, dayOfWeek).
		First(&shiftpattern).Error; err != nil {
		tx.Rollback()
		return err
	}

	if shiftpattern.Isdayoff {
		tx.Rollback()
		return nil
	}

	var attendancelog model.AttendanceLog
	err := tx.Where("employee_id =? AND status_attendance_log_id =?", input.EmployeeID, 1).Order("id DESC").First(&attendancelog).Error
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
		lat, err := strconv.ParseFloat(input.BranchLatitude, 64)
		if err != nil {
			tx.Rollback()
			return err
		}
		log, err := strconv.ParseFloat(input.BranchLongitude, 64)
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
		endTime, _ := time.Parse("15:04:05", session.EndTime)
		currentDate := time.Now().Format("2006-01-02")
		shiftEnd := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), 0, now.Location())
		isLeftEarly := 0
		if now.Before(shiftEnd) {
			isLeftEarly = 1
		}
		attendancelog.StatusAttendanceLogID = 2
		attendancelog.CheckDate = currentDate
		if err := tx.Save(&attendancelog).Error; err != nil {
			tx.Rollback()
			return err
		}
		newattendancerecord := model.AttendanceRecord{
			AttendanceLogID: attendancelog.ID,
			ShiftSessionID:  session.ID,
			CheckTime:       now,
			IsLate:          nil,
			IsLeftEarly:     &isLeftEarly,
			Latitude:        lat,
			Logitude:        log,
			Note:            input.Note,
			Iszoone:         inzone,
			Type:            "OUT",
		}
		if err := tx.Create(&newattendancerecord).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
