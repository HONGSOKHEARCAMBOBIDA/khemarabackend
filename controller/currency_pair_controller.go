package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CurrencyPairController struct {
	service service.CurrencyPairService
}

func NewCurrencyPairController() CurrencyPairController {
	return CurrencyPairController{
		service: service.NewCurrencyPairService(),
	}
}

func (cpc CurrencyPairController) CreateCurrencyPair(c *gin.Context) {
	var input request.CurrencyPairRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cpc.service.CreateCurrencyPair(input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Currency pair created")
}

func (cpc CurrencyPairController) GetCurrencypair(c *gin.Context) {
	currencypair, err := cpc.service.GetCurrencypair()
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, currencypair)

}

func (cpc CurrencyPairController) UpdateCurrencyPaire(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var input request.CurrencyPairRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cpc.service.UpdateCurrencyPaire(id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Update Success")
}

func (cpc CurrencyPairController) ChangeStatusCurrencyPair(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := cpc.service.ChangeStatusCurrencyPair(id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "change status success")
}
