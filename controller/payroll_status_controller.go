package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PayrollStatusController struct {
	service service.PayrollStatusService
}

func NewPayrollStatusController() PayrollStatusController {
	return PayrollStatusController{
		service: service.NewPayrollStatusService(),
	}
}

func (cr *PayrollStatusController) GetPayrollStatus(c *gin.Context) {
	data, err := cr.service.GetPayrollStatus()
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
