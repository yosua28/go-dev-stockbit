package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetListProductBankAccountAdmin(c echo.Context) error {

	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

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

	items := []string{"prod_bankacc_key", "product_code", "product_name_alt", "bank_account_name", "bank_account_purpose", "bank_fullname", "account_no", "account_holder_name"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var orderByJoin string
			orderByJoin = "pba.prod_bankacc_key"
			if orderBy == "prod_bankacc_key" {
				orderByJoin = "pba.prod_bankacc_key"
			}
			if orderBy == "product_code" {
				orderByJoin = "p.product_code"
			}
			if orderBy == "product_name_alt" {
				orderByJoin = "p.product_name_alt"
			}
			if orderBy == "bank_account_name" {
				orderByJoin = "pba.bank_account_name"
			}
			if orderBy == "bank_account_purpose" {
				orderByJoin = "bank_account_purpose.lkp_name"
			}
			if orderBy == "bank_fullname" {
				orderByJoin = "bank.bank_fullname"
			}
			if orderBy == "account_no" {
				orderByJoin = "ba.account_no"
			}
			if orderBy == "account_holder_name" {
				orderByJoin = "ba.account_holder_name"
			}

			params["orderBy"] = orderByJoin
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "pba.prod_bankacc_key"
		params["orderType"] = "ASC"
	}

	var searchData *string

	search := c.QueryParam("search")
	if search != "" {
		searchData = &search
	}

	productkey := c.QueryParam("product_key")
	if productkey != "" {
		productkeyCek, err := strconv.ParseUint(productkey, 10, 64)
		if err == nil && productkeyCek > 0 {
			params["p.product_key"] = productkey
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
		}
	}

	//mapping product bank
	var productBankAccountList []models.AdminMsProductBankAccountList
	status, err = models.AdminGetAllMsProductBankAccount(&productBankAccountList, limit, offset, params, noLimit, searchData)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminCountDataGetAllMsProductBankAccount(&countData, params, searchData)
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
	response.Data = productBankAccountList

	return c.JSON(http.StatusOK, response)
}

func GetProductBankAccountDetailAdmin(c echo.Context) error {

	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	//cek product bank account
	var productBankAccount models.MsProductBankAccount
	status, err = models.GetMsProductBankAccount(&productBankAccount, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	//cek bank account
	strBankAccountKey := strconv.FormatUint(*productBankAccount.BankAccountKey, 10)
	var bankAccount models.MsBankAccount
	status, err = models.GetBankAccount(&bankAccount, strBankAccountKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	//cek product
	strProductKey := strconv.FormatUint(*productBankAccount.ProductKey, 10)
	var product models.MsProduct
	status, err = models.GetMsProduct(&product, strProductKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.MsProductBankAccountDetailAdmin

	var lookupIds []string

	if _, ok := lib.Find(lookupIds, strconv.FormatUint(productBankAccount.BankAccountPurpose, 10)); !ok {
		lookupIds = append(lookupIds, strconv.FormatUint(productBankAccount.BankAccountPurpose, 10))
	}

	if _, ok := lib.Find(lookupIds, strconv.FormatUint(bankAccount.BankAccountType, 10)); !ok {
		lookupIds = append(lookupIds, strconv.FormatUint(bankAccount.BankAccountType, 10))
	}

	//gen lookup
	var lookupProductBank []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupProductBank, lookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupProductBank {
		gData[gen.LookupKey] = gen
	}

	responseData.ProdBankaccKey = productBankAccount.ProdBankaccKey

	//set product
	var pro models.MsProductInfo
	pro.ProductKey = product.ProductKey
	pro.ProductCode = product.ProductCode
	pro.ProductName = product.ProductName
	pro.ProductNameAlt = product.ProductNameAlt
	responseData.Product = &pro

	//set bank
	strBankKey := strconv.FormatUint(bankAccount.BankKey, 10)
	var bank models.MsBank
	status, err = models.GetMsBank(&bank, strBankKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	} else {
		log.Println("hahahah")
		var b models.MsBankList
		b.BankKey = bank.BankKey
		b.BankCode = bank.BankCode
		b.BankName = bank.BankName
		if bank.BankFullname != nil {
			b.BankFullname = *bank.BankFullname
		}
		responseData.Bank = &b
	}

	responseData.AccountNo = bankAccount.AccountNo
	responseData.AccountHolderName = bankAccount.AccountHolderName
	responseData.BranchName = bankAccount.BranchName

	//set currency
	strCurrencyKey := strconv.FormatUint(bankAccount.CurrencyKey, 10)
	var curr models.MsCurrency
	status, err = models.GetMsCurrency(&curr, strCurrencyKey)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	} else {
		log.Println("hahahah")
		var c models.MsCurrencyInfo
		c.CurrencyKey = curr.CurrencyKey
		c.Code = curr.Code
		c.Symbol = curr.Symbol
		c.Name = curr.Name
		c.FlagBase = curr.FlagBase
		responseData.Currency = &c
	}

	if n, ok := gData[bankAccount.BankAccountType]; ok {
		var trc models.LookupTrans

		trc.LookupKey = n.LookupKey
		trc.LkpGroupKey = n.LkpGroupKey
		trc.LkpCode = n.LkpCode
		trc.LkpName = n.LkpName
		responseData.BankAccountType = trc
	}

	responseData.SwiftCode = bankAccount.SwiftCode
	responseData.BankAccountName = productBankAccount.BankAccountName

	if n, ok := gData[productBankAccount.BankAccountPurpose]; ok {
		var trc models.LookupTrans

		trc.LookupKey = n.LookupKey
		trc.LkpGroupKey = n.LkpGroupKey
		trc.LkpCode = n.LkpCode
		trc.LkpName = n.LkpName
		responseData.BankAccountPurpose = trc
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
