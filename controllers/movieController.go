package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	_ "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func SearchMovies(c echo.Context) error {

	searchword := c.Param("searchword")
	if searchword == "" {
		log.Error("Missing required parameter: searchword")
		return lib.CustomError(http.StatusBadRequest, "searchword can not be blank", "searchword can not be blank")
	}
	pagination := c.Param("pagination")
	if pagination == "" {
		log.Error("Missing required parameter: pagination")
		return lib.CustomError(http.StatusBadRequest, "pagination can not be blank", "pagination can not be blank")
	} else {
		n, err := strconv.ParseUint(pagination, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: pagination")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: pagination", "Wrong input for parameter: pagination")
		}
	}

	urlparam := "&s=" + searchword + "&page" + pagination
	status, res, _ := requestToOmdbapi(config.API_SEARCH, urlparam)

	var dataBody map[string]interface{}
	_ = json.Unmarshal([]byte(res), &dataBody)

	if status != http.StatusOK {
		return lib.CustomError(status, dataBody["Error"].(string), dataBody["Error"].(string))
	}

	var response lib.Response
	response.Status.Code = uint(status)
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = dataBody
	return c.JSON(http.StatusOK, response)
}

func DetailMovies(c echo.Context) error {

	imdbid := c.Param("imdbid")

	if imdbid == "" {
		log.Error("Missing required parameter: imdbid")
		return lib.CustomError(http.StatusBadRequest, "imdbid can not be blank", "imdbid can not be blank")
	}
	urlparam := "&i=" + imdbid

	status, res, _ := requestToOmdbapi(config.API_DETAIL, urlparam)
	var dataBody map[string]interface{}
	_ = json.Unmarshal([]byte(res), &dataBody)
	if status != http.StatusOK {
		return lib.CustomError(status, dataBody["Error"].(string), dataBody["Error"].(string))
	}

	var response lib.Response
	response.Status.Code = uint(status)
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = dataBody
	return c.JSON(http.StatusOK, response)
}

func requestToOmdbapi(apiname string, param string) (int, string, error) {

	url := config.SandBox + "?apikey=" + config.OMDBKey + param
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)

	//CREATE LOG
	paramLog := make(map[string]string)
	dateLayout := "2006-01-02 15:04:05"
	paramLog["name"] = apiname
	paramLog["request_method"] = "GET"
	paramLog["url"] = url
	paramLog["status"] = strconv.FormatUint(uint64(res.StatusCode), 10)
	paramLog["response"] = string(body)
	paramLog["created_date"] = time.Now().Format(dateLayout)
	_, err = models.CreateLog(paramLog)
	if err != nil {
		log.Error("Error create log" + err.Error())
	}
	return res.StatusCode, string(body), err
}
