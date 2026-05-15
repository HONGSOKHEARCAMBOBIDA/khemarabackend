package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BonusTypeController struct {
	service service.BonusTypeService
}

func NewBonusTypeController() BonusTypeController {
	return BonusTypeController{
		service: service.NewBonusTypeService(),
	}
}

func (cr *BonusTypeController) GetBonusType(c *gin.Context) {
	data, err := cr.service.GetBonusType()
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
