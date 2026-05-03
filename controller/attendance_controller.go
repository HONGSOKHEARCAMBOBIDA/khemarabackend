package controller

import (
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	service service.AttendanceService
}

func NewAttendanceController() AttendanceController {
	return AttendanceController{
		service: service.NewAttendanceService(),
	}
}

func (cr *AttendanceController) CheckIn(c *gin.Context) {
	var input request.LocationRequest
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "Please Login")
		return
	}
	if err := cr.service.CheckIn(userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "check-in success")
}

func (cr *AttendanceController) CheckOut(c *gin.Context) {
	var input request.LocationRequest
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "Please Login")
		return
	}
	if err := cr.service.CheckOut(userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "check-out success")
}

func (cr *AttendanceController) GetAttendance(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	filter := map[string]string{
		"name":               c.Query("name"),
		"department_id":      c.Query("department_id"),
		"employee_id":        c.Query("employee_id"),
		"office_id":          c.Query("office_id"),
		"check_in_early":     c.Query("check_in_early"),
		"check_in_on_time":   c.Query("check_in_on_time"),
		"is_late":            c.Query("is_late"),
		"is_left_early":      c.Query("is_left_early"),
		"check_out_on_time":  c.Query("check_out_on_time"),
		"check_out_overtime": c.Query("check_out_overtime"),
		"check_date_from":    c.Query("check_date_from"),
		"check_date_to":      c.Query("check_date_to"),
	}
	attendance, metadata, err := cr.service.GetAttendance(filter, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       attendance,
		"pagination": metadata,
	})
}
