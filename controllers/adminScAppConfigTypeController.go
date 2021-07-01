package controllers

import (
	"api/lib"
	"api/models"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func AdminGetListDropdownScAppConfigType(c echo.Context) error {

	var err error
	var status int

	var configType []models.ListDropdownScAppConfigType

	status, err = models.AdminGetListDropdownScAppConfigType(&configType)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(configType) < 1 {
		log.Error("Config Type not found")
		return lib.CustomError(http.StatusNotFound, "Config Type not found", "Config Type not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = configType

	return c.JSON(http.StatusOK, response)
}
