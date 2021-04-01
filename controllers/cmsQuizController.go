package controllers

import (
	"api/lib"
	"api/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func GetCmsQuiz(c echo.Context) error {
	var err error
	var status int

	typeKeyStr := c.QueryParam("type_key")
	typeKey, _ := strconv.ParseUint(typeKeyStr, 10, 64)
	if typeKey == 0 {
		log.Error("Type should be number")
		return lib.CustomError(http.StatusBadRequest, "Type should be number", "Type should be number")
	}
	params := make(map[string]string)
	params["rec_status"] = "1"
	params["quiz_type_key"] = typeKeyStr
	params["orderBy"] = "quiz_header_key"
	params["orderType"] = "DESC"

	var headerDB []models.CmsQuizHeader
	status, err = models.GetAllCmsQuizHeader(&headerDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(headerDB) < 1 {
		log.Error("Quiz not found")
		return lib.CustomError(http.StatusNotFound, "Quiz not found", "Quiz not found")
	}
	header := headerDB[0]
	params = make(map[string]string)
	params["orderBy"] = "rec_order"
	params["quiz_header_key"] = strconv.FormatUint(header.QuizHeaderKey, 10)
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

	responseData.QuizHeaderKey = header.QuizHeaderKey
	if header.QuizName != nil {
		responseData.QuizName = *header.QuizName
	}
	if header.QuizDesc != nil {
		responseData.QuizDesc = *header.QuizDesc
	}

	responseData.Questions = &questionData

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func PostQuizAnswer(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	requestKey := m["request_key"].(string)
	fmt.Println(requestKey)
	data := m["quiz"].([]interface{})
	var bindVar []interface{}
	var score uint64 = 0
	for _, val := range data {

		var row []string
		valueMap := val.(map[string]interface{})
		row = append(row, requestKey)
		row = append(row, valueMap["question_key"].(string))
		row = append(row, valueMap["option_key"].(string))
		row = append(row, valueMap["score"].(string))
		row = append(row, "1")
		s, err := strconv.ParseUint(valueMap["score"].(string), 10, 64)
		if err == nil {
			score += s
		}
		bindVar = append(bindVar, row)
	}

	var riskProfile models.MsRiskProfile
	scoreStr := strconv.FormatUint(score, 10)
	status, err = models.GetMsRiskProfileScore(&riskProfile, scoreStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data risk profile")
	}

	params := make(map[string]string)
	params["oa_request_key"] = requestKey
	params["risk_profile_key"] = strconv.FormatUint(riskProfile.RiskProfileKey, 10)
	params["score_result"] = scoreStr
	params["rec_status"] = "1"

	status, err = models.CreateOaRiskProfile(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	status, err = models.CreateMultipleOaRiskProfileQuiz(bindVar)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	var responseData models.MsRiskProfileInfo

	responseData.RiskCode = riskProfile.RiskCode
	responseData.RiskName = riskProfile.RiskName
	responseData.RiskDesc = riskProfile.RiskDesc
	responseData.Score = score

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetQuizAnswer(c echo.Context) error {

	var err error
	var status int

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}
	var requestDB []models.OaRequest
	paramRequest := make(map[string]string)
	paramRequest["customer_key"] = strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	paramRequest["orderBy"] = "oa_request_key"
	paramRequest["orderType"] = "DESC"
	_, err = models.GetAllOaRequest(&requestDB, 1, 0, false, paramRequest)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	request := requestDB[0]
	var quizDB []models.OaRiskProfileQuiz
	paramQuiz := make(map[string]string)
	paramQuiz["oa_request_key"] = strconv.FormatUint(request.OaRequestKey, 10)
	paramQuiz["orderBy"] = "oa_request_key"
	paramQuiz["orderType"] = "DESC"
	_, err = models.GetAllOaRiskProfileQuiz(&quizDB, 100, 0, paramQuiz, true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var risk models.OaRiskProfile
	_, err = models.GetOaRiskProfile(&risk, strconv.FormatUint(request.OaRequestKey, 10), "oa_request_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var riskProfile models.MsRiskProfile
	_, err = models.GetMsRiskProfile(&riskProfile, strconv.FormatUint(risk.RiskProfileKey, 10))
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	responseData := make(map[string]interface{})
	responseData["score_result"] = risk.ScoreResult
	responseData["risk_code"] = riskProfile.RiskCode
	responseData["risk_name"] = riskProfile.RiskName
	responseData["risk_desc"] = riskProfile.RiskDesc
	var quizData []interface{}
	for _, q := range quizDB {
		quiz := make(map[string]interface{})
		quiz["question_key"] = q.QuizQuestionKey
		quiz["option_key"] = q.QuizOptionKey
		quiz["option_score"] = q.QuizOptionScore
		quizData = append(quizData, quiz)
	}
	responseData["quiz"] = quizData

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
