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
	pagination := c.Param("pagination")
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
	defer res.Body.Close()
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
		log.Error("Error create log endpoint motion pay")
		log.Error(err.Error())
	}
	// fmt.Println(res)
	// fmt.Println(string(body))
	return res.StatusCode, string(body), err
}
