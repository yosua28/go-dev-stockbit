package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
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

	productkey := c.QueryParam("product_key")
	if productkey != "" {
		productkeyCek, err := strconv.ParseUint(productkey, 10, 64)
		if err == nil && productkeyCek > 0 {
			params["pf.product_key"] = productkey
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
		}
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
	decimal.MarshalJSONWithoutQuotes = true

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

func DeleteProductFeeAdmin(c echo.Context) error {
	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	feekey := c.FormValue("fee_key")
	if feekey == "" {
		log.Error("Missing required parameter: fee_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_key", "Missing required parameter: fee_key")
	}

	feekeyCek, err := strconv.ParseUint(feekey, 10, 64)
	if err == nil && feekeyCek > 0 {
		params["fee_key"] = feekey
	} else {
		log.Error("Wrong input for parameter: fee_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_key", "Missing required parameter: fee_key")
	}

	var productFee models.MsProductFee
	status, err := models.GetMsProductFee(&productFee, feekey)
	if err != nil {
		log.Error("Product Fee not found hahahah")
		return lib.CustomError(status)
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateMsProductFee(params)
	if err != nil {
		log.Error("Failed update request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed update data")
	}

	paramsFeeItems := make(map[string]string)
	paramsFeeItems["rec_status"] = "0"
	paramsFeeItems["rec_deleted_date"] = time.Now().Format(dateLayout)
	paramsFeeItems["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateMsProductFeeItemByField(paramsFeeItems, feekey, "product_fee_key")
	if err != nil {
		log.Error("Failed delete request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func CreateAdminMsProductFee(c echo.Context) error {
	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	//product_key
	productkey := c.FormValue("product_key")
	if productkey == "" {
		log.Error("Missing required parameter: product_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key cann't be blank", "Missing required parameter: product_key cann't be blank")
	}
	strproductkey, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && strproductkey > 0 {
		params["product_key"] = productkey
	} else {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	//fee_type
	feetype := c.FormValue("fee_type")
	if feetype != "" {
		strfeetype, err := strconv.ParseUint(feetype, 10, 64)
		if err == nil && strfeetype > 0 {
			params["fee_type"] = feetype
		} else {
			log.Error("Wrong input for parameter: fee_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_type", "Missing required parameter: fee_type")
		}
	}

	//fee_code
	feecode := c.FormValue("fee_code")
	if feecode != "" {
		params["fee_code"] = feecode
	}

	//flag_show_ontnc
	flagshowontnc := c.FormValue("flag_show_ontnc")
	var flagshowontncBool bool
	if flagshowontnc != "" {
		flagshowontncBool, err = strconv.ParseBool(flagshowontnc)
		if err != nil {
			log.Error("flag_show_ontnc parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_show_ontnc parameter should be true/false", "flag_show_ontnc parameter should be true/false")
		}
		if flagshowontncBool == true {
			params["flag_show_ontnc"] = "1"
		} else {
			params["flag_show_ontnc"] = "0"
		}
	} else {
		params["flag_show_ontnc"] = "0"
	}

	//fee_annotation
	feeannotation := c.FormValue("fee_annotation")
	if feeannotation != "" {
		params["fee_annotation"] = feeannotation
	}

	//fee_desc
	feedesc := c.FormValue("fee_desc")
	if feedesc != "" {
		params["fee_desc"] = feedesc
	}

	//fee_date_start
	feedatestart := c.FormValue("fee_date_start")
	if feedatestart != "" {
		params["fee_date_start"] = feedatestart + " 00:00:00"
	}

	//fee_date_thru
	feedatethru := c.FormValue("fee_date_thru")
	if feedatethru != "" {
		params["fee_date_thru"] = feedatethru + " 00:00:00"
	}

	//fee_nominal_type
	feenominaltype := c.FormValue("fee_nominal_type")
	if feenominaltype != "" {
		strfeenominaltype, err := strconv.ParseUint(feenominaltype, 10, 64)
		if err == nil && strfeenominaltype > 0 {
			params["fee_nominal_type"] = feenominaltype
		} else {
			log.Error("Wrong input for parameter: fee_nominal_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_nominal_type", "Missing required parameter: fee_nominal_type")
		}
	}

	//enabled_min_amount
	enabledminamount := c.FormValue("enabled_min_amount")
	var enabledminamountBool bool
	if enabledminamount != "" {
		enabledminamountBool, err = strconv.ParseBool(enabledminamount)
		if err != nil {
			log.Error("enabled_min_amount parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "enabled_min_amount parameter should be true/false", "enabled_min_amount parameter should be true/false")
		}
		if enabledminamountBool == true {
			params["enabled_min_amount"] = "1"
		} else {
			params["enabled_min_amount"] = "0"
		}
	} else {
		log.Error("enabled_min_amount parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "enabled_min_amount parameter should be true/false", "enabled_min_amount parameter should be true/false")
	}

	//fee_min_amount
	feeminamount := c.FormValue("fee_min_amount")
	if feeminamount == "" {
		if enabledminamountBool == true {
			log.Error("Missing required parameter: fee_min_amount cann't be blank")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount cann't be blank", "Missing required parameter: fee_min_amount cann't be blank")
		}
	} else {
		feeminamountFloat, err := strconv.ParseFloat(feeminamount, 64)
		if err == nil {
			if feeminamountFloat < 0 {
				log.Error("Wrong input for parameter: fee_min_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount must cann't negatif", "Missing required parameter: fee_min_amount cann't negatif")
			}
			params["fee_min_amount"] = feeminamount
		} else {
			log.Error("Wrong input for parameter: fee_min_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount must number", "Missing required parameter: fee_min_amount number")
		}
	}

	//enabled_max_amount
	enabledmaxamount := c.FormValue("enabled_max_amount")
	var enabledmaxamountBool bool
	if enabledmaxamount != "" {
		enabledmaxamountBool, err = strconv.ParseBool(enabledmaxamount)
		if err != nil {
			log.Error("enabled_max_amount parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "enabled_max_amount parameter should be true/false", "enabled_max_amount parameter should be true/false")
		}
		if enabledmaxamountBool == true {
			params["enabled_max_amount"] = "1"
		} else {
			params["enabled_max_amount"] = "0"
		}
	} else {
		log.Error("enabled_max_amount parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "enabled_max_amount parameter should be true/false", "enabled_max_amount parameter should be true/false")
	}

	//fee_max_amount
	feemaxamount := c.FormValue("fee_max_amount")
	if feemaxamount == "" {
		if enabledmaxamountBool == true {
			log.Error("Missing required parameter: fee_max_amount cann't be blank")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount cann't be blank", "Missing required parameter: fee_max_amount cann't be blank")
		}
	} else {
		feemaxamountFloat, err := strconv.ParseFloat(feemaxamount, 64)
		if err == nil {
			if feemaxamountFloat < 0 {
				log.Error("Wrong input for parameter: fee_max_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount must cann't negatif", "Missing required parameter: fee_max_amount cann't negatif")
			}
			params["fee_max_amount"] = feemaxamount
		} else {
			log.Error("Wrong input for parameter: fee_max_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount must number", "Missing required parameter: fee_max_amount number")
		}
	}

	//fee_calc_method
	feecalcmethod := c.FormValue("fee_calc_method")
	if feecalcmethod != "" {
		strfeecalcmethod, err := strconv.ParseUint(feecalcmethod, 10, 64)
		if err == nil && strfeecalcmethod > 0 {
			params["fee_calc_method"] = feecalcmethod
		} else {
			log.Error("Wrong input for parameter: fee_calc_method")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_calc_method", "Missing required parameter: fee_calc_method")
		}
	}

	//calculation_baseon
	calculationbaseon := c.FormValue("calculation_baseon")
	if calculationbaseon != "" {
		strcalculationbaseon, err := strconv.ParseUint(calculationbaseon, 10, 64)
		if err == nil && strcalculationbaseon > 0 {
			params["calculation_baseon"] = calculationbaseon
		} else {
			log.Error("Wrong input for parameter: calculation_baseon")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: calculation_baseon", "Missing required parameter: calculation_baseon")
		}
	}

	//period_hold
	periodhold := c.FormValue("period_hold")
	if periodhold == "" {
		log.Error("Missing required parameter: period_hold cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: period_hold cann't be blank", "Missing required parameter: period_hold cann't be blank")
	}
	strperiodhold, err := strconv.ParseUint(periodhold, 10, 64)
	if err == nil && strperiodhold > 0 {
		params["period_hold"] = periodhold
	} else {
		log.Error("Wrong input for parameter: period_hold")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: period_hold", "Missing required parameter: period_hold")
	}

	//days_inyear
	daysinyear := c.FormValue("days_inyear")
	if daysinyear != "" {
		strdaysinyear, err := strconv.ParseUint(daysinyear, 10, 64)
		if err == nil && strdaysinyear > 0 {
			params["days_inyear"] = daysinyear
		} else {
			log.Error("Wrong input for parameter: days_inyear")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: days_inyear", "Missing required parameter: days_inyear")
		}
	}

	//json_fee_items
	var feeItems []map[string]interface{}
	jsonfeeitems := c.FormValue("json_fee_items")
	if jsonfeeitems != "" {

		jsonfeeitems = strings.Replace(jsonfeeitems, "\\", "", -1)
		if err := json.Unmarshal([]byte(jsonfeeitems), &feeItems); err != nil {
			log.Error(err)
			log.Error("Wrong input for parameter: json_fee_items")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: json_fee_items is invalid", "Missing required parameter: json_fee_items is invalid")
		}

		log.Println("===========================")
		fmt.Println(jsonfeeitems)
		fmt.Println(feeItems)
		log.Println("===========================")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "1"
	params["rec_order"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err, lastID := models.CreateMsProductFee(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	if len(feeItems) > 0 {
		for idx, val := range feeItems {

			idxStr := strconv.FormatInt(int64(idx), 10)
			paramsFeeItems := make(map[string]string)
			paramsFeeItems["product_fee_key"] = lastID
			paramsFeeItems["item_seqno"] = idxStr
			if (len(feeItems) - 1) == idx {
				paramsFeeItems["row_max"] = "1"
			} else {
				paramsFeeItems["row_max"] = "0"
			}

			strPrincipleLimit := fmt.Sprintf("%g", val["principle_limit"])
			paramsFeeItems["principle_limit"] = strPrincipleLimit

			strFeeValue := fmt.Sprintf("%g", val["fee_value"])
			paramsFeeItems["fee_value"] = strFeeValue

			paramsFeeItems["item_notes"] = val["item_notes"].(string)
			paramsFeeItems["rec_status"] = "1"
			paramsFeeItems["rec_order"] = "1"
			paramsFeeItems["rec_created_date"] = time.Now().Format(dateLayout)
			paramsFeeItems["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

			status, err := models.CreateMsProductFeeItem(paramsFeeItems)
			if err != nil {
				log.Error("Failed create request data: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed input data")
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminMsProductFee(c echo.Context) error {
	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	//fee_key
	feekey := c.FormValue("fee_key")
	if feekey == "" {
		log.Error("Missing required parameter: fee_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_key cann't be blank", "Missing required parameter: fee_key cann't be blank")
	}
	strfeekey, err := strconv.ParseUint(feekey, 10, 64)
	if err == nil && strfeekey > 0 {
		params["fee_key"] = feekey
	} else {
		log.Error("Wrong input for parameter: fee_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_key", "Missing required parameter: fee_key")
	}

	var productFee models.MsProductFee
	status, err = models.GetMsProductFee(&productFee, feekey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	//product_key
	productkey := c.FormValue("product_key")
	if productkey == "" {
		log.Error("Missing required parameter: product_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key cann't be blank", "Missing required parameter: product_key cann't be blank")
	}
	strproductkey, err := strconv.ParseUint(productkey, 10, 64)
	if err == nil && strproductkey > 0 {
		params["product_key"] = productkey
	} else {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	//fee_type
	feetype := c.FormValue("fee_type")
	if feetype != "" {
		strfeetype, err := strconv.ParseUint(feetype, 10, 64)
		if err == nil && strfeetype > 0 {
			params["fee_type"] = feetype
		} else {
			log.Error("Wrong input for parameter: fee_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_type", "Missing required parameter: fee_type")
		}
	}

	//fee_code
	feecode := c.FormValue("fee_code")
	if feecode != "" {
		params["fee_code"] = feecode
	}

	//flag_show_ontnc
	flagshowontnc := c.FormValue("flag_show_ontnc")
	var flagshowontncBool bool
	if flagshowontnc != "" {
		flagshowontncBool, err = strconv.ParseBool(flagshowontnc)
		if err != nil {
			log.Error("flag_show_ontnc parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "flag_show_ontnc parameter should be true/false", "flag_show_ontnc parameter should be true/false")
		}
		if flagshowontncBool == true {
			params["flag_show_ontnc"] = "1"
		} else {
			params["flag_show_ontnc"] = "0"
		}
	} else {
		params["flag_show_ontnc"] = "0"
	}

	//fee_annotation
	feeannotation := c.FormValue("fee_annotation")
	if feeannotation != "" {
		params["fee_annotation"] = feeannotation
	}

	//fee_desc
	feedesc := c.FormValue("fee_desc")
	if feedesc != "" {
		params["fee_desc"] = feedesc
	}

	//fee_date_start
	feedatestart := c.FormValue("fee_date_start")
	if feedatestart != "" {
		params["fee_date_start"] = feedatestart + " 00:00:00"
	}

	//fee_date_thru
	feedatethru := c.FormValue("fee_date_thru")
	if feedatethru != "" {
		params["fee_date_thru"] = feedatethru + " 00:00:00"
	}

	//fee_nominal_type
	feenominaltype := c.FormValue("fee_nominal_type")
	if feenominaltype != "" {
		strfeenominaltype, err := strconv.ParseUint(feenominaltype, 10, 64)
		if err == nil && strfeenominaltype > 0 {
			params["fee_nominal_type"] = feenominaltype
		} else {
			log.Error("Wrong input for parameter: fee_nominal_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_nominal_type", "Missing required parameter: fee_nominal_type")
		}
	}

	//enabled_min_amount
	enabledminamount := c.FormValue("enabled_min_amount")
	var enabledminamountBool bool
	if enabledminamount != "" {
		enabledminamountBool, err = strconv.ParseBool(enabledminamount)
		if err != nil {
			log.Error("enabled_min_amount parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "enabled_min_amount parameter should be true/false", "enabled_min_amount parameter should be true/false")
		}
		if enabledminamountBool == true {
			params["enabled_min_amount"] = "1"
		} else {
			params["enabled_min_amount"] = "0"
		}
	} else {
		log.Error("enabled_min_amount parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "enabled_min_amount parameter should be true/false", "enabled_min_amount parameter should be true/false")
	}

	//fee_min_amount
	feeminamount := c.FormValue("fee_min_amount")
	if feeminamount == "" {
		if enabledminamountBool == true {
			log.Error("Missing required parameter: fee_min_amount cann't be blank")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount cann't be blank", "Missing required parameter: fee_min_amount cann't be blank")
		}
	} else {
		feeminamountFloat, err := strconv.ParseFloat(feeminamount, 64)
		if err == nil {
			if feeminamountFloat < 0 {
				log.Error("Wrong input for parameter: fee_min_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount must cann't negatif", "Missing required parameter: fee_min_amount cann't negatif")
			}
			params["fee_min_amount"] = feeminamount
		} else {
			log.Error("Wrong input for parameter: fee_min_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_min_amount must number", "Missing required parameter: fee_min_amount number")
		}
	}

	//enabled_max_amount
	enabledmaxamount := c.FormValue("enabled_max_amount")
	var enabledmaxamountBool bool
	if enabledmaxamount != "" {
		enabledmaxamountBool, err = strconv.ParseBool(enabledmaxamount)
		if err != nil {
			log.Error("enabled_max_amount parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "enabled_max_amount parameter should be true/false", "enabled_max_amount parameter should be true/false")
		}
		if enabledmaxamountBool == true {
			params["enabled_max_amount"] = "1"
		} else {
			params["enabled_max_amount"] = "0"
		}
	} else {
		log.Error("enabled_max_amount parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "enabled_max_amount parameter should be true/false", "enabled_max_amount parameter should be true/false")
	}

	//fee_max_amount
	feemaxamount := c.FormValue("fee_max_amount")
	if feemaxamount == "" {
		if enabledmaxamountBool == true {
			log.Error("Missing required parameter: fee_max_amount cann't be blank")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount cann't be blank", "Missing required parameter: fee_max_amount cann't be blank")
		}
	} else {
		feemaxamountFloat, err := strconv.ParseFloat(feemaxamount, 64)
		if err == nil {
			if feemaxamountFloat < 0 {
				log.Error("Wrong input for parameter: fee_max_amount cann't negatif")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount must cann't negatif", "Missing required parameter: fee_max_amount cann't negatif")
			}
			params["fee_max_amount"] = feemaxamount
		} else {
			log.Error("Wrong input for parameter: fee_max_amount number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_max_amount must number", "Missing required parameter: fee_max_amount number")
		}
	}

	//fee_calc_method
	feecalcmethod := c.FormValue("fee_calc_method")
	if feecalcmethod != "" {
		strfeecalcmethod, err := strconv.ParseUint(feecalcmethod, 10, 64)
		if err == nil && strfeecalcmethod > 0 {
			params["fee_calc_method"] = feecalcmethod
		} else {
			log.Error("Wrong input for parameter: fee_calc_method")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_calc_method", "Missing required parameter: fee_calc_method")
		}
	}

	//calculation_baseon
	calculationbaseon := c.FormValue("calculation_baseon")
	if calculationbaseon != "" {
		strcalculationbaseon, err := strconv.ParseUint(calculationbaseon, 10, 64)
		if err == nil && strcalculationbaseon > 0 {
			params["calculation_baseon"] = calculationbaseon
		} else {
			log.Error("Wrong input for parameter: calculation_baseon")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: calculation_baseon", "Missing required parameter: calculation_baseon")
		}
	}

	//period_hold
	periodhold := c.FormValue("period_hold")
	if periodhold == "" {
		log.Error("Missing required parameter: period_hold cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: period_hold cann't be blank", "Missing required parameter: period_hold cann't be blank")
	}
	strperiodhold, err := strconv.ParseUint(periodhold, 10, 64)
	if err == nil && strperiodhold > 0 {
		params["period_hold"] = periodhold
	} else {
		log.Error("Wrong input for parameter: period_hold")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: period_hold", "Missing required parameter: period_hold")
	}

	//days_inyear
	daysinyear := c.FormValue("days_inyear")
	if daysinyear != "" {
		strdaysinyear, err := strconv.ParseUint(daysinyear, 10, 64)
		if err == nil && strdaysinyear > 0 {
			params["days_inyear"] = daysinyear
		} else {
			log.Error("Wrong input for parameter: days_inyear")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: days_inyear", "Missing required parameter: days_inyear")
		}
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateMsProductFee(params)
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

func CreateAdminMsProductFeeItem(c echo.Context) error {
	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	//product_fee_key
	productfeekey := c.FormValue("product_fee_key")
	if productfeekey == "" {
		log.Error("Missing required parameter: product_fee_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_key cann't be blank", "Missing required parameter: product_fee_key cann't be blank")
	}
	strproductfeekey, err := strconv.ParseUint(productfeekey, 10, 64)
	if err == nil && strproductfeekey > 0 {
		params["product_fee_key"] = productfeekey
	} else {
		log.Error("Wrong input for parameter: product_fee_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_key", "Missing required parameter: product_fee_key")
	}

	var productFee models.MsProductFee
	status, err = models.GetMsProductFee(&productFee, productfeekey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	//principle_limit
	principlelimit := c.FormValue("principle_limit")
	if principlelimit == "" {
		log.Error("Missing required parameter: principle_limit cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: principle_limit cann't be blank", "Missing required parameter: principle_limit cann't be blank")
	}
	_, err = strconv.ParseFloat(principlelimit, 64)
	if err == nil {
		params["principle_limit"] = principlelimit
	} else {
		log.Error("Wrong input for parameter: principle_limit")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: principle_limit", "Missing required parameter: principle_limit")
	}

	//fee_value
	feevalue := c.FormValue("fee_value")
	if feevalue == "" {
		log.Error("Missing required parameter: fee_value cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_value cann't be blank", "Missing required parameter: fee_value cann't be blank")
	}
	_, err = strconv.ParseFloat(feevalue, 64)
	if err == nil {
		params["fee_value"] = feevalue
	} else {
		log.Error("Wrong input for parameter: fee_value")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_value", "Missing required parameter: fee_value")
	}

	//item_notes
	itemnotes := c.FormValue("item_notes")
	if itemnotes != "" {
		params["item_notes"] = itemnotes
	}

	//get lastitem_seqno
	var productFeeItem models.MsProductFeeItem
	status, err = models.GetLastMsProductFeeItemByFeeKey(&productFeeItem, productfeekey, "item_seqno", "DESC")
	if err != nil {
		params["item_seqno"] = "0"
	} else {
		seqNoStr := strconv.FormatUint((productFeeItem.ItemSeqno + 1), 10)
		params["item_seqno"] = seqNoStr
	}

	params["row_max"] = "0"

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "1"
	params["rec_order"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.CreateMsProductFeeItem(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//update row_max_all
	var productFeeItems []models.MsProductFeeItem
	paramsUpdate := make(map[string]string)
	paramsUpdate["product_fee_key"] = productfeekey
	paramsUpdate["rec_status"] = "1"
	paramsUpdate["orderBy"] = "principle_limit"
	paramsUpdate["orderType"] = "ASC"
	status, err = models.GetAllMsProductFeeItem(&productFeeItems, paramsUpdate)
	if err == nil {
		for idx, item := range productFeeItems {
			paramsFeeItems := make(map[string]string)
			if len(productFeeItems)-1 == idx {
				paramsFeeItems["row_max"] = "1"
			} else {
				paramsFeeItems["row_max"] = "0"
			}
			paramsFeeItems["rec_modified_date"] = time.Now().Format(dateLayout)
			paramsFeeItems["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

			keyFeeItem := strconv.FormatUint(item.ProductFeeItemKey, 10)
			status, err = models.UpdateMsProductFeeItemByField(paramsFeeItems, keyFeeItem, "product_fee_item_key")
			if err != nil {
				log.Error("Failed delete request data: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed delete data")
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminMsProductFeeItem(c echo.Context) error {
	var err error
	var status int

	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	//product_fee_key
	productfeekey := c.FormValue("product_fee_key")
	if productfeekey == "" {
		log.Error("Missing required parameter: product_fee_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_key cann't be blank", "Missing required parameter: product_fee_key cann't be blank")
	}
	strproductfeekey, err := strconv.ParseUint(productfeekey, 10, 64)
	if err == nil && strproductfeekey > 0 {
	} else {
		log.Error("Wrong input for parameter: product_fee_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_key", "Missing required parameter: product_fee_key")
	}

	//product_fee_item_key
	productfeeitemkey := c.FormValue("product_fee_item_key")
	if productfeeitemkey == "" {
		log.Error("Missing required parameter: product_fee_item_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_item_key cann't be blank", "Missing required parameter: product_fee_item_key cann't be blank")
	}
	strproductfeeitemkey, err := strconv.ParseUint(productfeeitemkey, 10, 64)
	if err == nil && strproductfeeitemkey > 0 {
		params["product_fee_item_key"] = productfeeitemkey
	} else {
		log.Error("Wrong input for parameter: product_fee_item_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_item_key", "Missing required parameter: product_fee_item_key")
	}

	//cek data
	var productFeeItems []models.MsProductFeeItem
	paramsCheck := make(map[string]string)
	paramsCheck["product_fee_key"] = productfeekey
	paramsCheck["product_fee_item_key"] = productfeeitemkey
	paramsCheck["rec_status"] = "1"
	status, err = models.GetAllMsProductFeeItem(&productFeeItems, paramsCheck)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	//principle_limit
	principlelimit := c.FormValue("principle_limit")
	if principlelimit == "" {
		log.Error("Missing required parameter: principle_limit cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: principle_limit cann't be blank", "Missing required parameter: principle_limit cann't be blank")
	}
	_, err = strconv.ParseFloat(principlelimit, 64)
	if err == nil {
		params["principle_limit"] = principlelimit
	} else {
		log.Error("Wrong input for parameter: principle_limit")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: principle_limit", "Missing required parameter: principle_limit")
	}

	//fee_value
	feevalue := c.FormValue("fee_value")
	if feevalue == "" {
		log.Error("Missing required parameter: fee_value cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_value cann't be blank", "Missing required parameter: fee_value cann't be blank")
	}
	_, err = strconv.ParseFloat(feevalue, 64)
	if err == nil {
		params["fee_value"] = feevalue
	} else {
		log.Error("Wrong input for parameter: fee_value")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fee_value", "Missing required parameter: fee_value")
	}

	//item_notes
	itemnotes := c.FormValue("item_notes")
	if itemnotes != "" {
		params["item_notes"] = itemnotes
	}

	params["row_max"] = "0"

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "1"
	params["rec_order"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateMsProductFeeItemByField(params, productfeeitemkey, "product_fee_item_key")
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//update row_max_all
	var productFeeItemsUpdate []models.MsProductFeeItem
	paramsUpdate := make(map[string]string)
	paramsUpdate["product_fee_key"] = productfeekey
	paramsUpdate["rec_status"] = "1"
	paramsUpdate["orderBy"] = "principle_limit"
	paramsUpdate["orderType"] = "ASC"
	status, err = models.GetAllMsProductFeeItem(&productFeeItemsUpdate, paramsUpdate)
	if err == nil {
		for idx, item := range productFeeItemsUpdate {
			paramsFeeItems := make(map[string]string)
			if len(productFeeItemsUpdate)-1 == idx {
				paramsFeeItems["row_max"] = "1"
			} else {
				paramsFeeItems["row_max"] = "0"
			}
			paramsFeeItems["rec_modified_date"] = time.Now().Format(dateLayout)
			paramsFeeItems["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

			keyFeeItem := strconv.FormatUint(item.ProductFeeItemKey, 10)
			status, err = models.UpdateMsProductFeeItemByField(paramsFeeItems, keyFeeItem, "product_fee_item_key")
			if err != nil {
				log.Error("Failed delete request data: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed delete data")
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func DeleteAdminMsProductFeeItem(c echo.Context) error {
	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	//product_fee_key
	productfeekey := c.FormValue("product_fee_key")
	if productfeekey == "" {
		log.Error("Missing required parameter: product_fee_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_key cann't be blank", "Missing required parameter: product_fee_key cann't be blank")
	}
	strproductfeekey, err := strconv.ParseUint(productfeekey, 10, 64)
	if err == nil && strproductfeekey > 0 {
	} else {
		log.Error("Wrong input for parameter: product_fee_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_key", "Missing required parameter: product_fee_key")
	}

	//product_fee_item_key
	productfeeitemkey := c.FormValue("product_fee_item_key")
	if productfeeitemkey == "" {
		log.Error("Missing required parameter: product_fee_item_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_item_key cann't be blank", "Missing required parameter: product_fee_item_key cann't be blank")
	}
	strproductfeeitemkey, err := strconv.ParseUint(productfeeitemkey, 10, 64)
	if err == nil && strproductfeeitemkey > 0 {
		params["product_fee_item_key"] = productfeeitemkey
	} else {
		log.Error("Wrong input for parameter: product_fee_item_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_fee_item_key", "Missing required parameter: product_fee_item_key")
	}

	//cek data
	var productFeeItems []models.MsProductFeeItem
	paramsCheck := make(map[string]string)
	paramsCheck["product_fee_key"] = productfeekey
	paramsCheck["product_fee_item_key"] = productfeeitemkey
	paramsCheck["rec_status"] = "1"
	status, err = models.GetAllMsProductFeeItem(&productFeeItems, paramsCheck)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["row_max"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateMsProductFeeItemByField(params, productfeeitemkey, "product_fee_item_key")
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//update row_max_all
	var productFeeItemsUpdate []models.MsProductFeeItem
	paramsUpdate := make(map[string]string)
	paramsUpdate["product_fee_key"] = productfeekey
	paramsUpdate["rec_status"] = "1"
	paramsUpdate["orderBy"] = "principle_limit"
	paramsUpdate["orderType"] = "ASC"
	status, err = models.GetAllMsProductFeeItem(&productFeeItemsUpdate, paramsUpdate)
	if err == nil {
		for idx, item := range productFeeItemsUpdate {
			paramsFeeItems := make(map[string]string)
			if len(productFeeItemsUpdate)-1 == idx {
				paramsFeeItems["row_max"] = "1"
			} else {
				paramsFeeItems["row_max"] = "0"
			}
			paramsFeeItems["rec_modified_date"] = time.Now().Format(dateLayout)
			paramsFeeItems["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

			keyFeeItem := strconv.FormatUint(item.ProductFeeItemKey, 10)
			status, err = models.UpdateMsProductFeeItemByField(paramsFeeItems, keyFeeItem, "product_fee_item_key")
			if err != nil {
				log.Error("Failed delete request data: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed delete data")
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func DetailAdminMsProductFeeItem(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	//product_fee_item_key
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	//cek data
	var productFeeItem models.MsProductFeeItem
	_, err = models.GetMsProductFeeItem(&productFeeItem, keyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	var responseData models.MsProductFeeItemDetailList

	responseData.ProductFeeItemKey = productFeeItem.ProductFeeItemKey
	responseData.PrincipleLimit = productFeeItem.PrincipleLimit
	responseData.FeeValue = productFeeItem.FeeValue
	responseData.ItemNotes = productFeeItem.ItemNotes

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)

}
