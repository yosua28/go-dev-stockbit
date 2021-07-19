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

func AdminGetListMsBankCharges(c echo.Context) error {

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

	items := []string{"bcharges_key", "bank_network_type", "bank_name", "custodian_name", "min_nominal_trx", "value_type", "charges_value"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var ord string
			if orderBy == "bcharges_key" {
				ord = "bc.bcharges_key"
			} else if orderBy == "bank_network_type" {
				ord = "net.lkp_name"
			} else if orderBy == "bank_name" {
				ord = "b.bank_name"
			} else if orderBy == "custodian_name" {
				ord = "cb.custodian_short_name"
			} else if orderBy == "min_nominal_trx" {
				ord = "bc.min_nominal_trx"
			} else if orderBy == "value_type" {
				ord = "bc.value_type"
			} else if orderBy == "charges_value" {
				ord = "bc.charges_value"
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
		params["orderBy"] = "bc.bcharges_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	bankKey := c.QueryParam("bank_key")
	if bankKey != "" {
		params["bc.bank_key"] = bankKey
	}

	custodianKey := c.QueryParam("custodian_key")
	if custodianKey != "" {
		params["bc.custodian_key"] = custodianKey
	}

	var bankCharges []models.ListBankChargesAdmin

	status, err = models.AdminGetListBankCharges(&bankCharges, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(bankCharges) < 1 {
		log.Error("Bank Charges not found")
		return lib.CustomError(http.StatusNotFound, "Bank Charges not found", "Bank Charges not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListBankCharges(&countData, params, searchLike)
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
	response.Data = bankCharges

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMsBankCharges(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("bcharges_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: bcharges_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bcharges_key", "Missing required parameter: bcharges_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["bcharges_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsBankCharges(params)
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

func AdminCreateMsBankCharges(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	bankNetworkType := c.FormValue("bank_network_type")
	if bankNetworkType == "" {
		log.Error("Missing required parameter: bank_network_type")
		return lib.CustomError(http.StatusBadRequest, "bank_network_type can not be blank", "bank_network_type can not be blank")
	} else {
		n, err := strconv.ParseUint(bankNetworkType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bank_network_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_network_type", "Wrong input for parameter: bank_network_type")
		}
		params["bank_network_type"] = bankNetworkType
	}

	bankKey := c.FormValue("bank_key")
	if bankKey != "" {
		n, err := strconv.ParseUint(bankKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bank_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_key", "Wrong input for parameter: bank_key")
		}
		params["bank_key"] = bankKey
	}

	custodianKey := c.FormValue("custodian_key")
	if custodianKey != "" {
		n, err := strconv.ParseUint(custodianKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: custodian_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: custodian_key", "Wrong input for parameter: custodian_key")
		}
		params["custodian_key"] = custodianKey
	}

	minNominalTrx := c.FormValue("min_nominal_trx")
	if minNominalTrx != "" {
		_, err := decimal.NewFromString(minNominalTrx)
		if err != nil {
			log.Error("Wrong input for parameter: min_nominal_trx")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: min_nominal_trx", "Wrong input for parameter: min_nominal_trx")
		}
		params["min_nominal_trx"] = minNominalTrx
	} else {
		params["min_nominal_trx"] = "0"
	}

	valueType := c.FormValue("value_type")
	if valueType != "" {
		_, err := decimal.NewFromString(valueType)
		if err != nil {
			log.Error("Wrong input for parameter: value_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: value_type", "Wrong input for parameter: value_type")
		}
		params["value_type"] = valueType
	} else {
		params["value_type"] = "0"
	}

	chargesValue := c.FormValue("charges_value")
	if chargesValue != "" {
		_, err := decimal.NewFromString(chargesValue)
		if err != nil {
			log.Error("Wrong input for parameter: charges_value")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: charges_value", "Wrong input for parameter: charges_value")
		}
		params["charges_value"] = chargesValue
	} else {
		params["charges_value"] = "0"
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

	status, err = models.CreateMsBankCharges(params)
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

func AdminUpdateMsBankCharges(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	bchargesKey := c.FormValue("bcharges_key")
	if bchargesKey != "" {
		n, err := strconv.ParseUint(bchargesKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bcharges_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bcharges_key", "Wrong input for parameter: bcharges_key")
		}
		params["bcharges_key"] = bchargesKey
	} else {
		log.Error("Missing required parameter: bcharges_key")
		return lib.CustomError(http.StatusBadRequest, "bcharges_key can not be blank", "bcharges_key can not be blank")
	}

	bankNetworkType := c.FormValue("bank_network_type")
	if bankNetworkType == "" {
		log.Error("Missing required parameter: bank_network_type")
		return lib.CustomError(http.StatusBadRequest, "bank_network_type can not be blank", "bank_network_type can not be blank")
	} else {
		n, err := strconv.ParseUint(bankNetworkType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bank_network_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_network_type", "Wrong input for parameter: bank_network_type")
		}
		params["bank_network_type"] = bankNetworkType
	}

	bankKey := c.FormValue("bank_key")
	if bankKey != "" {
		n, err := strconv.ParseUint(bankKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bank_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_key", "Wrong input for parameter: bank_key")
		}
		params["bank_key"] = bankKey
	}

	chargesValue := c.FormValue("charges_value")
	if chargesValue != "" {
		_, err := decimal.NewFromString(chargesValue)
		if err != nil {
			log.Error("Wrong input for parameter: charges_value")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: charges_value", "Wrong input for parameter: charges_value")
		}
		params["charges_value"] = chargesValue
	} else {
		params["charges_value"] = "0"
	}

	minNominalTrx := c.FormValue("min_nominal_trx")
	if minNominalTrx != "" {
		_, err := decimal.NewFromString(minNominalTrx)
		if err != nil {
			log.Error("Wrong input for parameter: min_nominal_trx")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: min_nominal_trx", "Wrong input for parameter: min_nominal_trx")
		}
		params["min_nominal_trx"] = minNominalTrx
	} else {
		params["min_nominal_trx"] = "0"
	}

	valueType := c.FormValue("value_type")
	if valueType != "" {
		_, err := decimal.NewFromString(valueType)
		if err != nil {
			log.Error("Wrong input for parameter: value_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: value_type", "Wrong input for parameter: value_type")
		}
		params["value_type"] = valueType
	} else {
		params["value_type"] = "0"
	}

	custodianKey := c.FormValue("custodian_key")
	if custodianKey != "" {
		n, err := strconv.ParseUint(custodianKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: custodian_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: custodian_key", "Wrong input for parameter: custodian_key")
		}
		params["custodian_key"] = custodianKey
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

	status, err = models.UpdateMsBankCharges(params)
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

func AdminDetailBankCharges(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	bchargesKey := c.Param("bcharges_key")
	if bchargesKey == "" {
		log.Error("Missing required parameter: bcharges_key")
		return lib.CustomError(http.StatusBadRequest, "bcharges_key can not be blank", "bcharges_key can not be blank")
	} else {
		n, err := strconv.ParseUint(bchargesKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bcharges_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bcharges_key", "Wrong input for parameter: bcharges_key")
		}
	}
	var bank models.MsBankCharges
	_, err = models.GetMsBankCharges(&bank, bchargesKey)
	if err != nil {
		log.Error("Bank Charges not found")
		return lib.CustomError(http.StatusBadRequest, "Bank Charges not found", "Bank Charges not found")
	}

	responseData := make(map[string]interface{})
	responseData["bcharges_key"] = bank.BchargesKey
	responseData["bank_network_type"] = bank.BankNetworkType
	if bank.BankKey != nil {
		responseData["bank_key"] = *bank.BankKey
	} else {
		responseData["bank_key"] = ""
	}
	if bank.CustodianKey != nil {
		responseData["custodian_key"] = *bank.CustodianKey
	} else {
		responseData["custodian_key"] = ""
	}
	if bank.MinNominalTrx != nil {
		responseData["min_nominal_trx"] = *bank.MinNominalTrx
	} else {
		responseData["min_nominal_trx"] = ""
	}
	responseData["value_type"] = bank.ValueType
	if bank.ChargesValue != nil {
		responseData["charges_value"] = *bank.ChargesValue
	} else {
		responseData["charges_value"] = ""
	}
	if bank.RecOrder != nil {
		responseData["rec_order"] = *bank.RecOrder
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
