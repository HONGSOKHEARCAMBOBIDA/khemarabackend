package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

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
	if err := cr.service.CheckIn(input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "check-in success")
}
