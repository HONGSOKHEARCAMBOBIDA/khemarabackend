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

type PayrollController struct {
	service service.PayrollService
}

func NewPayrollController() PayrollController {
	return PayrollController{
		service: service.NewPayrollService(),
	}
}

func (cr *PayrollController) CreatePayroll(c *gin.Context) {
	var input []request.PayrollRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "Please Login")
		return
	}
	if err := cr.service.CreatePayroll(userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "payroll created success")
}

func (cr *PayrollController) DeletePayroll(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.DeletePayroll(id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "payroll deleted")
}

func (cr *PayrollController) GetDraftPayroll(c *gin.Context) {
	currencyParam := c.Query("currency")
	branchParam := c.Query("branch")
	payrolltype := c.Query("payroll")
	currencyID, err := strconv.Atoi(currencyParam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, "invalid currency id")
		return
	}

	branchID, err := strconv.Atoi(branchParam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, "invalid branch id")
		return
	}

	payrolltypeID, err := strconv.Atoi(payrolltype)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, "invalid branch id")
		return
	}

	data, err := cr.service.GetDraftPayroll(branchID, currencyID, payrolltypeID)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, data)
}

func (cr *PayrollController) GetPayroll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	filters := map[string]string{
		"branch_id":     c.Query("branch_id"),
		"name":          c.Query("name"),
		"position_id":   c.Query("position_id"),
		"status_id":     c.Query("status_id"),
		"office_id":     c.Query("office_id"),
		"department_id": c.Query("department_id"),
		"start_date":    c.Query("start_date"),
		"end_date":      c.Query("end_date"),
	}
	payroll, metadata, err := cr.service.GetPayroll(userID, filters, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       payroll,
		"pagination": metadata,
	})
}
