package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VillageController struct {
	service service.VillageService
}

func NewVillageController() VillageController {
	return VillageController{
		service: service.NewVillageService(),
	}
}

func (cr VillageController) GetVillage(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	village, err := cr.service.GetVillage(id)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, village)
}
