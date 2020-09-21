package controllers

import (
	"api/models"
	"api/lib"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

func GetMsCityList(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)
	field := c.Param("field")
	if field == "" {
		log.Error("Missing required parameters")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameters","Missing required parameters")
	}
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	params[field] = keyStr
	params["orderBy"] = "city_name"
	params["orderType"] = "ASC"
	params["rec_status"] = "1"
	var cityDB []models.MsCity
	status, err = models.GetAllMsCity(&cityDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(cityDB) < 1 {
		log.Error("Data not found")
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}
	var responseData []models.MsCityList
	
	for _, city := range cityDB {
		var data models.MsCityList

		data.CityKey = city.CityKey
		if city.ParentKey != nil {
			data.ParentKey = *city.ParentKey
		}
		data.CityCode = city.CityCode
		data.CityName = city.CityName
		data.CityLevel = city.CityLevel
		if city.PostalCode != nil {
			data.PostalCode = *city.PostalCode
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