package controllers

import (
	"api/models"
	"api/lib"
	"api/config"
	"net/http"
	"github.com/labstack/echo"
)

func GetMsFundTypeList(c echo.Context) error {

	var err error
	var status int

	params := make(map[string]string)
	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy!=""{
		if (orderBy == "rec_order") || (orderBy == "fund_type_code") || (orderBy == "fund_type_name") {
			params["orderBy"] = orderBy
		}else{
			return lib.CustomError(http.StatusBadRequest)
		}
	}
	// Get parameter order_type
	orderType := c.QueryParam("order_type")
	if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
		params["orderType"] = orderType
	}
	var fundType []models.MsFundType
	status, err = models.GetAllMsFundType(&fundType, 0, 0, params, true)
	if err != nil {
		return lib.CustomError(status)
	}
	if len(fundType) < 1 {
		return lib.CustomError(http.StatusNotFound)
	}
	var responseData []models.MsFundTypeList
	for _, fund := range fundType {
		var data models.MsFundTypeList
		data.FundTypeKey = fund.FundTypeKey
		if fund.FundTypeCode != nil {
			data.FundTypeCode = *fund.FundTypeCode
		}
		if fund.FundTypeName != nil {
			data.FundTypeName = *fund.FundTypeName
		}
		if fund.FundTypeDesc != nil {
			data.FundTypeDesc = *fund.FundTypeDesc
		}
		if fund.RecOrder != nil {
			data.RecOrder = *fund.RecOrder
		}
		if fund.RecImage1 != nil && *fund.RecImage1 != "" {
			data.RecImage1 = config.BaseUrl + "/images/fund_type/" + *fund.RecImage1
		}else{
			data.RecImage1 = config.BaseUrl + "/images/fund_type/default.png"
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