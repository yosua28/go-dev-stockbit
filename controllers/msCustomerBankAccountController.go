package controllers

import (
	"api/lib"
	"api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func GetAllCustomerBankAccount(c echo.Context) error {
	
	var err error
	var status int

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}

	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)

	var customerBankAcc []models.MsCustomerBankAccount
	customerBankAccParams := make(map[string]string)
	customerBankAccParams["customer_key"] = customerKey
	customerBankAccParams["rec_status"] = "1"
	status, err = models.GetAllMsCustomerBankAccount(&customerBankAcc, customerBankAccParams)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var bankAccountIDs []string
	var priority uint64
	if len(customerBankAcc) > 0 {
		for _, val := range customerBankAcc {
			bankAccountIDs = append(bankAccountIDs, strconv.FormatUint(val.BankAccountKey, 10))
			if val.FlagPriority == 1 {
				priority = val.BankAccountKey
			}
		}
	} 

	var bankAccountDB []models.MsBankAccount
	var bankAccountDatas []interface{}
	_, err = models.GetMsBankAccountIn(&bankAccountDB, bankAccountIDs, "bank_account_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var bankIDs []string
	bankAccDatas := make(map[uint64]models.MsBankAccount)
	for _, val := range bankAccountDB {
		bankIDs = append(bankIDs, strconv.FormatUint(val.BankKey, 10))
		bankAccDatas[val.BankAccountKey] = val
	}
	var bankDB []models.MsBank
	_, err = models.GetMsBankIn(&bankDB, bankIDs, "bank_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	bankDatas := make(map[uint64]string)
	for _, val := range bankDB {
		bankDatas[val.BankKey] = val.BankName
	}

	for _, val := range customerBankAcc {
		bankAccount := make(map[string]interface{})
		bankAcc := bankAccDatas[val.BankAccountKey]
		bankAccount["cust_bankacc_key"] = val.CustBankaccKey
		bankAccount["bank_account_key"] = val.BankAccountKey
		bankAccount["bank_key"] = bankDatas[bankAcc.BankKey]
		bankAccount["account_no"] = bankAcc.AccountNo
		bankAccount["account_holder_name"] = bankAcc.AccountHolderName
		bankAccount["branch_name"] = bankAcc.BranchName
		bankAccount["flag_priority"] = 0
		if val.BankAccountKey == priority {
			bankAccount["flag_priority"] = 1
		}
		bankAccountDatas = append(bankAccountDatas, bankAccount)
	}
	

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = bankAccountDatas

	return c.JSON(http.StatusOK, response)
}

func CustomerBankAccountPriority(c echo.Context) error {
	var err error
	var status int

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}

	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)

	customerBankAccountKeyStr := c.FormValue("customer_bankacc_key")
	var customerBankAccountKey uint64
	if customerBankAccountKeyStr != "" {
		customerBankAccountKey, err = strconv.ParseUint(customerBankAccountKeyStr, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: customer_bankacc_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_bankacc_key", "Wrong input for parameter: customer_bankacc_key")
		}
	} else {
		log.Error("Missing required parameter: customer_bankacc_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_bankacc_key", "Missing required parameter: customer_bankacc_key")
	}

	var customerBankAcc []models.MsCustomerBankAccount
	customerBankAccParams := make(map[string]string)
	customerBankAccParams["customer_key"] = customerKey
	customerBankAccParams["rec_status"] = "1"
	status, err = models.GetAllMsCustomerBankAccount(&customerBankAcc, customerBankAccParams)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	var custBankAccountIDs []string
	var priority uint64
	own := false
	if len(customerBankAcc) > 0 {
		for _, val := range customerBankAcc {
			custBankAccountIDs = append(custBankAccountIDs, strconv.FormatUint(val.CustBankaccKey, 10))
			if val.FlagPriority == 1 {
				priority = val.CustBankaccKey
			}
			if customerBankAccountKey == val.CustBankaccKey {
				own = true
			}
		}
	} else {
		log.Error("Bank account not found")
		return lib.CustomError(http.StatusNotFound, "Bank account not found", "Bank account not found")
	}

	if !own {
		log.Error("Bank account not found")
		return lib.CustomError(http.StatusNotFound, "Bank account not found", "Bank account not found")
	}

	if customerBankAccountKey != priority {
		customerBankAccParams := make(map[string]string)
		customerBankAccParams["flag_priority"] = "1"
		status, err = models.UpdateDataByField(customerBankAccParams, "cust_bankacc_key", strconv.FormatUint(customerBankAccountKey, 10))
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed update data")
		}
		customerBankAccParams["flag_priority"] = "0"
		status, err = models.UpdateDataByField(customerBankAccParams, "cust_bankacc_key", strconv.FormatUint(priority, 10))
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed update data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil

	return c.JSON(http.StatusOK, response)
}