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
	// var roleKeyBranchEntry uint64
	// roleKeyBranchEntry = 7
	// if lib.Profile.RoleKey == roleKeyBranchEntry {
	// 	log.Println(lib.Profile)
	// 	if lib.Profile.BranchKey != nil {
	// 		strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
	// 		params["d.branch_key"] = strBranchKey
	// 	} else {
	// 		log.Error("User Branch haven't Branch")
	// 		return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
	// 	}
	// }

	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		params["c.openacc_branch_key"] = branchKey
	} else {
		//if user category  = 3 -> user branch, 2 = user HO
		var userCategory uint64
		userCategory = 3
		if lib.Profile.UserCategoryKey == userCategory {
			log.Println(lib.Profile)
			if lib.Profile.BranchKey != nil {
				strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
				params["c.openacc_branch_key"] = strBranchKey
			} else {
				log.Error("User Branch haven't Branch")
				return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
			}
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

func GetBankProductTransaction(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	//product_key
	productKey := c.QueryParam("product_key")
	if productKey == "" {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}
	productKeyCek, err := strconv.ParseUint(productKey, 10, 64)
	if err == nil && productKeyCek > 0 {
		params["t.product_key"] = productKey
	} else {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	//nav_date
	navdate := c.QueryParam("nav_date")
	if navdate == "" {
		log.Error("Wrong input for parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}
	params["t.nav_date"] = navdate

	//trans_type
	transtype := c.QueryParam("trans_type")
	if transtype == "" {
		log.Error("Wrong input for parameter: trans_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_type", "Missing required parameter: trans_type")
	}
	if (transtype == "1") || (transtype == "2") {
		params["t.trans_type_key"] = transtype
	} else {
		log.Error("Wrong input for parameter: trans_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_type", "Missing required parameter: trans_type")
	}

	var bankTransaction []models.BankProductTransactionReport
	status, err = models.AdminGetBankProductTransactionReport(&bankTransaction, params)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	if len(bankTransaction) < 1 {
		log.Error("Bank not found")
		return lib.CustomError(http.StatusNotFound, "Bank not found", "Bank not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = bankTransaction

	return c.JSON(http.StatusOK, response)
}

func GetTransactionReportSubscribeDaily(c echo.Context) error {
	return getResultReportDaily("1", c)
}

func GetTransactionReportRedemptionDaily(c echo.Context) error {
	return getResultReportDaily("2", c)
}

func getResultReportDaily(trans_type string, c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	params := make(map[string]string)

	//product_key
	productKey := c.QueryParam("product_key")
	if productKey == "" {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}
	productKeyCek, err := strconv.ParseUint(productKey, 10, 64)
	if err == nil && productKeyCek > 0 {
		params["t.product_key"] = productKey
	} else {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	//nav_date
	navdate := c.QueryParam("nav_date")
	if navdate == "" {
		log.Error("Wrong input for parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}
	params["t.nav_date"] = navdate

	//prod_bankacc_key
	prodbankacckey := c.QueryParam("prod_bankacc_key")
	if prodbankacckey == "" {
		log.Error("Wrong input for parameter: prod_bankacc_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_bankacc_key", "Missing required parameter: prod_bankacc_key")
	}
	prodbankacckeyCek, err := strconv.ParseUint(prodbankacckey, 10, 64)
	if err == nil && prodbankacckeyCek > 0 {
		params["ba.prod_bankacc_key"] = prodbankacckey
	} else {
		log.Error("Wrong input for parameter: prod_bankacc_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: prod_bankacc_key", "Missing required parameter: prod_bankacc_key")
	}

	params["t.trans_type_key"] = trans_type

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

	items := []string{"transaction_key", "nav_date", "full_name"}
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "t.transaction_key"
			if orderBy == "transaction_key" {
				orderByJoin = "t.transaction_key"
			} else if orderBy == "nav_date" {
				orderByJoin = "t.nav_date"
			} else if orderBy == "full_name" {
				orderByJoin = "c.full_name"
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
		params["orderBy"] = "t.transaction_key"
		params["orderType"] = "ASC"
	}

	//get data header
	var header models.HeaderDailySubsRedmBatchForm
	_, err = models.AdminGetHeaderDailySubsRedmBatchForm(&header, params)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, "Failed get data", "Failed get data")
		}
	}

	//get data list
	var datas []models.DailySubsRedmBatchForm
	_, err = models.AdminGetDailySubsRedmBatchForm(&datas, limit, offset, params, noLimit)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, "Failed get data", "Failed get data")
		}
	}

	var totalUnits decimal.Decimal
	var totalAmount decimal.Decimal
	var totalFeeAmount decimal.Decimal
	var totalNetSub decimal.Decimal

	var responseData []models.ResponseDailySubsRedmBatchForm
	for _, tr := range datas {

		var trFormater models.ResponseDailySubsRedmBatchForm
		trFormater.Sid = tr.Sid
		trFormater.IfuaNo = tr.IfuaNo
		trFormater.AccountNo = tr.AccountNo
		trFormater.UnitHolderIDNo = tr.UnitHolderIDNo
		trFormater.FullName = tr.FullName
		trFormater.Amount = tr.Amount.Truncate(0)
		trFormater.FeeAmount = tr.FeeAmount.Truncate(0)
		trFormater.NettAmount = tr.NettAmount.Truncate(0)
		trFormater.BankFullName = tr.BankFullName
		trFormater.NoRekening = tr.NoRekening
		trFormater.TypeDescription = tr.TypeDescription

		if trans_type == "2" {
			unit := tr.Unit.Truncate(2)
			trFormater.Unit = &unit

			var noteRedm models.NotesRedemption
			_, err = models.AdminGetNotesRedemption(&noteRedm, tr.CustomerKey, tr.ProductKey, navdate)
			if err == nil {

				trFormater.Notes1 = noteRedm.Note1

				note2 := "Unit : " + noteRedm.Unit.Truncate(2).String() + "  -  Amount : " + noteRedm.Amount.Truncate(2).String()
				trFormater.Notes2 = &note2

				trFormater.Notes3 = noteRedm.Note3

			}

			if tr.PaymentDate != nil {
				trFormater.PaymentDate = tr.PaymentDate
			} else {
				paymentdate := ""
				trFormater.PaymentDate = &paymentdate
			}

		}

		responseData = append(responseData, trFormater)

		//count
		totalAmount = totalAmount.Add(tr.Amount).Truncate(0)
		totalFeeAmount = totalFeeAmount.Add(tr.FeeAmount).Truncate(0)
		totalNetSub = totalNetSub.Add(tr.NettAmount).Truncate(0)

		if trans_type == "2" { //REDM
			totalUnits = totalUnits.Add(tr.Amount).Truncate(0)
		}

	}

	//get count
	var responseCount models.CountNominal
	if trans_type == "2" {
		responseCount.CountUnit = &totalUnits
	}
	responseCount.CountAmount = totalAmount
	responseCount.CountFeeAmount = totalFeeAmount
	responseCount.CountNettAmount = totalNetSub

	var result models.ResponseDailySubscriptionBatchForm
	result.Header = header
	result.Data = &responseData
	result.Count = responseCount

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err := models.AdminCountDailySubsRedmBatchForm(&countData, params)
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
	response.Data = result

	return c.JSON(http.StatusOK, response)
}

func GetDailyTransactionReport(c echo.Context) error {

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

	productkey := c.QueryParam("product_key")
	if productkey != "" {
		params["t.product_key"] = productkey
	}

	customerkey := c.QueryParam("customer_type")
	if customerkey != "" {
		params["ct.lookup_key"] = customerkey
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

	division := c.QueryParam("division")
	if division != "" {
		params["division.lookup_key"] = division
	}

	branch := c.QueryParam("branch")
	if branch != "" {
		params["t.branch_key"] = branch
	}

	sales := c.QueryParam("sales")
	if sales != "" {
		params["a.agent_key"] = sales
	}

	items := []string{"client_name", "product", "category", "division", "branch", "sales"}

	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "t.nav_date"
			if orderBy == "client_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "product" {
				orderByJoin = "p.product_code"
			} else if orderBy == "category" {
				orderByJoin = "ct.lkp_name"
			} else if orderBy == "division" {
				orderByJoin = "division.lkp_name"
			} else if orderBy == "branch" {
				orderByJoin = "b.branch_name"
			} else if orderBy == "sales" {
				orderByJoin = "a.agent_name"
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
		params["orderType"] = "DESC"
	}

	var transaksi []models.DailyTransactionReportField

	status, err = models.DailyTransactionReport(&transaksi, limit, offset, params, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(transaksi) < 1 {
		log.Error("Transaksi not found")
		return lib.CustomError(http.StatusNotFound, "Transaksi not found", "Transaksi not found")
	}

	var transaksiTotal models.DailyTransactionReportTotalField

	status, err = models.DailyTransactionReportTotal(&transaksiTotal, params)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var responseData models.DailyTransactionReportResponse
	responseData.Total = &transaksiTotal
	responseData.DailyTransactionReport = &transaksi

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.DailyTransactionReportCountRow(&countData, params)
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

func GetSubscriptionBatchConfirmation(c echo.Context) error {

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

	productkey := c.QueryParam("product_key")
	if productkey != "" {
		params["t.product_key"] = productkey
	} else {
		log.Error("product_key parameter tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "product_key parameter tidak boleh kosong", "product_key parameter tidak boleh kosong")
	}

	date := c.QueryParam("date")
	if date == "" {
		log.Error("date parameter tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "date parameter tidak boleh kosong", "date parameter tidak boleh kosong")
	} else {
		params["t.nav_date"] = date
	}

	paymentMethod := c.QueryParam("payment_method")
	if paymentMethod != "" {
		params["t.payment_method"] = paymentMethod
	}

	items := []string{"no_sid", "account_no", "bk_unit_holder", "investor_name", "total_amount", "fee_percent",
		"nett_amount", "unit", "bank", "payment_method"}

	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "t.nav_date"
			if orderBy == "no_sid" {
				orderByJoin = "c.sid_no"
			} else if orderBy == "account_no" {
				orderByJoin = "c.unit_holder_idno"
			} else if orderBy == "bk_unit_holder" {
				orderByJoin = "ba.account_holder_name"
			} else if orderBy == "investor_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "total_amount" {
				orderByJoin = "t.total_amount"
			} else if orderBy == "fee_percent" {
				orderByJoin = "t.trans_fee_percent"
			} else if orderBy == "nett_amount" {
				orderByJoin = "t.trans_amount"
			} else if orderBy == "unit" {
				orderByJoin = "tc.confirmed_unit"
			} else if orderBy == "bank" {
				orderByJoin = "bank.bank_name"
			} else if orderBy == "payment_method" {
				orderByJoin = "p_method.lkp_name"
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
		params["orderType"] = "DESC"
	}

	var listTransaksi []models.SubscriptionBatchConfirmationField
	status, err = models.SubscriptionBatchConfirmation(&listTransaksi, limit, offset, params, noLimit)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	if len(listTransaksi) < 1 {
		log.Error("Transaksi not found")
		return lib.CustomError(http.StatusNotFound, "Transaksi not found", "Transaksi not found")
	}

	var totalAmount decimal.Decimal
	var totalFeeAmount decimal.Decimal
	var totalNettAmount decimal.Decimal

	tradeDate := listTransaksi[0].NavDate
	productName := listTransaksi[0].ProductName
	var nab decimal.Decimal
	nab = listTransaksi[0].NavValue.Truncate(2)
	referenceNo := listTransaksi[0].BatchDisplayNo

	var trFormaterList []models.SubscriptionBatchConfirmationField
	for _, tr := range listTransaksi {
		var trFormater models.SubscriptionBatchConfirmationField
		trFormater.NoSid = tr.NoSid
		trFormater.AccountNo = tr.AccountNo
		trFormater.BkUnitHolder = tr.BkUnitHolder
		trFormater.InvestorName = tr.InvestorName
		trFormater.Amount = tr.Amount.Truncate(0)
		trFormater.FeePercent = tr.FeePercent
		trFormater.FeeAmount = tr.FeeAmount.Truncate(0)
		trFormater.NettAmount = tr.NettAmount.Truncate(0)
		trFormater.Unit = tr.Unit.Truncate(2)
		trFormater.Bank = tr.Bank
		trFormater.TransType = tr.TransType
		trFormater.PaymentMethod = tr.PaymentMethod
		trFormaterList = append(trFormaterList, trFormater)

		totalAmount = totalAmount.Add(tr.Amount).Truncate(0)
		totalFeeAmount = totalFeeAmount.Add(tr.FeeAmount).Truncate(0)
		totalNettAmount = totalNettAmount.Add(tr.NettAmount).Truncate(0)
	}

	var result models.SubscriptionBatchConfirmationResponse
	result.ProductName = productName
	result.TradeDate = &tradeDate
	result.Nab = &nab
	result.ReferenceNo = &referenceNo
	result.TotalAmount = totalAmount
	result.TotalFeeAmount = totalFeeAmount
	result.TotalUnit = totalNettAmount
	result.Data = &trFormaterList

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.SubscriptionBatchConfirmationCount(&countData, params)
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
	response.Pagination = pagination
	response.Data = result

	return c.JSON(http.StatusOK, response)
}

func GetRedemptionBatchConfirmation(c echo.Context) error {

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

	productkey := c.QueryParam("product_key")
	if productkey != "" {
		params["t.product_key"] = productkey
	} else {
		log.Error("product_key parameter tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "product_key parameter tidak boleh kosong", "product_key parameter tidak boleh kosong")
	}

	date := c.QueryParam("date")
	if date == "" {
		log.Error("date parameter tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "date parameter tidak boleh kosong", "date parameter tidak boleh kosong")
	} else {
		params["t.nav_date"] = date
	}

	paymentMethod := c.QueryParam("payment_method")
	if paymentMethod != "" {
		params["t.payment_method"] = paymentMethod
	}

	items := []string{"no_sid", "account_no", "bk_unit_holder", "investor_name", "total_amount", "fee_percent",
		"nett_amount", "unit", "bank_account", "payment_method", "bank_account_name", "bank_name", "bank_branch"}

	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "t.nav_date"
			if orderBy == "no_sid" {
				orderByJoin = "c.sid_no"
			} else if orderBy == "account_no" {
				orderByJoin = "c.unit_holder_idno"
			} else if orderBy == "bk_unit_holder" {
				orderByJoin = "ba.account_holder_name"
			} else if orderBy == "investor_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "total_amount" {
				orderByJoin = "t.total_amount"
			} else if orderBy == "fee_percent" {
				orderByJoin = "t.trans_fee_percent"
			} else if orderBy == "nett_amount" {
				orderByJoin = "t.trans_amount"
			} else if orderBy == "unit" {
				orderByJoin = "tc.confirmed_unit"
			} else if orderBy == "bank_account" {
				orderByJoin = "ba.account_no"
			} else if orderBy == "payment_method" {
				orderByJoin = "p_method.lkp_name"
			} else if orderBy == "bank_account_name" {
				orderByJoin = "ba.account_holder_name"
			} else if orderBy == "bank_name" {
				orderByJoin = "bank.bank_name"
			} else if orderBy == "bank_branch" {
				orderByJoin = "ba.branch_name"
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
		params["orderType"] = "DESC"
	}

	var listTransaksi []models.RedemptionBatchConfirmationField
	status, err = models.RedemptionBatchConfirmation(&listTransaksi, limit, offset, params, noLimit)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	if len(listTransaksi) < 1 {
		log.Error("Transaksi not found")
		return lib.CustomError(http.StatusNotFound, "Transaksi not found", "Transaksi not found")
	}

	var totalAmount decimal.Decimal
	var totalFeeAmount decimal.Decimal
	var totalNettAmount decimal.Decimal

	tradeDate := listTransaksi[0].NavDate
	productName := listTransaksi[0].ProductName
	var nab decimal.Decimal
	nab = listTransaksi[0].NavValue.Truncate(2)
	referenceNo := listTransaksi[0].BatchDisplayNo

	var trFormaterList []models.RedemptionBatchConfirmationField
	for _, tr := range listTransaksi {
		var trFormater models.RedemptionBatchConfirmationField
		trFormater.NoSid = tr.NoSid
		trFormater.AccountNo = tr.AccountNo
		trFormater.BkUnitHolder = tr.BkUnitHolder
		trFormater.InvestorName = tr.InvestorName
		trFormater.Amount = tr.Amount.Truncate(0)
		trFormater.FeePercent = tr.FeePercent
		trFormater.FeeAmount = tr.FeeAmount.Truncate(0)
		trFormater.NettAmount = tr.NettAmount.Truncate(0)
		trFormater.Unit = tr.Unit.Truncate(2)
		trFormater.PaymentDate = tr.PaymentDate
		trFormater.BankAccount = tr.BankAccount
		trFormater.BankAccountName = tr.BankAccountName
		trFormater.BankName = tr.BankName
		trFormater.BankBranch = tr.BankBranch
		trFormater.PaymentMethod = tr.PaymentMethod
		trFormaterList = append(trFormaterList, trFormater)

		totalAmount = totalAmount.Add(tr.Amount).Truncate(0)
		totalFeeAmount = totalFeeAmount.Add(tr.FeeAmount).Truncate(0)
		totalNettAmount = totalNettAmount.Add(tr.NettAmount).Truncate(0)
	}

	var result models.RedemptionBatchConfirmationResponse
	result.ProductName = productName
	result.TradeDate = &tradeDate
	result.Nab = &nab
	result.ReferenceNo = &referenceNo
	result.TotalAmount = totalAmount
	result.TotalFeeAmount = totalFeeAmount
	result.TotalUnit = totalNettAmount
	result.Data = &trFormaterList

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.RedemptionBatchConfirmationCount(&countData, params)
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
	response.Pagination = pagination
	response.Data = result

	return c.JSON(http.StatusOK, response)
}
