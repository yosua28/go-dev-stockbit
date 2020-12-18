package controllers

import (
	"api/lib"
	"api/models"
	"net/http"
	"strconv"
	"time"

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
	params["orderType"] = "DESC"
	params["orderBy"] = "umessage_receipt_date"
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
	} else {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	read := c.FormValue("read")
	if !(read == "0" || read == "1") {
		log.Error("Missing required parameter: read")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: read", "Missing required parameter: read")
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

	params := make(map[string]string)
	params["umessage_key"] = keyStr
	params["flag_read"] = read

	status, err = models.UpdateScUserMessage(params)
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

func ArchiveMessage(c echo.Context) error {
	var err error
	var status int

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
	} else {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusNotFound, "Missing required parameter: key", "Missing required parameter: key")
	}

	var message models.ScUserMessage
	status, err = models.GetScUserMessage(&message, keyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Notifikaasi tidak ditemukan")
	}

	if message.UmessageRecipientKey != lib.Profile.UserID {
		log.Error("Notifikaasi tidak ditemukan")
		return lib.CustomError(http.StatusBadRequest, "Notifikaasi tidak ditemukan", "Notifikaasi tidak ditemukan")
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)
	params := make(map[string]string)
	params["umessage_key"] = keyStr
	params["flag_archieved"] = "1"
	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	status, err = models.UpdateScUserMessage(params)
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

func ArchiveAllMessage(c echo.Context) error {
	var err error
	var status int

	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)
	params := make(map[string]string)
	params["umessage_recipient_key"] = strIDUserLogin
	params["rec_status"] = "1"
	params["flag_archieved"] = "0"

	var countData models.CountData
	_, err = models.GetCountUserMessage(&countData, params)
	if err != nil {
		log.Error("Notifikaasi tidak ditemukan")
		return lib.CustomError(http.StatusNotFound, "Notifikaasi tidak ditemukan", "Notifikaasi tidak ditemukan")
	}

	if countData.CountData == 0 {
		log.Error("Notifikaasi tidak ditemukan")
		return lib.CustomError(http.StatusNotFound, "Notifikaasi tidak ditemukan", "Notifikaasi tidak ditemukan")
	}

	dateLayout := "2006-01-02 15:04:05"
	paramsUpdate := make(map[string]string)
	paramsUpdate["flag_archieved"] = "1"
	paramsUpdate["rec_modified_by"] = strIDUserLogin
	paramsUpdate["rec_modified_date"] = time.Now().Format(dateLayout)

	status, err = models.UpdateScUserMessageByField(paramsUpdate, "umessage_recipient_key", strIDUserLogin)
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
