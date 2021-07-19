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

func GetListTrCurrencyRate(c echo.Context) error {
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

	params := make(map[string]string)

	items := []string{"curr_rate_key", "rate_date", "rate_type", "currency_code", "currency_name", "rate_value"}

	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var ord string
			if orderBy == "rate_date" {
				ord = "cr.rate_date"
			} else if orderBy == "curr_rate_key" {
				ord = "cr.curr_rate_key"
			} else if orderBy == "rate_type" {
				ord = "ty.lkp_name"
			} else if orderBy == "currency_code" {
				ord = "c.code"
			} else if orderBy == "currency_name" {
				ord = "c.name"
			} else if orderBy == "rate_value" {
				ord = "cr.rate_value"
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
		params["orderBy"] = "rate_date"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	rateDate := c.QueryParam("rate_date")
	if rateDate != "" {
		params["cr.rate_date"] = rateDate
	}

	rateType := c.QueryParam("rate_type")
	if rateType != "" {
		params["cr.rate_type"] = rateType
	}

	currencyKey := c.QueryParam("currency_key")
	if currencyKey != "" {
		params["cr.currency_key"] = currencyKey
	}

	var currency []models.ListCurrencyRate
	status, err = models.AdminGetListCurrencyRate(&currency, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(currency) < 1 {
		log.Error("Currency Rate not found")
		return lib.CustomError(http.StatusNotFound, "Currency Rate not found", "Currency Rate not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetCurrencyRate(&countData, params, searchLike)
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
	response.Data = currency

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteTrCurrencyRate(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("curr_rate_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: curr_rate_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: curr_rate_key", "Missing required parameter: curr_rate_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["curr_rate_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateTrCurrenctyRate(params)
	if err != nil {
		log.Error("Error delete tr_currency_rate")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateTrCurrencyRate(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	rateDate := c.FormValue("rate_date")
	if rateDate == "" {
		log.Error("Missing required parameter: rate_date")
		return lib.CustomError(http.StatusBadRequest, "rate_date can not be blank", "rate_date can not be blank")
	} else {
		params["rate_date"] = rateDate
	}

	rateType := c.FormValue("rate_type")
	if rateType != "" {
		n, err := strconv.ParseUint(rateType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rate_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rate_type", "Wrong input for parameter: rate_type")
		}
		params["rate_type"] = rateType
	} else {
		log.Error("Missing required parameter: rate_type")
		return lib.CustomError(http.StatusBadRequest, "rate_type can not be blank", "rate_type can not be blank")
	}

	rateValue := c.FormValue("rate_value")
	if rateValue != "" {
		_, err := decimal.NewFromString(rateValue)
		if err != nil {
			log.Error("Wrong input for parameter: rate_value")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rate_value", "Wrong input for parameter: rate_value")
		}
		params["rate_value"] = rateValue
	} else {
		params["rate_value"] = "0"
	}

	currencyKey := c.FormValue("currency_key")
	if currencyKey != "" {
		n, err := strconv.ParseUint(currencyKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: currency_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: currency_key", "Wrong input for parameter: currency_key")
		}
		params["currency_key"] = currencyKey
	} else {
		log.Error("Missing required parameter: currency_key")
		return lib.CustomError(http.StatusBadRequest, "currency_key can not be blank", "currency_key can not be blank")
	}

	//check duplikat date, rate_type, currency_key
	var countData models.CountData
	status, err = models.CountTrCurrencyRateValidateUniqueDateRateCurrency(&countData, rateDate, rateType, currencyKey, "")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countData.CountData) > int(0) {
		log.Error("Data with rate_date, rate_type, and currency_key already exist")
		return lib.CustomError(http.StatusBadRequest, "Data with rate_date, rate_type, and currency_key already exist", "Data with rate_date, rate_type, and currency_key already exist")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.CreateTrCurrenctyRate(params)
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

func AdminUpdateTrCurrencyRate(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	currRateKey := c.FormValue("curr_rate_key")
	if currRateKey != "" {
		n, err := strconv.ParseUint(currRateKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: curr_rate_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: curr_rate_key", "Wrong input for parameter: curr_rate_key")
		}
		params["curr_rate_key"] = currRateKey
	} else {
		log.Error("Missing required parameter: curr_rate_key")
		return lib.CustomError(http.StatusBadRequest, "curr_rate_key can not be blank", "curr_rate_key can not be blank")
	}

	rateDate := c.FormValue("rate_date")
	if rateDate == "" {
		log.Error("Missing required parameter: rate_date")
		return lib.CustomError(http.StatusBadRequest, "rate_date can not be blank", "rate_date can not be blank")
	} else {
		params["rate_date"] = rateDate
	}

	rateType := c.FormValue("rate_type")
	if rateType != "" {
		n, err := strconv.ParseUint(rateType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rate_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rate_type", "Wrong input for parameter: rate_type")
		}
		params["rate_type"] = rateType
	} else {
		log.Error("Missing required parameter: rate_type")
		return lib.CustomError(http.StatusBadRequest, "rate_type can not be blank", "rate_type can not be blank")
	}

	rateValue := c.FormValue("rate_value")
	if rateValue != "" {
		_, err := decimal.NewFromString(rateValue)
		if err != nil {
			log.Error("Wrong input for parameter: rate_value")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rate_value", "Wrong input for parameter: rate_value")
		}
		params["rate_value"] = rateValue
	} else {
		params["rate_value"] = "0"
	}

	currencyKey := c.FormValue("currency_key")
	if currencyKey != "" {
		n, err := strconv.ParseUint(currencyKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: currency_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: currency_key", "Wrong input for parameter: currency_key")
		}
		params["currency_key"] = currencyKey
	} else {
		log.Error("Missing required parameter: currency_key")
		return lib.CustomError(http.StatusBadRequest, "currency_key can not be blank", "currency_key can not be blank")
	}

	//check duplikat date, rate_type, currency_key
	var countData models.CountData
	status, err = models.CountTrCurrencyRateValidateUniqueDateRateCurrency(&countData, rateDate, rateType, currencyKey, currRateKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countData.CountData) > int(0) {
		log.Error("Data with rate_date, rate_type, and currency_key already exist")
		return lib.CustomError(http.StatusBadRequest, "Data with rate_date, rate_type, and currency_key already exist", "Data with rate_date, rate_type, and currency_key already exist")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.UpdateTrCurrenctyRate(params)
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

func AdminDetailTrCurrencyRate(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	currRateKey := c.Param("curr_rate_key")
	if currRateKey == "" {
		log.Error("Missing required parameter: curr_rate_key")
		return lib.CustomError(http.StatusBadRequest, "curr_rate_key can not be blank", "curr_rate_key can not be blank")
	} else {
		n, err := strconv.ParseUint(currRateKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: curr_rate_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: curr_rate_key", "Wrong input for parameter: curr_rate_key")
		}
	}
	var currency models.TrCurrencyRate
	_, err = models.GetTrCurrencyRate(&currency, currRateKey)
	if err != nil {
		log.Error("Currency Rate not found")
		return lib.CustomError(http.StatusBadRequest, "Currency Rate not found", "Currency Rate not found")
	}

	responseData := make(map[string]interface{})
	responseData["curr_rate_key"] = currency.CurrRateKey
	responseData["rate_date"] = currency.RateDate
	responseData["rate_type"] = currency.RateType
	responseData["rate_value"] = currency.RateValue
	responseData["currency_key"] = currency.CurrencyKey

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
