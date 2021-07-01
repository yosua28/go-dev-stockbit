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

func AdminGetListScAppConfig(c echo.Context) error {

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

	items := []string{"config_type_code", "app_config_key", "app_config_code", "app_config_name", "app_config_desc", "app_config_datatype", "app_config_value"}

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
		params["orderBy"] = "app_config_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	configTypeKey := c.QueryParam("config_type_key")
	if configTypeKey != "" {
		params["c.config_type_key"] = configTypeKey
	}

	var config []models.ListScAppConfig

	status, err = models.AdminGetListScAppConfig(&config, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(config) < 1 {
		log.Error("Config not found")
		return lib.CustomError(http.StatusNotFound, "Config not found", "Config not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListScAppConfig(&countData, params, searchLike)
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
	response.Data = config

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteScAppConfig(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("app_config_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: app_config_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: app_config_key", "Missing required parameter: app_config_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["app_config_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateScAppConfig(params)
	if err != nil {
		log.Error("Error delete sc_app_config")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateScAppConfig(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	configTypeKey := c.FormValue("config_type_key")
	if configTypeKey == "" {
		log.Error("Missing required parameter: config_type_key")
		return lib.CustomError(http.StatusBadRequest, "config_type_key can not be blank", "config_type_key can not be blank")
	} else {
		n, err := strconv.ParseUint(configTypeKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: config_type_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: config_type_key", "Wrong input for parameter: config_type_key")
		}
		params["config_type_key"] = configTypeKey
	}

	appConfigCode := c.FormValue("app_config_code")
	if appConfigCode == "" {
		log.Error("Missing required parameter: app_config_code")
		return lib.CustomError(http.StatusBadRequest, "app_config_code can not be blank", "app_config_code can not be blank")
	} else {
		//validate unique app_config_code
		var countData models.CountData
		status, err = models.CountScAppConfigValidateUnique(&countData, "app_config_code", appConfigCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: app_config_code")
			return lib.CustomError(http.StatusBadRequest, "app_config_code already used", "app_config_code already used")
		}
		params["app_config_code"] = appConfigCode
	}

	appConfigName := c.FormValue("app_config_name")
	if appConfigName == "" {
		log.Error("Missing required parameter: app_config_name")
		return lib.CustomError(http.StatusBadRequest, "app_config_name can not be blank", "app_config_name can not be blank")
	} else {
		//validate unique app_config_name
		var countData models.CountData
		status, err = models.CountScAppConfigValidateUnique(&countData, "app_config_name", appConfigName, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: app_config_name")
			return lib.CustomError(http.StatusBadRequest, "app_config_name already used", "app_config_name already used")
		}
		params["app_config_name"] = appConfigName
	}

	configDesc := c.FormValue("app_config_desc")
	if configDesc != "" {
		params["app_config_desc"] = configDesc
	}

	dataType := c.FormValue("app_config_datatype")
	if dataType != "" {
		params["app_config_datatype"] = dataType
	}

	value := c.FormValue("app_config_value")
	if value == "" {
		log.Error("Missing required parameter: app_config_value")
		return lib.CustomError(http.StatusBadRequest, "app_config_value can not be blank", "app_config_value can not be blank")
	}
	params["app_config_value"] = value

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

	status, err = models.CreateScAppConfig(params)
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

func AdminUpdateScAppConfig(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	configKey := c.FormValue("app_config_key")
	if configKey == "" {
		log.Error("Missing required parameter: app_config_key")
		return lib.CustomError(http.StatusBadRequest, "app_config_key can not be blank", "app_config_key can not be blank")
	} else {
		n, err := strconv.ParseUint(configKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: app_config_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: app_config_key", "Wrong input for parameter: app_config_key")
		}
		var config models.ScAppConfig
		status, err = models.GetScAppConfig(&config, configKey)
		if err != nil {
			log.Error("Config not found")
			return lib.CustomError(http.StatusBadRequest, "Config not found", "Config not found")
		}
		params["app_config_key"] = configKey
	}

	configTypeKey := c.FormValue("config_type_key")
	if configTypeKey == "" {
		log.Error("Missing required parameter: config_type_key")
		return lib.CustomError(http.StatusBadRequest, "config_type_key can not be blank", "config_type_key can not be blank")
	} else {
		n, err := strconv.ParseUint(configTypeKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: config_type_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: config_type_key", "Wrong input for parameter: config_type_key")
		}
		params["config_type_key"] = configTypeKey
	}

	appConfigCode := c.FormValue("app_config_code")
	if appConfigCode == "" {
		log.Error("Missing required parameter: app_config_code")
		return lib.CustomError(http.StatusBadRequest, "app_config_code can not be blank", "app_config_code can not be blank")
	} else {
		//validate unique app_config_code
		var countData models.CountData
		status, err = models.CountScAppConfigValidateUnique(&countData, "app_config_code", appConfigCode, configKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: app_config_code")
			return lib.CustomError(http.StatusBadRequest, "app_config_code already used", "app_config_code already used")
		}
		params["app_config_code"] = appConfigCode
	}

	appConfigName := c.FormValue("app_config_name")
	if appConfigName == "" {
		log.Error("Missing required parameter: app_config_name")
		return lib.CustomError(http.StatusBadRequest, "app_config_name can not be blank", "app_config_name can not be blank")
	} else {
		//validate unique app_config_name
		var countData models.CountData
		status, err = models.CountScAppConfigValidateUnique(&countData, "app_config_name", appConfigName, configKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: app_config_name")
			return lib.CustomError(http.StatusBadRequest, "app_config_name already used", "app_config_name already used")
		}
		params["app_config_name"] = appConfigName
	}

	configDesc := c.FormValue("app_config_desc")
	if configDesc != "" {
		params["app_config_desc"] = configDesc
	}

	dataType := c.FormValue("app_config_datatype")
	if dataType != "" {
		params["app_config_datatype"] = dataType
	}

	value := c.FormValue("app_config_value")
	if value == "" {
		log.Error("Missing required parameter: app_config_value")
		return lib.CustomError(http.StatusBadRequest, "app_config_value can not be blank", "app_config_value can not be blank")
	}
	params["app_config_value"] = value

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
	params["rec_status"] = "1"

	status, err = models.UpdateScAppConfig(params)
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

func AdminDetailScAppConfig(c echo.Context) error {
	var err error

	configKey := c.Param("app_config_key")
	if configKey == "" {
		log.Error("Missing required parameter: app_config_key")
		return lib.CustomError(http.StatusBadRequest, "app_config_key can not be blank", "app_config_key can not be blank")
	} else {
		n, err := strconv.ParseUint(configKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: app_config_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: app_config_key", "Wrong input for parameter: app_config_key")
		}
	}

	var config models.ScAppConfig
	_, err = models.GetScAppConfig(&config, configKey)
	if err != nil {
		log.Error("Config not found")
		return lib.CustomError(http.StatusBadRequest, "Config not found", "Config not found")
	}

	responseData := make(map[string]interface{})
	responseData["app_config_key"] = config.AppConfigKey
	responseData["config_type_key"] = config.ConfigTypeKey
	responseData["app_config_code"] = *config.AppConfigCode
	responseData["app_config_name"] = *config.AppConfigName
	if config.AppConfigDesc != nil {
		responseData["app_config_desc"] = *config.AppConfigDesc
	} else {
		responseData["app_config_desc"] = ""
	}
	if config.AppConfigDatatype != nil {
		responseData["app_config_datatype"] = *config.AppConfigDatatype
	} else {
		responseData["app_config_datatype"] = ""
	}
	responseData["app_config_value"] = *config.AppConfigValue
	if config.RecOrder != nil {
		responseData["rec_order"] = *config.RecOrder
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
