package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func AdminGetListScUserDept(c echo.Context) error {
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

	params := make(map[string]string)

	items := []string{"role_key", "role_code", "role_name", "role_desc"}

	// Get parameter order_by
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
	}

	params["rec_status"] = "1"

	var scUserDept []models.ScUserDept
	status, err = models.GetAllScUserDept(&scUserDept, limit, offset, params, noLimit)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var responseData []models.ScUserDeptInfo

	for _, dept := range scUserDept {
		var data models.ScUserDeptInfo
		data.UserDeptKey = dept.UserDeptKey
		data.UserDeptCode = dept.UserDeptCode
		data.UserDeptName = dept.UserDeptName
		data.UserDeptDesc = dept.UserDeptDesc

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
func GetListScUserDeptAdmin(c echo.Context) error {

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

	items := []string{"user_dept_key", "parent_dept", "user_dept_code", "user_dept_name", "role_privileges", "branch_name"}

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
		params["orderBy"] = "user_dept_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		params["ud.branch_key"] = branchKey
	}
	rolePrivilages := c.QueryParam("role_privileges")
	if rolePrivilages != "" {
		params["ud.role_privileges"] = rolePrivilages
	}

	var userDept []models.ListUserDeptAdmin

	status, err = models.AdminGetListScUserDept(&userDept, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(userDept) < 1 {
		log.Error("User Dept not found")
		return lib.CustomError(http.StatusNotFound, "User Dept not found", "User Dept not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListScUserDept(&countData, params, searchLike)
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
	response.Data = userDept

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteScUserDept(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("user_dept_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: user_dept_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: user_dept_key", "Missing required parameter: user_dept_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["user_dept_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScUserDept(params)
	if err != nil {
		log.Error("Error delete sc_user_dept")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateUserDept(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	userDeptParent := c.FormValue("user_dept_parent")
	if userDeptParent != "" {
		n, err := strconv.ParseUint(userDeptParent, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: user_dept_parent")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: user_dept_parent", "Wrong input for parameter: user_dept_parent")
		}
		params["user_dept_parent"] = userDeptParent
	}

	userDeptCode := c.FormValue("user_dept_code")
	if userDeptCode == "" {
		log.Error("Missing required parameter: user_dept_code")
		return lib.CustomError(http.StatusBadRequest, "user_dept_code can not be blank", "user_dept_code can not be blank")
	} else {
		//validate unique user_dept_code
		var countData models.CountData
		status, err = models.CountScUserDeptValidateUnique(&countData, "user_dept_code", userDeptCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: user_dept_code")
			return lib.CustomError(http.StatusBadRequest, "user_dept_code already used", "user_dept_code already used")
		}
		params["user_dept_code"] = userDeptCode
	}

	userDeptName := c.FormValue("user_dept_name")
	if userDeptName == "" {
		log.Error("Missing required parameter: user_dept_name")
		return lib.CustomError(http.StatusBadRequest, "user_dept_name can not be blank", "user_dept_name can not be blank")
	} else {
		//validate unique user_dept_name
		var countData models.CountData
		status, err = models.CountScUserDeptValidateUnique(&countData, "user_dept_name", userDeptName, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: user_dept_name")
			return lib.CustomError(http.StatusBadRequest, "user_dept_name already used", "user_dept_name already used")
		}
		params["user_dept_name"] = userDeptName
	}

	userDeptDesc := c.FormValue("user_dept_desc")
	if userDeptDesc != "" {
		params["user_dept_desc"] = userDeptDesc
	}

	userDeptEmailAddress := c.FormValue("user_dept_email_address")
	if userDeptEmailAddress != "" {
		params["user_dept_email_address"] = userDeptEmailAddress
	}

	rolePrivileges := c.FormValue("role_privileges")
	if rolePrivileges == "" {
		log.Error("Missing required parameter: role_privileges")
		return lib.CustomError(http.StatusBadRequest, "role_privileges can not be blank", "role_privileges can not be blank")
	} else {
		n, err := strconv.ParseUint(rolePrivileges, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: role_privileges")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: role_privileges", "Wrong input for parameter: role_privileges")
		}
		params["role_privileges"] = rolePrivileges
	}

	branchKey := c.FormValue("branch_key")
	if branchKey != "" {
		n, err := strconv.ParseUint(branchKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
		params["branch_key"] = branchKey
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.CreateScUserDept(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminUpdateUserDept(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	userDeptKey := c.FormValue("user_dept_key")
	if userDeptKey == "" {
		log.Error("Missing required parameter: user_dept_key")
		return lib.CustomError(http.StatusBadRequest, "user_dept_key can not be blank", "user_dept_key can not be blank")
	} else {
		n, err := strconv.ParseUint(userDeptKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: user_dept_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: user_dept_key", "Wrong input for parameter: user_dept_key")
		}
		params["user_dept_key"] = userDeptKey
	}

	userDeptParent := c.FormValue("user_dept_parent")
	if userDeptParent != "" {
		n, err := strconv.ParseUint(userDeptParent, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: user_dept_parent")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: user_dept_parent", "Wrong input for parameter: user_dept_parent")
		}
		params["user_dept_parent"] = userDeptParent
	}

	userDeptCode := c.FormValue("user_dept_code")
	if userDeptCode == "" {
		log.Error("Missing required parameter: user_dept_code")
		return lib.CustomError(http.StatusBadRequest, "user_dept_code can not be blank", "user_dept_code can not be blank")
	} else {
		//validate unique user_dept_code
		var countData models.CountData
		status, err = models.CountScUserDeptValidateUnique(&countData, "user_dept_code", userDeptCode, userDeptKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: user_dept_code")
			return lib.CustomError(http.StatusBadRequest, "user_dept_code already used", "user_dept_code already used")
		}
		params["user_dept_code"] = userDeptCode
	}

	userDeptName := c.FormValue("user_dept_name")
	if userDeptName == "" {
		log.Error("Missing required parameter: user_dept_name")
		return lib.CustomError(http.StatusBadRequest, "user_dept_name can not be blank", "user_dept_name can not be blank")
	} else {
		//validate unique user_dept_name
		var countData models.CountData
		status, err = models.CountScUserDeptValidateUnique(&countData, "user_dept_name", userDeptName, userDeptKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: user_dept_name")
			return lib.CustomError(http.StatusBadRequest, "user_dept_name already used", "user_dept_name already used")
		}
		params["user_dept_name"] = userDeptName
	}

	userDeptDesc := c.FormValue("user_dept_desc")
	if userDeptDesc != "" {
		params["user_dept_desc"] = userDeptDesc
	}

	userDeptEmailAddress := c.FormValue("user_dept_email_address")
	if userDeptEmailAddress != "" {
		params["user_dept_email_address"] = userDeptEmailAddress
	}

	rolePrivileges := c.FormValue("role_privileges")
	if rolePrivileges == "" {
		log.Error("Missing required parameter: role_privileges")
		return lib.CustomError(http.StatusBadRequest, "role_privileges can not be blank", "role_privileges can not be blank")
	} else {
		n, err := strconv.ParseUint(rolePrivileges, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: role_privileges")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: role_privileges", "Wrong input for parameter: role_privileges")
		}
		params["role_privileges"] = rolePrivileges
	}

	branchKey := c.FormValue("branch_key")
	if branchKey != "" {
		n, err := strconv.ParseUint(branchKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
		params["branch_key"] = branchKey
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		n, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateScUserDept(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}
func AdminDetailUserDept(c echo.Context) error {
	var err error

	userDeptKey := c.Param("user_dept_key")
	if userDeptKey == "" {
		log.Error("Missing required parameter: user_dept_key")
		return lib.CustomError(http.StatusBadRequest, "user_dept_key can not be blank", "user_dept_key can not be blank")
	} else {
		n, err := strconv.ParseUint(userDeptKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: user_dept_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: user_dept_key", "Wrong input for parameter: user_dept_key")
		}
	}

	var userDept models.ScUserDept
	_, err = models.GetScUserDept(&userDept, userDeptKey)
	if err != nil {
		log.Error("User Dept not found")
		return lib.CustomError(http.StatusBadRequest, "User Dept not found", "User Dept not found")
	}

	responseData := make(map[string]interface{})
	responseData["menu_key"] = userDept.UserDeptKey
	if userDept.UserDeptParent != nil {
		responseData["user_dept_parent"] = *userDept.UserDeptParent
	} else {
		responseData["user_dept_parent"] = ""
	}
	responseData["user_dept_code"] = userDept.UserDeptCode
	responseData["user_dept_name"] = userDept.UserDeptName
	if userDept.UserDeptDesc != nil {
		responseData["user_dept_desc"] = *userDept.UserDeptDesc
	} else {
		responseData["user_dept_desc"] = ""
	}
	if userDept.UserDeptEmailAddress != nil {
		responseData["user_dept_email_address"] = *userDept.UserDeptEmailAddress
	} else {
		responseData["user_dept_email_address"] = ""
	}
	if userDept.RolePrivileges != nil {
		responseData["role_privileges"] = *userDept.RolePrivileges
	} else {
		responseData["role_privileges"] = ""
	}
	if userDept.BranchKey != nil {
		responseData["branch_key"] = *userDept.BranchKey
	} else {
		responseData["branch_key"] = ""
	}
	if userDept.RecOrder != nil {
		responseData["rec_order"] = *userDept.RecOrder
	} else {
		responseData["rec_order"] = ""
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
