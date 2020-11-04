package controllers

import (
	"api/lib"
	"api/models"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetTransactionType(c echo.Context) error {

	var err error
	var status int

	params := make(map[string]string)
	params["type_domain"] = "FRONT"

	var trTransactionType []models.TrTransactionType
	status, err = models.GetAllMsTransactionTypeByCondition(&trTransactionType, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(trTransactionType) < 1 {
		log.Error("transaction type not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var responseData []models.TrTransactionTypeList
	for _, tr := range trTransactionType {
		var data models.TrTransactionTypeList
		data.TransTypeKey = tr.TransTypeKey
		data.TypeCode = tr.TypeCode
		data.TypeDescription = tr.TypeDescription

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
