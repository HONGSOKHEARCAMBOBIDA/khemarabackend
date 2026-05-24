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

type AuthController struct {
	service service.AuthService
}

func NewAuthController() AuthController {
	return AuthController{
		service: service.NewAuthService(),
	}
}

func (cr *AuthController) Login(c *gin.Context) {
	var input request.AuthRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, 400, err.Error())
		return
	}
	result, err := cr.service.Login(input, c)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	//share.ResponeSuccess(c, 200, result.Token)
	share.RespondDate(c, http.StatusOK, result)
}

func (cr *AuthController) Register(c *gin.Context) {
	var input request.RegisterRequest
	if err := c.ShouldBind(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	id, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "Please Login")
		return
	}
	if err := cr.service.Register(id, input, c); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusCreated, "user create")
}

func (cr *AuthController) GetUserByBranch(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := cr.service.GetUserByBranch(id)
	if err != nil {
		share.ResponseError(c, http.StatusNoContent, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) UpdateUser(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusBadRequest, "please login")
		return
	}
	var input request.UserRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateUser(userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "updated")
}

func (cr *AuthController) ChangePassword(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusBadRequest, "please login")
		return
	}
	var input request.NewPassword
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.ChangePassword(userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "password changed")
}

func (cr *AuthController) GetUserByID(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := cr.service.GetUserByID(id)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
