package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeTypeController struct {
	service service.EmployeeTypeService
}

func NewEmployeeTypeController() EmployeeTypeController {
	return EmployeeTypeController{
		service: service.NewEmployeeTypeService(),
	}
}

func (cr *EmployeeTypeController) GetEmployeeType(c *gin.Context) {
	data, err := cr.service.GetEmployeeType()
	if err != nil {
		share.ResponseError(c, http.StatusNoContent, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr EmployeeTypeController) CreateEmployeeType(c *gin.Context) {
	var input request.EmployeeTypeRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateEmployeeType(input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "employee type create")
}

func (cr *EmployeeTypeController) UpdateEmployeeType(c *gin.Context) {
	idparam := c.Param("id")

	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.EmployeeTypeRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateEmployeeType(id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "employee type updated")
}

func (cr *EmployeeTypeController) ChangeStatusEmployeeType(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.ChangeStatusEmployeeType(id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "employee has change status")
}
