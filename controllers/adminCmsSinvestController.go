package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func DownloadOaRequestFormatSinvest(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int
	var offset uint64
	// var limit uint64

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

	var oaRequestDB []models.OaRequest
	status, err = models.GetAllOaRequestDoTransaction(&oaRequestDB, config.LimitQuery, offset, true, params, oaStatusIn, "oa_status")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(oaRequestDB) < 1 {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var oaRequestLookupIds []string
	var oaRequestIds []string
	for _, oareq := range oaRequestDB {
		if _, ok := lib.Find(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10)); !ok {
			oaRequestIds = append(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10))
		}

		if oareq.OaRiskLevel != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10))
			}
		}
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
	var postalAddressIds []string
	var nasionalityIds []string

	for _, oapersonal := range oaPersonalData {
		pdData[oapersonal.OaRequestKey] = oapersonal

		if _, ok := lib.Find(nasionalityIds, strconv.FormatUint(oapersonal.Nationality, 10)); !ok {
			nasionalityIds = append(nasionalityIds, strconv.FormatUint(oapersonal.Nationality, 10))
		}

		if oapersonal.Gender != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.Gender, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.Gender, 10))
			}
		}
		if oapersonal.MaritalStatus != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.MaritalStatus, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.MaritalStatus, 10))
			}
		}
		if oapersonal.Religion != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.Religion, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.Religion, 10))
			}
		}
		if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(oapersonal.Nationality, 10)); !ok {
			oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(oapersonal.Nationality, 10))
		}
		if oapersonal.Education != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.Education, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.Education, 10))
			}
		}
		if oapersonal.OccupJob != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.OccupJob, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.OccupJob, 10))
			}
		}
		if oapersonal.AnnualIncome != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.AnnualIncome, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.AnnualIncome, 10))
			}
		}
		if oapersonal.SourceofFund != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.SourceofFund, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.SourceofFund, 10))
			}
		}
		if oapersonal.InvesmentObjectives != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.InvesmentObjectives, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.InvesmentObjectives, 10))
			}
		}
		if oapersonal.RelationType != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oapersonal.RelationType, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oapersonal.RelationType, 10))
			}
		}

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

	//gen lookup
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

	//postal data
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
	var oaCityIds []string
	for _, posAdd := range oaPstalAddressList {
		postalData[posAdd.PostalAddressKey] = posAdd

		if posAdd.KabupatenKey != nil {
			if _, ok := lib.Find(oaCityIds, strconv.FormatUint(*posAdd.KabupatenKey, 10)); !ok {
				oaCityIds = append(oaCityIds, strconv.FormatUint(*posAdd.KabupatenKey, 10))
			}
		}
		if posAdd.KecamatanKey != nil {
			if _, ok := lib.Find(oaCityIds, strconv.FormatUint(*posAdd.KecamatanKey, 10)); !ok {
				oaCityIds = append(oaCityIds, strconv.FormatUint(*posAdd.KecamatanKey, 10))
			}
		}
	}

	var cityList []models.MsCity
	if len(oaCityIds) > 0 {
		status, err = models.GetMsCityIn(&cityList, oaCityIds, "city_key")
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

		if _, ok := lib.Find(nasionalityIds, strconv.FormatUint(city.CountryKey, 10)); !ok {
			nasionalityIds = append(nasionalityIds, strconv.FormatUint(city.CountryKey, 10))
		}
	}

	var countryList []models.MsCountry
	status, err = models.GetMsCountryIn(&countryList, nasionalityIds, "country_key")
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	countryData := make(map[uint64]models.MsCountry)
	for _, country := range countryList {
		countryData[country.CountryKey] = country
	}

	var responseData []models.OaRequestCsvFormat

	for _, oareq := range oaRequestDB {
		if n, ok := pdData[oareq.OaRequestKey]; ok {
			var data models.OaRequestCsvFormat
			data.Type = "1"
			strKey := strconv.FormatUint(oareq.OaRequestKey, 10)
			data.SACode = strKey
			data.SID = ""
			data.FirstName = n.FullName
			data.MiddleName = ""
			data.LastName = ""

			data.CountryOfNationality = ""
			data.CountryOfBirth = ""
			if n, ok := countryData[n.Nationality]; ok {
				data.CountryOfNationality = n.CouCode
				data.CountryOfBirth = n.CouCode
			}

			data.IDno = n.IDcardNo

			data.IDExpirationDate = ""
			if n.IDcardExpiredDate != nil {
				exCard := *n.IDcardExpiredDate
				data.IDExpirationDate = exCard
			}

			data.NpwpNo = ""
			data.NpwpRegistrationDate = ""
			data.PlaceOfBirth = n.PlaceBirth

			layout := "2006-01-02 15:04:05"
			newLayout := "20060102"
			date, _ := time.Parse(layout, n.DateBirth)
			data.DateOfBirth = date.Format(newLayout)

			data.Gender = ""
			if n.Gender != nil {
				if n, ok := gData[*n.Gender]; ok {
					gen := *n.LkpCode
					data.Gender = gen
				}
			}

			data.EducationalBackground = ""
			if n.Education != nil {
				if n, ok := gData[*n.Education]; ok {
					edu := *n.LkpCode
					data.EducationalBackground = edu
				}
			}

			data.MotherMaidenName = n.MotherMaidenName

			data.Religion = ""
			if n.Religion != nil {
				if n, ok := gData[*n.Religion]; ok {
					rel := *n.LkpCode
					data.Religion = rel
				}
			}

			data.Occupation = ""
			if n.OccupJob != nil {
				if n, ok := gData[*n.OccupJob]; ok {
					occ := *n.LkpCode
					data.Occupation = occ
				}
			}

			data.IncomeLevel = ""
			if n.AnnualIncome != nil {
				if n, ok := gData[*n.AnnualIncome]; ok {
					income := *n.LkpCode
					data.IncomeLevel = income
				}
			}

			data.MaritalStatus = ""
			if n.MaritalStatus != nil {
				if n, ok := gData[*n.MaritalStatus]; ok {
					marital := *n.LkpCode
					data.MaritalStatus = marital
				}
			}

			data.SpouseName = ""
			if n.RelationType != nil {
				strRelation := strconv.FormatUint(*n.RelationType, 10)
				lookupRelationPasangan := strconv.FormatUint(uint64(87), 10) // suami/istri
				if strRelation == lookupRelationPasangan {
					spouseName := *n.RelationFullName
					data.SpouseName = spouseName
				}
			}

			data.InvestorRiskProfile = ""
			if oareq.OaRiskLevel != nil {
				if g, ok := gData[*oareq.OaRiskLevel]; ok {
					irp := *g.LkpCode
					data.InvestorRiskProfile = irp
				}
			}

			data.InvestmentObjective = ""
			if n.InvesmentObjectives != nil {
				if g, ok := gData[*n.InvesmentObjectives]; ok {
					io := *g.LkpCode
					data.InvestmentObjective = io
				}
			}

			data.SourceOfFund = ""
			if n.SourceofFund != nil {
				if g, ok := gData[*n.SourceofFund]; ok {
					sof := *g.LkpCode
					data.SourceOfFund = sof
				}
			}

			data.AssetOwner = "1"

			data.KTPAddress = ""
			data.KTPCityCode = ""
			data.KTPPostalCode = ""
			data.CorrespondenceAddress = ""
			data.CorrespondenceCityCode = ""
			data.CorrespondenceCityName = ""
			data.CountryOfCorrespondence = ""
			if n.IDcardAddressKey != nil {
				if a, ok := postalData[*n.IDcardAddressKey]; ok {
					//set alamat KTP
					if a.AddressLine1 != nil {
						ktpAddress := *a.AddressLine1
						data.KTPAddress = ktpAddress
					}
					if a.KabupatenKey != nil {
						if c, ok := cityData[*a.KabupatenKey]; ok {
							data.KTPCityCode = c.CityCode
						}
					}

					if a.PostalCode != nil {
						postal := *a.PostalCode
						data.KTPPostalCode = postal
					}

					//set Correspondence Address untuk sementara samain KTP dulu
					if a.AddressLine1 != nil {
						ktpAddress := *a.AddressLine1
						data.CorrespondenceAddress = ktpAddress
					}

					if a.KabupatenKey != nil {
						if c, ok := cityData[*a.KabupatenKey]; ok {
							if co, ok := countryData[c.CountryKey]; ok {
								data.CountryOfCorrespondence = co.CouCode
							}
							data.CorrespondenceCityCode = c.CityCode
							data.CorrespondenceCityName = c.CityName
						}
					}
					if a.PostalCode != nil {
						postal := *a.PostalCode
						data.CorrespondencePostalCode = postal
					}
				}
			}

			if n.DomicileAddressKey != nil {
				if a, ok := postalData[*n.DomicileAddressKey]; ok {
					//set alamat KTP
					if a.AddressLine1 != nil {
						domiAddress := *a.AddressLine1
						data.DomicileAddress = domiAddress
					}

					if a.KabupatenKey != nil {
						if c, ok := cityData[*a.KabupatenKey]; ok {
							if co, ok := countryData[c.CountryKey]; ok {
								data.CountryOfDomicile = co.CouCode
							}
							data.DomicileCityCode = c.CityCode
							data.DomicileCityName = c.CityName
						}
					}

					if a.PostalCode != nil {
						postal := *a.PostalCode
						data.DomicilePostalCode = postal
					}

				}
			}

			homePhone := n.PhoneHome
			data.HomePhone = homePhone

			mobilePhone := n.PhoneMobile
			data.MobilePhone = mobilePhone
			data.Facsimile = ""
			data.Email = n.EmailAddress

			data.StatementType = "2"
			data.FATCA = ""
			data.ForeignTIN = ""
			data.ForeignTINIssuanceCountry = ""
			data.REDMPaymentBankBICCode1 = ""
			data.REDMPaymentBankBIMemberCode1 = ""
			data.REDMPaymentBankName1 = ""
			data.REDMPaymentBankCountry1 = ""
			data.REDMPaymentBankBranch1 = ""
			data.REDMPaymentACCcy1 = "IDR"
			data.REDMPaymentACNo1 = ""
			data.REDMPaymentACName1 = ""
			data.REDMPaymentBankBICCode2 = ""
			data.REDMPaymentBankBIMemberCode2 = ""
			data.REDMPaymentBankName2 = ""
			data.REDMPaymentBankCountry2 = ""
			data.REDMPaymentBankBranch2 = ""
			data.REDMPaymentACCcy2 = "IDR"
			data.REDMPaymentACNo2 = ""
			data.REDMPaymentACName2 = ""
			data.REDMPaymentBankBICCode3 = ""
			data.REDMPaymentBankBIMemberCode3 = ""
			data.REDMPaymentBankName3 = ""
			data.REDMPaymentBankCountry3 = ""
			data.REDMPaymentBankBranch3 = ""
			data.REDMPaymentACCcy3 = "IDR"
			data.REDMPaymentACNo3 = ""
			data.REDMPaymentACName3 = ""
			data.ClientCode = ""

			responseData = append(responseData, data)
		}
	}

	if len(responseData) > 0 {
		paramsUpdate := make(map[string]string)

		strOaStatus := strconv.FormatUint(261, 10) //customer builf, proses upload data to Sinvest
		paramsUpdate["oa_status"] = strOaStatus
		dateLayout := "2006-01-02 15:04:05"
		paramsUpdate["rec_modified_date"] = time.Now().Format(dateLayout)
		strKey := strconv.FormatUint(lib.Profile.UserID, 10)
		paramsUpdate["rec_modified_by"] = strKey

		_, err = models.UpdateOaRequestByKeyIn(paramsUpdate, oaRequestIds, "oa_request_key")
		if err != nil {
			log.Error("Error update oa request")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
