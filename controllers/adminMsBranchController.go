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

func GetListBranchDropdown(c echo.Context) error {

	var err error
	var status int

	var msBranch []models.MsBranchDropdown

	status, err = models.GetMsBranchDropdown(&msBranch)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msBranch) < 1 {
		log.Error("Branch not found")
		return lib.CustomError(http.StatusNotFound, "Branch not found", "Branch not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = msBranch

	return c.JSON(http.StatusOK, response)
}

func AdminGetListMsBranch(c echo.Context) error {

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

	items := []string{"participant_name", "branch_code", "branch_name", "branch_category", "city_name"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var ord string
			if orderBy == "participant_name" {
				ord = "p.participant_name"
			} else if orderBy == "branch_code" {
				ord = "b.branch_code"
			} else if orderBy == "branch_name" {
				ord = "b.branch_name"
			} else if orderBy == "branch_category" {
				ord = "cat.lkp_name"
			} else if orderBy == "city_name" {
				ord = "c.city_name"
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
		params["orderBy"] = "b.branch_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")

	var branch []models.ListMsBranch

	status, err = models.AdminGetListMsBranch(&branch, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(branch) < 1 {
		log.Error("Branch not found")
		return lib.CustomError(http.StatusNotFound, "Branch not found", "Branch not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListMsBranch(&countData, params, searchLike)
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
	response.Data = branch

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMsBranch(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("branch_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: branch_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: branch_key", "Missing required parameter: branch_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["branch_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsBranch(params)
	if err != nil {
		log.Error("Error delete ms_branch")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func AdminCreateMsBranch(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	participanKey := c.FormValue("participant_key")
	if participanKey != "" {
		n, err := strconv.ParseUint(participanKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: participant_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: participant_key", "Wrong input for parameter: participant_key")
		}
		params["participant_key"] = participanKey
	}

	branchCode := c.FormValue("branch_code")
	if branchCode == "" {
		log.Error("Missing required parameter: branch_code")
		return lib.CustomError(http.StatusBadRequest, "branch_code can not be blank", "branch_code can not be blank")
	} else {
		//validate unique branch_code
		var countData models.CountData
		status, err = models.CountMsBranchValidateUnique(&countData, "branch_code", branchCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: branch_code")
			return lib.CustomError(http.StatusBadRequest, "branch_code already used", "branch_code already used")
		}
		params["branch_code"] = branchCode
	}

	branchName := c.FormValue("branch_name")
	if branchName == "" {
		log.Error("Missing required parameter: branch_name")
		return lib.CustomError(http.StatusBadRequest, "branch_name can not be blank", "branch_name can not be blank")
	} else {
		//validate unique branch_name
		var countData models.CountData
		status, err = models.CountMsBranchValidateUnique(&countData, "branch_name", branchName, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: branch_name")
			return lib.CustomError(http.StatusBadRequest, "branch_name already used", "branch_name already used")
		}
		params["branch_name"] = branchName
	}

	branchCategory := c.FormValue("branch_category")
	if branchCategory != "" {
		n, err := strconv.ParseUint(branchCategory, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_category")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_category", "Wrong input for parameter: branch_category")
		}
		params["branch_category"] = branchCategory
	}

	cityKey := c.FormValue("city_key")
	if cityKey != "" {
		n, err := strconv.ParseUint(cityKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: city_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: city_key", "Wrong input for parameter: city_key")
		}
		params["city_key"] = cityKey
	}

	branchAddress := c.FormValue("branch_address")
	if branchAddress != "" {
		params["branch_address"] = branchAddress
	}

	branchEstablished := c.FormValue("branch_established")
	if branchEstablished != "" {
		params["branch_established"] = branchEstablished
	}

	branchPicName := c.FormValue("branch_pic_name")
	if branchPicName != "" {
		params["branch_pic_name"] = branchPicName
	}

	branchPicEmail := c.FormValue("branch_pic_email")
	if branchPicEmail != "" {
		params["branch_pic_email"] = branchPicEmail
	}

	branchPicPhoneno := c.FormValue("branch_pic_phoneno")
	if branchPicPhoneno != "" {
		if len(branchPicPhoneno) > 20 {
			log.Error("branch_pic_phoneno must maximal 20 character")
			return lib.CustomError(http.StatusBadRequest, "branch_pic_phoneno must maximal 20 character", "branch_pic_phoneno must maximal 20 character")
		}
		params["branch_pic_phoneno"] = branchPicPhoneno
	}

	branchCostCenter := c.FormValue("branch_cost_center")
	if branchCostCenter != "" {
		if len(branchCostCenter) > 10 {
			log.Error("branch_cost_center must maximal 10 character")
			return lib.CustomError(http.StatusBadRequest, "branch_cost_center must maximal 10 character", "branch_cost_center must maximal 10 character")
		}
		params["branch_cost_center"] = branchCostCenter
	}

	branchProfitCenter := c.FormValue("branch_profit_center")
	if branchProfitCenter != "" {
		if len(branchProfitCenter) > 10 {
			log.Error("branch_profit_center must maximal 10 character")
			return lib.CustomError(http.StatusBadRequest, "branch_profit_center must maximal 10 character", "branch_profit_center must maximal 10 character")
		}
		params["branch_profit_center"] = branchProfitCenter
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

	status, err = models.CreateMsBranch(params)
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

func AdminUpdateMsBranch(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	branchKey := c.FormValue("branch_key")
	if branchKey == "" {
		log.Error("Missing required parameter: branch_key")
		return lib.CustomError(http.StatusBadRequest, "branch_key can not be blank", "branch_key can not be blank")
	} else {
		n, err := strconv.ParseUint(branchKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
		var brnch models.MsBranch
		status, err = models.GetMsBranch(&brnch, branchKey)
		if err != nil {
			log.Error("Config not found")
			return lib.CustomError(http.StatusBadRequest, "Config not found", "Config not found")
		}
		params["branch_key"] = branchKey
	}

	participanKey := c.FormValue("participant_key")
	if participanKey != "" {
		n, err := strconv.ParseUint(participanKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: participant_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: participant_key", "Wrong input for parameter: participant_key")
		}
		params["participant_key"] = participanKey
	}

	branchCode := c.FormValue("branch_code")
	if branchCode == "" {
		log.Error("Missing required parameter: branch_code")
		return lib.CustomError(http.StatusBadRequest, "branch_code can not be blank", "branch_code can not be blank")
	} else {
		//validate unique branch_code
		var countData models.CountData
		status, err = models.CountMsBranchValidateUnique(&countData, "branch_code", branchCode, branchKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: branch_code")
			return lib.CustomError(http.StatusBadRequest, "branch_code already used", "branch_code already used")
		}
		params["branch_code"] = branchCode
	}

	branchName := c.FormValue("branch_name")
	if branchName == "" {
		log.Error("Missing required parameter: branch_name")
		return lib.CustomError(http.StatusBadRequest, "branch_name can not be blank", "branch_name can not be blank")
	} else {
		//validate unique branch_name
		var countData models.CountData
		status, err = models.CountMsBranchValidateUnique(&countData, "branch_name", branchName, branchKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: branch_name")
			return lib.CustomError(http.StatusBadRequest, "branch_name already used", "branch_name already used")
		}
		params["branch_name"] = branchName
	}

	branchCategory := c.FormValue("branch_category")
	if branchCategory != "" {
		n, err := strconv.ParseUint(branchCategory, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_category")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_category", "Wrong input for parameter: branch_category")
		}
		params["branch_category"] = branchCategory
	}

	cityKey := c.FormValue("city_key")
	if cityKey != "" {
		n, err := strconv.ParseUint(cityKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: city_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: city_key", "Wrong input for parameter: city_key")
		}
		params["city_key"] = cityKey
	}

	branchAddress := c.FormValue("branch_address")
	if branchAddress != "" {
		params["branch_address"] = branchAddress
	}

	branchEstablished := c.FormValue("branch_established")
	if branchEstablished != "" {
		params["branch_established"] = branchEstablished
	}

	branchPicName := c.FormValue("branch_pic_name")
	if branchPicName != "" {
		params["branch_pic_name"] = branchPicName
	}

	branchPicEmail := c.FormValue("branch_pic_email")
	if branchPicEmail != "" {
		params["branch_pic_email"] = branchPicEmail
	}

	branchPicPhoneno := c.FormValue("branch_pic_phoneno")
	if branchPicPhoneno != "" {
		if len(branchPicPhoneno) > 20 {
			log.Error("branch_pic_phoneno must maximal 20 character")
			return lib.CustomError(http.StatusBadRequest, "branch_pic_phoneno must maximal 20 character", "branch_pic_phoneno must maximal 20 character")
		}
		params["branch_pic_phoneno"] = branchPicPhoneno
	}

	branchCostCenter := c.FormValue("branch_cost_center")
	if branchCostCenter != "" {
		if len(branchCostCenter) > 10 {
			log.Error("branch_cost_center must maximal 10 character")
			return lib.CustomError(http.StatusBadRequest, "branch_cost_center must maximal 10 character", "branch_cost_center must maximal 10 character")
		}
		params["branch_cost_center"] = branchCostCenter
	}

	branchProfitCenter := c.FormValue("branch_profit_center")
	if branchProfitCenter != "" {
		if len(branchProfitCenter) > 10 {
			log.Error("branch_profit_center must maximal 10 character")
			return lib.CustomError(http.StatusBadRequest, "branch_profit_center must maximal 10 character", "branch_profit_center must maximal 10 character")
		}
		params["branch_profit_center"] = branchProfitCenter
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
	params["rec_status"] = "1"

	status, err = models.UpdateMsBranch(params)
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

func AdminDetailMsBranch(c echo.Context) error {
	var err error

	branchKey := c.Param("branch_key")
	if branchKey == "" {
		log.Error("Missing required parameter: branch_key")
		return lib.CustomError(http.StatusBadRequest, "branch_key can not be blank", "branch_key can not be blank")
	} else {
		n, err := strconv.ParseUint(branchKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: branch_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: branch_key", "Wrong input for parameter: branch_key")
		}
	}

	var branch models.MsBranch
	_, err = models.GetMsBranch(&branch, branchKey)
	if err != nil {
		log.Error("Branch not found")
		return lib.CustomError(http.StatusBadRequest, "Branch not found", "Branch not found")
	}

	responseData := make(map[string]interface{})
	responseData["branch_key"] = branch.BranchKey
	if branch.ParticipantKey != nil {
		responseData["participant_key"] = *branch.ParticipantKey
	} else {
		responseData["participant_key"] = ""
	}
	responseData["branch_code"] = branch.BranchCode
	responseData["branch_name"] = branch.BranchName
	if branch.BranchCategory != nil {
		responseData["branch_category"] = *branch.BranchCategory
	} else {
		responseData["branch_category"] = ""
	}
	if branch.CityKey != nil {
		responseData["city_key"] = *branch.CityKey
	} else {
		responseData["city_key"] = ""
	}
	if branch.BranchAddress != nil {
		responseData["branch_address"] = *branch.BranchAddress
	} else {
		responseData["branch_address"] = ""
	}
	if branch.BranchEstablished != nil {
		responseData["branch_established"] = *branch.BranchEstablished
	} else {
		responseData["branch_established"] = ""
	}
	if branch.BranchPicName != nil {
		responseData["branch_pic_name"] = *branch.BranchPicName
	} else {
		responseData["branch_pic_name"] = ""
	}
	if branch.BranchPicEmail != nil {
		responseData["branch_pic_email"] = *branch.BranchPicEmail
	} else {
		responseData["branch_pic_email"] = ""
	}
	if branch.BranchPicPhoneno != nil {
		responseData["branch_pic_phoneno"] = *branch.BranchPicPhoneno
	} else {
		responseData["branch_pic_phoneno"] = ""
	}
	if branch.BranchCostCenter != nil {
		responseData["branch_cost_center"] = *branch.BranchCostCenter
	} else {
		responseData["branch_cost_center"] = ""
	}
	if branch.BranchProfitCenter != nil {
		responseData["branch_cost_center"] = *branch.BranchProfitCenter
	} else {
		responseData["branch_cost_center"] = ""
	}
	if branch.RecOrder != nil {
		responseData["rec_order"] = *branch.RecOrder
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
