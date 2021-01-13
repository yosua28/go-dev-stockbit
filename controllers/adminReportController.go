package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	b64 "encoding/base64"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func GetTransactionHistory(c echo.Context) error {

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

	//if user admin role 7 branch
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	if lib.Profile.RoleKey == roleKeyBranchEntry {
		log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["t.branch_key"] = strBranchKey
		} else {
			log.Error("User Branch haven't Branch")
			return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
		}
	}

	productkey := c.QueryParam("product_key")
	if productkey != "" {
		params["t.product_key"] = productkey
	}

	customerkey := c.QueryParam("customer_key")
	if customerkey != "" {
		params["c.customer_key"] = customerkey
	}

	datefrom := c.QueryParam("date_from")
	if datefrom == "" {
		log.Error("date_from parameter tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "date_from parameter tidak boleh kosong", "date_from parameter tidak boleh kosong")
	}

	dateto := c.QueryParam("date_to")
	if dateto == "" {
		log.Error("date_to parameter tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "date_to parameter tidak boleh kosong", "date_to parameter tidak boleh kosong")
	}

	layoutISO := "2006-01-02"

	from, _ := time.Parse(layoutISO, datefrom)
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)

	to, _ := time.Parse(layoutISO, dateto)
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)

	params["dateFrom"] = datefrom
	params["dateTo"] = dateto

	if from.Before(to) {
		params["dateFrom"] = datefrom
		params["dateTo"] = dateto
	}

	if from.After(to) {
		params["dateFrom"] = dateto
		params["dateTo"] = datefrom
	}

	items := []string{"product_name", "nav_date", "full_name", "cif", "sid"}

	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "t.nav_date"
			if orderBy == "product_name" {
				orderByJoin = "p.product_name"
			} else if orderBy == "nav_date" {
				orderByJoin = "t.nav_date"
			} else if orderBy == "full_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "cif" {
				orderByJoin = "c.unit_holder_idno"
			} else if orderBy == "sid" {
				orderByJoin = "c.sid_no"
			}

			params["orderBy"] = orderByJoin
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
		params["orderBy"] = "t.nav_date"
		params["orderType"] = "ASC"
	}

	paramsLike := make(map[string]string)

	customername := c.QueryParam("customer_name")
	if customername != "" {
		paramsLike["c.full_name"] = customername
	}

	var transaksi []models.TransactionCustomerHistory

	status, err = models.AdminGetTransactionCustomerHistory(&transaksi, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(transaksi) < 1 {
		log.Error("Transaksi not found")
		return lib.CustomError(http.StatusNotFound, "Transaksi not found", "Transaksi not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminCountTransactionCustomerHistory(&countData, params, paramsLike)
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
	response.Data = transaksi

	return c.JSON(http.StatusOK, response)
}

func GetListCustomerDropDown(c echo.Context) error {

	var err error
	var status int

	params := make(map[string]string)

	//if user admin role 7 branch
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	if lib.Profile.RoleKey == roleKeyBranchEntry {
		log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["d.branch_key"] = strBranchKey
		} else {
			log.Error("User Branch haven't Branch")
			return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
		}
	}

	paramsLike := make(map[string]string)

	customername := c.QueryParam("customer_name")
	if customername != "" {
		paramsLike["c.full_name"] = customername
	}

	var customer []models.CustomerDropdown

	status, err = models.GetCustomerDropdown(&customer, params, paramsLike)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(customer) < 1 {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = customer

	return c.JSON(http.StatusOK, response)
}

func GetDetailCustomerProduct(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	params := make(map[string]string)
	paramsList := make(map[string]string)
	paramsLike := make(map[string]string)

	param := c.Param("param")
	if param == "" {
		log.Error("Wrong input for parameter param")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter param", "Wrong input for parameter param")
	}

	raw, err := b64.StdEncoding.DecodeString(strings.Replace(param, "%3D", "=", 3))

	if err != nil {
		log.Error("Wrong input for parameter param")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter param", "Wrong input for parameter param")
	}

	s := strings.Split(string(raw), ",")

	var customerKey string
	var productKey string
	var dateFrom string
	var dateTo string

	for idx, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			if idx == 0 {
				customerKey = value
			}
			if idx == 1 {
				productKey = value
			}
			if idx == 2 {
				dateFrom = value
			}
			if idx == 3 {
				dateTo = value
			}
		} else {
			log.Error("Wrong input for parameter param")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter param", "Wrong input for parameter param")
			break
		}
	}

	params["a.customer_key"] = customerKey
	params["a.product_key"] = productKey

	paramsList["t.customer_key"] = customerKey
	paramsList["t.product_key"] = productKey
	paramsList["orderBy"] = "t.transaction_key"
	paramsList["orderType"] = "ASC"

	//if user admin role 7 branch
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	if lib.Profile.RoleKey == roleKeyBranchEntry {
		log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["d.branch_key"] = strBranchKey
			paramsList["t.branch_key"] = strBranchKey
		} else {
			log.Error("User Branch haven't Branch")
			return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
		}
	}

	//get header detail
	var header models.DetailHeaderTransaksiCustomer
	status, err = models.AdminGetDetailHeaderTransaksiCustomer(&header, dateFrom, dateTo, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, "Failed get data", "Failed get data")
	}

	var listTransaksi []models.TransactionConsumenProduct
	status, err = models.AdminGetTransactionConsumenProduct(&listTransaksi, paramsList, paramsLike, dateFrom, dateTo)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	var totalSubs decimal.Decimal
	var totalRedm decimal.Decimal
	var totalNetSub decimal.Decimal

	var trFormaterList []models.TransactionConsumenProduct
	for _, tr := range listTransaksi {
		var trFormater models.TransactionConsumenProduct
		trFormater.TransactionKey = tr.TransactionKey
		trFormater.TransTypeKey = tr.TransTypeKey
		trFormater.NavDate = tr.NavDate
		trFormater.TypeDescription = tr.TypeDescription
		trFormater.NavValue = tr.NavValue.Truncate(2)
		trFormater.Unit = tr.Unit.Truncate(2)
		trFormater.GrossAmount = tr.GrossAmount.Truncate(0)
		trFormater.FeeAmount = tr.FeeAmount.Truncate(0)
		trFormater.NetAmount = tr.NetAmount.Truncate(0)
		trFormaterList = append(trFormaterList, trFormater)

		if (trFormater.TransTypeKey) == 1 || (trFormater.TransTypeKey == 4) {
			totalSubs = totalSubs.Add(tr.NetAmount).Truncate(0)
		}

		if (trFormater.TransTypeKey) == 2 || (trFormater.TransTypeKey == 3) {
			totalRedm = totalRedm.Add(tr.NetAmount).Truncate(0)
		}

	}

	totalNetSub = totalSubs.Sub(totalRedm).Truncate(0)

	var result models.DataDetailTransaksiCustomerProduct
	result.DataTransaksi = header
	result.DetailTransaksi = &trFormaterList
	result.CountSubscription = totalSubs
	result.CountRedemption = totalRedm
	result.CountNetSub = totalNetSub

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = result

	return c.JSON(http.StatusOK, response)
}
