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

func initAuthHoIt() error {
	var roleKeyHoIt uint64
	roleKeyHoIt = 15

	if lib.Profile.RoleKey != roleKeyHoIt {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func GetListProductAdmin(c echo.Context) error {

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

	items := []string{"product_key", "product_id", "product_code", "product_name", "launch_date", "inception_date", "isin_code", "flag_syariah", "sinvest_fund_code"}

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
		params["orderBy"] = "product_key"
		params["orderType"] = "ASC"
	}

	params["rec_status"] = "1"

	paramsLike := make(map[string]string)

	productName := c.QueryParam("product_name")
	if productName != "" {
		paramsLike["product_name"] = productName
	}
	productCode := c.QueryParam("product_code")
	if productCode != "" {
		paramsLike["product_code"] = productCode
	}
	isinCode := c.QueryParam("isin_code")
	if isinCode != "" {
		paramsLike["isin_code"] = isinCode
	}

	var msProduct []models.MsProduct

	status, err = models.AdminGetAllMsProductWithLike(&msProduct, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msProduct) < 1 {
		log.Error("product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}

	var currencyIds []string
	var productCategoryIds []string
	var productTypeIds []string
	var genLookupIds []string
	var custodianIds []string
	for _, pro := range msProduct {

		if pro.CurrencyKey != nil {
			if _, ok := lib.Find(currencyIds, strconv.FormatUint(*pro.CurrencyKey, 10)); !ok {
				currencyIds = append(currencyIds, strconv.FormatUint(*pro.CurrencyKey, 10))
			}
		}

		if pro.ProductCategoryKey != nil {
			if _, ok := lib.Find(productCategoryIds, strconv.FormatUint(*pro.ProductCategoryKey, 10)); !ok {
				productCategoryIds = append(productCategoryIds, strconv.FormatUint(*pro.ProductCategoryKey, 10))
			}
		}

		if pro.ProductTypeKey != nil {
			if _, ok := lib.Find(productTypeIds, strconv.FormatUint(*pro.ProductTypeKey, 10)); !ok {
				productTypeIds = append(productTypeIds, strconv.FormatUint(*pro.ProductTypeKey, 10))
			}
		}

		if pro.CustodianKey != nil {
			if _, ok := lib.Find(custodianIds, strconv.FormatUint(*pro.CustodianKey, 10)); !ok {
				custodianIds = append(custodianIds, strconv.FormatUint(*pro.CustodianKey, 10))
			}
		}
	}

	//mapping currency
	var msCurrency []models.MsCurrency
	if len(currencyIds) > 0 {
		status, err = models.GetMsCurrencyIn(&msCurrency, currencyIds, "currency_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	currencyData := make(map[uint64]models.MsCurrency)
	for _, b := range msCurrency {
		currencyData[b.CurrencyKey] = b
	}

	//mapping MsProductCategory
	var msProductCategory []models.MsProductCategory
	if len(productCategoryIds) > 0 {
		status, err = models.GetMsProductCategoryIn(&msProductCategory, productCategoryIds, "product_category_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	proCatData := make(map[uint64]models.MsProductCategory)
	for _, a := range msProductCategory {
		proCatData[a.ProductCategoryKey] = a
	}

	//mapping product_type
	var msProductType []models.MsProductType
	if len(productTypeIds) > 0 {
		status, err = models.GetMsProductTypeIn(&msProductType, productTypeIds, "product_type_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	productTypeData := make(map[uint64]models.MsProductType)
	for _, p := range msProductType {
		productTypeData[p.ProductTypeKey] = p
	}

	//gen lookup
	var lookup []models.GenLookup
	if len(genLookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookup, genLookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookup {
		gData[gen.LookupKey] = gen
	}

	//mapping parent custodian
	var msCustoBank []models.MsCustodianBank
	if len(custodianIds) > 0 {
		status, err = models.GetMsCustodianBankIn(&msCustoBank, custodianIds, "custodian_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	custoData := make(map[uint64]models.MsCustodianBank)
	for _, cus := range msCustoBank {
		custoData[cus.CustodianKey] = cus
	}

	var responseData []models.AdminMsProductList
	for _, pro := range msProduct {
		var data models.AdminMsProductList

		data.ProductKey = pro.ProductKey
		data.ProductID = pro.ProductID
		data.ProductCode = pro.ProductCode
		data.ProductName = pro.ProductName
		data.ProductNameAlt = pro.ProductNameAlt
		if pro.CurrencyKey != nil {
			if n, ok := currencyData[*pro.CurrencyKey]; ok {
				data.CurrencyName = n.Name
			}
		}
		if pro.ProductCategoryKey != nil {
			if n, ok := proCatData[*pro.ProductCategoryKey]; ok {
				data.ProductCategoryName = n.CategoryName
			}
		}
		if pro.ProductTypeKey != nil {
			if n, ok := productTypeData[*pro.ProductTypeKey]; ok {
				data.ProductTypeName = n.ProductTypeName
			}
		}
		if pro.RiskProfileKey != nil {
			if n, ok := gData[*pro.RiskProfileKey]; ok {
				data.RiskProfileName = n.LkpName
			}
		}
		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		if pro.LaunchDate != nil {
			date, err := time.Parse(layout, *pro.LaunchDate)
			if err == nil {
				oke := date.Format(newLayout)
				data.LaunchDate = &oke
			}
		}
		if pro.InceptionDate != nil {
			date, _ := time.Parse(layout, *pro.InceptionDate)
			if err == nil {
				oke := date.Format(newLayout)
				data.InceptionDate = &oke
			}
		}

		data.IsinCode = pro.IsinCode

		if pro.FlagSyariah == 1 {
			data.Syariah = "Ya"
		} else {
			data.Syariah = "Tidak"
		}

		if pro.CustodianKey != nil {
			if n, ok := custoData[*pro.CustodianKey]; ok {
				data.CustodianFullName = n.CustodianFullName
			}
		}

		data.SinvestFundCode = pro.SinvestFundCode

		if pro.FlagEnabled == 1 {
			data.Enabled = "Ya"
		} else {
			data.Enabled = "Tidak"
		}

		if pro.FlagSubscription == 1 {
			data.Subscription = "Ya"
		} else {
			data.Subscription = "Tidak"
		}

		if pro.FlagRedemption == 1 {
			data.Redemption = "Ya"
		} else {
			data.Redemption = "Tidak"
		}

		if pro.FlagSwitchOut == 1 {
			data.SwitchOut = "Ya"
		} else {
			data.SwitchOut = "Tidak"
		}

		if pro.FlagSwitchIn == 1 {
			data.SwitchIn = "Ya"
		} else {
			data.SwitchIn = "Tidak"
		}

		responseData = append(responseData, data)
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminGetCountMsProductWithLike(&countData, params, paramsLike)
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
