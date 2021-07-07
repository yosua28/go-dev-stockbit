package controllers

import (
	"api/lib"
	"api/models"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func AdminGetListDropdownMsParticipant(c echo.Context) error {

	var err error
	var status int

	var participant []models.ListDropdownMsParticipant

	status, err = models.AdminGetListDropdownMsParticipant(&participant)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(participant) < 1 {
		log.Error("Participant Type not found")
		return lib.CustomError(http.StatusNotFound, "Participant Type not found", "Participant Type not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = participant

	return c.JSON(http.StatusOK, response)
}
