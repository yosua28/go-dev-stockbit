package controllers

import (
	"api/lib"
	"api/models"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetListAgentDropdown(c echo.Context) error {

	var err error
	var status int

	var msAgent []models.MsAgentDropdown

	status, err = models.GetMsAgentDropdown(&msAgent)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msAgent) < 1 {
		log.Error("Agent not found")
		return lib.CustomError(http.StatusNotFound, "Agent not found", "Agent not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = msAgent

	return c.JSON(http.StatusOK, response)
}
