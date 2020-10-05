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

func GetOaRequestList(c echo.Context) error {
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

	statusStr := c.QueryParam("status")
	if statusStr != "" {
		_, err := strconv.ParseUint(statusStr, 10, 64)
		if err == nil {
			params["oa_status"] = statusStr
		} else {
			log.Error("Status should be number")
			return lib.CustomError(http.StatusBadRequest, "Status should be number", "Status should be number")
		}
	}

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
		if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(oapersonal.Nationality, 10)); !ok {
			personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(oapersonal.Nationality, 10))
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
		if oapersonal.Correspondence != nil {
			if _, ok := lib.Find(personalDataLookupIds, strconv.FormatUint(*oapersonal.Correspondence, 10)); !ok {
				personalDataLookupIds = append(personalDataLookupIds, strconv.FormatUint(*oapersonal.Correspondence, 10))
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
		if n, ok := pData[oapersonal.Nationality]; ok {
			responseData.Nationality = n.LkpName
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
		if oapersonal.Correspondence != nil {
			if n, ok := pData[*oapersonal.Correspondence]; ok {
				responseData.Correspondence = n.LkpName
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

func UpdateStatusApproval(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	oastatus := c.FormValue("oa_status")
	if oastatus == "" {
		log.Error("Missing required parameter: oa_status")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		n, err := strconv.ParseUint(oastatus, 10, 64)
		if err == nil && n > 0 {
			params["oa_status"] = oastatus
		} else {
			log.Error("Wrong input for parameter: oa_status")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_status", "Wrong input for parameter: oa_status")
		}
	}

	oarequestkey := c.FormValue("oa_request_key")
	if oarequestkey == "" {
		log.Error("Missing required parameter: oa_request_key")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		n, err := strconv.ParseUint(oarequestkey, 10, 64)
		if err == nil && n > 0 {
			params["oa_request_key"] = oarequestkey
		} else {
			log.Error("Wrong input for parameter: oa_request_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: oa_request_key", "Wrong input for parameter: oa_request_key")
		}
	}

	var oareq models.OaRequest
	status, err = models.GetOaRequest(&oareq, oarequestkey)
	if err != nil {
		return lib.CustomError(status)
	}

	_, err = models.UpdateOaRequest(params)
	if err != nil {
		log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	log.Info("Success send otp")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
