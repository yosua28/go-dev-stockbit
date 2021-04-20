package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func GetListPromo(c echo.Context) error {

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

	items := []string{"promo_key", "promo_code", "promo_title", "promo_category", "promo_nominal", "promo_max_nominal", "promo_valid_date1", "promo_notif_start"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
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
		params["orderBy"] = "promo_key"
		params["orderType"] = "DESC"
	}

	startvalid := c.QueryParam("start_valid")

	endvalid := c.QueryParam("end_valid")

	layoutISO := "2006-01-02"

	if startvalid != "" && endvalid != "" {
		from, _ := time.Parse(layoutISO, startvalid)
		from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)

		to, _ := time.Parse(layoutISO, endvalid)
		to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)

		params["start_valid"] = startvalid
		params["end_valid"] = endvalid

		if from.Before(to) {
			params["start_valid"] = startvalid
			params["end_valid"] = endvalid
		}

		if from.After(to) {
			params["start_valid"] = endvalid
			params["end_valid"] = startvalid
		}
	}

	promocategory := c.QueryParam("promo_category")
	if promocategory != "" {
		params["promo_category"] = promocategory
	}

	var promoList []models.TrPromoData

	status, err = models.AdminGetAllTrPromo(&promoList, limit, offset, params, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(promoList) < 1 {
		log.Error("Promo not found")
		return lib.CustomError(http.StatusNotFound, "Promo not found", "Promo not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminGetCountTrPromo(&countData, params)
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
	response.Data = promoList

	return c.JSON(http.StatusOK, response)
}

func CreateAdminTrPromo(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	//promo_code
	promocode := c.FormValue("promo_code")
	if promocode == "" {
		log.Error("Missing required parameter: promo_code cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_code cann't be blank", "Missing required parameter: promo_code cann't be blank")
	}

	//cek duplikat promo code
	var promo models.TrPromo
	status, err = models.GetTrPromo(&promo, "promo_code", promocode)
	if err == nil {
		log.Error("Missing required parameter: promo_code already used")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_code already used", "Missing required parameter: promo_code already used")
	}
	log.Println(promo)
	params["promo_code"] = promocode

	//promo_title
	promootitle := c.FormValue("promo_title")
	if promootitle == "" {
		log.Error("Missing required parameter: promo_title cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_title cann't be blank", "Missing required parameter: promo_title cann't be blank")
	}
	params["promo_title"] = promootitle

	//promo_category
	promocategory := c.FormValue("promo_category")
	if promocategory != "" {
		strpromocategory, err := strconv.ParseUint(promocategory, 10, 64)
		if err == nil && strpromocategory > 0 {
			params["promo_category"] = promocategory
		} else {
			log.Error("Wrong input for parameter: promo_category")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_category", "Missing required parameter: promo_category")
		}
	} else {
		log.Error("Missing required parameter: promo_category cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_category cann't be blank", "Missing required parameter: promo_category cann't be blank")
	}

	//promo_nominal
	promonominalStr := c.FormValue("promo_nominal")
	if promonominalStr != "" {
		_, err := strconv.ParseFloat(promonominalStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: promo_nominal")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: promo_nominal", "Wrong input for parameter: promo_nominal")
		}
		params["promo_nominal"] = promonominalStr
	} else {
		log.Error("Missing required parameter: promo_nominal cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_nominal cann't be blank", "Missing required parameter: promo_nominal cann't be blank")
	}

	//promo_max_nominal
	promomaxnominalStr := c.FormValue("promo_max_nominal")
	if promomaxnominalStr != "" {
		_, err := strconv.ParseFloat(promomaxnominalStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: promo_max_nominal")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: promo_max_nominal", "Wrong input for parameter: promo_max_nominal")
		}
		params["promo_max_nominal"] = promomaxnominalStr
	} else {
		log.Error("Missing required parameter: promo_max_nominal cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_max_nominal cann't be blank", "Missing required parameter: promo_max_nominal cann't be blank")
	}

	//promo_values_type
	promovaluestype := c.FormValue("promo_values_type")
	if promovaluestype != "" {
		strpromovaluestype, err := strconv.ParseUint(promovaluestype, 10, 64)
		if err == nil && strpromovaluestype > 0 {
			params["promo_values_type"] = promovaluestype
		} else {
			log.Error("Wrong input for parameter: promo_values_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_values_type", "Missing required parameter: promo_values_type")
		}
	} else {
		log.Error("Missing required parameter: promo_values_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_values_type cann't be blank", "Missing required parameter: promo_values_type cann't be blank")
	}

	//promo_max_user
	promomaxuser := c.FormValue("promo_maxuser")
	if promomaxuser != "" {
		strpromomaxuser, err := strconv.ParseUint(promomaxuser, 10, 64)
		if err == nil && strpromomaxuser > 0 {
			params["promo_maxuser"] = promomaxuser
		} else {
			log.Error("Wrong input for parameter: promo_maxuser")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_maxuser", "Missing required parameter: promo_maxuser")
		}
	} else {
		log.Error("Missing required parameter: promo_maxuser cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_maxuser cann't be blank", "Missing required parameter: promo_maxuser cann't be blank")
	}

	//promo_stay_periode
	promostayperiode := c.FormValue("promo_stay_periode")
	if promostayperiode != "" {
		strpromostayperiode, err := strconv.ParseUint(promostayperiode, 10, 64)
		if err == nil && strpromostayperiode > 0 {
			params["promo_stay_periode"] = promostayperiode
		} else {
			log.Error("Wrong input for parameter: promo_stay_periode")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_stay_periode", "Missing required parameter: promo_stay_periode")
		}
	} else {
		log.Error("Missing required parameter: promo_stay_periode cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_stay_periode cann't be blank", "Missing required parameter: promo_stay_periode cann't be blank")
	}

	//promo_flag_uniq_user
	promoflaguniquser := c.FormValue("promo_flag_uniq_user")
	var promoflaguniquserBool bool
	if promoflaguniquser != "" {
		promoflaguniquserBool, err = strconv.ParseBool(promoflaguniquser)
		if err != nil {
			log.Error("promo_flag_uniq_user parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "promo_flag_uniq_user parameter should be true/false", "promo_flag_uniq_user parameter should be true/false")
		}
		if promoflaguniquserBool == true {
			params["promo_flag_uniq_user"] = "1"
		} else {
			params["promo_flag_uniq_user"] = "0"
		}
	} else {
		log.Error("promo_flag_uniq_user parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "promo_flag_uniq_user parameter should be true/false", "promo_flag_uniq_user parameter should be true/false")
	}

	//promo_valid_date1
	promovaliddate1 := c.FormValue("promo_valid_date1")
	if promovaliddate1 == "" {
		log.Error("Missing required parameter: promo_valid_date1 cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_valid_date1 cann't be blank", "Missing required parameter: promo_valid_date1 cann't be blank")
	}
	params["promo_valid_date1"] = promovaliddate1 + " 00:00:00"

	//promo_valid_date2
	promovaliddate2 := c.FormValue("promo_valid_date2")
	if promovaliddate2 == "" {
		log.Error("Missing required parameter: promo_valid_date2 cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_valid_date2 cann't be blank", "Missing required parameter: promo_valid_date2 cann't be blank")
	}
	params["promo_valid_date2"] = promovaliddate2 + " 23:59:59"

	//promo_notif_start
	promonotifstart := c.FormValue("promo_notif_start")
	if promonotifstart == "" {
		log.Error("Missing required parameter: promo_notif_start cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_notif_start cann't be blank", "Missing required parameter: promo_notif_start cann't be blank")
	}
	params["promo_notif_start"] = promonotifstart + " 00:00:00"

	//promo_notif_end
	promonotifend := c.FormValue("promo_notif_end")
	if promonotifend == "" {
		log.Error("Missing required parameter: promo_notif_end cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_notif_end cann't be blank", "Missing required parameter: promo_notif_end cann't be blank")
	}
	params["promo_notif_end"] = promonotifend + " 23:59:59"

	//promo_notif_type
	promonotiftype := c.FormValue("promo_notif_type")
	if promonotiftype != "" {
		strpromonotiftype, err := strconv.ParseUint(promonotiftype, 10, 64)
		if err == nil && strpromonotiftype > 0 {
			params["promo_notif_type"] = promonotiftype
		} else {
			log.Error("Wrong input for parameter: promo_notif_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_notif_type", "Missing required parameter: promo_notif_type")
		}
	} else {
		log.Error("Missing required parameter: promo_notif_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_notif_type cann't be blank", "Missing required parameter: promo_notif_type cann't be blank")
	}

	//promo_description
	promodescription := c.FormValue("promo_description")
	if promodescription == "" {
		log.Error("Missing required parameter: promo_description cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_description cann't be blank", "Missing required parameter: promo_description cann't be blank")
	}
	params["promo_description"] = promodescription

	//promo_tnc
	promotnc := c.FormValue("promo_tnc")
	if promotnc == "" {
		log.Error("Missing required parameter: promo_tnc cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_tnc cann't be blank", "Missing required parameter: promo_tnc cann't be blank")
	}
	params["promo_tnc"] = promotnc
	params["rec_status"] = "1"

	//promo_product_items
	promoproductitems := c.FormValue("promo_product_items")
	if promoproductitems == "" {
		log.Error("Missing required parameter: product cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product cann't be blank", "Missing required parameter: product cann't be blank")
	}

	var file *multipart.FileHeader
	file, err = c.FormFile("image")
	if file != nil {
		err = os.MkdirAll(config.BasePath+"/images/promo", 0755)
		if err != nil {
			log.Error(err.Error())
		} else {
			// Get file extension
			extension := filepath.Ext(file.Filename)
			log.Println("-------------extension---------------")
			log.Println(extension)
			log.Println("-------------extension---------------")
			// Generate filename
			filename := lib.RandStringBytesMaskImprSrc(20)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/promo/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image1"] = filename + extension
		}
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err, lastID := models.CreateTrPromo(params)
	if err != nil {
		log.Error("Failed create request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	s := strings.Split(promoproductitems, ",")

	var bindVarPromoPruduct []interface{}
	for _, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			var row []string
			row = append(row, lastID)                                     //promo_key
			row = append(row, is)                                         //product_key
			row = append(row, "1")                                        //flag_allowed
			row = append(row, "1")                                        //rec_status
			row = append(row, time.Now().Format(dateLayout))              //rec_created_date
			row = append(row, strconv.FormatUint(lib.Profile.UserID, 10)) //rec_created_by
			bindVarPromoPruduct = append(bindVarPromoPruduct, row)
		}
	}

	_, err = models.CreateMultiplePromoProduct(bindVarPromoPruduct)
	if err != nil {
		log.Error("Failed create promo product: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func DeletePromo(c echo.Context) error {
	var err error
	params := make(map[string]string)

	promoKey := c.FormValue("promo_key")
	if promoKey == "" {
		log.Error("Missing required parameter: promo_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_key", "Missing required parameter: promo_key")
	}

	promoKeyCek, err := strconv.ParseUint(promoKey, 10, 64)
	if err == nil && promoKeyCek > 0 {
		params["promo_key"] = promoKey
	} else {
		log.Error("Wrong input for parameter: promo_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_key", "Missing required parameter: promo_key")
	}

	var promo models.TrPromo
	status, err := models.GetTrPromo(&promo, "promo_key", promoKey)
	if err != nil {
		log.Error("Promo not found")
		return lib.CustomError(status)
	}

	var trans models.TrTransaction
	status, err = models.GetTrTransactionByField(&trans, "promo_code", *promo.PromoCode)
	if err == nil {
		log.Error("Promo already used in transaction, cann't delete promo")
		return lib.CustomError(http.StatusBadRequest, "Promo already used in transaction, cann't delete promo", "Promo already used in transaction, cann't delete promo")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateTrPromo(params)
	if err != nil {
		log.Error("Failed delete data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed delete data")
	}

	paramsPromo := make(map[string]string)
	paramsPromo["rec_status"] = "0"
	paramsPromo["rec_deleted_date"] = time.Now().Format(dateLayout)
	paramsPromo["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UpdateTrPromoProductByField(paramsPromo, "promo_key", promoKey)
	if err != nil {
		log.Error("Failed delete data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func UpdateAdminTrPromo(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	//promo_key
	promokey := c.FormValue("promo_key")
	if promokey != "" {
		strpromokey, err := strconv.ParseUint(promokey, 10, 64)
		if err == nil && strpromokey > 0 {
			params["promo_key"] = promokey
		} else {
			log.Error("Wrong input for parameter: promo_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_key", "Missing required parameter: promo_key")
		}
	} else {
		log.Error("Missing required parameter: promo_key cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_key cann't be blank", "Missing required parameter: promo_key cann't be blank")
	}

	//promo_code
	promocode := c.FormValue("promo_code")
	if promocode == "" {
		log.Error("Missing required parameter: promo_code cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_code cann't be blank", "Missing required parameter: promo_code cann't be blank")
	}

	var promoData models.TrPromo
	status, err = models.GetTrPromo(&promoData, "promo_key", promokey)
	if err != nil {
		log.Error("Promo not found")
		return lib.CustomError(status)
	}

	//cek duplikat promo code
	var promo models.TrPromo
	status, err = models.GetTrPromoValidasiDuplikat(&promo, "promo_code", promocode, promokey)
	if err == nil {
		log.Error("Missing required parameter: promo_code already used")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_code already used", "Missing required parameter: promo_code already used")
	}

	//cek promo used or not

	if promocode != *promoData.PromoCode {
		var trans models.TrTransaction
		status, err = models.GetTrTransactionByField(&trans, "promo_code", *promoData.PromoCode)
		if err == nil {
			log.Error("Promo already used in transaction, cann't update promo code")
			return lib.CustomError(http.StatusBadRequest, "Promo already used in transaction, cann't update promo code", "Promo already used in transaction, cann't update promo code")
		}
	}

	//promo_title
	promootitle := c.FormValue("promo_title")
	if promootitle == "" {
		log.Error("Missing required parameter: promo_title cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_title cann't be blank", "Missing required parameter: promo_title cann't be blank")
	}
	params["promo_title"] = promootitle

	//promo_category
	promocategory := c.FormValue("promo_category")
	if promocategory != "" {
		strpromocategory, err := strconv.ParseUint(promocategory, 10, 64)
		if err == nil && strpromocategory > 0 {
			params["promo_category"] = promocategory
		} else {
			log.Error("Wrong input for parameter: promo_category")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_category", "Missing required parameter: promo_category")
		}
	} else {
		log.Error("Missing required parameter: promo_category cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_category cann't be blank", "Missing required parameter: promo_category cann't be blank")
	}

	//promo_nominal
	promonominalStr := c.FormValue("promo_nominal")
	if promonominalStr != "" {
		_, err := strconv.ParseFloat(promonominalStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: promo_nominal")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: promo_nominal", "Wrong input for parameter: promo_nominal")
		}
		params["promo_nominal"] = promonominalStr
	} else {
		log.Error("Missing required parameter: promo_nominal cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_nominal cann't be blank", "Missing required parameter: promo_nominal cann't be blank")
	}

	//promo_max_nominal
	promomaxnominalStr := c.FormValue("promo_max_nominal")
	if promomaxnominalStr != "" {
		_, err := strconv.ParseFloat(promomaxnominalStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: promo_max_nominal")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: promo_max_nominal", "Wrong input for parameter: promo_max_nominal")
		}
		params["promo_max_nominal"] = promomaxnominalStr
	} else {
		log.Error("Missing required parameter: promo_max_nominal cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_max_nominal cann't be blank", "Missing required parameter: promo_max_nominal cann't be blank")
	}

	//promo_values_type
	promovaluestype := c.FormValue("promo_values_type")
	if promovaluestype != "" {
		strpromovaluestype, err := strconv.ParseUint(promovaluestype, 10, 64)
		if err == nil && strpromovaluestype > 0 {
			params["promo_values_type"] = promovaluestype
		} else {
			log.Error("Wrong input for parameter: promo_values_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_values_type", "Missing required parameter: promo_values_type")
		}
	} else {
		log.Error("Missing required parameter: promo_values_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_values_type cann't be blank", "Missing required parameter: promo_values_type cann't be blank")
	}

	//promo_max_user
	promomaxuser := c.FormValue("promo_maxuser")
	if promomaxuser != "" {
		strpromomaxuser, err := strconv.ParseUint(promomaxuser, 10, 64)
		if err == nil && strpromomaxuser > 0 {
			params["promo_maxuser"] = promomaxuser
		} else {
			log.Error("Wrong input for parameter: promo_maxuser")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_maxuser", "Missing required parameter: promo_maxuser")
		}
	} else {
		log.Error("Missing required parameter: promo_maxuser cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_maxuser cann't be blank", "Missing required parameter: promo_maxuser cann't be blank")
	}

	//promo_stay_periode
	promostayperiode := c.FormValue("promo_stay_periode")
	if promostayperiode != "" {
		strpromostayperiode, err := strconv.ParseUint(promostayperiode, 10, 64)
		if err == nil && strpromostayperiode > 0 {
			params["promo_stay_periode"] = promostayperiode
		} else {
			log.Error("Wrong input for parameter: promo_stay_periode")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_stay_periode", "Missing required parameter: promo_stay_periode")
		}
	} else {
		log.Error("Missing required parameter: promo_stay_periode cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_stay_periode cann't be blank", "Missing required parameter: promo_stay_periode cann't be blank")
	}

	//promo_flag_uniq_user
	promoflaguniquser := c.FormValue("promo_flag_uniq_user")
	var promoflaguniquserBool bool
	if promoflaguniquser != "" {
		promoflaguniquserBool, err = strconv.ParseBool(promoflaguniquser)
		if err != nil {
			log.Error("promo_flag_uniq_user parameter should be true/false")
			return lib.CustomError(http.StatusBadRequest, "promo_flag_uniq_user parameter should be true/false", "promo_flag_uniq_user parameter should be true/false")
		}
		if promoflaguniquserBool == true {
			params["promo_flag_uniq_user"] = "1"
		} else {
			params["promo_flag_uniq_user"] = "0"
		}
	} else {
		log.Error("promo_flag_uniq_user parameter should be true/false")
		return lib.CustomError(http.StatusBadRequest, "promo_flag_uniq_user parameter should be true/false", "promo_flag_uniq_user parameter should be true/false")
	}

	//promo_valid_date1
	promovaliddate1 := c.FormValue("promo_valid_date1")
	if promovaliddate1 == "" {
		log.Error("Missing required parameter: promo_valid_date1 cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_valid_date1 cann't be blank", "Missing required parameter: promo_valid_date1 cann't be blank")
	}
	params["promo_valid_date1"] = promovaliddate1 + " 00:00:00"

	//promo_valid_date2
	promovaliddate2 := c.FormValue("promo_valid_date2")
	if promovaliddate2 == "" {
		log.Error("Missing required parameter: promo_valid_date2 cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_valid_date2 cann't be blank", "Missing required parameter: promo_valid_date2 cann't be blank")
	}
	params["promo_valid_date2"] = promovaliddate2 + " 23:59:59"

	//promo_notif_start
	promonotifstart := c.FormValue("promo_notif_start")
	if promonotifstart == "" {
		log.Error("Missing required parameter: promo_notif_start cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_notif_start cann't be blank", "Missing required parameter: promo_notif_start cann't be blank")
	}
	params["promo_notif_start"] = promonotifstart + " 00:00:00"

	//promo_notif_end
	promonotifend := c.FormValue("promo_notif_end")
	if promonotifend == "" {
		log.Error("Missing required parameter: promo_notif_end cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_notif_end cann't be blank", "Missing required parameter: promo_notif_end cann't be blank")
	}
	params["promo_notif_end"] = promonotifend + " 23:59:59"

	//promo_notif_type
	promonotiftype := c.FormValue("promo_notif_type")
	if promonotiftype != "" {
		strpromonotiftype, err := strconv.ParseUint(promonotiftype, 10, 64)
		if err == nil && strpromonotiftype > 0 {
			params["promo_notif_type"] = promonotiftype
		} else {
			log.Error("Wrong input for parameter: promo_notif_type")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_notif_type", "Missing required parameter: promo_notif_type")
		}
	} else {
		log.Error("Missing required parameter: promo_notif_type cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_notif_type cann't be blank", "Missing required parameter: promo_notif_type cann't be blank")
	}

	//promo_description
	promodescription := c.FormValue("promo_description")
	if promodescription == "" {
		log.Error("Missing required parameter: promo_description cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_description cann't be blank", "Missing required parameter: promo_description cann't be blank")
	}
	params["promo_description"] = promodescription

	//promo_tnc
	promotnc := c.FormValue("promo_tnc")
	if promotnc == "" {
		log.Error("Missing required parameter: promo_tnc cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_tnc cann't be blank", "Missing required parameter: promo_tnc cann't be blank")
	}
	params["promo_tnc"] = promotnc
	params["rec_status"] = "1"

	//promo_product_items
	promoproductitems := c.FormValue("promo_product_items")
	if promoproductitems == "" {
		log.Error("Missing required parameter: product cann't be blank")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product cann't be blank", "Missing required parameter: product cann't be blank")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	var file *multipart.FileHeader
	file, err = c.FormFile("image")
	if file != nil {
		err = os.MkdirAll(config.BasePath+"/images/promo", 0755)
		if err != nil {
			log.Error(err.Error())
		} else {
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			filename := lib.RandStringBytesMaskImprSrc(20)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/promo/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image1"] = filename + extension
		}
	}

	status, err = models.UpdateTrPromo(params)
	if err != nil {
		log.Error("Failed update request data: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed update data")
	}

	s := strings.Split(promoproductitems, ",")

	var productIds []string

	for _, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			if _, ok := lib.Find(productIds, is); !ok {
				productIds = append(productIds, is)
			}
		}
	}

	//get promo product and delete
	if len(productIds) > 0 {
		var promoProductDelete []models.TrPromoProduct
		status, err = models.AdminGetPromoProductInNotIn(&promoProductDelete, productIds, "product_key", promokey, "NOT IN")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("Error Get Promo Product")
				return lib.CustomError(status)
			}
		}
		if len(promoProductDelete) > 0 {
			var ppKey []string
			for _, pp := range promoProductDelete {
				strPpKey := strconv.FormatUint(pp.PromoProductKey, 10)
				if _, ok := lib.Find(ppKey, strPpKey); !ok {
					ppKey = append(ppKey, strPpKey)
				}

			}
			if len(ppKey) > 0 {
				paramsPromo := make(map[string]string)
				paramsPromo["rec_status"] = "0"
				paramsPromo["rec_deleted_date"] = time.Now().Format(dateLayout)
				paramsPromo["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

				status, err = models.UpdateTrPromoProductByFieldIn(paramsPromo, "promo_product_key", ppKey)
				if err != nil {
					log.Error("Failed delete data: " + err.Error())
					return lib.CustomError(status, err.Error(), "failed delete data")
				}
			}
		}
	}

	//get new product and create
	if len(productIds) > 0 {
		var promoProductExist []models.TrPromoProduct
		status, err = models.AdminGetPromoProductInNotIn(&promoProductExist, productIds, "product_key", promokey, "IN")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("Error Get Promo Product")
				return lib.CustomError(status)
			}
		}

		var productKeyNew []string
		if len(promoProductExist) > 0 {
			var productKeyExist []string
			for _, pp := range promoProductExist {
				ssProductKey := strconv.FormatUint(pp.ProductKey, 10)
				if _, ok := lib.Find(productKeyExist, ssProductKey); !ok {
					productKeyExist = append(productKeyExist, ssProductKey)
				}
			}

			if len(productKeyExist) > 0 {
				for _, value := range productIds {
					is := strings.TrimSpace(value)
					if _, ok := lib.Find(productKeyExist, is); !ok {
						productKeyNew = append(productKeyNew, is)
					}
				}
			} else {
				productKeyNew = productIds
			}
		} else {
			productKeyNew = productIds
		}

		var bindVarPromoPruduct []interface{}
		if len(productKeyNew) > 0 {
			for _, value := range productKeyNew {
				is := strings.TrimSpace(value)
				if is != "" {
					var row []string
					row = append(row, promokey)                                   //promo_key
					row = append(row, is)                                         //product_key
					row = append(row, "1")                                        //flag_allowed
					row = append(row, "1")                                        //rec_status
					row = append(row, time.Now().Format(dateLayout))              //rec_created_date
					row = append(row, strconv.FormatUint(lib.Profile.UserID, 10)) //rec_created_by
					bindVarPromoPruduct = append(bindVarPromoPruduct, row)
				}
			}
			_, err = models.CreateMultiplePromoProduct(bindVarPromoPruduct)
			if err != nil {
				log.Error("Failed create promo product: " + err.Error())
				return lib.CustomError(status, err.Error(), "failed input data")
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func DetailPromo(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var promo models.TrPromo
	status, err = models.GetTrPromo(&promo, "promo_key", keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	var responseData models.TrPromoDetail

	var lookupIds []string

	if _, ok := lib.Find(lookupIds, strconv.FormatUint(promo.PromoCategory, 10)); !ok {
		lookupIds = append(lookupIds, strconv.FormatUint(promo.PromoCategory, 10))
	}

	if _, ok := lib.Find(lookupIds, strconv.FormatUint(promo.PromoValuesType, 10)); !ok {
		lookupIds = append(lookupIds, strconv.FormatUint(promo.PromoValuesType, 10))
	}

	if _, ok := lib.Find(lookupIds, strconv.FormatUint(promo.PromoNotifType, 10)); !ok {
		lookupIds = append(lookupIds, strconv.FormatUint(promo.PromoNotifType, 10))
	}

	//gen lookup oa request
	var lookupOaReq []models.GenLookup
	if len(lookupIds) > 0 {
		status, err = models.GetGenLookupIn(&lookupOaReq, lookupIds, "lookup_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}

	gData := make(map[uint64]models.GenLookup)
	for _, gen := range lookupOaReq {
		gData[gen.LookupKey] = gen
	}

	responseData.PromoKey = promo.PromoKey
	responseData.PromoCode = promo.PromoCode
	responseData.PromoTitle = promo.PromoTitle

	dir := config.BaseUrl + "/images/promo/"

	if promo.RecImage1 != nil {
		path := dir + *promo.RecImage1
		responseData.Image = &path
	}

	if n, ok := gData[promo.PromoCategory]; ok {
		var trc models.LookupTrans

		trc.LookupKey = n.LookupKey
		trc.LkpGroupKey = n.LkpGroupKey
		trc.LkpCode = n.LkpCode
		trc.LkpName = n.LkpName
		responseData.PromoCategory = trc
	}

	responseData.PromoNominal = promo.PromoNominal.Truncate(0)
	responseData.PromoMaxNominal = promo.PromoMaxNominal.Truncate(0)

	if n, ok := gData[promo.PromoValuesType]; ok {
		var trc models.LookupTrans

		trc.LookupKey = n.LookupKey
		trc.LkpGroupKey = n.LkpGroupKey
		trc.LkpCode = n.LkpCode
		trc.LkpName = n.LkpName
		responseData.PromoValuesType = trc
	}

	dateLayout := "2006-01-02 15:04:05"
	newlayout := "2006-01-02"

	responseData.PromoMaxuser = promo.PromoMaxuser
	responseData.PromoStayPeriode = promo.PromoStayPeriode
	responseData.PromoFlagUniqUser = promo.PromoFlagUniqUser

	t, _ := time.Parse(dateLayout, promo.PromoValidDate1)
	responseData.PromoValidDate1 = t.Format(newlayout)

	t, _ = time.Parse(dateLayout, promo.PromoValidDate2)
	responseData.PromoValidDate2 = t.Format(newlayout)

	t, _ = time.Parse(dateLayout, promo.PromoNotifStart)
	responseData.PromoNotifStart = t.Format(newlayout)

	t, _ = time.Parse(dateLayout, promo.PromoNotifEnd)
	responseData.PromoNotifEnd = t.Format(newlayout)

	if n, ok := gData[promo.PromoNotifType]; ok {
		var trc models.LookupTrans

		trc.LookupKey = n.LookupKey
		trc.LkpGroupKey = n.LkpGroupKey
		trc.LkpCode = n.LkpCode
		trc.LkpName = n.LkpName
		responseData.PromoNotifType = trc
	}

	responseData.PromoNotifDescription = promo.PromoNotifDescription
	responseData.PromoTnc = promo.PromoTnc

	var promoProductData []models.TrPromoProductData
	status, err = models.AdminGetPromoProductByPromoKey(&promoProductData, keyStr)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	responseData.PromoProduct = promoProductData

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func CheckPromo(c echo.Context) error {
	//promo_code
	promoCode := c.FormValue("promo_code")
	if promoCode == "" {
		log.Error("Wrong input for parameter: promo_code")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: promo_code", "Missing required parameter: promo_code")
	}

	//customer_key
	customerKey := c.FormValue("customer_key")
	if customerKey == "" {
		log.Error("Wrong input for parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	//product_key
	productKey := c.FormValue("product_key")
	if productKey == "" {
		log.Error("Wrong input for parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	err, enable, text, _ := validatePromo(promoCode, customerKey, productKey)

	if err != nil {
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
	}

	var responseData models.CheckPromo
	responseData.PromoEnabled = enable
	responseData.Message = text

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func validatePromo(promoCode string, customerKey string, productKey string) (error, bool, string, *string) {
	var err error

	var promoEnabled bool
	var message string
	var promoKey *string

	//1. cek apakah promo ada
	var promo models.TrPromo
	_, err = models.GetTrPromoProductActive(&promo, promoCode, productKey)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return err, false, "Failed get data", nil
		} else {
			promoEnabled = false
			message = "Promo tidak ditemukan."
			promoKey = nil
		}
	} else {
		//cek apakah promo unik user - cek udah pernah pake apa belom
		checkCountUsedPromo := 0
		promoKeyStr := strconv.FormatUint(promo.PromoKey, 10)
		if promo.PromoFlagUniqUser == 1 {
			var countData models.CountData
			_, err = models.AdminGetCountPromoUsed(&countData, &promoKeyStr, &customerKey)
			if err != nil {
				log.Error(err.Error())
				return err, false, "Failed get data", nil
			}
			if int(countData.CountData) > int(0) {
				promoEnabled = false
				message = "Promo tidak dapat digunakan lagi, kamu sudah pernah menggunakannya."
				promoKey = nil
			} else {
				checkCountUsedPromo = 1
			}
		} else {
			checkCountUsedPromo = 1
		}

		if checkCountUsedPromo == 1 {
			var countData models.CountData
			_, err = models.AdminGetCountPromoUsed(&countData, &promoKeyStr, nil)
			if err != nil {
				log.Error(err.Error())
				return err, false, "Failed get data", nil
			}
			if int(countData.CountData) >= int(promo.PromoMaxuser) {
				promoEnabled = false
				message = "Promo sudah melebihi jumlah batas penggunaan user."
				promoKey = nil
			} else {
				promoEnabled = true
				message = promo.PromoNotifDescription
				promoKey = &promoKeyStr
			}
		}
	}
	return nil, promoEnabled, message, promoKey
}
