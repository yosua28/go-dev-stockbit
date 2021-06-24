package controllers

import (
	"api/lib"
	"api/models"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func AdminGetListScMenuTypeDropdown(c echo.Context) error {

	var menuType []models.ListMenuTypeDropdown

	status, err := models.AdminGetListMenuTypeDropdown(&menuType)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(menuType) < 1 {
		log.Error("Menu Type not found")
		return lib.CustomError(http.StatusNotFound, "Menu Type not found", "Menu Type not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = menuType

	return c.JSON(http.StatusOK, response)
}
