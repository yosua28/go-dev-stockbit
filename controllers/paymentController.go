package controllers

import (
	"api/lib"
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func SpinCreateOrder(c echo.Context) error {

	responseData := make(map[string]interface{})

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}
	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	params := make(map[string]string)
	params["reference_code"] = lib.GenerateReference("MDSP", customerKey)
	params["amount"] = "100000"
	params["description"] = "Transaction " + params["reference_code"]
	responseData["reference_code"] = params["reference_code"]
	status, res, err := lib.Spin(params["reference_code"], "CREATE_ORDER", params)
	if err != nil {
		log.Error(status, err.Error())
		return lib.CustomError(status, err.Error(), "Create order failed")
	}
	if status != 200 {
		log.Error(status, "Create order failed")
		return lib.CustomError(status, "Create order failed", "Create order failed")
	}
	var sec map[string]interface{}
	if err := json.Unmarshal([]byte(res), &sec); err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadGateway, err.Error(), "Parsing data failed")
	}
	messageAction := sec["message_action"].(string)
	messageData := sec["message_data"].(map[string]interface{})
	orderID := messageData["order_id"].(string)
	responseData["order_id"] = orderID
	if messageAction == "CREATE_ORDER_SUCCESS" {
		params := make(map[string]string)
		params["order_id"] = orderID
		params["phone"] = lib.Profile.PhoneNumber

		status, _, err := lib.Spin(orderID, "CREATE_OTP", params)
		if err != nil {
			log.Error(status, err.Error())
			return lib.CustomError(status, err.Error(), "Create otp failed")
		}
		if status != 200 {
			log.Error(status, "Create otp failed")
			return lib.CustomError(status, "Create otp failed", "Create otp failed")
		}
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func FMNotif(c echo.Context) error {
	var u map[string]interface{}
	if err := c.Bind(&u); err != nil {
		return lib.CustomError(http.StatusBadRequest, err.Error(), "No data")
	}
	fmt.Println(u)
	return c.JSON(http.StatusOK, "transaksi_valid")
}

func FMThankYou(c echo.Context) error {
	return c.JSON(http.StatusOK, "success payment, thank you")
}
