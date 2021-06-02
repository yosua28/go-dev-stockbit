package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func GetPromoList(c echo.Context) error {

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
	params["rec_status"] = "1"

	var promoDB []models.TrPromo
	status, err = models.GetAllTrPromoActive(&promoDB, limit, offset, params, noLimit)
	if err != nil {
		return lib.CustomError(status, "Promo tidak ditemukan", "Promo tidak ditemukan")
	}

	var responseData []map[string]interface{}
	for _, promo := range promoDB {
		data := make(map[string]interface{})

		data["promo_key"] = promo.PromoKey
		data["promo_code"] = promo.PromoCode
		data["promo_title"] = promo.PromoTitle
		data["description"] = promo.PromoNotifDescription
		data["tnc"] = promo.PromoTnc

		dir := config.BaseUrl + "/images/promo/"

		if promo.RecImage1 != nil {
			path := dir + *promo.RecImage1
			data["image"] = &path
		}

		dateLayout := "2006-01-02 15:04:05"
		newlayout := "02 Jan 2006"
		month := make(map[string]string)
		month["Jan"] = "Januari"
		month["Feb"] = "Februari"
		month["Mar"] = "Maret"
		month["Apr"] = "April"
		month["May"] = "Mei"
		month["Jun"] = "Juni"
		month["Jul"] = "Juli"
		month["Aug"] = "Agustus"
		month["Sep"] = "September"
		month["Oct"] = "Oktober"
		month["Nov"] = "November"
		month["Dec"] = "Desember"
		t, _ := time.Parse(dateLayout, promo.PromoValidDate2)
		m := t.Month().String()[:3]
		data["valid_date"] = strings.Replace(t.Format(newlayout), m, month[m], 1)

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
