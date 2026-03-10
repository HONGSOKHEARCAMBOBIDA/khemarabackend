package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PositionController struct {
	service service.PositionService
}

func NewPositionController() PositionController {
	return PositionController{
		service: service.NewPositionService(),
	}
}

func (cr *PositionController) GetAllPosition(c *gin.Context) {
	data, err := cr.service.GetAllPosition()
	if err != nil {
		share.ResponseError(c, http.StatusNoContent, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *PositionController) GetByDepartmentID(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := cr.service.GetByDepartmentID(id)
	if err != nil {
		share.ResponseError(c, http.StatusNoContent, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *PositionController) CreatePosition(c *gin.Context) {
	var input request.PositionRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreatePosition(input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "position created")
}

func (cr *PositionController) UpdatePosition(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.PositionRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdatePosition(id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "position updated")

}

func (cr *PositionController) ChangeStatusPosition(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.ChangeStatusPosition(id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "status change")
}
