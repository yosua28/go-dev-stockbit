package controllers

import (
	"api/config"
	"api/db"
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
	errorAuth := initAuthCsKycFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int
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

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13

	log.Println(lib.Profile.RoleKey)

	strOaKey := strconv.FormatUint(*oareq.Oastatus, 10)

	if lib.Profile.RoleKey == roleKeyCs {
		oaStatusCs := strconv.FormatUint(uint64(258), 10)
		if strOaKey != oaStatusCs {
			log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		oaStatusKyc := strconv.FormatUint(uint64(259), 10)
		if strOaKey != oaStatusKyc {
			log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}

	if lib.Profile.RoleKey == roleKeyFundAdmin {
		oaStatusFundAdmin1 := strconv.FormatUint(uint64(260), 10)
		oaStatusFundAdmin2 := strconv.FormatUint(uint64(261), 10)
		if (strOaKey != oaStatusFundAdmin1) && (strOaKey != oaStatusFundAdmin2) {
			log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
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
					responseData.BankAccount.BankKey = bank.BankFullname
				}
				responseData.BankAccount.AccountNo = &bankaccount.AccountNo
				responseData.BankAccount.AccountHolderName = &bankaccount.AccountHolderName
				responseData.BankAccount.BranchName = bankaccount.BranchName

				var lookup models.GenLookup
				strlookup := strconv.FormatUint(bankaccount.BankAccountType, 10)
				status, err = models.GetGenLookup(&lookup, strlookup)
				if err != nil {
					if err != sql.ErrNoRows {
						return lib.CustomError(status)
					}
				} else {
					responseData.BankAccount.BankAccountType = lookup.LkpName
				}
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

	oastatus := c.FormValue("oa_status") //259
	if oastatus == "" {
		log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err := strconv.ParseUint(oastatus, 10, 64)
	if err == nil && n > 0 {
		params["oa_status"] = oastatus
	} else {
		log.Error("Wrong input for parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
	}

	check1notes := c.FormValue("notes")
	params["check1_notes"] = check1notes

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

	dateLayout := "2006-01-02 15:04:05"
	params["check1_date"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["check1_flag"] = "1"
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
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

	oastatus := c.FormValue("oa_status") //260
	if oastatus == "" {
		log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	n, err := strconv.ParseUint(oastatus, 10, 64)
	if err == nil && n > 0 {
		params["oa_status"] = oastatus
	} else {
		log.Error("Wrong input for parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
	}

	check2notes := c.FormValue("notes")
	params["check2_notes"] = check2notes

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

	dateLayout := "2006-01-02 15:04:05"
	params["check2_date"] = time.Now().Format(dateLayout)
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["check2_flag"] = "1"
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
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

	//update oa request
	_, err = models.UpdateOaRequest(params)
	if err != nil {
		tx.Rollback()
		log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}
	log.Info("Success update approved Compliance Transaction")

	//create customer
	var oapersonal models.OaPersonalData
	strKeyOa := strconv.FormatUint(oareq.OaRequestKey, 10)
	status, err = models.GetOaPersonalDataByOaRequestKey(&oapersonal, strKeyOa)
	if err != nil {
		tx.Rollback()
		log.Error("Error Personal Data not Found")
		return lib.CustomError(status, err.Error(), "Personal data not found")
	}

	paramsCustomer := make(map[string]string)
	paramsCustomer["id_customer"] = "0"
	paramsCustomer["unit_holder_idno"] = "202010000001"
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
	if oareq.UserLoginKey != nil {
		strUserLoginKey := strconv.FormatUint(*oareq.UserLoginKey, 10)
		paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
	} else {
		paramsUserMessage["umessage_recipient_key"] = "0"
	}
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sender_key"] = strKey
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	paramsUserMessage["umessage_subject"] = "Pembuatan Akun DUIT"
	paramsUserMessage["umessage_body"] = "Selamat !!! User anda telah disetujui dan sekarang anda telah menjadi Customer."
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

	tx.Commit()

	log.Info("Success create customer")

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
