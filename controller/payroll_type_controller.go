package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PayRollTypeController struct {
	service service.PayRollTypeService
}

func NewPayRollTypeController() PayRollTypeController {
	return PayRollTypeController{
		service: service.NewPayRollTypeService(),
	}
}

func (cr *PayRollTypeController) GetPayrollType(c *gin.Context) {
	data, err := cr.service.GetPayrollType()
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
