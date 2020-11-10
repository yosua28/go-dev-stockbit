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

func GetListProductFeeAdmin(c echo.Context) error {

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

	items := []string{"fee_key", "fee_code", "product_key", "product_code", "product_name", "feetypename", "fee_date_start", "fee_date_thru", "period_hold"}

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

	var searchData *string

	search := c.QueryParam("search")
	if search != "" {
		searchData = &search
	}

	//mapping parent custodian
	var msProductFee []models.AdminListMsProductFee
	status, err = models.AdminGetAllMsProductFee(&msProductFee, limit, offset, params, noLimit, searchData)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminCountDataGetAllMsProductFee(&countData, params, searchData)
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
	response.Data = msProductFee

	return c.JSON(http.StatusOK, response)
}

func GetProductFeeDetailAdmin(c echo.Context) error {

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

	var productFee models.MsProductFee
	status, err = models.GetMsProductFee(&productFee, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.MsProductFeeDetailAdmin

	var lookupIds []string

	if productFee.FeeType != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*productFee.FeeType, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*productFee.FeeType, 10))
		}
	}

	if productFee.FeeNominalType != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*productFee.FeeNominalType, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*productFee.FeeNominalType, 10))
		}
	}

	if productFee.FeeCalcMethod != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*productFee.FeeCalcMethod, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*productFee.FeeCalcMethod, 10))
		}
	}

	if productFee.CalculationBaseon != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*productFee.CalculationBaseon, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*productFee.CalculationBaseon, 10))
		}
	}

	if productFee.DaysInyear != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*productFee.DaysInyear, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*productFee.DaysInyear, 10))
		}
	}

	//gen lookup oa request
	var lookupProductFee []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupProductFee, lookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupProductFee {
		gData[gen.LookupKey] = gen
	}

	responseData.FeeKey = productFee.FeeKey

	//product
	var product models.MsProduct
	strProductKey := strconv.FormatUint(productFee.ProductKey, 10)
	status, err = models.GetMsProduct(&product, strProductKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		var pro models.MsProductListDropdown
		pro.ProductKey = product.ProductKey
		pro.ProductCode = product.ProductCode
		pro.ProductName = product.ProductName
		responseData.Product = pro
	}

	if productFee.FeeType != nil {
		if n, ok := gData[*productFee.FeeType]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.FeeType = &trc
		}
	}

	responseData.FeeCode = productFee.FeeCode

	if productFee.FlagShowOntnc != nil {
		if *productFee.FlagShowOntnc == uint8(1) {
			responseData.FlagShowOntnc = true
		} else {
			responseData.FlagShowOntnc = false
		}
	} else {
		responseData.FlagShowOntnc = false
	}

	responseData.FeeAnnotation = productFee.FeeAnnotation
	responseData.FeeDesc = productFee.FeeDesc
	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	if productFee.FeeDateStart != nil {
		date, err := time.Parse(layout, *productFee.FeeDateStart)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.FeeDateStart = &oke
		}
	}
	if productFee.FeeDateThru != nil {
		date, err := time.Parse(layout, *productFee.FeeDateThru)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.FeeDateThru = &oke
		}
	}

	if productFee.FeeNominalType != nil {
		if n, ok := gData[*productFee.FeeNominalType]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.FeeNominalType = &trc
		}
	}

	if productFee.EnabledMinAmount == uint8(1) {
		responseData.EnabledMinAmount = true
	} else {
		responseData.EnabledMinAmount = false
	}

	responseData.FeeMinAmount = productFee.FeeMinAmount

	if productFee.EnabledMaxAmount == uint8(1) {
		responseData.EnabledMaxAmount = true
	} else {
		responseData.EnabledMaxAmount = false
	}

	responseData.FeeMaxAmount = productFee.FeeMaxAmount

	if productFee.FeeCalcMethod != nil {
		if n, ok := gData[*productFee.FeeCalcMethod]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.FeeCalcMethod = &trc
		}
	}

	if productFee.CalculationBaseon != nil {
		if n, ok := gData[*productFee.CalculationBaseon]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.CalculationBaseon = &trc
		}
	}

	responseData.PeriodHold = productFee.PeriodHold

	if productFee.DaysInyear != nil {
		if n, ok := gData[*productFee.DaysInyear]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.DaysInyear = &trc
		}
	}

	params := make(map[string]string)
	strProductFee := strconv.FormatUint(productFee.FeeKey, 10)
	params["product_fee_key"] = strProductFee
	params["rec_status"] = "1"
	var productFeeItems []models.MsProductFeeItem
	status, err = models.GetAllMsProductFeeItem(&productFeeItems, params)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	var productFeeList []models.MsProductFeeItemDetailList
	for _, feeItem := range productFeeItems {
		var data models.MsProductFeeItemDetailList

		data.ProductFeeItemKey = feeItem.ProductFeeItemKey
		data.PrincipleLimit = feeItem.PrincipleLimit
		data.FeeValue = feeItem.FeeValue
		data.ItemNotes = feeItem.ItemNotes

		productFeeList = append(productFeeList, data)
	}

	responseData.ProductFeeItems = &productFeeList

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
