package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"bytes"
	"crypto/tls"
	"database/sql"
	"html/template"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func CreateTransaction(c echo.Context) error {
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

	typeKeyStr := c.FormValue("type_key")
	if typeKeyStr != "" {
		typeKey, err := strconv.ParseUint(typeKeyStr, 10, 64)
		if err != nil || !(typeKey == 1 || typeKey == 2 || typeKey == 3 || typeKey == 4) {
			log.Error("Wrong input for parameter: type_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: type_key", "Wrong input for parameter: type_key")
		}
	} else {
		log.Error("Missing required parameter: type_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: type_key", "Missing required parameter: type_key")
	}

	redemAllStr := c.FormValue("redem_all")
	flagRedemAll := false
	if redemAllStr != "" {
		redemAll, err := strconv.ParseUint(redemAllStr, 10, 64)
		if err == nil && (redemAll == 1 || redemAll == 0) {
			flagRedemAll = true
		} else {
			log.Error("Wrong input for parameter: redem_all")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: redem_all", "Wrong input for parameter: redem_all")
		}
	}
	parentKeyStr := "NULL"
	if typeKeyStr == "4" {
		parentKeyStr = c.FormValue("parent_key")
		if parentKeyStr != "" {
			parentKey, err := strconv.ParseUint(parentKeyStr, 10, 64)
			if err != nil || parentKey < 1 {
				log.Error("Wrong input for parameter: parent_key")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: parent_key", "Wrong input for parameter: parent_key")
			}
			var transaction models.TrTransaction
			status, err = models.GetTrTransaction(&transaction, parentKeyStr)
			if err != nil {
				log.Error("Parent transaction not found")
				return lib.CustomError(http.StatusBadRequest, "Parent transaction not found", "Parent transaction not found")
			}
			if transaction.TransTypeKey != 3 {
				log.Error("Parent transaction not found")
				return lib.CustomError(http.StatusBadRequest, "Parent transaction not found", "Parent transaction not found")
			}
		} else {
			log.Error("Missing required parameter: parent_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: parent_key", "Missing required parameter: parent_key")
		}
	}

	transUnitStr := c.FormValue("trans_unit")
	if transUnitStr != "" {
		_, err := strconv.ParseFloat(transUnitStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: trans_unit")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_unit", "Wrong input for parameter: trans_unit")
		}
	} else {
		log.Error("Missing required parameter: trans_unit")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_unit", "Missing required parameter: trans_unit")
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
	} else {
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
	} else {
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
		_, err = models.GetAllMsHoliday(&holiday, getParams)
		if (err == nil && len(holiday) < 1) || err != nil {
			params["nav_date"] = now.Format(layout) + " 00:00:00"
			break
		} else {
			i++
		}
	}
	params["entry_mode"] = "140"
	params["trans_amount"] = transAmountStr
	params["trans_unit"] = transUnitStr
	params["trans_type_key"] = typeKeyStr
	if flagRedemAll {
		params["flag_redem_all"] = redemAllStr
	}
	params["parent_key"] = parentKeyStr
	params["flag_newsub"] = "0"

	paramsTr := make(map[string]string)
	paramsTr["customer_key"] = customerKey
	paramsTr["product_key"] = productKeyStr
	paramsTr["trans_type_key"] = "1"
	var transactionDB []models.TrTransaction
	status, err = models.GetAllTrTransaction(&transactionDB, paramsTr)
	if err == nil && len(transactionDB) < 1 {
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

	if typeKeyStr == "1" {
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
	}

	status, err, transactionID := models.CreateTrTransaction(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}
	responseData := make(map[string]string)
	responseData["transaction_key"] = transactionID
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
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
				_, err = os.Stat(config.BasePath + "/images/user/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/transfer/" + filename + extension)
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
		log.Error(err.Error())
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
	params["orderBy"] = "transaction_key"
	params["orderType"] = "DESC"
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
		status, err = models.GetTrTransactionDateRange(&transactionDB, params, "'"+startDate+"'", "'"+endDate+"'")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get transaction data")
			}
		}
		if len(transactionDB) < 1 {
			log.Error("Transaction not found")
			return lib.CustomError(status, "Transaction not found", "Transaction not found")
		}
	} else if trStatus == "process" {
		status, err = models.GetTrTransactionOnProcess(&transactionDB, params)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get transaction data")
			}
		}
		if len(transactionDB) < 1 {
			log.Error("Transaction not found")
			return lib.CustomError(status, "Transaction not found", "Transaction not found")
		}
	} else {
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
		// return lib.CustomError(http.StatusNotFound, "Nav data not found", "Nav data not found")
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
			if product, ok := pData[transaction.ProductKey]; ok {
				data.ProductName = product.ProductNameAlt
			}
			if status, ok := sData[transaction.TransStatusKey]; ok {
				data.TransStatus = *status.StatusCode
			}
			if typ, ok := tData[transaction.TransTypeKey]; ok {
				// data.TransType = *typ.TypeCode
				data.TransType = *typ.TypeDescription
			}
			if nav, ok := nData[transaction.NavDate]; ok {
				data.NavValue = nav.NavValue
			}
			data.TransDate = transaction.TransDate
			data.NavDate = transaction.NavDate
			data.TransAmount = transaction.TransAmount
			data.TransUnit = float32(math.Floor(float64(transaction.TransUnit)*10000) / 10000)
			data.TotalAmount = transaction.TotalAmount
			if transaction.FileUploadDate != nil {
				data.Uploaded = true
				data.DateUploaded = transaction.FileUploadDate
			} else {
				data.Uploaded = false
			}
			if transaction.TransBankKey != nil {
				if bank, ok := bData[*transaction.TransBankKey]; ok {
					data.BankName = &bank.BankName
					data.BankAccName = transaction.TransBankaccName
					data.BankAccNo = transaction.TransBankAccNo
				}
			}
			if transaction.TransTypeKey == 4 {
				if product, ok := pData[transaction.ProductKey]; ok {
					data.ProductIn = &product.ProductName
				}
				if transaction.ParentKey != nil {
					if swot, ok := switchout[*transaction.ParentKey]; ok {
						if product, ok := pData[swot.ProductKey]; ok {
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

func SendEmailPortofolio(c echo.Context) error {

	// Create new PDF generator
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtml.OrientationLandscape)
	pdfg.Grayscale.Set(false)

	// Create a new input page from an URL
	page := wkhtml.NewPage(config.BasePath + "/mail/account-statement-" + strconv.FormatUint(lib.Profile.UserID, 10) + ".html")

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)
	page.Allow.Set(config.BasePath + "/mail/images")

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}
	err = os.MkdirAll(config.BasePath+"/files/"+strconv.FormatUint(lib.Profile.UserID, 10), 0755)
	// Write buffer contents to file on disk
	err = pdfg.WriteFile(config.BasePath + "/files/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/account-statement.pdf")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}
	log.Info("Success create file")
	t := template.New("index-portofolio.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-portofolio.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, ""); err != nil {
		log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", lib.Profile.Email)
	mailer.SetHeader("Subject", "[MNCduit] Laporan Akun Statement")
	mailer.SetBody("text/html", result)
	mailer.Attach(config.BasePath + "/files/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/account-statement.pdf")

	dialer := gomail.NewDialer(
		config.EmailSMTPHost,
		int(config.EmailSMTPPort),
		config.EmailFrom,
		config.EmailFromPassword,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Error(err)
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}
	log.Info("Email sent")
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}

func SendEmailTransaction(c echo.Context) error {

	var err error
	var status int
	params := make(map[string]string)
	trxHistory := make(map[string]interface{})
	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"
	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}
	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	params["customer_key"] = customerKey
	params["rec_status"] = "1"
	params["orderBy"] = "transaction_key"
	params["orderType"] = "DESC"
	productKeyStr := c.FormValue("product_key")
	var product models.MsProduct
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			params["product_key"] = productKeyStr

			status, err = models.GetMsProduct(&product, productKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Product not found")
			}
			trxHistory["ProductName"] = product.ProductNameAlt
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	var transactionDB []models.TrTransaction
	startDate := c.FormValue("start_date")
	startDate += " 00:00:00"
	if startDate != "" {
		date, err := time.Parse(layout, startDate)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: start_date")
		}
		trxHistory["DateStart"] = date.Format(newLayout)
	} else {
		log.Error("Missing required parameter: start_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: start_date", "Missing required parameter: start_date")
	}

	endDate := c.FormValue("end_date")
	endDate += " 23:59:59"
	if endDate != "" {
		date, err := time.Parse(layout, endDate)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Wrong input for parameter: end_date")
		}
		trxHistory["DateEnd"] = date.Format(newLayout)
	} else {
		log.Error("Missing required parameter: end_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: end_date", "Missing required parameter: end_date")
	}

	params["trans_status_key"] = "9"
	status, err = models.GetTrTransactionDateRange(&transactionDB, params, "'"+startDate+"'", "'"+endDate+"'")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get transaction data")
	}
	if len(transactionDB) < 1 {
		log.Error("Transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var statusIDs []string
	var typeIDs []string
	var bankIDs []string
	var transactionIDs []string
	var navDates []string
	for _, transaction := range transactionDB {
		transactionIDs = append(transactionIDs, strconv.FormatUint(transaction.TransactionKey, 10))
		statusIDs = append(statusIDs, strconv.FormatUint(transaction.TransStatusKey, 10))
		typeIDs = append(typeIDs, strconv.FormatUint(transaction.TransTypeKey, 10))
		if transaction.TransBankKey != nil {
			bankIDs = append(bankIDs, strconv.FormatUint(*transaction.TransBankKey, 10))
		}
		navDates = append(navDates, "'"+transaction.NavDate+"'")
	}

	var tcDB []models.TrTransactionConfirmation
	status, err = models.GetTrTransactionConfirmationIn(&tcDB, transactionIDs, "transaction_key")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get TC data")
	}
	if len(tcDB) < 1 {
		log.Error("TC data not found")
		return lib.CustomError(http.StatusNotFound, "TC data not found", "TC data not found")
	}

	tcData := make(map[uint64]models.TrTransactionConfirmation)
	for _, tc := range tcDB {
		tcData[tc.TransactionKey] = tc
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

	var trDatas []map[string]interface{}
	var i uint64 = 0
	for _, transaction := range transactionDB {
		if transaction.TransTypeKey != 3 {
			data := make(map[string]interface{})
			i++
			data["No"] = i
			date, _ := time.Parse(layout, transaction.TransDate)
			data["Date"] = date.Format(newLayout)
			if typ, ok := tData[transaction.TransTypeKey]; ok {
				data["Type"] = *typ.TypeDescription
			}
			fee := transaction.TransFeeAmount + transaction.ChargesFeeAmount + transaction.ServicesFeeAmount
			data["Fee"] = fee
			if tc, ok := tcData[transaction.TransactionKey]; ok {
				data["Unit"] = tc.ConfirmedUnit
				data["Amount"] = tc.ConfirmedAmount
				data["Total"] = tc.ConfirmedAmount + fee
			}
			if nav, ok := nData[transaction.NavDate]; ok {
				data["Nav"] = nav.NavValue
			}

			trDatas = append(trDatas, data)
		}
	}
	trxHistory["Datas"] = trDatas
	t := template.New("transaction-history-template.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/transaction-history-template.html")
	if err != nil {
		log.Println(err)
	}
	f, err := os.Create(config.BasePath + "/mail/transaction-history-" + strconv.FormatUint(lib.Profile.UserID, 10) + ".html")
	if err != nil {
		log.Println("create file: ", err)
	}
	if err := t.Execute(f, trxHistory); err != nil {
		log.Println(err)
	}

	f.Close()

	// Create new PDF generator
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtml.OrientationLandscape)
	pdfg.Grayscale.Set(false)

	// Create a new input page from an URL
	page := wkhtml.NewPage(config.BasePath + "/mail/transaction-history-" + strconv.FormatUint(lib.Profile.UserID, 10) + ".html")

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)
	page.Allow.Set(config.BasePath + "/mail/images")

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}
	err = os.MkdirAll(config.BasePath+"/files/"+strconv.FormatUint(lib.Profile.UserID, 10), 0755)
	// Write buffer contents to file on disk
	err = pdfg.WriteFile(config.BasePath + "/files/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/transaction-history.pdf")
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}
	log.Info("Success create file")

	t = template.New("index-transaction.html")

	t, err = t.ParseFiles(config.BasePath + "/mail/index-transaction.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, ""); err != nil {
		log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", lib.Profile.Email)
	mailer.SetHeader("Subject", "[MNCduit] Transaction History")
	mailer.SetBody("text/html", result)
	mailer.Attach(config.BasePath + "/files/" + strconv.FormatUint(lib.Profile.UserID, 10) + "/transaction-history.pdf")

	dialer := gomail.NewDialer(
		config.EmailSMTPHost,
		int(config.EmailSMTPPort),
		config.EmailFrom,
		config.EmailFromPassword,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Error(err)
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed send email")
	}
	log.Info("Email sent")
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)
}
