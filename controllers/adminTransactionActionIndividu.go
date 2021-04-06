package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"bytes"
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func GetTransactionSubscription(c echo.Context) error {
	var trType []string
	trType = append(trType, "1")
	return getListAdminTransaction(c, trType)
}
func GetTransactionRedemption(c echo.Context) error {
	var trType []string
	trType = append(trType, "2")
	return getListAdminTransaction(c, trType)
}
func GetTransactionSwitching(c echo.Context) error {
	var trType []string
	trType = append(trType, "3")
	trType = append(trType, "4")
	return getListAdminTransaction(c, trType)
}

func getListAdminTransaction(c echo.Context, trType []string) error {

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

	items := []string{"transaction_key", "branch_key", "agent_key", "customer_key", "product_key", "trans_date", "trans_amount", "trans_bank_key"}

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
		params["orderBy"] = "transaction_key"
		params["orderType"] = "ASC"
	}

	productKey := c.QueryParam("product_key")
	if productKey != "" {
		productKeyCek, err := strconv.ParseUint(productKey, 10, 64)
		if err == nil && productKeyCek > 0 {
			params["product_key"] = productKey
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
		}
	}

	transstatuskey := c.QueryParam("trans_status_key")
	if transstatuskey != "" {
		transstatuskeyCek, err := strconv.ParseUint(transstatuskey, 10, 64)
		if err == nil && transstatuskeyCek > 0 {
			params["trans_status_key"] = transstatuskey
		} else {
			log.Error("Wrong input for parameter: trans_status_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_status_key", "Missing required parameter: trans_status_key")
		}
	}

	navdate := c.QueryParam("nav_date")
	if navdate != "" {
		params["nav_date"] = navdate
	}

	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	if lib.Profile.RoleKey == roleKeyBranchEntry {
		log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["branch_key"] = strBranchKey
		} else {
			log.Error("User Branch haven't Branch")
			return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
		}
	}

	var trTransaction []models.TrTransaction

	status, err = models.AdminGetAllTrTransaction(&trTransaction, limit, offset, noLimit, params, trType, "trans_type_key", true)

	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(trTransaction) < 1 {
		log.Error("transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var branchIds []string
	var agentIds []string
	var customerIds []string
	var productIds []string
	var transTypeIds []string
	var transStatusIds []string
	var transactionIds []string
	for _, tr := range trTransaction {

		if tr.BranchKey != nil {
			if _, ok := lib.Find(branchIds, strconv.FormatUint(*tr.BranchKey, 10)); !ok {
				branchIds = append(branchIds, strconv.FormatUint(*tr.BranchKey, 10))
			}
		}
		if tr.AgentKey != nil {
			if _, ok := lib.Find(agentIds, strconv.FormatUint(*tr.AgentKey, 10)); !ok {
				agentIds = append(agentIds, strconv.FormatUint(*tr.AgentKey, 10))
			}
		}
		if _, ok := lib.Find(transactionIds, strconv.FormatUint(tr.TransactionKey, 10)); !ok {
			transactionIds = append(transactionIds, strconv.FormatUint(tr.TransactionKey, 10))
		}
		if _, ok := lib.Find(customerIds, strconv.FormatUint(tr.CustomerKey, 10)); !ok {
			customerIds = append(customerIds, strconv.FormatUint(tr.CustomerKey, 10))
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(tr.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(tr.ProductKey, 10))
		}
		if _, ok := lib.Find(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10)); !ok {
			transTypeIds = append(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10))
		}
		if _, ok := lib.Find(transStatusIds, strconv.FormatUint(tr.TransStatusKey, 10)); !ok {
			transStatusIds = append(transStatusIds, strconv.FormatUint(tr.TransStatusKey, 10))
		}
	}

	//mapping branch
	var msBranch []models.MsBranch
	if len(branchIds) > 0 {
		status, err = models.GetMsBranchIn(&msBranch, branchIds, "branch_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	branchData := make(map[uint64]models.MsBranch)
	for _, b := range msBranch {
		branchData[b.BranchKey] = b
	}

	//mapping agent
	var msAgent []models.MsAgent
	if len(agentIds) > 0 {
		status, err = models.GetMsAgentIn(&msAgent, agentIds, "agent_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	agentData := make(map[uint64]models.MsAgent)
	for _, a := range msAgent {
		agentData[a.AgentKey] = a
	}

	//mapping customer
	var msCustomer []models.MsCustomer
	if len(customerIds) > 0 {
		status, err = models.GetMsCustomerIn(&msCustomer, customerIds, "customer_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	customerData := make(map[uint64]models.MsCustomer)
	for _, c := range msCustomer {
		customerData[c.CustomerKey] = c
	}

	//mapping product
	var msProduct []models.MsProduct
	if len(productIds) > 0 {
		status, err = models.GetMsProductIn(&msProduct, productIds, "product_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	productData := make(map[uint64]models.MsProduct)
	for _, p := range msProduct {
		productData[p.ProductKey] = p
	}

	//mapping Trans type
	var transactionType []models.TrTransactionType
	if len(transTypeIds) > 0 {
		status, err = models.GetMsTransactionTypeIn(&transactionType, transTypeIds, "trans_type_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	transactionTypeData := make(map[uint64]models.TrTransactionType)
	for _, t := range transactionType {
		transactionTypeData[t.TransTypeKey] = t
	}

	//mapping trans status
	var trTransactionStatus []models.TrTransactionStatus
	if len(transStatusIds) > 0 {
		status, err = models.GetMsTransactionStatusIn(&trTransactionStatus, transStatusIds, "trans_status_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	transStatusData := make(map[uint64]models.TrTransactionStatus)
	for _, ts := range trTransactionStatus {
		transStatusData[ts.TransStatusKey] = ts
	}

	//mapping tc confirmation
	var transConf []models.TrTransactionConfirmation
	if len(transactionIds) > 0 {
		status, err = models.GetTrTransactionConfirmationIn(&transConf, transactionIds, "transaction_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get TC data")
			}
		}
	}
	tcData := make(map[uint64]models.TrTransactionConfirmation)
	for _, tc := range transConf {
		tcData[tc.TransactionKey] = tc
	}

	var responseData []models.AdminTrTransactionInquiryList
	for _, tr := range trTransaction {
		var data models.AdminTrTransactionInquiryList

		data.TransactionKey = tr.TransactionKey
		data.CustomerKey = tr.CustomerKey
		data.ProductKey = tr.ProductKey

		if tr.BranchKey != nil {
			if n, ok := branchData[*tr.BranchKey]; ok {
				data.BranchName = n.BranchName
			}
		}

		if tr.AgentKey != nil {
			if n, ok := agentData[*tr.AgentKey]; ok {
				data.AgentName = n.AgentName
			}
		}

		if n, ok := customerData[tr.CustomerKey]; ok {
			data.CustomerName = n.FullName
		}

		if n, ok := productData[tr.ProductKey]; ok {
			data.ProductName = n.ProductNameAlt
		}

		if n, ok := transStatusData[tr.TransStatusKey]; ok {
			data.TransStatus = *n.StatusCode
		}

		layout := "2006-01-02 15:04:05"
		newLayout := "02 Jan 2006"
		date, _ := time.Parse(layout, tr.TransDate)
		data.TransDate = date.Format(newLayout)
		date, _ = time.Parse(layout, tr.NavDate)
		data.NavDate = date.Format(newLayout)

		if n, ok := transactionTypeData[tr.TransTypeKey]; ok {
			data.TransType = *n.TypeDescription
		}

		if tc, ok := tcData[tr.TransactionKey]; ok {
			data.TransAmount = tc.ConfirmedAmount
			data.TransUnit = tc.ConfirmedUnit
		} else {
			data.TransAmount = tr.TransAmount
			data.TransUnit = tr.TransUnit
		}

		data.TotalAmount = tr.TotalAmount

		responseData = append(responseData, data)
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminGetCountTrTransaction(&countData, params, trType, "trans_type_key")
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
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func CreateTransactionSubscription(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	params := make(map[string]string)

	customerKeyStr := c.FormValue("customer_key")
	var cus models.MsCustomer
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err == nil && customerKey > 0 {
			status, err = models.GetMsCustomer(&cus, customerKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Customer tidak ditemukan")
			}
			if cus.CifSuspendFlag == uint8(1) {
				log.Error("Customer Suspended")
				return lib.CustomError(http.StatusBadRequest, "Customer Suspended", "Customer Suspended")
			}
			params["customer_key"] = customerKeyStr
		} else {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	productKeyStr := c.FormValue("product_key")
	var product models.MsProduct
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			params["product_key"] = productKeyStr
			status, err = models.GetMsProduct(&product, productKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
			}
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
		value, err := decimal.NewFromString(transAmountStr)
		if err != nil {
			log.Error("Wrong input for parameter: trans_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_amount", "Wrong input for parameter: trans_amount")
		}
		if value.Cmp(product.MinSubAmount) == -1 {
			log.Error("sub amount < minimum sub")
			return lib.CustomError(http.StatusBadRequest, "sub amount < minum sub", "Minumum subscription untuk product ini adalah: "+product.MinSubAmount.String())
		}
		if transAmountStr == "0" {
			log.Error("Wrong input for parameter: trans_amount")
			return lib.CustomError(http.StatusBadRequest, "trans_amount harus lebih dari 0", "trans_amount harus lebih dari 0")
		}
		params["trans_amount"] = transAmountStr
	} else {
		log.Error("Missing required parameter: trans_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_amount", "Missing required parameter: trans_amount")
	}

	totalAmountStr := c.FormValue("total_amount")
	if totalAmountStr != "" {
		_, err := strconv.ParseFloat(totalAmountStr, 64)
		if err != nil {
			log.Error("Wrong input for parameter: total_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: total_amount", "Wrong input for parameter: total_amount")
		}
		params["total_amount"] = transAmountStr
	} else {
		log.Error("Missing required parameter: total_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: total_amount", "Missing required parameter: total_amount")
	}

	paymentStr := c.FormValue("payment_method")
	if paymentStr != "" {
		paymentKey, err := strconv.ParseUint(paymentStr, 10, 64)
		if err == nil && paymentKey > 0 {
			params["payment_method"] = paymentStr
		} else {
			log.Error("Missing required parameter: payment_method")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: payment_method", "Missing required parameter: payment_method")
		}
	} else {
		log.Error("Missing required parameter: payment_method")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: payment_method", "Missing required parameter: payment_method")
	}

	bankStr := c.FormValue("bank_transaction")
	if bankStr != "" {
		bankKey, err := strconv.ParseUint(bankStr, 10, 64)
		if err != nil || bankKey == 0 {
			log.Error("Missing required parameter: bank_transaction")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_transaction", "Missing required parameter: bank_transaction")
		}
	} else {
		log.Error("Missing required parameter: bank_transaction")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_transaction", "Missing required parameter: bank_transaction")
	}

	var file *multipart.FileHeader
	file, err = c.FormFile("transfer_pic")
	if file == nil {
		log.Error("Missing required parameter: transfer_pic")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: transfer_pic", "Missing required parameter: transfer_pic")
	}

	transRemark := c.FormValue("trans_remarks")
	params["trans_remarks"] = transRemark

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	//cek tr_account / save
	var accKey string
	paramsAcc := make(map[string]string)
	paramsAcc["customer_key"] = customerKeyStr
	paramsAcc["product_key"] = productKeyStr
	paramsAcc["rec_status"] = "1"
	var trAccountDB []models.TrAccount
	status, err = models.GetAllTrAccount(&trAccountDB, paramsAcc)
	if len(trAccountDB) > 0 {
		params["flag_newsub"] = "0"
		accKey = strconv.FormatUint(trAccountDB[0].AccKey, 10)
		if trAccountDB[0].SubSuspendFlag != nil && *trAccountDB[0].SubSuspendFlag == 1 {
			log.Error("Account suspended")
			return lib.CustomError(status, "Account suspended", "Account suspended")
		}
	} else {
		params["flag_newsub"] = "1"
		paramsAcc["acc_status"] = "204"
		paramsAcc["rec_created_date"] = time.Now().Format(dateLayout)
		paramsAcc["rec_created_by"] = strIDUserLogin
		status, err, accKey = models.CreateTrAccount(paramsAcc)
		if err != nil {
			log.Error("Failed create account product data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}
	//cek tr_account_agent / save
	paramsAccAgent := make(map[string]string)
	paramsAccAgent["acc_key"] = accKey
	var agentCustomerDB models.MsAgentCustomer
	status, err = models.GetLastAgenCunstomer(&agentCustomerDB, customerKeyStr)
	if err != nil {
		log.Error("Failed get data agent: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	paramsAccAgent["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	paramsAccAgent["rec_status"] = "1"

	var acaKey string
	var accountAgentDB []models.TrAccountAgent
	status, err = models.GetAllTrAccountAgent(&accountAgentDB, paramsAccAgent)
	if len(accountAgentDB) > 0 {
		acaKey = strconv.FormatUint(accountAgentDB[0].AcaKey, 10)
	} else {
		paramsCreateAccAgent := make(map[string]string)
		paramsCreateAccAgent["acc_key"] = accKey
		paramsCreateAccAgent["eff_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_by"] = strIDUserLogin
		paramsCreateAccAgent["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
		paramsCreateAccAgent["rec_status"] = "1"
		status, err, acaKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
		if err != nil {
			log.Error("Failed create account agent data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}
	//save tr_transaction
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	if lib.Profile.RoleKey == roleKeyBranchEntry {
		log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["branch_key"] = strBranchKey
		} else {
			params["branch_key"] = "1"
		}
	} else {
		params["branch_key"] = "1"
	}

	params["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	params["trans_status_key"] = "2"
	params["trans_date"] = time.Now().Format(dateLayout)
	params["trans_type_key"] = "1"
	params["trx_code"] = "137"
	params["trans_calc_method"] = "288"
	layout := "2006-01-02"
	now := time.Now()
	nowString := now.Format(layout)
	t, _ := time.Parse(layout, now.AddDate(0, 0, -1).Format(layout))
	dateBursa := SettDate(t, int(1))
	if nowString == dateBursa && (now.Hour() == 12 && now.Minute() > 0) || now.Hour() > 12 {
		t, _ := time.Parse(layout, dateBursa)
		params["nav_date"] = SettDate(t, int(1)) + " 00:00:00"
	} else {
		params["nav_date"] = dateBursa + " 00:00:00"
	}
	params["entry_mode"] = "140"
	params["trans_unit"] = "0"
	params["trans_fee_percent"] = "0"
	params["charges_fee_amount"] = "0"
	var scApp models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&scApp, "BANK_CHARGES")
	if err != nil {
		params["trans_fee_amount"] = "0"
	} else {
		if scApp.AppConfigValue != nil {
			params["trans_fee_amount"] = *scApp.AppConfigValue
		} else {
			params["trans_fee_amount"] = "0"
		}
	}
	var scApp2 models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&scApp2, "SERVICE_CHARGES")
	if err != nil {
		params["services_fee_amount"] = "0"
	} else {
		if scApp.AppConfigValue != nil {
			params["services_fee_amount"] = *scApp.AppConfigValue
		} else {
			params["services_fee_amount"] = "0"
		}
	}
	params["aca_key"] = acaKey
	params["trans_source"] = "141"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strIDUserLogin
	params["risk_waiver"] = "0"

	var riskProfil models.RiskProfilCustomer
	status, err = models.GetRiskProfilCustomer(&riskProfil, customerKeyStr)
	if err != nil {
		if product.RiskProfileKey != nil {
			if riskProfil.RiskProfileKey < *product.RiskProfileKey {
				params["risk_waiver"] = "1"
			}
		}
	}

	var userData models.ScUserLogin
	status, err = models.GetScUserLoginByCustomerKey(&userData, customerKeyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	err = os.MkdirAll(config.BasePath+"/images/user/"+strconv.FormatUint(userData.UserLoginKey, 10)+"/transfer", 0755)
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
				err = lib.UploadImage(file, config.BasePath+"/images/user/"+strconv.FormatUint(userData.UserLoginKey, 10)+"/transfer/"+filename+extension)
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

	status, err, transactionID := models.CreateTrTransaction(params)

	//save tr_transaction_bank_account
	paramsBankTransaction := make(map[string]string)
	paramsBankTransaction["transaction_key"] = transactionID
	paramsBankTransaction["prod_bankacc_key"] = bankStr
	paramsBankTransaction["rec_status"] = "1"

	var customerBankDB []models.MsCustomerBankAccount
	paramCustomerBank := make(map[string]string)
	paramCustomerBank["customer_key"] = customerKeyStr
	paramCustomerBank["flag_priority"] = "1"
	paramCustomerBank["orderBy"] = "cust_bankacc_key"
	paramCustomerBank["orderType"] = "DESC"
	status, err = models.GetAllMsCustomerBankAccount(&customerBankDB, paramCustomerBank)
	if err != nil {
		log.Error(err.Error())
		paramsBankTransaction["cust_bankacc_key"] = "1"
	} else {
		paramsBankTransaction["cust_bankacc_key"] = strconv.FormatUint(customerBankDB[0].CustBankaccKey, 10)
	}
	paramsBankTransaction["rec_created_date"] = time.Now().Format(dateLayout)
	paramsBankTransaction["rec_created_by"] = strIDUserLogin
	status, err = models.CreateTrTransactionBankAccount(paramsBankTransaction)
	if err != nil {
		log.Error(err.Error())
	}

	//create message
	//create push notif
	customerUserLoginKey := strconv.FormatUint(userData.UserLoginKey, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	paramsUserMessage["umessage_recipient_key"] = customerUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	var subject string
	var body string
	var typ string
	if params["flag_newsub"] == "1" {
		typ = "subscription"
		subject = "Subscription sedang Diproses"
		body = "Terima kasih telah melakukan subscription. Kami sedang memproses transaksi kamu."
	} else {
		typ = "topup"
		subject = "Top Up sedang Diproses"
		body = "Terima kasih telah melakukan transaksi top up. Kami sedang memproses transaksi kamu."
	}

	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body

	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strIDUserLogin

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		log.Error("Error create user message")
	} else {
		log.Error("Sukses insert user message")
	}
	lib.CreateNotifCustomerFromAdminByCustomerId(customerKeyStr, subject, body, "TRANSACTION")

	//send email
	params["product_name"] = product.ProductNameAlt
	params["currency"] = strconv.FormatUint(*product.CurrencyKey, 10)
	params["parrent"] = transactionID
	err = mailTransactionManual(typ, params, customerKeyStr, userData.UloginEmail)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func GetTopupData(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true

	customerKeyStr := c.Param("customer_key")
	var cus models.MsCustomer
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err == nil && customerKey > 0 {
			status, err = models.GetMsCustomer(&cus, customerKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Customer tidak ditemukan")
			}
			if cus.CifSuspendFlag == uint8(1) {
				log.Error("Customer Suspended")
				return lib.CustomError(http.StatusBadRequest, "Customer Suspended", "Customer Suspended")
			}
		} else {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	params := make(map[string]string)
	//if user admin role 7 branch
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	if lib.Profile.RoleKey == roleKeyBranchEntry {
		log.Println(lib.Profile)
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["d.branch_key"] = strBranchKey
		} else {
			log.Error("User Branch haven't Branch")
			return lib.CustomError(http.StatusBadRequest, "Wrong User Branch haven't Branch", "Wrong User Branch haven't Branch")
		}
	}
	params["c.customer_key"] = customerKeyStr
	paramsLike := make(map[string]string)
	var customerList []models.CustomerDropdown
	status, err = models.GetCustomerDropdown(&customerList, params, paramsLike)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(customerList) < 1 {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusNotFound, "Customer not found", "Customer not found")
	}

	var customer models.CustomerDropdown
	customer = customerList[0]

	productKeyStr := c.Param("product_key")
	var product models.ProductSubscription
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			status, err = models.AdminGetProductSubscriptionByProductKey(&product, productKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
			}
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	var productResponse models.ProductSubscription
	productResponse.ProductKey = product.ProductKey
	productResponse.FundTypeName = product.FundTypeName
	productResponse.ProductName = product.ProductName
	productResponse.NavDate = product.NavDate
	productResponse.NavValue = product.NavValue.Truncate(2)
	productResponse.PerformD1 = product.PerformD1.Truncate(2)
	productResponse.PerformM1 = product.PerformM1.Truncate(2)
	productResponse.PerformY1 = product.PerformY1.Truncate(2)
	productResponse.ProductImage = product.ProductImage
	productResponse.MinSubAmount = product.MinSubAmount.Truncate(2)
	productResponse.MinRedAmount = product.MinRedAmount.Truncate(2)
	productResponse.MinRedUnit = product.MinRedUnit.Truncate(2)
	productResponse.ProspectusLink = product.ProspectusLink
	productResponse.FfsLink = product.FfsLink
	productResponse.RiskName = product.RiskName
	productResponse.FeeService = product.FeeService.Truncate(0)
	productResponse.FeeTransfer = product.FeeTransfer.Truncate(0)

	var responseData models.AdminTopupData

	responseData.Customer = customer
	responseData.Product = productResponse

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func DeleteTransactionAdmin(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	params := make(map[string]string)

	keyStr := c.FormValue("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		log.Error("Missing required parameter: key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: key", "Missing required parameter: key")
	}

	var transaction models.TrTransaction
	_, err = models.GetTrTransaction(&transaction, keyStr)
	if err != nil {
		log.Error("Transaction not found")
		return lib.CustomError(http.StatusBadRequest, "Transaction not found", "Transaction not found")
	}

	if transaction.TransStatusKey != uint64(2) { //cek sudah diproses belum
		log.Error("Transaction in process, can't delete data.")
		return lib.CustomError(http.StatusBadRequest, "Transaction in process, can't delete data.", "Transaction in process, can't delete data.")
	}

	if transaction.TransSource != nil {
		if *transaction.TransSource != uint64(141) { //cek transaction hanya manual transaksi oleh admin
			log.Error("Can't delete data.")
			return lib.CustomError(http.StatusBadRequest, "Can't delete data.", "Can't delete data.")
		}
	} else {
		log.Error("Can't delete data.")
		return lib.CustomError(http.StatusBadRequest, "Can't delete data.", "Can't delete data.")
	}

	dateLayout := "2006-01-02 15:04:05"
	params["transaction_key"] = keyStr
	params["rec_status"] = "0"
	params["rec_deleted_date"] = time.Now().Format(dateLayout)
	params["rec_deleted_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = nil
	return c.JSON(http.StatusOK, response)

}

func GetCustomerBankAccountRedemption(c echo.Context) error {

	var err error
	var status int

	keyStr := c.Param("customer_key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var customerBankAccountInfo []models.MsCustomerBankAccountInfo
	status, err = models.GetAllMsCustomerBankAccountTransaction(&customerBankAccountInfo, keyStr)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = customerBankAccountInfo

	return c.JSON(http.StatusOK, response)
}

func CreateTransactionRedemption(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	params := make(map[string]string)

	customerKeyStr := c.FormValue("customer_key")
	var cus models.MsCustomer
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err == nil && customerKey > 0 {
			status, err = models.GetMsCustomer(&cus, customerKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Customer tidak ditemukan")
			}
			if cus.CifSuspendFlag == uint8(1) {
				log.Error("Customer Suspended")
				return lib.CustomError(http.StatusBadRequest, "Customer Suspended", "Customer Suspended")
			}
			params["customer_key"] = customerKeyStr
		} else {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	productKeyStr := c.FormValue("product_key")
	var product models.MsProduct
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			params["product_key"] = productKeyStr
			status, err = models.GetMsProduct(&product, productKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
			}
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	transAmountStr := c.FormValue("trans_amount")
	if transAmountStr == "" {
		log.Error("Missing required parameter: trans_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_amount", "Missing required parameter: trans_amount")
	}

	transUnitStr := c.FormValue("trans_unit")
	if transUnitStr == "" {
		log.Error("Missing required parameter: trans_unit")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_unit", "Missing required parameter: trans_unit")
	}

	var productIds []string
	productIds = append(productIds, productKeyStr)
	var productNotAllow []models.ProductCheckAllowRedmSwtching
	status, err = models.CheckProductAllowRedmOrSwitching(&productNotAllow, customerKeyStr, productIds)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data.")
		}
	}

	if len(productNotAllow) > 0 {
		log.Error("Tidak dapat melakukan redemption. Sedang ada proses untuk produk ini.")
		return lib.CustomError(http.StatusBadRequest, "Tidak dapat melakukan redemption. Sedang ada proses untuk produk ini.", "Tidak dapat melakukan redemption. Sedang ada proses untuk produk ini.")
	}

	zero := decimal.NewFromInt(0)
	var balance models.SumBalanceUnit
	status, err = models.GetBalanceUnitByCustomerAndProduct(&balance, customerKeyStr, productKeyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
	} else {
		if balance.Unit.Cmp(zero) == -1 {
			log.Error("Balance Unit 0")
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
		}
	}
	unitTersedia := balance.Unit.Truncate(2)

	var navDB []models.TrNav
	status, err = models.GetLastNavIn(&navDB, productIds)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	nominalTersedia := balance.Unit.Mul(navDB[0].NavValue).Truncate(0)

	metodePerhitungan := c.FormValue("metode_perhitungan")
	if productKeyStr != "" {
		if metodePerhitungan == "1" { //all unit
			params["flag_redempt_all"] = "1"
			params["trans_amount"] = "0"
			value, err := decimal.NewFromString(transUnitStr)
			if err != nil {
				log.Error("Wrong input for parameter: trans_unit")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_unit", "Wrong input for parameter: trans_unit")
			}
			if value.Cmp(zero) == 0 {
				log.Error("trans_unit harus lebih besar dari 0")
				return lib.CustomError(http.StatusBadRequest, "trans_unit harus lebih besar dari 0", "trans_unit harus lebih besar dari 0")
			}
			if value.Cmp(product.MinRedUnit) == -1 {
				log.Error("red unit < minimum red unit ")
				return lib.CustomError(http.StatusBadRequest, "red unit < minum red unit", "Minumum redemption unit untuk product ini adalah: "+product.MinRedUnit.String())
			}

			if value.Cmp(unitTersedia) == 1 {
				log.Error("red unit > unit tersedia")
				return lib.CustomError(http.StatusBadRequest, "red unit > unit tersedia", "Redemption unit tidak boleh lebih besar dari unit tersedia. Unit tersedia saat ini adalah: "+balance.Unit.String())
			}

			params["trans_unit"] = transUnitStr
			params["total_amount"] = "0"
		} else if metodePerhitungan == "2" { //unit penyertaan
			params["trans_amount"] = "0"
			value, err := decimal.NewFromString(transUnitStr)
			if err != nil {
				log.Error("Wrong input for parameter: trans_unit")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_unit", "Wrong input for parameter: trans_unit")
			}
			if value.Cmp(zero) == 0 {
				log.Error("trans_unit harus lebih besar dari 0")
				return lib.CustomError(http.StatusBadRequest, "trans_unit harus lebih besar dari 0", "trans_unit harus lebih besar dari 0")
			}
			if value.Cmp(product.MinRedUnit) == -1 {
				log.Error("red unit < minimum red unit ")
				return lib.CustomError(http.StatusBadRequest, "red unit < minum red unit", "Minumum redemption unit untuk product ini adalah: "+product.MinRedUnit.String())
			}

			if value.Cmp(unitTersedia) == 1 {
				log.Error("red unit > unit tersedia")
				return lib.CustomError(http.StatusBadRequest, "red unit > unit tersedia", "Redemption unit tidak boleh lebih besar dari unit tersedia. Unit tersedia saat ini adalah: "+balance.Unit.String())
			}

			sisaUnitAfterRed := unitTersedia.Sub(value).Truncate(2)
			minSisa := product.MinUnitAfterRed.Truncate(2)

			if sisaUnitAfterRed.Cmp(minSisa) == -1 {
				log.Error("Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All")
				return lib.CustomError(http.StatusBadRequest, "Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All. Sisa unit harus minimal : "+minSisa.String(), "Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All. Sisa unit harus minimal : "+minSisa.String())
			}

			params["trans_unit"] = transUnitStr
			params["total_amount"] = "0"
		} else if metodePerhitungan == "3" { //Nominal
			params["trans_unit"] = "0"
			value, err := decimal.NewFromString(transAmountStr)
			if err != nil {
				log.Error("Wrong input for parameter: trans_amount")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_amount", "Wrong input for parameter: trans_amount")
			}
			if value.Cmp(zero) == 0 {
				log.Error("trans_amount harus lebih besar dari 0")
				return lib.CustomError(http.StatusBadRequest, "trans_amount harus lebih besar dari 0", "trans_amount harus lebih besar dari 0")
			}
			if value.Cmp(product.MinRedAmount) == -1 {
				log.Error("red amount < minimum red amount ")
				return lib.CustomError(http.StatusBadRequest, "red amount < minum red amount", "Minumum redemption amount untuk product ini adalah: "+product.MinRedAmount.String())
			}
			if value.Cmp(nominalTersedia) == 1 {
				log.Error("red nominal > nominal tersedia")
				return lib.CustomError(http.StatusBadRequest, "red amount > nominal amount tersedia", "Redemption amount tidak boleh lebih besar dari amount tersedia. Amount tersedia saat ini adalah: "+nominalTersedia.String())
			}

			unitTerpakai := value.Div(navDB[0].NavValue).Truncate(2)
			sisaUnitAfterRed := unitTersedia.Sub(unitTerpakai).Truncate(2)
			minSisa := product.MinUnitAfterRed.Truncate(2)

			if sisaUnitAfterRed.Cmp(minSisa) == -1 {
				log.Error("Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All")
				return lib.CustomError(http.StatusBadRequest, "Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All. Sisa unit harus minimal : "+minSisa.String(), "Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All. Sisa unit harus minimal : "+minSisa.String())
			}

			params["trans_amount"] = transAmountStr
			params["total_amount"] = transAmountStr
		} else {
			log.Error("Missing required parameter: metode_perhitungan")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: metode_perhitungan", "Missing required parameter: metode_perhitungan")
		}
	} else {
		log.Error("Missing required parameter: metode_perhitungan")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: metode_perhitungan", "Missing required parameter: metode_perhitungan")
	}

	bankStr := c.FormValue("bank_redemption")
	if bankStr != "" {
		bankKey, err := strconv.ParseUint(bankStr, 10, 64)
		if err != nil || bankKey == 0 {
			log.Error("Missing required parameter: bank_redemption")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_redemption", "Missing required parameter: bank_redemption")
		}
	} else {
		log.Error("Missing required parameter: bank_redemption")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: bank_redemption", "Missing required parameter: bank_redemption")
	}

	transRemark := c.FormValue("trans_remarks")
	params["trans_remarks"] = transRemark

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	//cek tr_account / save
	var accKey string
	paramsAcc := make(map[string]string)
	paramsAcc["customer_key"] = customerKeyStr
	paramsAcc["product_key"] = productKeyStr
	paramsAcc["rec_status"] = "1"
	var trAccountDB []models.TrAccount
	status, err = models.GetAllTrAccount(&trAccountDB, paramsAcc)
	if len(trAccountDB) > 0 {
		accKey = strconv.FormatUint(trAccountDB[0].AccKey, 10)
		if trAccountDB[0].RedSuspendFlag != nil && *trAccountDB[0].RedSuspendFlag == 1 {
			log.Error("Account suspended")
			return lib.CustomError(status, "Account suspended", "Account suspended")
		}
	} else {
		paramsAcc["acc_status"] = "204"
		paramsAcc["rec_created_date"] = time.Now().Format(dateLayout)
		paramsAcc["rec_created_by"] = strIDUserLogin
		status, err, accKey = models.CreateTrAccount(paramsAcc)
		if err != nil {
			log.Error("Failed create account product data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}
	//cek tr_account_agent / save
	paramsAccAgent := make(map[string]string)
	paramsAccAgent["acc_key"] = accKey
	var agentCustomerDB models.MsAgentCustomer
	status, err = models.GetLastAgenCunstomer(&agentCustomerDB, customerKeyStr)
	if err != nil {
		log.Error("Failed get data agent: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	paramsAccAgent["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	paramsAccAgent["rec_status"] = "1"

	var acaKey string
	var accountAgentDB []models.TrAccountAgent
	status, err = models.GetAllTrAccountAgent(&accountAgentDB, paramsAccAgent)
	if len(accountAgentDB) > 0 {
		acaKey = strconv.FormatUint(accountAgentDB[0].AcaKey, 10)
	} else {
		paramsCreateAccAgent := make(map[string]string)
		paramsCreateAccAgent["acc_key"] = accKey
		paramsCreateAccAgent["eff_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_by"] = strIDUserLogin
		paramsCreateAccAgent["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
		paramsCreateAccAgent["rec_status"] = "1"
		status, err, acaKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
		if err != nil {
			log.Error("Failed create account agent data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}
	//save tr_transaction
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	if lib.Profile.RoleKey == roleKeyBranchEntry {
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["branch_key"] = strBranchKey
		} else {
			params["branch_key"] = "1"
		}
	} else {
		params["branch_key"] = "1"
	}

	params["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	params["trans_status_key"] = "2"
	params["trans_date"] = time.Now().Format(dateLayout)
	params["trans_type_key"] = "2"
	params["trx_code"] = "137"
	params["payment_method"] = "284"
	params["trans_calc_method"] = "288"
	layout := "2006-01-02"
	now := time.Now()
	nowString := now.Format(layout)
	t, _ := time.Parse(layout, now.AddDate(0, 0, -1).Format(layout))
	dateBursa := SettDate(t, int(1))
	if nowString == dateBursa && (now.Hour() == 12 && now.Minute() > 0) || now.Hour() > 12 {
		t, _ := time.Parse(layout, dateBursa)
		params["nav_date"] = SettDate(t, int(1)) + " 00:00:00"
	} else {
		params["nav_date"] = dateBursa + " 00:00:00"
	}
	params["entry_mode"] = "140"
	params["trans_fee_percent"] = "0"
	params["charges_fee_amount"] = "0"
	var scApp models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&scApp, "BANK_CHARGES")
	if err != nil {
		params["trans_fee_amount"] = "0"
	} else {
		if scApp.AppConfigValue != nil {
			params["trans_fee_amount"] = *scApp.AppConfigValue
		} else {
			params["trans_fee_amount"] = "0"
		}
	}
	var scApp2 models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&scApp2, "SERVICE_CHARGES")
	if err != nil {
		params["services_fee_amount"] = "0"
	} else {
		if scApp.AppConfigValue != nil {
			params["services_fee_amount"] = *scApp.AppConfigValue
		} else {
			params["services_fee_amount"] = "0"
		}
	}
	params["aca_key"] = acaKey
	params["trans_source"] = "141"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strIDUserLogin
	params["risk_waiver"] = "0"

	var userData models.ScUserLogin
	status, err = models.GetScUserLoginByCustomerKey(&userData, customerKeyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	status, err, transactionID := models.CreateTrTransaction(params)

	//save tr_transaction_bank_account
	paramsBankTransaction := make(map[string]string)
	paramsBankTransaction["transaction_key"] = transactionID
	paramsBankTransaction["cust_bankacc_key"] = bankStr
	paramsBankTransaction["rec_status"] = "1"

	purpose := "270"
	paramsProBankAcc := make(map[string]string)
	paramsProBankAcc["bank_account_purpose"] = purpose
	paramsProBankAcc["product_key"] = productKeyStr
	var productBankDB []models.MsProductBankAccount
	status, err = models.GetAllMsProductBankAccount(&productBankDB, paramsProBankAcc)
	if err != nil {
		log.Error(err.Error())
		paramsBankTransaction["prod_bankacc_key"] = "1"
	} else {
		paramsBankTransaction["prod_bankacc_key"] = strconv.FormatUint(productBankDB[0].ProdBankaccKey, 10)
	}
	paramsBankTransaction["rec_created_date"] = time.Now().Format(dateLayout)
	paramsBankTransaction["rec_created_by"] = strIDUserLogin
	status, err = models.CreateTrTransactionBankAccount(paramsBankTransaction)
	if err != nil {
		log.Error(err.Error())
	}

	//create message
	customerUserLoginKey := strconv.FormatUint(userData.UserLoginKey, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	paramsUserMessage["umessage_recipient_key"] = customerUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"

	subject := "Redemption sedang Diproses"
	body := "Redemption kamu telah kami terima. Kami akan memproses transaksi kamu."

	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body

	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strIDUserLogin

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		log.Error("Error create user message")
	} else {
		log.Error("Sukses insert user message")
	}

	//create push notif
	lib.CreateNotifCustomerFromAdminByCustomerId(customerKeyStr, subject, body, "TRANSACTION")

	//send email
	var customerBankAccountInfo models.MsCustomerBankAccountInfo
	status, err = models.GetMsCustomerBankAccountTransactionByCustBankaccKey(&customerBankAccountInfo, bankStr)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	params["product_name"] = product.ProductNameAlt
	params["currency"] = strconv.FormatUint(*product.CurrencyKey, 10)
	params["parrent"] = transactionID
	params["BankName"] = customerBankAccountInfo.BankName
	params["AccountNo"] = customerBankAccountInfo.AccountNo
	params["AccountHolderName"] = customerBankAccountInfo.AccountName
	if customerBankAccountInfo.BranchName != nil {
		params["BranchName"] = *customerBankAccountInfo.BranchName
	} else {
		params["BranchName"] = "-"
	}
	err = mailTransactionManual("redemption", params, customerKeyStr, userData.UloginEmail)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func CreateTransactionSwitching(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	params := make(map[string]string)
	paramsSwIn := make(map[string]string)

	customerKeyStr := c.FormValue("customer_key")
	var cus models.MsCustomer
	if customerKeyStr != "" {
		customerKey, err := strconv.ParseUint(customerKeyStr, 10, 64)
		if err == nil && customerKey > 0 {
			status, err = models.GetMsCustomer(&cus, customerKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Customer tidak ditemukan")
			}
			if cus.CifSuspendFlag == uint8(1) {
				log.Error("Customer Suspended")
				return lib.CustomError(http.StatusBadRequest, "Customer Suspended", "Customer Suspended")
			}
			params["customer_key"] = customerKeyStr
			paramsSwIn["customer_key"] = customerKeyStr
		} else {
			log.Error("Wrong input for parameter: customer_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_key", "Wrong input for parameter: customer_key")
		}
	} else {
		log.Error("Missing required parameter: customer_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_key", "Missing required parameter: customer_key")
	}

	productKeyStr := c.FormValue("product_from")
	var product models.MsProduct
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
			params["product_key"] = productKeyStr
			status, err = models.GetMsProduct(&product, productKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
			}
		} else {
			log.Error("Wrong input for parameter: product_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_key", "Wrong input for parameter: product_key")
		}
	} else {
		log.Error("Missing required parameter: product_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_key", "Missing required parameter: product_key")
	}

	transAmountStr := c.FormValue("trans_amount")
	if transAmountStr == "" {
		log.Error("Missing required parameter: trans_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_amount", "Missing required parameter: trans_amount")
	}

	transUnitStr := c.FormValue("trans_unit")
	if transUnitStr == "" {
		log.Error("Missing required parameter: trans_unit")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_unit", "Missing required parameter: trans_unit")
	}

	var productIds []string
	productIds = append(productIds, productKeyStr)
	var productNotAllow []models.ProductCheckAllowRedmSwtching
	status, err = models.CheckProductAllowRedmOrSwitching(&productNotAllow, customerKeyStr, productIds)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			log.Error(err.Error())
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data.")
		}
	}

	if len(productNotAllow) > 0 {
		log.Error("Tidak dapat melakukan switching. Sedang ada proses untuk produk ini.")
		return lib.CustomError(http.StatusBadRequest, "Tidak dapat melakukan switching. Sedang ada proses untuk produk ini.", "Tidak dapat melakukan switching. Sedang ada proses untuk produk ini.")
	}

	productToKeyStr := c.FormValue("product_to")
	var productTo models.MsProduct
	if productToKeyStr != "" {
		productToKey, err := strconv.ParseUint(productToKeyStr, 10, 64)
		if err == nil && productToKey > 0 {
			paramsSwIn["product_key"] = productToKeyStr
			status, err = models.GetMsProduct(&productTo, productToKeyStr)
			if err != nil {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest, err.Error(), "Product Tujuan tidak ditemukan")
			}
		} else {
			log.Error("Wrong input for parameter: product_to")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_to", "Wrong input for parameter: product_to")
		}
	} else {
		log.Error("Missing required parameter: product_to")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_to", "Missing required parameter: product_to")
	}

	zero := decimal.NewFromInt(0)
	var balance models.SumBalanceUnit
	status, err = models.GetBalanceUnitByCustomerAndProduct(&balance, customerKeyStr, productKeyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
	} else {
		if balance.Unit.Cmp(zero) == -1 {
			log.Error("Balance Unit 0")
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Product tidak ditemukan")
		}
	}
	unitTersedia := balance.Unit.Truncate(2)

	//NAV Product FROM
	var navDB []models.TrNav
	status, err = models.GetLastNavIn(&navDB, productIds)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	nominalTersedia := balance.Unit.Mul(navDB[0].NavValue).Truncate(0)

	//validasi product from
	metodePerhitungan := c.FormValue("metode_perhitungan")
	if productKeyStr != "" {
		if metodePerhitungan == "1" { //all unit
			params["flag_redempt_all"] = "1"
			params["trans_amount"] = "0"
			value, err := decimal.NewFromString(transUnitStr)
			if err != nil {
				log.Error("Wrong input for parameter: trans_unit")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_unit", "Wrong input for parameter: trans_unit")
			}
			if value.Cmp(zero) == 0 {
				log.Error("trans_unit harus lebih besar dari 0")
				return lib.CustomError(http.StatusBadRequest, "trans_unit harus lebih besar dari 0", "trans_unit harus lebih besar dari 0")
			}
			if value.Cmp(product.MinRedUnit) == -1 {
				log.Error("switching unit < minimum switching unit ")
				return lib.CustomError(http.StatusBadRequest, "switching unit < minum switching unit", "Minumum Switching unit untuk product ini adalah: "+product.MinRedUnit.String())
			}

			if value.Cmp(unitTersedia) == 1 {
				log.Error("switching unit > unit tersedia")
				return lib.CustomError(http.StatusBadRequest, "switching unit > unit tersedia", "Switching unit tidak boleh lebih besar dari unit tersedia. Unit tersedia saat ini adalah: "+balance.Unit.String())
			}

			params["trans_unit"] = transUnitStr
			params["total_amount"] = "0"
		} else if metodePerhitungan == "2" { //unit penyertaan
			params["trans_amount"] = "0"
			value, err := decimal.NewFromString(transUnitStr)
			if err != nil {
				log.Error("Wrong input for parameter: trans_unit")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_unit", "Wrong input for parameter: trans_unit")
			}
			if value.Cmp(zero) == 0 {
				log.Error("trans_unit harus lebih besar dari 0")
				return lib.CustomError(http.StatusBadRequest, "trans_unit harus lebih besar dari 0", "trans_unit harus lebih besar dari 0")
			}
			if value.Cmp(product.MinRedUnit) == -1 {
				log.Error("switching unit < minimum switching unit ")
				return lib.CustomError(http.StatusBadRequest, "switching unit < minum switching unit", "Minumum switching unit untuk product ini adalah: "+product.MinRedUnit.String())
			}

			if value.Cmp(unitTersedia) == 1 {
				log.Error("switching unit > unit tersedia")
				return lib.CustomError(http.StatusBadRequest, "switching unit > unit tersedia", "Switching unit tidak boleh lebih besar dari unit tersedia. Unit tersedia saat ini adalah: "+balance.Unit.String())
			}

			sisaUnitAfterRed := unitTersedia.Sub(value).Truncate(2)
			minSisa := product.MinUnitAfterRed.Truncate(2)

			if sisaUnitAfterRed.Cmp(minSisa) == -1 {
				log.Error("Sisa unit setelah switching kurang dari minimal unit, Silahkan switch All")
				return lib.CustomError(http.StatusBadRequest, "Sisa unit setelah switching kurang dari minimal unit, Silahkan switching All. Sisa unit harus minimal : "+minSisa.String(), "Sisa unit setelah switching kurang dari minimal unit, Silahkan switching All. Sisa unit harus minimal : "+minSisa.String())
			}

			params["trans_unit"] = transUnitStr
			params["total_amount"] = "0"
		} else if metodePerhitungan == "3" { //Nominal
			params["trans_unit"] = "0"
			value, err := decimal.NewFromString(transAmountStr)
			if err != nil {
				log.Error("Wrong input for parameter: trans_amount")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_amount", "Wrong input for parameter: trans_amount")
			}
			if value.Cmp(zero) == 0 {
				log.Error("trans_amount harus lebih besar dari 0")
				return lib.CustomError(http.StatusBadRequest, "trans_amount harus lebih besar dari 0", "trans_amount harus lebih besar dari 0")
			}
			if value.Cmp(product.MinRedAmount) == -1 {
				log.Error("switching amount < minimum switching amount ")
				return lib.CustomError(http.StatusBadRequest, "switching amount < minum switching amount", "Minumum switching amount untuk product ini adalah: "+product.MinRedAmount.String())
			}
			if value.Cmp(nominalTersedia) == 1 {
				log.Error("red nominal > nominal tersedia")
				return lib.CustomError(http.StatusBadRequest, "switching amount > nominal amount tersedia", "Switching amount tidak boleh lebih besar dari amount tersedia. Amount tersedia saat ini adalah: "+nominalTersedia.String())
			}

			unitTerpakai := value.Div(navDB[0].NavValue).Truncate(2)
			sisaUnitAfterRed := unitTersedia.Sub(unitTerpakai).Truncate(2)
			minSisa := product.MinUnitAfterRed.Truncate(2)

			if sisaUnitAfterRed.Cmp(minSisa) == -1 {
				log.Error("Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All")
				return lib.CustomError(http.StatusBadRequest, "Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All. Sisa unit harus minimal : "+minSisa.String(), "Sisa unit setelah redemption kurang dari minimal unit, Silahkan redemption All. Sisa unit harus minimal : "+minSisa.String())
			}
			params["trans_amount"] = transAmountStr
			params["total_amount"] = transAmountStr
		} else {
			log.Error("Missing required parameter: metode_perhitungan")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: metode_perhitungan", "Missing required parameter: metode_perhitungan")
		}
	} else {
		log.Error("Missing required parameter: metode_perhitungan")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: metode_perhitungan", "Missing required parameter: metode_perhitungan")
	}

	//NAV Product TO
	var productToIds []string
	productToIds = append(productToIds, productToKeyStr)
	var navProductToDB []models.TrNav
	status, err = models.GetLastNavIn(&navProductToDB, productToIds)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}

	//validasi product to
	if metodePerhitungan == "3" { //nominal
		minSubNewProd := productTo.MinSubAmount.Truncate(0)
		value, err := decimal.NewFromString(transAmountStr)
		if err != nil {
			log.Error("Wrong input for parameter: trans_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_amount", "Wrong input for parameter: trans_amount")
		}
		jumlahSub := value.Truncate(0)

		if jumlahSub.Cmp(minSubNewProd) == -1 {
			log.Error("switching nominal < minimal switching product tujuan")
			return lib.CustomError(http.StatusBadRequest, "switching nominal < minimal switching product tujuan", "Switching amount tidak boleh kurang dari minimal switching product tujuan. Product tujuan memiliki minimal switching : "+productTo.MinSubAmount.String())
		}
	} else { //unit penyertaan/unit all
		value, err := decimal.NewFromString(transUnitStr)
		if err != nil {
			log.Error("Wrong input for parameter: trans_unit")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_unit", "Wrong input for parameter: trans_unit")
		}

		minSubNewProd := productTo.MinSubAmount.Truncate(0)
		jumlahSubNominal := value.Mul(navDB[0].NavValue)

		if jumlahSubNominal.Cmp(minSubNewProd) == -1 {
			log.Error("switching nominal < minimal switching product tujuan")
			return lib.CustomError(http.StatusBadRequest, "switching nominal < minimal switching product tujuan", "Switching amount tidak boleh kurang dari minimal switching product baru. Product tujuan memiliki minimal switching : "+productTo.MinSubAmount.String())
		}

	}

	transRemark := c.FormValue("trans_remarks")
	params["trans_remarks"] = transRemark
	paramsSwIn["trans_remarks"] = transRemark

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	//SAVE PRODUCT FROM
	//cek tr_account / save
	var accKey string
	paramsAcc := make(map[string]string)
	paramsAcc["customer_key"] = customerKeyStr
	paramsAcc["product_key"] = productKeyStr
	paramsAcc["rec_status"] = "1"
	var trAccountDB []models.TrAccount
	status, err = models.GetAllTrAccount(&trAccountDB, paramsAcc)
	if len(trAccountDB) > 0 {
		accKey = strconv.FormatUint(trAccountDB[0].AccKey, 10)
		if trAccountDB[0].RedSuspendFlag != nil && *trAccountDB[0].RedSuspendFlag == 1 {
			log.Error("Product Asal suspended")
			return lib.CustomError(status, "Product Asal suspended", "Product Asal suspended")
		}
	} else {
		paramsAcc["acc_status"] = "204"
		paramsAcc["rec_created_date"] = time.Now().Format(dateLayout)
		paramsAcc["rec_created_by"] = strIDUserLogin
		status, err, accKey = models.CreateTrAccount(paramsAcc)
		if err != nil {
			log.Error("Failed create account product data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}
	//cek tr_account_agent / save
	paramsAccAgent := make(map[string]string)
	paramsAccAgent["acc_key"] = accKey
	var agentCustomerDB models.MsAgentCustomer
	status, err = models.GetLastAgenCunstomer(&agentCustomerDB, customerKeyStr)
	if err != nil {
		log.Error("Failed get data agent: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	paramsAccAgent["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	paramsAccAgent["rec_status"] = "1"

	var acaKey string
	var accountAgentDB []models.TrAccountAgent
	status, err = models.GetAllTrAccountAgent(&accountAgentDB, paramsAccAgent)
	if len(accountAgentDB) > 0 {
		acaKey = strconv.FormatUint(accountAgentDB[0].AcaKey, 10)
	} else {
		paramsCreateAccAgent := make(map[string]string)
		paramsCreateAccAgent["acc_key"] = accKey
		paramsCreateAccAgent["eff_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_by"] = strIDUserLogin
		paramsCreateAccAgent["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
		paramsCreateAccAgent["rec_status"] = "1"
		status, err, acaKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
		if err != nil {
			log.Error("Failed create account agent data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}
	//save tr_transaction
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	if lib.Profile.RoleKey == roleKeyBranchEntry {
		if lib.Profile.BranchKey != nil {
			strBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
			params["branch_key"] = strBranchKey
			paramsSwIn["branch_key"] = strBranchKey
		} else {
			params["branch_key"] = "1"
			paramsSwIn["branch_key"] = "1"
		}
	} else {
		params["branch_key"] = "1"
		paramsSwIn["branch_key"] = "1"
	}

	params["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	params["trans_status_key"] = "2"
	params["trans_date"] = time.Now().Format(dateLayout)
	params["trans_type_key"] = "3"
	params["trx_code"] = "137"
	params["payment_method"] = "284"
	params["trans_calc_method"] = "288"

	layout := "2006-01-02"
	now := time.Now()
	nowString := now.Format(layout)
	t, _ := time.Parse(layout, now.AddDate(0, 0, -1).Format(layout))
	dateBursa := SettDate(t, int(1))
	if nowString == dateBursa && (now.Hour() == 12 && now.Minute() > 0) || now.Hour() > 12 {
		t, _ := time.Parse(layout, dateBursa)
		params["nav_date"] = SettDate(t, int(1)) + " 00:00:00"
		paramsSwIn["nav_date"] = SettDate(t, int(1)) + " 00:00:00"
	} else {
		params["nav_date"] = dateBursa + " 00:00:00"
		paramsSwIn["nav_date"] = dateBursa + " 00:00:00"
	}
	params["entry_mode"] = "140"
	params["trans_fee_percent"] = "0"
	params["charges_fee_amount"] = "0"

	var scApp models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&scApp, "BANK_CHARGES")
	if err != nil {
		params["trans_fee_amount"] = "0"
		paramsSwIn["trans_fee_amount"] = "0"
	} else {
		if scApp.AppConfigValue != nil {
			params["trans_fee_amount"] = *scApp.AppConfigValue
			paramsSwIn["trans_fee_amount"] = *scApp.AppConfigValue
		} else {
			params["trans_fee_amount"] = "0"
			paramsSwIn["trans_fee_amount"] = "0"
		}
	}
	var scApp2 models.ScAppConfig
	status, err = models.GetScAppConfigByCode(&scApp2, "SERVICE_CHARGES")
	if err != nil {
		params["services_fee_amount"] = "0"
		paramsSwIn["services_fee_amount"] = "0"
	} else {
		if scApp.AppConfigValue != nil {
			params["services_fee_amount"] = *scApp.AppConfigValue
			paramsSwIn["services_fee_amount"] = *scApp.AppConfigValue
		} else {
			params["services_fee_amount"] = "0"
			paramsSwIn["services_fee_amount"] = "0"
		}
	}
	params["aca_key"] = acaKey
	params["trans_source"] = "141"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strIDUserLogin
	params["risk_waiver"] = "0"
	params["flag_newsub"] = "0"

	var userData models.ScUserLogin
	status, err = models.GetScUserLoginByCustomerKey(&userData, customerKeyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	status, err, transactionID := models.CreateTrTransaction(params)

	//SAVE PRODUCT TO
	//cek tr_account / save
	var accNewKey string
	paramsNewProdAcc := make(map[string]string)
	paramsNewProdAcc["customer_key"] = customerKeyStr
	paramsNewProdAcc["product_key"] = productToKeyStr
	paramsNewProdAcc["rec_status"] = "1"
	var trAccountNewDB []models.TrAccount
	status, err = models.GetAllTrAccount(&trAccountNewDB, paramsNewProdAcc)
	if len(trAccountNewDB) > 0 {
		accNewKey = strconv.FormatUint(trAccountNewDB[0].AccKey, 10)
		if trAccountNewDB[0].SubSuspendFlag != nil && *trAccountNewDB[0].SubSuspendFlag == 1 {
			log.Error("Product Tujuan suspended")
			return lib.CustomError(status, "Product Tujuan suspended", "Product Tujuan suspended")
		}
	} else {
		paramsNewProdAcc["acc_status"] = "204"
		paramsNewProdAcc["rec_created_date"] = time.Now().Format(dateLayout)
		paramsNewProdAcc["rec_created_by"] = strIDUserLogin
		status, err, accNewKey = models.CreateTrAccount(paramsNewProdAcc)
		if err != nil {
			log.Error("Failed create account product data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}
	//cek tr_account_agent / save
	paramsNewProdAccAgent := make(map[string]string)
	paramsNewProdAccAgent["acc_key"] = accNewKey
	paramsNewProdAccAgent["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	paramsNewProdAccAgent["rec_status"] = "1"

	var acaNewProdKey string
	var accountNewProdAgentDB []models.TrAccountAgent
	status, err = models.GetAllTrAccountAgent(&accountNewProdAgentDB, paramsNewProdAccAgent)
	if len(accountNewProdAgentDB) > 0 {
		acaNewProdKey = strconv.FormatUint(accountNewProdAgentDB[0].AcaKey, 10)
	} else {
		paramsCreateAccAgent := make(map[string]string)
		paramsCreateAccAgent["acc_key"] = accKey
		paramsCreateAccAgent["eff_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_date"] = time.Now().Format(dateLayout)
		paramsCreateAccAgent["rec_created_by"] = strIDUserLogin
		paramsCreateAccAgent["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
		paramsCreateAccAgent["rec_status"] = "1"
		status, err, acaNewProdKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
		if err != nil {
			log.Error("Failed create account agent data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	paramsSwIn["agent_key"] = strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	paramsSwIn["trans_status_key"] = "2"
	paramsSwIn["trans_date"] = time.Now().Format(dateLayout)
	paramsSwIn["trans_type_key"] = "4"
	paramsSwIn["trx_code"] = "137"
	paramsSwIn["payment_method"] = "284"
	paramsSwIn["trans_calc_method"] = "288"

	paramsSwIn["entry_mode"] = "140"
	paramsSwIn["trans_fee_percent"] = "0"
	paramsSwIn["charges_fee_amount"] = "0"

	paramsSwIn["aca_key"] = acaNewProdKey
	paramsSwIn["trans_source"] = "141"
	paramsSwIn["rec_status"] = "1"
	paramsSwIn["rec_created_date"] = time.Now().Format(dateLayout)
	paramsSwIn["rec_created_by"] = strIDUserLogin
	paramsSwIn["risk_waiver"] = "0"
	paramsSwIn["parent_key"] = transactionID
	paramsSwIn["trans_amount"] = "0"
	paramsSwIn["trans_unit"] = "0"
	paramsSwIn["total_amount"] = "0"
	paramsSwIn["flag_newsub"] = "0"

	status, err, _ = models.CreateTrTransaction(paramsSwIn)

	//create message
	customerUserLoginKey := strconv.FormatUint(userData.UserLoginKey, 10)
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"
	paramsUserMessage["umessage_recipient_key"] = customerUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"

	subject := "Switching sedang Diproses"
	body := "Switching kamu telah kami terima. Kami sedang memproses transaksi kamu."

	paramsUserMessage["umessage_subject"] = subject
	paramsUserMessage["umessage_body"] = body

	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strIDUserLogin

	status, err = models.CreateScUserMessage(paramsUserMessage)
	if err != nil {
		log.Error("Error create user message")
	} else {
		log.Error("Sukses insert user message")
	}

	//create push notif
	lib.CreateNotifCustomerFromAdminByCustomerId(customerKeyStr, subject, body, "TRANSACTION")

	//send email
	params["product_name"] = product.ProductNameAlt
	params["currency"] = strconv.FormatUint(*product.CurrencyKey, 10)
	params["parrent"] = transactionID
	params["ProductOutNameAlt"] = productTo.ProductNameAlt
	err = mailTransactionManual("switching", params, customerKeyStr, userData.UloginEmail)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func mailTransactionManual(typ string, params map[string]string, customerKey string, email string) error {
	var err error
	var mailTemp, subject string
	decimal.MarshalJSONWithoutQuotes = true
	ac0 := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	mailParam := make(map[string]string)
	if params["currency"] == "1" {
		mailParam["Symbol"] = "Rp. "
	} else if params["currency"] == "2" {
		mailParam["Symbol"] = "$"
	}
	val, _ := decimal.NewFromString(params["trans_fee_amount"])
	mailParam["Fee"] = ac0.FormatMoneyDecimal(val.Truncate(0))
	if typ == "subscription" || typ == "topup" {
		mailTemp = "index-" + typ + "-complete.html"
		s := "Subscription"
		if typ == "topup" {
			s = "Top Up"
		}
		subject = s + " Kamu sedang Diproses"
	} else if typ == "redemption" {
		mailTemp = "index-" + typ + ".html"
		subject = "Redemption Kamu sedang Diproses"

		mailParam["BankName"] = params["BankName"]
		mailParam["BankAccNo"] = params["AccountNo"]
		mailParam["AccHolderName"] = params["AccountHolderName"]
		mailParam["Branch"] = params["BranchName"]

	} else if typ == "switching" {
		mailTemp = "index-" + typ + ".html"
		subject = "Switching Kamu sedang Diproses"
		mailParam["ProductOut"] = params["ProductOutNameAlt"]
	} else {
		log.Error("Failed send mail: type not valid")
		return err
	}

	value, _ := strconv.ParseFloat(params["trans_unit"], 64)
	if value > 0 {
		mailParam["Symbol"] = ""
		mailParam["Amount"] = fmt.Sprintf("%.2f", value)
		mailParam["Str"] = " Unit"
	} else {
		val, _ := decimal.NewFromString(params["trans_amount"])
		mailParam["Amount"] = ac0.FormatMoneyDecimal(val)
	}
	var customer models.MsCustomer
	_, err = models.GetMsCustomer(&customer, customerKey)
	if err != nil {
		log.Error("Failed send mail: " + err.Error())
		return err
	}
	mailParam["Name"] = customer.FullName
	mailParam["Cif"] = customer.UnitHolderIDno
	layout := "2006-01-02 15:04:05"
	dateLayout := "02 Jan 2006"
	date, _ := time.Parse(layout, params["trans_date"])
	mailParam["Date"] = date.Format(dateLayout)
	hr, min, _ := date.Clock()
	mailParam["Time"] = strconv.Itoa(hr) + "." + strconv.Itoa(min) + " WIB"

	mailParam["ProductName"] = params["product_name"]
	mailParam["ProductIn"] = params["product_name"]

	mailParam["FileUrl"] = config.FileUrl + "/images/mail"

	t := template.New(mailTemp)

	t, err = t.ParseFiles(config.BasePath + "/mail/" + mailTemp)
	if err != nil {
		log.Error("Failed send mail: " + err.Error())
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mailParam); err != nil {
		log.Error("Failed send mail: " + err.Error())
		return err
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "[MNC Duit] "+subject)
	mailer.SetBody("text/html", result)

	dialer := gomail.NewDialer(
		config.EmailSMTPHost,
		int(config.EmailSMTPPort),
		config.EmailFrom,
		config.EmailFromPassword,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Error("Failed send mail: " + err.Error())
		return err
	}
	log.Info("Email sent")
	return nil
}
