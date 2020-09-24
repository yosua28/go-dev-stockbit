package controllers

import (
	"api/models"
	"api/lib"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

func GetCmsQuiz(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	headerKeyStr := c.QueryParam("header_key")
	headerKey, _ := strconv.ParseUint(headerKeyStr, 10, 64)
	if headerKey == 0 {
		log.Error("Header key should be number")
		return lib.CustomError(http.StatusNotFound,"Header key should be number","Header key should be number")
	}

	var headerDB models.CmsQuizHeader
	status, err = models.GetCmsQuizHeader(&headerDB, headerKeyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	params["orderBy"] = "rec_status"
	params["quiz_header_key"] = headerKeyStr
	params["orderType"] = "ASC"
	params["rec_status"] = "1"
	var questionDB []models.CmsQuizQuestion
	status, err = models.GetAllCmsQuizQuestion(&questionDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(questionDB) < 1 {
		log.Error("Data not found")
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}

	var questionIDs []string
	for _, question := range questionDB {
		questionIDs = append(questionIDs, strconv.FormatUint(question.QuizQuestionKey, 10))
	}
	var optionDB []models.CmsQuizOptions
	status, err = models.GetCmsQuizOptionsIn(&optionDB, questionIDs, "quiz_question_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(optionDB) < 1 {
		log.Error("Data not found")
		return lib.CustomError(http.StatusNotFound, "Data not found", "Data not found")
	}

	optionData := make(map[uint64][]models.CmsQuizOptionsInfo)

	for _, option := range optionDB {
		var data models.CmsQuizOptionsInfo

		data.QuizOptionKey = option.QuizOptionKey
		if option.QuizOptionLabel != nil {
			data.QuizOptionLabel = *option.QuizOptionLabel
		}
		if option.QuizOptionTitle != nil {
			data.QuizOptionTitle = *option.QuizOptionTitle
		}
		if option.QuizOptionScore != nil {
			data.QuizOptionScore = *option.QuizOptionScore
		}
		if option.QuizOptionDefault != nil {
			data.QuizOptionDefault = *option.QuizOptionDefault
		}
		
		optionData[option.QuizQuestionKey] = append(optionData[option.QuizQuestionKey], data)
	}

	var questionData []models.CmsQuizQuestionInfo
	
	for _, question := range questionDB {
		var data models.CmsQuizQuestionInfo
		data.QuizQuestionKey = question.QuizQuestionKey
		if question.QuizTitle != nil {
			data.QuizTitle = *question.QuizTitle
		}
		data.FileImageAllowed = false
		if question.FileImageAllowed == 1 {
			data.FileImageAllowed = true
		}
		data.QuizOptionType = question.QuizOptionType
		if question.QuizOptionDefault != nil {
			data.QuizOptionDefault = *question.QuizOptionDefault
		}
		if opt, ok := optionData[question.QuizQuestionKey]; ok {
			data.Options = &opt
		}
		
		questionData = append(questionData, data)
	}
	

	var responseData models.CmsQuizHeaderData
	
	responseData.QuizHeaderKey = headerDB.QuizHeaderKey
	if headerDB.QuizName != nil {
		responseData.QuizName = *headerDB.QuizName
	}
	if headerDB.QuizDesc != nil {
		responseData.QuizDesc = *headerDB.QuizDesc
	}
	
	responseData.Questions = &questionData

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	
	return c.JSON(http.StatusOK, response)
}