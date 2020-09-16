package controllers

import(
	"api/lib"
	_ "api/models"
	_ "api/config"
	"net/http"
	"strconv"
	_ "mime/multipart"
	_ "path/filepath"

	"github.com/labstack/echo"
	"github.com/badoux/checkmail"
	log "github.com/sirupsen/logrus"
)

func CreateOaPersonalData(c echo.Context) error {
	var err error
	params := make(map[string]string)

	// Address KTP Parameters
	addressKtpParams := make(map[string]string)

	addressKtpParams["address_type"] = "ktp"

	addressKtp := c.FormValue("address_ktp")
	if addressKtp == "" {
		log.Error("Missing required parameter: address_ktp")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter: address_ktp","Missing required parameter: address_ktp")
	}
	addressKtpParams["address_line1"] = addressKtp

	kabupatenKtp := c.FormValue("kabupaten_ktp")
	if kabupatenKtp == "" {
		log.Error("Missing required parameter: kabupaten_ktp")
		return lib.CustomError(http.StatusBadRequest,)
	}else{
		city, err := strconv.ParseUint(kabupatenKtp, 10, 64)
		if err == nil && city > 0{
			addressKtpParams["kabupaten_key"] = kabupatenKtp
		}else{
			log.Error("Wrong input for parameter: kabupaten_ktp")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	kecamatanKtp := c.FormValue("kecamatan_ktp")
	if kecamatanKtp == "" {
		log.Error("Missing required parameter: kecamatan_ktp")
		return lib.CustomError(http.StatusBadRequest)
	}else{
		city, err := strconv.ParseUint(kecamatanKtp, 10, 64)
		if err == nil && city > 0{
			addressKtpParams["kecamatan_key"] = kecamatanKtp
		}else{
			log.Error("Wrong input for parameter: kecamatan_ktp")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	postalKtp := c.FormValue("postal_ktp")
	if postalKtp == "" {
		log.Error("Missing required parameter: postal_ktp")
		return lib.CustomError(http.StatusBadRequest)
	}else{
		city, err := strconv.ParseUint(postalKtp, 10, 64)
		if err == nil && city > 0{
			addressKtpParams["postal_code"] = postalKtp
		}else{
			log.Error("Wrong input for parameter: postal_ktp")
			return lib.CustomError(http.StatusBadRequest)
		}
	}
	
	// Address Domicile Parameters
	addressDomicileParams := make(map[string]string)

	addressDomicileParams["address_type"] = "domicile"

	addressDomicile := c.FormValue("address_domicile")
	if addressDomicile == "" {
		log.Error("Missing required parameter: address_domicile")
		return lib.CustomError(http.StatusBadRequest)
	}
	addressDomicileParams["address_line1"] = addressDomicile

	kabupatenDomicile := c.FormValue("kabupaten_domicile")
	if kabupatenDomicile == "" {
		log.Error("Missing required parameter: kabupaten_domicile")
		return lib.CustomError(http.StatusBadRequest)
	}else{
		city, err := strconv.ParseUint(kabupatenDomicile, 10, 64)
		if err == nil && city > 0{
			addressDomicileParams["kabupaten_key"] = kabupatenDomicile
		}else{
			log.Error("Wrong input for parameter: kabupaten_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	kecamatanDomicile := c.FormValue("kecamatan_domicile")
	if kecamatanDomicile == "" {
		log.Error("Missing required parameter: kecamatan_domicile")
		return lib.CustomError(http.StatusBadRequest)
	}else{
		city, err := strconv.ParseUint(kecamatanDomicile, 10, 64)
		if err == nil && city > 0{
			addressDomicileParams["kecamatan_key"] = kecamatanDomicile
		}else{
			log.Error("Wrong input for parameter: kecamatan_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	postalDomicile := c.FormValue("postal_ktp")
	if postalDomicile == "" {
		log.Error("Missing required parameter: postal_ktp")
		return lib.CustomError(http.StatusBadRequest)
	}else{
		city, err := strconv.ParseUint(postalDomicile, 10, 64)
		if err == nil && city > 0{
			addressDomicileParams["postal_code"] = postalDomicile
		}else{
			log.Error("Wrong input for parameter: postal_domicile")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

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
	params["email"] = email

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

	// Get parameter banner desktop
	// var file *multipart.FileHeader
	// file, err = c.FormFile("pic_ktp")
	// if file != nil {
	// 	if err != nil {
	// 		return lib.CustomError(http.StatusBadRequest)
	// 	}
	// 	// Get file extension
	// 	extension := filepath.Ext(file.Filename)
	// 	// Generate filename
	// 	var filename string
	// 	for {
	// 		filename = lib.RandStringBytesMaskImprSrc(20)
	// 		log.Println("Generate filename:", filename)
	// 		var personalData []models.OaPersonalData
	// 		getParams := make(map[string]string)
	// 		getParams["pic_ktp"] = filename + extension
	// 		total, _ := models.GetAllOaPersonalData(&personalData, 1, 0, getParams,false)
	// 		if total == 0 {
	// 			break
	// 		}
	// 	}
	// 	// Upload image and move to proper directory
	// 	err = lib.UploadImage(file, config.BasePath+"/images/user/"+filename+extension)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return lib.CustomError(http.StatusInternalServerError)
	// 	}
	// } 
	

	// Validate email
	err = checkmail.ValidateFormat(email)
	if err != nil {
		log.Error("Email format is not valid")
		return lib.CustomError(http.StatusBadRequest, "Email format is not valid", "Email format is not valid")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}