package controllers

import (
	"api/lib"
	"api/models"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetTransactionStatus(c echo.Context) error {

	var err error
	var status int

	params := make(map[string]string)
	params["rec_status"] = "1"

	var trTransactionStatus []models.TrTransactionStatus
	status, err = models.GetAllMsTransactionStatus(&trTransactionStatus, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(trTransactionStatus) < 1 {
		log.Error("transaction status not found")
		return lib.CustomError(http.StatusNotFound, "Transaction Status not found", "Transaction Status not found")
	}

	var responseData []models.TrTransactionStatusDropdown
	for _, tr := range trTransactionStatus {
		var data models.TrTransactionStatusDropdown
		data.TransStatusKey = tr.TransStatusKey
		data.StatusCode = tr.StatusCode

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
