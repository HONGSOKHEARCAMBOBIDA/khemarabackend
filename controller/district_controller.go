package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DistrictController struct {
	service service.DistrictService
}

func NewDistrictController() *DistrictController {
	return &DistrictController{
		service: service.NewDistrictService(),
	}
}

func (cr *DistrictController) GetDistrict(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	district, err := cr.service.GetDistrict(id)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, district)
}
