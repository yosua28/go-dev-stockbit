package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func MotionPayStatus(c echo.Context) error {
	decimal.MarshalJSONWithoutQuotes = true

	responseData := make(map[string]interface{})
	responseData["status"] = lib.STATUS_NON_LINKED
	responseData["max_amount"] = decimal.NewFromInt(lib.MAX_AMOUNT_MOTION_PAY)

	var linkage models.ScLinkage
	_, err := models.GetLinkageByField(&linkage, strconv.FormatUint(lib.Profile.UserID, 10), "user_login_key")
	if err != nil {
		responseData["status"] = lib.STATUS_NON_LINKED
	} else {
		if linkage.LinkedStatus == nil {
			responseData["status"] = lib.STATUS_NON_LINKED
		} else {
			if *linkage.LinkedStatus != lib.STATUS_LINKED {
				responseData["status"] = lib.STATUS_NON_LINKED
			} else {
				responseData["status"] = lib.STATUS_LINKED
				//get saldo cash
				responseData["cash_balance"] = getBalance(*linkage.UserToken, lib.CASH_BALANCE)

				//get saldo point
				responseData["point_balance"] = getBalance(*linkage.UserToken, lib.POINT_BALANCE)
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

func LinkingMotionPay(c echo.Context) error {
	var err error
	var status int

	params := make(map[string]interface{})

	phone := c.FormValue("phone")
	if phone == "" {
		log.Error("Wrong input for parameter: phone")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: phone", "Wrong input for parameter: phone")
	}
	linkedId := ""

	var linkage models.ScLinkage
	_, err = models.GetLinkageByField(&linkage, strconv.FormatUint(lib.Profile.UserID, 10), "user_login_key")
	if err == nil {
		linkedId = strconv.FormatUint(linkage.LinkedKey, 10)
		if linkage.LinkedStatus != nil && *linkage.LinkedStatus == lib.STATUS_LINKED {
			log.Error("Proses linked tidak bisa. User sudah terlinked")
			return lib.CustomError(http.StatusBadRequest, "Anda sudah terlinked", "Anda sudah terlinked")
		}
	}

	var dataMotion models.RegisPremiumMotionPay
	_, err = models.GetDataRegisPremiumMotionPay(&dataMotion, strconv.FormatUint(lib.Profile.UserID, 10))
	if err != nil {
		log.Println(dataMotion.Fullname)
		log.Println(err.Error())
		log.Error("User data tidak ditemukan")
		return lib.CustomError(http.StatusBadRequest, "Data tidak ditemukan. Silakan menghubungi Customer Service untuk informasi lebih lanjut.", "Data tidak ditemukan. Silakan menghubungi Customer Service untuk informasi lebih lanjut.")
	}

	log.Println(dataMotion.Fullname)

	params["terminal_id"] = ""
	params["fullname"] = dataMotion.Fullname
	params["phone"] = phone
	if dataMotion.Email != nil {
		params["email"] = *dataMotion.Email
	} else {
		params["email"] = ""
	}
	params["mother_maid_name"] = dataMotion.MotherMaidName
	params["id_type"] = dataMotion.IdType
	params["id_number"] = dataMotion.IdNumber
	params["birth_date"] = dataMotion.BirthDate
	params["birth_city"] = dataMotion.BirthCity
	params["gender"] = dataMotion.Gender
	if dataMotion.IdType != "EKTP" {
		if dataMotion.ProvinceCode != nil {
			params["province_code"] = *dataMotion.ProvinceCode
		} else {
			params["province_code"] = "11"
		}
		if dataMotion.DistrictCode != nil {
			params["district_code"] = *dataMotion.DistrictCode
		} else {
			params["district_code"] = "01"
		}
		if dataMotion.SubDistrictCode != nil {
			params["sub_district_code"] = *dataMotion.SubDistrictCode
		} else {
			params["sub_district_code"] = "01"
		}
	}
	params["address"] = dataMotion.Address
	params["address_area"] = ""
	if dataMotion.Nationality != nil {
		params["nationality"] = *dataMotion.Nationality
	} else {
		params["nationality"] = "ID"
	}
	if dataMotion.Occupation != nil {
		params["occupation"] = *dataMotion.Occupation
	} else {
		params["occupation"] = "LAINNYA"
	}
	params["remark"] = *dataMotion.Remark
	params["referral_code"] = ""
	id, err := machineid.ID()
	if err == nil {
		params["device_id"] = id
	} else {
		params["device_id"] = "47658790707jhfj678oihkt98769"
	}
	params["timezone"] = "Asia/Jakarta"

	paramsFds := make(map[string]interface{})
	paramsFds["ip_address"] = ""
	paramsFds["sdk_version"] = ""
	paramsFds["device_type"] = ""
	paramsFds["device_os"] = ""
	paramsFds["user_agent"] = ""
	paramsFds["geo_location"] = ""
	params["fds"] = paramsFds

	status, res, header, err := requestToMotionPay(
		phone,
		lib.REGISTRATION_PREMIUM_NO_OTP,
		lib.PATH_REGISTRATION_PREMIUM_NO_OTP,
		"POST",
		params,
		false,
		true,
	)
	if err == nil && status == http.StatusOK {
		var dataBody map[string]interface{}
		err := json.Unmarshal([]byte(res), &dataBody)
		log.Println(err)
		messageData := dataBody["message_data"].(map[string]interface{})

		log.Println(messageData)
		log.Println(header)
		log.Println(linkedId)

		//save sc_linkage
		paramsLinkage := make(map[string]string)
		paramsLinkage["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
		paramsLinkage["settle_channel"] = "299"
		paramsLinkage["linked_mobileno"] = params["phone"].(string)
		paramsLinkage["linked_full_name"] = params["fullname"].(string)
		paramsLinkage["linked_member_type"] = lib.PREMIUM
		paramsLinkage["mother_maid_name"] = params["mother_maid_name"].(string)
		paramsLinkage["id_type"] = params["id_type"].(string)
		paramsLinkage["id_number"] = params["id_number"].(string)
		paramsLinkage["birth_date"] = params["birth_date"].(string)
		paramsLinkage["birth_city"] = params["birth_city"].(string)
		paramsLinkage["gender"] = params["gender"].(string)
		paramsLinkage["address"] = params["address"].(string)
		paramsLinkage["nationality"] = params["nationality"].(string)
		paramsLinkage["occupation"] = params["occupation"].(string)
		paramsLinkage["remark"] = params["remark"].(string)
		paramsLinkage["device_id"] = params["device_id"].(string)
		paramsLinkage["linked_status"] = lib.STATUS_LINKED
		paramsLinkage["linked_name"] = messageData["status"].(string)
		paramsLinkage["user_token"] = messageData["user_token"].(string)
		paramsLinkage["rec_attribute_id1"] = header
		jsonString, err := json.Marshal(params)
		if err == nil {
			paramsLinkage["rec_attribute_id2"] = string(jsonString)
		}
		paramsLinkage["rec_attribute_id3"] = res
		paramsLinkage["rec_order"] = "0"
		paramsLinkage["rec_status"] = "1"

		dateLayout := "2006-01-02 15:04:05"
		if linkedId != "" {
			paramsLinkage["linked_key"] = linkedId
			paramsLinkage["rec_modified_date"] = time.Now().Format(dateLayout)
			paramsLinkage["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		} else {
			paramsLinkage["rec_created_date"] = time.Now().Format(dateLayout)
			paramsLinkage["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		}

		if linkedId != "" {
			status, err = models.UpdateScLinkage(paramsLinkage)
			if err != nil {
				log.Error("Error save sc_linkage")
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Permintaan tidak dapat di proses. Silakan menghubungi Customer Service untuk informasi lebih lanjut.")
			}
		} else {
			status, err = models.CreateScLinkage(paramsLinkage)
			if err != nil {
				log.Error("Error save sc_linkage")
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Permintaan tidak dapat di proses. Silakan menghubungi Customer Service untuk informasi lebih lanjut.")
			}
		}
	} else {
		log.Error("Error response dari motion pay")
		log.Error(status, err.Error())
		return lib.CustomError(status, err.Error(), "Permintaan tidak dapat di proses. Silakan menghubungi Customer Service untuk informasi lebih lanjut.")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func UnlinkingMotionPay(c echo.Context) error {
	var err error
	var status int

	paramsLinkage := make(map[string]string)
	paramsLinkage["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramsLinkage["rec_status"] = "0"
	paramsLinkage["linked_status"] = lib.STATUS_UNLINKED
	dateLayout := "2006-01-02 15:04:05"
	paramsLinkage["rec_modified_date"] = time.Now().Format(dateLayout)
	paramsLinkage["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	status, err = models.UnlinkedUser(paramsLinkage, "user_login_key", strconv.FormatUint(lib.Profile.UserID, 10))
	if err != nil {
		log.Error("Error save sc_linkage")
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Permintaan tidak dapat di proses. Silakan menghubungi Customer Service untuk informasi lebih lanjut.")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""

	return c.JSON(http.StatusOK, response)
}

func getBalance(userToken string, apiName string) decimal.Decimal {
	decimal.MarshalJSONWithoutQuotes = true
	params := make(map[string]interface{})
	var val decimal.Decimal
	val = decimal.NewFromInt(0)
	var path string
	if apiName == lib.CASH_BALANCE {
		path = lib.PATH_CASH_BALANCE
	} else {
		path = lib.PATH_POINT_BALANCE
	}
	status, res, _, err := requestToMotionPay(
		userToken,
		apiName,
		path,
		"GET",
		params,
		true,
		false,
	)
	if err == nil && status == http.StatusOK {
		var dataBody map[string]interface{}
		err := json.Unmarshal([]byte(res), &dataBody)
		log.Println(err)
		messageData := dataBody["message_data"].(map[string]interface{})

		if apiName == lib.CASH_BALANCE {
			params := fmt.Sprintf("%v", messageData["cash_balance"])
			taxRate, _ := decimal.NewFromString(params)
			val = taxRate
		} else if apiName == lib.POINT_BALANCE {
			params := fmt.Sprintf("%v", messageData["points_balance"])
			taxRate, _ := decimal.NewFromString(params)
			val = taxRate
		}
	} else {
		log.Error(status, err.Error())
	}

	return val
}

func requestToMotionPay(
	userTokenOrNoHp string,
	apiName string,
	path string,
	requestMethod string,
	params map[string]interface{},
	useHeaderUserToken bool,
	useHeaderContenType bool) (int, string, string, error) {

	paramLog := make(map[string]string)

	url := config.SandBox + path
	dateLayout := "2006-01-02 15:04:05"
	paramLog["merchant"] = "MOTION PAY"
	paramLog["endpoint_name"] = apiName
	paramLog["request_method"] = requestMethod
	paramLog["url"] = url
	paramLog["created_date"] = time.Now().Format(dateLayout)
	paramLog["created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
	paramLog["note"] = "GET DATA MOTION PAY " + apiName

	jsonString, err := json.Marshal(params)
	payload := strings.NewReader(string(jsonString))
	req, err := http.NewRequest(requestMethod, url, payload)
	if err != nil {
		log.Error("Error1", err.Error())
		return http.StatusBadGateway, "", paramLog["header"], err
	}
	req.Header.Add("auth-merchant", config.MERCHANT_ID)
	req.Header.Add("auth-partner", config.PARTNER_ID)
	req.Header.Add("auth-signature", generateSignature(apiName, userTokenOrNoHp))
	req.Header.Add("cache-control", "no-cache")
	if useHeaderUserToken {
		req.Header.Add("auth-user-token", userTokenOrNoHp)
	}
	if useHeaderContenType {
		req.Header.Add("content-type", "application/json")
	}

	paramLog["header"] = lib.FormatRequest(req)
	paramLog["body"] = string(jsonString)

	res, err := http.DefaultClient.Do(req)
	log.Info(res.StatusCode)
	paramLog["status"] = strconv.FormatUint(uint64(res.StatusCode), 10)

	if res.StatusCode != http.StatusOK {
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			log.Error("Error create log endpoint motion pay")
			log.Error(err.Error())
		}
		return res.StatusCode, "", paramLog["header"], err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		paramLog["response"] = err.Error()
		_, err = models.CreateEndpoint3rdPartyLog(paramLog)
		if err != nil {
			log.Error("Error create log endpoint motion pay")
			log.Error(err.Error())
		}
		log.Error("Error3", err.Error())
		return http.StatusBadGateway, "", paramLog["header"], err
	}
	paramLog["response"] = string(body)
	_, err = models.CreateEndpoint3rdPartyLog(paramLog)
	if err != nil {
		log.Error("Error create log endpoint motion pay")
		log.Error(err.Error())
	}

	return http.StatusOK, string(body), paramLog["header"], nil
}

func generateSignature(apiName string, userTokenOrNoHp string) string {
	sig := config.MERCHANT_ID + "||" +
		config.PARTNER_ID + "||" +
		config.TOKEN + "||" +
		userTokenOrNoHp + "||" +
		apiName

	encryptedPasswordByte := sha256.Sum256([]byte(sig))
	signature := hex.EncodeToString(encryptedPasswordByte[:])

	return signature

}
