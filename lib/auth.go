package lib

import (
	"api/models"
	"api/config"
	"strings"
	"strconv"
	"net/http"
	"fmt"

	log "github.com/sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type CProfile struct {
	UserID      uint64  `json:"user_id"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	Admin       *bool   `json:"admin,omitempty"`
}

var Profile CProfile

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var tokenString string
		request := c.Request()
		authorization := request.Header["Authorization"]
		if authorization != nil {
			if strings.HasPrefix(authorization[0], "Bearer ") == true {
				tokenString = authorization[0][7:]
				log.Info(tokenString)
			}
		}
		token, err := VerifyToken(tokenString)
		if err != nil {
			log.Error(err.Error())
			return CustomError(http.StatusForbidden, err.Error(), "Authentication failed : cannot verified user")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		log.Info(claims)
		if ok && token.Valid {
			accessUuid, ok := claims["uuid"].(string)
			if !ok {
				log.Error("Cannot get uuid")
				return CustomError(http.StatusForbidden, "Cannot get uuid", "Authentication failed : cannot verified user")
			}
			params := make(map[string]string)
			params["session_id"] = accessUuid
			var loginSession []models.ScLoginSession
			_, err := models.GetAllScLoginSession(&loginSession, 0, 0, params, true)
			if err != nil {
				log.Error("Error get email")
				return CustomError(http.StatusForbidden, "Forbidden", "you have to login first")
			}
			if len(loginSession) < 1 {
				log.Error("No matching token " + tokenString)
				return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
			}
		
			paramsUser := make(map[string]string)
			paramsUser["user_login_key"] = strconv.FormatUint(loginSession[0].UserLoginKey, 10)
			var userLogin []models.ScUserLogin
			_, err = models.GetAllScUserLogin(&userLogin, 0, 0, paramsUser, true)
			if err != nil {
				log.Error("Error get email")
				return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
			}
			if len(userLogin) < 1 {
				log.Error("No user login")
				return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
			}

			Profile.UserID = userLogin[0].UserLoginKey
			Profile.Email = userLogin[0].UloginEmail
			Profile.PhoneNumber = *userLogin[0].UloginMobileno

		} else {
			log.Error("Invalid token")
			return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
		}

		return next(c)
	}
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   //Make sure that the token method conform to "SigningMethodHMAC"
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		log.Error("unexpected signing method: %v", token.Header["alg"])
		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	   }
	   return []byte(config.Secret), nil
	})
	if err != nil {
	   return nil, err
	}
	return token, nil
}
