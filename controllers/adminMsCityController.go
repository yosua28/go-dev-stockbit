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

func AdminGetListMsCity(c echo.Context) error {

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

	items := []string{"cou_name", "city_parent", "city_name", "city_code", "city_level", "postal_code"}

	params := make(map[string]string)
	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		_, found := lib.Find(items, orderBy)
		if found {
			var ord string
			if orderBy == "cou_name" {
				ord = "cou.cou_name"
			} else if orderBy == "city_parent" {
				ord = "par.city_name"
			} else if orderBy == "city_name" {
				ord = "c.city_name"
			} else if orderBy == "city_code" {
				ord = "c.city_code"
			} else if orderBy == "city_level" {
				ord = "cl.lkp_name"
			} else if orderBy == "postal_code" {
				ord = "c.postal_code"
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
		params["orderBy"] = "c.city_level"
		params["orderType"] = "ASC"
	}

	searchLike := c.QueryParam("search_like")

	countryKey := c.QueryParam("country_key")
	if countryKey != "" {
		params["c.country_key"] = countryKey
	}

	parentKey := c.QueryParam("parent_key")
	if parentKey != "" {
		params["c.parent_key"] = parentKey
	}

	cityLevel := c.QueryParam("city_level")
	if cityLevel != "" {
		params["c.city_level"] = cityLevel
	}

	var city []models.ListCity

	status, err = models.AdminGetListCity(&city, limit, offset, params, searchLike, noLimit)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(city) < 1 {
		log.Error("City not found")
		return lib.CustomError(http.StatusNotFound, "City not found", "City not found")
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.CountAdminGetCity(&countData, params, searchLike)
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
	response.Data = city

	return c.JSON(http.StatusOK, response)
}

func AdminDeleteMsCity(c echo.Context) error {
	var err error

	params := make(map[string]string)

	keyStr := c.FormValue("city_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: city_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: city_key", "Missing required parameter: city_key")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["city_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateMsCity(params)
	if err != nil {
		log.Error("Error delete ms_city")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed delete data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func GetCityLevel(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	params["lkp_group_key"] = "47"
	params["rec_status"] = "1"

	var lookupDB []models.GenLookup
	status, err = models.GetAllGenLookup(&lookupDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var responseData []models.GenLookupDataInfo
	for _, lkp := range lookupDB {
		if lkp.LkpName != nil && lkp.LkpCode != nil {
			var data models.GenLookupDataInfo
			data.Name = *lkp.LkpName
			data.Value = *lkp.LkpCode
			responseData = append(responseData, data)
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func AdminCreateMsCity(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	countryKey := c.FormValue("country_key")
	if countryKey == "" {
		log.Error("Missing required parameter: country_key")
		return lib.CustomError(http.StatusBadRequest, "country_key can not be blank", "country_key can not be blank")
	} else {
		n, err := strconv.ParseUint(countryKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: country_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: country_key", "Wrong input for parameter: country_key")
		}
		params["country_key"] = countryKey
	}

	cityLevel := c.FormValue("city_level")
	if cityLevel == "" {
		log.Error("Missing required parameter: city_level")
		return lib.CustomError(http.StatusBadRequest, "city_level can not be blank", "city_level can not be blank")
	} else {
		n, err := strconv.ParseUint(cityLevel, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: city_level")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: city_level", "Wrong input for parameter: city_level")
		}
		params["city_level"] = cityLevel
	}

	regionArea := c.FormValue("region_area")
	province := c.FormValue("province")
	kabKodya := c.FormValue("kab_kodya")

	if cityLevel == "1" {

	} else if cityLevel == "2" {
		if regionArea == "" {
			log.Error("Missing required parameter: region_area")
			return lib.CustomError(http.StatusBadRequest, "region_area can not be blank", "region_area can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: region_area")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: region_area", "Wrong input for parameter: region_area")
			}
		}
		params["parent_key"] = regionArea
	} else if cityLevel == "3" {
		if regionArea == "" {
			log.Error("Missing required parameter: region_area")
			return lib.CustomError(http.StatusBadRequest, "region_area can not be blank", "region_area can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: region_area")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: region_area", "Wrong input for parameter: region_area")
			}
		}

		if province == "" {
			log.Error("Missing required parameter: province")
			return lib.CustomError(http.StatusBadRequest, "province can not be blank", "province can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: province")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: province", "Wrong input for parameter: province")
			}
		}
		params["parent_key"] = province
	} else if cityLevel == "4" {
		if regionArea == "" {
			log.Error("Missing required parameter: region_area")
			return lib.CustomError(http.StatusBadRequest, "region_area can not be blank", "region_area can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: region_area")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: region_area", "Wrong input for parameter: region_area")
			}
		}

		if province == "" {
			log.Error("Missing required parameter: province")
			return lib.CustomError(http.StatusBadRequest, "province can not be blank", "province can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: province")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: province", "Wrong input for parameter: province")
			}
		}

		if kabKodya == "" {
			log.Error("Missing required parameter: kab_kodya")
			return lib.CustomError(http.StatusBadRequest, "kab_kodya can not be blank", "kab_kodya can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: kab_kodya")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kab_kodya", "Wrong input for parameter: kab_kodya")
			}
		}
		params["parent_key"] = kabKodya
	}

	cityCode := c.FormValue("city_code")
	if cityCode == "" {
		log.Error("Missing required parameter: city_code")
		return lib.CustomError(http.StatusBadRequest, "city_code can not be blank", "city_code can not be blank")
	} else {
		//validate unique city_code
		var countData models.CountData
		status, err = models.CountMsCityValidateUnique(&countData, "city_code", cityCode, "")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: city_code")
			return lib.CustomError(http.StatusBadRequest, "city_code already used", "city_code already used")
		}
		params["city_code"] = cityCode
	}

	cityName := c.FormValue("city_name")
	if cityName == "" {
		log.Error("Missing required parameter: city_name")
		return lib.CustomError(http.StatusBadRequest, "city_name can not be blank", "city_name can not be blank")
	} else {
		params["city_name"] = cityName
	}

	postalCode := c.FormValue("postal_code")
	if postalCode != "" {
		params["postal_code"] = postalCode
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		_, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.CreateMsCity(params)
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

func AdminUpdateMsCity(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]string)

	cityKey := c.FormValue("city_key")
	if cityKey != "" {
		n, err := strconv.ParseUint(cityKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: city_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: city_key", "Wrong input for parameter: city_key")
		}
		params["city_key"] = cityKey
	} else {
		log.Error("Missing required parameter: city_key")
		return lib.CustomError(http.StatusBadRequest, "city_key can not be blank", "city_key can not be blank")
	}

	countryKey := c.FormValue("country_key")
	if countryKey == "" {
		log.Error("Missing required parameter: country_key")
		return lib.CustomError(http.StatusBadRequest, "country_key can not be blank", "country_key can not be blank")
	} else {
		n, err := strconv.ParseUint(countryKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: country_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: country_key", "Wrong input for parameter: country_key")
		}
		params["country_key"] = countryKey
	}

	cityLevel := c.FormValue("city_level")
	if cityLevel == "" {
		log.Error("Missing required parameter: city_level")
		return lib.CustomError(http.StatusBadRequest, "city_level can not be blank", "city_level can not be blank")
	} else {
		n, err := strconv.ParseUint(cityLevel, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: city_level")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: city_level", "Wrong input for parameter: city_level")
		}
		params["city_level"] = cityLevel
	}

	regionArea := c.FormValue("region_area")
	province := c.FormValue("province")
	kabKodya := c.FormValue("kab_kodya")

	if cityLevel == "1" {

	} else if cityLevel == "2" {
		if regionArea == "" {
			log.Error("Missing required parameter: region_area")
			return lib.CustomError(http.StatusBadRequest, "region_area can not be blank", "region_area can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: region_area")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: region_area", "Wrong input for parameter: region_area")
			}
		}
		params["parent_key"] = regionArea
	} else if cityLevel == "3" {
		if regionArea == "" {
			log.Error("Missing required parameter: region_area")
			return lib.CustomError(http.StatusBadRequest, "region_area can not be blank", "region_area can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: region_area")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: region_area", "Wrong input for parameter: region_area")
			}
		}

		if province == "" {
			log.Error("Missing required parameter: province")
			return lib.CustomError(http.StatusBadRequest, "province can not be blank", "province can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: province")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: province", "Wrong input for parameter: province")
			}
		}
		params["parent_key"] = province
	} else if cityLevel == "4" {
		if regionArea == "" {
			log.Error("Missing required parameter: region_area")
			return lib.CustomError(http.StatusBadRequest, "region_area can not be blank", "region_area can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: region_area")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: region_area", "Wrong input for parameter: region_area")
			}
		}

		if province == "" {
			log.Error("Missing required parameter: province")
			return lib.CustomError(http.StatusBadRequest, "province can not be blank", "province can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: province")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: province", "Wrong input for parameter: province")
			}
		}

		if kabKodya == "" {
			log.Error("Missing required parameter: kab_kodya")
			return lib.CustomError(http.StatusBadRequest, "kab_kodya can not be blank", "kab_kodya can not be blank")
		} else {
			n, err := strconv.ParseUint(cityLevel, 10, 64)
			if err != nil || n == 0 {
				log.Error("Wrong input for parameter: kab_kodya")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: kab_kodya", "Wrong input for parameter: kab_kodya")
			}
		}
		params["parent_key"] = kabKodya
	}

	cityCode := c.FormValue("city_code")
	if cityCode == "" {
		log.Error("Missing required parameter: city_code")
		return lib.CustomError(http.StatusBadRequest, "city_code can not be blank", "city_code can not be blank")
	} else {
		//validate unique city_code
		var countData models.CountData
		status, err = models.CountMsCityValidateUnique(&countData, "city_code", cityCode, cityKey)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if int(countData.CountData) > int(0) {
			log.Error("Missing required parameter: city_code")
			return lib.CustomError(http.StatusBadRequest, "city_code already used", "city_code already used")
		}
		params["city_code"] = cityCode
	}

	cityName := c.FormValue("city_name")
	if cityName == "" {
		log.Error("Missing required parameter: city_name")
		return lib.CustomError(http.StatusBadRequest, "city_name can not be blank", "city_name can not be blank")
	} else {
		params["city_name"] = cityName
	}

	postalCode := c.FormValue("postal_code")
	if postalCode != "" {
		params["postal_code"] = postalCode
	}

	recOrder := c.FormValue("rec_order")
	if recOrder != "" {
		_, err := strconv.ParseUint(recOrder, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: rec_order")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: rec_order", "Wrong input for parameter: rec_order")
		}
		params["rec_order"] = recOrder
	}

	dateLayout := "2006-01-02 15:04:05"
	params["rec_modified_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	params["rec_status"] = "1"

	status, err = models.UpdateMsCity(params)
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

func AdminDetailMsCity(c echo.Context) error {
	var err error

	cityKey := c.Param("city_key")
	if cityKey == "" {
		log.Error("Missing required parameter: city_key")
		return lib.CustomError(http.StatusBadRequest, "city_key can not be blank", "city_key can not be blank")
	} else {
		n, err := strconv.ParseUint(cityKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: city_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: city_key", "Wrong input for parameter: city_key")
		}
	}
	var city models.MsCity
	_, err = models.GetMsCity(&city, cityKey)
	if err != nil {
		log.Error("City not found")
		return lib.CustomError(http.StatusBadRequest, "City not found", "City not found")
	}

	if city.CityLevel == uint64(0) || city.CityLevel > uint64(4) {
		log.Error("City not found")
		return lib.CustomError(http.StatusBadRequest, "City not found", "City not found")
	}

	responseData := make(map[string]interface{})

	responseData["city_key"] = city.CityKey
	responseData["country_key"] = city.CountryKey
	responseData["city_code"] = city.CityCode
	responseData["city_name"] = city.CityName
	responseData["city_level"] = city.CityLevel
	if city.PostalCode != nil {
		responseData["postal_code"] = *city.PostalCode
	} else {
		responseData["postal_code"] = ""
	}
	if city.RecOrder != nil {
		responseData["rec_order"] = *city.RecOrder
	} else {
		responseData["rec_order"] = ""
	}

	if city.CityLevel == uint64(2) { //province
		responseData["region_area"] = city.ParentKey
	} else if city.CityLevel == uint64(3) { //kab_kodya
		if city.ParentKey != nil {
			responseData["province"] = *city.ParentKey

			var regionArea models.MsCity
			_, err = models.GetMsCity(&regionArea, strconv.FormatUint(*city.ParentKey, 10))
			if err != nil {
				log.Error("City (Region Area) not found")
				return lib.CustomError(http.StatusBadRequest, "City (Region Area) not found", "City (Region Area) not found")
			}

			if regionArea.ParentKey != nil {
				responseData["region_area"] = *regionArea.ParentKey
			}
		}
	} else if city.CityLevel == uint64(4) { //kecamatan
		if city.ParentKey != nil {
			responseData["kab_kodya"] = *city.ParentKey

			var province models.MsCity
			_, err = models.GetMsCity(&province, strconv.FormatUint(*city.ParentKey, 10))
			if err != nil {
				log.Error("City (Province) not found")
				return lib.CustomError(http.StatusBadRequest, "City (Province) not found", "City (Province) not found")
			}

			if province.ParentKey != nil {
				responseData["province"] = *province.ParentKey

				var regionArea models.MsCity
				_, err = models.GetMsCity(&regionArea, strconv.FormatUint(*province.ParentKey, 10))
				if err != nil {
					log.Error("City (Region Area) not found")
					return lib.CustomError(http.StatusBadRequest, "City (Region Area) not found", "City (Region Area) not found")
				}

				if regionArea.ParentKey != nil {
					responseData["region_area"] = *regionArea.ParentKey
				}
			}
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func GetCityParent(c echo.Context) error {
	var err error
	var status int

	var city []models.ListParent
	status, err = models.AdminGetListParent(&city)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = city

	return c.JSON(http.StatusOK, response)
}
