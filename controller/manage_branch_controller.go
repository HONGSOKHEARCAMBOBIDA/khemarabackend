package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ManageBranchController struct {
	service service.ManageBranchService
}

func NewManageBranchController() *ManageBranchController {
	return &ManageBranchController{
		service: service.NewManageBranchService(),
	}
}

func (cr ManageBranchController) GetManageBranch(c *gin.Context) {
	data, err := cr.service.GetManageBranch()
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
