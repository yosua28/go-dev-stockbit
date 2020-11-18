package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func LogoutAdmin(c echo.Context) error {
	var err error

	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	dateLayout := "2006-01-02 15:04:05"
	paramsSession := make(map[string]string)
	paramsSession["user_login_key"] = strIDUserLogin
	paramsSession["logout_date"] = time.Now().Format(dateLayout)
	paramsSession["login_session_key"] = ""

	_, err = models.UpdateScLoginSession(paramsSession)
	if err != nil {
		log.Error("Error update session in logout")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func GetListScUserLoginAdmin(c echo.Context) error {

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var err error
	var status int

	//Get parameter limit
	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err == nil {
			if (limit == 0) || (limit > config.LimitQuery) {
				limit = config.LimitQuery
			}
		} else {
			log.Error("Limit should be number")
			return lib.CustomError(http.StatusBadRequest, "Limit should be number", "Limit should be number")
		}
	} else {
		limit = config.LimitQuery
	}
	// Get parameter page
	pageStr := c.QueryParam("page")
	var page uint64
	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			if page == 0 {
				page = 1
			}
		} else {
			log.Error("Page should be number")
			return lib.CustomError(http.StatusBadRequest, "Page should be number", "Page should be number")
		}
	} else {
		page = 1
	}
	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}

	noLimitStr := c.QueryParam("nolimit")
	var noLimit bool
	if noLimitStr != "" {
		noLimit, err = strconv.ParseBool(noLimitStr)
		if err != nil {
			log.Error("Nolimit parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "Nolimit parameter should be true/false", "Nolimit parameter should be true/false")
		}
	} else {
		noLimit = false
	}

	items := []string{"user_login_key", "ucategory_name", "user_dept_name", "ulogin_name", "ulogin_full_name", "ulogin_email", "role_name", "rec_created_date"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var orderByJoin string
			orderByJoin = "u.user_login_key"
			if orderBy == "user_login_key" {
				orderByJoin = "u.user_login_key"
			}
			if orderBy == "ucategory_name" {
				orderByJoin = "cat.ucategory_name"
			}
			if orderBy == "user_dept_name" {
				orderByJoin = "dept.user_dept_name"
			}
			if orderBy == "ulogin_name" {
				orderByJoin = "u.ulogin_name"
			}
			if orderBy == "ulogin_full_name" {
				orderByJoin = "u.ulogin_full_name"
			}
			if orderBy == "ulogin_email" {
				orderByJoin = "u.ulogin_email"
			}
			if orderBy == "role_name" {
				orderByJoin = "role.role_name"
			}
			if orderBy == "rec_created_date" {
				orderByJoin = "u.rec_created_date"
			}

			params["orderBy"] = orderByJoin
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "u.user_login_key"
		params["orderType"] = "ASC"
	}

	var searchData *string

	search := c.QueryParam("search")
	if search != "" {
		searchData = &search
	}

	rolekey := c.QueryParam("role_key")
	if rolekey != "" {
		rolekeyCek, err := strconv.ParseUint(rolekey, 10, 64)
		if err == nil && rolekeyCek > 0 {
			params["role.role_key"] = rolekey
		} else {
			log.Error("Wrong input for parameter: role_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key", "Missing required parameter: role_key")
		}
	}

	//mapping scUserLogin
	var scUserLogin []models.AdminListScUserLogin
	status, err = models.AdminGetAllScUserLogin(&scUserLogin, limit, offset, params, noLimit, searchData)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminCountDataGetAllScUserlogin(&countData, params, searchData)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) < int(limit) {
			pagination = 1
		} else {
			calc := math.Ceil(float64(countData.CountData) / float64(limit))
			pagination = int(calc)
		}
	} else {
		pagination = 1
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = scUserLogin

	return c.JSON(http.StatusOK, response)
}

func GetDetailScUserLoginAdmin(c echo.Context) error {

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var err error
	var status int

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var scUserLogin models.ScUserLogin
	status, err = models.GetScUserLoginByKey(&scUserLogin, keyStr)
	if err != nil {
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.AdminDetailScUserLogin

	responseData.UserLoginKey = scUserLogin.UserLoginKey

	var scUserCategory models.ScUserCategory
	strUCKey := strconv.FormatUint(scUserLogin.UserCategoryKey, 10)
	status, err = models.GetScUserCategory(&scUserCategory, strUCKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		var ucat models.ScUserCategoryInfo
		ucat.UserCategoryKey = scUserCategory.UserCategoryKey
		ucat.UcategoryCode = scUserCategory.UcategoryCode
		ucat.UcategoryName = scUserCategory.UcategoryName
		ucat.UcategoryDesc = scUserCategory.UcategoryDesc
		responseData.UserCategory = ucat
	}

	if scUserLogin.UserDeptKey != nil {
		var scUserDept models.ScUserDept
		strUDept := strconv.FormatUint(*scUserLogin.UserDeptKey, 10)
		status, err = models.GetScUserDept(&scUserDept, strUDept)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var udept models.ScUserDeptInfo
			udept.UserDeptKey = scUserDept.UserDeptKey
			udept.UserDeptCode = scUserDept.UserDeptCode
			udept.UserDeptName = scUserDept.UserDeptName
			udept.UserDeptDesc = scUserDept.UserDeptDesc
			responseData.UserDept = &udept
		}
	}

	responseData.UloginName = scUserLogin.UloginName
	responseData.UloginFullName = scUserLogin.UloginFullName
	responseData.UloginEmail = scUserLogin.UloginEmail

	if scUserLogin.RoleKey != nil {
		var scRole models.ScRole
		strRoleKey := strconv.FormatUint(*scUserLogin.RoleKey, 10)
		status, err = models.GetScRole(&scRole, strRoleKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var role models.ScRoleInfoLogin
			role.RoleKey = scRole.RoleKey
			role.RoleCode = scRole.RoleCode
			role.RoleName = scRole.RoleName
			role.RoleDesc = scRole.RoleDesc
			responseData.Role = &role
		}
	}

	if scUserLogin.RecStatus == uint8(1) {
		responseData.Enabled = true
	} else {
		responseData.Enabled = false
	}

	if scUserLogin.UloginLocked == uint8(1) {
		responseData.Locked = true
	} else {
		responseData.Locked = false
	}

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	if scUserLogin.RecCreatedDate != nil {
		date, err := time.Parse(layout, *scUserLogin.RecCreatedDate)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.CreatedDate = &oke
		}
	}

	if scUserLogin.RecImage1 != nil && *scUserLogin.RecImage1 != "" {
		responseData.RecImage = config.BaseUrl + "/images/user/" + strconv.FormatUint(scUserLogin.UserLoginKey, 10) + "/profile/" + *scUserLogin.RecImage1
	} else {
		responseData.RecImage = config.BaseUrl + "/images/post/default.png"
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
