package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHasPermissionController struct {
	service service.RoleHasPermissionService
}

func NewRoleHasPermissionController() RoleHasPermissionController {
	return RoleHasPermissionController{
		service: service.NewRoleHasPermissionService(),
	}
}

func (cr *RoleHasPermissionController) CreateRoleHasPermission(c *gin.Context) {
	var input request.RoleHasPermissionRequestCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateRoleHasPermission(input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.ResponseSuccess(c, http.StatusCreated, "Permission assigned to role")
}

func (cr *RoleHasPermissionController) DeleteRoleHasPermission(c *gin.Context) {
	var input request.RoleHasPermissionRequestDelete

	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.DeleteRoleHasPermission(input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.ResponseSuccess(c, http.StatusOK, "Permission removed from role")
}
