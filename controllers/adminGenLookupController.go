package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func AdminGetListLookup(c echo.Context) error {

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

	items := []string{"lkp_group_code", "lkp_code", "lkp_name", "lkp_desc", "lkp_group_key"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var ord string
			if orderBy == "lkp_group_code" {
				ord = "gg.lkp_group_code"
			} else {
				ord = "g." + orderBy
			}
			params["orderBy"] = ord
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "g.lkp_group_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	lkpGroupKey := c.QueryParam("lkp_group_key")
	if lkpGroupKey != "" {
		params["g.lkp_group_key"] = lkpGroupKey
	}

	var lookup []models.ListLookup

	status, err = models.AdminGetLookup(&lookup, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(lookup) < 1 {
		log.Error("Lookup not found")
		return lib.CustomError(http.StatusNotFound, "Lookup not found", "Lookup not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetLookup(&countData, params, searchLike)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) < int(limit) {
			pagination = 1
		} else {
			calc := math.Ceil(float64(countData.CountData) / float64(limit))
			pagination = int(calc)
		}
	} else {
		pagination = 1
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = lookup

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteLookup(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("lookup_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: lookup_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: lookup_key", "Missing required parameter: lookup_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["lookup_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateLookup(params)
	if err != nil {
		log.Error("Error delete gen_lookup")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateLookup(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	lkpGroupKey := c.FormValue("lkp_group_key")
	if lkpGroupKey == "" {
		log.Error("Missing required parameter: lkp_group_key")
		return lib.CustomError(http.StatusBadRequest, "lkp_group_key can not be blank", "lkp_group_key can not be blank")
	} else {
		n, err := strconv.ParseUint(lkpGroupKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: lkp_group_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: lkp_group_key", "Wrong input for parameter: lkp_group_key")
		}
		params["lkp_group_key"] = lkpGroupKey
	}

	lkpCode := c.FormValue("lkp_code")
	if lkpCode == "" {
		log.Error("Missing required parameter: lkp_code")
		return lib.CustomError(http.StatusBadRequest, "lkp_code can not be blank", "lkp_code can not be blank")
	} else {
		params["lkp_code"] = lkpCode
	}

	lkpName := c.FormValue("lkp_name")
	if lkpName == "" {
		log.Error("Missing required parameter: lkp_name")
		return lib.CustomError(http.StatusBadRequest, "lkp_name can not be blank", "lkp_name can not be blank")
	} else {
		params["lkp_name"] = lkpName
	}

	lkpDesc := c.FormValue("lkp_desc")
	if lkpDesc != "" {
		params["lkp_desc"] = lkpDesc
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.CreateLookup(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminUpdateLookup(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	lookupKey := c.FormValue("lookup_key")
	if lookupKey == "" {
		log.Error("Missing required parameter: lookup_key")
		return lib.CustomError(http.StatusBadRequest, "lookup_key can not be blank", "lookup_key can not be blank")
	} else {
		n, err := strconv.ParseUint(lookupKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: lookup_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: lookup_key", "Wrong input for parameter: lookup_key")
		}
		var lookup models.GenLookup
		status, err = models.GetLookup(&lookup, lookupKey)
		if err != nil {
			log.Error("Lookup not found")
			return lib.CustomError(http.StatusBadRequest, "Lookup not found", "Lookup not found")
		}
		params["lookup_key"] = lookupKey
	}

	lkpGroupKey := c.FormValue("lkp_group_key")
	if lkpGroupKey == "" {
		log.Error("Missing required parameter: lkp_group_key")
		return lib.CustomError(http.StatusBadRequest, "lkp_group_key can not be blank", "lkp_group_key can not be blank")
	} else {
		n, err := strconv.ParseUint(lkpGroupKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: lkp_group_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: lkp_group_key", "Wrong input for parameter: lkp_group_key")
		}
		params["lkp_group_key"] = lkpGroupKey
	}

	lkpCode := c.FormValue("lkp_code")
	if lkpCode == "" {
		log.Error("Missing required parameter: lkp_code")
		return lib.CustomError(http.StatusBadRequest, "lkp_code can not be blank", "lkp_code can not be blank")
	} else {
		params["lkp_code"] = lkpCode
	}

	lkpName := c.FormValue("lkp_name")
	if lkpName == "" {
		log.Error("Missing required parameter: lkp_name")
		return lib.CustomError(http.StatusBadRequest, "lkp_name can not be blank", "lkp_name can not be blank")
	} else {
		params["lkp_name"] = lkpName
	}

	lkpDesc := c.FormValue("lkp_desc")
	if lkpDesc != "" {
		params["lkp_desc"] = lkpDesc
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.UpdateLookup(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminDetailLookup(c echo.Context) error {
	var err error

	lookupKey := c.Param("lookup_key")
	if lookupKey == "" {
		log.Error("Missing required parameter: lookup_key")
		return lib.CustomError(http.StatusBadRequest, "lookup_key can not be blank", "lookup_key can not be blank")
	} else {
		n, err := strconv.ParseUint(lookupKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: lookup_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: lookup_key", "Wrong input for parameter: lookup_key")
		}
	}

	var lookup models.GenLookup
	_, err = models.GetLookup(&lookup, lookupKey)
	if err != nil {
		log.Error("Lookup not found")
		return lib.CustomError(http.StatusBadRequest, "Lookup not found", "Lookup not found")
	}

	responseData := make(map[string]interface{})
	responseData["lookup_key"] = lookup.LookupKey
	responseData["lkp_group_key"] = lookup.LkpGroupKey
	responseData["lkp_code"] = lookup.LkpCode
	responseData["lkp_name"] = lookup.LkpName
	if lookup.LkpDesc != nil {
		responseData["lkp_desc"] = *lookup.LkpDesc
	} else {
		responseData["lkp_desc"] = ""
	}
	if lookup.RecOrder != nil {
		responseData["rec_order"] = *lookup.RecOrder
	} else {
		responseData["rec_order"] = ""
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
