package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeaveTypeController struct {
	service service.LeaveTypeService
}

func NewLeaveTypeController() LeaveTypeController {
	return LeaveTypeController{
		service: service.NewLeaveTypeService(),
	}
}

func (cr *LeaveTypeController) GetLeaveType(c *gin.Context) {
	data, err := cr.service.GetLeaveType()
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
