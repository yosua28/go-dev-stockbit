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

func AdminGetListMmMailMaster(c echo.Context) error {

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

	items := []string{"mail_master_type", "mail_master_category", "template_name", "description"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var ord string
			if orderBy == "mail_master_type" {
				ord = "ty.lkp_name"
			} else if orderBy == "mail_master_category" {
				ord = "ct.lkp_name"
			} else if orderBy == "template_name" {
				ord = "m.mail_template_name"
			} else if orderBy == "description" {
				ord = "m.mail_template_desc"
			}
			params["orderBy"] = ord
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "m.mail_master_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")
	mailMasterType := c.QueryParam("mail_master_type")
	if mailMasterType != "" {
		params["m.mail_master_type"] = mailMasterType
	}

	var mail []models.ListMmMailMaster

	status, err = models.AdminGetListMmMailMaster(&mail, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(mail) < 1 {
		log.Error("Mail Master not found")
		return lib.CustomError(http.StatusNotFound, "Mail Master not found", "Mail Master not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetMmMailMaster(&countData, params, searchLike)
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
	response.Data = mail

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMmMailMaster(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("mail_master_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: mail_master_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: mail_master_key", "Missing required parameter: mail_master_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["mail_master_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMmMailMaster(params)
	if err != nil {
		log.Error("Error delete mm_mail_master")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
