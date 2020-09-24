package controllers

import (
	"api/models"
	"api/lib"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

func GetGenLookup(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	groupKeyStr := c.QueryParam("group_key")
	groupKey, _ := strconv.ParseUint(groupKeyStr, 10, 64)
	if groupKey == 0 {
		log.Error("Group key should be number")
		return lib.CustomError(http.StatusNotFound,"Group key should be number","Group key should be number")
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
		if lkp.LkpText1 != nil{
			data.Name = *lkp.LkpText1
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