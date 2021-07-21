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

func AdminGetListMsCountry(c echo.Context) error {

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

	items := []string{"cou_code", "cou_name", "currency_code", "currency_name", "currency_symbol"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var ord string
			if orderBy == "cou_code" {
				ord = "c.cou_code"
			} else if orderBy == "cou_name" {
				ord = "c.cou_name"
			} else if orderBy == "currency_code" {
				ord = "cur.code"
			} else if orderBy == "currency_name" {
				ord = "cur.name"
			} else if orderBy == "currency_symbol" {
				ord = "cur.symbol"
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
		params["orderBy"] = "c.country_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	var country []models.ListCountry

	status, err = models.AdminGetListCountry(&country, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(country) < 1 {
		log.Error("Country Charges not found")
		return lib.CustomError(http.StatusNotFound, "Country Charges not found", "Country Charges not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetCountry(&countData, params, searchLike)
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
	response.Data = country

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMsCountry(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("country_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: country_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: country_key", "Missing required parameter: country_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["country_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsCountry(params)
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

func AdminCreateMsCountry(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	couCode := c.FormValue("cou_code")
	if couCode == "" {
		log.Error("Missing required parameter: cou_code")
		return lib.CustomError(http.StatusBadRequest, "cou_code can not be blank", "cou_code can not be blank")
	} else {
		if len(couCode) > 5 {
			log.Error("cou_code must maximal 5 character")
			return lib.CustomError(http.StatusBadRequest, "cou_code must maximal 5 character", "cou_code must maximal 5 character")
		}
		//validate unique bank_code
		var countData models.CountData
		status, err = models.CountMsCountryValidateUnique(&countData, "cou_code", couCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: cou_code")
			return lib.CustomError(http.StatusBadRequest, "cou_code already used", "cou_code already used")
		}
		params["cou_code"] = couCode
	}

	couName := c.FormValue("cou_name")
	if couName == "" {
		log.Error("Missing required parameter: cou_name")
		return lib.CustomError(http.StatusBadRequest, "cou_name can not be blank", "cou_name can not be blank")
	} else {
		if len(couName) > 50 {
			log.Error("cou_name must maximal 50 character")
			return lib.CustomError(http.StatusBadRequest, "cou_name must maximal 50 character", "cou_name must maximal 50 character")
		}
		//validate unique bank_code
		var countData models.CountData
		status, err = models.CountMsCountryValidateUnique(&countData, "cou_name", couName, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: cou_name")
			return lib.CustomError(http.StatusBadRequest, "cou_name already used", "cou_name already used")
		}
		params["cou_name"] = couName
	}

	shortName := c.FormValue("short_name")
	if shortName != "" {
		if len(shortName) > 30 {
			log.Error("short_name must maximal 30 character")
			return lib.CustomError(http.StatusBadRequest, "short_name must maximal 30 character", "short_name must maximal 30 character")
		}
		params["short_name"] = shortName
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

	currencyKey := c.FormValue("currency_key")
	if currencyKey != "" {
		n, err := strconv.ParseUint(currencyKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: currency_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: currency_key", "Wrong input for parameter: currency_key")
		}
		params["currency_key"] = currencyKey
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

	status, err = models.CreateMsCountry(params)
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

func AdminUpdateMsCountry(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	countryKey := c.FormValue("country_key")
	if countryKey != "" {
		n, err := strconv.ParseUint(countryKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: country_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: country_key", "Wrong input for parameter: country_key")
		}
		params["country_key"] = countryKey
	} else {
		log.Error("Missing required parameter: country_key")
		return lib.CustomError(http.StatusBadRequest, "country_key can not be blank", "country_key can not be blank")
	}

	couCode := c.FormValue("cou_code")
	if couCode == "" {
		log.Error("Missing required parameter: cou_code")
		return lib.CustomError(http.StatusBadRequest, "cou_code can not be blank", "cou_code can not be blank")
	} else {
		if len(couCode) > 5 {
			log.Error("cou_code must maximal 5 character")
			return lib.CustomError(http.StatusBadRequest, "cou_code must maximal 5 character", "cou_code must maximal 5 character")
		}
		//validate unique bank_code
		var countData models.CountData
		status, err = models.CountMsCountryValidateUnique(&countData, "cou_code", couCode, countryKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: cou_code")
			return lib.CustomError(http.StatusBadRequest, "cou_code already used", "cou_code already used")
		}
		params["cou_code"] = couCode
	}

	couName := c.FormValue("cou_name")
	if couName == "" {
		log.Error("Missing required parameter: cou_name")
		return lib.CustomError(http.StatusBadRequest, "cou_name can not be blank", "cou_name can not be blank")
	} else {
		if len(couName) > 50 {
			log.Error("cou_name must maximal 50 character")
			return lib.CustomError(http.StatusBadRequest, "cou_name must maximal 50 character", "cou_name must maximal 50 character")
		}
		//validate unique bank_code
		var countData models.CountData
		status, err = models.CountMsCountryValidateUnique(&countData, "cou_name", couName, countryKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: cou_name")
			return lib.CustomError(http.StatusBadRequest, "cou_name already used", "cou_name already used")
		}
		params["cou_name"] = couName
	}

	shortName := c.FormValue("short_name")
	if shortName != "" {
		if len(shortName) > 30 {
			log.Error("short_name must maximal 30 character")
			return lib.CustomError(http.StatusBadRequest, "short_name must maximal 30 character", "short_name must maximal 30 character")
		}
		params["short_name"] = shortName
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

	currencyKey := c.FormValue("currency_key")
	if currencyKey != "" {
		n, err := strconv.ParseUint(currencyKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: currency_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: currency_key", "Wrong input for parameter: currency_key")
		}
		params["currency_key"] = currencyKey
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

	status, err = models.UpdateMsCountry(params)
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

func AdminDetailMsCountry(c echo.Context) error {
	var err error

	countryKey := c.Param("country_key")
	if countryKey == "" {
		log.Error("Missing required parameter: country_key")
		return lib.CustomError(http.StatusBadRequest, "country_key can not be blank", "country_key can not be blank")
	} else {
		n, err := strconv.ParseUint(countryKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: country_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: country_key", "Wrong input for parameter: country_key")
		}
	}
	var country models.MsCountry
	_, err = models.GetMsCountry(&country, countryKey)
	if err != nil {
		log.Error("Country not found")
		return lib.CustomError(http.StatusBadRequest, "Country not found", "Country not found")
	}

	responseData := make(map[string]interface{})
	responseData["country_key"] = country.CountryKey
	responseData["cou_code"] = country.CouCode
	responseData["cou_name"] = country.CouName
	if country.ShortName != nil {
		responseData["short_name"] = *country.ShortName
	} else {
		responseData["short_name"] = ""
	}
	responseData["flag_base"] = country.FlagBase
	if country.CurrencyKey != nil {
		responseData["currency_key"] = *country.CurrencyKey
	} else {
		responseData["currency_key"] = ""
	}
	if country.RecOrder != nil {
		responseData["rec_order"] = *country.RecOrder
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
