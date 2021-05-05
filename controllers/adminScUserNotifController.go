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

func GetListUserNotif(c echo.Context) error {

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

	items := []string{"notif_hdr_key", "notif_category", "notif_date_sent", "umessage_subject", "umessage_body", "alert_notif_type"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var orderByJoin string
			orderByJoin = "s.notif_hdr_key"
			if orderBy == "notif_hdr_key" {
				orderByJoin = "s.notif_hdr_key"
			} else if orderBy == "notif_category" {
				orderByJoin = "cat.lkp_name"
			} else if orderBy == "notif_date_sent" {
				orderByJoin = "s.notif_date_sent"
			} else if orderBy == "umessage_subject" {
				orderByJoin = "s.umessage_subject"
			} else if orderBy == "umessage_body" {
				orderByJoin = "s.umessage_body"
			} else if orderBy == "alert_notif_type" {
				orderByJoin = "ty.lkp_name"
			}

			params["orderBy"] = orderByJoin
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
		params["orderBy"] = "s.notif_hdr_key"
		params["orderType"] = "DESC"
	}

	notifCategory := c.QueryParam("notif_category")
	if notifCategory != "" {
		params["s.notif_category"] = notifCategory
	}

	paramLike := ""

	filterData := c.QueryParam("filter_data")
	if filterData != "" {
		paramLike = filterData
	}

	var notifList []models.UserNotifField

	status, err = models.AdminGetAllUserNotif(&notifList, limit, offset, params, paramLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(notifList) < 1 {
		log.Error("Notif User not found")
		return lib.CustomError(http.StatusNotFound, "Notif User not found", "Notif User not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetAllUserNotif(&countData, params, paramLike)
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
	response.Data = notifList

	return c.JSON(http.StatusOK, response)
}

func CreateAdminScUserNotif(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	//notif_category
	notifcategory := c.FormValue("notif_category")
	if notifcategory != "" {
		strnotifcategory, err := strconv.ParseUint(notifcategory, 10, 64)
		if err == nil && strnotifcategory > 0 {
			params["notif_category"] = notifcategory
		} else {
			log.Error("Wrong input for parameter: notif_category")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_category", "Missing required parameter: notif_category")
		}
	} else {
		log.Error("Missing required parameter: notif_category cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_category cann't be blank", "Missing required parameter: notif_category cann't be blank")
	}

	//notif_date_sent
	notifdatesent := c.FormValue("notif_date_sent")
	if notifdatesent != "" {
		params["notif_date_sent"] = notifdatesent
	} else {
		log.Error("Missing required parameter: notif_date_sent cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_date_sent cann't be blank", "Missing required parameter: notif_date_sent cann't be blank")
	}

	//umessage_subject
	umessagesubject := c.FormValue("umessage_subject")
	if umessagesubject != "" {
		params["umessage_subject"] = umessagesubject
	} else {
		log.Error("Missing required parameter: umessage_subject cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: umessage_subject cann't be blank", "Missing required parameter: umessage_subject cann't be blank")
	}

	//umessage_body
	umessagebody := c.FormValue("umessage_body")
	if umessagebody != "" {
		params["umessage_body"] = umessagebody
	} else {
		log.Error("Missing required parameter: umessage_body cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: umessage_body cann't be blank", "Missing required parameter: umessage_body cann't be blank")
	}

	//alert_notif_type
	alertnotiftype := c.FormValue("alert_notif_type")
	if alertnotiftype != "" {
		stralertnotiftype, err := strconv.ParseUint(alertnotiftype, 10, 64)
		if err == nil && stralertnotiftype > 0 {
			params["alert_notif_type"] = alertnotiftype
		} else {
			log.Error("Wrong input for parameter: alert_notif_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: alert_notif_type", "Missing required parameter: alert_notif_type")
		}
	} else {
		log.Error("Missing required parameter: alert_notif_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: alert_notif_type cann't be blank", "Missing required parameter: alert_notif_type cann't be blank")
	}

	params["rec_status"] = "1"

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.CreateScUserNotif(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func DeleteUserNotif(c echo.Context) error {
	var err error
	params := make(map[string]string)

	notifHdrKey := c.FormValue("notif_hdr_key")
	if notifHdrKey == "" {
		log.Error("Missing required parameter: notif_hdr_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_hdr_key", "Missing required parameter: notif_hdr_key")
	}

	keyCek, err := strconv.ParseUint(notifHdrKey, 10, 64)
	if err == nil && keyCek > 0 {
		params["notif_hdr_key"] = notifHdrKey
	} else {
		log.Error("Wrong input for parameter: notif_hdr_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_hdr_key", "Missing required parameter: notif_hdr_key")
	}

	var notif models.ScUserNotif
	status, err := models.GetScUserNotif(&notif, notifHdrKey)
	if err != nil {
		log.Error("Notif not found")
		return lib.CustomError(status)
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateScUserNotif(params)
	if err != nil {
		log.Error("Failed delete data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminScUserNotif(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	//notif_hdr_key
	notifhdrkey := c.FormValue("notif_hdr_key")
	if notifhdrkey != "" {
		strnotifhdrkey, err := strconv.ParseUint(notifhdrkey, 10, 64)
		if err == nil && strnotifhdrkey > 0 {
			params["notif_hdr_key"] = notifhdrkey
		} else {
			log.Error("Wrong input for parameter: notif_hdr_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_hdr_key", "Missing required parameter: notif_hdr_key")
		}
	} else {
		log.Error("Missing required parameter: notif_hdr_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_hdr_key cann't be blank", "Missing required parameter: notif_hdr_key cann't be blank")
	}

	//notif_category
	notifcategory := c.FormValue("notif_category")
	if notifcategory != "" {
		strnotifcategory, err := strconv.ParseUint(notifcategory, 10, 64)
		if err == nil && strnotifcategory > 0 {
			params["notif_category"] = notifcategory
		} else {
			log.Error("Wrong input for parameter: notif_category")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_category", "Missing required parameter: notif_category")
		}
	} else {
		log.Error("Missing required parameter: notif_category cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_category cann't be blank", "Missing required parameter: notif_category cann't be blank")
	}

	//notif_date_sent
	notifdatesent := c.FormValue("notif_date_sent")
	if notifdatesent != "" {
		params["notif_date_sent"] = notifdatesent
	} else {
		log.Error("Missing required parameter: notif_date_sent cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: notif_date_sent cann't be blank", "Missing required parameter: notif_date_sent cann't be blank")
	}

	//umessage_subject
	umessagesubject := c.FormValue("umessage_subject")
	if umessagesubject != "" {
		params["umessage_subject"] = umessagesubject
	} else {
		log.Error("Missing required parameter: umessage_subject cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: umessage_subject cann't be blank", "Missing required parameter: umessage_subject cann't be blank")
	}

	//umessage_body
	umessagebody := c.FormValue("umessage_body")
	if umessagebody != "" {
		params["umessage_body"] = umessagebody
	} else {
		log.Error("Missing required parameter: umessage_body cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: umessage_body cann't be blank", "Missing required parameter: umessage_body cann't be blank")
	}

	//alert_notif_type
	alertnotiftype := c.FormValue("alert_notif_type")
	if alertnotiftype != "" {
		stralertnotiftype, err := strconv.ParseUint(alertnotiftype, 10, 64)
		if err == nil && stralertnotiftype > 0 {
			params["alert_notif_type"] = alertnotiftype
		} else {
			log.Error("Wrong input for parameter: alert_notif_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: alert_notif_type", "Missing required parameter: alert_notif_type")
		}
	} else {
		log.Error("Missing required parameter: alert_notif_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: alert_notif_type cann't be blank", "Missing required parameter: alert_notif_type cann't be blank")
	}

	params["rec_status"] = "1"

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateScUserNotif(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func DetailUserNotif(c echo.Context) error {
	var err error
	var status int

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var notif models.UserNotifField
	status, err = models.AdminGetDetailUserNotif(&notif, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = notif

	return c.JSON(http.StatusOK, response)
}
