package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func CreateSubscribeTransaction(c echo.Context) error {
	var err error
	var status int
	params := make(map[string]string)

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	} 
	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	params["customer_key"] = customerKey
	params["rec_status"] = "1"

	productKeyStr := c.FormValue("product_key")
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			params["product_key"] = productKeyStr
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}
	
	transAmountStr := c.FormValue("trans_amount")
	if transAmountStr != "" {
		_, err := strconv.ParseFloat(transAmountStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: trans_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_amount", "Wrong input for parameter: trans_amount")
		}
	} else {
		log.Error("Missing required parameter: trans_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_amount", "Missing required parameter: trans_amount")
	}

	transFeePercentStr := c.FormValue("trans_fee_percent")
	if transFeePercentStr != "" {
		_, err := strconv.ParseFloat(transFeePercentStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: trans_fee_percent")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_fee_percent", "Wrong input for parameter: trans_fee_percent")
		}
	} else {
		log.Error("Missing required parameter: trans_fee_percent")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_fee_percent", "Missing required parameter: trans_fee_percent")
	}

	transFeeAmountStr := c.FormValue("trans_fee_amount")
	if transFeeAmountStr != "" {
		_, err := strconv.ParseFloat(transFeeAmountStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: trans_fee_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_fee_amount", "Wrong input for parameter: trans_fee_amount")
		}
	} else {
		log.Error("Missing required parameter: trans_fee_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_fee_amount", "Missing required parameter: trans_fee_amount")
	}

	chargesFeeAmountStr := c.FormValue("charges_fee_amount")
	if chargesFeeAmountStr != "" {
		_, err := strconv.ParseFloat(chargesFeeAmountStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: charges_fee_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: charges_fee_amount", "Wrong input for parameter: charges_fee_amount")
		}
	} else {
		log.Error("Missing required parameter: charges_fee_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: charges_fee_amount", "Missing required parameter: charges_fee_amount")
	}

	servicesFeeAmountStr := c.FormValue("services_fee_amount")
	if servicesFeeAmountStr != "" {
		_, err := strconv.ParseFloat(servicesFeeAmountStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: services_fee_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: services_fee_amount", "Wrong input for parameter: services_fee_amount")
		}
	} else {
		log.Error("Missing required parameter: services_fee_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: services_fee_amount", "Missing required parameter: services_fee_amount")
	}

	totalAmountStr := c.FormValue("total_amount")
	if totalAmountStr != "" {
		_, err := strconv.ParseFloat(totalAmountStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: total_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: total_amount", "Wrong input for parameter: total_amount")
		}
	} else {
		log.Error("Missing required parameter: total_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: total_amount", "Missing required parameter: total_amount")
	}

	transSourceStr := c.FormValue("trans_source")
	if transSourceStr != "" {
		_, err := strconv.ParseUint(transSourceStr, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: trans_source")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_source", "Wrong input for parameter: trans_source")
		}
	} else {
		log.Error("Missing required parameter: trans_source")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_source", "Missing required parameter: trans_source")
	}

	riskWaiver := c.FormValue("risk_waiver")
	if riskWaiver != "" {
		if !(riskWaiver == "1" || riskWaiver == "0") {
			log.Error("Wrong input for parameter: risk_waiver")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: risk_waiver", "Wrong input for parameter: risk_waiver")
		}
	} else {
		log.Error("Missing required parameter: risk_waiver")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: risk_waiver", "Missing required parameter: risk_waiver")
	}

	var accKey string
	var trAccountDB []models.TrAccount
	status, err = models.GetAllTrAccount(&trAccountDB, params)
	if len(trAccountDB) > 0 {
		accKey = strconv.FormatUint(trAccountDB[0].AccKey, 10)
	}else{
		params["acc_status"] = "204"
		status, err, accKey = models.CreateTrAccount(params)
		if err != nil {
			log.Error("Failed create account product data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	params = make(map[string]string)
	params["acc_key"] = accKey
	var agentCustomerDB models.MsAgentCustomer
	status, err = models.GetLastAgenCunstomer(&agentCustomerDB, customerKey)
	if err != nil {
		log.Error("Failed get data agent: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	params["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	date := time.Now()
	dateLayout := "2006-01-02 15:04:05"
	dateStr := date.Format(dateLayout)
	
	params["rec_status"] = "1"

	var acaKey string
	var accountAgentDB []models.TrAccountAgent
	status, err = models.GetAllTrAccountAgent(&accountAgentDB, params)
	if len(accountAgentDB) > 0 {
		acaKey = strconv.FormatUint(accountAgentDB[0].AcaKey, 10)
	}else{
		params["eff_date"] = dateStr
		status, err, acaKey = models.CreateTrAccountAgent(params)
		if err != nil {
			log.Error("Failed create account agent data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	params = make(map[string]string)
	params["branch_key"] = "1"
	params["agent_key"] = "1"
	params["customer_key"] = customerKey
	params["product_key"] = productKeyStr
	params["trans_status_key"] = "2"
	params["trans_date"] = dateStr
	params["trans_type_key"] = "1"
	params["trx_code"] = "137"
	params["trans_calc_method"] = "288"

	i := 0
	layout := "2006-01-02"
	for {
		now := time.Now().AddDate(0, 0, i)
		var holiday []models.MsHoliday
		getParams := make(map[string]string)
		getParams["holiday_date"] = now.Format(layout) + " 00:00:00"
		_, err = models.GetAllMsHoliday(&holiday,params)
		if (err == nil && len(holiday) < 1) || err != nil {
			params["nav_date"] = now.Format(layout) + " 00:00:00"
			break
		} else {
			i++
		}
	}
	params["entry_mode"] = "140"
	params["trans_amount"] = transAmountStr
	params["trans_unit"] = "0"
	params["flag_newsub"] = "0"

	paramsTr := make(map[string]string)
	paramsTr["customer_key"] = customerKey
	paramsTr["product_key"] = productKeyStr
	paramsTr["trans_type_key"] = "1"
	var transactionDB []models.TrTransaction
	status, err = models.GetAllTrTransaction(&transactionDB, paramsTr)
	if (err == nil && len(transactionDB) < 1) {
		params["flag_newsub"] = "1"
	}
	
	params["trans_fee_percent"] = transFeePercentStr
	params["trans_fee_amount"] = transFeeAmountStr
	params["charges_fee_amount"] = chargesFeeAmountStr
	params["services_fee_amount"] = servicesFeeAmountStr
	params["total_amount"] = totalAmountStr

	transRemark := c.FormValue("trans_remarks")
	params["trans_remarks"] = transRemark

	params["risk_waiver"] = riskWaiver
	params["trans_source"] = transSourceStr
	params["payment_method"] = "284"
	params["aca_key"] = acaKey
	params["rec_status"] = "1"

	err = os.MkdirAll(config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/transfer", 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("transfer_pic")
		if file != nil {
			if err == nil {
				// Get file extension
				extension := filepath.Ext(file.Filename)
				// Generate filename
				var filename string
				for {
					filename = lib.RandStringBytesMaskImprSrc(20)
					log.Println("Generate filename:", filename)
					var trans []models.TrTransaction
					getParams := make(map[string]string)
					getParams["rec_image1"] = filename + extension
					_, err := models.GetAllTrTransaction(&trans, getParams)
					if (err == nil && len(trans) < 1) || err != nil {
						break
					}
				}
				// Upload image and move to proper directory
				err = lib.UploadImage(file, config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/transfer/"+filename+extension)
				if err != nil {
					log.Println(err)
					return lib.CustomError(http.StatusInternalServerError)
				}
				params["rec_image1"] = filename + extension
				dateLayout := "2006-01-02 15:04:05"
				params["file_upload_date"] = time.Now().Format(dateLayout)
			}
		}
	}

	status, err = models.CreateTrTransaction(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
		
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}


func UploadTransferPic(c echo.Context) error {
	var err error
	var status int
	params := make(map[string]string)
	transactionKeyStr := c.FormValue("transaction_key")
	if transactionKeyStr != "" {
		productKey, err := strconv.ParseUint(transactionKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			params["transaction_key"] = transactionKeyStr
		} else {
			log.Error("Wrong input for parameter: transaction_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
		}
	} else {
		log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: transaction_key", "Missing required parameter: transaction_key")
	}

	var transactionDB []models.TrTransaction
	status, err = models.GetAllTrTransaction(&transactionDB, params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Transaction not found")
	}
	if len(transactionDB) == 0 {
		log.Error("Transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	err = os.MkdirAll(config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/transfer", 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("transfer_pic")
		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Missing required parameter: transfer_pic")
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			// Generate filename
			var filename string
			for {
				filename = lib.RandStringBytesMaskImprSrc(20)
				log.Println("Generate filename:", filename)
				var trans []models.TrTransaction
				getParams := make(map[string]string)
				getParams["rec_image1"] = filename + extension
				_, err = os.Stat(config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/transfer/"+filename+extension)
				if err != nil {
					if os.IsNotExist(err) {
						_, err = models.GetAllTrTransaction(&trans, getParams)
						if (err == nil && len(trans) < 1) || err != nil {
							break
						}
					}
				}
			}
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/images/user/"+strconv.FormatUint(lib.Profile.UserID, 10)+"/transfer/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}
			params["rec_image1"] = filename + extension
			dateLayout := "2006-01-02 15:04:05"
			params["file_upload_date"] = time.Now().Format(dateLayout)
		}
	}

	status, err = models.UpdateTrTransaction(params)
	if err != nil {
		log.Error(err.Error)
		return lib.CustomError(status, err.Error(), "Failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func GetTransactionList(c echo.Context) error {
	var err error
	var status int
	params := make(map[string]string)

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	} 
	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	params["customer_key"] = customerKey
	params["rec_status"] = "1"
	productKeyStr := c.FormValue("product_key")
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			params["product_key"] = productKeyStr
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	}
	
	trStatus := c.FormValue("status")

	var transactionDB []models.TrTransaction
	if trStatus == "posted" {
		startDate := c.FormValue("start_date")
		startDate += " 00:00:00"
		endDate := c.FormValue("end_date")
		endDate += " 23:59:59"
		params["trans_status_key"] = "9"
		status, err = models.GetTrTransactionDateRange(&transactionDB, params, "'"+startDate+"'","'"+ endDate + "'")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get transaction data")
		}
		if len(transactionDB) < 1 {
			log.Error("Transaction not found")
			return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
		}
	}else if trStatus == "process" {
		status, err = models.GetTrTransactionOnProcess(&transactionDB, params)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get transaction data")
		}
		if len(transactionDB) < 1 {
			log.Error("Transaction not found")
			return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
		}
	}else{
		log.Error("Wrong input for parameter: status")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: status", "Wrong input for parameter: status")
	}

	var productIDs []string
	var statusIDs []string
	var typeIDs []string
	var bankIDs []string
	var navDates []string
	switchout := make(map[uint64]models.TrTransaction) 
	for _, transaction := range transactionDB {
		if transaction.TransTypeKey == 3 {
			switchout[transaction.TransactionKey] = transaction
		}
		productIDs = append(productIDs, strconv.FormatUint(transaction.ProductKey, 10))
		statusIDs = append(statusIDs, strconv.FormatUint(transaction.TransStatusKey, 10))
		typeIDs = append(typeIDs, strconv.FormatUint(transaction.TransTypeKey, 10))
		if transaction.TransBankKey != nil {
			bankIDs = append(bankIDs, strconv.FormatUint(*transaction.TransBankKey, 10))
		}
		navDates = append(navDates, "'"+transaction.NavDate+"'")
	}

	var productDB []models.MsProduct
	status, err = models.GetMsProductIn(&productDB, productIDs, "product_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get product data")
	}
	if len(transactionDB) < 1 {
		log.Error("Product not found")
		return lib.CustomError(http.StatusNotFound, "Product not found", "Product not found")
	}
	pData := make(map[uint64]models.MsProduct)
	for _, product := range productDB {
		pData[product.ProductKey] = product
	}	

	var statusDB []models.TrTransactionStatus
	status, err = models.GetMsTransactionStatusIn(&statusDB, statusIDs, "trans_status_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get status data")
	}
	if len(statusDB) < 1 {
		log.Error("Status key not found")
		return lib.CustomError(http.StatusNotFound, "Status key not found", "Status key not found")
	}
	sData := make(map[uint64]models.TrTransactionStatus)
	for _, status := range statusDB {
		sData[status.TransStatusKey] = status
	}	

	var typeDB []models.TrTransactionType
	status, err = models.GetMsTransactionTypeIn(&typeDB, typeIDs, "trans_type_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get type data")
	}
	if len(typeDB) < 1 {
		log.Error("Type key not found")
		return lib.CustomError(http.StatusNotFound, "Type key not found", "Type key not found")
	}
	tData := make(map[uint64]models.TrTransactionType)
	for _, typ := range typeDB {
		tData[typ.TransTypeKey] = typ
	}	

	var bankDB []models.MsBank
	bData := make(map[uint64]models.MsBank)
	if len(bankIDs) > 0 {
		status, err = models.GetMsBankIn(&bankDB, bankIDs, "bank_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get bank data")
		}
		
		for _, bank := range bankDB {
			bData[bank.BankKey] = bank
		}
	}	

	var navDB []models.TrNav
	status, err = models.GetTrNavIn(&navDB, navDates, "nav_date")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get nav data")
	}
	if len(navDB) < 1 {
		log.Error("Nav data not found")
		return lib.CustomError(http.StatusNotFound, "Nav data not found", "Nav data not found")
	}
	nData := make(map[string]models.TrNav)
	for _, nav := range navDB {
		nData[nav.NavDate] = nav
	}
	
	var responseData []models.TrTransactionList
	for _, transaction := range transactionDB {
		if transaction.TransTypeKey != 3 {
			var data models.TrTransactionList

			data.TransactionKey = transaction.TransactionKey
			if product, ok := pData[transaction.ProductKey]; ok{
				data.ProductName = product.ProductName
			}
			if status, ok := sData[transaction.TransStatusKey]; ok{
				data.TransStatus =  *status.StatusCode
			}
			if typ, ok := tData[transaction.TransTypeKey]; ok{
				data.TransType =  *typ.TypeCode
			}
			if nav, ok := nData[transaction.NavDate]; ok{
				data.NavValue =  nav.NavValue
			}
			data.TransDate = transaction.TransDate
			data.NavDate = transaction.NavDate
			data.TransAmount = transaction.TransAmount
			data.TransUnit = transaction.TransUnit
			data.TotalAmount = transaction.TotalAmount
			if transaction.FileUploadDate != nil {
				data.Uploaded = true
				data.DateUploaded = transaction.FileUploadDate
			}else{
				data.Uploaded = false
			}
			if transaction.TransBankKey != nil {
				if bank, ok := bData[*transaction.TransBankKey]; ok{
					data.BankName =  &bank.BankName
					data.BankAccName =  transaction.TransBankaccName
					data.BankAccNo =  transaction.TransBankAccNo
				}
			}
			if transaction.TransTypeKey == 4 {
				if product, ok := pData[transaction.ProductKey]; ok{
					data.ProductIn = &product.ProductName
				}
				if transaction.ParentKey != nil{
					if swot, ok := switchout[*transaction.ParentKey]; ok{
						if product, ok := pData[swot.ProductKey]; ok{
							data.ProductOut = &product.ProductName
						}
					}
				}
			}
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