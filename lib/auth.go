package lib

import (
	"api/config"
	"api/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type CProfile struct {
	UserID          uint64  `json:"user_id"`
	Email           string  `json:"email"`
	PhoneNumber     string  `json:"phone_number"`
	RoleKey         uint64  `json:"role_key"`
	RoleCategoryKey uint64  `json:"role_category_key"`
	RecImage1       string  `json:"rec_image1"`
	CustomerKey     *uint64 `json:"customer_key"`
	UserCategoryKey uint64  `json:"user_category_key"`
	RolePrivileges  *uint64 `json:"role_privileges"`
	BranchKey       *uint64 `json:"branch_key"`
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
			_, err := models.GetAllScLoginSession(&loginSession, config.LimitQuery, 0, params, true)
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
			_, err = models.GetAllScUserLogin(&userLogin, config.LimitQuery, 0, paramsUser, true)
			if err != nil {
				log.Error("Error get email")
				return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
			}
			if len(userLogin) < 1 {
				log.Error("No user login")
				return CustomError(http.StatusForbidden, "Forbidden", "You have to login first")
			}

			user := userLogin[0]
			if user.RoleKey != nil && *user.RoleKey > 0 {
				Profile.RoleKey = *user.RoleKey
				paramsRole := make(map[string]string)
				paramsRole["role_key"] = strconv.FormatUint(*user.RoleKey, 10)
				var role []models.ScRole
				_, err = models.GetAllScRole(&role, config.LimitQuery, 0, paramsRole, true)
				if err != nil {
					log.Error(err.Error())
				} else if len(role) > 0 {
					if role[0].RoleCategoryKey != nil && *role[0].RoleCategoryKey > 0 {
						Profile.RoleCategoryKey = *role[0].RoleCategoryKey
					}
				}

				if user.UserDeptKey != nil {
					var dept models.ScUserDept
					strDept := strconv.FormatUint(*user.UserDeptKey, 10)
					_, err = models.GetScUserDept(&dept, strDept)
					if err != nil {
						log.Error(err.Error())
					} else {
						Profile.RolePrivileges = dept.RolePrivileges
						Profile.BranchKey = dept.BranchKey
					}
				}
			}

			Profile.UserID = user.UserLoginKey
			Profile.Email = user.UloginEmail
			Profile.PhoneNumber = *user.UloginMobileno
			Profile.CustomerKey = user.CustomerKey
			Profile.UserCategoryKey = user.UserCategoryKey
			if user.RecImage1 != nil && *user.RecImage1 != "" {
				Profile.RecImage1 = config.BaseUrl + "/user/" + strconv.FormatUint(user.UserLoginKey, 10) + "/profile/" + *user.RecImage1
			} else {
				Profile.RecImage1 = config.BaseUrl + "/user/default.png"
			}

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
