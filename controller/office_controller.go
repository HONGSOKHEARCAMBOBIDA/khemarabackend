package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OfficeController struct {
	service service.OfficeService
}

func NewOfficeController() *OfficeController {
	return &OfficeController{
		service: service.NewOfficeService(),
	}
}

func (cr *OfficeController) GetAllOffice(c *gin.Context) {
	office, err := cr.service.GetAllOffice()
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, office)
}
