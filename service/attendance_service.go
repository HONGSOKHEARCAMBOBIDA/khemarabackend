package service

import (
	"fmt"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AttendanceService interface {
	CheckIn(id int, input request.LocationRequest) error
	CheckOut(id int, input request.LocationRequest) error
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

	startTime, _ := time.Parse("15:04:05", session.StartTime)
	isLate := 0
	if now.Hour() > startTime.Hour() ||
		(now.Hour() == startTime.Hour() && now.Minute() > startTime.Minute()) {
		isLate = 1
	}

	newLog := model.AttendanceLog{
		EmployeeID:            user.EmployeeID,
		CheckDate:             currentDate,
		Note:                  "",
		BranchID:              input.BranchID,
		StatusAttendanceLogID: 1,
		ShiftSessionOrder:     shiftOrder,
	}

	if err := tx.Create(&newLog).Error; err != nil {
		tx.Rollback()
		return err
	}

	newRecord := model.AttendanceRecord{
		AttendanceLogID: newLog.ID,
		ShiftSessionID:  session.ID,
		CheckTime:       now,
		IsLate:          &isLate,
		IsLeftEarly:     nil,
		Latitude:        userLat,
		Logitude:        userLng,
		Note:            input.Note,
		Iszoone:         inzone,
		Type:            "IN",
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
			Latitude:        latput,
			Logitude:        logput,
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
