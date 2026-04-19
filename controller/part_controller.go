package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PartController struct {
	service service.PartService
}

func NewPartController() PartController {
	return PartController{
		service: service.NewPartService(),
	}
}

func (cr *PartController) GetPart(c *gin.Context) {
	data, err := cr.service.GetPart()
	if err != nil {
		share.ResponseError(c, http.StatusNoContent, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
