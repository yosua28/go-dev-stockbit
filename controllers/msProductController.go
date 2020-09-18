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
	fundTypeKey, _ := strconv.ParseUint(fundTypeKeyStr, 10, 64)
	if fundTypeKey == 0 {
		log.Error("Fund Type should be number")
		return lib.CustomError(http.StatusNotFound,"Fund Type should be number","Fund Type should be number")
	}

	params["fund_type_key"] = fundTypeKeyStr

	// Get parameter order_by
	orderBy := c.QueryParam("order_by")
	if orderBy!=""{
		if (orderBy == "post_title") || (orderBy == "post_publish_thru") || (orderBy == "post_publish_start") {
			params["orderBy"] = orderBy
		}else{
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	}
	// Get parameter order_type
	orderType := c.QueryParam("order_type")
	if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
		params["orderType"] = orderType
	}

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
	for _, product := range productDB {
		if _, ok := lib.Find(productIDs, strconv.FormatUint(product.ProductKey, 10)); !ok {
			productIDs = append(productIDs, strconv.FormatUint(product.ProductKey, 10))
		}
		if _, ok := lib.Find(fundTypeIDs, strconv.FormatUint(*product.FundTypeKey, 10)); !ok {
			fundTypeIDs = append(fundTypeIDs, strconv.FormatUint(*product.FundTypeKey, 10))
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
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var performanceDB []models.FfsNavPerformance
	status, err = models.GetLastNavPerformanceIn(&performanceDB, productIDs)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	nData := make(map[uint64]models.TrNav)
	fData := make(map[uint64]models.MsFundType)
	pData := make(map[uint64]models.FfsNavPerformance)

	for _, nav := range navDB {
		nData[nav.ProductKey] = nav
	}
	for _, fundType := range fundTypeDB {
		fData[fundType.FundTypeKey] = fundType
	}
	for _, performance := range performanceDB {
		pData[performance.ProductKey] = performance
	}
	
	
	var productData []models.MsProductList
	for _, product := range productDB {
		var data models.MsProductList
	
		data.ProductKey = product.ProductKey
		data.ProductID = product.ProductID
		data.ProductCode = product.ProductCode
		data.ProductName = product.ProductName
		data.ProductNameAlt = product.ProductNameAlt
		data.MinSubAmount = product.MinSubAmount
		if product.RecImage1 != nil && *product.RecImage1 != ""{
			data.RecImage1 = config.BaseUrl + "/images/product/" + *product.RecImage1
		}else{
			data.RecImage1 = config.BaseUrl + "/images/product/default.png"
		}

		var fundType models.MsFundTypeInfo
		fundType.FundTypeKey = *product.FundTypeKey
		if _, ok := fData[*product.FundTypeKey]; ok {
			fundType.FundTypeCode = fData[*product.FundTypeKey].FundTypeCode
			fundType.FundTypeName = fData[*product.FundTypeKey].FundTypeName
		}
		data.FundType = &fundType

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		
		var nav models.TrNavInfo
		if _, ok := nData[product.ProductKey]; ok {
			date, _ := time.Parse(layout, nData[product.ProductKey].NavDate)
			fmt.Println(date)
			nav.NavDate = date.Format(newLayout)
			nav.NavValue = nData[product.ProductKey].NavValue
		}
		data.Nav = &nav
		
		var perform models.FfsNavPerformanceInfo
		if _, ok := pData[product.ProductKey]; ok {
			date, _ := time.Parse(layout, pData[product.ProductKey].NavDate)
			perform.NavDate = date.Format(newLayout)
			perform.D1 = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformD1) + `%`
			perform.MTD = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformMtd) + `%`
			perform.M1 = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformM1) + `%`
			perform.M3 = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformM3) + `%`
			perform.M6 = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformM6) + `%`
			perform.Y1 = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformY1) + `%`
			perform.Y3 = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformY3) + `%`
			perform.Y5 = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformY5) + `%`
			perform.YTD = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformYtd) + `%`
			perform.CAGR = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformCagr) + `%`
			perform.ALL = fmt.Sprintf("%.3f", pData[product.ProductKey].PerformAll) + `%`
		} 
		data.NavPerformance = &perform
	
		productData = append(productData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = productData
	
	return c.JSON(http.StatusOK, response)
}

func GetMsProductData(c echo.Context) error {
	var err error
	var status int

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

	var data models.MsProductData
	
	data.ProductKey = product.ProductKey
	data.ProductID = product.ProductID
	data.ProductCode = product.ProductCode
	data.ProductName = product.ProductName
	data.ProductNameAlt = product.ProductNameAlt
	data.MinSubAmount = product.MinSubAmount
	if product.RecImage1 != nil && *product.RecImage1 != ""{
		data.RecImage1 = config.BaseUrl + "/images/product/" + *product.RecImage1
	}else{
		data.RecImage1 = config.BaseUrl + "/images/product/default.png"
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
