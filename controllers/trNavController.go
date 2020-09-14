package controllers

import (
	"api/models"
	"api/lib"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

func GetTrNavProduct(c echo.Context) error {
	var err error
	var status int
	
	productKeyStr := c.Param("product_key")
	productKey, _ := strconv.ParseUint(productKeyStr, 10, 64)
	if productKey == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	duration := c.Param("duration")
	if duration == "" {
		log.Error("Missing required parameters")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameters","Missing required parameters")
	}

	layout := "2006-01-02"
	now := time.Now()
	var date1, date2 string
	date1 = now.Format(layout + " 00:00:00")
	if duration == "d1" {
		date2 = now.Format(layout) + " 00:00:00"
	}else if duration == "m1" {
		date2 = now.AddDate(0,-1,0).Format(layout) + " 00:00:00"
	}else if duration == "m3" {
		date2 = now.AddDate(0,-3,0).Format(layout) + " 00:00:00"
	}else if duration == "y1" {
		date2 = now.AddDate(-1,0,0).Format(layout) + " 00:00:00"
	}else if duration == "y3" {
		date2 = now.AddDate(-3,0,0).Format(layout) + " 00:00:00"
	}else if duration == "y5" {
		date2 = now.AddDate(-5,0,0).Format(layout) + " 00:00:00"
	}else if duration == "ytd" {
		date2 = strconv.Itoa(now.Year()-1) + "-12-31 00:00:00"
	}else{
		log.Error("Missing required parameters")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameters","Missing required parameters")
	}

	var navDB []models.TrNav
	var productIDs []string
	productIDs = append(productIDs, productKeyStr)
	status, err = models.GetAllTrNavBetween(&navDB, date2, date1, productIDs)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(navDB) < 1 {
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}
	
	var navData []models.TrNavInfo
	for _, nav := range navDB {
		var data models.TrNavInfo
		date, _ := time.Parse("2006-01-02 15:04:05", nav.NavDate)
		data.NavDate = date.Format("02 Jan 2006")
		data.NavValue = nav.NavValue
		
		navData = append(navData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = navData
	
	return c.JSON(http.StatusOK, response)
}