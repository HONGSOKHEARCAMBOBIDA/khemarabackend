package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct {
	service service.DepartmentService
}

func NewDepartmentController() *DepartmentController {
	return &DepartmentController{
		service: service.NewDepartmentService(),
	}
}

func (cr *DepartmentController) GetDepartment(c *gin.Context) {
	data, err := cr.service.GetDepartment()
	if err != nil {
		share.ResponseError(c, http.StatusNoContent, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, data)
}

func (cr *DepartmentController) CreateDepartment(c *gin.Context) {
	var input request.DepartmentRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateDepartment(input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "department created")
}

func (cr *DepartmentController) UpdateDepartment(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.DepartmentRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateDepartment(id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "department updated")
}

func (cr *DepartmentController) ChangeStatusDepartment(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.ChangeStatusDepartment(id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "department status change")
}
