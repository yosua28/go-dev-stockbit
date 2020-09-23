package controllers

import (
	"api/models"
	"api/lib"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

func GetMsCountryList(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)
	params["orderBy"] = "cou_name"
	params["orderType"] = "ASC"
	params["rec_status"] = "1"
	var countryDB []models.MsCountry
	status, err = models.GetAllMsCountry(&countryDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(countryDB) < 1 {
		log.Error("Data not found")
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}
	var responseData []models.MsCountryList
	
	for _, country := range countryDB {
		var data models.MsCountryList

		data.CountryKey = country.CountryKey
		data.CouCode = country.CouCode
		data.CouName = country.CouName
		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	
	return c.JSON(http.StatusOK, response)
}