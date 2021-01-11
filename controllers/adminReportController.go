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
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func GetTransactionHistory(c echo.Context) error {

	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

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

	params := make(map[string]string)

	productkey := c.QueryParam("product_key")
	if productkey != "" {
		params["t.product_key"] = productkey
	}

	customerkey := c.QueryParam("customer_key")
	if customerkey != "" {
		params["c.customer_key"] = customerkey
	}

	datefrom := c.QueryParam("date_from")
	if datefrom == "" {
		log.Error("date_from parameter tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "date_from parameter tidak boleh kosong", "date_from parameter tidak boleh kosong")
	}

	dateto := c.QueryParam("date_to")
	if dateto == "" {
		log.Error("date_to parameter tidak boleh kosong")
		return lib.CustomError(http.StatusBadRequest, "date_to parameter tidak boleh kosong", "date_to parameter tidak boleh kosong")
	}

	layoutISO := "2006-01-02"

	from, _ := time.Parse(layoutISO, datefrom)
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)

	to, _ := time.Parse(layoutISO, dateto)
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)

	params["dateFrom"] = datefrom
	params["dateTo"] = dateto

	if from.Before(to) {
		params["dateFrom"] = datefrom
		params["dateTo"] = dateto
	}

	if from.After(to) {
		params["dateFrom"] = dateto
		params["dateTo"] = datefrom
	}

	items := []string{"product_name", "nav_date", "full_name", "cif", "sid"}

	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {

			var orderByJoin string
			orderByJoin = "t.nav_date"
			if orderBy == "product_name" {
				orderByJoin = "p.product_name"
			} else if orderBy == "nav_date" {
				orderByJoin = "t.nav_date"
			} else if orderBy == "full_name" {
				orderByJoin = "c.full_name"
			} else if orderBy == "cif" {
				orderByJoin = "c.unit_holder_idno"
			} else if orderBy == "sid" {
				orderByJoin = "c.sid_no"
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
		params["orderBy"] = "t.nav_date"
		params["orderType"] = "ASC"
	}

	paramsLike := make(map[string]string)

	customername := c.QueryParam("customer_name")
	if customername != "" {
		paramsLike["c.full_name"] = customername
	}

	var transaksi []models.TransactionCustomerHistory

	status, err = models.AdminGetTransactionCustomerHistory(&transaksi, limit, offset, params, paramsLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(transaksi) < 1 {
		log.Error("Transaksi not found")
		return lib.CustomError(http.StatusNotFound, "Transaksi not found", "Transaksi not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminCountTransactionCustomerHistory(&countData, params, paramsLike)
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
	response.Data = transaksi

	return c.JSON(http.StatusOK, response)
}
