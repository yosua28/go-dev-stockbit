package controllers

import (
	"api/lib"
	"api/models"
	"api/config"
	"net/http"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	_ "encoding/base64"
	"time"
	"unicode"
	"strconv"
	"strings"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/badoux/checkmail"
	log "github.com/sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	jwt "github.com/dgrijalva/jwt-go"
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
		return lib.CustomError(http.StatusBadRequest, "Email "+email+" already registered", "Email kamu sudah terdaftar \n Silahkan masukkan email lainnya \n Atau hubungi Customer Service untuk informasi lebih lanjut")
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
	params["ulogin_full_name"] = email 
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
		return lib.CustomError(http.StatusInternalServerError, "Token is expired", "Token is expired")
	}
	log.Info("Success verify email")
	// Set expired for otp
	date := time.Now().Add(1*time.Minute)
	expiredOTP := date.Format(dateLayout)

	// Send otp
	otp, err := sendOTP("0", *accountData.UloginMobileno)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, "Failed send otp", "Failed send otp")
	}
	if otp == "" {
		log.Error("Failed send otp")
		return lib.CustomError(http.StatusInternalServerError, "Failed send otp", "Failed send otp")
	}

	params["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	params["otp_number"] = otp
	params["otp_number_expired"] = expiredOTP
	params["verified_email"] = "1"
	params["last_verified_email"] = time.Now().Format(dateLayout)
	//params["string_token"] = ""

	_ ,err = models.UpdateScUserLogin(params)
	if err != nil {
		log.Error("Error update user data")
		return lib.CustomError(http.StatusInternalServerError,err.Error(),"Failed update data")
	}
	
	log.Info("Success send otp")
	
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func VerifyOtp(c echo.Context) error {
	var err error
	var status int
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
		return lib.CustomError(http.StatusInternalServerError,"OTP is expired","OTP is expired")
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

	// Create session key
	uuid := uuid.Must(uuid.NewV4(), nil)
	uuidString := uuid.String()

	atClaims := jwt.MapClaims{}
	paramsRequest := make(map[string]string)
	paramsRequest["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsRequest["orderBy"] = "oa_request_key"
	paramsRequest["orderType"] = "DESC"
	var request []models.OaRequest
	status, err = models.GetAllOaRequest(&request, config.LimitQuery, 0, true, paramsRequest)
	if err != nil {
		log.Error(err.Error)
	} else if len(request) > 0 {
		if request[0].Oastatus != nil && *request[0].Oastatus > 0 {
			var lookup models.GenLookup
			status, err = models.GetGenLookup(&lookup, strconv.FormatUint(*request[0].Oastatus, 10))
			if err != nil {
				log.Error(err.Error())
			} else {
				if lookup.LkpName != nil && *lookup.LkpName != "" {
					atClaims["oa_status"] = *lookup.LkpName
				}
			}
		}
	}
	if accountData.RoleKey != nil && *accountData.RoleKey > 0 {
		atClaims["role_key"] = *accountData.RoleKey
		paramsRole := make(map[string]string)
		paramsRole["role_key"] = strconv.FormatUint(*accountData.RoleKey,10)
		var role []models.ScRole
		_, err = models.GetAllScRole(&role, config.LimitQuery, 0, paramsRole, true)
		if err != nil {
			log.Error(err.Error())
		}else if len(role) > 0 {
			if role[0].RoleCategoryKey != nil && *role[0].RoleCategoryKey > 0 {
				atClaims["role_category_key"] = *role[0].RoleCategoryKey
			}
		}
	}
	atClaims["uuid"] = uuidString
	atClaims["exp"] = time.Now().Add(time.Minute * 50).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.Secret))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusUnauthorized,err.Error(),"Login failed")
	}

	// Check previous login
	var loginSession []models.ScLoginSession
	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	status, err = models.GetAllScLoginSession(&loginSession, 0, 0, params, true)
	paramsSession["session_id"] = uuidString
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
	data.SessionID = token
	log.Info(data)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
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
		return lib.CustomError(http.StatusUnauthorized,"Email not registered","Email not registered")
	}

	accountData := userLogin[0]
	log.Info(accountData)

	if *accountData.VerifiedEmail != 1 || accountData.VerifiedMobileno != 1 {
		log.Error("Email or Mobile number not verified")
		return lib.CustomError(http.StatusUnauthorized,"Email or Mobile number not verified","Email or Mobile number not verified")
	}

	// Check valid password
	encryptedPasswordByte := sha256.Sum256([]byte(password))
	encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])
	if encryptedPassword != accountData.UloginPassword {
		log.Error("Wrong password")
		return lib.CustomError(http.StatusUnauthorized,"Wrong password","Wrong password")
	}

	// Create session key
	uuid := uuid.Must(uuid.NewV4(), nil)
	uuidString := uuid.String()
	
	atClaims := jwt.MapClaims{}
	paramsRequest := make(map[string]string)
	paramsRequest["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	paramsRequest["orderBy"] = "oa_request_key"
	paramsRequest["orderType"] = "DESC"
	var request []models.OaRequest
	status, err = models.GetAllOaRequest(&request, config.LimitQuery, 0, true, paramsRequest)
	if err != nil {
		log.Error(err.Error)
	} else if len(request) > 0 {
		if request[0].Oastatus != nil && *request[0].Oastatus > 0 {
			var lookup models.GenLookup
			status, err = models.GetGenLookup(&lookup, strconv.FormatUint(*request[0].Oastatus, 10))
			if err != nil {
				log.Error(err.Error())
			} else {
				if lookup.LkpName != nil && *lookup.LkpName != "" {
					atClaims["oa_status"] = *lookup.LkpName
				}
			}
		}
	}
	if accountData.RoleKey != nil && *accountData.RoleKey > 0 {
		atClaims["role_key"] = *accountData.RoleKey
		paramsRole := make(map[string]string)
		paramsRole["role_key"] = strconv.FormatUint(*accountData.RoleKey,10)
		var role []models.ScRole
		_, err = models.GetAllScRole(&role, config.LimitQuery, 0, paramsRole, true)
		if err != nil {
			log.Error(err.Error())
		}else if len(role) > 0 {
			if role[0].RoleCategoryKey != nil && *role[0].RoleCategoryKey > 0 {
				atClaims["role_category_key"] = *role[0].RoleCategoryKey
			}
		}
	}
	atClaims["uuid"] = uuidString
	atClaims["exp"] = time.Now().Add(time.Minute * 50).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.Secret))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusUnauthorized,err.Error(),"Login failed")
	}

	// sessionKey := base64.StdEncoding.EncodeToString([]byte(uuidString))
	dateLayout := "2006-01-02 15:04:05"
	// expired := date.Add(time.Second * time.Duration(config.SessionExpired)).Format(dateLayout)

	// Check previous login
	var loginSession []models.ScLoginSession
	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strconv.FormatUint(accountData.UserLoginKey, 10)
	status, err = models.GetAllScLoginSession(&loginSession, 0, 0, params, true)
	paramsSession["session_id"] = uuidString
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
	data.SessionID = token
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
	
		// Send otp
		otp, err := sendOTP("0", *accountData.UloginMobileno)
		if err != nil || otp == "" {
			log.Error(err.Error())
			return lib.CustomError(http.StatusInternalServerError, "Failed send otp", "Failed send otp")
		}
	
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

func GetUserLogin(c echo.Context) error {
	var err error

	var oaRequestDB []models.OaRequest
	params := make(map[string]string)
	params["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["orderBy"] = "oa_request_key"
	params["orderType"] = "DESC"
	_, err = models.GetAllOaRequest(&oaRequestDB, 0, 0, true, params)
	if err != nil {
		log.Error(err.Error())
	}
	var requestKey string
	if len(oaRequestDB) > 0 {
		requestKey = strconv.FormatUint(oaRequestDB[0].OaRequestKey, 10)
	}

	var personalDataDB models.OaPersonalData
	var riskProfileDB models.OaRiskProfile
	if requestKey != "" {
		_, err = models.GetOaPersonalData(&personalDataDB, requestKey, "oa_request_key")
		if err != nil {
			log.Error(err.Error())
		}
		_, err = models.GetOaRiskProfile(&riskProfileDB, requestKey, "oa_request_key")
		if err != nil {
			log.Error(err.Error())
		}
	}

	var bankAccountDB models.MsBankAccount
	if personalDataDB.BankAccountKey != nil && *personalDataDB.BankAccountKey > 0 {
		_, err = models.GetBankAccount(&bankAccountDB, strconv.FormatUint(*personalDataDB.BankAccountKey, 10))
		if err != nil {
			log.Error(err.Error())
		}
	}

	var riskDB models.MsRiskProfile
	if riskProfileDB.RiskProfileKey > 0 {
		_, err = models.GetMsRiskProfile(&riskDB, strconv.FormatUint(riskProfileDB.RiskProfileKey, 10))
		if err != nil {
			log.Error(err.Error())
		}
	}

	var bankDB models.MsBank
	if bankAccountDB.BankKey > 0 {
		_, err = models.GetMsBank(&bankDB, strconv.FormatUint(bankAccountDB.BankKey, 10))
		if err != nil {
			log.Error(err.Error())
		}
	}

	var customerDB models.MsCustomer
	if lib.Profile.CustomerKey != nil && *lib.Profile.CustomerKey > 0 {
		_, err = models.GetMsCustomer(&customerDB, strconv.FormatUint(*lib.Profile.CustomerKey, 10))
		if err != nil {
			log.Error(err.Error())
		}
	}

	var responseData models.UserProfile
	responseData.FullName = personalDataDB.FullName
	if customerDB.SidNo != nil {
		responseData.SID = *customerDB.SidNo
	}
	responseData.Email = personalDataDB.EmailAddress
	responseData.PhoneNumber = personalDataDB.PhoneMobile
	responseData.RiskProfile.RiskProfileKey = riskDB.RiskProfileKey
	responseData.RiskProfile.RiskCode = riskDB.RiskCode
	responseData.RiskProfile.RiskName = riskDB.RiskName
	responseData.RiskProfile.RiskDesc = riskDB.RiskDesc
	if riskProfileDB.ScoreResult != nil {
		responseData.RiskProfile.Score = *riskProfileDB.ScoreResult
	}
	responseData.RecImage1 = lib.Profile.RecImage1
	responseData.BankAcc.BankName = bankDB.BankName
	responseData.BankAcc.AccountNo = bankAccountDB.AccountNo
	responseData.BankAcc.AccountHolderName = bankAccountDB.AccountHolderName
	responseData.BankAcc.BranchName = bankAccountDB.BranchName

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
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

func sendOTP(gateway, phone string) (string, error){
	curlParam := make(map[string]string)
	curlParam["gateway"] = gateway
	curlParam["msisdn"] = phone
	jsonString, err := json.Marshal(curlParam)
	payload := strings.NewReader(string(jsonString))
	
	req, err := http.NewRequest("POST", config.CitcallUrl, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Apikey 7f837aea98ceea9efcd33ca1d435c9cf")
		
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var sec map[string]interface{}
	if err = json.Unmarshal(body, &sec); err != nil {
	    return "", err
	}
	log.Info(string(body))
	var otp string
	if sec["rc"].(float64) == 0 {
		token := sec["token"].(string)
		otp =  token[len(token)-4:]
	}
	return otp, nil
}