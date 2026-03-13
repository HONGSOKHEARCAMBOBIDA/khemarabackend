package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DayOfWeekController struct {
	service service.DayOfWeekService
}

func NewDayOfWeekController() *DayOfWeekController {
	return &DayOfWeekController{
		service: service.NewDayOfWeekService(),
	}
}

func (cr *DayOfWeekController) GetDayOfWeek(c *gin.Context) {
	data, err := cr.service.GetDayOfWeek()
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
