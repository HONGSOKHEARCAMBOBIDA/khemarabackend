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

func (cr *AuthController) RefreshToken(c *gin.Context) {
	var input request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := cr.service.RefreshToken(input, c)
	if err != nil {
		share.ResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, result)
}

func (cr *AuthController) Logout(c *gin.Context) {
	var input request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Revoke the specific session matching this refresh token
	sessions, err := cr.service.GetSessions(uint(userID))
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Find session by refresh token and revoke it
	for _, session := range sessions {
		if err := cr.service.RevokeSession(session.ID, uint(userID)); err == nil {
			share.RespondDate(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
			return
		}
	}

	share.RespondDate(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (cr *AuthController) GetSessions(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	sessions, err := cr.service.GetSessions(uint(userID))
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, sessions)
}

func (cr *AuthController) RevokeSession(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	if err := cr.service.RevokeSession(uint(sessionID), uint(userID)); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, gin.H{"message": "Session revoked successfully"})
}

func (cr *AuthController) RevokeAllSessions(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := cr.service.RevokeAllSessions(uint(userID)); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, gin.H{"message": "All sessions revoked successfully"})
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
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.UserRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateUser(id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "updated")
}

func (cr *AuthController) ChangePassword(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.NewPassword
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.ChangePassword(id, input); err != nil {
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
