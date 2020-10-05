package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetOaPersonalDataList(c echo.Context) error {
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

	items := []string{"personal_data_key", "full_name", "gender", "email_address", "rec_order", "rec_status"}

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
	}

	var oaPersonalData []models.OaPersonalData
	status, err = models.GetAllOaPersonalData(&oaPersonalData, limit, offset, params, noLimit)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaPersonalData) < 1 {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var lookupIds []string
	var bankAcountIds []string
	var countryIds []string
	var postalAddressIds []string
	var oaRequestIds []string
	for _, oapersonal := range oaPersonalData {
		//gen lookup
		if oapersonal.AnnualIncome != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.AnnualIncome, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.AnnualIncome, 10))
			}
		}
		if oapersonal.BeneficialRelation != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.BeneficialRelation, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.BeneficialRelation, 10))
			}
		}
		if oapersonal.CorrespondAddress != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.CorrespondAddress, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.CorrespondAddress, 10))
			}
		}
		if oapersonal.Correspondence != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.Correspondence, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.Correspondence, 10))
			}
		}
		if oapersonal.Education != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.Education, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.Education, 10))
			}
		}
		if oapersonal.EmergencyRelation != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.EmergencyRelation, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.EmergencyRelation, 10))
			}
		}
		if oapersonal.Gender != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.Gender, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.Gender, 10))
			}
		}
		if oapersonal.IDcardType != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.IDcardType, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.IDcardType, 10))
			}
		}
		if oapersonal.InvesmentObjectives != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.InvesmentObjectives, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.InvesmentObjectives, 10))
			}
		}
		if oapersonal.MaritalStatus != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.MaritalStatus, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.MaritalStatus, 10))
			}
		}
		if oapersonal.OccupBusinessFields != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.OccupBusinessFields, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.OccupBusinessFields, 10))
			}
		}
		if oapersonal.OccupJob != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.OccupJob, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.OccupJob, 10))
			}
		}
		if oapersonal.OccupPosition != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.OccupPosition, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.OccupPosition, 10))
			}
		}
		if oapersonal.RelationBusinessFields != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.RelationBusinessFields, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.RelationBusinessFields, 10))
			}
		}
		if oapersonal.RelationOccupation != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.RelationOccupation, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.RelationOccupation, 10))
			}
		}
		if oapersonal.RelationType != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.RelationType, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.RelationType, 10))
			}
		}
		if oapersonal.Religion != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.Religion, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.Religion, 10))
			}
		}
		if oapersonal.SourceofFund != nil {
			if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oapersonal.SourceofFund, 10)); !ok {
				lookupIds = append(lookupIds, strconv.FormatUint(*oapersonal.SourceofFund, 10))
			}
		}

		//bank account
		if oapersonal.BankAccountKey != nil {
			if _, ok := lib.Find(bankAcountIds, strconv.FormatUint(*oapersonal.BankAccountKey, 10)); !ok {
				bankAcountIds = append(bankAcountIds, strconv.FormatUint(*oapersonal.BankAccountKey, 10))
			}
		}

		//country
		if _, ok := lib.Find(countryIds, strconv.FormatUint(oapersonal.Nationality, 10)); !ok {
			countryIds = append(countryIds, strconv.FormatUint(oapersonal.Nationality, 10))
		}

		//postal address
		if oapersonal.DomicileAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.DomicileAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.DomicileAddressKey, 10))
			}
		}
		if oapersonal.IDcardAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.IDcardAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.IDcardAddressKey, 10))
			}
		}
		if oapersonal.OccupAddressKey != nil {
			if _, ok := lib.Find(postalAddressIds, strconv.FormatUint(*oapersonal.OccupAddressKey, 10)); !ok {
				postalAddressIds = append(postalAddressIds, strconv.FormatUint(*oapersonal.OccupAddressKey, 10))
			}
		}

		//oa request
		if _, ok := lib.Find(oaRequestIds, strconv.FormatUint(oapersonal.OaRequestKey, 10)); !ok {
			oaRequestIds = append(oaRequestIds, strconv.FormatUint(oapersonal.OaRequestKey, 10))
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

	//mapping bank account
	var msBankAccount []models.MsBankAccount
	if len(bankAcountIds) > 0 {
		status, err = models.GetMsBankAccountIn(&msBankAccount, bankAcountIds, "bank_account_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	bankData := make(map[uint64]models.MsBankAccount)
	for _, bank := range msBankAccount {
		bankData[bank.BankAccountKey] = bank
	}

	//mapping country
	var msCountry []models.MsCountry
	if len(bankAcountIds) > 0 {
		status, err = models.GetMsCountryIn(&msCountry, bankAcountIds, "country_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	countryData := make(map[uint64]models.MsCountry)
	for _, coun := range msCountry {
		countryData[coun.CountryKey] = coun
	}

	//mapping postal address
	var oaPostalAddress []models.OaPostalAddress
	if len(bankAcountIds) > 0 {
		status, err = models.GetOaPostalAddressIn(&oaPostalAddress, postalAddressIds, "postal_address_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	postalAddressData := make(map[uint64]models.OaPostalAddress)
	for _, pa := range oaPostalAddress {
		postalAddressData[pa.PostalAddressKey] = pa
	}

	//mapping oa request
	var oaRequest []models.OaRequest
	if len(bankAcountIds) > 0 {
		status, err = models.GetOaRequestsIn(&oaRequest, oaRequestIds, "oa_request_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	oaRequestData := make(map[uint64]models.OaRequest)
	for _, oareq := range oaRequest {
		oaRequestData[oareq.OaRequestKey] = oareq
	}

	var responseData []models.OaPersonalDataResponse

	for _, oaPD := range oaPersonalData {
		var data models.OaPersonalDataResponse

		// data.OaRequestKey = oaPD.OaRequestKey
		// data.Nationality = oaPD.Nationality
		// data.IDcardType = oaPD.IDcardType
		// data.Gender = oaPD.Gender
		// data.MaritalStatus = oaPD.MaritalStatus
		// data.IDcardAddressKey = oaPD.IDcardAddressKey
		// data.DomicileAddressKey = oaPD.DomicileAddressKey
		// data.CorrespondAddress = oaPD.CorrespondAddress
		// data.Religion = oaPD.Religion
		// data.Education = oaPD.Education
		// data.AnnualIncome = oaPD.AnnualIncome
		// data.BeneficialRelation = oaPD.BeneficialRelation
		// data.OccupAddressKey = oaPD.OccupAddressKey
		// data.Correspondence = oaPD.Correspondence
		// data.EmergencyRelation = oaPD.EmergencyRelation
		// data.InvesmentObjectives = oaPD.InvesmentObjectives
		// data.OccupBusinessFields = oaPD.OccupBusinessFields
		// data.OccupJob = oaPD.OccupJob
		// data.OccupPosition = oaPD.OccupPosition
		// data.RelationBusinessFields = oaPD.RelationBusinessFields
		// data.RelationOccupation = oaPD.RelationOccupation
		// data.RelationType = oaPD.RelationType
		// data.SourceofFund = oaPD.SourceofFund
		// data.BankAccountKey = oaPD.BankAccountKey

		data.PersonalDataKey = oaPD.PersonalDataKey
		data.FullName = oaPD.FullName
		data.PlaceBirth = oaPD.PlaceBirth
		data.DateBirth = oaPD.DateBirth
		data.IDcardNo = oaPD.IDcardNo
		data.IDcardExpiredDate = oaPD.IDcardExpiredDate
		data.IDcardNeverExpired = oaPD.IDcardNeverExpired
		data.PhoneHome = oaPD.PhoneHome
		data.PhoneMobile = oaPD.PhoneMobile
		data.EmailAddress = oaPD.EmailAddress
		data.PicSelfie = oaPD.PicSelfie
		data.PicKtp = oaPD.PicKtp
		data.PicSelfieKtp = oaPD.PicSelfieKtp
		data.GeolocName = oaPD.GeolocName
		data.GeolocLongitude = oaPD.GeolocLongitude
		data.GeolocLatitude = oaPD.GeolocLatitude
		data.OccupCompany = oaPD.OccupCompany
		data.OccupPhone = oaPD.OccupPhone
		data.OccupWebURL = oaPD.OccupWebUrl
		data.RelationFullName = oaPD.RelationFullName
		data.MotherMaidenName = oaPD.MotherMaidenName
		data.EmergencyFullName = oaPD.EmergencyFullName
		data.EmergencyPhoneNo = oaPD.EmergencyPhoneNo
		data.BeneficialFullName = oaPD.BeneficialFullName
		data.RecOrder = oaPD.RecOrder
		data.RecStatus = oaPD.RecStatus

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

// func GetOaPersonalData(c echo.Context) error {
// 	var err error
// 	var status int
// 	//Get parameter limit
// 	keyStr := c.Param("key")
// 	key, _ := strconv.ParseUint(keyStr, 10, 64)
// 	if key == 0 {
// 		return lib.CustomError(http.StatusNotFound)
// 	}

// 	var oareq models.OaRequest
// 	status, err = models.GetOaRequest(&oareq, keyStr)
// 	if err != nil {
// 		return lib.CustomError(status)
// 	}

// 	var lookupIds []string

// 	if oareq.OaRequestType != nil {
// 		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.OaRequestType, 10)); !ok {
// 			lookupIds = append(lookupIds, strconv.FormatUint(*oareq.OaRequestType, 10))
// 		}
// 	}
// 	if oareq.OaRiskLevel != nil {
// 		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10)); !ok {
// 			lookupIds = append(lookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10))
// 		}
// 	}

// 	if oareq.Oastatus != nil {
// 		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10)); !ok {
// 			lookupIds = append(lookupIds, strconv.FormatUint(*oareq.Oastatus, 10))
// 		}
// 	}

// 	//mapping lookup
// 	var genLookup []models.GenLookup
// 	if len(lookupIds) > 0 {
// 		status, err = models.GetGenLookupIn(&genLookup, lookupIds, "lookup_key")
// 		if err != nil {
// 			log.Error(err.Error())
// 			return lib.CustomError(status, err.Error(), "Failed get data")
// 		}
// 	}
// 	gData := make(map[uint64]models.GenLookup)
// 	for _, gen := range genLookup {
// 		gData[gen.LookupKey] = gen
// 	}

// 	//maping response
// 	var responseData models.OaRequestDataResponse

// 	if oareq.OaRequestType != nil {
// 		if n, ok := gData[*oareq.OaRequestType]; ok {
// 			responseData.OaRequestType = n.LkpText1
// 		}
// 	}

// 	if oareq.OaRiskLevel != nil {
// 		if n, ok := gData[*oareq.OaRiskLevel]; ok {
// 			responseData.OaRiskLevel = n.LkpText1
// 		}
// 	}

// 	if oareq.Oastatus != nil {
// 		if n, ok := gData[*oareq.Oastatus]; ok {
// 			responseData.Oastatus = *n.LkpText1
// 		}
// 	}

// 	if oareq.CustomerKey != nil {
// 		var msCustomer models.MsCustomer
// 		var t uint64 = *oareq.CustomerKey
// 		str := strconv.FormatUint(t, 10)
// 		status, err = models.GetMsCustomer(&msCustomer, str)
// 		if err != nil {
// 			return lib.CustomError(status)
// 		}
// 		responseData.Customer = &msCustomer.FullName
// 	}

// 	if oareq.UserLoginKey != nil {
// 		var scUserLogin models.ScUserLogin
// 		var t uint64 = *oareq.UserLoginKey
// 		str := strconv.FormatUint(t, 10)
// 		status, err = models.GetScUserLoginByKey(&scUserLogin, str)
// 		if err != nil {
// 			return lib.CustomError(status)
// 		}
// 		responseData.UserLoginName = &scUserLogin.UloginName
// 		responseData.UserLoginFullName = &scUserLogin.UloginFullName
// 	}

// 	responseData.OaRequestKey = oareq.OaRequestKey
// 	responseData.OaEntryStart = oareq.OaEntryStart
// 	responseData.OaEntryEnd = oareq.OaEntryEnd
// 	responseData.Check1Date = oareq.Check1Date
// 	responseData.Check1Flag = oareq.Check1Flag
// 	responseData.Check1References = oareq.Check1References
// 	responseData.Check1Notes = oareq.Check1Notes
// 	responseData.Check2Date = oareq.Check2Date
// 	responseData.Check2Flag = oareq.Check2Flag
// 	responseData.Check2References = oareq.Check2References
// 	responseData.Check2Notes = oareq.Check2Notes
// 	responseData.RecOrder = oareq.RecOrder
// 	responseData.RecStatus = oareq.RecStatus

// 	var response lib.Response
// 	response.Status.Code = http.StatusOK
// 	response.Status.MessageServer = "OK"
// 	response.Status.MessageClient = "OK"
// 	response.Data = responseData

// 	return c.JSON(http.StatusOK, response)
// }
