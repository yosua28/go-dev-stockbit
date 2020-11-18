package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"bufio"
	"database/sql"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	var customerIds []string
	for _, oareq := range oaRequestDB {
		if _, ok := lib.Find(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10)); !ok {
			oaRequestIds = append(oaRequestIds, strconv.FormatUint(oareq.OaRequestKey, 10))
		}

		if oareq.OaRiskLevel != nil {
			if _, ok := lib.Find(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10)); !ok {
				oaRequestLookupIds = append(oaRequestLookupIds, strconv.FormatUint(*oareq.OaRiskLevel, 10))
			}
		}

		if oareq.CustomerKey != nil {
			if _, ok := lib.Find(customerIds, strconv.FormatUint(*oareq.CustomerKey, 10)); !ok {
				customerIds = append(customerIds, strconv.FormatUint(*oareq.CustomerKey, 10))
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
	var bankAccountIds []string

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

		if oapersonal.BankAccountKey != nil {
			if _, ok := lib.Find(bankAccountIds, strconv.FormatUint(*oapersonal.BankAccountKey, 10)); !ok {
				bankAccountIds = append(bankAccountIds, strconv.FormatUint(*oapersonal.BankAccountKey, 10))
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

	//customer
	var customer []models.MsCustomer
	if len(customerIds) > 0 {
		status, err = models.GetMsCustomerIn(&customer, customerIds, "customer_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	customerData := make(map[uint64]models.MsCustomer)
	for _, cus := range customer {
		customerData[cus.CustomerKey] = cus
	}

	//bank account
	var bankAccount []models.MsBankAccount
	if len(bankAccountIds) > 0 {
		status, err = models.GetMsBankAccountIn(&bankAccount, bankAccountIds, "bank_account_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	bankData := make(map[uint64]models.MsBankAccount)
	var bankIds []string
	var currencyIds []string
	for _, b := range bankAccount {
		bankData[b.BankAccountKey] = b

		if _, ok := lib.Find(bankIds, strconv.FormatUint(b.BankKey, 10)); !ok {
			bankIds = append(bankIds, strconv.FormatUint(b.BankKey, 10))
		}

		if _, ok := lib.Find(currencyIds, strconv.FormatUint(b.CurrencyKey, 10)); !ok {
			currencyIds = append(currencyIds, strconv.FormatUint(b.CurrencyKey, 10))
		}
	}

	//ms bank
	var msBank []models.MsBank
	if len(bankIds) > 0 {
		status, err = models.GetMsBankIn(&msBank, bankIds, "bank_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	msBankData := make(map[uint64]models.MsBank)
	for _, b := range msBank {
		msBankData[b.BankKey] = b
	}

	//ms currency
	var msCurrency []models.MsCurrency
	if len(currencyIds) > 0 {
		status, err = models.GetMsCurrencyIn(&msCurrency, currencyIds, "currency_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	currencyData := make(map[uint64]models.MsCurrency)
	for _, c := range msCurrency {
		currencyData[c.CurrencyKey] = c
	}

	var responseData []models.OaRequestCsvFormatFiksTxt

	var scApp models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&scApp, "LAST_CLIENT_CODE")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data LAST_CLIENT_CODE")
	}

	log.Println("last = " + *scApp.AppConfigValue)

	last, _ := strconv.ParseUint(*scApp.AppConfigValue, 10, 64)
	if last == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var lastValue string

	txtHeader := "Type|SA Code|SID|First Name|Middle Name| Last Name|Country of Nationality|ID no |ID Expiration Date| NPWP no|NPWP Registration Date|Country of Birth|Place of Birth|Date of Birth |Gender|Educational Background|Mothers Maiden Name|Religion|Occupation|Income Level (IDR)|Marital Status|Spouses Name|Investors Risk Profile|Investment Objective|Source of Fund|Asset Owner|KTP Address|KTP City Code|KTP Postal Code|Correspondence Address|Correspondence City Code|Correspondence City Name|Correspondence Postal Code|Country of Correspondence|Domicile Address|Domicile City Code|Domicile City Name|Domicile Postal Code|Country of Domicile|Home Phone|Mobile Phone|Facsimile|Email|Statement Type|FATCA (Status)|TIN / Foreign TIN|TIN / Foreign TIN Issuance Country|REDM Payment Bank BIC Code 1|REDM Payment Bank BI Member Code 1|REDM Payment Bank Name 1| REDM Payment Bank Country 1| REDM Payment Bank Branch 1|REDM Payment A/C CCY 1|REDM Payment A/C No. 1|REDM Payment A/C Name 1|REDM Payment Bank BIC Code 2|REDM Payment Bank BI Member Code 2|REDM Payment Bank Name 2|REDM Payment Bank Country 2|REDM Payment Bank Branch 2|REDM Payment A/C CCY 2|REDM Payment A/C No. 2|REDM Payment A/C Name 2|REDM Payment Bank BIC Code 3|REDM Payment Bank BI Member Code 3| REDM Payment Bank Name 3|REDM Payment Bank Country 3|REDM Payment Bank Branch 3|REDM Payment A/C CCY 3|REDM Payment A/C No. 3|REDM Payment A/C Name 3|Client Code"
	var dataRow models.OaRequestCsvFormatFiksTxt
	dataRow.DataRow = txtHeader
	responseData = append(responseData, dataRow)

	for _, oareq := range oaRequestDB {
		if n, ok := pdData[oareq.OaRequestKey]; ok {
			var data models.OaRequestCsvFormat

			strType := strconv.FormatUint(*oareq.OaRequestType, 10)
			if strType == "127" {
				data.Type = "1"
			} else {
				data.Type = "2"
			}
			data.SACode = "EP002"
			data.SID = ""
			data.FirstName = ""
			data.MiddleName = ""
			data.LastName = ""
			if oareq.CustomerKey != nil {
				if c, ok := customerData[*oareq.CustomerKey]; ok {
					if c.SidNo != nil {
						data.SID = *c.SidNo
					}
					if c.FirstName != nil {
						data.FirstName = *c.FirstName
					}
					if c.MiddleName != nil {
						data.MiddleName = *c.MiddleName
					}
					if c.LastName != nil {
						data.LastName = *c.LastName
					}
				}
			}

			data.CountryOfNationality = ""
			data.CountryOfBirth = ""
			if co, ok := countryData[n.Nationality]; ok {
				strCountry := strconv.FormatUint(co.CountryKey, 10)
				if strCountry == "97" { // indonesia WNI
					data.CountryOfNationality = "1"
				} else {
					data.CountryOfNationality = "2"
				}
				data.CountryOfBirth = co.CouCode
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
				}
			}

			data.CorrespondenceAddress = ""
			data.CorrespondenceCityCode = ""
			data.CorrespondenceCityName = ""
			data.CorrespondencePostalCode = ""
			data.CountryOfCorrespondence = ""
			if n.DomicileAddressKey != nil {
				if a, ok := postalData[*n.DomicileAddressKey]; ok {
					//set alamat KTP
					if a.AddressLine1 != nil {
						domiAddress := *a.AddressLine1
						data.DomicileAddress = domiAddress
						data.CorrespondenceAddress = domiAddress
					}

					if a.KabupatenKey != nil {
						if c, ok := cityData[*a.KabupatenKey]; ok {
							if co, ok := countryData[c.CountryKey]; ok {
								data.CountryOfDomicile = co.CouCode
								data.CountryOfCorrespondence = co.CouCode
							}
							data.DomicileCityCode = c.CityCode
							data.DomicileCityName = c.CityName
							data.CorrespondenceCityCode = c.CityCode
							data.CorrespondenceCityName = c.CityName
						}
					}

					if a.PostalCode != nil {
						postal := *a.PostalCode
						data.DomicilePostalCode = postal
						data.CorrespondencePostalCode = postal
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
			data.REDMPaymentACCcy1 = ""
			data.REDMPaymentACNo1 = ""
			data.REDMPaymentACName1 = ""

			if bankAccount, ok := bankData[*n.BankAccountKey]; ok {
				if msBank, ok := msBankData[bankAccount.BankKey]; ok {
					data.REDMPaymentBankBIMemberCode1 = msBank.BankCode
					data.REDMPaymentBankName1 = msBank.BankName
				}
				data.REDMPaymentBankCountry1 = "ID"

				if cur, ok := currencyData[bankAccount.CurrencyKey]; ok {
					data.REDMPaymentACCcy1 = cur.Code
				}
				data.REDMPaymentACNo1 = bankAccount.AccountNo
				data.REDMPaymentACName1 = bankAccount.AccountHolderName
			}

			data.REDMPaymentBankBICCode2 = ""
			data.REDMPaymentBankBIMemberCode2 = ""
			data.REDMPaymentBankName2 = ""
			data.REDMPaymentBankCountry2 = ""
			data.REDMPaymentBankBranch2 = ""
			data.REDMPaymentACCcy2 = ""
			data.REDMPaymentACNo2 = ""
			data.REDMPaymentACName2 = ""
			data.REDMPaymentBankBICCode3 = ""
			data.REDMPaymentBankBIMemberCode3 = ""
			data.REDMPaymentBankName3 = ""
			data.REDMPaymentBankCountry3 = ""
			data.REDMPaymentBankBranch3 = ""
			data.REDMPaymentACCcy3 = ""
			data.REDMPaymentACNo3 = ""
			data.REDMPaymentACName3 = ""
			data.ClientCode = ""

			//start update client_code if new customer
			if strType == "127" { //type NEW
				last = last + 1
				paramsCustomer := make(map[string]string)
				var convLast string
				convLast = strconv.FormatUint(uint64(last), 10)
				clientCode := lib.PadLeft(convLast, "0", 6)
				paramsCustomer["client_code"] = clientCode
				dateLayout := "2006-01-02 15:04:05"
				paramsCustomer["rec_modified_date"] = time.Now().Format(dateLayout)
				strKeyLogin := strconv.FormatUint(lib.Profile.UserID, 10)
				paramsCustomer["rec_modified_by"] = strKeyLogin
				strCustomerKey := strconv.FormatUint(*oareq.CustomerKey, 10)
				paramsCustomer["customer_key"] = strCustomerKey
				_, err = models.UpdateMsCustomer(paramsCustomer)
				if err != nil {
					log.Error("Error update oa request")
					return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
				}
				lastValue = paramsCustomer["client_code"]

				data.ClientCode = paramsCustomer["client_code"]
			}
			//end

			txtData := data.Type + "|" +
				data.SACode + "|" +
				data.SID + "|" +
				data.FirstName + "|" +
				data.MiddleName + "|" +
				data.LastName + "|" +
				data.CountryOfNationality + "|" +
				data.IDno + "|" +
				data.IDExpirationDate + "|" +
				data.NpwpNo + "|" +
				data.NpwpRegistrationDate + "|" +
				data.CountryOfBirth + "|" +
				data.PlaceOfBirth + "|" +
				data.DateOfBirth + "|" +
				data.Gender + "|" +
				data.EducationalBackground + "|" +
				data.MotherMaidenName + "|" +
				data.Religion + "|" +
				data.Occupation + "|" +
				data.IncomeLevel + "|" +
				data.MaritalStatus + "|" +
				data.SpouseName + "|" +
				data.InvestorRiskProfile + "|" +
				data.InvestmentObjective + "|" +
				data.SourceOfFund + "|" +
				data.AssetOwner + "|" +
				data.KTPAddress + "|" +
				data.KTPCityCode + "|" +
				data.KTPPostalCode + "|" +
				data.CorrespondenceAddress + "|" +
				data.CorrespondenceCityCode + "|" +
				data.CorrespondenceCityName + "|" +
				data.CorrespondencePostalCode + "|" +
				data.CountryOfCorrespondence + "|" +
				data.DomicileAddress + "|" +
				data.DomicileCityCode + "|" +
				data.DomicileCityName + "|" +
				data.DomicilePostalCode + "|" +
				data.CountryOfDomicile + "|" +
				data.HomePhone + "|" +
				data.MobilePhone + "|" +
				data.Facsimile + "|" +
				data.Email + "|" +
				data.StatementType + "|" +
				data.FATCA + "|" +
				data.ForeignTIN + "|" +
				data.ForeignTINIssuanceCountry + "|" +
				data.REDMPaymentBankBICCode1 + "|" +
				data.REDMPaymentBankBIMemberCode1 + "|" +
				data.REDMPaymentBankName1 + "|" +
				data.REDMPaymentBankCountry1 + "|" +
				data.REDMPaymentBankBranch1 + "|" +
				data.REDMPaymentACCcy1 + "|" +
				data.REDMPaymentACNo1 + "|" +
				data.REDMPaymentACName1 + "|" +
				data.REDMPaymentBankBICCode2 + "|" +
				data.REDMPaymentBankBIMemberCode2 + "|" +
				data.REDMPaymentBankName2 + "|" +
				data.REDMPaymentBankCountry2 + "|" +
				data.REDMPaymentBankBranch2 + "|" +
				data.REDMPaymentACCcy2 + "|" +
				data.REDMPaymentACNo2 + "|" +
				data.REDMPaymentACName2 + "|" +
				data.REDMPaymentBankBICCode3 + "|" +
				data.REDMPaymentBankBIMemberCode3 + "|" +
				data.REDMPaymentBankName3 + "|" +
				data.REDMPaymentBankCountry3 + "|" +
				data.REDMPaymentBankBranch3 + "|" +
				data.REDMPaymentACCcy3 + "|" +
				data.REDMPaymentACNo3 + "|" +
				data.REDMPaymentACName3 + "|" +
				data.ClientCode

			var txt models.OaRequestCsvFormatFiksTxt
			txt.DataRow = txtData

			responseData = append(responseData, txt)
		}
	}

	//value awal = 009995 ----------------- update app_config
	if lastValue != "" {
		paramsConfig := make(map[string]string)
		paramsConfig["app_config_value"] = lastValue
		dateLayout := "2006-01-02 15:04:05"
		paramsConfig["rec_modified_date"] = time.Now().Format(dateLayout)
		strKeyLogin := strconv.FormatUint(lib.Profile.UserID, 10)
		paramsConfig["rec_modified_by"] = strKeyLogin
		_, err = models.UpdateMsCustomerByConfigCode(paramsConfig, "LAST_CLIENT_CODE")
		if err != nil {
			log.Error("Error update App Config")
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
		}
	}
	//end

	if len(responseData) > 0 {
		paramsUpdate := make(map[string]string)

		strOaStatus := strconv.FormatUint(261, 10) //customer build, proses upload data to Sinvest
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

func UploadOaRequestFormatSinvest(c echo.Context) error {

	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error

	err = os.MkdirAll(config.BasePath+"/oa_upload/sinvest/", 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("file")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			log.Println(extension)
			roles := []string{".txt", ".TXT"}
			_, found := lib.Find(roles, extension)
			if !found {
				return lib.CustomError(http.StatusUnauthorized, "Format file must .txt", "Format file must .txt")
			}
			// Generate filename
			var filename string
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename+extension)
			// Upload txt and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/oa_upload/sinvest/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}

			fileTxt, err := os.Open(config.BasePath + "/oa_upload/sinvest/" + filename + extension)

			if err != nil {
				log.Println("failed to open txt")
				log.Println(err)
				// log.Fatalf("failed to open")
			}

			scanner := bufio.NewScanner(fileTxt)

			scanner.Split(bufio.ScanLines)
			var text []string

			for scanner.Scan() {
				text = append(text, scanner.Text())
			}

			fileTxt.Close()

			dateLayout := "2006-01-02 15:04:05"

			var customerIds []string
			for idx, ea := range text {
				if idx > 0 {

					s := strings.Split(ea, "|")

					sidNo := strings.TrimSpace(s[2])
					ifuaNo := strings.TrimSpace(s[3])
					ifuaName := strings.TrimSpace(s[4])
					clientCode := strings.TrimSpace(s[5])

					//get ms_customer by clientCode
					var customer models.MsCustomer
					_, err := models.GetMsCustomerByClientCode(&customer, clientCode)
					if err != nil {
						log.Error("get customer error : client_code = " + clientCode + ". Error : " + err.Error())
						continue
					}

					strCusKey := strconv.FormatUint(customer.CustomerKey, 10)
					if _, ok := lib.Find(customerIds, strCusKey); !ok {
						customerIds = append(customerIds, strCusKey)
					}

					//update ms_customer
					params := make(map[string]string)
					params["sid_no"] = sidNo
					params["customer_key"] = strCusKey
					params["rec_modified_date"] = time.Now().Format(dateLayout)
					params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
					_, err = models.UpdateMsCustomer(params)
					if err != nil {
						log.Error("Error update sid_no ms_customer")
						continue
					}

					//update tr_account_all
					paramsTrAccount := make(map[string]string)
					paramsTrAccount["ifua_no"] = ifuaNo
					paramsTrAccount["ifua_name"] = ifuaName
					paramsTrAccount["rec_modified_date"] = time.Now().Format(dateLayout)
					paramsTrAccount["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
					_, err = models.UpdateTrAccountUploadSinvest(paramsTrAccount, "customer_key", strCusKey)
					if err != nil {
						log.Error("Error update ifua_no, ifua_name tr_account")
						continue
					}
				}
			}

			//update oa_status di oa_request by customer_key
			if len(customerIds) > 0 {
				paramsOa := make(map[string]string)
				paramsOa["oa_status"] = "262"
				paramsOa["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsOa["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
				_, err := models.UpdateOaRequestByFieldIn(paramsOa, customerIds, "customer_key")
				if err != nil {
					log.Error("Error update oa_status in oa_request : " + err.Error())
				}
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
