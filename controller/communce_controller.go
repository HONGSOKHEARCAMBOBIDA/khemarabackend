package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommunceController struct {
	service service.CommunceService
}

func NewCommunceController() *CommunceController {
	return &CommunceController{
		service: service.NewCommunceService(),
	}
}

func (cr *CommunceController) GetCommunce(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := cr.service.GetCommunce(id)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
