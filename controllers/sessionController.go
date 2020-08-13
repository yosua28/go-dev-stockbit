package controllers

import (
	"api/lib"
	"api/models"
	"api/config"
	"net/http"
	_ "crypto/sha256"
	"crypto/tls"
	_ "encoding/hex"
	_ "time"
	"unicode"

	"github.com/labstack/echo"
	"github.com/badoux/checkmail"
	log "github.com/sirupsen/logrus"
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
		log.Error("Error get " + email)
		return lib.CustomError(status, "Error get email", "Error get email")
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
	// encryptedPasswordByte := sha256.Sum256([]byte(password))
	// encryptedPassword := hex.EncodeToString(encryptedPasswordByte[:])

	// Set expired for token
	// date := time.Now().AddDate(0, 0, 1)
	// dateLayout := "2006-01-02T15:04:05+00:00"
	// expired := date.Format(dateLayout)

	// Generate OTP
	otp := lib.EncodeToString(4) 

	// Input to database
	// params["ulogin_email"] = email 
	// params["ulogin_password"] = encryptedPassword 
	// params["ulogin_mobileno"] = phone 
	// params["user_category_key"] = "1" 

	// status, err = models.CreateScUserLogin(params)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed create user")
	// }

	// Send email
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "[MNCduit] Verify your email address")
	mailer.SetBody("text/html", "your OTP : " + otp + "<br><br>Thank you,<br>MNCduit")

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