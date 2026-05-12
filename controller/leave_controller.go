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

type LeaveController struct {
	service service.LeaveService
}

func NewLeaveController() LeaveController {
	return LeaveController{
		service: service.NewLeaveService(),
	}
}

func (cr *LeaveController) CreateLeave(c *gin.Context) {
	var input request.LeaveCreate
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "login please")
		return
	}
	if err := cr.service.CreateLeave(userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "leave created")
}

func (cr *LeaveController) GetLeave(c *gin.Context) {
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
		"employee_id":     c.Query("employee_id"),
		"branch_id":       c.Query("branch_id"),
		"office_id":       c.Query("office_id"),
		"status_leave_id": c.Query("status_leave_id"),
		"leave_type_id":   c.Query("leave_type_id"),
		"start_date":      c.Query("start_date"),
		"end_date":        c.Query("end_date"),
		"search":          c.Query("search"),
	}

	leave, metadata, err := cr.service.GetLeave(userID, filter, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       leave,
		"pagination": metadata,
	})
}

func (cr *LeaveController) ApproveLeave(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.LeaveApproveRequest
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	if err := cr.service.ApproveLeave(id, input, userID); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "approved")
}

func (cr *LeaveController) UpdateLeave(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.LeaveUpdate
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateLeave(id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "leave updated")
}
