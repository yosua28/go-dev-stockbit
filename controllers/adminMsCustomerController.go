package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func GetListCustomerIndividuInquiry(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

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

	items := []string{"cif", "full_name", "sid", "date_birth", "customer_key", "phone_mobile", "cif_suspend_flag", "mother_maiden_name", "ktp"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "c.customer_key"
			if orderBy == "cif" {
				orderByJoin = "c.unit_holder_idno"
			} else if orderBy == "full_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "sid" {
				orderByJoin = "c.sid"
			} else if orderBy == "date_birth" {
				orderByJoin = "pd.date_birth"
			} else if orderBy == "phone_mobile" {
				orderByJoin = "pd.phone_mobile"
			} else if orderBy == "cif_suspend_flag" {
				orderByJoin = "c.cif_suspend_flag"
			} else if orderBy == "mother_maiden_name" {
				orderByJoin = "pd.mother_maiden_name"
			} else if orderBy == "ktp" {
				orderByJoin = "pd.idcard_no"
			}

			params["orderBy"] = orderByJoin
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
		params["orderBy"] = "c.customer_key"
		params["orderType"] = "DESC"
	}

	params["c.investor_type"] = "263"

	paramsLike := make(map[string]string)

	cif := c.QueryParam("cif")
	if cif != "" {
		paramsLike["c.unit_holder_idno"] = cif
	}
	fullname := c.QueryParam("full_name")
	if fullname != "" {
		paramsLike["c.full_name"] = fullname
	}
	datebirth := c.QueryParam("date_birth")
	if datebirth != "" {
		paramsLike["pd.date_birth"] = datebirth
	}
	ktp := c.QueryParam("ktp")
	if ktp != "" {
		paramsLike["pd.idcard_no"] = ktp
	}
	mothermaidenname := c.QueryParam("mother_maiden_name")
	if mothermaidenname != "" {
		paramsLike["pd.mother_maiden_name"] = mothermaidenname
	}

	var customers []models.CustomerIndividuInquiry

	status, err = models.AdminGetAllCustomerIndividuInquery(&customers, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(customers) < 1 {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetAllCustomerIndividuInquery(&countData, params, paramsLike)
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
	response.Data = customers

	return c.JSON(http.StatusOK, response)
}

func GetListCustomerInstitutionInquiry(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

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

	items := []string{"customer_key", "full_name", "sid", "npwp", "cif", "institution", "cif_suspend_flag"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "c.customer_key"
			if orderBy == "cif" {
				orderByJoin = "c.unit_holder_idno"
			} else if orderBy == "full_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "sid" {
				orderByJoin = "c.sid"
			} else if orderBy == "institution" {
				orderByJoin = "pd.insti_full_name"
			} else if orderBy == "cif_suspend_flag" {
				orderByJoin = "c.cif_suspend_flag"
			} else if orderBy == "npwp" {
				orderByJoin = "pd.npwp_no"
			} else if orderBy == "ktp" {
				orderByJoin = "pd.idcard_no"
			}

			params["orderBy"] = orderByJoin
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
		params["orderBy"] = "c.customer_key"
		params["orderType"] = "DESC"
	}

	params["c.investor_type"] = "264"

	paramsLike := make(map[string]string)

	cif := c.QueryParam("cif")
	if cif != "" {
		paramsLike["c.unit_holder_idno"] = cif
	}
	fullname := c.QueryParam("full_name")
	if fullname != "" {
		paramsLike["c.full_name"] = fullname
	}
	npwp := c.QueryParam("npwp")
	if npwp != "" {
		paramsLike["pd.npwp_no"] = npwp
	}

	var customers []models.CustomerInstituionInquiry

	status, err = models.AdminGetAllCustomerInstitutionInquery(&customers, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(customers) < 1 {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetAllCustomerInstitutionInquery(&countData, params, paramsLike)
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
	response.Data = customers

	return c.JSON(http.StatusOK, response)
}

func GetDetailCustomerIndividu(c echo.Context) error {
	var err error

	keyStr := c.Param("key")
	if keyStr == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var customer models.CustomerIndividuInquiry
	_, err = models.AdminGetHeaderCustomerIndividu(&customer, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var oaCustomer []models.OaCustomer
	_, err = models.AdminGetAllOaByCustomerKey(&oaCustomer, keyStr)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(http.StatusNotFound)
		}
	}

	var responseData models.DetailCustomerIndividuInquiry
	responseData.Header = customer
	responseData.PersonalData = &oaCustomer

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetDetailCustomerInstitution(c echo.Context) error {
	var err error

	keyStr := c.Param("key")
	if keyStr == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var customer models.CustomerInstituionInquiry
	_, err = models.AdminGetHeaderCustomerInstitution(&customer, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var oaCustomer []models.OaCustomer
	_, err = models.AdminGetAllOaByCustomerKey(&oaCustomer, keyStr)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(http.StatusNotFound)
		}
	}

	var responseData models.DetailCustomerInstitutionInquiry
	responseData.Header = customer
	responseData.PersonalData = &oaCustomer

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetDetailCustomerInquiry(c echo.Context) error {
	var err error

	keyStr := c.Param("key")
	if keyStr == "" {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var customer models.DetailCustomerInquiry
	_, err = models.AdminGetHeaderDetailCustomer(&customer, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var customerHeader models.DetailHeaderCustomerInquiry

	customerHeader.InvestorType = customer.InvestorType
	customerHeader.CustomerKey = customer.CustomerKey
	customerHeader.Cif = customer.Cif
	customerHeader.FullName = customer.FullName
	customerHeader.SidNo = customer.SidNo
	customerHeader.CifSuspendFlag = customer.CifSuspendFlag

	if customer.InvestorType == "263" {
		customerHeader.DateBirth = customer.DateBirth
		customerHeader.IDcardNo = customer.IDcardNo
		customerHeader.PhoneMobile = customer.PhoneMobile
		customerHeader.MotherMaidenName = customer.MotherMaidenName
	} else {
		customerHeader.Npwp = customer.Npwp
		customerHeader.Institusion = customer.Institusion
	}

	var oaCustomer []models.OaCustomer
	_, err = models.AdminGetAllOaByCustomerKey(&oaCustomer, keyStr)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(http.StatusNotFound)
		}
	}

	var responseData models.DetailCustomerInquiryResponse
	responseData.Header = customerHeader
	responseData.PersonalData = &oaCustomer

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func DetailPersonalDataCustomerIndividu(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	//Get parameter limit
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	if oareq.CustomerKey == nil {
		log.Println("data belum jadi customer")
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.DetailPersonalDataCustomerIndividu

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"

	responseData.OaRequestKey = oareq.OaRequestKey
	date, _ := time.Parse(layout, oareq.OaEntryStart)
	responseData.OaEntryStart = date.Format(newLayout)
	date, _ = time.Parse(layout, oareq.OaEntryEnd)
	responseData.OaEntryEnd = date.Format(newLayout)

	var oaRequestLookupIds []string

	if oareq.OaRequestType != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRequestType, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRequestType, 10))
		}
	}
	if oareq.Oastatus != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
		}
	}
	if oareq.OaRiskLevel != nil {
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10))
		}
	}

	//gen lookup oa request
	var lookupOaReq []models.GenLookup
	if len(oaRequestLookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, oaRequestLookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupOaReq {
		gData[gen.LookupKey] = gen
	}

	if oareq.OaRequestType != nil {
		if n, ok := gData[*oareq.OaRequestType]; ok {
			responseData.OaRequestType = n.LkpName
		}
	}

	if oareq.OaRiskLevel != nil {
		if n, ok := gData[*oareq.OaRiskLevel]; ok {
			responseData.OaRiskLevel = n.LkpName
		}
	}

	if oareq.Oastatus != nil {
		if n, ok := gData[*oareq.Oastatus]; ok {
			responseData.Oastatus = *n.LkpName
		}
	}

	//check personal data by oa request key
	var oapersonal models.OaPersonalData
	strKey := strconv.FormatUint(oareq.OaRequestKey, 10)
	status, err = models.GetOaPersonalDataByOaRequestKey(&oapersonal, strKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.FullName = oapersonal.FullName
		responseData.IDCardNo = oapersonal.IDcardNo
		date, _ = time.Parse(layout, oapersonal.DateBirth)
		responseData.DateBirth = date.Format(newLayout)
		responseData.PhoneNumber = oapersonal.PhoneMobile
		responseData.EmailAddress = oapersonal.EmailAddress
		responseData.PlaceBirth = oapersonal.PlaceBirth
		responseData.PhoneHome = oapersonal.PhoneHome

		dir := config.BaseUrl + "/images/user/" + strconv.FormatUint(*oareq.UserLoginKey, 10) + "/"

		if oapersonal.PicKtp != nil && *oapersonal.PicKtp != "" {
			path := dir + *oapersonal.PicKtp
			responseData.PicKtp = &path
		}

		if oapersonal.PicSelfie != nil && *oapersonal.PicSelfie != "" {
			path := dir + *oapersonal.PicSelfie
			responseData.PicSelfie = &path
		}

		if oapersonal.PicSelfieKtp != nil && *oapersonal.PicSelfieKtp != "" {
			path := dir + *oapersonal.PicSelfieKtp
			responseData.PicSelfieKtp = &path
		}

		responseData.OccupCompany = oapersonal.OccupCompany
		responseData.OccupPhone = oapersonal.OccupPhone
		responseData.OccupWebURL = oapersonal.OccupWebUrl
		responseData.MotherMaidenName = oapersonal.MotherMaidenName
		responseData.BeneficialFullName = oapersonal.BeneficialFullName

		//mapping gen lookup
		var personalDataLookupIds []string
		if oapersonal.Gender != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Gender, 10))
			}
		}
		if oapersonal.MaritalStatus != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.MaritalStatus, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.MaritalStatus, 10))
			}
		}
		if oapersonal.Religion != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Religion, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Religion, 10))
			}
		}
		if oapersonal.Education != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Education, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Education, 10))
			}
		}
		if oapersonal.OccupJob != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupJob, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupJob, 10))
			}
		}
		if oapersonal.OccupPosition != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupPosition, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupPosition, 10))
			}
		}
		if oapersonal.AnnualIncome != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.AnnualIncome, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.AnnualIncome, 10))
			}
		}
		if oapersonal.SourceofFund != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.SourceofFund, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.SourceofFund, 10))
			}
		}
		if oapersonal.InvesmentObjectives != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.InvesmentObjectives, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.InvesmentObjectives, 10))
			}
		}
		if oapersonal.Correspondence != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Correspondence, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Correspondence, 10))
			}
		}
		if oapersonal.BeneficialRelation != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.BeneficialRelation, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.BeneficialRelation, 10))
			}
		}
		if oapersonal.EmergencyRelation != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.EmergencyRelation, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.EmergencyRelation, 10))
			}
		}
		if oapersonal.RelationType != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationType, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationType, 10))
			}
		}
		if oapersonal.RelationOccupation != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationOccupation, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationOccupation, 10))
			}
		}
		if oapersonal.RelationBusinessFields != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationBusinessFields, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.RelationBusinessFields, 10))
			}
		}
		if oapersonal.OccupBusinessFields != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupBusinessFields, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.OccupBusinessFields, 10))
			}
		}
		//gen lookup personal data
		var lookupPersonData []models.GenLookup
		if len(personalDataLookupIds) > 0 {
			status, err = models.GetGenLookupIn(&lookupPersonData, personalDataLookupIds, "lookup_key")
			if err != nil {
				if err != sql.ErrNoRows {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}
		}

		pData := make(map[uint64]models.GenLookup)
		for _, genLook := range lookupPersonData {
			pData[genLook.LookupKey] = genLook
		}

		if oapersonal.Gender != nil {
			if n, ok := pData[*oapersonal.Gender]; ok {
				responseData.Gender = n.LkpName
			}
		}
		if oapersonal.MaritalStatus != nil {
			if n, ok := pData[*oapersonal.MaritalStatus]; ok {
				responseData.MaritalStatus = n.LkpName
			}
		}
		if oapersonal.Religion != nil {
			if n, ok := pData[*oapersonal.Religion]; ok {
				responseData.Religion = n.LkpName
			}
		}

		var country models.MsCountry

		strCountry := strconv.FormatUint(oapersonal.Nationality, 10)
		status, err = models.GetMsCountry(&country, strCountry)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("Error Personal Data not Found")
				return lib.CustomError(status, err.Error(), "Personal data not found")
			}
		} else {
			responseData.Nationality = &country.CouName
		}

		if oapersonal.Education != nil {
			if n, ok := pData[*oapersonal.Education]; ok {
				responseData.Education = n.LkpName
			}
		}
		if oapersonal.OccupJob != nil {
			if n, ok := pData[*oapersonal.OccupJob]; ok {
				responseData.OccupJob = n.LkpName
			}
		}
		if oapersonal.OccupPosition != nil {
			if n, ok := pData[*oapersonal.OccupPosition]; ok {
				responseData.OccupPosition = n.LkpName
			}
		}
		if oapersonal.AnnualIncome != nil {
			if n, ok := pData[*oapersonal.AnnualIncome]; ok {
				responseData.AnnualIncome = n.LkpName
			}
		}
		if oapersonal.SourceofFund != nil {
			if n, ok := pData[*oapersonal.SourceofFund]; ok {
				responseData.SourceofFund = n.LkpName
			}
		}
		if oapersonal.InvesmentObjectives != nil {
			if n, ok := pData[*oapersonal.InvesmentObjectives]; ok {
				responseData.InvesmentObjectives = n.LkpName
			}
		}
		if oapersonal.Correspondence != nil {
			if n, ok := pData[*oapersonal.Correspondence]; ok {
				responseData.Correspondence = n.LkpName
			}
		}
		if oapersonal.BeneficialRelation != nil {
			if n, ok := pData[*oapersonal.BeneficialRelation]; ok {
				responseData.BeneficialRelation = n.LkpName
			}
		}
		if oapersonal.OccupBusinessFields != nil {
			if n, ok := pData[*oapersonal.OccupBusinessFields]; ok {
				responseData.OccupBusinessFields = n.LkpName
			}
		}

		//mapping idcard address &  domicile
		var postalAddressIds []string
		if oapersonal.IDcardAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.IDcardAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.IDcardAddressKey, 10))
			}
		}
		if oapersonal.DomicileAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.DomicileAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.DomicileAddressKey, 10))
			}
		}
		if oapersonal.OccupAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.OccupAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.OccupAddressKey, 10))
			}
		}
		var oaPstalAddressList []models.OaPostalAddress
		if len(postalAddressIds) > 0 {
			status, err = models.GetOaPostalAddressIn(&oaPstalAddressList, postalAddressIds, "postal_address_key")
			if err != nil {
				if err != sql.ErrNoRows {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}
		}

		postalData := make(map[uint64]models.OaPostalAddress)
		for _, posAdd := range oaPstalAddressList {
			postalData[posAdd.PostalAddressKey] = posAdd
		}

		if len(postalData) > 0 {
			if oapersonal.IDcardAddressKey != nil {
				if p, ok := postalData[*oapersonal.IDcardAddressKey]; ok {
					responseData.IDcardAddress.Address = p.AddressLine1
					responseData.IDcardAddress.PostalCode = p.PostalCode

					var cityIds []string
					if p.KabupatenKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KabupatenKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KabupatenKey, 10))
						}
					}
					if p.KecamatanKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KecamatanKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KecamatanKey, 10))
						}
					}

					var cityList []models.MsCity
					if len(cityIds) > 0 {
						status, err = models.GetMsCityIn(&cityList, cityIds, "city_key")
						if err != nil {
							if err != sql.ErrNoRows {
								log.Error(err.Error())
								return lib.CustomError(status, err.Error(), "Failed get data")
							}
						}
					}
					cityData := make(map[uint64]models.MsCity)
					for _, city := range cityList {
						cityData[city.CityKey] = city
					}
					if c, ok := cityData[*p.KabupatenKey]; ok {
						responseData.IDcardAddress.Kabupaten = &c.CityName
					}
					if c, ok := cityData[*p.KecamatanKey]; ok {
						responseData.IDcardAddress.Kecamatan = &c.CityName
					}
				}
			}
			if oapersonal.DomicileAddressKey != nil {
				if p, ok := postalData[*oapersonal.DomicileAddressKey]; ok {
					responseData.DomicileAddress.Address = p.AddressLine1
					responseData.DomicileAddress.PostalCode = p.PostalCode

					var cityIds []string
					if p.KabupatenKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KabupatenKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KabupatenKey, 10))
						}
					}
					if p.KecamatanKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KecamatanKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KecamatanKey, 10))
						}
					}

					var cityList []models.MsCity
					if len(cityIds) > 0 {
						status, err = models.GetMsCityIn(&cityList, cityIds, "city_key")
						if err != nil {
							if err != sql.ErrNoRows {
								log.Error(err.Error())
								return lib.CustomError(status, err.Error(), "Failed get data")
							}
						}
					}
					cityData := make(map[uint64]models.MsCity)
					for _, city := range cityList {
						cityData[city.CityKey] = city
					}
					if p.KabupatenKey != nil {
						if c, ok := cityData[*p.KabupatenKey]; ok {
							responseData.DomicileAddress.Kabupaten = &c.CityName
						}
					}
					if p.KecamatanKey != nil {
						if c, ok := cityData[*p.KecamatanKey]; ok {
							responseData.DomicileAddress.Kecamatan = &c.CityName
						}
					}
				}
			}
			if oapersonal.OccupAddressKey != nil {
				if p, ok := postalData[*oapersonal.OccupAddressKey]; ok {
					responseData.OccupAddressKey.Address = p.AddressLine1
					responseData.OccupAddressKey.PostalCode = p.PostalCode

					var cityIds []string
					if p.KabupatenKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KabupatenKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KabupatenKey, 10))
						}
					}
					if p.KecamatanKey != nil {
						if _, ok := lib.Find(cityIds, strconv.FormatUint(*p.KecamatanKey, 10)); !ok {
							cityIds = append(cityIds, strconv.FormatUint(*p.KecamatanKey, 10))
						}
					}

					var cityList []models.MsCity
					if len(cityIds) > 0 {
						status, err = models.GetMsCityIn(&cityList, cityIds, "city_key")
						if err != nil {
							if err != sql.ErrNoRows {
								log.Error(err.Error())
								return lib.CustomError(status, err.Error(), "Failed get data")
							}
						}
					}
					cityData := make(map[uint64]models.MsCity)
					for _, city := range cityList {
						cityData[city.CityKey] = city
					}
					if p.KabupatenKey != nil {
						if c, ok := cityData[*p.KabupatenKey]; ok {
							responseData.DomicileAddress.Kabupaten = &c.CityName
						}
					}
					if p.KecamatanKey != nil {
						if c, ok := cityData[*p.KecamatanKey]; ok {
							responseData.DomicileAddress.Kecamatan = &c.CityName
						}
					}
				}
			}
		}

		//mapping bank account
		if oapersonal.BankAccountKey != nil {
			var bankaccount models.MsBankAccount
			strBank := strconv.FormatUint(*oapersonal.BankAccountKey, 10)
			status, err = models.GetBankAccount(&bankaccount, strBank)
			if err != nil {
				if err != sql.ErrNoRows {
					return lib.CustomError(status)
				}
			} else {
				var bank models.MsBank
				strBankKey := strconv.FormatUint(bankaccount.BankKey, 10)
				status, err = models.GetMsBank(&bank, strBankKey)
				if err != nil {
					if err != sql.ErrNoRows {
						return lib.CustomError(status)
					}
				} else {
					responseData.BankAccount.BankName = bank.BankName
				}
				responseData.BankAccount.AccountNo = bankaccount.AccountNo
				responseData.BankAccount.AccountHolderName = bankaccount.AccountHolderName
				responseData.BankAccount.BranchName = bankaccount.BranchName
			}
		}

		//mapping relation
		if oapersonal.RelationType != nil {
			if n, ok := pData[*oapersonal.RelationType]; ok {
				responseData.Relation.RelationType = n.LkpName
			}
		}
		responseData.Relation.RelationFullName = oapersonal.RelationFullName
		if oapersonal.RelationOccupation != nil {
			if n, ok := pData[*oapersonal.RelationOccupation]; ok {
				responseData.Relation.RelationOccupation = n.LkpName
			}
		}
		if oapersonal.RelationBusinessFields != nil {
			if n, ok := pData[*oapersonal.RelationBusinessFields]; ok {
				responseData.Relation.RelationBusinessFields = n.LkpName
			}
		}

		//mapping emergency
		responseData.Emergency.EmergencyFullName = oapersonal.EmergencyFullName
		if oapersonal.EmergencyRelation != nil {
			if n, ok := pData[*oapersonal.EmergencyRelation]; ok {
				responseData.Emergency.EmergencyRelation = n.LkpName
			}
		}
		responseData.Emergency.EmergencyPhoneNo = oapersonal.EmergencyPhoneNo

		var oaRiskProfile []models.AdminOaRiskProfile
		status, err = models.AdminGetOaRiskProfile(&oaRiskProfile, strKey)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
		responseData.RiskProfile = oaRiskProfile

		//mapping oa risk profile quiz
		var oaRiskProfileQuiz []models.AdminOaRiskProfileQuiz
		status, err = models.AdminGetOaRiskProfileQuizByOaRequestKey(&oaRiskProfileQuiz, strKey)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
		if len(oaRiskProfileQuiz) > 0 {
			var questionIDs []string
			for _, quiz := range oaRiskProfileQuiz {
				if _, ok := lib.Find(questionIDs, strconv.FormatUint(quiz.QuizQuestionKey, 10)); !ok {
					questionIDs = append(questionIDs, strconv.FormatUint(quiz.QuizQuestionKey, 10))
				}
			}
			var optionDB []models.CmsQuizOptions
			status, err = models.GetCmsQuizOptionsIn(&optionDB, questionIDs, "quiz_question_key")
			if err != nil {
				if err != sql.ErrNoRows {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}

			optionData := make(map[uint64][]models.CmsQuizOptionsInfo)
			optionUserData := make(map[uint64]models.CmsQuizOptions)
			if len(optionDB) > 0 {
				for _, option := range optionDB {

					optionUserData[option.QuizOptionKey] = option

					var data models.CmsQuizOptionsInfo

					data.QuizOptionKey = option.QuizOptionKey
					if option.QuizOptionLabel != nil {
						data.QuizOptionLabel = *option.QuizOptionLabel
					}
					if option.QuizOptionTitle != nil {
						data.QuizOptionTitle = *option.QuizOptionTitle
					}
					if option.QuizOptionScore != nil {
						data.QuizOptionScore = *option.QuizOptionScore
					}
					if option.QuizOptionDefault != nil {
						data.QuizOptionDefault = *option.QuizOptionDefault
					}

					optionData[option.QuizQuestionKey] = append(optionData[option.QuizQuestionKey], data)
				}
			}

			var riskQuiz []models.RiskProfileQuiz

			for _, oaRiskQuiz := range oaRiskProfileQuiz {
				var risk models.RiskProfileQuiz

				risk.RiskProfileQuizKey = oaRiskQuiz.RiskProfileQuizKey
				if n, ok := optionUserData[oaRiskQuiz.QuizOptionKeyUser]; ok {
					risk.QuizOptionUser.QuizOptionKey = n.QuizOptionKey
					if n.QuizOptionLabel != nil {
						risk.QuizOptionUser.QuizOptionLabel = *n.QuizOptionLabel
					}
					if n.QuizOptionTitle != nil {
						risk.QuizOptionUser.QuizOptionTitle = *n.QuizOptionTitle
					}
					if n.QuizOptionScore != nil {
						risk.QuizOptionUser.QuizOptionScore = *n.QuizOptionScore
					}
					if n.QuizOptionDefault != nil {
						risk.QuizOptionUser.QuizOptionDefault = *n.QuizOptionDefault
					}
				}
				risk.QuizOptionScoreUser = oaRiskQuiz.QuizOptionScoreUser
				risk.QuizQuestionKey = oaRiskQuiz.QuizQuestionKey
				risk.HeaderQuizName = *oaRiskQuiz.HeaderQuizName
				risk.QuizTitle = oaRiskQuiz.QuizTitle

				if opt, ok := optionData[oaRiskQuiz.QuizQuestionKey]; ok {
					risk.Options = opt
				}

				riskQuiz = append(riskQuiz, risk)
			}
			responseData.RiskProfileQuiz = riskQuiz
		}

		//add response field Sinvest
		if oareq.CustomerKey != nil {
			var customer models.MsCustomer
			strCustomerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
			status, err = models.GetMsCustomer(&customer, strCustomerKey)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed get data")
				}
			}

			responseData.FirstName = customer.FirstName
			responseData.MiddleName = customer.MiddleName
			responseData.LastName = customer.LastName
			responseData.ClientCode = customer.ClientCode
			responseData.TinNumber = customer.TinNumber

			if customer.TinIssuanceDate != nil {
				layout := "2006-01-02 15:04:05"
				newLayout := "02 Jan 2006"
				date, _ := time.Parse(layout, *customer.TinIssuanceDate)
				oke := date.Format(newLayout)
				responseData.TinIssuanceDate = &oke
			}

			if customer.FatcaStatus != nil {
				var fatca models.GenLookup
				strLookKey := strconv.FormatUint(*customer.FatcaStatus, 10)
				status, err = models.GetGenLookup(&fatca, strLookKey)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Error(err.Error())
						return lib.CustomError(status, err.Error(), "Failed get data")
					}
				}
				responseData.FatcaStatus = fatca.LkpName
			}

			if customer.TinIssuanceCountry != nil {
				var country models.MsCountry
				strCountryKey := strconv.FormatUint(*customer.TinIssuanceCountry, 10)
				status, err = models.GetMsCountry(&country, strCountryKey)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Error(err.Error())
						return lib.CustomError(status, err.Error(), "Failed get data")
					}
				}
				responseData.TinIssuanceCountry = &country.CouName
			}
		} else {
			sliceName := strings.Fields(oapersonal.FullName)
			if len(sliceName) > 0 {
				responseData.FirstName = &sliceName[0]
				if len(sliceName) > 1 {
					responseData.MiddleName = &sliceName[1]
					if len(sliceName) > 2 {
						lastName := strings.Join(sliceName[2:len(sliceName)], " ")
						responseData.LastName = &lastName
					}
				}
			}
		}
	}

	//set customer
	var customer models.CustomerDetailPersonalData
	strCustomerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
	_, err = models.GetCustomerDetailPersonalData(&customer, strCustomerKey)
	if err == nil {
		responseData.Customer = &customer
	}

	if customer.InvestorType != "263" { //individu
		return lib.CustomError(http.StatusNotFound)
	}

	//mapping user approval
	var userApprovalIds []string
	if oareq.Check1References != nil {

		usercs, _ := strconv.ParseUint(*oareq.Check1References, 10, 64)
		if usercs > 0 {
			if _, ok := lib.Find(userApprovalIds, strconv.FormatUint(usercs, 10)); !ok {
				userApprovalIds = append(userApprovalIds, strconv.FormatUint(usercs, 10))
			}
		}
	}
	if oareq.Check2References != nil {
		userkyc, _ := strconv.ParseUint(*oareq.Check2References, 10, 64)
		if userkyc > 0 {
			if _, ok := lib.Find(userApprovalIds, strconv.FormatUint(userkyc, 10)); !ok {
				userApprovalIds = append(userApprovalIds, strconv.FormatUint(userkyc, 10))
			}
		}
	}

	//gen lookup personal data
	var userappr []models.ScUserLogin
	if len(userApprovalIds) > 0 {
		status, err = models.GetScUserLoginIn(&userappr, userApprovalIds, "user_login_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	usrData := make(map[uint64]models.ScUserLogin)
	for _, usr := range userappr {
		usrData[usr.UserLoginKey] = usr
	}

	//set approv cs
	if oareq.Check1References != nil {
		usercs, _ := strconv.ParseUint(*oareq.Check1References, 10, 64)
		if usercs > 0 {
			if n, ok := usrData[usercs]; ok {
				var approvecs models.ApprovalData
				approvecs.ApproveUser = n.UloginFullName
				if oareq.Check1Date != nil {
					date, _ = time.Parse(layout, *oareq.Check1Date)
					oke := date.Format(newLayout)
					approvecs.ApproveDate = &oke
				}
				approvecs.ApproveNotes = oareq.Check1Notes

				responseData.ApproveCS = &approvecs
			}
		}
	}

	//set approv kyc
	if oareq.Check2References != nil {
		userkyc, _ := strconv.ParseUint(*oareq.Check2References, 10, 64)
		if userkyc > 0 {
			if n, ok := usrData[userkyc]; ok {
				var approvekyc models.ApprovalData
				approvekyc.ApproveUser = n.UloginFullName
				if oareq.Check1Date != nil {
					date, _ = time.Parse(layout, *oareq.Check2Date)
					oke := date.Format(newLayout)
					approvekyc.ApproveDate = &oke
				}
				approvekyc.ApproveNotes = oareq.Check2Notes

				responseData.ApproveKYC = &approvekyc
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
