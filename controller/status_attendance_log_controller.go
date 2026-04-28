package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusAttendanceLogController struct {
	service service.StatusAttendanceLogService
}

func NewStatusAttendanceLogController() StatusAttendanceLogController {
	return StatusAttendanceLogController{
		service: service.NewStatusAttendanceLogService(),
	}
}

func (cr *StatusAttendanceLogController) GetStatusAttendanceLogService(c *gin.Context) {
	data, err := cr.service.GetStatusAttendanceLogService()
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
