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

type LoanController struct {
	service service.LoanService
}

func NewLoanController() LoanController {
	return LoanController{
		service: service.NewLoanService(),
	}
}

func (cr *LoanController) CreateLoan(c *gin.Context) {
	var input request.LoanRequest
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "login please")
		return
	}
	if err := cr.service.CreateLoan(userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "loan created")
}

func (cr *LoanController) GetLoan(c *gin.Context) {
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
	filter := map[string]string{
		"search":      c.Query("search"),
		"employee_id": c.Query("employee_id"),
		"branch_id":   c.Query("branch_id"),
		"status":      c.Query("status"),
	}

	loan, metadata, err := cr.service.GetLoan(userID, filter, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       loan,
		"pagination": metadata,
	})

}
