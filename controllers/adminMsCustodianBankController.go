package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetListMsCustodianBank(c echo.Context) error {
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

	items := []string{"custodian_key", "custodian_code", "custodian_short_name", "custodian_full_name"}

	// Get parameter order_by
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
	}

	var custodianBank []models.MsCustodianBank
	status, err = models.AdminGetListMsCustodianBank(&custodianBank, limit, offset, params, noLimit)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var responseData []models.MsCustodianBankInfoList

	for _, cus := range custodianBank {
		var data models.MsCustodianBankInfoList
		data.CustodianKey = cus.CustodianKey
		data.CustodianCode = cus.CustodianCode
		data.CustodianShortName = cus.CustodianShortName
		data.CustodianFullName = cus.CustodianFullName

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func AdminGetListMsCustodianBank(c echo.Context) error {

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

	items := []string{"custodian_code", "custodian_short_name", "custodian_full_name", "bi_member_code", "swift_code"}

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
		params["orderBy"] = "custodian_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	var bank []models.ListCustodianBankAdmin

	status, err = models.AdminGetListCustodianBank(&bank, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(bank) < 1 {
		log.Error("Custodian Bank not found")
		return lib.CustomError(http.StatusNotFound, "Custodian Bank not found", "Custodian Bank not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListCustodianBank(&countData, params, searchLike)
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
	response.Data = bank

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMsCustodianBank(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("custodian_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: custodian_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: custodian_key", "Missing required parameter: custodian_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["custodian_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsCustodianBank(params)
	if err != nil {
		log.Error("Error delete ms_custodian_bank")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateMsCustodianBank(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	custodianCode := c.FormValue("custodian_code")
	if custodianCode == "" {
		log.Error("Missing required parameter: custodian_code")
		return lib.CustomError(http.StatusBadRequest, "custodian_code can not be blank", "custodian_code can not be blank")
	} else {
		//validate unique custodian_code
		var countData models.CountData
		status, err = models.CountMsCustodianBankValidateUnique(&countData, "custodian_code", custodianCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: custodian_code")
			return lib.CustomError(http.StatusBadRequest, "custodian_code already used", "custodian_code already used")
		}
		params["custodian_code"] = custodianCode
	}

	custodianShortName := c.FormValue("custodian_short_name")
	if custodianShortName == "" {
		log.Error("Missing required parameter: custodian_short_name")
		return lib.CustomError(http.StatusBadRequest, "custodian_short_name can not be blank", "custodian_short_name can not be blank")
	} else {
		//validate unique custodian_short_name
		var countData models.CountData
		status, err = models.CountMsCustodianBankValidateUnique(&countData, "custodian_short_name", custodianShortName, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: custodian_short_name")
			return lib.CustomError(http.StatusBadRequest, "custodian_short_name already used", "custodian_short_name already used")
		}
		params["custodian_short_name"] = custodianShortName
	}

	custodianFullName := c.FormValue("custodian_full_name")
	if custodianFullName == "" {
		log.Error("Missing required parameter: custodian_full_name")
		return lib.CustomError(http.StatusBadRequest, "custodian_full_name can not be blank", "custodian_full_name can not be blank")
	} else {
		params["custodian_full_name"] = custodianFullName
	}

	biMemberCode := c.FormValue("bi_member_code")
	if biMemberCode != "" {
		//validate unique bi_member_code
		var countData models.CountData
		status, err = models.CountMsCustodianBankValidateUnique(&countData, "bi_member_code", biMemberCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: bi_member_code")
			return lib.CustomError(http.StatusBadRequest, "bi_member_code already used", "bi_member_code already used")
		}
		params["bi_member_code"] = biMemberCode
	}

	swiftCode := c.FormValue("swift_code")
	if swiftCode != "" {
		//validate unique swift_code
		var countData models.CountData
		status, err = models.CountMsCustodianBankValidateUnique(&countData, "swift_code", swiftCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: swift_code")
			return lib.CustomError(http.StatusBadRequest, "swift_code already used", "swift_code already used")
		}
		params["swift_code"] = swiftCode
	}

	bankLocal := c.FormValue("bank_local")
	if bankLocal == "" {
		log.Error("Missing required parameter: bank_local")
		return lib.CustomError(http.StatusBadRequest, "bank_local can not be blank", "bank_local can not be blank")
	} else {
		if bankLocal != "1" && bankLocal != "0" {
			log.Error("Missing required parameter: bank_local")
			return lib.CustomError(http.StatusBadRequest, "bank_local must 1 / 0", "bank_local must 1 / 0")
		}
		params["flag_local"] = bankLocal
	}

	bankGovernment := c.FormValue("bank_government")
	if bankGovernment == "" {
		log.Error("Missing required parameter: bank_government")
		return lib.CustomError(http.StatusBadRequest, "bank_government can not be blank", "bank_government can not be blank")
	} else {
		if bankGovernment != "1" && bankGovernment != "0" {
			log.Error("Missing required parameter: bank_government")
			return lib.CustomError(http.StatusBadRequest, "bank_government must 1 / 0", "bank_government must 1 / 0")
		}
		params["flag_government"] = bankGovernment
	}

	bankWebUrl := c.FormValue("bank_web_url")
	if bankWebUrl != "" {
		params["bank_web_url"] = bankWebUrl
	}

	custodianProfile := c.FormValue("custodian_profile")
	if custodianProfile != "" {
		params["custodian_profile"] = custodianProfile
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.CreateMsCustodianBank(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminUpdateMsCustodianBank(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	custodianKey := c.FormValue("custodian_key")
	if custodianKey != "" {
		n, err := strconv.ParseUint(custodianKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: custodian_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: custodian_key", "Wrong input for parameter: custodian_key")
		}
		params["custodian_key"] = custodianKey
	} else {
		log.Error("Missing required parameter: custodian_key")
		return lib.CustomError(http.StatusBadRequest, "custodian_key can not be blank", "custodian_key can not be blank")
	}

	custodianCode := c.FormValue("custodian_code")
	if custodianCode == "" {
		log.Error("Missing required parameter: custodian_code")
		return lib.CustomError(http.StatusBadRequest, "custodian_code can not be blank", "custodian_code can not be blank")
	} else {
		//validate unique custodian_code
		var countData models.CountData
		status, err = models.CountMsCustodianBankValidateUnique(&countData, "custodian_code", custodianCode, custodianKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: custodian_code")
			return lib.CustomError(http.StatusBadRequest, "custodian_code already used", "custodian_code already used")
		}
		params["custodian_code"] = custodianCode
	}

	custodianShortName := c.FormValue("custodian_short_name")
	if custodianShortName == "" {
		log.Error("Missing required parameter: custodian_short_name")
		return lib.CustomError(http.StatusBadRequest, "custodian_short_name can not be blank", "custodian_short_name can not be blank")
	} else {
		//validate unique custodian_short_name
		var countData models.CountData
		status, err = models.CountMsCustodianBankValidateUnique(&countData, "custodian_short_name", custodianShortName, custodianKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: custodian_short_name")
			return lib.CustomError(http.StatusBadRequest, "custodian_short_name already used", "custodian_short_name already used")
		}
		params["custodian_short_name"] = custodianShortName
	}

	custodianFullName := c.FormValue("custodian_full_name")
	if custodianFullName == "" {
		log.Error("Missing required parameter: custodian_full_name")
		return lib.CustomError(http.StatusBadRequest, "custodian_full_name can not be blank", "custodian_full_name can not be blank")
	} else {
		params["custodian_full_name"] = custodianFullName
	}

	biMemberCode := c.FormValue("bi_member_code")
	if biMemberCode != "" {
		//validate unique bi_member_code
		var countData models.CountData
		status, err = models.CountMsCustodianBankValidateUnique(&countData, "bi_member_code", biMemberCode, custodianKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: bi_member_code")
			return lib.CustomError(http.StatusBadRequest, "bi_member_code already used", "bi_member_code already used")
		}
		params["bi_member_code"] = biMemberCode
	}

	swiftCode := c.FormValue("swift_code")
	if swiftCode != "" {
		//validate unique swift_code
		var countData models.CountData
		status, err = models.CountMsCustodianBankValidateUnique(&countData, "swift_code", swiftCode, custodianKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: swift_code")
			return lib.CustomError(http.StatusBadRequest, "swift_code already used", "swift_code already used")
		}
		params["swift_code"] = swiftCode
	}

	bankLocal := c.FormValue("bank_local")
	if bankLocal == "" {
		log.Error("Missing required parameter: bank_local")
		return lib.CustomError(http.StatusBadRequest, "bank_local can not be blank", "bank_local can not be blank")
	} else {
		if bankLocal != "1" && bankLocal != "0" {
			log.Error("Missing required parameter: bank_local")
			return lib.CustomError(http.StatusBadRequest, "bank_local must 1 / 0", "bank_local must 1 / 0")
		}
		params["flag_local"] = bankLocal
	}

	bankGovernment := c.FormValue("bank_government")
	if bankGovernment == "" {
		log.Error("Missing required parameter: bank_government")
		return lib.CustomError(http.StatusBadRequest, "bank_government can not be blank", "bank_government can not be blank")
	} else {
		if bankGovernment != "1" && bankGovernment != "0" {
			log.Error("Missing required parameter: bank_government")
			return lib.CustomError(http.StatusBadRequest, "bank_government must 1 / 0", "bank_government must 1 / 0")
		}
		params["flag_government"] = bankGovernment
	}

	bankWebUrl := c.FormValue("bank_web_url")
	if bankWebUrl != "" {
		params["bank_web_url"] = bankWebUrl
	}

	custodianProfile := c.FormValue("custodian_profile")
	if custodianProfile != "" {
		params["custodian_profile"] = custodianProfile
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.UpdateMsCustodianBank(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminDetailMsCustodianBank(c echo.Context) error {
	var err error

	custodianKey := c.Param("custodian_key")
	if custodianKey == "" {
		log.Error("Missing required parameter: custodian_key")
		return lib.CustomError(http.StatusBadRequest, "custodian_key can not be blank", "custodian_key can not be blank")
	} else {
		n, err := strconv.ParseUint(custodianKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: custodian_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: custodian_key", "Wrong input for parameter: custodian_key")
		}
	}
	var bank models.MsCustodianBank
	_, err = models.GetMsCustodianBank(&bank, custodianKey)
	if err != nil {
		log.Error("Custodian Bank not found")
		return lib.CustomError(http.StatusBadRequest, "Custodian Bank not found", "Custodian Bank not found")
	}

	responseData := make(map[string]interface{})
	responseData["custodian_key"] = bank.CustodianKey
	responseData["custodian_code"] = bank.CustodianCode
	responseData["custodian_short_name"] = bank.CustodianShortName
	if bank.CustodianFullName != nil {
		responseData["custodian_full_name"] = *bank.CustodianFullName
	} else {
		responseData["custodian_full_name"] = ""
	}
	if bank.BiMemberCode != nil {
		responseData["bi_member_code"] = *bank.BiMemberCode
	} else {
		responseData["bi_member_code"] = ""
	}
	if bank.SwiftCode != nil {
		responseData["swift_code"] = *bank.SwiftCode
	} else {
		responseData["swift_code"] = ""
	}
	responseData["bank_local"] = bank.FlagLocal
	responseData["bank_government"] = bank.FlagGoverment
	if bank.BankWebUrl != nil {
		responseData["bank_web_url"] = *bank.BankWebUrl
	} else {
		responseData["bank_web_url"] = ""
	}
	if bank.CustodianProfile != nil {
		responseData["custodian_profile"] = *bank.CustodianProfile
	} else {
		responseData["custodian_profile"] = ""
	}
	if bank.RecOrder != nil {
		responseData["rec_order"] = *bank.RecOrder
	} else {
		responseData["rec_order"] = ""
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
