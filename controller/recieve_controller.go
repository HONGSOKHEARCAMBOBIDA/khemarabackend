package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecieveController struct {
	service service.RecieveService
}

func NewRecieveController() RecieveController {
	return RecieveController{
		service: service.NewRecieveService(),
	}
}

func (cr *RecieveController) GetRecieve(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := cr.service.GetRecieve(id)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
