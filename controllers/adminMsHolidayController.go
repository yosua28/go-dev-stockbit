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
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func AdminGetListMsHoliday(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
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

	items := []string{"stock_market", "holiday_date", "holiday_name"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var ord string
			if orderBy == "stock_market" {
				ord = "stock.lkp_name"
			} else if orderBy == "holiday_date" {
				ord = "h.holiday_date"
			} else if orderBy == "holiday_name" {
				ord = "h.holiday_name"
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
		params["orderBy"] = "h.holiday_date"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	var holiday []models.ListHoliday

	status, err = models.AdminGetListHoliday(&holiday, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(holiday) < 1 {
		log.Error("Country Charges not found")
		return lib.CustomError(http.StatusNotFound, "Country Charges not found", "Country Charges not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetHoliday(&countData, params, searchLike)
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
	response.Data = holiday

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMsHoliday(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("holiday_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: holiday_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: holiday_key", "Missing required parameter: holiday_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["holiday_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsHoliday(params)
	if err != nil {
		log.Error("Error delete ms_bank")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateMsHoliday(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	stockMarket := c.FormValue("stock_market")
	if stockMarket == "" {
		log.Error("Missing required parameter: stock_market")
		return lib.CustomError(http.StatusBadRequest, "stock_market can not be blank", "stock_market can not be blank")
	} else {
		n, err := strconv.ParseUint(stockMarket, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: stock_market")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: stock_market", "Wrong input for parameter: stock_market")
		}
		params["stock_market"] = stockMarket
	}

	holidayDate := c.FormValue("holiday_date")
	if holidayDate == "" {
		log.Error("Missing required parameter: holiday_date")
		return lib.CustomError(http.StatusBadRequest, "holiday_date can not be blank", "holiday_date can not be blank")
	} else {
		//validate unique holiday_date
		var countData models.CountData
		status, err = models.CountMsHolidayValidateUnique(&countData, "holiday_date", holidayDate, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: holiday_date")
			return lib.CustomError(http.StatusBadRequest, "holiday_date already used", "holiday_date already used")
		}
		params["holiday_date"] = holidayDate
	}

	shortName := c.FormValue("holiday_name")
	if shortName != "" {
		if len(shortName) > 30 {
			log.Error("holiday_name must maximal 30 character")
			return lib.CustomError(http.StatusBadRequest, "holiday_name must maximal 30 character", "holiday_name must maximal 30 character")
		}
		params["holiday_name"] = shortName
	} else {
		log.Error("Missing required parameter: holiday_name")
		return lib.CustomError(http.StatusBadRequest, "holiday_name can not be blank", "holiday_name can not be blank")
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		_, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.CreateMsHoliday(params)
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

func AdminUpdateMsHoliday(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	holidayKey := c.FormValue("holiday_key")
	if holidayKey != "" {
		n, err := strconv.ParseUint(holidayKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: holiday_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: holiday_key", "Wrong input for parameter: holiday_key")
		}
		params["holiday_key"] = holidayKey
	} else {
		log.Error("Missing required parameter: holiday_key")
		return lib.CustomError(http.StatusBadRequest, "holiday_key can not be blank", "holiday_key can not be blank")
	}

	stockMarket := c.FormValue("stock_market")
	if stockMarket == "" {
		log.Error("Missing required parameter: stock_market")
		return lib.CustomError(http.StatusBadRequest, "stock_market can not be blank", "stock_market can not be blank")
	} else {
		n, err := strconv.ParseUint(stockMarket, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: stock_market")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: stock_market", "Wrong input for parameter: stock_market")
		}
		params["stock_market"] = stockMarket
	}

	holidayDate := c.FormValue("holiday_date")
	if holidayDate == "" {
		log.Error("Missing required parameter: holiday_date")
		return lib.CustomError(http.StatusBadRequest, "holiday_date can not be blank", "holiday_date can not be blank")
	} else {
		//validate unique holiday_date
		var countData models.CountData
		status, err = models.CountMsHolidayValidateUnique(&countData, "holiday_date", holidayDate, holidayKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: holiday_date")
			return lib.CustomError(http.StatusBadRequest, "holiday_date already used", "holiday_date already used")
		}
		params["holiday_date"] = holidayDate
	}

	shortName := c.FormValue("holiday_name")
	if shortName != "" {
		if len(shortName) > 30 {
			log.Error("holiday_name must maximal 30 character")
			return lib.CustomError(http.StatusBadRequest, "holiday_name must maximal 30 character", "holiday_name must maximal 30 character")
		}
		params["holiday_name"] = shortName
	} else {
		log.Error("Missing required parameter: holiday_name")
		return lib.CustomError(http.StatusBadRequest, "holiday_name can not be blank", "holiday_name can not be blank")
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		_, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.UpdateMsHoliday(params)
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

func AdminDetailMsHoliday(c echo.Context) error {
	var err error

	holidayKey := c.Param("holiday_key")
	if holidayKey == "" {
		log.Error("Missing required parameter: holiday_key")
		return lib.CustomError(http.StatusBadRequest, "holiday_key can not be blank", "holiday_key can not be blank")
	} else {
		n, err := strconv.ParseUint(holidayKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: holiday_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: holiday_key", "Wrong input for parameter: holiday_key")
		}
	}
	var holiday models.MsHoliday
	_, err = models.GetMsHoliday(&holiday, holidayKey)
	if err != nil {
		log.Error("Holiday not found")
		return lib.CustomError(http.StatusBadRequest, "Holiday not found", "Holiday not found")
	}

	responseData := make(map[string]interface{})
	responseData["holiday_key"] = holiday.HolidayKey
	responseData["stock_market"] = holiday.StickMarketKey
	responseData["holiday_date"] = holiday.HolidayDate
	if holiday.HolidayName != nil {
		responseData["holiday_name"] = *holiday.HolidayName
	} else {
		responseData["holiday_name"] = ""
	}
	if holiday.RecOrder != nil {
		responseData["rec_order"] = *holiday.RecOrder
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
