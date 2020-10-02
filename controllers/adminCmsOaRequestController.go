package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetOaRequestList(c echo.Context) error {
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

	items := []string{"oa_request_key", "oa_request_type", "oa_entry_start", "oa_entry_end", "oa_status", "rec_order", "rec_status", "oa_risk_level"}

	params := make(map[string]string)
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

	statusStr := c.QueryParam("status")
	if statusStr != "" {
		_, err := strconv.ParseUint(statusStr, 10, 64)
		if err == nil {
			params["oa_status"] = statusStr
		} else {
			log.Error("Status should be number")
			return lib.CustomError(http.StatusBadRequest, "Status should be number", "Status should be number")
		}
	}

	var oaRequestDB []models.OaRequest
	status, err = models.GetAllOaRequest(&oaRequestDB, limit, offset, noLimit, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaRequestDB) < 1 {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var lookupIds []string
	var customerIds []string
	var userLoginIds []string
	for _, oareq := range oaRequestDB {

		if oareq.OaRequestType != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.OaRequestType, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oareq.OaRequestType, 10))
			}
		}
		if oareq.OaRiskLevel != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10))
			}
		}

		if oareq.Oastatus != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
			}
		}
		if oareq.CustomerKey != nil {
			if _, ok := lib.Find(customerIds, strconv.FormatUint(*oareq.CustomerKey, 10)); !ok {
				customerIds = append(customerIds, strconv.FormatUint(*oareq.CustomerKey, 10))
			}
		}
		if oareq.UserLoginKey != nil {
			if _, ok := lib.Find(userLoginIds, strconv.FormatUint(*oareq.UserLoginKey, 10)); !ok {
				userLoginIds = append(userLoginIds, strconv.FormatUint(*oareq.UserLoginKey, 10))
			}
		}
	}

	//mapping lookup
	var genLookup []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&genLookup, lookupIds, "lookup_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	gData := make(map[uint64]models.GenLookup)
	for _, gen := range genLookup {
		gData[gen.LookupKey] = gen
	}

	//mapping customer
	var msCustomer []models.MsCustomer
	if len(customerIds) > 0 {
		status, err = models.GetMsCustomerIn(&msCustomer, customerIds, "customer_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	cData := make(map[uint64]models.MsCustomer)
	for _, cus := range msCustomer {
		cData[cus.CustomerKey] = cus
	}

	//mapping user login
	var scUserLogin []models.ScUserLogin
	if len(userLoginIds) > 0 {
		status, err = models.GetScUserLoginIn(&scUserLogin, userLoginIds, "user_login_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	uData := make(map[uint64]models.ScUserLogin)
	for _, ul := range scUserLogin {
		uData[ul.UserLoginKey] = ul
	}

	var responseData []models.OaRequestListResponse
	for _, oareq := range oaRequestDB {
		var data models.OaRequestListResponse

		if oareq.OaRequestType != nil {
			if n, ok := gData[*oareq.OaRequestType]; ok {
				data.OaRequestType = n.LkpText1
			}
		}

		if oareq.OaRiskLevel != nil {
			if n, ok := gData[*oareq.OaRiskLevel]; ok {
				data.OaRiskLevel = n.LkpText1
			}
		}

		if oareq.Oastatus != nil {
			if n, ok := gData[*oareq.Oastatus]; ok {
				data.Oastatus = *n.LkpText1
			}
		}

		if oareq.CustomerKey != nil {
			if n, ok := cData[*oareq.CustomerKey]; ok {
				data.Customer = &n.FullName
			}
		}

		if oareq.UserLoginKey != nil {
			if n, ok := uData[*oareq.UserLoginKey]; ok {
				data.UserLoginName = &n.UloginName
				data.UserLoginFullName = &n.UloginFullName
			}
		}

		data.OaRequestKey = oareq.OaRequestKey
		data.OaEntryStart = oareq.OaEntryStart
		data.OaEntryEnd = oareq.OaEntryEnd
		data.Check1Date = oareq.Check1Date
		data.Check1Flag = oareq.Check1Flag
		data.Check1References = oareq.Check1References
		data.Check1Notes = oareq.Check1Notes
		data.Check2Date = oareq.Check2Date
		data.Check2Flag = oareq.Check2Flag
		data.Check2References = oareq.Check2References
		data.Check2Notes = oareq.Check2Notes
		data.RecOrder = oareq.RecOrder
		data.RecStatus = oareq.RecStatus
		responseData = append(responseData, data)
	}

	var countData models.OaRequestCountData
	var pagination int
	if limit > 0 {
		status, err = models.GetCountOaRequest(&countData, params)
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
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetOaRequestData(c echo.Context) error {
	var err error
	var status int
	//Get parameter limit
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var lookupIds []string

	if oareq.OaRequestType != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.OaRequestType, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oareq.OaRequestType, 10))
		}
	}
	if oareq.OaRiskLevel != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10))
		}
	}

	if oareq.Oastatus != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
		}
	}

	//mapping lookup
	var genLookup []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&genLookup, lookupIds, "lookup_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	gData := make(map[uint64]models.GenLookup)
	for _, gen := range genLookup {
		gData[gen.LookupKey] = gen
	}

	//maping response
	var responseData models.OaRequestDataResponse

	if oareq.OaRequestType != nil {
		if n, ok := gData[*oareq.OaRequestType]; ok {
			responseData.OaRequestType = n.LkpText1
		}
	}

	if oareq.OaRiskLevel != nil {
		if n, ok := gData[*oareq.OaRiskLevel]; ok {
			responseData.OaRiskLevel = n.LkpText1
		}
	}

	if oareq.Oastatus != nil {
		if n, ok := gData[*oareq.Oastatus]; ok {
			responseData.Oastatus = *n.LkpText1
		}
	}

	if oareq.CustomerKey != nil {
		var msCustomer models.MsCustomer
		var t uint64 = *oareq.CustomerKey
		str := strconv.FormatUint(t, 10)
		status, err = models.GetMsCustomer(&msCustomer, str)
		if err != nil {
			return lib.CustomError(status)
		}
		responseData.Customer = &msCustomer.FullName
	}

	if oareq.UserLoginKey != nil {
		var scUserLogin models.ScUserLogin
		var t uint64 = *oareq.UserLoginKey
		str := strconv.FormatUint(t, 10)
		status, err = models.GetScUserLoginByKey(&scUserLogin, str)
		if err != nil {
			return lib.CustomError(status)
		}
		responseData.UserLoginName = &scUserLogin.UloginName
		responseData.UserLoginFullName = &scUserLogin.UloginFullName
	}

	responseData.OaRequestKey = oareq.OaRequestKey
	responseData.OaEntryStart = oareq.OaEntryStart
	responseData.OaEntryEnd = oareq.OaEntryEnd
	responseData.Check1Date = oareq.Check1Date
	responseData.Check1Flag = oareq.Check1Flag
	responseData.Check1References = oareq.Check1References
	responseData.Check1Notes = oareq.Check1Notes
	responseData.Check2Date = oareq.Check2Date
	responseData.Check2Flag = oareq.Check2Flag
	responseData.Check2References = oareq.Check2References
	responseData.Check2Notes = oareq.Check2Notes
	responseData.RecOrder = oareq.RecOrder
	responseData.RecStatus = oareq.RecStatus

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
