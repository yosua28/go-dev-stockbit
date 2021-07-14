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

func GetListAgentDropdown(c echo.Context) error {

	var err error
	var status int

	var msAgent []models.MsAgentDropdown

	status, err = models.GetMsAgentDropdown(&msAgent)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(msAgent) < 1 {
		log.Error("Agent not found")
		return lib.CustomError(http.StatusNotFound, "Agent not found", "Agent not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = msAgent

	return c.JSON(http.StatusOK, response)
}

func GetListAgentLastBranch(c echo.Context) error {

	var err error
	var status int

	keyStr := c.Param("branch_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var agent []models.MsAgentLastBranch

	status, err = models.GetMsAgentLastBranch(&agent, keyStr)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(agent) < 1 {
		log.Error("Agent not found")
		return lib.CustomError(http.StatusNotFound, "Agent not found", "Agent not found")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = agent

	return c.JSON(http.StatusOK, response)
}

func AdminGetListMsAgent(c echo.Context) error {

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

	items := []string{"branch_name", "agent_id", "agent_code", "agent_name", "agent_category", "agent_channel"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var ord string
			if orderBy == "branch_name" {
				ord = "b.branch_name"
			} else if orderBy == "agent_id" {
				ord = "a.agent_id"
			} else if orderBy == "agent_code" {
				ord = "a.agent_code"
			} else if orderBy == "agent_name" {
				ord = "a.agent_name"
			} else if orderBy == "agent_category" {
				ord = "cat.lkp_name"
			} else if orderBy == "agent_channel" {
				ord = "cha.lkp_name"
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
		params["orderBy"] = "a.agent_key"
		params["orderType"] = "DESC"
	}

	searchLike := c.QueryParam("search_like")
	branchKey := c.QueryParam("branch_key")
	if branchKey != "" {
		params["b.branch_key"] = branchKey
	}

	var agent []models.ListAgentAdmin

	status, err = models.AdminGetListAgent(&agent, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(agent) < 1 {
		log.Error("Agent not found")
		return lib.CustomError(http.StatusNotFound, "Agent not found", "Agent not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetListAgent(&countData, params, searchLike)
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
	response.Data = agent

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMsAgent(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("agent_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: agent_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: agent_key", "Missing required parameter: agent_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["agent_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsAgent(params)
	if err != nil {
		log.Error("Error delete ms_branch")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	//delete ms_agent_branch
	paramsAgentBranch := make(map[string]string)
	paramsAgentBranch["rec_status"] = "0"
	paramsAgentBranch["rec_deleted_date"] = time.Now().Format(dateLayout)
	paramsAgentBranch["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateDeleteBranchAgent(paramsAgentBranch, "agent_key", keyStr)
	if err != nil {
		log.Error("Error delete ms_agent_branch")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
