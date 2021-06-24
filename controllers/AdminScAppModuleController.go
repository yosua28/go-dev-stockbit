package controllers

import (
	"api/lib"
	"api/models"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func AdminGetListScAppModuleDropdown(c echo.Context) error {

	var appModule []models.ListAppModuleDropdown

	status, err := models.AdminGetListAppModuleDropdown(&appModule)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(appModule) < 1 {
		log.Error("App Module not found")
		return lib.CustomError(http.StatusNotFound, "App Module not found", "App Module not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = appModule

	return c.JSON(http.StatusOK, response)
}
