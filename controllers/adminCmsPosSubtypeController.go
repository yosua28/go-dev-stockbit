package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func AdminGetListPostSubtype(c echo.Context) error {
	var err error
	var status int

	//Get parameter limit
	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err == nil {
			if (limit == 0) || (limit > config.LimitQuery) {
				limit = config.LimitQuery
			}
		} else {
			log.Error("Limit should be number")
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
	} else {
		limit = config.LimitQuery
	}
	// Get parameter page
	pageStr := c.QueryParam("page")
	var page uint64
	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			if page == 0 {
				page = 1
			}
		} else {
			log.Error("Page should be number")
			return lib.CustomError(http.StatusBadRequest, "Page should be number", "Page should be number")
		}
	} else {
		page = 1
	}
	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}

	noLimitStr := c.QueryParam("nolimit")
	var noLimit bool
	if noLimitStr != "" {
		noLimit, err = strconv.ParseBool(noLimitStr)
		if err != nil {
			log.Error("Nolimit parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "Nolimit parameter should be true/false", "Nolimit parameter should be true/false")
		}
	} else {
		noLimit = false
	}

	params := make(map[string]string)

	items := []string{"post_subtype_key", "post_subtype_code", "post_subtype_name"}

	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			params["orderBy"] = orderBy
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	}

	params["rec_status"] = "1"

	var subtype []models.CmsPostSubtype
	status, err = models.GetAllCmsPostSubtype(&subtype, limit, offset, params, noLimit)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var responseData []models.CmsPostSubtypeDropdown

	for _, pt := range subtype {
		var data models.CmsPostSubtypeDropdown
		data.PostSubtypeKey = pt.PostSubtypeKey
		data.PostSubtypeCode = pt.PostSubtypeCode
		data.PostSubtypeName = pt.PostSubtypeName

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func AdminGetListPostSubtypeByType(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)
	posttype := c.Param("post_type")
	if posttype == "" {
		log.Error("Missing required parameter: post_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_type", "Missing required parameter: post_type")
	}
	sub, err := strconv.ParseUint(posttype, 10, 64)
	if err == nil && sub > 0 {
		params["post_type_key"] = posttype
	} else {
		log.Error("Wrong input for parameter: post_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: post_type", "Missing required parameter: post_type")
	}

	params["rec_status"] = "1"

	var subtype []models.CmsPostSubtype
	status, err = models.GetAllCmsPostSubtype(&subtype, 0, 0, params, true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var responseData []models.CmsPostSubtypeDropdown

	for _, pt := range subtype {
		var data models.CmsPostSubtypeDropdown
		data.PostSubtypeKey = pt.PostSubtypeKey
		data.PostSubtypeCode = pt.PostSubtypeCode
		data.PostSubtypeName = pt.PostSubtypeName

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
