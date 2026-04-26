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

type EmployeeController struct {
	service service.EmployeeService
}

func NewEmployeeController() EmployeeController {
	return EmployeeController{
		service: service.NewEmployeeService(),
	}
}

func (cr *EmployeeController) GetEmployee(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	filter := map[string]string{
		"name":          c.Query("name"),
		"branch_id":     c.Query("branch_id"),
		"department_id": c.Query("department_id"),
		"position_id":   c.Query("position_id"),
		"office_id":     c.Query("office_id"),
		"is_promote":    c.Query("is_promote"),
	}
	employee, metadata, err := cr.service.GetEmployee(filter, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       employee,
		"pagination": metadata,
	})
}

func (cr *EmployeeController) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var input request.EmployeeEmpoyeeProfileRequestUpdate
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateEmployee(id, input, c, userID); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "update success")
}

func (cr *EmployeeController) UpdateEmployeeEducation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.EmployeeEducationRequestUpdate
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateEmployeeEducation(id, input, c); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "updated")
}

func (cr *EmployeeController) CreateEmployeeEducation(c *gin.Context) {
	var input request.EmployeeEducationRequestCreate
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateEmployeeEducation(input, c); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "created")
}

func (cr *EmployeeController) UpdateEmployeeWorkExperience(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.EmployeeWorkExperienceRequestUpdate
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateEmployeeWorkExperience(id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Update")
}
