package controllers

import (
	"api/models"
	"api/config"
	"api/lib"
	"net/http"
	"strconv"
	"time"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

func GetMsProductList(c echo.Context) error {
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
			return lib.CustomError(http.StatusBadRequest,"Page should be number", "Page should be number")
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
			return lib.CustomError(http.StatusBadRequest,"Nolimit parameter should be true/false","Nolimit parameter should be true/false")
		}
	} else {
		noLimit = false
	}

	params := make(map[string]string)
	fundTypeKeyStr := c.QueryParam("fund_type")

	if fundTypeKeyStr != ""{
		fundTypeKey, _ := strconv.ParseUint(fundTypeKeyStr, 10, 64)
		if fundTypeKey == 0 {
			log.Error("Fund Type should be number")
			return lib.CustomError(http.StatusNotFound,"Fund Type should be number","Fund Type should be number")
		}
		params["fund_type_key"] = fundTypeKeyStr
	}

	var performanceDB []models.FfsNavPerformance
	var responseOrder []string
	performData := false
	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy!=""{
		if (orderBy == "post_title") || (orderBy == "post_publish_thru") || (orderBy == "post_publish_start") {
			params["orderBy"] = orderBy
			// Get parameter order_type
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		}else if (orderBy == "cagr") || (orderBy == "y5") || (orderBy == "y3") || (orderBy == "y1") || (orderBy == "m6") || (orderBy == "m3") || (orderBy == "ytd"){
			params := make(map[string]string)
			params["orderBy"] = "perform_"+orderBy
			params["orderType"] = "DESC"
			status, err = models.GetAllLastNavPerformance(&performanceDB, params)
			if err != nil {
				log.Error("Get performance: " + err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}

			for _, perform := range performanceDB {
				responseOrder = append(responseOrder, strconv.FormatUint(perform.ProductKey, 10))
			}
			performData = true
		}else{
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	}

	params["flag_enabled"] = "1"
	params["rec_status"] = "1"

	var productDB []models.MsProduct
	status, err = models.GetAllMsProduct(&productDB, limit, offset, params, noLimit)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(productDB) < 1 {
		log.Error("product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}

	var productIDs []string
	var fundTypeIDs []string
	var riskIDs []string
	productData := make(map[string]models.MsProduct)
	for _, product := range productDB {
		productData[strconv.FormatUint(product.ProductKey, 10)] = product
		if _, ok := lib.Find(productIDs, strconv.FormatUint(product.ProductKey, 10)); !ok {
			productIDs = append(productIDs, strconv.FormatUint(product.ProductKey, 10))
		}
		if _, ok := lib.Find(fundTypeIDs, strconv.FormatUint(*product.FundTypeKey, 10)); !ok {
			fundTypeIDs = append(fundTypeIDs, strconv.FormatUint(*product.FundTypeKey, 10))
		}
		if _, ok := lib.Find(riskIDs, strconv.FormatUint(*product.RiskProfileKey, 10)); !ok {
			riskIDs = append(riskIDs, strconv.FormatUint(*product.RiskProfileKey, 10))
		}
	}

	var navDB []models.TrNav
	status, err = models.GetLastNavIn(&navDB, productIDs)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var fundTypeDB []models.MsFundType
	status, err = models.GetMsFundTypeIn(&fundTypeDB, fundTypeIDs, "fund_type_key")
	if err != nil {
		log.Error("Get fund type: " + err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if !performData {
		status, err = models.GetLastNavPerformanceIn(&performanceDB, productIDs)
		if err != nil {
			log.Error("Get performance: " + err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		responseOrder = productIDs
	}
	var riskDB []models.MsRiskProfile
	status, err = models.GetMsRiskProfileIn(&riskDB, riskIDs)
	if err != nil {
		log.Error("Get risk profile: "+err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	nData := make(map[uint64]models.TrNav)
	fData := make(map[uint64]models.MsFundType)
	pData := make(map[uint64]models.FfsNavPerformance)
	rData := make(map[uint64]models.MsRiskProfile)

	for _, nav := range navDB {
		nData[nav.ProductKey] = nav
	}
	for _, fundType := range fundTypeDB {
		fData[fundType.FundTypeKey] = fundType
	}
	for _, performance := range performanceDB {
		pData[performance.ProductKey] = performance
	}
	for _, risk := range riskDB {
		rData[risk.RiskProfileKey] = risk
	}
	
	
	var responseData []models.MsProductList
	for _, order := range responseOrder {
		if product, ok := productData[order]; ok {
			var data models.MsProductList
		
			data.ProductKey = product.ProductKey
			data.ProductID = product.ProductID
			data.ProductCode = product.ProductCode
			data.ProductName = product.ProductName
			data.ProductNameAlt = product.ProductNameAlt
			data.MinSubAmount = product.MinSubAmount
			data.MinSubAmount = product.MinSubAmount

			if product.RecImage1 != nil && *product.RecImage1 != ""{
				data.RecImage1 = config.BaseUrl + "/images/product/" + *product.RecImage1
			}else{
				data.RecImage1 = config.BaseUrl + "/images/product/default.png"
			}

			var fundType models.MsFundTypeInfo
			fundType.FundTypeKey = *product.FundTypeKey
			if ft, ok := fData[*product.FundTypeKey]; ok {
				fundType.FundTypeCode = ft.FundTypeCode
				fundType.FundTypeName = ft.FundTypeName
			}
			data.FundType = &fundType

			var risk models.MsRiskProfileInfo
			if r, ok := rData[*product.RiskProfileKey]; ok {
				risk.RiskCode = r.RiskCode
				risk.RiskName = r.RiskName
				risk.RiskDesc = r.RiskDesc
			}
			data.RiskProfile = &risk

			layout := "2006-01-02 15:04:05"
			newLayout := "02 Jan 2006"
			
			var nav models.TrNavInfo
			if n, ok := nData[product.ProductKey]; ok {
				date, _ := time.Parse(layout, n.NavDate)
				nav.NavDate = date.Format(newLayout)
				nav.NavValue = n.NavValue
			}
			data.Nav = &nav
			
			var perform models.FfsNavPerformanceInfo
			if p, ok := pData[product.ProductKey]; ok {
				date, _ := time.Parse(layout, p.NavDate)
				perform.NavDate = date.Format(newLayout)
				perform.D1 = fmt.Sprintf("%.3f", p.PerformD1) + `%`
				perform.MTD = fmt.Sprintf("%.3f", p.PerformMtd) + `%`
				perform.M1 = fmt.Sprintf("%.3f", p.PerformM1) + `%`
				perform.M3 = fmt.Sprintf("%.3f", p.PerformM3) + `%`
				perform.M6 = fmt.Sprintf("%.3f", p.PerformM6) + `%`
				perform.Y1 = fmt.Sprintf("%.3f", p.PerformY1) + `%`
				perform.Y3 = fmt.Sprintf("%.3f", p.PerformY3) + `%`
				perform.Y5 = fmt.Sprintf("%.3f", p.PerformY5) + `%`
				perform.YTD = fmt.Sprintf("%.3f", p.PerformYtd) + `%`
				perform.CAGR = fmt.Sprintf("%.3f", p.PerformCagr) + `%`
				perform.ALL = fmt.Sprintf("%.3f", p.PerformAll) + `%`
			} 
			data.NavPerformance = &perform
		
			responseData = append(responseData, data)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	
	return c.JSON(http.StatusOK, response)
}

func GetMsProductData(c echo.Context) error {
	var err error
	var status int
	var data models.MsProductData

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var product models.MsProduct
	status, err = models.GetMsProduct(&product, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	productIDs := []string{strconv.FormatUint(product.ProductKey, 10)}

	var navDB []models.TrNav
	status, err = models.GetLastNavIn(&navDB, productIDs)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var ffsDB []models.FfsPublish
	status, err = models.GetLastFfsIn(&ffsDB, productIDs)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var performanceDB []models.FfsNavPerformance
	status, err = models.GetLastNavPerformanceIn(&performanceDB, productIDs)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var custodianDB models.MsCustodianBank
	status, err = models.GetMsCustodianBank(&custodianDB, strconv.FormatUint(*product.CustodianKey, 10))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var riskDB models.MsRiskProfile
	status, err = models.GetMsRiskProfile(&riskDB, strconv.FormatUint(*product.RiskProfileKey, 10))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var feeDB []models.MsProductFee
	params := make(map[string]string)
	params["product_key"] = strconv.FormatUint(product.ProductKey, 10)
	status, err = models.GetAllMsProductFee(&feeDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(feeDB) > 0 {
		var feeIDs []string
		for _, fee := range feeDB {
			feeIDs = append(feeIDs, strconv.FormatUint(fee.FeeKey, 10))
		}

		var feeItemDB []models.MsProductFeeItem
		params = make(map[string]string)
		status, err = models.GetMsProductFeeItemIn(&feeItemDB, feeIDs, "product_fee_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}

		feeItemData := make(map[uint64][]models.MsProductFeeItemInfo)
		if len(feeItemDB) > 0 {
			for _, feeItem := range feeItemDB {
				var data models.MsProductFeeItemInfo
				data.ItemSeqno = feeItem.ItemSeqno
				data.RowMax = feeItem.RowMax
				data.PrincipleLimit = feeItem.PrincipleLimit
				data.FeeValue = feeItem.FeeValue
				if feeItem.ItemNotes != nil {
					data.ItemNotes = *feeItem.ItemNotes
				}

				feeItemData[feeItem.ProductFeeKey] = append(feeItemData[feeItem.ProductFeeKey], data)
			}
		}

		var feeData []models.MsProductFeeInfo
		for _, fee := range feeDB {
			var data models.MsProductFeeInfo
			if fee.FeeAnnotation != nil {
				data.FeeAnnotation = *fee.FeeAnnotation
			}
			if fee.FeeDesc != nil {
				data.FeeDesc = *fee.FeeDesc
			}
			if fee.FeeCode != nil {
				data.FeeCode = *fee.FeeCode
			}
			if item, ok := feeItemData[fee.FeeKey]; ok {
				data.FeeItem = item
			}

			feeData = append(feeData, data)
		}
		data.ProductFee = feeData
	}

	var productBankDB []models.MsProductBankAccount
	params = make(map[string]string)
	params["product_key"] = strconv.FormatUint(product.ProductKey, 10)
	status, err = models.GetAllMsProductBankAccount(&productBankDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(productBankDB) > 0 {
		var accountIDs []string
		for _, bank := range productBankDB {
			if bank.BankAccountKey != nil {
				accountIDs = append(accountIDs, strconv.FormatUint(*bank.BankAccountKey, 10))
			}
		}

		var bankAccountDB []models.MsBankAccount
		status, err = models.GetMsBankAccountIn(&bankAccountDB, accountIDs, "bank_account_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}

		var bankIDs []string
		for _, account := range bankAccountDB {
			bankIDs = append(bankIDs, strconv.FormatUint(account.BankKey, 10))
		}

		var bankDB []models.MsBank
		status, err = models.GetMsBankIn(&bankDB, bankIDs, "bank_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		bankData := make(map[uint64]models.MsBank)
		for _, bank := range bankDB {
			bankData[bank.BankKey] = bank
		}

		accountData := make(map[uint64]models.BankAccount)
		for _, account := range bankAccountDB {
			var data models.BankAccount
			data.AccountNo = account.AccountNo
			data.AccountHolderName = account.AccountHolderName
			data.BranchName = account.BranchName
			if acc, ok := bankData[account.BankKey]; ok {
				data.BankName = acc.BankName
			}

			accountData[account.BankAccountKey] = data
		}

		var productBankData []models.MsProductBankAccountInfo
		for _, bank := range productBankDB {
			var data models.MsProductBankAccountInfo
			data.BankAccountName = bank.BankAccountName
			data.BankAccountPurpose = bank.BankAccountPurpose
			if bank.BankAccountKey != nil {
				if acc, ok := accountData[*bank.BankAccountKey]; ok {
					data.BankAccount = acc
				}
			}
			
			productBankData = append(productBankData, data)
		}
		data.BankAcc = productBankData
	}

	params = make(map[string]string)
	params["config_type_key"] = "6"
	var chargesDB []models.ScAppConfig
	status, err = models.GetAllScAppConfig(&chargesDB, params)
	if err != nil {
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	if len(chargesDB) > 0 {
		for _, charges := range chargesDB {
			if charges.AppConfigValue != nil {
				if charges.AppConfigKey == 21 { 
					data.FeeTransfer = *charges.AppConfigValue
				}
				if charges.AppConfigKey == 22 {
					data.FeeService = *charges.AppConfigValue
				}
			}else{
				data.FeeTransfer = "0"
				data.FeeService = "0"
			}
		}
	}
	
	data.ProductKey = product.ProductKey
	data.ProductID = product.ProductID
	data.ProductCode = product.ProductCode
	data.ProductName = product.ProductName
	data.ProductNameAlt = product.ProductNameAlt
	data.MinSubAmount = product.MinSubAmount
	if ffsDB[0].FfsLink != nil && *ffsDB[0].FfsLink != ""{
		data.FundFactSheet = *ffsDB[0].FfsLink
	}else{
		data.FundFactSheet = "#"
	}
	if product.ProspectusLink != nil && *product.ProspectusLink != ""{
		data.ProspectusLink = *product.ProspectusLink
	}else{
		data.ProspectusLink = "#"
	}
	if product.RecImage1 != nil && *product.RecImage1 != ""{
		data.RecImage1 = config.BaseUrl + "/images/product/" + *product.RecImage1
	}else{
		data.RecImage1 = config.BaseUrl + "/images/product/default.png"
	}
	data.FlagSubscription = false
	if product.FlagSubscription == 1 {
		data.FlagSubscription = true
	}
	data.FlagRedemption = false 
	if product.FlagRedemption == 1 {
		data.FlagRedemption = true
	}
	data.FlagSwitchOut = false 
	if product.FlagSwitchOut == 1 {
		data.FlagSwitchOut = true
	}
	data.FlagSwitchIn = false 
	if product.FlagSwitchIn == 1 {
		data.FlagSwitchIn = true
	}
	
	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"

	var nav models.TrNavInfo
	date, _ := time.Parse(layout, navDB[0].NavDate)
	nav.NavDate = date.Format(newLayout)
	nav.NavValue = navDB[0].NavValue
	
	data.Nav = &nav

	var perform models.FfsNavPerformanceInfo
	date, _ = time.Parse(layout, performanceDB[0].NavDate)
	perform.NavDate = date.Format(newLayout)
	perform.D1 = fmt.Sprintf("%.3f", performanceDB[0].PerformD1) + `%`
	perform.MTD = fmt.Sprintf("%.3f", performanceDB[0].PerformMtd) + `%`
	perform.M1 = fmt.Sprintf("%.3f", performanceDB[0].PerformM1) + `%`
	perform.M3 = fmt.Sprintf("%.3f", performanceDB[0].PerformM3) + `%`
	perform.M6 = fmt.Sprintf("%.3f", performanceDB[0].PerformM6) + `%`
	perform.Y1 = fmt.Sprintf("%.3f", performanceDB[0].PerformY1) + `%`
	perform.Y3 = fmt.Sprintf("%.3f", performanceDB[0].PerformY3) + `%`
	perform.Y5 = fmt.Sprintf("%.3f", performanceDB[0].PerformY5) + `%`
	perform.YTD = fmt.Sprintf("%.3f", performanceDB[0].PerformYtd) + `%`
	perform.CAGR = fmt.Sprintf("%.3f", performanceDB[0].PerformCagr) + `%`
	perform.ALL = fmt.Sprintf("%.3f", performanceDB[0].PerformAll) + `%`
	data.NavPerformance = &perform

	var risk models.MsRiskProfileInfo
	risk.RiskCode = riskDB.RiskCode
	risk.RiskName = riskDB.RiskName
	risk.RiskDesc = riskDB.RiskDesc

	data.RiskProfile = &risk

	var custodian models.MsCustodianBankInfo
	custodian.CustodianCode = custodianDB.CustodianCode
	custodian.CustodianShortName = custodianDB.CustodianShortName
	custodian.CustodianFullName = custodianDB.CustodianFullName

	data.CustodianBank = &custodian

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data

	return c.JSON(http.StatusOK, response)
}
