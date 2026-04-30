package controller

import (
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShiftSessionController struct {
	service service.ShiftSessionService
}

func NewShiftSessionController() ShiftSessionController {
	return ShiftSessionController{
		service: service.Newshiftsessionservice(),
	}
}

func (cr *ShiftSessionController) GetAllShiftSession(c *gin.Context) {
	data, err := cr.service.GetAllShiftSession()
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *ShiftSessionController) GetByShiftID(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := cr.service.GetByShiftID(id)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *ShiftSessionController) CreateShiftSession(c *gin.Context) {
	var input request.ShiftSessionRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateShiftSession(input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusCreated, "shift session created")
}

func (cr *ShiftSessionController) UpdateShiftSession(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.ShiftSessionRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateShiftSession(id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "shift session updated")
}

func (cr *ShiftSessionController) ChangeStatusShiftSession(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.ChangeStatusShiftSession(id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "shift session update status")
}

func (cr *ShiftSessionController) GetShiftSessionV2(c *gin.Context) {
	id, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	data, err := cr.service.GetShiftSessionV2(id)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
