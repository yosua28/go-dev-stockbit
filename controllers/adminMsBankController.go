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

func AdminGetListMsBank(c echo.Context) error {

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

	items := []string{"bank_code", "bank_name", "bank_fullname", "bi_member_code", "swift_code", "flag_local", "flag_government"}

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
		params["orderBy"] = "bank_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	bankLocal := c.QueryParam("bank_local")
	if bankLocal != "" {
		params["flag_local"] = bankLocal
	}

	bankGovernment := c.QueryParam("bank_government")
	if bankGovernment != "" {
		params["flag_government"] = bankGovernment
	}

	var bank []models.ListBankAdmin

	status, err = models.AdminGetListBank(&bank, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(bank) < 1 {
		log.Error("Bank not found")
		return lib.CustomError(http.StatusNotFound, "Bank not found", "Bank not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListBank(&countData, params, searchLike)
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

func AdminDeleteMsBank(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("bank_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key", "Missing required parameter: bank_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["bank_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsBank(params)
	if err != nil {
		log.Error("Error delete ms_bank")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateMsBank(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	bankCode := c.FormValue("bank_code")
	if bankCode == "" {
		log.Error("Missing required parameter: bank_code")
		return lib.CustomError(http.StatusBadRequest, "bank_code can not be blank", "bank_code can not be blank")
	} else {
		//validate unique bank_code
		var countData models.CountData
		status, err = models.CountMsBankValidateUnique(&countData, "bank_code", bankCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: bank_code")
			return lib.CustomError(http.StatusBadRequest, "bank_code already used", "bank_code already used")
		}
		params["bank_code"] = bankCode
	}

	bankName := c.FormValue("bank_name")
	if bankName == "" {
		log.Error("Missing required parameter: bank_name")
		return lib.CustomError(http.StatusBadRequest, "bank_name can not be blank", "bank_name can not be blank")
	} else {
		//validate unique bank_name
		var countData models.CountData
		status, err = models.CountMsBankValidateUnique(&countData, "bank_name", bankName, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: bank_name")
			return lib.CustomError(http.StatusBadRequest, "bank_name already used", "bank_name already used")
		}
		params["bank_name"] = bankName
	}

	bankFullName := c.FormValue("bank_fullname")
	if bankFullName == "" {
		log.Error("Missing required parameter: bank_fullname")
		return lib.CustomError(http.StatusBadRequest, "bank_fullname can not be blank", "bank_fullname can not be blank")
	} else {
		params["bank_fullname"] = bankFullName
	}

	biMemberCode := c.FormValue("bi_member_code")
	if biMemberCode != "" {
		//validate unique bi_member_code
		var countData models.CountData
		status, err = models.CountMsBankValidateUnique(&countData, "bi_member_code", biMemberCode, "")
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
	if swiftCode == "" {
		log.Error("Missing required parameter: swift_code")
		return lib.CustomError(http.StatusBadRequest, "swift_code can not be blank", "swift_code can not be blank")
	} else {
		//validate unique swift_code
		var countData models.CountData
		status, err = models.CountMsBankValidateUnique(&countData, "swift_code", swiftCode, "")
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

	bankIbankUrl := c.FormValue("bank_ibank_url")
	if bankIbankUrl != "" {
		params["bank_ibank_url"] = bankIbankUrl
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

	status, err = models.CreateMsBank(params)
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

func AdminUpdateMsBank(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	bankKey := c.FormValue("bank_key")
	if bankKey != "" {
		n, err := strconv.ParseUint(bankKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bank_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_key", "Wrong input for parameter: bank_key")
		}
		params["bank_key"] = bankKey
	} else {
		log.Error("Missing required parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "bank_key can not be blank", "bank_key can not be blank")
	}

	bankCode := c.FormValue("bank_code")
	if bankCode == "" {
		log.Error("Missing required parameter: bank_code")
		return lib.CustomError(http.StatusBadRequest, "bank_code can not be blank", "bank_code can not be blank")
	} else {
		//validate unique bank_code
		var countData models.CountData
		status, err = models.CountMsBankValidateUnique(&countData, "bank_code", bankCode, bankKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: bank_code")
			return lib.CustomError(http.StatusBadRequest, "bank_code already used", "bank_code already used")
		}
		params["bank_code"] = bankCode
	}

	bankName := c.FormValue("bank_name")
	if bankName == "" {
		log.Error("Missing required parameter: bank_name")
		return lib.CustomError(http.StatusBadRequest, "bank_name can not be blank", "bank_name can not be blank")
	} else {
		//validate unique bank_name
		var countData models.CountData
		status, err = models.CountMsBankValidateUnique(&countData, "bank_name", bankName, bankKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: bank_name")
			return lib.CustomError(http.StatusBadRequest, "bank_name already used", "bank_name already used")
		}
		params["bank_name"] = bankName
	}

	bankFullName := c.FormValue("bank_fullname")
	if bankFullName == "" {
		log.Error("Missing required parameter: bank_fullname")
		return lib.CustomError(http.StatusBadRequest, "bank_fullname can not be blank", "bank_fullname can not be blank")
	} else {
		params["bank_fullname"] = bankFullName
	}

	biMemberCode := c.FormValue("bi_member_code")
	if biMemberCode != "" {
		//validate unique bi_member_code
		var countData models.CountData
		status, err = models.CountMsBankValidateUnique(&countData, "bi_member_code", biMemberCode, bankKey)
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
	if swiftCode == "" {
		log.Error("Missing required parameter: swift_code")
		return lib.CustomError(http.StatusBadRequest, "swift_code can not be blank", "swift_code can not be blank")
	} else {
		//validate unique swift_code
		var countData models.CountData
		status, err = models.CountMsBankValidateUnique(&countData, "swift_code", swiftCode, bankKey)
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

	bankIbankUrl := c.FormValue("bank_ibank_url")
	if bankIbankUrl != "" {
		params["bank_ibank_url"] = bankIbankUrl
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

	status, err = models.UpdateMsBank(params)
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

func AdminDetailBank(c echo.Context) error {
	var err error

	bankKey := c.Param("bank_key")
	if bankKey == "" {
		log.Error("Missing required parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "bank_key can not be blank", "bank_key can not be blank")
	} else {
		n, err := strconv.ParseUint(bankKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: bank_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: bank_key", "Wrong input for parameter: bank_key")
		}
	}
	var bank models.MsBank
	_, err = models.GetMsBank(&bank, bankKey)
	if err != nil {
		log.Error("Bank not found")
		return lib.CustomError(http.StatusBadRequest, "Bank not found", "Bank not found")
	}

	responseData := make(map[string]interface{})
	responseData["bank_key"] = bank.BankKey
	responseData["bank_code"] = bank.BankCode
	responseData["bank_name"] = bank.BankName
	if bank.BankFullname != nil {
		responseData["bank_fullname"] = *bank.BankFullname
	} else {
		responseData["bank_fullname"] = ""
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
	if bank.BankIbankUrl != nil {
		responseData["bank_ibank_url"] = *bank.BankIbankUrl
	} else {
		responseData["bank_ibank_url"] = ""
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
