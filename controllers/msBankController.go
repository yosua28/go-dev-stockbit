package controllers

import (
	"api/models"
	"api/lib"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

func GetMsBankList(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)
	params["orderBy"] = "bank_name"
	params["orderType"] = "ASC"
	params["rec_status"] = "1"
	var bankDB []models.MsBank
	status, err = models.GetAllMsBank(&bankDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(bankDB) < 1 {
		log.Error("Data not found")
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}
	var responseData []models.MsBankList
	
	for _, bank := range bankDB {
		var data models.MsBankList

		data.BankKey = bank.BankKey
		data.BankCode = bank.BankCode
		data.BankName = bank.BankName
		if bank.BankFullname != nil {
			data.BankFullname = *bank.BankFullname
		}
		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	
	return c.JSON(http.StatusOK, response)
}