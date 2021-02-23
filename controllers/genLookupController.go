package controllers

import (
	"api/lib"
	"api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetGenLookup(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	groupKeyStr := c.QueryParam("group_key")
	groupKey, _ := strconv.ParseUint(groupKeyStr, 10, 64)
	if groupKey == 0 {
		log.Error("Group key should be number")
		return lib.CustomError(http.StatusNotFound, "Group key should be number", "Group key should be number")
	}

	params["lkp_group_key"] = groupKeyStr

	var lookupDB []models.GenLookup
	status, err = models.GetAllGenLookup(&lookupDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	lookupData := make(map[uint64][]models.GenLookupInfo)

	for _, lkp := range lookupDB {
		var data models.GenLookupInfo
		if lkp.LkpName != nil {
			data.Name = *lkp.LkpName
		}
		data.Value = lkp.LookupKey
		lookupData[lkp.LkpGroupKey] = append(lookupData[lkp.LkpGroupKey], data)
	}

	var lkpGroupDB []models.GenLookupGroup
	status, err = models.GetAllGenLookupGroup(&lkpGroupDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var responseData []models.GenLookupGroupList

	for _, lkpGroup := range lkpGroupDB {
		var data models.GenLookupGroupList

		data.GroupName = lkpGroup.LkpGroupName

		if lkp, ok := lookupData[lkpGroup.LkpGroupKey]; ok {
			data.Lookup = &lkp
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

func GetMetodePerhitungan(c echo.Context) error {
	var responseData []models.MetodePerhitungan

	var data1 models.MetodePerhitungan
	data1.Key = "1"
	data1.Name = "All Units"
	var data2 models.MetodePerhitungan
	data2.Key = "2"
	data2.Name = "Unit Penyertaan"
	var data3 models.MetodePerhitungan
	data3.Key = "3"
	data3.Name = "Nominal"

	responseData = append(responseData, data1)
	responseData = append(responseData, data2)
	responseData = append(responseData, data3)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
