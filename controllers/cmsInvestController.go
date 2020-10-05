package controllers

import (
	"api/models"
	"api/lib"
	"api/config"
	"net/http"
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

func GetCmsInvestPurpose(c echo.Context) error {

	var err error
	var status int

	params := make(map[string]string)
	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy!=""{
		if (orderBy == "rec_order") || (orderBy == "purpose_code") || (orderBy == "purpose_name") {
			params["orderBy"] = orderBy
		}else{
			log.Error("Wrong input for parameter: order_by")
			return lib.CustomError(http.StatusBadRequest)
		}
	}
	// Get parameter order_type
	orderType := c.QueryParam("order_type")
	if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
		params["orderType"] = orderType
	}
	params["rec_status"] = "1"
	var purposeDB []models.CmsInvestPurpose
	status, err = models.GetAllCmsInvestPurpose(&purposeDB, 0, 0, params, true)
	if err != nil {
		return lib.CustomError(status)
	}
	if len(purposeDB) < 1 {
		return lib.CustomError(http.StatusNotFound)
	}
	var responseData []models.CmsInvestPurposeList
	for _, purpose := range purposeDB {
		var data models.CmsInvestPurposeList
		data.InvestPurposeKey = purpose.InvestPurposeKey
		data.PurposeCode = purpose.PurposeCode
		if purpose.PurposeName != nil {
			data.PurposeName = *purpose.PurposeName
		}
		if purpose.PurposeDesc != nil {
			data.PurposeDesc = *purpose.PurposeDesc
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

func GetCmsInvestParter(c echo.Context) error {
	
	var err error
	var status int

	params := make(map[string]string)
	// Get parameter order_by
	purposeKey := c.QueryParam("purpose_key")
	if purposeKey == "" {
		log.Error("Missing required parameter: purpose_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: purpose_key", "Missing required parameter: purpose_key")
	}
	params["invest_purpose_key"] = purposeKey
	params["rec_status"] = "1"
	var partnerDB []models.CmsInvestPartner
	status, err = models.GetAllCmsInvestPartner(&partnerDB, 0, 0, params, true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status)
	}
	if len(partnerDB) < 1 {
		return lib.CustomError(http.StatusNotFound)
	}
	var responseData []models.CmsInvestPartnerList
	for _, partner := range partnerDB {
		var data models.CmsInvestPartnerList
		data.InvestPartnerKey = partner.InvestPartnerKey
		data.PartnerCode = partner.PartnerCode
		if partner.PartnerBusinessName != nil {
			data.PartnerBusinessName = *partner.PartnerBusinessName
		}
		if partner.PartnerUrl != nil {
			data.PartnerUrl = *partner.PartnerUrl
		}
		if partner.RecImage1 != nil {
			data.RecImage1 = config.BaseUrl + "/images/partner/" + *partner.RecImage1
		}else{
			data.RecImage1 = config.BaseUrl + "/images/partner/default.png"
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