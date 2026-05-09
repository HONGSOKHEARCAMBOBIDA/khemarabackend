package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeaveDurationController struct {
	service service.LeaveDurationService
}

func NewLeaveDurationController() LeaveDurationController {
	return LeaveDurationController{
		service: service.NewLeaveDurationService(),
	}
}

func (cr *LeaveDurationController) GetLeaveDuration(c *gin.Context) {
	data, err := cr.service.GetLeaveDuration()
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
