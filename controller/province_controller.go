package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProvinceController struct {
	service service.ProvinceService
}

func NewProvinceController() *ProvinceController {
	return &ProvinceController{
		service: service.NewProvinceService(),
	}
}

func (cr *ProvinceController) GetProvince(c *gin.Context) {

	province, err := cr.service.GetProvince()
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, province)
}
