package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"math"
	"net/http"
	"strconv"
	"strings"
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

func AdminCreateMmMailMaster(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	mailMasterType := c.FormValue("mail_master_type")
	if mailMasterType != "" {
		n, err := strconv.ParseUint(mailMasterType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: mail_master_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: mail_master_type", "Wrong input for parameter: mail_master_type")
		}
		params["mail_master_type"] = mailMasterType
	} else {
		log.Error("Missing required parameter: mail_master_type")
		return lib.CustomError(http.StatusBadRequest, "mail_master_type can not be blank", "mail_master_type can not be blank")
	}

	mailMasterCategory := c.FormValue("mail_master_category")
	if mailMasterCategory != "" {
		n, err := strconv.ParseUint(mailMasterCategory, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: mail_master_category")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: mail_master_category", "Wrong input for parameter: mail_master_category")
		}
		params["mail_master_category"] = mailMasterCategory
	}

	mailTemplateName := c.FormValue("mail_template_name")
	if mailTemplateName == "" {
		log.Error("Missing required parameter: mail_template_name")
		return lib.CustomError(http.StatusBadRequest, "mail_template_name can not be blank", "mail_template_name can not be blank")
	} else {
		//validate unique mail_template_name
		var countData models.CountData
		status, err = models.CountMmMailMasterValidateUnique(&countData, "mail_template_name", mailTemplateName, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: mail_template_name")
			return lib.CustomError(http.StatusBadRequest, "mail_template_name already used", "mail_template_name already used")
		}
		params["mail_template_name"] = mailTemplateName
	}

	mailTemplateDesc := c.FormValue("mail_template_desc")
	if mailTemplateDesc != "" {
		params["mail_template_desc"] = mailTemplateDesc
	} else {
		log.Error("Missing required parameter: mail_template_desc")
		return lib.CustomError(http.StatusBadRequest, "mail_template_desc can not be blank", "mail_template_desc can not be blank")
	}

	mailSubject := c.FormValue("mail_subject")
	if mailSubject != "" {
		params["mail_subject"] = mailSubject
	} else {
		log.Error("Missing required parameter: mail_subject")
		return lib.CustomError(http.StatusBadRequest, "mail_subject can not be blank", "mail_subject can not be blank")
	}

	mailTo := c.FormValue("mail_to")
	if mailTo != "" {
		params["mail_to_email_param"] = mailTo
	} else {
		log.Error("Missing required parameter: mail_to_email_param")
		return lib.CustomError(http.StatusBadRequest, "mail_to_email_param can not be blank", "mail_to_email_param can not be blank")
	}

	mailCc := c.FormValue("mail_cc")
	if mailCc != "" {
		params["mail_cc_email_param"] = mailCc
	}

	mailBody := c.FormValue("mail_body")
	if mailBody != "" {
		params["mail_body"] = mailBody
	} else {
		log.Error("Missing required parameter: mail_body")
		return lib.CustomError(http.StatusBadRequest, "mail_body can not be blank", "mail_body can not be blank")
	}

	mailParameter := c.FormValue("mail_parameter")

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["mail_has_attachment"] = "0"
	params["mail_flag_html"] = "1"
	params["rec_status"] = "1"

	status, err, lastID := models.CreateMmMailMaster(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	s := strings.Split(mailParameter, ",")
	var mailParams []string

	for _, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			if _, ok := lib.Find(mailParams, is); !ok {
				mailParams = append(mailParams, is)
			}
		}
	}
	if len(mailParams) > 0 {
		var bindVar []interface{}
		for _, val := range mailParams {

			var row []string
			row = append(row, lastID)
			row = append(row, val)
			row = append(row, "1")
			row = append(row, time.Now().Format(dateLayout))
			row = append(row, strconv.FormatUint(lib.Profile.UserID, 10))
			bindVar = append(bindVar, row)
		}
		status, err = models.CreateMultipleMmMailMasterParamenter(bindVar)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed input data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminUpdateMmMailMaster(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	mailMasterKey := c.FormValue("mail_master_key")
	if mailMasterKey != "" {
		n, err := strconv.ParseUint(mailMasterKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: mail_master_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: mail_master_key", "Wrong input for parameter: mail_master_key")
		}
		params["mail_master_key"] = mailMasterKey
	} else {
		log.Error("Missing required parameter: mail_master_key")
		return lib.CustomError(http.StatusBadRequest, "mail_master_key can not be blank", "mail_master_key can not be blank")
	}

	mailMasterType := c.FormValue("mail_master_type")
	if mailMasterType != "" {
		n, err := strconv.ParseUint(mailMasterType, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: mail_master_type")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: mail_master_type", "Wrong input for parameter: mail_master_type")
		}
		params["mail_master_type"] = mailMasterType
	} else {
		log.Error("Missing required parameter: mail_master_type")
		return lib.CustomError(http.StatusBadRequest, "mail_master_type can not be blank", "mail_master_type can not be blank")
	}

	mailMasterCategory := c.FormValue("mail_master_category")
	if mailMasterCategory != "" {
		n, err := strconv.ParseUint(mailMasterCategory, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: mail_master_category")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: mail_master_category", "Wrong input for parameter: mail_master_category")
		}
		params["mail_master_category"] = mailMasterCategory
	}

	mailTemplateName := c.FormValue("mail_template_name")
	if mailTemplateName == "" {
		log.Error("Missing required parameter: mail_template_name")
		return lib.CustomError(http.StatusBadRequest, "mail_template_name can not be blank", "mail_template_name can not be blank")
	} else {
		//validate unique mail_template_name
		var countData models.CountData
		status, err = models.CountMmMailMasterValidateUnique(&countData, "mail_template_name", mailTemplateName, mailMasterKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: mail_template_name")
			return lib.CustomError(http.StatusBadRequest, "mail_template_name already used", "mail_template_name already used")
		}
		params["mail_template_name"] = mailTemplateName
	}

	mailTemplateDesc := c.FormValue("mail_template_desc")
	if mailTemplateDesc != "" {
		params["mail_template_desc"] = mailTemplateDesc
	} else {
		log.Error("Missing required parameter: mail_template_desc")
		return lib.CustomError(http.StatusBadRequest, "mail_template_desc can not be blank", "mail_template_desc can not be blank")
	}

	mailSubject := c.FormValue("mail_subject")
	if mailSubject != "" {
		params["mail_subject"] = mailSubject
	} else {
		log.Error("Missing required parameter: mail_subject")
		return lib.CustomError(http.StatusBadRequest, "mail_subject can not be blank", "mail_subject can not be blank")
	}

	mailTo := c.FormValue("mail_to")
	if mailTo != "" {
		params["mail_to_email_param"] = mailTo
	} else {
		log.Error("Missing required parameter: mail_to_email_param")
		return lib.CustomError(http.StatusBadRequest, "mail_to_email_param can not be blank", "mail_to_email_param can not be blank")
	}

	mailCc := c.FormValue("mail_cc")
	if mailCc != "" {
		params["mail_cc_email_param"] = mailCc
	}

	mailBody := c.FormValue("mail_body")
	if mailBody != "" {
		params["mail_body"] = mailBody
	} else {
		log.Error("Missing required parameter: mail_body")
		return lib.CustomError(http.StatusBadRequest, "mail_body can not be blank", "mail_body can not be blank")
	}

	mailParameter := c.FormValue("mail_parameter")

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.UpdateMmMailMaster(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	s := strings.Split(mailParameter, ",")
	var mailParams []string

	for _, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			if _, ok := lib.Find(mailParams, is); !ok {
				mailParams = append(mailParams, is)
			}
		}
	}

	if len(mailParams) > 0 {

		paramsPar := make(map[string]string)
		paramsPar["mail_master_key"] = mailMasterKey
		paramsPar["rec_status"] = "1"
		var mmParam []models.MmMailMasterParameter
		status, err = models.GetAllMmMailParametergent(&mmParam, paramsPar)

		var paramKeyDelete []string
		var paramCodeExisting []string

		if err == nil {
			if len(mailParams) > 0 {
				for _, val := range mmParam {
					_, found := lib.Find(mailParams, *val.MailParamCode)
					if found {
						paramCodeExisting = append(paramCodeExisting, *val.MailParamCode)
					} else {
						paramKeyDelete = append(paramKeyDelete, strconv.FormatUint(val.MailParameterKey, 10))
					}
				}
			}
		}

		if len(paramKeyDelete) > 0 {
			paramsParDelete := make(map[string]string)
			paramsParDelete["rec_deleted_date"] = time.Now().Format(dateLayout)
			paramsParDelete["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
			paramsParDelete["rec_status"] = "0"
			status, err = models.UpdateDeleteAllParameter("mail_parameter_key", paramsParDelete, paramKeyDelete)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed update data")
			}
		}

		var bindVar []interface{}
		for _, val := range mailParams {
			_, found := lib.Find(paramCodeExisting, val)
			if !found {
				var row []string
				row = append(row, mailMasterKey)
				row = append(row, val)
				row = append(row, "1")
				row = append(row, time.Now().Format(dateLayout))
				row = append(row, strconv.FormatUint(lib.Profile.UserID, 10))
				bindVar = append(bindVar, row)
			}

		}

		if len(bindVar) > 0 {
			status, err = models.CreateMultipleMmMailMasterParamenter(bindVar)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed input data")
			}
		}
	} else {
		//delete semua param by mail_master_key
		paramsDataDelete := make(map[string]string)
		paramsDataDelete["rec_deleted_date"] = time.Now().Format(dateLayout)
		paramsDataDelete["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		paramsDataDelete["rec_status"] = "0"
		var paramKeyDelete []string
		paramKeyDelete = append(paramKeyDelete, mailMasterKey)
		status, err = models.UpdateDeleteAllParameter("mail_master_key", paramsDataDelete, paramKeyDelete)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed update data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func AdminDetailMmMailMaster(c echo.Context) error {
	var err error

	mailMasterKey := c.Param("mail_master_key")
	if mailMasterKey == "" {
		log.Error("Missing required parameter: mail_master_key")
		return lib.CustomError(http.StatusBadRequest, "mail_master_key can not be blank", "mail_master_key can not be blank")
	} else {
		n, err := strconv.ParseUint(mailMasterKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: mail_master_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: mail_master_key", "Wrong input for parameter: mail_master_key")
		}
	}

	var mail models.MmMailMaster
	_, err = models.GetMmMailMaster(&mail, "mail_master_key", mailMasterKey)
	if err != nil {
		log.Error("Mail Master not found")
		return lib.CustomError(http.StatusBadRequest, "Mail Master not found", "Mail Master not found")
	}

	responseData := make(map[string]interface{})
	responseData["mail_master_key"] = mail.MailMasterKey
	responseData["mail_master_type"] = mail.MailMasterType
	if mail.MailMasterCategory != nil {
		responseData["mail_master_category"] = *mail.MailMasterCategory
	} else {
		responseData["mail_master_category"] = ""
	}
	responseData["mail_template_name"] = mail.MailTemplateName
	if mail.MailTemplateDesc != nil {
		responseData["mail_template_desc"] = *mail.MailTemplateDesc
	} else {
		responseData["mail_template_desc"] = ""
	}
	if mail.MailSubject != nil {
		responseData["mail_subject"] = *mail.MailSubject
	} else {
		responseData["mail_subject"] = ""
	}
	if mail.MailToEmailParam != nil {
		responseData["mail_to"] = *mail.MailToEmailParam
	} else {
		responseData["mail_to"] = ""
	}
	if mail.MailCcEmailParam != nil {
		responseData["mail_cc"] = *mail.MailCcEmailParam
	} else {
		responseData["mail_cc"] = ""
	}
	if mail.MailBody != nil {
		responseData["mail_body"] = *mail.MailBody
	} else {
		responseData["mail_body"] = ""
	}

	paramsPar := make(map[string]string)
	paramsPar["mail_master_key"] = mailMasterKey
	paramsPar["rec_status"] = "1"
	var mmParam []models.MmMailMasterParameter
	_, err = models.GetAllMmMailParametergent(&mmParam, paramsPar)

	var paramCodeExisting []string
	if err == nil {
		if len(mmParam) > 0 {
			for _, val := range mmParam {
				if val.MailParamCode != nil {
					paramCodeExisting = append(paramCodeExisting, *val.MailParamCode)
				}
			}
		}
	}

	responseData["mail_parameter"] = paramCodeExisting

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
