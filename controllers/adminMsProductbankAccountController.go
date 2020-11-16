package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"net/http"
	"strconv"
	"time"

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

func DeleteProductBankAccountAdmin(c echo.Context) error {
	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	key := c.FormValue("key")
	if key == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	keyCek, err := strconv.ParseUint(key, 10, 64)
	if err == nil && keyCek > 0 {
		params["prod_bankacc_key"] = key
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	//cek product bank account
	var productBankAccount models.MsProductBankAccount
	status, err := models.GetMsProductBankAccount(&productBankAccount, key)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateMsProductBankAccount(params)
	if err != nil {
		log.Error("Failed update request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed update data")
	}

	if productBankAccount.BankAccountKey != nil {
		strBankAccount := strconv.FormatUint(*productBankAccount.BankAccountKey, 10)
		paramsBA := make(map[string]string)
		paramsBA["bank_account_key"] = strBankAccount
		paramsBA["rec_status"] = "0"
		paramsBA["rec_deleted_date"] = time.Now().Format(dateLayout)
		paramsBA["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		status, err = models.UpdateMsBankAccount(paramsBA)
		if err != nil {
			log.Error("Failed update request data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed update data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func CreateAdminMsProductBankAccount(c echo.Context) error {
	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	paramsProBankAcc := make(map[string]string)
	paramsBankAcc := make(map[string]string)

	//product_key
	productkey := c.FormValue("product_key")
	if productkey == "" {
		log.Error("Missing required parameter: product_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key cann't be blank", "Missing required parameter: product_key cann't be blank")
	}
	strproductkey, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && strproductkey > 0 {
		paramsProBankAcc["product_key"] = productkey
	} else {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	//bank_key
	bankkey := c.FormValue("bank_key")
	if bankkey == "" {
		log.Error("Missing required parameter: bank_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key cann't be blank", "Missing required parameter: bank_key cann't be blank")
	}
	strbankkey, err := strconv.ParseUint(bankkey, 10, 64)
	if err == nil && strbankkey > 0 {
		paramsBankAcc["bank_key"] = bankkey
	} else {
		log.Error("Wrong input for parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key", "Missing required parameter: bank_key")
	}

	//account_no
	accountno := c.FormValue("account_no")
	if accountno == "" {
		log.Error("Missing required parameter: account_no cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_no cann't be blank", "Missing required parameter: account_no cann't be blank")
	}
	paramsBankAcc["account_no"] = accountno

	//account_holder_name
	accountholdername := c.FormValue("account_holder_name")
	if accountholdername == "" {
		log.Error("Missing required parameter: account_holder_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_holder_name cann't be blank", "Missing required parameter: account_holder_name cann't be blank")
	}
	paramsBankAcc["account_holder_name"] = accountholdername

	//branch_name
	branchname := c.FormValue("branch_name")
	if branchname != "" {
		paramsBankAcc["branch_name"] = branchname
	}

	//currency_key
	currencykey := c.FormValue("currency_key")
	if currencykey == "" {
		log.Error("Missing required parameter: currency_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key cann't be blank", "Missing required parameter: currency_key cann't be blank")
	}
	strcurrencykey, err := strconv.ParseUint(currencykey, 10, 64)
	if err == nil && strcurrencykey > 0 {
		paramsBankAcc["currency_key"] = currencykey
	} else {
		log.Error("Wrong input for parameter: currency_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key", "Missing required parameter: currency_key")
	}

	//bank_account_type
	bankaccounttype := c.FormValue("bank_account_type")
	if bankaccounttype == "" {
		log.Error("Missing required parameter: bank_account_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_type cann't be blank", "Missing required parameter: bank_account_type cann't be blank")
	}
	strbankaccounttype, err := strconv.ParseUint(bankaccounttype, 10, 64)
	if err == nil && strbankaccounttype > 0 {
		paramsBankAcc["bank_account_type"] = bankaccounttype
	} else {
		log.Error("Wrong input for parameter: bank_account_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_type", "Missing required parameter: bank_account_type")
	}

	paramsBankAcc["rec_domain"] = "132"

	//swift_code
	swiftcode := c.FormValue("swift_code")
	if swiftcode != "" {
		paramsBankAcc["swift_code"] = swiftcode
	}

	//bank_account_name
	bankaccountname := c.FormValue("bank_account_name")
	if bankaccountname == "" {
		log.Error("Missing required parameter: bank_account_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_name cann't be blank", "Missing required parameter: bank_account_name cann't be blank")
	}
	paramsProBankAcc["bank_account_name"] = bankaccountname

	//bank_account_type
	bankaccountpurpose := c.FormValue("bank_account_purpose")
	if bankaccountpurpose == "" {
		log.Error("Missing required parameter: bank_account_purpose cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_purpose cann't be blank", "Missing required parameter: bank_account_purpose cann't be blank")
	}
	strbankaccountpurpose, err := strconv.ParseUint(bankaccountpurpose, 10, 64)
	if err == nil && strbankaccountpurpose > 0 {
		paramsProBankAcc["bank_account_purpose"] = bankaccountpurpose
	} else {
		log.Error("Wrong input for parameter: bank_account_purpose")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_purpose", "Missing required parameter: bank_account_purpose")
	}

	dateLayout := "2006-01-02 15:04:05"
	paramsBankAcc["rec_status"] = "1"
	paramsBankAcc["rec_order"] = "0"
	paramsBankAcc["rec_created_date"] = time.Now().Format(dateLayout)
	paramsBankAcc["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	//save bank account
	status, err, lastID := models.CreateMsBankAccount(paramsBankAcc)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//save product bank account
	paramsProBankAcc["rec_status"] = "1"
	paramsProBankAcc["bank_account_key"] = lastID
	paramsProBankAcc["rec_order"] = "0"
	paramsProBankAcc["rec_created_date"] = time.Now().Format(dateLayout)
	paramsProBankAcc["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.CreateMsProductBankAccount(paramsProBankAcc)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminMsProductBankAccount(c echo.Context) error {
	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	paramsProBankAcc := make(map[string]string)
	paramsBankAcc := make(map[string]string)

	//key /prod_bankacc_key
	key := c.FormValue("key")
	if key == "" {
		log.Error("Missing required parameter: key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key cann't be blank", "Missing required parameter: key cann't be blank")
	}
	strkey, err := strconv.ParseUint(key, 10, 64)
	if err == nil && strkey > 0 {
		paramsProBankAcc["prod_bankacc_key"] = key
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	//cek product bank account
	var productBankAccount models.MsProductBankAccount
	status, err = models.GetMsProductBankAccount(&productBankAccount, key)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	//product_key
	productkey := c.FormValue("product_key")
	if productkey == "" {
		log.Error("Missing required parameter: product_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key cann't be blank", "Missing required parameter: product_key cann't be blank")
	}
	strproductkey, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && strproductkey > 0 {
		paramsProBankAcc["product_key"] = productkey
	} else {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	//bank_key
	bankkey := c.FormValue("bank_key")
	if bankkey == "" {
		log.Error("Missing required parameter: bank_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key cann't be blank", "Missing required parameter: bank_key cann't be blank")
	}
	strbankkey, err := strconv.ParseUint(bankkey, 10, 64)
	if err == nil && strbankkey > 0 {
		paramsBankAcc["bank_key"] = bankkey
	} else {
		log.Error("Wrong input for parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key", "Missing required parameter: bank_key")
	}

	//account_no
	accountno := c.FormValue("account_no")
	if accountno == "" {
		log.Error("Missing required parameter: account_no cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_no cann't be blank", "Missing required parameter: account_no cann't be blank")
	}
	paramsBankAcc["account_no"] = accountno

	//account_holder_name
	accountholdername := c.FormValue("account_holder_name")
	if accountholdername == "" {
		log.Error("Missing required parameter: account_holder_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_holder_name cann't be blank", "Missing required parameter: account_holder_name cann't be blank")
	}
	paramsBankAcc["account_holder_name"] = accountholdername

	//branch_name
	branchname := c.FormValue("branch_name")
	if branchname != "" {
		paramsBankAcc["branch_name"] = branchname
	}

	//currency_key
	currencykey := c.FormValue("currency_key")
	if currencykey == "" {
		log.Error("Missing required parameter: currency_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key cann't be blank", "Missing required parameter: currency_key cann't be blank")
	}
	strcurrencykey, err := strconv.ParseUint(currencykey, 10, 64)
	if err == nil && strcurrencykey > 0 {
		paramsBankAcc["currency_key"] = currencykey
	} else {
		log.Error("Wrong input for parameter: currency_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: currency_key", "Missing required parameter: currency_key")
	}

	//bank_account_type
	bankaccounttype := c.FormValue("bank_account_type")
	if bankaccounttype == "" {
		log.Error("Missing required parameter: bank_account_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_type cann't be blank", "Missing required parameter: bank_account_type cann't be blank")
	}
	strbankaccounttype, err := strconv.ParseUint(bankaccounttype, 10, 64)
	if err == nil && strbankaccounttype > 0 {
		paramsBankAcc["bank_account_type"] = bankaccounttype
	} else {
		log.Error("Wrong input for parameter: bank_account_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_type", "Missing required parameter: bank_account_type")
	}

	//swift_code
	swiftcode := c.FormValue("swift_code")
	if swiftcode != "" {
		paramsBankAcc["swift_code"] = swiftcode
	}

	//bank_account_name
	bankaccountname := c.FormValue("bank_account_name")
	if bankaccountname == "" {
		log.Error("Missing required parameter: bank_account_name cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_name cann't be blank", "Missing required parameter: bank_account_name cann't be blank")
	}
	paramsProBankAcc["bank_account_name"] = bankaccountname

	//bank_account_type
	bankaccountpurpose := c.FormValue("bank_account_purpose")
	if bankaccountpurpose == "" {
		log.Error("Missing required parameter: bank_account_purpose cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_purpose cann't be blank", "Missing required parameter: bank_account_purpose cann't be blank")
	}
	strbankaccountpurpose, err := strconv.ParseUint(bankaccountpurpose, 10, 64)
	if err == nil && strbankaccountpurpose > 0 {
		paramsProBankAcc["bank_account_purpose"] = bankaccountpurpose
	} else {
		log.Error("Wrong input for parameter: bank_account_purpose")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_account_purpose", "Missing required parameter: bank_account_purpose")
	}

	strconv.FormatUint(lib.Profile.UserID, 10)

	dateLayout := "2006-01-02 15:04:05"
	paramsBankAcc["bank_account_key"] = strconv.FormatUint(*productBankAccount.BankAccountKey, 10)
	paramsBankAcc["rec_modified_date"] = time.Now().Format(dateLayout)
	paramsBankAcc["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	//save update bank account
	status, err = models.UpdateMsBankAccount(paramsBankAcc)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//save update product bank account
	paramsProBankAcc["rec_modified_date"] = time.Now().Format(dateLayout)
	paramsProBankAcc["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.UpdateMsProductBankAccount(paramsProBankAcc)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}
