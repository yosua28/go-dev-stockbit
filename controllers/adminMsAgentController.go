package controllers

import (
	"api/lib"
	"api/models"
	"net/http"
	"strconv"

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

func GetListAgentLastBranch(c echo.Context) error {

	var err error
	var status int

	keyStr := c.Param("branch_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var agent []models.MsAgentLastBranch

	status, err = models.GetMsAgentLastBranch(&agent, keyStr)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(agent) < 1 {
		log.Error("Agent not found")
		return lib.CustomError(http.StatusNotFound, "Agent not found", "Agent not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = agent

	return c.JSON(http.StatusOK, response)
}
