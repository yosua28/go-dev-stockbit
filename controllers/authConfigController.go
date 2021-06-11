package controllers

import (
	"api/lib"
	"api/models"
	_ "encoding/base64"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetUserConfig(c echo.Context) error {

	responseData := make(map[string]interface{})
	var scApp models.ScAppConfig
	status, err := models.GetScAppConfigByCode(&scApp, "IDLE_TIME_MOBILE")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data IDLE_TIME_MOBILE")
	}
	idle, _ := strconv.ParseUint(*scApp.AppConfigValue, 10, 64)
	if idle == 0 {
		return lib.CustomError(http.StatusNotFound)
	}
	responseData["idle_time"] = idle

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}
