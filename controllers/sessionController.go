package controllers

import (
	"api/lib"
	"api/models"
	"api/config"
	"net/http"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/base64"
	"time"
	"unicode"
	"strconv"

	"github.com/labstack/echo"
	"github.com/badoux/checkmail"
	log "github.com/sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gomail.v2"
)

func Register(c echo.Context) error {
	var err error
	var status int
	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}
	password := c.FormValue("password")
	if password == "" {
		log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}
	phone := c.FormValue("phone")
	if phone == "" {
		log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}

	// Validate email
	err = checkmail.ValidateFormat(email)
	if err != nil {
		log.Error("Email format is not valid")
		return lib.CustomError(http.StatusBadRequest, "Email format is not valid", "Email format is not valid")
	}
	// err = checkmail.ValidateHost(email)
	// if err != nil {
	// 	log.Error("Email host is not valid")
	// 	return ctx.TextResponse("Email is not valid", fasthttp.StatusBadRequest)
	// }
	// if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
	// 	log.Error("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
	// 	return ctx.TextResponse("Email is not valid", fasthttp.StatusBadRequest)
	// }
	var user []models.ScUserLogin
	params := make(map[string]string)
	params["ulogin_email"] = email 
	status, err = models.GetAllScUserLogin(&user, 0, 0, params, true)
	if err != nil {
		log.Error("Error get email " + email)
		return lib.CustomError(status, err.Error(), "Error get email")
	}
	if len(user) > 0 {
		log.Error("Email " + email + " already registered")
		return lib.CustomError(http.StatusBadRequest, "Email "+email+" already registered", "Email "+email+" already registered")
	}

	// Validate password
	length, number, upper, special := verifyPassword(password)
	if length == false || number == false || upper == false || special == false {
		log.Error("Password does meet the criteria")
		return lib.CustomError(http.StatusBadRequest, "Password does meet the criteria", "Your password need at least 8 character length, has lower and upper case letter, has numeric letter, and has special character")
	}

	// Encrypt password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])

	// Set expired for token
	date := time.Now().AddDate(0, 0, 1)
	dateLayout := "2006-01-02 15:04:05"
	expired := date.Format(dateLayout)

	// Generate verify key
	verifyKeyByte := sha256.Sum256([]byte(email + "_" + expired))
	verifyKey := hex.EncodeToString(verifyKeyByte[:])

	// Input to database
	params["ulogin_email"] = email 
	params["ulogin_name"] = email 
	params["ulogin_must_changepwd"] = "0" 
	params["user_category_key"] = "1" 
	params["user_dept_key"] = "1" 
	params["last_password_changed"] = time.Now().Format(dateLayout) 
	params["ulogin_password"] = encryptedPassword 
	params["verified_email"] = "0" 
	params["verified_mobileno"] = "0" 
	params["ulogin_mobileno"] = phone 
	params["ulogin_enabled"] = "1" 
	params["ulogin_locked"] = "0" 
	params["ulogin_failed_count"] = "0" 
	params["user_category_key"] = "1" 
	params["last_access"] = time.Now().Format(dateLayout) 
	params["accept_login_tnc"] = "1" 
	params["allowed_sharing_login"] = "1" 
	params["string_token"] = verifyKey 
	params["token_expired"] = expired 
	params["rec_status"] = "1" 

	status, err = models.CreateScUserLogin(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create user")
	}

	// Send email
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "[MNCduit] Verify your email address")
	mailer.SetBody("text/html", "Hi,<br><br>To complete register MNCduit account, please verify your email:<br><br>https://devapi.mncasset.com/api/verifyemail?token="+verifyKey+"<br><br>Thank you,<br>MNCduit Team")

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
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Error send email")
	}
	log.Info("Email sent")
	var data models.ScUserLoginRegister
	data.UloginEmail = email
	data.UloginMobileno = phone

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	return c.JSON(http.StatusOK, response)
}

func VerifyEmail(c echo.Context) error {
	var err error
	// Get parameter key
	token := c.QueryParam("token")
	if token == ""{
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter","Missing required parameter")
	}
	params := make(map[string]string)
	params["string_token"] = token
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		log.Error("Error get email")
		return lib.CustomError(http.StatusBadRequest,"Error get email","Error get email")
	}
	if len(userLogin) < 1 {
		log.Error("No matching token " + token)
		return lib.CustomError(http.StatusBadRequest,"Token not found","Token not found")
	}

	accountData := userLogin[0]
	log.Info("Found account with email " + accountData.UloginEmail)
	
	// Check if token is expired
	dateLayout := "2006-01-02 15:04:05"
	expired, err := time.Parse(dateLayout, *accountData.TokenExpired)
	if err != nil {
		log.Error("Error parsing data")
		return lib.CustomError(http.StatusInternalServerError,err.Error(),"Error parsing data")
	}
	now := time.Now()
	if now.After(expired) {
		log.Error("Token is expired")
		return lib.CustomError(http.StatusInternalServerError,err.Error(),"Token is expired")
	}
	log.Info("Success verify email")
	// Set expired for otp
	date := time.Now().Add(1*time.Minute)
	expiredOTP := date.Format(dateLayout)

	// Generate otp
	otp := lib.EncodeToString(6)

	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["otp_number"] = otp
	params["otp_number_expired"] = expiredOTP
	params["verified_email"] = "1"
	params["last_verified_email"] = time.Now().Format(dateLayout)
	params["string_token"] = ""

	_ ,err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError,err.Error(),"Failed update data")
	}
	
	// Send otp
	//
	//
	//
	//
	//
	
	log.Info("Success send otp")
	
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
func VerifyOtp(c echo.Context) error {
	var err error
	// Get parameter key
	otp := c.FormValue("otp")
	if otp == ""{
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter","Missing required parameter")
	}
	params := make(map[string]string)
	params["otp_number"] = otp
	var userLogin []models.ScUserLogin
	_, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		log.Error("No matching otp " + otp)
		return lib.CustomError(http.StatusBadRequest,"OTP not found","OTP not found")
	}

	accountData := userLogin[0]
	log.Info("Found account with email " + accountData.UloginEmail)
	
	// Check if token is expired
	dateLayout := "2006-01-02 15:04:05"
	expired, err := time.Parse(dateLayout, *accountData.OtpNumberExpired)
	if err != nil {
		log.Error("Error parsing data")
		return lib.CustomError(http.StatusInternalServerError,err.Error(),"Error parsing data")
	}
	now := time.Now()
	if now.After(expired) {
		log.Error("OTP is expired")
		return lib.CustomError(http.StatusInternalServerError,err.Error(),"OTP is expired")
	}
	log.Info("Success verify OTP")

	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["otp_number"] = ""
	params["ulogin_enabled"] = "1"
	params["verified_mobileno"] = "1"
	params["last_verified_mobileno"] = time.Now().Format(dateLayout)

	_ ,err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError,err.Error(),"Failed update data")
	}
	
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func Login(c echo.Context) error {

	var err error
	var status int
	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter","Missing required parameter")
	}
	password := c.FormValue("password")
	if password == "" {
		log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest,"Missing required parameter","Missing required parameter")
	}

	// Check valid email
	params := make(map[string]string)
	params["ulogin_email"] = email
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		log.Error("Error get email")
		return lib.CustomError(status,"Error get email","Error get email")
	}
	if len(userLogin) < 1 {
		log.Error("Email not registered")
		return lib.CustomError(status,"Email not registered","Email not registered")
	}

	accountData := userLogin[0]
	log.Info(accountData)

	// Check valid password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	if encryptedPassword != accountData.UloginPassword {
		log.Error("Wrong password")
		return lib.CustomError(http.StatusUnauthorized,"Wrong password","Wrong password")
	}

	// Create session key
	date := time.Now()
	uuid := uuid.Must(uuid.NewV4(), nil)
	uuidString := uuid.String()
	sessionKey := base64.StdEncoding.EncodeToString([]byte(uuidString))
	dateLayout := "2006-01-02 15:04:05"
	expired := date.Add(time.Second * time.Duration(config.SessionExpired)).Format(dateLayout)

	// Check previous login
	var loginSession []models.ScLoginSession
	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	status, err = models.GetAllScLoginSession(&loginSession, 0, 0, params, true)
	paramsSession["session_id"] = sessionKey
	paramsSession["login_date"] = time.Now().Format(dateLayout)
	paramsSession["rec_status"] = "1"
	if err == nil && len(loginSession) > 0 {
		log.Info("Active session for previous login, overwrite whit new session")
		if len(loginSession) > 1 {

		}
		paramsSession["login_session_key"] = strconv.FormatUint(loginSession[0].LoginSessionKey, 10)

		status, err = models.UpdateScLoginSession(paramsSession)
		if err != nil {
			log.Error("Error update session")
			return lib.CustomError(status,"Error update session","Login failed")
		}
	} else {
		status, err = models.CreateScLoginSession(paramsSession)
		if err != nil {
			log.Error("Error create session")
			return lib.CustomError(status,"Error create session","Login failed")
		}
	}

	log.Info("Success login")

	var data models.ScLoginSessionInfo
	data.SessionID = sessionKey
	data.Expired = expired
	data.Email = email
	log.Info(data)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	log.Info(response)
	return c.JSON(http.StatusOK, response)
}

func ResendVerification(c echo.Context) error {
	var err error
	var status int
	// Check parameters
	email := c.FormValue("email")
	if email == "" {
		log.Error("Missing required parameter")
		return lib.CustomError(http.StatusBadRequest)
	}

	params := make(map[string]string)
	params["ulogin_email"] = email
	var userLogin []models.ScUserLogin
	status, err = models.GetAllScUserLogin(&userLogin, 0, 0, params, true)
	if err != nil {
		log.Error("Error get email")
		return lib.CustomError(status, err.Error(), "Failed get email")
	}
	if len(userLogin) < 1 {
		log.Error("No matching email " + email)
		return lib.CustomError(http.StatusBadRequest,"Email not registered","Email not registered")
	}

	accountData := userLogin[0]
	log.Info("Found account with email " + accountData.UloginEmail)

	dateLayout := "2006-01-02 15:04:05"
	if accountData.VerifiedEmail != nil && *accountData.VerifiedEmail == 1 {
		date := time.Now().Add(1*time.Minute)
		expiredOTP := date.Format(dateLayout)
	
		// Generate otp
		otp := lib.EncodeToString(6)
	
		params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
		params["otp_number"] = otp
		params["otp_number_expired"] = expiredOTP
		params["verified_email"] = "1"
		params["last_verified_email"] = time.Now().Format(dateLayout)
		params["string_token"] = ""
	
		_ ,err = models.UpdateScUserLogin(params)
		if err != nil {
			log.Error("Error update user data")
			return lib.CustomError(http.StatusInternalServerError,err.Error(),"Failed update data")
		}
		
		// Send otp
		//
		//
		//
		//
		//
		
		log.Info("Success send otp")
	} else {
		// Set expired for token
		date := time.Now().AddDate(0, 0, 1)
		expired := date.Format(dateLayout)

		// Generate verify key
		verifyKeyByte := sha256.Sum256([]byte(email + "_" + expired))
		verifyKey := hex.EncodeToString(verifyKeyByte[:])

		// Update token
		params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
		params["string_token"] = verifyKey 
		params["token_expired"] = expired 

		status, err = models.UpdateScUserLogin(params)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed update token")
		}

		// Send email
		mailer := gomail.NewMessage()
		mailer.SetHeader("From", config.EmailFrom)
		mailer.SetHeader("To", email)
		mailer.SetHeader("Subject", "[MNCduit] Verify your email address")
		mailer.SetBody("text/html", "Hi,<br><br>To complete register MNCduit account, please verify your email:<br><br>https://api.mncduit.id/emailverification?key="+verifyKey+"<br><br>Thank you,<br>MNCduit Team")

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
			return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
		}
		log.Info("Email sent")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func verifyPassword(s string) (length, number, upper, special bool) {
	var letter bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c):
			letter = true
		default:
			//return false, false, false, false
		}
	}
	length = letter && len(s) >= 8

	return
}