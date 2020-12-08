package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	fundTypeKeyStr := c.QueryParam("fund_type")

	if fundTypeKeyStr != "" {
		fundTypeKey, _ := strconv.ParseUint(fundTypeKeyStr, 10, 64)
		if fundTypeKey == 0 {
			log.Error("Fund Type should be number")
			return lib.CustomError(http.StatusNotFound, "Fund Type should be number", "Fund Type should be number")
		}
		params["fund_type_key"] = fundTypeKeyStr
	}

	transaction := c.QueryParam("transaction")
	if !(transaction == "sub" || transaction == "red" || transaction == "all") {
		transaction = "all"
	}

	exceptStr := c.QueryParam("except")

	var except uint64 = 0
	if exceptStr != "" {
		except, _ = strconv.ParseUint(exceptStr, 10, 64)
	}
	var userProduct []string
	if lib.Profile.CustomerKey != nil && *lib.Profile.CustomerKey > 0 {
		paramsAcc := make(map[string]string)
		paramsAcc["customer_key"] = strconv.FormatUint(*lib.Profile.CustomerKey, 10)
		paramsAcc["rec_status"] = "1"

		var accDB []models.TrAccount
		status, err = models.GetAllTrAccount(&accDB, paramsAcc)
		if err != nil {
			log.Error(err.Error())
		}

		var accIDs []string
		accProduct := make(map[uint64]uint64)
		acaProduct := make(map[uint64]uint64)
		var acaDB []models.TrAccountAgent
		if len(accDB) > 0 {
			for _, acc := range accDB {
				accIDs = append(accIDs, strconv.FormatUint(acc.AccKey, 10))
				accProduct[acc.AccKey] = acc.ProductKey
			}
			status, err = models.GetTrAccountAgentIn(&acaDB, accIDs, "acc_key")
			if err != nil {
				log.Error(err.Error())
			}
			if len(acaDB) > 0 {
				var acaIDs []string
				for _, aca := range acaDB {
					acaIDs = append(acaIDs, strconv.FormatUint(aca.AcaKey, 10))
					acaProduct[aca.AcaKey] = aca.AccKey
				}
				var balanceDB []models.TrBalance
				status, err = models.GetLastBalanceIn(&balanceDB, acaIDs)
				if err != nil {
					log.Error(err.Error())
				}
				if len(balanceDB) > 0 {
					for _, balance := range balanceDB {
						if accKey, ok := acaProduct[balance.AcaKey]; ok {
							if balance.BalanceUnit > 0 {
								if productKey, ok := accProduct[accKey]; ok {
									userProduct = append(userProduct, strconv.FormatUint(productKey, 10))
								}
							}
						}
					}
				}
			}
		}
	}

	var performanceDB []models.FfsNavPerformance
	var responseOrder []string
	performData := false
	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		if (orderBy == "post_title") || (orderBy == "post_publish_thru") || (orderBy == "post_publish_start") || (orderBy == "rec_order") {
			params["orderBy"] = orderBy
			// Get parameter order_type
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else if (orderBy == "cagr") || (orderBy == "y5") || (orderBy == "y3") || (orderBy == "y1") || (orderBy == "m6") || (orderBy == "m3") || (orderBy == "ytd") {
			params := make(map[string]string)
			params["orderBy"] = "perform_" + orderBy
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
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "rec_order"
		params["orderType"] = "DESC"
	}

	params["flag_enabled"] = "1"
	params["rec_status"] = "1"
	params["flag_enabled"] = "1"

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
		_, ok := lib.Find(userProduct, strconv.FormatUint(product.ProductKey, 10))
		if ((!ok && transaction == "sub") || (ok && transaction == "red") || transaction == "all") && product.ProductKey != except {
			productData[strconv.FormatUint(product.ProductKey, 10)] = product
		}
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
		log.Error("Get risk profile: " + err.Error())
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

	var productNotAllow []models.ProductCheckAllowRedmSwtching
	if lib.Profile.CustomerKey != nil {
		if len(productIDs) > 0 {
			strCustomer := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
			status, err = models.CheckProductAllowRedmOrSwitching(&productNotAllow, strCustomer, productIDs)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	var productIDNotAllowed []string
	for _, pKey := range productNotAllow {
		if _, ok := lib.Find(productIDNotAllowed, strconv.FormatUint(pKey.ProductKey, 10)); !ok {
			productIDNotAllowed = append(productIDNotAllowed, strconv.FormatUint(pKey.ProductKey, 10))
		}
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

			if product.RecImage1 != nil && *product.RecImage1 != "" {
				data.RecImage1 = config.BaseUrl + "/images/product/" + *product.RecImage1
			} else {
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
				risk.RiskProfileKey = r.RiskProfileKey
				risk.RiskCode = r.RiskCode
				risk.RiskName = r.RiskName
				risk.RiskDesc = r.RiskDesc
			}
			data.RiskProfile = &risk

			// layout := "2006-01-02 15:04:05"
			// newLayout := "02 Jan 2006"

			var nav models.TrNavInfo
			if n, ok := nData[product.ProductKey]; ok {
				// date, _ := time.Parse(layout, n.NavDate)
				// nav.NavDate = date.Format(newLayout)
				nav.NavDate = n.NavDate
				nav.NavValue = n.NavValue
			}
			data.Nav = &nav

			var perform models.FfsNavPerformanceInfo
			if p, ok := pData[product.ProductKey]; ok {
				// date, _ := time.Parse(layout, p.NavDate)
				// perform.NavDate = date.Format(newLayout)
				perform.NavDate = p.NavDate
				perform.D1 = fmt.Sprintf("%.2f", p.PerformD1) + `%`
				perform.MTD = fmt.Sprintf("%.2f", p.PerformMtd) + `%`
				perform.M1 = fmt.Sprintf("%.2f", p.PerformM1) + `%`
				perform.M3 = fmt.Sprintf("%.2f", p.PerformM3) + `%`
				perform.M6 = fmt.Sprintf("%.2f", p.PerformM6) + `%`
				perform.Y1 = fmt.Sprintf("%.2f", p.PerformY1) + `%`
				perform.Y3 = fmt.Sprintf("%.2f", p.PerformY3) + `%`
				perform.Y5 = fmt.Sprintf("%.2f", p.PerformY5) + `%`
				perform.YTD = fmt.Sprintf("%.2f", p.PerformYtd) + `%`
				perform.CAGR = fmt.Sprintf("%.2f", p.PerformCagr) + `%`
				perform.ALL = fmt.Sprintf("%.2f", p.PerformAll) + `%`
			}
			data.NavPerformance = &perform

			data.IsAllowRedemption = true
			data.IsAllowSwitchin = true

			if _, ok := lib.Find(productIDNotAllowed, strconv.FormatUint(product.ProductKey, 10)); !ok {
				data.IsAllowRedemption = true
				data.IsAllowSwitchin = true
			} else {
				data.IsAllowRedemption = false
				data.IsAllowSwitchin = false
			}

			if product.FlagSwitchIn == 1 {
				data.IsAllowProductDestination = true
			} else {
				data.IsAllowProductDestination = false
			}

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

	var currencyDB models.MsCurrency
	status, err = models.GetMsCurrency(&currencyDB, strconv.FormatUint(*product.CurrencyKey, 10))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var currency models.MsCurrencyInfo
	currency.CurrencyKey = currencyDB.CurrencyKey
	currency.Code = currencyDB.Code
	currency.Symbol = currencyDB.Symbol
	currency.Name = currencyDB.Name
	currency.FlagBase = currencyDB.FlagBase

	data.Currency = currency

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
			if fee.FlagShowOntnc != nil && *fee.FlagShowOntnc == 1 {
				var data models.MsProductFeeInfo
				if fee.FeeAnnotation != nil {
					data.FeeAnnotation = *fee.FeeAnnotation
				}
				
				data.FlagShowOntnc = *fee.FlagShowOntnc
				
				if fee.FeeType != nil {
					data.FeeType = *fee.FeeType
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
			if bank.BankAccountPurpose == 269 {
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
			} else {
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
	data.MinRedAmount = product.MinRedAmount
	data.MinRedUnit = product.MinRedUnit
	data.MinUnitAfterRed = product.MinUnitAfterRed
	if ffsDB[0].FfsLink != nil && *ffsDB[0].FfsLink != "" {
		data.FundFactSheet = *ffsDB[0].FfsLink
	} else {
		data.FundFactSheet = "#"
	}
	if product.ProspectusLink != nil && *product.ProspectusLink != "" {
		data.ProspectusLink = *product.ProspectusLink
	} else {
		data.ProspectusLink = "#"
	}
	if product.RecImage1 != nil && *product.RecImage1 != "" {
		data.RecImage1 = config.BaseUrl + "/images/product/" + *product.RecImage1
	} else {
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

	// layout := "2006-01-02 15:04:05"
	// newLayout := "02 Jan 2006"

	var nav models.TrNavInfo
	// date, _ := time.Parse(layout, navDB[0].NavDate)
	// nav.NavDate = date.Format(newLayout)
	nav.NavDate = navDB[0].NavDate
	nav.NavValue = float32(math.Floor(float64(navDB[0].NavValue)*100) / 100)

	data.SubSuspend = false
	data.RedSuspend = false
	if lib.Profile.CustomerKey != nil && *lib.Profile.CustomerKey > 0 {
		paramsAcc := make(map[string]string)
		paramsAcc["customer_key"] = strconv.FormatUint(*lib.Profile.CustomerKey, 10)
		paramsAcc["product_key"] = strconv.FormatUint(product.ProductKey, 10)
		paramsAcc["rec_status"] = "1"

		var accDB []models.TrAccount
		status, err = models.GetAllTrAccount(&accDB, paramsAcc)
		if err != nil {
			log.Error(err.Error())
		}

		var accIDs []string
		accProduct := make(map[uint64]uint64)
		acaProduct := make(map[uint64]uint64)
		var acaDB []models.TrAccountAgent
		if len(accDB) > 0 {
			for _, acc := range accDB {
				accIDs = append(accIDs, strconv.FormatUint(acc.AccKey, 10))
				accProduct[acc.AccKey] = acc.ProductKey
				if acc.SubSuspendFlag != nil && *acc.SubSuspendFlag == 1 {
					data.SubSuspend = true
				}
				if acc.RedSuspendFlag != nil && *acc.RedSuspendFlag == 1 {
					data.RedSuspend = true
				}
			}
			status, err = models.GetTrAccountAgentIn(&acaDB, accIDs, "acc_key")
			if err != nil {
				log.Error(err.Error())
			}
			if len(acaDB) > 0 {
				var acaIDs []string
				for _, aca := range acaDB {
					acaIDs = append(acaIDs, strconv.FormatUint(aca.AcaKey, 10))
					acaProduct[aca.AcaKey] = aca.AccKey
				}
				var balanceDB []models.TrBalance
				status, err = models.GetLastBalanceIn(&balanceDB, acaIDs)
				if err != nil {
					log.Error(err.Error())
				}
				if len(balanceDB) > 0 {
					for _, balance := range balanceDB {
						data.BalanceUnit += balance.BalanceUnit
					}
					var invest float32
					invest = nav.NavValue * data.BalanceUnit
					data.InvestValue = fmt.Sprintf("%v", math.Trunc(float64(invest)))
				}
			}
		}
	}
	data.BalanceUnit = float32(math.Floor(float64(data.BalanceUnit)*100) / 100)
	data.Nav = &nav

	var perform models.FfsNavPerformanceInfo
	// date, _ = time.Parse(layout, performanceDB[0].NavDate)
	// perform.NavDate = date.Format(newLayout)
	perform.NavDate = performanceDB[0].NavDate
	perform.D1 = fmt.Sprintf("%.2f", performanceDB[0].PerformD1) + `%`
	perform.MTD = fmt.Sprintf("%.2f", performanceDB[0].PerformMtd) + `%`
	perform.M1 = fmt.Sprintf("%.2f", performanceDB[0].PerformM1) + `%`
	perform.M3 = fmt.Sprintf("%.2f", performanceDB[0].PerformM3) + `%`
	perform.M6 = fmt.Sprintf("%.2f", performanceDB[0].PerformM6) + `%`
	perform.Y1 = fmt.Sprintf("%.2f", performanceDB[0].PerformY1) + `%`
	perform.Y3 = fmt.Sprintf("%.2f", performanceDB[0].PerformY3) + `%`
	perform.Y5 = fmt.Sprintf("%.2f", performanceDB[0].PerformY5) + `%`
	perform.YTD = fmt.Sprintf("%.2f", performanceDB[0].PerformYtd) + `%`
	perform.CAGR = fmt.Sprintf("%.2f", performanceDB[0].PerformCagr) + `%`
	perform.ALL = fmt.Sprintf("%.2f", performanceDB[0].PerformAll) + `%`
	data.NavPerformance = &perform

	var risk models.MsRiskProfileInfo
	risk.RiskProfileKey = riskDB.RiskProfileKey
	risk.RiskCode = riskDB.RiskCode
	risk.RiskName = riskDB.RiskName
	risk.RiskDesc = riskDB.RiskDesc

	data.RiskProfile = &risk

	var custodian models.MsCustodianBankInfo
	custodian.CustodianCode = custodianDB.CustodianCode
	custodian.CustodianShortName = custodianDB.CustodianShortName
	custodian.CustodianFullName = custodianDB.CustodianFullName

	data.CustodianBank = &custodian

	var countData models.CountData

	if lib.Profile.CustomerKey != nil {
		paramsCekTrans := make(map[string]string)
		paramsCekTrans["rec_status"] = "1"
		paramsCekTrans["product_key"] = strconv.FormatUint(product.ProductKey, 10)
		customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
		paramsCekTrans["customer_key"] = customerKey
		var transTypeKey []string
		transTypeKey = append(transTypeKey, "1")
		transTypeKey = append(transTypeKey, "4")
		status, err = models.AdminGetCountTrTransaction(&countData, paramsCekTrans, transTypeKey, "trans_type_key")
		if err == nil {
			if int(countData.CountData) > 0 {
				data.IsNew = false
				data.TncIsNew = ""
			} else {
				data.IsNew = true
				var scApp models.ScAppConfig
				status, err = models.GetScAppConfigByCode(&scApp, "NEW_PRODUCT_SUBSCRIBE")
				if err != nil {
					data.TncIsNew = ""
				} else {
					str1 := scApp.AppConfigValue
					res1 := strings.Replace(*str1, "#ProductName#", product.ProductName, 1)
					data.TncIsNew = res1
				}
			}
		} else {
			data.IsNew = false
			data.TncIsNew = ""
		}
	} else {
		data.IsNew = false
		data.TncIsNew = ""
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data

	return c.JSON(http.StatusOK, response)
}

func Portofolio(c echo.Context) error {
	var err error
	var status int

	responseData := make(map[string]interface{})
	params := make(map[string]string)

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}
	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	params["customer_key"] = customerKey
	params["trans_status_key"] = "9"
	var transactionDB []models.TrTransaction
	status, err = models.GetAllTrTransaction(&transactionDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get transaction data")
	}
	if len(transactionDB) < 1 {
		log.Error("Transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var transactionIDs []string
	var productIDs []string
	for _, transaction := range transactionDB {
		transactionIDs = append(transactionIDs, strconv.FormatUint(transaction.TransactionKey, 10))
		productIDs = append(productIDs, strconv.FormatUint(transaction.ProductKey, 10))
	}
	var productDB []models.MsProduct
	status, err = models.GetMsProductIn(&productDB, productIDs, "product_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get product data")
	}
	if len(productDB) < 1 {
		log.Error("Product data not found")
		return lib.CustomError(http.StatusNotFound, "Product data not found", "Product data not found")
	}

	var currencyIDs []string
	productData := make(map[uint64]uint64)
	netSubProduct := make(map[uint64]float32)
	totalProduct := make(map[uint64]float32)
	for _, product := range productDB {
		currencyIDs = append(currencyIDs, strconv.FormatUint(*product.CurrencyKey, 10))
		productData[product.ProductKey] = *product.CurrencyKey
		netSubProduct[product.ProductKey] = 0
		totalProduct[product.ProductKey] = 0
	}

	var currencyDB []models.MsCurrency
	status, err = models.GetMsCurrencyIn(&currencyDB, currencyIDs, "currency_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get currency data")
	}

	ccy := make(map[uint64]string)
	for _, currency := range currencyDB {
		ccy[currency.CurrencyKey] = currency.Code
	}

	var rateDB []models.TrCurrencyRate
	status, err = models.GetLastCurrencyIn(&rateDB, currencyIDs)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get currency rate data")
	}

	rateData := make(map[uint64]float32)
	rateData[1] = 1
	for _, rate := range rateDB {
		rateData[rate.CurrencyKey] = rate.RateValue
	}

	var tcDB []models.TrTransactionConfirmation
	status, err = models.GetTrTransactionConfirmationIn(&tcDB, transactionIDs, "transaction_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get TC data")
	}
	if len(tcDB) < 1 {
		log.Error("TC data not found")
		return lib.CustomError(http.StatusNotFound, "TC data not found", "TC data not found")
	}

	tcData := make(map[uint64]models.TrTransactionConfirmation)
	for _, tc := range tcDB {
		tcData[tc.TransactionKey] = tc
	}
	var netSub float32 = 0
	for _, transaction := range transactionDB {
		if tc, ok := tcData[transaction.TransactionKey]; ok {
			if transaction.TransTypeKey == 1 || transaction.TransTypeKey == 4 {
				netSub += (tc.ConfirmedAmount * rateData[productData[transaction.ProductKey]])
				netSubProduct[transaction.ProductKey] += (tc.ConfirmedAmount * rateData[productData[transaction.ProductKey]])
				log.Info("NETSUB + #", transaction.ProductKey, "#", fmt.Sprintf("%.3f", (tc.ConfirmedAmount*rateData[productData[transaction.ProductKey]])))
			} else {
				netSub -= (tc.ConfirmedAmount * rateData[productData[transaction.ProductKey]])
				netSubProduct[transaction.ProductKey] -= (tc.ConfirmedAmount * rateData[productData[transaction.ProductKey]])
				log.Info("NETSUB - #", transaction.ProductKey, "#", fmt.Sprintf("%.3f", (tc.ConfirmedAmount*rateData[productData[transaction.ProductKey]])))
			}
		}
	}

	responseData["net_sub"] = float32(math.Trunc(float64(netSub)))

	params = make(map[string]string)
	var userProduct []string
	balanceUnit := make(map[uint64]float32)
	avgNav := make(map[uint64]float32)
	suspend := make(map[uint64]bool)
	if lib.Profile.CustomerKey != nil && *lib.Profile.CustomerKey > 0 {
		paramsAcc := make(map[string]string)
		paramsAcc["customer_key"] = strconv.FormatUint(*lib.Profile.CustomerKey, 10)
		paramsAcc["rec_status"] = "1"

		var accDB []models.TrAccount
		status, err = models.GetAllTrAccount(&accDB, paramsAcc)
		if err != nil {
			log.Error(err.Error())
		}

		var accIDs []string
		accProduct := make(map[uint64]uint64)
		acaProduct := make(map[uint64]uint64)
		var acaDB []models.TrAccountAgent
		if len(accDB) > 0 {
			for _, acc := range accDB {
				accIDs = append(accIDs, strconv.FormatUint(acc.AccKey, 10))
				accProduct[acc.AccKey] = acc.ProductKey
				if (acc.SubSuspendFlag != nil && *acc.SubSuspendFlag == 1) || 
				(acc.RedSuspendFlag != nil && *acc.RedSuspendFlag == 1) {
					suspend[acc.ProductKey] = true
				}else{
					suspend[acc.ProductKey] = false
				}
			}
			status, err = models.GetTrAccountAgentIn(&acaDB, accIDs, "acc_key")
			if err != nil {
				log.Error(err.Error())
			}
			if len(acaDB) > 0 {
				var acaIDs []string
				for _, aca := range acaDB {
					acaIDs = append(acaIDs, strconv.FormatUint(aca.AcaKey, 10))
					acaProduct[aca.AcaKey] = aca.AccKey
				}
				var balanceDB []models.TrBalance
				status, err = models.GetLastBalanceIn(&balanceDB, acaIDs)
				if err != nil {
					log.Error(err.Error())
				}
				if len(balanceDB) > 0 {
					for _, balance := range balanceDB {
						log.Info(balance.BalanceKey, balance.AcaKey)
						if accKey, ok := acaProduct[balance.AcaKey]; ok {
							if balance.BalanceUnit > 0 {
								if productKey, ok := accProduct[accKey]; ok {
									if _, ok := balanceUnit[productKey]; ok {
										balanceUnit[productKey] += balance.BalanceUnit
									} else {
										balanceUnit[productKey] = balance.BalanceUnit
									}
									avgNav[productKey] = *balance.AvgNav
									userProduct = append(userProduct, strconv.FormatUint(productKey, 10))
								}
							}
						}
					}
				}
			}
		}
	}

	productDB = nil
	status, err = models.GetMsProductIn(&productDB, userProduct, "product_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(productDB) < 1 {
		log.Error("product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}
	productData = make(map[uint64]uint64)
	for _, product := range productDB {
		productData[product.ProductKey] = *product.CurrencyKey
	}

	var navDB2 []models.TrNav
	status, err = models.GetLastNavIn(&navDB2, userProduct)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	navData := make(map[uint64]models.TrNav)
	var total float32
	for _, nav := range navDB2 {
		navData[nav.ProductKey] = nav
		if b, ok := balanceUnit[nav.ProductKey]; ok {
			total += ((b * nav.NavValue) * rateData[productData[nav.ProductKey]])
			totalProduct[nav.ProductKey] += ((b * nav.NavValue) * rateData[productData[nav.ProductKey]])
			log.Info("TOTAL#", nav.ProductKey, "#", fmt.Sprintf("%.2f", ((b*nav.NavValue)*rateData[productData[nav.ProductKey]])))
		}
	}
	responseData["total_invest"] = float32(math.Trunc(float64(total)))

	imba := ((total - netSub) / netSub) * 100
	responseData["imba"] = fmt.Sprintf("%.2f", imba) + `%`
	var products []interface{}
	var portofolio models.Portofolio
	p := message.NewPrinter(language.Indonesian)
	portofolio.Total = p.Sprintf("%v", float32(math.Trunc(float64(total))))

	var portofolioDatas []models.ProductPortofolio
	var totalGainLoss float32
	for _, product := range productDB {
		data := make(map[string]interface{})
		var portofolioData models.ProductPortofolio

		data["product_key"] = product.ProductKey
		data["product_id"] = product.ProductID
		if _, ok := suspend[product.ProductKey]; ok {
			data["suspend"] = suspend[product.ProductKey]
		} else {
			data["suspend"] = false
		}
		data["product_code"] = product.ProductCode
		data["product_name"] = product.ProductName
		data["product_name_alt"] = product.ProductNameAlt
		imba := ((totalProduct[product.ProductKey] - netSubProduct[product.ProductKey]) / netSubProduct[product.ProductKey]) * 100
		data["imba"] = fmt.Sprintf("%.2f", imba) + `%`
		portofolioData.ProductName = product.ProductNameAlt
		portofolioData.CCY = ccy[*product.CurrencyKey]
		portofolioData.AvgNav = p.Sprintf("%.2f", avgNav[product.ProductKey])
		portofolioData.Kurs = p.Sprintf("%v", float32(math.Trunc(float64(rateData[*product.CurrencyKey]))))

		if product.RecImage1 != nil && *product.RecImage1 != "" {
			data["rec_image1"] = config.BaseUrl + "/images/product/" + *product.RecImage1
		} else {
			data["rec_image1"] = config.BaseUrl + "/images/product/default.png"
		}
		if n, ok := navData[product.ProductKey]; ok {
			portofolioData.Nav = p.Sprintf("%.2f", n.NavValue)
			if b, ok := balanceUnit[product.ProductKey]; ok {
				data["invest_value"] = float32(math.Trunc(float64((b * n.NavValue) * rateData[*product.CurrencyKey])))
				portofolioData.Amount = p.Sprintf("%v", float32(math.Trunc(float64(b*n.NavValue))))
				portofolioData.AmountIDR = p.Sprintf("%v", float32(math.Trunc(float64(data["invest_value"].(float32)))))
				portofolioData.Unit = p.Sprintf("%.2f", b)
				gainLoss := (avgNav[product.ProductKey] - n.NavValue) * b
				portofolioData.GainLoss = p.Sprintf("%v", float32(math.Trunc(float64(gainLoss))))
				totalGainLoss += (gainLoss * rateData[*product.CurrencyKey])
				portofolioData.GainLossIDR = p.Sprintf("%v", float32(math.Trunc(float64(gainLoss*rateData[*product.CurrencyKey]))))
				percent := (((b * n.NavValue) * rateData[*product.CurrencyKey]) / total) * 100
				data["percent"] = fmt.Sprintf("%.2f", percent) + `%`
			}
		}
		portofolioDatas = append(portofolioDatas, portofolioData)
		products = append(products, data)
	}

	// PDF Template
	portofolio.TotalGainLoss = p.Sprintf("%v", float32(math.Trunc(float64(totalGainLoss))))
	portofolio.Datas = portofolioDatas
	var customer models.MsCustomer
	status, err = models.GetMsCustomer(&customer, customerKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, "Failed get customer", "Customer not found")
	}
	dateLayout := "2006-01-02 15:04:05"
	portofolio.Date = time.Now().Format(dateLayout)
	portofolio.Cif = customer.UnitHolderIDno
	sid := ""
	if customer.SidNo != nil {
		sid = *customer.SidNo
	}
	portofolio.Sid = sid
	portofolio.Name = customer.FullName

	params = make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["orderBy"] = "oa_request_key"
	params["orderType"] = "DESC"
	var requestDB []models.OaRequest
	status, err = models.GetAllOaRequest(&requestDB, 100, 0, true, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, "Failed get request", "Failed get request")
	}

	request := requestDB[0]

	var personalData models.OaPersonalData
	status, err = models.GetOaPersonalData(&personalData, strconv.FormatUint(request.OaRequestKey, 10), "oa_request_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, "Failed get personal data", "Failed get personal data")
	}

	var country models.MsCountry
	status, err = models.GetMsCountry(&country, strconv.FormatUint(personalData.Nationality, 10))
	if err != nil {
		log.Error(err.Error())
	}

	portofolio.Country = country.CouName

	var address models.OaPostalAddress
	status, err = models.GetOaPostalAddress(&address, strconv.FormatUint(*personalData.IDcardAddressKey, 10))
	if err != nil {
		log.Error(err.Error())
	}

	portofolio.Address = *address.AddressLine1

	var city models.MsCity
	status, err = models.GetMsCity(&city, strconv.FormatUint(*address.KabupatenKey, 10))
	if err != nil {
		log.Error(err.Error())
	}

	postalcode := ""
	if address.PostalCode != nil {
		postalcode = *address.PostalCode
	}

	portofolio.City = city.CityName + " " + postalcode

	t := template.New("account-statement-template.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/account-statement-template.html")
	if err != nil {
		log.Println(err)
	}
	f, err := os.Create(config.BasePath + "/mail/account-statement-" + strconv.FormatUint(lib.Profile.UserID, 10) + ".html")
	if err != nil {
		log.Println("create file: ", err)
	}
	if err := t.Execute(f, portofolio); err != nil {
		log.Println(err)
	}

	f.Close()
	// End PDF Template

	responseData["product"] = products
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func ProductListMutasi(c echo.Context) error {
	var err error
	var status int
	params := make(map[string]string)

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}
	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	params["customer_key"] = customerKey
	params["trans_status_key"] = "9"
	var transactionDB []models.TrTransaction
	status, err = models.GetAllTrTransaction(&transactionDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get transaction data")
	}
	if len(transactionDB) < 1 {
		log.Error("Transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var productIDs []string
	for _, transaction := range transactionDB {
		productIDs = append(productIDs, strconv.FormatUint(transaction.ProductKey, 10))
	}

	var productDB []models.MsProduct
	status, err = models.GetMsProductIn(&productDB, productIDs, "product_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(productDB) < 1 {
		log.Error("product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}

	var products []interface{}
	for _, product := range productDB {
		data := make(map[string]interface{})

		data["product_key"] = product.ProductKey
		data["product_id"] = product.ProductID
		data["product_code"] = product.ProductCode
		data["product_name"] = product.ProductName
		data["product_name_alt"] = product.ProductNameAlt

		products = append(products, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = products

	return c.JSON(http.StatusOK, response)
}
