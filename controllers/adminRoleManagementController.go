package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetListRoleManagementAdmin(c echo.Context) error {

	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

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

	items := []string{"role_key", "role_category_code", "role_category_name", "role_code", "role_name", "role_desc"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var orderByJoin string
			orderByJoin = "role.role_key"
			if orderBy == "role_key" {
				orderByJoin = "role.role_key"
			}
			if orderBy == "role_category_code" {
				orderByJoin = "cat.role_category_code"
			}
			if orderBy == "role_category_name" {
				orderByJoin = "cat.role_category_name"
			}
			if orderBy == "role_code" {
				orderByJoin = "role.role_code"
			}
			if orderBy == "role_name" {
				orderByJoin = "role.role_name"
			}
			if orderBy == "role_desc" {
				orderByJoin = "role.role_desc"
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
		params["orderBy"] = "role.role_key"
		params["orderType"] = "ASC"
	}

	var searchData *string

	search := c.QueryParam("search")
	if search != "" {
		searchData = &search
	}

	//mapping role management
	var roleManagement []models.AdminRoleManagement
	status, err = models.AdminGetAllRoleManagement(&roleManagement, limit, offset, params, noLimit, searchData)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminCountDataRoleManagement(&countData, params, searchData)
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
	response.Data = roleManagement

	return c.JSON(http.StatusOK, response)
}

func GetListUserByRole(c echo.Context) error {

	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

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

	items := []string{"ulogin_name", "ulogin_full_name", "ulogin_email"}

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
		params["orderBy"] = "ulogin_name"
		params["orderType"] = "ASC"
	}

	var isNew bool
	isNew = true

	roleKey := c.QueryParam("role_key")
	if roleKey != "" {
		sub, err := strconv.ParseUint(roleKey, 10, 64)
		if err == nil && sub > 0 {
			params["role_key"] = roleKey
			isNew = false
		} else {
			log.Error("Wrong input for parameter: role_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key", "Missing required parameter: role_key")
		}
	}

	params["rec_status"] = "1"

	//mapping role management
	var users []models.ScUserLogin
	var countData models.CountData
	var pagination int
	var responseData []models.AdminListScUserLoginRole

	if isNew == false {
		status, err = models.GetAllScUserLogin(&users, limit, offset, params, noLimit)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
			}
		}

		for _, us := range users {
			var data models.AdminListScUserLoginRole
			data.UloginName = us.UloginName
			data.UloginFullName = us.UloginFullName
			data.UloginEmail = us.UloginEmail

			responseData = append(responseData, data)
		}

		if limit > 0 {
			status, err = models.GetCountScUserLogin(&countData, params)
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
		}
	} else {
		pagination = 1
	}

	var response lib.ResponseWithPagination
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Pagination = pagination
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetDetailRoleManagement(c echo.Context) error {
	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var role models.ScRole
	_, err = models.GetScRole(&role, keyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData models.AdminRoleManagementDetail

	responseData.RoleKey = role.RoleKey

	responseData.RoleCode = role.RoleCode
	responseData.RoleName = role.RoleName
	responseData.RoleDesc = role.RoleDesc

	if role.RoleCategoryKey != nil {
		var roleCategory models.ScRoleCategory
		strRoleCategory := strconv.FormatUint(*role.RoleCategoryKey, 10)
		_, err = models.GetScRoleCategory(&roleCategory, strRoleCategory)
		if err == nil {
			var rc models.ScRoleCategoryInfo
			rc.RoleCategoryKey = roleCategory.RoleCategoryKey
			rc.RoleCategoryCode = roleCategory.RoleCategoryCode
			rc.RoleCategoryName = roleCategory.RoleCategoryName
			rc.RoleCategoryDesc = roleCategory.RoleCategoryDesc
			responseData.RoleCategory = &rc
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetDetailMenuRoleManagement(c echo.Context) error {
	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var strRoleKey string

	roleKey := c.QueryParam("role_key")
	if roleKey != "" {
		sub, err := strconv.ParseUint(roleKey, 10, 64)
		if err == nil && sub > 0 {
			strRoleKey = roleKey
		} else {
			log.Error("Wrong input for parameter: role_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_key", "Missing required parameter: role_key")
		}
	} else {
		strRoleKey = ""
	}

	var parentMenu []models.ListMenuRoleManagement
	_, err = models.AdminGetListMenuRole(&parentMenu, strRoleKey, true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusNotFound)
	}

	var childMenu []models.ListMenuRoleManagement
	_, err = models.AdminGetListMenuRole(&childMenu, strRoleKey, false)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusNotFound)
	}

	var responseData []models.ScMenuDetail

	for _, parent := range parentMenu {
		var data models.ScMenuDetail
		data.MenuKey = parent.MenuKey
		data.ModuleName = parent.ModuleName
		data.MenuName = parent.MenuName
		data.MenuDesc = parent.MenuDesc

		var child []models.ScMenuDetailChild
		for _, c := range childMenu {

			if parent.MenuKey == *c.MenuParent {
				var cc models.ScMenuDetailChild
				cc.MenuKey = c.MenuKey
				cc.MenuName = c.MenuName
				cc.MenuDesc = c.MenuDesc
				if c.Checked == "1" {
					cc.IsChecked = true
				} else {
					cc.IsChecked = false
				}
				child = append(child, cc)
			}
		}
		data.ChildMenu = &child

		if len(child) > 0 {
			responseData = append(responseData, data)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetListRoleCategory(c echo.Context) error {

	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

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

	items := []string{"role_category_key", "role_category_code", "role_category_name", "role_category_desc"}

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
		params["orderBy"] = "role_category_key"
		params["orderType"] = "ASC"
	}

	params["rec_status"] = "1"

	//mapping role category
	var roleCategory []models.ScRoleCategory
	_, err = models.GetAllScRoleCategory(&roleCategory, limit, offset, params, noLimit)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	var responseData []models.ScRoleCategoryInfo
	for _, cat := range roleCategory {
		var data models.ScRoleCategoryInfo
		data.RoleCategoryKey = cat.RoleCategoryKey
		data.RoleCategoryCode = cat.RoleCategoryCode
		data.RoleCategoryName = cat.RoleCategoryName
		data.RoleCategoryDesc = cat.RoleCategoryDesc
		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func CreateAdminRoleManagement(c echo.Context) error {
	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	//role_category_key
	rolecategorykey := c.FormValue("role_category_key")
	if rolecategorykey == "" {
		log.Error("Missing required parameter: role_category_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_category_key cann't be blank", "Missing required parameter: role_category_key cann't be blank")
	}
	strrolecategorykey, err := strconv.ParseUint(rolecategorykey, 10, 64)
	if err == nil && strrolecategorykey > 0 {
		params["role_category_key"] = rolecategorykey
	} else {
		log.Error("Wrong input for parameter: role_category_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_category_key", "Missing required parameter: role_category_key")
	}

	//role_code
	rolecode := c.FormValue("role_code")
	if rolecode == "" {
		log.Error("Wrong input for parameter: role_code")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_code", "Missing required parameter: role_code")
	}
	params["role_code"] = rolecode
	paramCode := make(map[string]string)
	paramCode["role_code"] = rolecode
	paramCode["rec_status"] = "1"

	var countDataExisting models.CountData
	status, err := models.AdminGetValidateUniqueMsRole(&countDataExisting, paramCode, nil)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		log.Error("Missing required parameter: role_code already existing, use other role_code")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_code already existing, use other role_code", "Missing required parameter: role_code already existing, use other role_code")
	}

	//role_name
	rolename := c.FormValue("role_name")
	if rolename == "" {
		log.Error("Wrong input for parameter: role_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_name", "Missing required parameter: role_name")
	}
	params["role_name"] = rolename
	paramName := make(map[string]string)
	paramName["role_name"] = rolename
	paramName["rec_status"] = "1"

	status, err = models.AdminGetValidateUniqueMsRole(&countDataExisting, paramName, nil)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		log.Error("Missing required parameter: role_name already existing, use other role_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_name already existing, use other role_name", "Missing required parameter: role_name already existing, use other role_name")
	}

	//role_desc
	roledesc := c.FormValue("role_desc")
	if roledesc != "" {
		params["role_desc"] = roledesc
	}
	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err, lastID := models.CreateScRole(params)
	if err != nil {
		log.Error("Failed create role data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//list menu id
	menuids := c.FormValue("menu_ids")
	var transParamIds []string
	if menuids != "" {
		s := strings.Split(menuids, ",")

		for _, value := range s {
			is := strings.TrimSpace(value)
			if is != "" {
				if _, ok := lib.Find(transParamIds, is); !ok {
					transParamIds = append(transParamIds, is)
				}
			}
		}
	}

	//get endpoint menu
	if len(transParamIds) > 0 {
		var scEndpoint []models.ScEndpoint
		status, err = models.GetScEndpointIn(&scEndpoint, transParamIds, "menu_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
			}
		}

		for _, ep := range scEndpoint {
			paramsEpAuth := make(map[string]string)

			strMenuKey := strconv.FormatUint(*ep.MenuKey, 10)
			strEpKey := strconv.FormatUint(ep.EndpointKey, 10)

			paramsEpAuth["menu_key"] = strMenuKey
			paramsEpAuth["endpoint_key"] = strEpKey
			paramsEpAuth["role_key"] = lastID
			paramsEpAuth["rec_status"] = "1"
			paramsEpAuth["rec_created_date"] = time.Now().Format(dateLayout)
			paramsEpAuth["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			status, err := models.CreateScEndpointAuth(paramsEpAuth)
			if err != nil {
				log.Error("Failed create endpoint auth data: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed input data")
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

func DeleteRoleManagement(c echo.Context) error {
	var err error
	var status int

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	//role_key
	rolekey := c.FormValue("key")
	if rolekey == "" {
		log.Error("Missing required parameter: key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key cann't be blank", "Missing required parameter: key cann't be blank")
	}
	strrolekey, err := strconv.ParseUint(rolekey, 10, 64)
	if err == nil && strrolekey > 0 {
		params["role_key"] = rolekey
	} else {
		log.Error("Wrong input for parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var role models.ScRole
	_, err = models.GetScRole(&role, rolekey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusNotFound)
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	//delete role management
	status, err = models.UpdateScRole(params)
	if err != nil {
		log.Error("Failed delete role management : " + err.Error())
		return lib.CustomError(status, err.Error(), "failed delete role management")
	}

	//delete endpoint auth
	paramsEndpointAuth := make(map[string]string)
	paramsEndpointAuth["rec_status"] = "0"
	paramsEndpointAuth["rec_deleted_date"] = time.Now().Format(dateLayout)
	paramsEndpointAuth["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	status, err = models.UpdateEndpointAuthByField(paramsEndpointAuth, rolekey, "role_key")
	if err != nil {
		log.Error("Failed delete endpoint auth : " + err.Error())
		return lib.CustomError(status, err.Error(), "failed delete endpoint auth")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminRoleManagement(c echo.Context) error {
	var err error

	errorAuth := initAuthHoIt()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	params := make(map[string]string)

	keyStr := c.FormValue("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}
	params["role_key"] = keyStr

	var role models.ScRole
	_, err = models.GetScRole(&role, keyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusNotFound)
	}

	//role_category_key
	rolecategorykey := c.FormValue("role_category_key")
	if rolecategorykey == "" {
		log.Error("Missing required parameter: role_category_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_category_key cann't be blank", "Missing required parameter: role_category_key cann't be blank")
	}
	strrolecategorykey, err := strconv.ParseUint(rolecategorykey, 10, 64)
	if err == nil && strrolecategorykey > 0 {
		params["role_category_key"] = rolecategorykey
	} else {
		log.Error("Wrong input for parameter: role_category_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_category_key", "Missing required parameter: role_category_key")
	}

	//role_code
	rolecode := c.FormValue("role_code")
	if rolecode == "" {
		log.Error("Wrong input for parameter: role_code")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_code", "Missing required parameter: role_code")
	}
	params["role_code"] = rolecode
	paramCode := make(map[string]string)
	paramCode["role_code"] = rolecode
	paramCode["rec_status"] = "1"

	var countDataExisting models.CountData
	status, err := models.AdminGetValidateUniqueMsRole(&countDataExisting, paramCode, &keyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		log.Error("Missing required parameter: role_code already existing, use other role_code")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_code already existing, use other role_code", "Missing required parameter: role_code already existing, use other role_code")
	}

	//role_name
	rolename := c.FormValue("role_name")
	if rolename == "" {
		log.Error("Wrong input for parameter: role_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_name", "Missing required parameter: role_name")
	}
	params["role_name"] = rolename
	paramName := make(map[string]string)
	paramName["role_name"] = rolename
	paramName["rec_status"] = "1"

	status, err = models.AdminGetValidateUniqueMsRole(&countDataExisting, paramName, &keyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if int(countDataExisting.CountData) > 0 {
		log.Error("Missing required parameter: role_name already existing, use other role_name")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: role_name already existing, use other role_name", "Missing required parameter: role_name already existing, use other role_name")
	}

	//role_desc
	roledesc := c.FormValue("role_desc")
	if roledesc != "" {
		params["role_desc"] = roledesc
	}
	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateScRole(params)
	if err != nil {
		log.Error("Failed update role data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed update data")
	}

	//list menu id
	menuids := c.FormValue("menu_ids")
	var transParamIds []string
	if menuids != "" {
		s := strings.Split(menuids, ",")

		for _, value := range s {
			is := strings.TrimSpace(value)
			if is != "" {
				if _, ok := lib.Find(transParamIds, is); !ok {
					transParamIds = append(transParamIds, is)
				}
			}
		}
	}

	//get endpoint menu
	if len(transParamIds) > 0 {
		//add jika ada new endpoint menu
		var scEndpoint []models.ScEndpoint
		status, err = models.AdminGetEndpointNewInUpdateRole(&scEndpoint, keyStr, transParamIds)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
			}
		}

		for _, ep := range scEndpoint {
			paramsEpAuth := make(map[string]string)

			strMenuKey := strconv.FormatUint(*ep.MenuKey, 10)
			strEpKey := strconv.FormatUint(ep.EndpointKey, 10)

			paramsEpAuth["menu_key"] = strMenuKey
			paramsEpAuth["endpoint_key"] = strEpKey
			paramsEpAuth["role_key"] = keyStr
			paramsEpAuth["rec_status"] = "1"
			paramsEpAuth["rec_created_date"] = time.Now().Format(dateLayout)
			paramsEpAuth["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			status, err := models.CreateScEndpointAuth(paramsEpAuth)
			if err != nil {
				log.Error("Failed create endpoint auth data: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed input data")
			}
		}

		//update rec_status = 0 jika yg sudah pernah dipilih, tidak dipilih lagi
		var scEndpointAuth []models.ScEndpointAuth
		status, err = models.AdminGetEndpointAuthUncheckUpdate(&scEndpointAuth, keyStr, transParamIds)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
			}
		}

		for _, epa := range scEndpointAuth {
			paramsEpAuthUpdate := make(map[string]string)

			strEpAuthKey := strconv.FormatUint(epa.EpAuthKey, 10)

			paramsEpAuthUpdate["ep_auth_key"] = strEpAuthKey
			paramsEpAuthUpdate["rec_status"] = "0"
			paramsEpAuthUpdate["rec_deleted_date"] = time.Now().Format(dateLayout)
			paramsEpAuthUpdate["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

			status, err = models.UpdateEndpointAuthByField(paramsEpAuthUpdate, strEpAuthKey, "ep_auth_key")
			if err != nil {
				log.Error("Failed delete endpoint auth : " + err.Error())
			}
		}
	} else {
		//delete endpoint auth all
		paramsEndpointAuth := make(map[string]string)
		paramsEndpointAuth["rec_status"] = "0"
		paramsEndpointAuth["rec_deleted_date"] = time.Now().Format(dateLayout)
		paramsEndpointAuth["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		status, err = models.UpdateEndpointAuthByField(paramsEndpointAuth, keyStr, "role_key")
		if err != nil {
			log.Error("Failed delete endpoint auth : " + err.Error())
			return lib.CustomError(status, err.Error(), "failed delete endpoint auth")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
