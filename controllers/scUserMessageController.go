package controllers

import (
	"api/lib"
	"api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetMessageList(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)
	params["umessage_recipient_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"
	params["flag_archieved"] = "0"
	params["orderType"] = "ASC"
	params["orderBy"] = "flag_read"
	var messageDB []models.ScUserMessage
	status, err = models.GetAllScUserMessage(&messageDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(messageDB) < 1 {
		log.Error("Data not found")
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}

	var lookupDB []models.GenLookup
	status, err = models.GetGenLookupIn(&lookupDB, []string{"56", "57"}, "lkp_group_key")
	if err != nil {
		log.Error(err.Error())
	}
	lookupData := make(map[uint64]models.GenLookup)
	if len(lookupDB) > 0 {
		for _, lookup := range lookupDB {
			lookupData[lookup.LookupKey] = lookup
		}
	}
	var responseData []models.ScUserMessageData

	for _, message := range messageDB {
		var data models.ScUserMessageData

		data.UmessageKey = message.UmessageKey
		if message.UmessageType != nil {
			if l, ok := lookupData[*message.UmessageType]; ok {
				data.UmessageType.Value = l.LookupKey
				data.UmessageType.Name = *l.LkpName
			}
		}

		data.UmessageReceiptDate = message.UmessageReceiptDate
		data.FlagRead = message.FlagRead
		data.UmessageBody = message.UmessageBody
		data.UmessageSubject = message.UmessageSubject
		data.UmessageBody = message.UmessageBody
		data.UparentKey = message.UparentKey
		if message.UmessageCategory != nil {
			if l, ok := lookupData[*message.UmessageCategory]; ok {
				data.UmessageCategory.Value = l.LookupKey
				data.UmessageCategory.Name = *l.LkpName
			}
		}
		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetMessageData(c echo.Context) error {
	var err error
	var status int

	keyStr := c.Param("key")
	if keyStr != "" {
		key, err := strconv.ParseUint(keyStr, 10, 64)
		if err != nil {
			log.Error("Wrong value for parameter: key")
			return lib.CustomError(http.StatusBadRequest, "Wrong value for parameter: key", "Wrong value for parameter: key")
		}
		if key == 0 {
			return lib.CustomError(http.StatusNotFound)
		}
	} else {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var message models.ScUserMessage
	status, err = models.GetScUserMessage(&message, keyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Message not found")
	}

	if message.UmessageRecipientKey != lib.Profile.UserID {
		log.Error("Message not found")
		return lib.CustomError(status, "Message not found", "Message not found")
	}

	var lookupDB []models.GenLookup
	status, err = models.GetGenLookupIn(&lookupDB, []string{"56", "57"}, "lkp_group_key")
	if err != nil {
		log.Error(err.Error())
	}
	lookupData := make(map[uint64]models.GenLookup)
	if len(lookupDB) > 0 {
		for _, lookup := range lookupDB {
			lookupData[lookup.LookupKey] = lookup
		}
	}
	var data models.ScUserMessageData

	data.UmessageKey = message.UmessageKey
	if message.UmessageType != nil {
		if l, ok := lookupData[*message.UmessageType]; ok {
			data.UmessageType.Value = l.LookupKey
			data.UmessageType.Name = *l.LkpName
		}
	}

	data.UmessageReceiptDate = message.UmessageReceiptDate
	data.FlagRead = message.FlagRead
	data.UmessageSubject = message.UmessageSubject
	data.UmessageBody = message.UmessageBody
	data.UparentKey = message.UparentKey
	if message.UmessageCategory != nil {
		if l, ok := lookupData[*message.UmessageCategory]; ok {
			data.UmessageCategory.Value = l.LookupKey
			data.UmessageCategory.Name = *l.LkpName
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data

	return c.JSON(http.StatusOK, response)
}

func PatchMessage(c echo.Context) error {
	var err error
	var status int
	params := make(map[string]string)
	where := make(map[string]string)

	action := c.FormValue("action")
	if action != "" {
		value := c.FormValue("value")
		if !(value == "0" || value == "1") {
			log.Error("Missing required parameter: value")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: value", "Missing required parameter: value")
		}
		if action == "read" {
			params["flag_read"] = value
		} else if action == "delete" {
			params["flag_archieved"] = value
		} else {
			return lib.CustomError(http.StatusBadRequest)
		}
	} else {
		log.Error("Missing required parameter: action")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: action", "Missing required parameter: action")
	}

	all := c.FormValue("all")
	if all == "0" {
		keyStr := c.FormValue("key")
		if keyStr != "" {
			key, err := strconv.ParseUint(keyStr, 10, 64)
			if err != nil {
				log.Error("Wrong value for parameter: key")
				return lib.CustomError(http.StatusBadRequest, "Wrong value for parameter: key", "Wrong value for parameter: key")
			}
			if key == 0 {
				return lib.CustomError(http.StatusNotFound)
			}
			var message models.ScUserMessage
			status, err = models.GetScUserMessage(&message, keyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Message not found")
			}
			if message.UmessageRecipientKey != lib.Profile.UserID {
				log.Error("Message not found")
				return lib.CustomError(status, "Message not found", "Message not found")
			}
			where["umessage_key"] = keyStr
		} else {
			log.Error("Missing required parameter: key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
		}
	} else if all == "1" {
		where["umessage_recipient_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	} else {
		log.Error("Wrong input for parameter: all")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: all", "Wrong input for parameter: all")
	}

	status, err = models.UpdateScUserMessage(params, where)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil

	return c.JSON(http.StatusOK, response)
}

func GetCountMessageData(c echo.Context) error {
	var err error
	params := make(map[string]string)
	params["umessage_recipient_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"
	params["flag_archieved"] = "0"
	params["flag_read"] = "0"

	count := 0

	var countData models.CountData
	_, err = models.GetCountUserMessage(&countData, params)
	if err == nil {
		count = countData.CountData
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = count

	return c.JSON(http.StatusOK, response)
}
