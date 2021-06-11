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

func GetListTrAccount(c echo.Context) error {

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

	items := []string{"key", "product", "full_name", "sid", "cif", "sub_suspend_flag", "red_suspend_flag"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var orderByJoin string
			orderByJoin = ""
			if orderBy == "key" {
				orderByJoin = "a.acc_key"
			} else if orderBy == "product" {
				orderByJoin = "p.product_name_alt"
			} else if orderBy == "full_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "sid" {
				orderByJoin = "c.sid_no"
			} else if orderBy == "cif" {
				orderByJoin = "c.unit_holder_idno"
			} else if orderBy == "sub_suspend_flag" {
				orderByJoin = "a.sub_suspend_flag"
			} else if orderBy == "red_suspend_flag" {
				orderByJoin = "a.red_suspend_flag"
			}

			params["orderBy"] = orderByJoin
			orderType := c.QueryParam("order_type")
			if (orderType == "asc") || (orderType == "ASC") || (orderType == "desc") || (orderType == "DESC") {
				params["orderType"] = orderType
			}
		} else {
			log.Error("Wrong input for parameter order_by")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter order_by", "Wrong input for parameter order_by")
		}
	} else {
		params["orderBy"] = "a.acc_key"
		params["orderType"] = "DESC"
	}

	productKey := c.QueryParam("product_key")
	if productKey != "" {
		params["a.product_key"] = productKey
	}

	var account []models.TrAccountAdmin

	status, err = models.AdminGetAllTrAccount(&account, limit, offset, params, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(account) < 1 {
		log.Error("Account not found")
		return lib.CustomError(http.StatusNotFound, "Account not found", "Account not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetAllTrAccount(&countData, params)
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
	response.Data = account

	return c.JSON(http.StatusOK, response)
}
func AdminGetDetailAccount(c echo.Context) error {
	accKey := c.Param("acc_key")
	if accKey != "" {
		aKey, err := strconv.ParseUint(accKey, 10, 64)
		if err != nil || aKey == 0 {
			log.Error("Wrong input for parameter: acc_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: acc_key", "Wrong input for parameter: acc_key")
		}
	} else {
		log.Error("Missing required parameter: acc_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: acc_key", "Missing required parameter: acc_key")
	}

	var account models.TrAccountAdmin
	_, err := models.AdminGetDetailTrAccount(&account, accKey)
	if err != nil {
		log.Error("Error get data tr_account")
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = account
	return c.JSON(http.StatusOK, response)

}

func AdminUpdateTrAccount(c echo.Context) error {
	productKey := c.FormValue("product_key")
	if productKey != "" {
		pKey, err := strconv.ParseUint(productKey, 10, 64)
		if err != nil || pKey == 0 {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	customerKeyStr := c.FormValue("customer_key")
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err != nil || customerKey == 0 {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	}

	subsuspendflag := c.FormValue("sub_suspend_flag")
	if subsuspendflag == "" {
		log.Error("Missing required parameter: sub_suspend_flag")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: sub_suspend_flag", "Missing required parameter: sub_suspend_flag")
	}

	redsuspendflag := c.FormValue("red_suspend_flag")
	if redsuspendflag == "" {
		log.Error("Missing required parameter: red_suspend_flag")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: red_suspend_flag", "Missing required parameter: red_suspend_flag")
	}

	subsuspendreason := c.FormValue("sub_suspend_reason")
	if subsuspendflag == "1" {
		if subsuspendreason == "" {
			log.Error("Missing required parameter: sub_suspend_reason")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: sub_suspend_reason", "Missing required parameter: sub_suspend_reason")
		}
	}

	redsuspendreason := c.FormValue("red_suspend_reason")
	if redsuspendflag == "1" {
		if redsuspendreason == "" {
			log.Error("Missing required parameter: red_suspend_reason")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: red_suspend_reason", "Missing required parameter: red_suspend_reason")
		}
	}

	//update tr_account
	dateLayout := "2006-01-02 15:04:05"
	paramsTrAccount := make(map[string]string)
	paramsTrAccount["sub_suspend_flag"] = subsuspendflag
	paramsTrAccount["sub_suspend_modified_date"] = time.Now().Format(dateLayout)
	paramsTrAccount["sub_suspend_reason"] = subsuspendreason
	paramsTrAccount["sub_suspend_reference"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsTrAccount["red_suspend_flag"] = redsuspendflag
	paramsTrAccount["red_suspend_modified_date"] = time.Now().Format(dateLayout)
	paramsTrAccount["red_suspend_reason"] = redsuspendreason
	paramsTrAccount["red_suspend_reference"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsTrAccount["rec_modified_date"] = time.Now().Format(dateLayout)
	paramsTrAccount["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	if customerKeyStr == "" {
		_, err := models.UpdateTrAccountUploadSinvest(paramsTrAccount, "product_key", productKey)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		_, err := models.UpdateTrAccountByProductAndCustomer(paramsTrAccount, productKey, customerKeyStr)
		if err != nil {
			log.Println(err.Error())
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func AdminGetCustomerAccount(c echo.Context) error {
	productKey := c.Param("product_key")
	if productKey != "" {
		pKey, err := strconv.ParseUint(productKey, 10, 64)
		if err != nil || pKey == 0 {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	var customer []models.CustomerDropdown
	_, err := models.GetCustomerAccountByProduct(&customer, productKey)
	if err != nil {
		log.Error("Error get data tr_account")
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = customer
	return c.JSON(http.StatusOK, response)

}
