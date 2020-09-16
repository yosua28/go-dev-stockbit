package controllers

import (
	"api/models"
	"api/lib"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetCmsPostTypeList(c echo.Context) error {

	var err error
	var status int

	params := make(map[string]string)
	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy!=""{
		if (orderBy == "rec_order") || (orderBy == "fund_type_code") || (orderBy == "fund_type_name") {
			params["orderBy"] = orderBy
		}else{
			log.Error("Wrong input for parameter: order_by")
			return lib.CustomError(http.StatusBadRequest,"Wrong input for parameter: order_by","Wrong input for parameter: order_by")
		}
	}
	// Get parameter order_type
	orderType := c.QueryParam("order_type")
	if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
		params["orderType"] = orderType
	}
	var postType []models.CmsPostType
	status, err = models.GetAllCmsPostType(&postType, 0, 0, params, true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(),"Failed get data")
	}
	if len(postType) < 1 {
		log.Error("data not found")
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}
	var postSubtype []models.CmsPostSubtype
	status, err = models.GetAllCmsPostSubtype(&postSubtype, 0, 0, params, true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(),"Failed get data")
	}

	sData := make(map[uint64][]models.CmsPostSubtypeInfo)
	for _, subtype := range postSubtype {
		var info models.CmsPostSubtypeInfo
		info.PostSubtypeKey = subtype.PostSubtypeKey
		info.PostSubtypeCode = subtype.PostSubtypeCode
		if subtype.PostSubtypeName != nil {
			info.PostSubtypeName = *subtype.PostSubtypeName
		}
		sData[subtype.PostTypeKey] = append(sData[subtype.PostTypeKey], info)
	}

	var responseData []models.CmsPostTypeList
	for _, typ := range postType {
		var data models.CmsPostTypeList
		data.PostTypeKey = typ.PostTypeKey
		data.PostTypeCode = typ.PostTypeCode
		if typ.PostTypeName != nil {
			data.PostTypeName = *typ.PostTypeName
		}
		if typ.PostTypeDesc != nil {
			data.PostTypeDesc = *typ.PostTypeDesc
		}
		if typ.PostTypeGroup != nil {
			data.PostTypeGroup = *typ.PostTypeGroup
		}
		if _, ok := sData[typ.PostTypeKey]; ok {
			data.SubType = sData[typ.PostTypeKey]
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