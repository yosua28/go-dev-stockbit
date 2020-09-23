package controllers

import(
	"api/lib"
	"api/models"
	"api/config"
	"net/http"
	"strconv"
	"mime/multipart"
	"path/filepath"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/badoux/checkmail"
	log "github.com/sirupsen/logrus"
)

func CreateOaPersonalData(c echo.Context) error {
	var err error
	params := make(map[string]string)

	// Address ID Parameters
	addressIDParams := make(map[string]string)

	addressTypeID := c.FormValue("address_type_idcard")
	if addressTypeID == "" {
		log.Error("Missing required parameter: address_type_idcard")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: address_type_idcard","Missing required parameter: address_type_idcard")
	}
	addressIDParams["address_type"] = addressTypeID

	addressid := c.FormValue("address_idcard")
	if addressid == "" {
		log.Error("Missing required parameter: address_idcard")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: address_idcard","Missing required parameter: address_idcard")
	}
	addressIDParams["address_line1"] = addressid

	kabupatenid := c.FormValue("kabupaten_idcard")
	if kabupatenid != "" {
		city, err := strconv.ParseUint(kabupatenid, 10, 64)
		if err == nil && city > 0{
			addressIDParams["kabupaten_key"] = kabupatenid
		}else{
			log.Error("Wrong input for parameter: kabupaten_idcard")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	kecamatanid := c.FormValue("kecamatan_idcard")
	if kecamatanid != "" {
		city, err := strconv.ParseUint(kecamatanid, 10, 64)
		if err == nil && city > 0{
			addressIDParams["kecamatan_key"] = kecamatanid
		}else{
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
		log.Error("Failed create adrress data: "+err.Error())
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

	addressDomicileParams["address_type"] = "Domicile"

	addressDomicile := c.FormValue("address_domicile")
	if addressDomicile == "" {
		log.Error("Missing required parameter: address_domicile")
		return lib.CustomError(http.StatusBadRequest)
	}
	addressDomicileParams["address_line1"] = addressDomicile

	kabupatenDomicile := c.FormValue("kabupaten_domicile")
	if kabupatenDomicile != "" {
		city, err := strconv.ParseUint(kabupatenDomicile, 10, 64)
		if err == nil && city > 0{
			addressDomicileParams["kabupaten_key"] = kabupatenDomicile
		}else{
			log.Error("Wrong input for parameter: kabupaten_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	kecamatanDomicile := c.FormValue("kecamatan_domicile")
	if kecamatanDomicile != "" {
		city, err := strconv.ParseUint(kecamatanDomicile, 10, 64)
		if err == nil && city > 0{
			addressDomicileParams["kecamatan_key"] = kecamatanDomicile
		}else{
			log.Error("Wrong input for parameter: kecamatan_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	postalDomicile := c.FormValue("postal_domicile")
	if postalDomicile != "" {
		city, err := strconv.ParseUint(postalDomicile, 10, 64)
		if err == nil && city > 0{
			addressDomicileParams["postal_code"] = postalDomicile
		}else{
			log.Error("Wrong input for parameter: postal_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	addressDomicileParams["rec_status"] = "1"

	status, err, addressDomicileID := models.CreateOaPostalAddress(addressDomicileParams)
	if err != nil {
		log.Error("Failed create adrress data: "+err.Error())
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

	nationality := c.FormValue("nationality")
	if nationality == "" {
		log.Error("Missing required parameter: nationality")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["nationality"] = nationality

	idcardNumber := c.FormValue("idcard_number")
	if idcardNumber == "" {
		log.Error("Missing required parameter: idcard_number")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["idcard_no"] = idcardNumber

	gender := c.FormValue("gender")
	if gender == "" {
		log.Error("Missing required parameter: gender")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["gender"] = gender

	maritalStatus := c.FormValue("marital_status")
	if maritalStatus == "" {
		log.Error("Missing required parameter: marital_status")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["marital_status"] = maritalStatus

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
	if phoneHome == "" {
		log.Error("Missing required parameter: email")
		return lib.CustomError(http.StatusBadRequest)
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
		return lib.CustomError(http.StatusBadRequest)
	}
	params["religion"] = religion

	education := c.FormValue("education")
	if education == "" {
		log.Error("Missing required parameter: education")
		return lib.CustomError(http.StatusBadRequest)
	}
	params["education"] = education

	
	err = os.MkdirAll(config.BasePath + "/images/user/" + strconv.FormatUint(lib.Profile.UserID, 10), 0755)
	if err != nil {
		log.Fatal(err)
	}else{
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
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams,false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath + "/images/user/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/" + filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
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
				_, err := models.GetAllOaPersonalData(&personalData, 1, 0, getParams,false)
				if (err == nil && len(personalData) < 1) || err != nil {
					break
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/user/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/" + filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
		} 
	}
	

	job := c.FormValue("job")
	if job == "" {
		log.Error("Missing required parameter: job")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: job","Missing required parameter: job")
	}
	params["occup_job"] = job

	company := c.FormValue("company")
	if company == "" {
		log.Error("Missing required parameter: company")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: company","Missing required parameter: company")
	}
	params["occup_company"] = company
	
	position := c.FormValue("position")
	if position == "" {
		log.Error("Missing required parameter: position")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: position","Missing required parameter: position")
	}
	params["occup_position"] = position
	
	addressCompanyParams := make(map[string]string)
	companyAddress := c.FormValue("company_address")
	if companyAddress == "" {
		log.Error("Missing required parameter: company_address")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: company_address","Missing required parameter: company_address")
	}
	addressCompanyParams["address_type"] = "Company"
	addressCompanyParams["address_line1"] = companyAddress
	addressCompanyParams["rec_status"] = "1"

	status, err, addressCompanyID := models.CreateOaPostalAddress(addressCompanyParams)
	if err != nil {
		log.Error("Failed create adrress data: "+err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	addressID, err = strconv.ParseUint(addressCompanyID, 10, 64)
	if addressID == 0 {
		log.Error("Failed create adrress data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	params["occup_address_key"] = addressCompanyID
	
	businessField := c.FormValue("business_field")
	if businessField == "" {
		log.Error("Missing required parameter: business_field")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: business_field","Missing required parameter: business_field")
	}
	params["occup_business_fields"] = businessField
	
	annualIncome := c.FormValue("annual_income")
	if annualIncome == "" {
		log.Error("Missing required parameter: annual_income")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: annual_income","Missing required parameter: annual_income")
	}
	params["annual_income"] = annualIncome
	
	fundSource := c.FormValue("fund_source")
	if fundSource == "" {
		log.Error("Missing required parameter: fund_source")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: fund_source","Missing required parameter: fund_source")
	}
	params["sourceof_fund"] = fundSource
	
	objectives := c.FormValue("objectives")
	if objectives == "" {
		log.Error("Missing required parameter: objectives")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: objectives","Missing required parameter: objectives")
	}
	params["invesment_objectives"] = objectives
	
	corespondence := c.FormValue("corespondence")
	if corespondence == "" {
		log.Error("Missing required parameter: corespondence")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: corespondence","Missing required parameter: corespondence")
	}
	params["correspondence"] = corespondence
	
	relationName := c.FormValue("relation_name")
	if relationName == "" {
		log.Error("Missing required parameter: relation_name")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: relation_name","Missing required parameter: relation_name")
	}
	params["relation_full_name"] = relationName
	
	relationOccupation := c.FormValue("relation_occupation")
	if relationOccupation == "" {
		log.Error("Missing required parameter: relation_occupation")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: relation_occupation","Missing required parameter: relation_occupation")
	}
	params["relation_occupation"] = relationOccupation
	
	relationBusinessField := c.FormValue("relation_business_field")
	if relationBusinessField == "" {
		log.Error("Missing required parameter: relation_business_field")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: relation_business_field","Missing required parameter: relation_business_field")
	}
	params["relation_business_fields"] = relationBusinessField
	
	MotherMaidenName := c.FormValue("mother_maiden_name")
	if MotherMaidenName == "" {
		log.Error("Missing required parameter: mother_maiden_name")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: mother_maiden_name","Missing required parameter: mother_maiden_name")
	}
	params["mother_maiden_name"] = MotherMaidenName
	
	emergencyName := c.FormValue("emergency_name")
	if emergencyName == "" {
		log.Error("Missing required parameter: emergency_name")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: emergency_name","Missing required parameter: emergency_name")
	}
	params["emergency_full_name"] = emergencyName
	
	emergencyRelation := c.FormValue("emergency_relation")
	if emergencyRelation == "" {
		log.Error("Missing required parameter: emergency_relation")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: emergency_relation","Missing required parameter: emergency_relation")
	}
	params["emergency_relation"] = emergencyRelation
	
	emergencyPhone := c.FormValue("emergency_phone")
	if emergencyPhone == "" {
		log.Error("Missing required parameter: emergency_phone")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: emergency_phone","Missing required parameter: emergency_phone")
	}
	params["emergency_phone_no"] = emergencyPhone
	
	beneficialName := c.FormValue("beneficial_name")
	if beneficialName == "" {
		log.Error("Missing required parameter: beneficial_name")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: beneficial_name","Missing required parameter: beneficial_name")
	}
	params["beneficial_full_name"] = beneficialName
	
	beneficialRelation := c.FormValue("beneficial_relation")
	if beneficialRelation == "" {
		log.Error("Missing required parameter: beneficial_relation")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: beneficial_relation","Missing required parameter: beneficial_relation")
	}
	params["beneficial_relation"] = beneficialRelation
	
	paramsBank := make(map[string]string)
	bankKey := c.FormValue("bank_key")
	if bankKey == "" {
		log.Error("Missing required parameter: bank_key")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: bank_key","Missing required parameter: bank_key")
	}else{
		bank, err := strconv.ParseUint(bankKey, 10, 64)
		if err == nil && bank > 0{
			paramsBank["bank_key"] = bankKey
		}else{
			log.Error("Wrong input for parameter: bank_key")
			return lib.CustomError(http.StatusBadRequest)
		}
	}
	
	accountNo := c.FormValue("account_no")
	if accountNo == "" {
		log.Error("Missing required parameter: account_no")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: account_no","Missing required parameter: account_no")
	}
	paramsBank["account_no"] = accountNo
	
	accountName := c.FormValue("account_name")
	if accountName == "" {
		log.Error("Missing required parameter: account_name")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: account_name","Missing required parameter: account_name")
	}
	paramsBank["account_holder_name"] = accountName
	
	branchName := c.FormValue("branch_name")
	if branchName == "" {
		log.Error("Missing required parameter: branch_name")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: branch_name","Missing required parameter: branch_name")
	}
	paramsBank["branch_name"] = branchName
	paramsBank["currency_key"] = "1"
	paramsBank["bank_account_type"] = "1"
	paramsBank["rec_domain"] = "1"
	paramsBank["rec_status"] = "1"
	
	status, err, bankAccountID := models.CreateMsBankAccount(paramsBank)
	if err != nil {
		log.Error("Failed create adrress data: "+err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	accountID, err := strconv.ParseUint(bankAccountID, 10, 64)
	if accountID == 0 {
		log.Error("Failed create adrress data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	params["bank_account_key"] = bankAccountID

	// Create Request
	dateLayout := "2006-01-02 15:04:05"
	date := time.Now().Format(dateLayout)
	paramsRequest := make(map[string]string)
	paramsRequest["oa_status"] = "NEW"
	paramsRequest["user_login_key"] = strconv.FormatUint(lib.Profile.UserID,10) 
	paramsRequest["oa_entry_start"] = date 
	paramsRequest["oa_entry_end"] = date 
	paramsRequest["oa_request_type"] = "NEW" 
	paramsRequest["rec_status"] = "1" 
	status, err, requestID := models.CreateOaRequest(paramsRequest)
	if err != nil {
		log.Error("Failed create request data: "+err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	request, err := strconv.ParseUint(requestID, 10, 64)
	if request == 0 {
		log.Error("Failed create adrress data")
		return lib.CustomError(http.StatusBadGateway, "failed input data", "failed input data")
	}
	params["oa_request_key"] = requestID
	params["rec_status"] = "1"
	params["idcard_type"] = addressTypeID

	status, err = models.CreateOaPersonalData(params)
	if err != nil {
		log.Error("Failed create personal data: "+err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}