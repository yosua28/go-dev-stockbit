package controllers

import (
	"api/lib"
	"api/models"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetListBranchDropdown(c echo.Context) error {

	var err error
	var status int

	var msBranch []models.MsBranchDropdown

	status, err = models.GetMsBranchDropdown(&msBranch)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msBranch) < 1 {
		log.Error("Branch not found")
		return lib.CustomError(http.StatusNotFound, "Branch not found", "Branch not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = msBranch

	return c.JSON(http.StatusOK, response)
}
