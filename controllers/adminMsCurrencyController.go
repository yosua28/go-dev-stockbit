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

func GetListMsCurrency(c echo.Context) error {
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

	items := []string{"currency_key", "code", "symbol", "flag_base"}

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

	var currency []models.MsCurrency
	status, err = models.AdminGetListMsCurrency(&currency, limit, offset, params, noLimit)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var responseData []models.MsCurrencyInfo

	for _, cur := range currency {
		var data models.MsCurrencyInfo
		data.CurrencyKey = cur.CurrencyKey
		data.Code = cur.Code
		data.Symbol = cur.Symbol
		data.Name = cur.Name
		data.FlagBase = cur.FlagBase

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func AdminGetListMsCurrency(c echo.Context) error {

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

	items := []string{"code", "symbol", "name", "flag_base"}

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
	} else {
		params["orderBy"] = "currency_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	var currency []models.ListCurrency

	status, err = models.AdminGetListCurrency(&currency, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(currency) < 1 {
		log.Error("Currency not found")
		return lib.CustomError(http.StatusNotFound, "Currency not found", "Currency not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetCurrency(&countData, params, searchLike)
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

func AdminDeleteMsCurrency(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("currency_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: currency_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key", "Missing required parameter: currency_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["currency_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsCurrency(params)
	if err != nil {
		log.Error("Error delete ms_currency")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateMsCurrency(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	code := c.FormValue("code")
	if code == "" {
		log.Error("Missing required parameter: code")
		return lib.CustomError(http.StatusBadRequest, "code can not be blank", "code can not be blank")
	} else {
		params["code"] = code
	}

	symbol := c.FormValue("symbol")
	if symbol == "" {
		log.Error("Missing required parameter: symbol")
		return lib.CustomError(http.StatusBadRequest, "symbol can not be blank", "symbol can not be blank")
	} else {
		params["symbol"] = symbol
	}

	name := c.FormValue("name")
	if name == "" {
		log.Error("Missing required parameter: name")
		return lib.CustomError(http.StatusBadRequest, "name can not be blank", "name can not be blank")
	} else {
		params["name"] = name
	}

	flagBase := c.FormValue("flag_base")
	if flagBase == "" {
		log.Error("Missing required parameter: flag_base")
		return lib.CustomError(http.StatusBadRequest, "flag_base can not be blank", "flag_base can not be blank")
	} else {
		if flagBase != "1" && flagBase != "0" {
			log.Error("Missing required parameter: flag_base")
			return lib.CustomError(http.StatusBadRequest, "flag_base must 1 / 0", "flag_base must 1 / 0")
		}
		params["flag_base"] = flagBase
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

	status, err = models.CreateMsCurrency(params)
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

func AdminUpdateMsCurrency(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

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

	code := c.FormValue("code")
	if code == "" {
		log.Error("Missing required parameter: code")
		return lib.CustomError(http.StatusBadRequest, "code can not be blank", "code can not be blank")
	} else {
		params["code"] = code
	}

	symbol := c.FormValue("symbol")
	if symbol == "" {
		log.Error("Missing required parameter: symbol")
		return lib.CustomError(http.StatusBadRequest, "symbol can not be blank", "symbol can not be blank")
	} else {
		params["symbol"] = symbol
	}

	name := c.FormValue("name")
	if name == "" {
		log.Error("Missing required parameter: name")
		return lib.CustomError(http.StatusBadRequest, "name can not be blank", "name can not be blank")
	} else {
		params["name"] = name
	}

	flagBase := c.FormValue("flag_base")
	if flagBase == "" {
		log.Error("Missing required parameter: flag_base")
		return lib.CustomError(http.StatusBadRequest, "flag_base can not be blank", "flag_base can not be blank")
	} else {
		if flagBase != "1" && flagBase != "0" {
			log.Error("Missing required parameter: flag_base")
			return lib.CustomError(http.StatusBadRequest, "flag_base must 1 / 0", "flag_base must 1 / 0")
		}
		params["flag_base"] = flagBase
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

	status, err = models.UpdateMsCurrency(params)
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

func AdminDetailMsCurrency(c echo.Context) error {
	var err error

	currencyKey := c.Param("currency_key")
	if currencyKey == "" {
		log.Error("Missing required parameter: currency_key")
		return lib.CustomError(http.StatusBadRequest, "currency_key can not be blank", "currency_key can not be blank")
	} else {
		n, err := strconv.ParseUint(currencyKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: currency_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: currency_key", "Wrong input for parameter: currency_key")
		}
	}
	var currency models.MsCurrency
	_, err = models.GetMsCurrency(&currency, currencyKey)
	if err != nil {
		log.Error("Currency not found")
		return lib.CustomError(http.StatusBadRequest, "Currency not found", "Currency not found")
	}

	responseData := make(map[string]interface{})
	responseData["currency_key"] = currency.CurrencyKey
	responseData["code"] = currency.Code
	if currency.Symbol != nil {
		responseData["symbol"] = *currency.Symbol
	} else {
		responseData["symbol"] = ""
	}
	if currency.Name != nil {
		responseData["name"] = *currency.Name
	} else {
		responseData["name"] = ""
	}
	responseData["flag_base"] = currency.FlagBase
	if currency.RecOrder != nil {
		responseData["rec_order"] = *currency.RecOrder
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
