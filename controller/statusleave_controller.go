package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusLeaveController struct {
	service service.StatusLeaveService
}

func NewStatusLeaveController() StatusLeaveController {
	return StatusLeaveController{
		service: service.NewStatusLeaveService(),
	}
}

func (cr *StatusLeaveController) GetStatusLeave(c *gin.Context) {
	data, err := cr.service.GetStatusLeave()
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
