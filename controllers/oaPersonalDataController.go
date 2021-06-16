package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/badoux/checkmail"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func CreateOaPersonalData(c echo.Context) error {
	var err error
	params := make(map[string]string)
	var bindVar [][]string

	// Address ID Parameters
	addressIDParams := make(map[string]string)

	addressIDParams["address_type"] = "17"

	addressid := c.FormValue("address_idcard")
	if addressid == "" {
		log.Error("Missing required parameter: address_idcard")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: address_idcard", "Missing required parameter: address_idcard")
	}
	addressIDParams["address_line1"] = addressid

	kabupatenid := c.FormValue("kabupaten_idcard")
	if kabupatenid != "" {
		city, err := strconv.ParseUint(kabupatenid, 10, 64)
		if err == nil && city > 0 {
			addressIDParams["kabupaten_key"] = kabupatenid
		} else {
			log.Error("Wrong input for parameter: kabupaten_idcard")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	kecamatanid := c.FormValue("kecamatan_idcard")
	if kecamatanid != "" {
		city, err := strconv.ParseUint(kecamatanid, 10, 64)
		if err == nil && city > 0 {
			addressIDParams["kecamatan_key"] = kecamatanid
		} else {
			log.Error("Wrong input for parameter: kecamatan_idcard")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	postalid := c.FormValue("postal_idcard")
	if postalid != "" {
		addressIDParams["postal_code"] = postalid
	}

	addressIDParams["rec_status"] = "1"
	status, err, addressidID := models.CreateOaPostalAddress(addressIDParams)
	if err != nil {
		log.Error("Failed create adrress data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	addressID, err := strconv.ParseUint(addressidID, 10, 64)
	if addressID == 0 {
		log.Error("Failed create adrress data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	params["idcard_address_key"] = addressidID

	// Address Domicile Parameters
	addressDomicileParams := make(map[string]string)

	addressDomicileParams["address_type"] = "18"

	addressDomicile := c.FormValue("address_domicile")
	if addressDomicile == "" {
		log.Error("Missing required parameter: address_domicile")
		return lib.CustomError(http.StatusBadRequest)
	}
	addressDomicileParams["address_line1"] = addressDomicile

	kabupatenDomicile := c.FormValue("kabupaten_domicile")
	if kabupatenDomicile != "" {
		city, err := strconv.ParseUint(kabupatenDomicile, 10, 64)
		if err == nil && city > 0 {
			addressDomicileParams["kabupaten_key"] = kabupatenDomicile
		} else {
			log.Error("Wrong input for parameter: kabupaten_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	kecamatanDomicile := c.FormValue("kecamatan_domicile")
	if kecamatanDomicile != "" {
		city, err := strconv.ParseUint(kecamatanDomicile, 10, 64)
		if err == nil && city > 0 {
			addressDomicileParams["kecamatan_key"] = kecamatanDomicile
		} else {
			log.Error("Wrong input for parameter: kecamatan_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	postalDomicile := c.FormValue("postal_domicile")
	if postalDomicile != "" {
		city, err := strconv.ParseUint(postalDomicile, 10, 64)
		if err == nil && city > 0 {
			addressDomicileParams["postal_code"] = postalDomicile
		} else {
			log.Error("Wrong input for parameter: postal_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	addressDomicileParams["rec_status"] = "1"

	status, err, addressDomicileID := models.CreateOaPostalAddress(addressDomicileParams)
	if err != nil {
		log.Error("Failed create adrress data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	addressID, err = strconv.ParseUint(addressDomicileID, 10, 64)
	if addressID == 0 {
		log.Error("Failed create adrress data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}

	params["domicile_address_key"] = addressDomicileID

	// Check parameters
	fullName := c.FormValue("full_name")
	if fullName == "" {
		log.Error("Missing required parameter: full_name")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["full_name"] = fullName

	idcardType := c.FormValue("idcard_type")
	if idcardType == "" {
		log.Error("Missing required parameter: idcard_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: idcard_type", "Missing required parameter: idcard_type")
	} else {
		n, err := strconv.ParseUint(idcardType, 10, 64)
		if err == nil && n > 0 {
			params["idcard_type"] = idcardType
		} else {
			log.Error("Wrong input for parameter: idcard_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: idcard_type", "Wrong input for parameter: idcard_type")
		}
	}

	placeBirth := c.FormValue("place_birth")
	if placeBirth == "" {
		log.Error("Missing required parameter: place_birth")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["place_birth"] = placeBirth

	dateBirth := c.FormValue("date_birth")
	if dateBirth == "" {
		log.Error("Missing required parameter: date_birth")
		return lib.CustomError(http.StatusBadRequest)
	}
	log.Info("dateBirth: " + dateBirth)
	layout := "2006-01-02 15:04:05"
	dateBirth += " 00:00:00"
	date, err := time.Parse(layout, dateBirth)
	dateStr := date.Format(layout)
	log.Info("dateBirth: " + dateStr)
	params["date_birth"] = dateStr

	nationality := c.FormValue("nationality")
	if nationality == "" {
		log.Error("Missing required parameter: nationality")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		n, err := strconv.ParseUint(nationality, 10, 64)
		if err == nil && n > 0 {
			params["nationality"] = nationality

		} else {
			log.Error("Wrong input for parameter: nationality")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: nationality", "Wrong input for parameter: nationality")
		}
	}

	requestTypeStr := c.FormValue("request_type")
	if requestTypeStr == "" {
		log.Error("Missing required parameter: request_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: request_type", "Missing required parameter: request_type")
	} else {
		_, err := strconv.ParseUint(requestTypeStr, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: request_type")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	if requestTypeStr != "127" {
		var oaRequestDB []models.OaRequest
		paramsCek := make(map[string]string)
		paramsCek["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
		paramsCek["orderBy"] = "oa_request_key"
		paramsCek["rec_status"] = "1"
		paramsCek["orderType"] = "DESC"
		status, err := models.GetAllOaRequest(&oaRequestDB, 0, 0, true, paramsCek)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Oa Request not found")
		}

		if len(oaRequestDB) > 0 {
			if *oaRequestDB[0].Oastatus == uint64(258) || *oaRequestDB[0].Oastatus == uint64(259) {
				log.Error("last oa/pengkinian belum di approve")
				return lib.CustomError(http.StatusNotFound, "Data terakhir belum diapprove", "Data terakhir belum diapprove")
			}
		} else {
			log.Error("oa not found")
			return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
		}
	}

	idcardNumber := c.FormValue("idcard_number")
	if idcardNumber == "" {
		log.Error("Missing required parameter: idcard_number")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		if requestTypeStr == "127" {
			paramsPersonalData := make(map[string]string)
			paramsPersonalData["idcard_no"] = idcardNumber
			paramsPersonalData["rec_status"] = "1"
			var personalDataDB []models.OaPersonalData
			_, err = models.GetAllOaPersonalData(&personalDataDB, 0, 0, paramsPersonalData, true)
			if err != nil {
				log.Error("error get data")
				return lib.CustomError(http.StatusBadRequest, "Nomor kartu ID sudah pernah digunakan", "Nomor kartu ID sudah pernah digunakan")
			}
			if len(personalDataDB) > 0 {
				log.Error("idcard_number alredy used")
				return lib.CustomError(http.StatusBadRequest, "Nomor kartu ID sudah pernah digunakan", "Nomor kartu ID sudah pernah digunakan")
			}
		}
		params["idcard_no"] = idcardNumber
	}

	gender := c.FormValue("gender")
	if gender == "" {
		log.Error("Missing required parameter: gender")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		n, err := strconv.ParseUint(gender, 10, 64)
		if err == nil && n > 0 {
			params["gender"] = gender
		} else {
			log.Error("Wrong input for parameter: gender")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: gender", "Wrong input for parameter: gender")
		}
	}

	maritalStatus := c.FormValue("marital_status")
	if maritalStatus == "" {
		log.Error("Missing required parameter: marital_status")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		n, err := strconv.ParseUint(maritalStatus, 10, 64)
		if err == nil && n > 0 {
			params["marital_status"] = maritalStatus
		} else {
			log.Error("Wrong input for parameter: marital_status")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: marital_status", "Wrong input for parameter: marital_status")
		}
	}

	phoneHome := c.FormValue("phone_home")
	if phoneHome == "" {
		log.Error("Missing required parameter: phone_home")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["phone_home"] = phoneHome

	phoneMobile := c.FormValue("phone_mobile")
	if phoneHome == "" {
		log.Error("Missing required parameter: phone_home")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["phone_mobile"] = phoneMobile

	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter: email")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: email", "Missing required parameter: email")
	}
	// Validate email
	err = checkmail.ValidateFormat(email)
	if err != nil {
		log.Error("Email format is not valid")
		return lib.CustomError(http.StatusBadRequest, "Email format is not valid", "Email format is not valid")
	}
	params["email_address"] = email

	religion := c.FormValue("religion")
	if religion == "" {
		log.Error("Missing required parameter: religion")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: religion", "Missing required parameter: religion")
	} else {
		n, err := strconv.ParseUint(religion, 10, 64)
		if err == nil && n > 0 {
			params["religion"] = religion
		} else {
			log.Error("Wrong input for parameter: religion")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: religion", "Wrong input for parameter: religion")
		}
	}

	religionOther := c.FormValue("religion_other")
	if religionOther != "" {
		var row []string
		row = append(row, "1")
		row = append(row, "0")
		row = append(row, religionOther)
		bindVar = append(bindVar, row)
	}

	education := c.FormValue("education")
	if education == "" {
		log.Error("Missing required parameter: education")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		n, err := strconv.ParseUint(education, 10, 64)
		if err == nil && n > 0 {
			params["education"] = education
		} else {
			log.Error("Wrong input for parameter: education")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: education", "Wrong input for parameter: education")
		}
	}

	educationOther := c.FormValue("education_other")
	if educationOther != "" {
		var row []string
		row = append(row, "3")
		row = append(row, "0")
		row = append(row, educationOther)
		bindVar = append(bindVar, row)
	}

	err = os.MkdirAll(config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10), 0755)
	if err != nil {
		log.Error(err.Error())
	} else { //if requestTypeStr == "127"
		var file *multipart.FileHeader
		file, err = c.FormFile("pic_ktp")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				log.Println("Generate filename:", filename)
				var personalData []models.OaPersonalData
				getParams := make(map[string]string)
				getParams["pic_ktp"] = filename + extension
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["pic_ktp"] = filename + extension
		} else {
			picKtp := c.FormValue("pic_ktp_str")
			if picKtp == "" {
				log.Error("Missing required parameter: pic_ktp/pic_ktp_str")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: pic_ktp/pic_ktp_str", "Missing required parameter: pic_ktp_str")
			} else {
				params["pic_ktp"] = picKtp
			}
		}

		file, err = c.FormFile("pic_selfie_ktp")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				log.Println("Generate filename:", filename)
				var personalData []models.OaPersonalData
				getParams := make(map[string]string)
				getParams["pic_selfie_ktp"] = filename + extension
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["pic_selfie_ktp"] = filename + extension
		} else {
			picSelfie := c.FormValue("pic_selfie_ktp_str")
			if picSelfie == "" {
				log.Error("Missing required parameter: pic_selfie_ktp/pic_selfie_ktp_str")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: pic_selfie_ktp/pic_selfie_ktp_str", "Missing required parameter: pic_selfie_ktp/pic_selfie_ktp_str")
			} else {
				params["pic_selfie_ktp"] = picSelfie
			}
		}

		err = os.MkdirAll(config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/signature", 0755)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(http.StatusBadGateway, err.Error(), err.Error())
		}
		file, err = c.FormFile("signature")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: signature")
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				log.Println("Generate filename:", filename)
				var personalData []models.OaPersonalData
				getParams := make(map[string]string)
				getParams["rec_image1"] = filename + extension
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams, false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/signature/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image1"] = filename + extension
		} else {
			signature := c.FormValue("signature_str")
			if signature == "" {
				log.Error("Missing required parameter: signature/signature_str")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: signature/signature_str", "Missing required parameter: signature/signature_str")
			} else {
				params["rec_image1"] = signature
			}
		}

	}
	// else {
	// 	picKtp := c.FormValue("pic_ktp_str")
	// 	if picKtp == "" {
	// 		log.Error("Missing required parameter: pic_ktp_str")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: pic_ktp_str", "Missing required parameter: pic_ktp_str")
	// 	} else {
	// 		params["pic_ktp"] = picKtp
	// 	}
	// 	picSelfie := c.FormValue("pic_selfie_ktp_str")
	// 	if picSelfie == "" {
	// 		log.Error("Missing required parameter: pic_selfie_ktp_str")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: pic_selfie_ktp_str", "Missing required parameter: pic_selfie_ktp_str")
	// 	} else {
	// 		params["pic_selfie_ktp"] = picSelfie
	// 	}
	// 	signature := c.FormValue("signature_str")
	// 	if signature == "" {
	// 		log.Error("Missing required parameter: signature_str")
	// 		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: signature_str", "Missing required parameter: signature_str")
	// 	} else {
	// 		params["rec_image1"] = signature
	// 	}
	// }

	job := c.FormValue("job")
	if job == "" {
		log.Error("Missing required parameter: job")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: job", "Missing required parameter: job")
	} else {
		n, err := strconv.ParseUint(job, 10, 64)
		if err == nil && n > 0 {
			params["occup_job"] = job
		} else {
			log.Error("Wrong input for parameter: job")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: job", "Wrong input for parameter: job")
		}
	}

	jobOther := c.FormValue("job_other")
	if jobOther != "" {
		var row []string
		row = append(row, "2")
		row = append(row, "0")
		row = append(row, jobOther)
		bindVar = append(bindVar, row)
	}

	company := c.FormValue("company")
	if company != "" {
		params["occup_company"] = company
	}

	position := c.FormValue("position")
	if position != "" {
		n, err := strconv.ParseUint(job, 10, 64)
		if err == nil && n > 0 {
			params["occup_position"] = position
		} else {
			log.Error("Wrong input for parameter: position")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: position", "Wrong input for parameter: position")
		}
	}

	positionOther := c.FormValue("position_other")
	if positionOther != "" {
		var row []string
		row = append(row, "8")
		row = append(row, "0")
		row = append(row, positionOther)
		bindVar = append(bindVar, row)
	}

	addressCompanyParams := make(map[string]string)
	companyAddress := c.FormValue("company_address")
	if companyAddress != "" {
		addressCompanyParams["address_type"] = "19"
		addressCompanyParams["address_line1"] = companyAddress
		addressCompanyParams["rec_status"] = "1"

		status, err, addressCompanyID := models.CreateOaPostalAddress(addressCompanyParams)
		if err != nil {
			log.Error("Failed create adrress data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
		addressID, err = strconv.ParseUint(addressCompanyID, 10, 64)
		if addressID == 0 {
			log.Error("Failed create adrress data")
			return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
		}
		params["occup_address_key"] = addressCompanyID
	}

	businessField := c.FormValue("business_field")
	if businessField != "" {
		n, err := strconv.ParseUint(businessField, 10, 64)
		if err == nil && n > 0 {
			params["occup_business_fields"] = businessField
		} else {
			log.Error("Wrong input for parameter: business_field")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: business_field", "Wrong input for parameter: business_field")
		}
	}

	businessFieldOther := c.FormValue("business_field_other")
	if businessFieldOther != "" {
		var row []string
		row = append(row, "4")
		row = append(row, "0")
		row = append(row, businessFieldOther)
		bindVar = append(bindVar, row)
	}

	annualIncome := c.FormValue("annual_income")
	if annualIncome == "" {
		log.Error("Missing required parameter: annual_income")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: annual_income", "Missing required parameter: annual_income")
	} else {
		n, err := strconv.ParseUint(annualIncome, 10, 64)
		if err == nil && n > 0 {
			params["annual_income"] = annualIncome
		} else {
			log.Error("Wrong input for parameter: annual_income")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: annual_income", "Wrong input for parameter: annual_income")
		}
	}

	fundSource := c.FormValue("fund_source")
	if fundSource == "" {
		log.Error("Missing required parameter: fund_source")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: fund_source", "Missing required parameter: fund_source")
	} else {
		n, err := strconv.ParseUint(fundSource, 10, 64)
		if err == nil && n > 0 {
			params["sourceof_fund"] = fundSource
		} else {
			log.Error("Wrong input for parameter: fund_source")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: fund_source", "Wrong input for parameter: fund_source")
		}
	}

	fundSourceOther := c.FormValue("fund_source_other")
	if fundSourceOther != "" {
		var row []string
		row = append(row, "5")
		row = append(row, "0")
		row = append(row, fundSourceOther)
		bindVar = append(bindVar, row)
	}

	objectives := c.FormValue("objectives")
	if objectives == "" {
		log.Error("Missing required parameter: objectives")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: objectives", "Missing required parameter: objectives")
	} else {
		n, err := strconv.ParseUint(objectives, 10, 64)
		if err == nil && n > 0 {
			params["invesment_objectives"] = objectives
		} else {
			log.Error("Wrong input for parameter: objectives")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: objectives", "Wrong input for parameter: objectives")
		}
	}

	objectivesOther := c.FormValue("objectives_other")
	if objectivesOther != "" {
		var row []string
		row = append(row, "6")
		row = append(row, "0")
		row = append(row, objectivesOther)
		bindVar = append(bindVar, row)
	}

	corespondence := c.FormValue("corespondence")
	if corespondence == "" {
		log.Error("Missing required parameter: corespondence")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: corespondence", "Missing required parameter: corespondence")
	} else {
		n, err := strconv.ParseUint(corespondence, 10, 64)
		if err == nil && n > 0 {
			params["correspondence"] = corespondence
		} else {
			log.Error("Wrong input for parameter: corespondence")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: corespondence", "Wrong input for parameter: corespondence")
		}
	}

	relationName := c.FormValue("relation_name")
	if relationName == "" {
		log.Error("Missing required parameter: relation_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: relation_name", "Missing required parameter: relation_name")
	}
	params["relation_full_name"] = relationName

	relationOccupation := c.FormValue("relation_occupation")
	if relationOccupation == "" {
		log.Error("Missing required parameter: relation_occupation")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: relation_occupation", "Missing required parameter: relation_occupation")
	} else {
		n, err := strconv.ParseUint(corespondence, 10, 64)
		if err == nil && n > 0 {
			params["relation_occupation"] = relationOccupation
		} else {
			log.Error("Wrong input for parameter: relation_occupation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_occupation", "Wrong input for parameter: relation_occupation")
		}
	}

	relationOccupationOther := c.FormValue("relation_occupation_other")
	if relationOccupationOther != "" {
		var row []string
		row = append(row, "9")
		row = append(row, "0")
		row = append(row, relationOccupationOther)
		bindVar = append(bindVar, row)
	}

	relationBusinessField := c.FormValue("relation_business_field")
	if relationBusinessField != "" {
		n, err := strconv.ParseUint(corespondence, 10, 64)
		if err == nil && n > 0 {
			params["relation_business_fields"] = relationBusinessField
		} else {
			log.Error("Wrong input for parameter: relation_business_field")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_business_field", "Wrong input for parameter: relation_business_field")
		}
	}

	relationBusinessFieldOther := c.FormValue("relation_business_field_other")
	if relationBusinessFieldOther != "" {
		var row []string
		row = append(row, "10")
		row = append(row, "0")
		row = append(row, relationBusinessFieldOther)
		bindVar = append(bindVar, row)
	}

	MotherMaidenName := c.FormValue("mother_maiden_name")
	if MotherMaidenName == "" {
		log.Error("Missing required parameter: mother_maiden_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: mother_maiden_name", "Missing required parameter: mother_maiden_name")
	}
	params["mother_maiden_name"] = MotherMaidenName

	emergencyName := c.FormValue("emergency_name")
	if emergencyName == "" {
		log.Error("Missing required parameter: emergency_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: emergency_name", "Missing required parameter: emergency_name")
	}
	params["emergency_full_name"] = emergencyName

	emergencyRelation := c.FormValue("emergency_relation")
	if emergencyRelation == "" {
		log.Error("Missing required parameter: emergency_relation")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: emergency_relation", "Missing required parameter: emergency_relation")
	} else {
		n, err := strconv.ParseUint(emergencyRelation, 10, 64)
		if err == nil && n > 0 {
			params["emergency_relation"] = emergencyRelation
		} else {
			log.Error("Wrong input for parameter: emergency_relation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: emergency_relation", "Wrong input for parameter: emergency_relation")
		}
	}

	emergencyPhone := c.FormValue("emergency_phone")
	if emergencyPhone == "" {
		log.Error("Missing required parameter: emergency_phone")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: emergency_phone", "Missing required parameter: emergency_phone")
	}
	params["emergency_phone_no"] = emergencyPhone

	beneficialName := c.FormValue("beneficial_name")
	if beneficialName == "" {
		log.Error("Missing required parameter: beneficial_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: beneficial_name", "Missing required parameter: beneficial_name")
	}
	params["beneficial_full_name"] = beneficialName

	beneficialRelation := c.FormValue("beneficial_relation")
	if beneficialRelation == "" {
		log.Error("Missing required parameter: beneficial_relation")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: beneficial_relation", "Missing required parameter: beneficial_relation")
	} else {
		n, err := strconv.ParseUint(beneficialRelation, 10, 64)
		if err == nil && n > 0 {
			params["beneficial_relation"] = beneficialRelation
		} else {
			log.Error("Wrong input for parameter: beneficial_relation")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: beneficial_relation", "Wrong input for parameter: beneficial_relation")
		}
	}

	relationType := c.FormValue("relation_type")
	if relationType == "" {
		log.Error("Missing required parameter: relation_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: relation_type", "Missing required parameter: relation_type")
	} else {
		n, err := strconv.ParseUint(relationType, 10, 64)
		if err == nil && n > 0 {
			params["relation_type"] = relationType
		} else {
			log.Error("Wrong input for parameter: relation_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: relation_type", "Wrong input for parameter: relation_type")
		}
	}

	beneficialRelationOther := c.FormValue("beneficial_relation_other")
	if beneficialRelationOther != "" {
		var row []string
		row = append(row, "7")
		row = append(row, "0")
		row = append(row, beneficialRelationOther)
		bindVar = append(bindVar, row)
	}

	pepStatus := c.FormValue("pep_status")
	if pepStatus == "" {
		log.Error("Missing required parameter: pep_status")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: pep_status", "Missing required parameter: pep_status")
	} else {
		n, err := strconv.ParseUint(pepStatus, 10, 64)
		if err == nil && n > 0 {
			params["pep_status"] = pepStatus
		} else {
			log.Error("Wrong input for parameter: pep_status")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: pep_status", "Wrong input for parameter: pep_status")
		}
	}

	paramsBank := make(map[string]string)
	bankKey := c.FormValue("bank_key")
	if bankKey == "" {
		log.Error("Missing required parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_key", "Missing required parameter: bank_key")
	} else {
		bank, err := strconv.ParseUint(bankKey, 10, 64)
		if err == nil && bank > 0 {
			paramsBank["bank_key"] = bankKey
		} else {
			log.Error("Wrong input for parameter: bank_key")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	accountNo := c.FormValue("account_no")
	if accountNo == "" {
		log.Error("Missing required parameter: account_no")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_no", "Missing required parameter: account_no")
	}
	paramsBank["account_no"] = accountNo

	accountName := c.FormValue("account_name")
	if accountName == "" {
		log.Error("Missing required parameter: account_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: account_name", "Missing required parameter: account_name")
	}
	paramsBank["account_holder_name"] = accountName

	branchName := c.FormValue("branch_name")
	if branchName == "" {
		log.Error("Missing required parameter: branch_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: branch_name", "Missing required parameter: branch_name")
	}
	paramsBank["branch_name"] = branchName
	paramsBank["currency_key"] = "1"
	paramsBank["bank_account_type"] = "1"
	paramsBank["rec_domain"] = "1"
	paramsBank["rec_status"] = "1"

	status, err, bankAccountID := models.CreateMsBankAccount(paramsBank)
	if err != nil {
		log.Error("Failed create adrress data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	accountID, err := strconv.ParseUint(bankAccountID, 10, 64)
	if accountID == 0 {
		log.Error("Failed create adrress data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	params["bank_account_key"] = bankAccountID

	// Create Request
	dateNow := time.Now().Format(layout)
	paramsRequest := make(map[string]string)
	paramsRequest["oa_status"] = "258"
	paramsRequest["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsRequest["oa_entry_start"] = dateNow
	paramsRequest["oa_entry_end"] = dateNow
	paramsRequest["oa_request_type"] = requestTypeStr
	paramsRequest["branch_key"] = "1"
	paramsRequest["agent_key"] = "1"
	salesCode := c.FormValue("sales_code")
	if salesCode != "" {
		paramsRequest["sales_code"] = salesCode
		paramsAgent := make(map[string]string)
		paramsAgent["agent_code"] = salesCode
		paramsAgent["rec_status"] = "1"
		var agentDB []models.MsAgent
		_, err = models.GetAllMsAgent(&agentDB, 0, 0, paramsAgent, true)
		if err == nil && len(agentDB) > 0 {
			agentKeyStr := strconv.FormatUint(agentDB[0].AgentKey, 10)
			paramsAgentBranch := make(map[string]string)
			paramsAgentBranch["agent_key"] = agentKeyStr
			paramsAgentBranch["rec_status"] = "1"
			paramsAgentBranch["orderBy"] = "eff_date"
			paramsAgentBranch["orderType"] = "DESC"
			var agentBranchDB []models.MsAgentBranch
			_, err = models.GetAllMsAgentBranch(&agentBranchDB, 0, 0, paramsAgentBranch, true)
			if err == nil && len(agentDB) > 0 {
				paramsRequest["branch_key"] = strconv.FormatUint(agentBranchDB[0].BranchKey, 10)
				paramsRequest["agent_key"] = agentKeyStr
			}
		}
	}
	paramsRequest["rec_status"] = "1"
	paramsRequest["rec_created_date"] = dateNow
	paramsRequest["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err, requestID := models.CreateOaRequest(paramsRequest)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	request, err := strconv.ParseUint(requestID, 10, 64)
	if request == 0 {
		log.Error("Failed create adrress data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	params["oa_request_key"] = requestID
	params["rec_status"] = "1"
	params["rec_created_date"] = dateNow
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err, requestKey := models.CreateOaPersonalData(params)
	if err != nil {
		log.Error("Failed create personal data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var bindInterface []interface{}
	for i := 0; i < len(bindVar); i++ {
		bindVar[i][1] = requestKey
		bindInterface = append(bindInterface, bindVar[i])
	}

	status, err = models.CreateMultipleUdfValue(bindInterface)
	if err != nil {
		log.Error(err.Error())
	}

	responseData := make(map[string]string)
	responseData["request_key"] = requestID
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func GetOaPersonalData(c echo.Context) error {

	var oaRequestDB []models.OaRequest
	params := make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["orderBy"] = "oa_request_key"
	params["rec_status"] = "1"
	params["orderType"] = "DESC"
	status, err := models.GetAllOaRequest(&oaRequestDB, 0, 0, true, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Oa Request not found")
	}
	var requestKey string
	if len(oaRequestDB) > 0 {
		requestKey = strconv.FormatUint(oaRequestDB[0].OaRequestKey, 10)
	} else {
		log.Error("oa not found")
		return lib.CustomError(http.StatusNotFound, "Oa Request not found", "Oa Request not found")
	}

	var personalDataDB models.OaPersonalData
	if requestKey != "" {
		status, err = models.GetOaPersonalData(&personalDataDB, requestKey, "oa_request_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Oa Request not found")
		}
	}

	responseData := make(map[string]interface{})
	responseData["full_name"] = personalDataDB.FullName
	responseData["place_birth"] = personalDataDB.PlaceBirth
	responseData["date_birth"] = personalDataDB.DateBirth
	responseData["nationality"] = personalDataDB.Nationality
	responseData["idcard_type"] = personalDataDB.IDcardType
	responseData["idcard_no"] = personalDataDB.IDcardNo
	responseData["idcard_expired_date"] = personalDataDB.IDcardExpiredDate
	responseData["idcard_never_expired"] = personalDataDB.IDcardNeverExpired
	responseData["gender"] = personalDataDB.Gender
	if personalDataDB.PepStatus == nil {
		responseData["pep_status"] = uint64(322)
	} else {
		responseData["pep_status"] = personalDataDB.PepStatus
	}
	responseData["marital_status"] = personalDataDB.MaritalStatus
	var address models.OaPostalAddress
	_, err = models.GetOaPostalAddress(&address, strconv.FormatUint(*personalDataDB.IDcardAddressKey, 10))
	if err == nil {
		addressID := make(map[string]interface{})
		addressID["postal_address_key"] = address.PostalAddressKey
		addressID["kabupaten_key"] = address.KabupatenKey
		addressID["kecamatan_key"] = address.KecamatanKey
		addressID["address_line1"] = address.AddressLine1
		addressID["address_line2"] = address.AddressLine2
		addressID["address_line3"] = address.AddressLine3
		addressID["postal_code"] = address.PostalCode
		responseData["idcard_address"] = addressID
	}
	_, err = models.GetOaPostalAddress(&address, strconv.FormatUint(*personalDataDB.DomicileAddressKey, 10))
	if err == nil {
		addressID := make(map[string]interface{})
		addressID["postal_address_key"] = address.PostalAddressKey
		addressID["kabupaten_key"] = address.KabupatenKey
		addressID["kecamatan_key"] = address.KecamatanKey
		addressID["address_line1"] = address.AddressLine1
		addressID["address_line2"] = address.AddressLine2
		addressID["address_line3"] = address.AddressLine3
		addressID["postal_code"] = address.PostalCode
		responseData["domicile_address"] = addressID
	}
	responseData["phone_home"] = personalDataDB.PhoneHome
	responseData["phone_mobile"] = personalDataDB.PhoneMobile
	responseData["email"] = personalDataDB.EmailAddress
	responseData["religion"] = personalDataDB.Religion
	dir := config.BaseUrl + "/images/user/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/"
	selfie := make(map[string]interface{})
	if personalDataDB.PicSelfie != nil && *personalDataDB.PicSelfie != "" {
		selfie["file_name"] = *personalDataDB.PicSelfie
		selfie["full_url"] = dir + *personalDataDB.PicSelfie
	}
	responseData["pic_selfie"] = selfie

	ktp := make(map[string]interface{})
	if personalDataDB.PicKtp != nil && *personalDataDB.PicKtp != "" {
		ktp["file_name"] = *personalDataDB.PicKtp
		ktp["full_url"] = dir + *personalDataDB.PicKtp
	}
	responseData["pic_ktp"] = ktp

	selfie_ktp := make(map[string]interface{})
	if personalDataDB.PicSelfieKtp != nil && *personalDataDB.PicSelfieKtp != "" {
		selfie_ktp["file_name"] = *personalDataDB.PicSelfieKtp
		selfie_ktp["full_url"] = dir + *personalDataDB.PicSelfieKtp
	}
	responseData["pic_selfie_ktp"] = selfie_ktp

	signature := make(map[string]interface{})
	if personalDataDB.RecImage1 != nil && *personalDataDB.RecImage1 != "" {
		signature["file_name"] = *personalDataDB.RecImage1
		signature["full_url"] = dir + "signature/" + *personalDataDB.RecImage1
	}
	responseData["signature"] = signature

	responseData["education"] = personalDataDB.Education
	responseData["occup_job"] = personalDataDB.OccupJob
	responseData["occup_company"] = personalDataDB.OccupCompany
	responseData["occup_position"] = personalDataDB.OccupPosition
	_, err = models.GetOaPostalAddress(&address, strconv.FormatUint(*personalDataDB.OccupAddressKey, 10))
	if err == nil {
		addressID := make(map[string]interface{})
		addressID["postal_address_key"] = address.PostalAddressKey
		addressID["kabupaten_key"] = address.KabupatenKey
		addressID["kecamatan_key"] = address.KecamatanKey
		addressID["address_line1"] = address.AddressLine1
		addressID["address_line2"] = address.AddressLine2
		addressID["address_line3"] = address.AddressLine3
		addressID["postal_code"] = address.PostalCode
		responseData["occup_address"] = addressID
	}
	responseData["occup_business_field"] = personalDataDB.OccupBusinessFields
	responseData["occup_phone"] = personalDataDB.OccupPhone
	responseData["occup_web_url"] = personalDataDB.OccupWebUrl
	responseData["correspondence"] = personalDataDB.Correspondence
	responseData["annual_income"] = personalDataDB.AnnualIncome
	responseData["sourceof_fund"] = personalDataDB.SourceofFund
	responseData["invesment_objectives"] = personalDataDB.InvesmentObjectives
	responseData["relation_type"] = personalDataDB.RelationType
	responseData["relation_full_name"] = personalDataDB.RelationFullName
	responseData["relation_occupation"] = personalDataDB.RelationOccupation
	responseData["relation_business_fields"] = personalDataDB.RelationBusinessFields
	responseData["mother_maiden_name"] = personalDataDB.MotherMaidenName
	responseData["emergency_full_name"] = personalDataDB.EmergencyFullName
	responseData["emergency_relation"] = personalDataDB.EmergencyRelation
	responseData["emergency_phone_no"] = personalDataDB.EmergencyPhoneNo
	responseData["beneficial_full_name"] = personalDataDB.BeneficialFullName
	responseData["beneficial_relation"] = personalDataDB.BeneficialRelation
	var bankAccountDB models.MsBankAccount
	if personalDataDB.BankAccountKey != nil && *personalDataDB.BankAccountKey > 0 {
		_, err = models.GetBankAccount(&bankAccountDB, strconv.FormatUint(*personalDataDB.BankAccountKey, 10))
		if err == nil {
			bankAccount := make(map[string]interface{})
			bankAccount["bank_key"] = bankAccountDB.BankKey
			bankAccount["account_no"] = bankAccountDB.AccountNo
			bankAccount["account_holder_name"] = bankAccountDB.AccountHolderName
			bankAccount["branch_name"] = bankAccountDB.BranchName
			responseData["bank_account"] = bankAccount
		}
	}

	var requestDB []models.OaRequest
	paramRequest := make(map[string]string)
	paramRequest["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramRequest["orderBy"] = "oa_request_key"
	paramRequest["orderType"] = "DESC"
	_, err = models.GetAllOaRequest(&requestDB, 1, 0, false, paramRequest)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	request := requestDB[0]
	var quizDB []models.OaRiskProfileQuiz
	paramQuiz := make(map[string]string)
	paramQuiz["oa_request_key"] = strconv.FormatUint(request.OaRequestKey, 10)
	paramQuiz["orderBy"] = "oa_request_key"
	paramQuiz["orderType"] = "DESC"
	_, err = models.GetAllOaRiskProfileQuiz(&quizDB, 100, 0, paramQuiz, true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var risk models.OaRiskProfile
	_, err = models.GetOaRiskProfile(&risk, strconv.FormatUint(request.OaRequestKey, 10), "oa_request_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var riskProfile models.MsRiskProfile
	_, err = models.GetMsRiskProfile(&riskProfile, strconv.FormatUint(risk.RiskProfileKey, 10))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	riskProfileData := make(map[string]interface{})
	riskProfileData["score_result"] = risk.ScoreResult
	riskProfileData["risk_code"] = riskProfile.RiskCode
	riskProfileData["risk_name"] = riskProfile.RiskName
	riskProfileData["risk_desc"] = riskProfile.RiskDesc
	var quizData []interface{}
	for _, q := range quizDB {
		quiz := make(map[string]interface{})
		quiz["question_key"] = q.QuizQuestionKey
		quiz["option_key"] = q.QuizOptionKey
		quiz["score"] = q.QuizOptionScore
		quizData = append(quizData, quiz)
	}
	riskProfileData["quiz"] = quizData

	responseData["risk_profile"] = riskProfileData
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func IDCardNumberValidation(c echo.Context) error {
	idcardNumber := c.QueryParam("idcard_number")
	paramsPersonalData := make(map[string]string)
	paramsPersonalData["idcard_no"] = idcardNumber
	paramsPersonalData["rec_status"] = "1"
	var personalDataDB []models.OaPersonalData
	_, err := models.GetAllOaPersonalData(&personalDataDB, 0, 0, paramsPersonalData, true)
	if err != nil {
		log.Error("error get data")
		return lib.CustomError(http.StatusBadRequest, "No. Identitas kamu telah terdaftar.", "No. Identitas kamu telah terdaftar.")
	}
	if len(personalDataDB) > 0 {
		log.Error("idcard_number alredy used")
		return lib.CustomError(http.StatusBadRequest, "No. Identitas kamu telah terdaftar.", "No. Identitas kamu telah terdaftar.")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func SalesCodeValidation(c echo.Context) error {
	salesCode := c.FormValue("sales_code")
	if salesCode != "" {
		paramsAgent := make(map[string]string)
		paramsAgent["agent_code"] = salesCode
		paramsAgent["rec_status"] = "1"
		var agentDB []models.MsAgent
		_, err := models.GetAllMsAgent(&agentDB, 0, 0, paramsAgent, true)
		if err != nil {
			log.Error("error get data")
			return lib.CustomError(http.StatusNotFound, "Sales code tidak dapat digunakan", "Sales code tidak dapat digunakan")
		}
		if len(agentDB) == 0 {
			log.Error("error get data")
			return lib.CustomError(http.StatusNotFound, "Sales code tidak dapat digunakan", "Sales code tidak dapat digunakan")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
