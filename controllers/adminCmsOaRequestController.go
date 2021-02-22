package controllers

import (
	"api/config"
	"api/db"
	"api/lib"
	"api/models"
	"bytes"
	"crypto/tls"
	"database/sql"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func initAuthCs() error {
	var roleKeyCs uint64
	roleKeyCs = 11

	if lib.Profile.RoleKey != roleKeyCs {
		return lib.CustomError(http.StatusBadRequest, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthKyc() error {
	var roleKeyKyc uint64
	roleKeyKyc = 12

	if lib.Profile.RoleKey != roleKeyKyc {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthFundAdmin() error {
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	if lib.Profile.RoleKey != roleKeyFundAdmin {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthCsKyc() error {
	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	if (lib.Profile.RoleKey != roleKeyCs) && (lib.Profile.RoleKey != roleKeyKyc) {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthCsKycFundAdmin() error {
	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	if (lib.Profile.RoleKey != roleKeyCs) && (lib.Profile.RoleKey != roleKeyKyc) && (lib.Profile.RoleKey != roleKeyFundAdmin) {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func GetOaRequestList(c echo.Context) error {

	errorAuth := initAuthCsKyc()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	var err error
	var status int

	oaStatusCs := "258"
	oaStatusKyc := "259"

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

	items := []string{"oa_request_key", "oa_request_type", "oa_entry_start", "oa_entry_end", "oa_status", "rec_order", "rec_status", "oa_risk_level"}

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
		params["orderBy"] = "oa_request_key"
		params["orderType"] = "DESC"
	}

	//if user approval CS
	if lib.Profile.RoleKey == roleKeyCs {
		params["oa_status"] = oaStatusCs
	}
	//if user approval KYC / Complainer
	if lib.Profile.RoleKey == roleKeyKyc {
		params["oa_status"] = oaStatusKyc
	}
	params["rec_status"] = "1"

	var oaRequestDB []models.OaRequest
	status, err = models.GetAllOaRequest(&oaRequestDB, limit, offset, noLimit, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaRequestDB) < 1 {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var lookupIds []string
	var oaRequestIds []string
	for _, oareq := range oaRequestDB {

		if oareq.Oastatus != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
			}
		}

		if _, ok := lib.Find(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10)); !ok {
			oaRequestIds = append(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10))
		}
	}

	//mapping lookup
	var genLookup []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&genLookup, lookupIds, "lookup_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	gData := make(map[uint64]models.GenLookup)
	for _, gen := range genLookup {
		gData[gen.LookupKey] = gen
	}

	//mapping personal data
	var oaPersonalData []models.OaPersonalData
	if len(oaRequestIds) > 0 {
		status, err = models.GetOaPersonalDataIn(&oaPersonalData, oaRequestIds, "oa_request_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	pdData := make(map[uint64]models.OaPersonalData)
	for _, oaPD := range oaPersonalData {
		pdData[oaPD.OaRequestKey] = oaPD
	}

	var responseData []models.OaRequestListResponse
	for _, oareq := range oaRequestDB {
		var data models.OaRequestListResponse

		if oareq.Oastatus != nil {
			if n, ok := gData[*oareq.Oastatus]; ok {
				data.Oastatus = *n.LkpName
			}
		}

		data.OaRequestKey = oareq.OaRequestKey

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, _ := time.Parse(layout, oareq.OaEntryStart)
		data.OaEntryStart = date.Format(newLayout)
		date, _ = time.Parse(layout, oareq.OaEntryEnd)
		data.OaEntryEnd = date.Format(newLayout)

		if n, ok := pdData[oareq.OaRequestKey]; ok {
			data.EmailAddress = n.EmailAddress
			data.PhoneNumber = n.PhoneMobile
			date, _ = time.Parse(layout, n.DateBirth)
			data.DateBirth = date.Format(newLayout)
			data.FullName = n.FullName
			data.IDCardNo = n.IDcardNo
		}

		responseData = append(responseData, data)
	}

	var countData models.OaRequestCountData
	var pagination int
	if limit > 0 {
		status, err = models.GetCountOaRequest(&countData, params)
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

func GetOaRequestData(c echo.Context) error {
	keyStr := c.Param("key")
	return ResultOaRequestData(keyStr, c, false)
}

func GetLastHistoryOaRequestData(c echo.Context) error {
	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err := models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var lastKeyStr string

	if oareq.OaRequestType == nil {
		log.Error("OA Request Type Null")
		return lib.CustomError(http.StatusBadRequest)
	} else if *oareq.OaRequestType == 127 { //NEW error tidak ada history
		log.Error("OA Request Type NEW harusnya UPDATE")
		return lib.CustomError(http.StatusBadRequest)
	} else if *oareq.OaRequestType == 128 {
		if oareq.CustomerKey == nil { //Error jika belum jadi customer
			return lib.CustomError(http.StatusBadRequest)
		}
		var lastHistoryOareq models.OaRequestKeyLastHistory
		customerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
		status, err := models.AdminGetLastHistoryOaRequest(&lastHistoryOareq, customerKey, keyStr)
		if err != nil {
			return lib.CustomError(status)
		}
		lastKeyStr = strconv.FormatUint(lastHistoryOareq.OaRequestKey, 10)
	}

	if lastKeyStr == "" {
		return lib.CustomError(http.StatusBadRequest)
	}

	return ResultOaRequestData(lastKeyStr, c, true)
}

func ResultOaRequestData(keyStr string, c echo.Context, isHistory bool) error {
	errorAuth := initAuthCsKycFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	//Get parameter limit
	// keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	log.Println(lib.Profile.RoleKey)

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	if lib.Profile.RoleKey == roleKeyCs {
		if isHistory == false {
			oaStatusCs := strconv.FormatUint(uint64(258), 10)
			if strOaKey != oaStatusCs {
				log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		if isHistory == false {
			oaStatusKyc := strconv.FormatUint(uint64(259), 10)
			if strOaKey != oaStatusKyc {
				log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	if lib.Profile.RoleKey == roleKeyFundAdmin {
		if isHistory == false {
			oaStatusFundAdmin1 := strconv.FormatUint(uint64(260), 10)
			oaStatusFundAdmin2 := strconv.FormatUint(uint64(261), 10)
			if (strOaKey != oaStatusFundAdmin1) && (strOaKey != oaStatusFundAdmin2) {
				log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	var responseData models.OaRequestDetailResponse

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
					if p.KabupatenKey != nil {
						if c, ok := cityData[*p.KabupatenKey]; ok {
							responseData.IDcardAddress.Kabupaten = &c.CityName
						}
					}

					if p.KecamatanKey != nil {
						if c, ok := cityData[*p.KecamatanKey]; ok {
							responseData.IDcardAddress.Kecamatan = &c.CityName
						}
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

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func UpdateStatusApprovalCS(c echo.Context) error {
	errorAuth := initAuthCs()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int

	params := make(map[string]string)

	oastatus := c.FormValue("oa_status") //259 = approve --------- 258 = reject
	if oastatus == "" {
		log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err := strconv.ParseUint(oastatus, 10, 64)
	if err == nil && n > 0 {
		if (oastatus != "259") && (oastatus != "258") {
			log.Error("Wrong input for parameter: oa_status must 259/258")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
		}
		params["oa_status"] = oastatus
	} else {
		log.Error("Wrong input for parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
	}

	dateLayout := "2006-01-02 15:04:05"
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)

	check1notes := c.FormValue("notes")
	params["check1_notes"] = check1notes

	if oastatus != "259" { //jika reject
		if check1notes == "" {
			log.Error("Missing required parameter notes: Notes tidak boleh kosong")
			return lib.CustomError(http.StatusBadRequest, "Notes tidak boleh kosong", "Notes tidak boleh kosong")
		}
		params["rec_status"] = "0"
		params["rec_deleted_date"] = time.Now().Format(dateLayout)
		params["rec_deleted_by"] = strKey
	}

	oarequestkey := c.FormValue("oa_request_key")
	if oarequestkey == "" {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err = strconv.ParseUint(oarequestkey, 10, 64)
	if err == nil && n > 0 {
		params["oa_request_key"] = oarequestkey
	} else {
		log.Error("Wrong input for parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
	}

	params["check1_date"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["check1_flag"] = "1"
	params["check1_references"] = strKey
	params["rec_modified_by"] = strKey

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, oarequestkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	oaStatusCs := strconv.FormatUint(uint64(258), 10)
	if strOaKey != oaStatusCs {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	_, err = models.UpdateOaRequest(params)
	if err != nil {
		log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	log.Info("Success update approved CS")

	if oastatus == "259" { //jika approve
		//send email to KYC
		var oapersonal models.OaPersonalData
		strKeyOa := strconv.FormatUint(oareq.OaRequestKey, 10)
		status, err = models.GetOaPersonalDataByOaRequestKey(&oapersonal, strKeyOa)
		if err != nil {
			log.Error("Error Personal Data not Found")
			return lib.CustomError(status, err.Error(), "Personal data not found")
		}

		paramsScLogin := make(map[string]string)
		paramsScLogin["role_key"] = "12"
		paramsScLogin["rec_status"] = "1"
		var userLogin []models.ScUserLogin
		_, err = models.GetAllScUserLogin(&userLogin, 0, 0, paramsScLogin, true)
		if err != nil {
			log.Error("Error get email")
			log.Error(err)
		}

		for _, scLogin := range userLogin {
			strUserCat := strconv.FormatUint(scLogin.UserCategoryKey, 10)
			if (strUserCat == "2") || (strUserCat == "3") {
				mailer := gomail.NewMessage()
				mailer.SetHeader("From", config.EmailFrom)
				mailer.SetHeader("To", scLogin.UloginEmail)
				mailer.SetHeader("Subject", "[MNC Duit] Verifikasi Opening Account")
				mailer.SetBody("text/html", "Segera verifikasi opening account baru dengan nama : "+oapersonal.FullName)
				dialer := gomail.NewDialer(
					config.EmailSMTPHost,
					int(config.EmailSMTPPort),
					config.EmailFrom,
					config.EmailFromPassword,
				)
				dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

				err = dialer.DialAndSend(mailer)
				if err != nil {
					log.Error("Error send email")
					log.Error(err)
				}
			}
		}
		//end send email to KYC
	} else {
		//create user message
		paramsUserMessage := make(map[string]string)
		paramsUserMessage["umessage_type"] = "245"
		strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)
		if oareq.UserLoginKey != nil {
			paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
		} else {
			paramsUserMessage["umessage_recipient_key"] = "0"
		}
		paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["flag_read"] = "0"
		paramsUserMessage["umessage_sender_key"] = strKey
		paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["flag_sent"] = "1"
		subject := "Pembukaan Rekening kamu ditolak"
		body := check1notes + " Silakan menghubungi Customer Service untuk informasi lebih lanjut."
		paramsUserMessage["umessage_body"] = body
		paramsUserMessage["umessage_subject"] = subject
		paramsUserMessage["umessage_category"] = "248"
		paramsUserMessage["flag_archieved"] = "0"
		paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["rec_status"] = "1"
		paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["rec_created_by"] = strKey

		status, err = models.CreateScUserMessage(paramsUserMessage)
		if err != nil {
			log.Error("Error create user message")
		}
		lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func UpdateStatusApprovalCompliance(c echo.Context) error {
	errorAuth := initAuthKyc()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int

	params := make(map[string]string)

	oastatus := c.FormValue("oa_status") //260 = app --- 258 = reject
	if oastatus == "" {
		log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err := strconv.ParseUint(oastatus, 10, 64)
	if err == nil && n > 0 {
		if (oastatus != "260") && (oastatus != "258") {
			log.Error("Wrong input for parameter: oa_status must 260/258")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
		}
		if oastatus == "260" {
			params["oa_status"] = oastatus
		}
	} else {
		log.Error("Wrong input for parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
	}

	check2notes := c.FormValue("notes")
	params["check2_notes"] = check2notes

	dateLayout := "2006-01-02 15:04:05"
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)

	if oastatus != "260" { //jika reject
		if check2notes == "" {
			log.Error("Missing required parameter notes: Notes tidak boleh kosong")
			return lib.CustomError(http.StatusBadRequest, "Notes tidak boleh kosong", "Notes tidak boleh kosong")
		}
		params["rec_status"] = "0"
		params["rec_deleted_date"] = time.Now().Format(dateLayout)
		params["rec_deleted_by"] = strKey
	}

	oarequestkey := c.FormValue("oa_request_key")
	if oarequestkey == "" {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err = strconv.ParseUint(oarequestkey, 10, 64)
	if err == nil && n > 0 {
		params["oa_request_key"] = oarequestkey
	} else {
		log.Error("Wrong input for parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
	}

	oarisklevel := c.FormValue("oa_risk_level")
	if oarisklevel == "" {
		log.Error("Missing required parameter: oa_risk_level")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err = strconv.ParseUint(oarisklevel, 10, 64)
	if err == nil && n > 0 {
		params["oa_risk_level"] = oarisklevel
	} else {
		log.Error("Wrong input for parameter: oa_risk_level")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
	}

	params["check2_date"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["check2_flag"] = "1"
	params["check2_references"] = strKey
	params["rec_modified_by"] = strKey

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, oarequestkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	oaStatusKyc := strconv.FormatUint(uint64(259), 10)
	if strOaKey != oaStatusKyc {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	tx, err := db.Db.Begin()

	//cek rec order
	if oareq.CustomerKey == nil {
		params["rec_order"] = "0"
	} else {
		params["rec_order"] = "0"

		var lastHistoryOareq models.OaRequestKeyLastHistory
		customerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
		_, err := models.AdminGetLastHistoryOaRequest(&lastHistoryOareq, customerKey, oarequestkey)
		if err == nil {
			if lastHistoryOareq.RecOrder != nil {
				lastOrder := *lastHistoryOareq.RecOrder + 1
				params["rec_order"] = strconv.FormatUint(lastOrder, 10)
			}
		}
	}

	//update oa request
	_, err = models.UpdateOaRequest(params)
	if err != nil {
		tx.Rollback()
		log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}
	log.Info("Success update approved Compliance Transaction")

	if oastatus == "260" {
		var oapersonal models.OaPersonalData
		strKeyOa := strconv.FormatUint(oareq.OaRequestKey, 10)
		status, err = models.GetOaPersonalDataByOaRequestKey(&oapersonal, strKeyOa)
		if err != nil {
			tx.Rollback()
			log.Error("Error Personal Data not Found")
			return lib.CustomError(status, err.Error(), "Personal data not found")
		}

		if oareq.CustomerKey == nil { //NEW OA
			//create customer

			paramsCustomer := make(map[string]string)
			paramsCustomer["id_customer"] = "0"

			year, month, _ := time.Now().Date()

			var customer models.MsCustomer
			tahun := strconv.FormatUint(uint64(year), 10)
			bulan := strconv.FormatUint(uint64(month), 10)
			unitHolderLike := tahun + bulan
			status, err = models.GetLastUnitHolder(&customer, unitHolderLike)
			if err != nil {
				paramsCustomer["unit_holder_idno"] = unitHolderLike + "000001"
			} else {
				dgt := customer.UnitHolderIDno[len(customer.UnitHolderIDno)-6:]
				productKeyCek, _ := strconv.ParseUint(dgt, 10, 64)
				hasil := strconv.FormatUint(productKeyCek+1, 10)
				if len(hasil) == 1 {
					hasil = "00000" + hasil
				} else if len(hasil) == 2 {
					hasil = "0000" + hasil
				} else if len(hasil) == 3 {
					hasil = "000" + hasil
				} else if len(hasil) == 4 {
					hasil = "00" + hasil
				} else if len(hasil) == 5 {
					hasil = "0" + hasil
				}

				resultData := unitHolderLike + hasil
				paramsCustomer["unit_holder_idno"] = resultData
			}

			paramsCustomer["full_name"] = oapersonal.FullName
			paramsCustomer["investor_type"] = "263"
			paramsCustomer["customer_category"] = "265"
			paramsCustomer["cif_suspend_flag"] = "0"
			paramsCustomer["openacc_branch_key"] = "1"
			paramsCustomer["openacc_agent_key"] = "1"
			paramsCustomer["openacc_date"] = time.Now().Format(dateLayout)
			paramsCustomer["flag_employee"] = "0"
			paramsCustomer["flag_group"] = "0"
			paramsCustomer["merging_flag"] = "0"
			paramsCustomer["rec_status"] = "1"
			paramsCustomer["rec_created_date"] = time.Now().Format(dateLayout)
			paramsCustomer["rec_created_by"] = strKey

			sliceName := strings.Fields(oapersonal.FullName)

			if len(sliceName) > 0 {
				paramsCustomer["first_name"] = sliceName[0]
				if len(sliceName) > 1 {
					paramsCustomer["middle_name"] = sliceName[1]
					if len(sliceName) > 2 {
						lastName := strings.Join(sliceName[2:len(sliceName)], " ")
						paramsCustomer["last_name"] = lastName
					}
				}
			}

			strNationality := strconv.FormatUint(oapersonal.Nationality, 10)
			if strNationality == "97" {
				paramsCustomer["fatca_status"] = "278"
			} else if strNationality == "225" {
				paramsCustomer["fatca_status"] = "279"
			} else {
				paramsCustomer["fatca_status"] = "280"
			}

			status, err, requestID := models.CreateMsCustomer(paramsCustomer)
			if err != nil {
				tx.Rollback()
				log.Error("Error create customer")
				return lib.CustomError(status, err.Error(), "failed input data")
			}
			request, err := strconv.ParseUint(requestID, 10, 64)
			if request == 0 {
				tx.Rollback()
				log.Error("Failed create customer")
				return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
			}

			paramOaUpdate := make(map[string]string)
			paramOaUpdate["customer_key"] = requestID
			paramOaUpdate["oa_request_key"] = oarequestkey

			_, err = models.UpdateOaRequest(paramOaUpdate)
			if err != nil {
				tx.Rollback()
				log.Error("Error update oa request")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
			}

			//create user message
			paramsUserMessage := make(map[string]string)
			paramsUserMessage["umessage_type"] = "245"
			strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)
			if oareq.UserLoginKey != nil {
				paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
			} else {
				paramsUserMessage["umessage_recipient_key"] = "0"
			}
			paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_read"] = "0"
			paramsUserMessage["umessage_sender_key"] = strKey
			paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_sent"] = "1"
			subject := "Selamat! Pembukaan Rekening telah Disetujui"
			body := "Saat ini akun kamu sudah aktif dan bisa melakukan transaksi. Yuk mulai investasi sekarang juga."
			paramsUserMessage["umessage_subject"] = subject
			paramsUserMessage["umessage_body"] = body
			paramsUserMessage["umessage_category"] = "248"
			paramsUserMessage["flag_archieved"] = "0"
			paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_status"] = "1"
			paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_created_by"] = strKey

			status, err = models.CreateScUserMessage(paramsUserMessage)
			if err != nil {
				tx.Rollback()
				log.Error("Error create user message")
				return lib.CustomError(status, err.Error(), "failed input data")
			}
			lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body)

			//update sc user login
			paramsUserLogin := make(map[string]string)
			paramsUserLogin["customer_key"] = requestID
			paramsUserLogin["rec_modified_date"] = time.Now().Format(dateLayout)
			paramsUserLogin["rec_modified_by"] = strKey
			paramsUserLogin["role_key"] = "1"
			strUserLoginKeyOa := strconv.FormatUint(*oareq.UserLoginKey, 10)
			paramsUserLogin["user_login_key"] = strUserLoginKeyOa
			_, err = models.UpdateScUserLogin(paramsUserLogin)
			if err != nil {
				tx.Rollback()
				log.Error("Error update oa request")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
			}

			//create agent customer
			paramsAgentCustomer := make(map[string]string)
			paramsAgentCustomer["customer_key"] = requestID

			paramsAgentCustomer["agent_key"] = "1"
			if oareq.SalesCode != nil {
				var agent models.MsAgent
				var salcode string
				salcode = *oareq.SalesCode
				status, err = models.GetMsAgentByField(&agent, salcode, "agent_code")
				if err == nil {
					paramsAgentCustomer["agent_key"] = strconv.FormatUint(agent.AgentKey, 10)
				}
			}

			paramsAgentCustomer["rec_status"] = "1"
			paramsAgentCustomer["eff_date"] = oareq.OaEntryStart
			paramsAgentCustomer["rec_created_date"] = time.Now().Format(dateLayout)
			paramsAgentCustomer["rec_created_by"] = strKey
			status, err = models.CreateMsAgentCustomer(paramsAgentCustomer)
			if err != nil {
				tx.Rollback()
				log.Error("Error create agent customer")
				return lib.CustomError(status, err.Error(), "failed input data")
			}

			tx.Commit()

			log.Info("Success create customer")

			//send email to customer
			var userData models.ScUserLogin
			status, err = models.GetScUserLoginByCustomerKey(&userData, requestID)
			if err == nil {
				sendEmailApproveOa(oapersonal.FullName, userData.UloginEmail)
			}

			// send email to fund admin
			paramsScLogin := make(map[string]string)
			paramsScLogin["role_key"] = "13"
			paramsScLogin["rec_status"] = "1"
			var userLogin []models.ScUserLogin
			_, err = models.GetAllScUserLogin(&userLogin, 0, 0, paramsScLogin, true)
			if err != nil {
				log.Error("Error get email")
				log.Error(err)
			}

			for _, scLogin := range userLogin {
				strUserCat := strconv.FormatUint(scLogin.UserCategoryKey, 10)
				if (strUserCat == "2") || (strUserCat == "3") {
					mailer := gomail.NewMessage()
					mailer.SetHeader("From", config.EmailFrom)
					mailer.SetHeader("To", scLogin.UloginEmail)
					mailer.SetHeader("Subject", "[MNC Duit] Verifikasi Opening Account")
					mailer.SetBody("text/html", "Segera verifikasi opening account baru dengan nama : "+oapersonal.FullName)
					dialer := gomail.NewDialer(
						config.EmailSMTPHost,
						int(config.EmailSMTPPort),
						config.EmailFrom,
						config.EmailFromPassword,
					)
					dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

					err = dialer.DialAndSend(mailer)
					if err != nil {
						log.Error("Error send email")
						log.Error(err)
					}
				}
			}
			// end send email to fund admin

			// insert into table ms_customer_bank_account
			paramsCusBankAcc := make(map[string]string)
			strBankAccKey := strconv.FormatUint(*oapersonal.BankAccountKey, 10)
			paramsCusBankAcc["customer_key"] = requestID
			paramsCusBankAcc["bank_account_key"] = strBankAccKey
			//get ms abank account
			var bankaccount models.MsBankAccount
			status, err = models.GetBankAccount(&bankaccount, strBankAccKey)
			if err == nil {
				paramsCusBankAcc["bank_account_name"] = bankaccount.AccountHolderName
			}
			paramsCusBankAcc["flag_priority"] = "1"
			paramsCusBankAcc["rec_status"] = "1"
			paramsCusBankAcc["rec_created_date"] = time.Now().Format(dateLayout)
			paramsCusBankAcc["rec_created_by"] = strKey
			status, err = models.CreateMsCustomerBankAccount(paramsCusBankAcc)
			if err != nil {
				tx.Rollback()
				log.Error("Error create ms_customer_bank_account")
				return lib.CustomError(status, err.Error(), "failed input data")
			}
			//end insert into table ms_customer_bank_account
		} else { //create user message
			paramsUserMessage := make(map[string]string)
			paramsUserMessage["umessage_type"] = "245"
			strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)
			if oareq.UserLoginKey != nil {
				paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
			} else {
				paramsUserMessage["umessage_recipient_key"] = "0"
			}
			paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_read"] = "0"
			paramsUserMessage["umessage_sender_key"] = strKey
			paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["flag_sent"] = "1"
			subject := "Selamat! Pengkinian Data telah Disetujui"
			body := "Saat ini pengkinian data kamu sudah disetujui. Yuk investasi sekarang juga."
			paramsUserMessage["umessage_subject"] = subject
			paramsUserMessage["umessage_body"] = body
			paramsUserMessage["umessage_category"] = "248"
			paramsUserMessage["flag_archieved"] = "0"
			paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_status"] = "1"
			paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
			paramsUserMessage["rec_created_by"] = strKey

			status, err = models.CreateScUserMessage(paramsUserMessage)
			if err != nil {
				tx.Rollback()
				log.Error("Error create user message")
				return lib.CustomError(status, err.Error(), "failed input data")
			}
			lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body)

			// insert into table ms_customer_bank_account
			paramsCusBankAcc := make(map[string]string)
			strCustomerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
			strBankAccKeyNew := strconv.FormatUint(*oapersonal.BankAccountKey, 10)

			var bankaccountnew models.MsBankAccount
			status, err = models.GetBankAccount(&bankaccountnew, strBankAccKeyNew)
			if err == nil {
				//check bank is ready or no
				bankKeyNew := strconv.FormatUint(bankaccountnew.BankKey, 10)
				var cekBankCus models.CheckBankAccountPengkinianData
				status, err = models.CheckMsBankAccountPengkinianData(&cekBankCus, strCustomerKey, strBankAccKeyNew,
					bankKeyNew, bankaccountnew.AccountNo, bankaccountnew.AccountHolderName)
				if err != nil {
					if err == sql.ErrNoRows {
						paramsCusBankAcc["customer_key"] = strCustomerKey
						paramsCusBankAcc["bank_account_key"] = strBankAccKeyNew
						paramsCusBankAcc["bank_account_name"] = bankaccountnew.AccountHolderName
						paramsCusBankAcc["flag_priority"] = "0"
						paramsCusBankAcc["rec_status"] = "1"
						paramsCusBankAcc["rec_created_date"] = time.Now().Format(dateLayout)
						paramsCusBankAcc["rec_created_by"] = strKey
						status, err = models.CreateMsCustomerBankAccount(paramsCusBankAcc)
						if err != nil {
							tx.Rollback()
							log.Error("Error create ms_customer_bank_account")
							return lib.CustomError(status, err.Error(), "failed input data")
						}
						//end insert into table ms_customer_bank_account

					}
				}
			}

		}
	} else { //reject
		//create user message
		paramsUserMessage := make(map[string]string)
		paramsUserMessage["umessage_type"] = "245"
		strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)
		if oareq.UserLoginKey != nil {
			paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
		} else {
			paramsUserMessage["umessage_recipient_key"] = "0"
		}
		paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["flag_read"] = "0"
		paramsUserMessage["umessage_sender_key"] = strKey
		paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["flag_sent"] = "1"
		var subject string
		body := check2notes + " Silakan menghubungi Customer Service untuk informasi lebih lanjut."
		paramsUserMessage["umessage_body"] = body
		if oareq.CustomerKey == nil { //NEW OA
			paramsUserMessage["umessage_subject"] = "Pembukaan Rekening kamu ditolak"
			subject = "Pembukaan Rekening kamu ditolak"
		} else {
			paramsUserMessage["umessage_subject"] = "Pengkinian Data kamu ditolak"
			subject = "Pengkinian Data kamu ditolak"
		}
		paramsUserMessage["umessage_body"] = body
		paramsUserMessage["umessage_category"] = "248"
		paramsUserMessage["flag_archieved"] = "0"
		paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["rec_status"] = "1"
		paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["rec_created_by"] = strKey

		status, err = models.CreateScUserMessage(paramsUserMessage)
		if err != nil {
			tx.Rollback()
			log.Error("Error create user message")
		}
		lib.CreateNotifCustomerFromAdminByUserLoginKey(strUserLoginKey, subject, body)
		tx.Commit()
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func GetOaRequestListDoTransaction(c echo.Context) error {

	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

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

	items := []string{"oa_request_key", "oa_request_type", "oa_entry_start", "oa_entry_end", "oa_status", "rec_order", "rec_status", "oa_risk_level"}

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
		params["orderBy"] = "oa_request_key"
		params["orderType"] = "DESC"
	}

	params["rec_status"] = "1"

	var oaStatusIn []string
	oaStatusIn = append(oaStatusIn, "260")
	oaStatusIn = append(oaStatusIn, "261")

	var oaRequestDB []models.OaRequest
	status, err = models.GetAllOaRequestDoTransaction(&oaRequestDB, limit, offset, noLimit, params, oaStatusIn, "oa_status")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaRequestDB) < 1 {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var lookupIds []string
	var oaRequestIds []string
	for _, oareq := range oaRequestDB {

		if oareq.Oastatus != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
			}
		}

		if _, ok := lib.Find(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10)); !ok {
			oaRequestIds = append(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10))
		}
	}

	//mapping lookup
	var genLookup []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&genLookup, lookupIds, "lookup_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	gData := make(map[uint64]models.GenLookup)
	for _, gen := range genLookup {
		gData[gen.LookupKey] = gen
	}

	//mapping personal data
	var oaPersonalData []models.OaPersonalData
	if len(oaRequestIds) > 0 {
		status, err = models.GetOaPersonalDataIn(&oaPersonalData, oaRequestIds, "oa_request_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	pdData := make(map[uint64]models.OaPersonalData)
	for _, oaPD := range oaPersonalData {
		pdData[oaPD.OaRequestKey] = oaPD
	}

	var responseData []models.OaRequestListResponse
	for _, oareq := range oaRequestDB {
		var data models.OaRequestListResponse

		if oareq.Oastatus != nil {
			if n, ok := gData[*oareq.Oastatus]; ok {
				data.Oastatus = *n.LkpName
			}
		}

		data.OaRequestKey = oareq.OaRequestKey

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, _ := time.Parse(layout, oareq.OaEntryStart)
		data.OaEntryStart = date.Format(newLayout)
		date, _ = time.Parse(layout, oareq.OaEntryEnd)
		data.OaEntryEnd = date.Format(newLayout)

		if n, ok := pdData[oareq.OaRequestKey]; ok {
			data.EmailAddress = n.EmailAddress
			data.PhoneNumber = n.PhoneMobile
			date, _ = time.Parse(layout, n.DateBirth)
			data.DateBirth = date.Format(newLayout)
			data.FullName = n.FullName
			data.IDCardNo = n.IDcardNo
		}

		responseData = append(responseData, data)
	}

	var countData models.OaRequestCountData
	var pagination int
	if limit > 0 {
		status, err = models.GetCountOaRequestDoTransaction(&countData, params, oaStatusIn, "oa_status")
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

func sendEmailApproveOa(fullName string, email string) {
	// Send email
	t := template.New("email-sukses-verifikasi.html")

	t, err := t.ParseFiles(config.BasePath + "/mail/email-sukses-verifikasi.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl,
		struct {
			Name    string
			FileUrl string
		}{
			Name:    fullName,
			FileUrl: config.FileUrl + "/images/mail"}); err != nil {
		log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "[MNC Duit] Pembukaan Rekening Kamu telah Disetujui")
	mailer.SetBody("text/html", result)
	dialer := gomail.NewDialer(
		config.EmailSMTPHost,
		int(config.EmailSMTPPort),
		config.EmailFrom,
		config.EmailFromPassword,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Error(err)
	}
	log.Info("Email sent")
}
