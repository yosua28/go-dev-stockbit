package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func initAuthBranchEntryHoEntry() error {
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	var roleKeyHoEntry uint64
	roleKeyHoEntry = 10

	if (lib.Profile.RoleKey != roleKeyBranchEntry) && (lib.Profile.RoleKey != roleKeyHoEntry) {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	return nil
}

func initAuthTransactionAdmin() error {
	roles := []string{"7", "10", "11", "12", "13"}
	strRoleLogin := strconv.FormatUint(lib.Profile.RoleKey, 10)
	_, found := lib.Find(roles, strRoleLogin)
	if !found {
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	return nil
}

func GetTransactionApprovalList(c echo.Context) error {
	errorAuth := initAuthCsKyc()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	var transStatusKey []string

	//if user approval CS
	if lib.Profile.RoleKey == roleKeyCs {
		transStatusKey = append(transStatusKey, "2")
	}
	//if user approval KYC / Complainer
	if lib.Profile.RoleKey == roleKeyKyc {
		transStatusKey = append(transStatusKey, "4")
	}

	return getListAdmin(transStatusKey, c, nil)
}

func GetTransactionCutOffList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "5")

	return getListAdmin(transStatusKey, c, nil)
}

func GetTransactionBatchList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	//date
	postnavdate := c.QueryParam("nav_date")
	if postnavdate == "" {
		log.Error("Missing required parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "6")

	return getListAdmin(transStatusKey, c, &postnavdate)
}

func GetTransactionConfirmationList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "7")

	return getListAdmin(transStatusKey, c, nil)
}

func GetTransactionCorrectionAdminList(c echo.Context) error {
	errorAuth := initAuthBranchEntryHoEntry()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "1")

	return getListAdmin(transStatusKey, c, nil)
}

func GetTransactionPostingList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "8")

	return getListAdmin(transStatusKey, c, nil)
}

func getListAdmin(transStatusKey []string, c echo.Context, postnavdate *string) error {

	var err error
	var status int

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

	params["rec_status"] = "1"
	if postnavdate != nil {
		params["nav_date"] = *postnavdate
	}

	//if user admin role 7 branch
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
	status, err = models.AdminGetAllTrTransaction(&trTransaction, limit, offset, noLimit, params, transStatusKey, "trans_status_key", false)
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
	var bankIds []string
	var parentTransIds []string
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
		if _, ok := lib.Find(customerIds, strconv.FormatUint(tr.CustomerKey, 10)); !ok {
			customerIds = append(customerIds, strconv.FormatUint(tr.CustomerKey, 10))
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(tr.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(tr.ProductKey, 10))
		}
		if _, ok := lib.Find(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10)); !ok {
			transTypeIds = append(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10))
		}
		if tr.TransBankKey != nil {
			if _, ok := lib.Find(bankIds, strconv.FormatUint(*tr.TransBankKey, 10)); !ok {
				bankIds = append(bankIds, strconv.FormatUint(*tr.TransBankKey, 10))
			}
		}

		strTransTypeKey := strconv.FormatUint(tr.TransTypeKey, 10)
		if strTransTypeKey == "4" {
			if tr.ParentKey != nil {
				if _, ok := lib.Find(parentTransIds, strconv.FormatUint(*tr.ParentKey, 10)); !ok {
					parentTransIds = append(parentTransIds, strconv.FormatUint(*tr.ParentKey, 10))
				}
			}
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

	//mapping parent transaction
	var parentTrans []models.TrTransaction
	if len(parentTransIds) > 0 {
		status, err = models.GetTrTransactionIn(&parentTrans, parentTransIds, "transaction_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	parentTransData := make(map[uint64]models.TrTransaction)
	for _, pt := range parentTrans {
		parentTransData[pt.TransactionKey] = pt

		if _, ok := lib.Find(productIds, strconv.FormatUint(pt.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(pt.ProductKey, 10))
		}
	}

	//mapping product
	var msProduct []models.MsProduct
	if len(productIds) > 0 {
		status, err = models.GetMsProductIn(&msProduct, productIds, "product_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
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
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	transactionTypeData := make(map[uint64]models.TrTransactionType)
	for _, t := range transactionType {
		transactionTypeData[t.TransTypeKey] = t
	}

	//mapping ms bank
	var msBank []models.MsBank
	if len(bankIds) > 0 {
		status, err = models.GetMsBankIn(&msBank, bankIds, "bank_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	bankData := make(map[uint64]models.MsBank)
	for _, b := range msBank {
		bankData[b.BankKey] = b
	}

	//mapping trans status
	var trTransactionStatus []models.TrTransactionStatus
	if len(transStatusKey) > 0 {
		status, err = models.GetMsTransactionStatusIn(&trTransactionStatus, transStatusKey, "trans_status_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	transStatusData := make(map[uint64]models.TrTransactionStatus)
	for _, ts := range trTransactionStatus {
		transStatusData[ts.TransStatusKey] = ts
	}

	var responseData []models.AdminTrTransactionList
	for _, tr := range trTransaction {
		var data models.AdminTrTransactionList

		data.TransactionKey = tr.TransactionKey

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
			data.ProductName = n.ProductName
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

		data.TransAmount = tr.TransAmount
		data.TransUnit = tr.TransUnit
		data.TotalAmount = tr.TotalAmount

		if tr.TransBankKey != nil {
			if n, ok := bankData[*tr.TransBankKey]; ok {
				data.TransBankName = n.BankName
			}
		}

		data.TransBankAccNo = tr.TransBankAccNo
		data.TransBankaccName = tr.TransBankaccName

		strTransTypeKey := strconv.FormatUint(tr.TransTypeKey, 10)

		if strTransTypeKey == "4" {
			data.ProductIn = &data.ProductName
			if tr.ParentKey != nil {
				if n, ok := parentTransData[*tr.ParentKey]; ok {
					if pd, ok := productData[n.ProductKey]; ok {
						data.ProductOut = &pd.ProductName
					}
				}
			}
		}

		responseData = append(responseData, data)
	}

	var countData models.CountData
	var pagination int
	if limit > 0 {
		status, err = models.AdminGetCountTrTransaction(&countData, params, transStatusKey, "trans_status_key")
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

func GetTransactionDetail(c echo.Context) error {
	var err error
	var status int

	keyStr := c.Param("key")
	key, _ := strconv.ParseUint(keyStr, 10, 64)
	if key == 0 {
		return lib.CustomError(http.StatusNotFound)
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, keyStr)
	if err != nil {
		return lib.CustomError(status)
	}

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)
	strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

	if strTransTypeKey == "3" {
		log.Error("Data not found")
		return lib.CustomError(http.StatusUnauthorized, "Data not found", "Data not found")
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12
	var roleKeyFundAdmin uint64
	roleKeyFundAdmin = 13
	var roleKeyBranchEntry uint64
	roleKeyBranchEntry = 7
	var roleKeyHoEntry uint64
	roleKeyHoEntry = 10

	if lib.Profile.RoleKey == roleKeyCs {
		statusCs := strconv.FormatUint(uint64(2), 10)
		if statusCs != strTransStatusKey {
			log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}
	if lib.Profile.RoleKey == roleKeyKyc {
		statusKyc := strconv.FormatUint(uint64(4), 10)
		if statusKyc != strTransStatusKey {
			log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}
	if lib.Profile.RoleKey == roleKeyFundAdmin {
		status := []string{"5", "6", "7", "8"}
		_, found := lib.Find(status, strTransStatusKey)
		if !found {
			log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}
	}
	if (lib.Profile.RoleKey == roleKeyBranchEntry) || (lib.Profile.RoleKey == roleKeyHoEntry) {
		statusCorrected := strconv.FormatUint(uint64(1), 10)
		if statusCorrected != strTransStatusKey {
			log.Error("User Autorizer")
			return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
		}

		//if user role 7, check branch
		if lib.Profile.RoleKey == roleKeyBranchEntry {
			if lib.Profile.BranchKey != nil {
				strUserBranchKey := strconv.FormatUint(*lib.Profile.BranchKey, 10)
				strTransBranchKey := strconv.FormatUint(*transaction.BranchKey, 10)
				if strUserBranchKey != strTransBranchKey {
					log.Error("User Autorizer")
					return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
				}
			} else {
				log.Error("User Autorizer")
				return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
			}
		}
	}

	var responseData models.AdminTransactionDetail

	var lookupIds []string

	if transaction.TrxCode != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.TrxCode, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.TrxCode, 10))
		}
	}
	if transaction.EntryMode != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.EntryMode, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.EntryMode, 10))
		}
	}
	if transaction.PaymentMethod != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.PaymentMethod, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.PaymentMethod, 10))
		}
	}
	if transaction.TrxRiskLevel != nil {
		if _, ok := lib.Find(lookupIds, strconv.FormatUint(*transaction.TrxRiskLevel, 10)); !ok {
			lookupIds = append(lookupIds, strconv.FormatUint(*transaction.TrxRiskLevel, 10))
		}
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

	if transaction.TrxCode != nil {
		if n, ok := gData[*transaction.TrxCode]; ok {
			var trc models.LookupTrans

			trc.LookupKey = n.LookupKey
			trc.LkpGroupKey = n.LkpGroupKey
			trc.LkpCode = n.LkpCode
			trc.LkpName = n.LkpName
			responseData.TrxCode = &trc
		}
	}

	if transaction.EntryMode != nil {
		if n, ok := gData[*transaction.EntryMode]; ok {
			var entm models.LookupTrans

			entm.LookupKey = n.LookupKey
			entm.LkpGroupKey = n.LkpGroupKey
			entm.LkpCode = n.LkpCode
			entm.LkpName = n.LkpName
			responseData.EntryMode = &entm
		}
	}

	if transaction.PaymentMethod != nil {
		if n, ok := gData[*transaction.PaymentMethod]; ok {
			var pm models.LookupTrans
			pm.LookupKey = n.LookupKey
			pm.LkpGroupKey = n.LkpGroupKey
			pm.LkpCode = n.LkpCode
			pm.LkpName = n.LkpName
			responseData.PaymentMethod = &pm
		}
	}

	if transaction.TrxRiskLevel != nil {
		if n, ok := gData[*transaction.TrxRiskLevel]; ok {
			var risk models.LookupTrans

			risk.LookupKey = n.LookupKey
			risk.LkpGroupKey = n.LkpGroupKey
			risk.LkpCode = n.LkpCode
			risk.LkpName = n.LkpName
			responseData.TrxRiskLevel = &risk
		}
	}

	layout := "2006-01-02 15:04:05"
	newLayout := "02 Jan 2006"

	responseData.TransactionKey = transaction.TransactionKey
	date, _ := time.Parse(layout, transaction.TransDate)
	responseData.TransDate = date.Format(newLayout)
	date, _ = time.Parse(layout, transaction.NavDate)
	responseData.NavDate = date.Format(newLayout)
	if transaction.RecCreatedDate != nil {
		date, err = time.Parse(layout, *transaction.RecCreatedDate)
		if err == nil {
			oke := date.Format(newLayout)
			responseData.RecCreatedDate = &oke
		}
	}
	responseData.RecCreatedBy = transaction.RecCreatedBy
	responseData.TransAmount = transaction.TransAmount
	responseData.TransUnit = transaction.TransUnit
	responseData.TransUnitPercent = transaction.TransUnitPercent
	if transaction.FlagRedemtAll != nil {
		if int(*transaction.FlagRedemtAll) > 0 {
			responseData.FlagRedemtAll = true
		}
	}
	if transaction.FlagNewSub != nil {
		if int(*transaction.FlagNewSub) > 0 {
			responseData.FlagNewSub = true
		}
	}
	responseData.TransFeePercent = transaction.TransFeePercent
	responseData.TransFeeAmount = transaction.TransFeeAmount
	responseData.ChargesFeeAmount = transaction.ChargesFeeAmount
	responseData.ServicesFeeAmount = transaction.ServicesFeeAmount
	responseData.TotalAmount = transaction.TotalAmount
	responseData.SettlementDate = transaction.SettlementDate
	responseData.TransBankAccNo = transaction.TransBankAccNo
	responseData.TransBankaccName = transaction.TransBankaccName
	responseData.TransRemarks = transaction.TransRemarks
	responseData.TransReferences = transaction.TransReferences
	responseData.PromoCode = transaction.PromoCode
	responseData.SalesCode = transaction.SalesCode
	if transaction.RiskWaiver > 0 {
		responseData.RiskWaiver = true
	}
	responseData.FileUploadDate = transaction.FileUploadDate
	responseData.ProceedDate = transaction.ProceedDate
	responseData.ProceedAmount = transaction.ProceedAmount
	responseData.SentDate = transaction.SentDate
	responseData.SentReferences = transaction.SentReferences
	responseData.ConfirmedDate = transaction.ConfirmedDate
	responseData.PostedDate = transaction.PostedDate
	responseData.PostedUnits = transaction.PostedUnits
	responseData.SettledDate = transaction.SettledDate

	var userData models.ScUserLogin
	strCustomer := strconv.FormatUint(transaction.CustomerKey, 10)
	status, err = models.GetScUserLoginByCustomerKey(&userData, strCustomer)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	}

	dir := config.BaseUrl + "/images/user/" + strconv.FormatUint(userData.UserLoginKey, 10) + "/transfer/"

	if transaction.RecImage1 != nil {
		path := dir + *transaction.RecImage1
		responseData.RecImage1 = &path
	}

	if transaction.BranchKey != nil {
		var branch models.MsBranch
		strBranch := strconv.FormatUint(*transaction.BranchKey, 10)
		status, err = models.GetMsBranch(&branch, strBranch)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var br models.BranchTrans
			br.BranchKey = branch.BranchKey
			br.BranchCode = branch.BranchCode
			br.BranchName = branch.BranchName
			responseData.Branch = &br
		}
	}

	//check agent
	if transaction.AgentKey != nil {
		var agent models.MsAgent
		strAgent := strconv.FormatUint(*transaction.AgentKey, 10)
		status, err = models.GetMsAgent(&agent, strAgent)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var a models.AgentTrans
			a.AgentKey = agent.AgentKey
			a.AgentCode = agent.AgentCode
			a.AgentName = agent.AgentName
			responseData.Agent = &a
		}
	}

	//check customer
	var customer models.MsCustomer
	strCus := strconv.FormatUint(transaction.CustomerKey, 10)
	status, err = models.GetMsCustomer(&customer, strCus)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.Customer.CustomerKey = customer.CustomerKey
		responseData.Customer.FullName = customer.FullName
		responseData.Customer.SidNo = customer.SidNo
		responseData.Customer.UnitHolderIDno = customer.UnitHolderIDno
	}

	//check product
	var product models.MsProduct
	strPro := strconv.FormatUint(transaction.ProductKey, 10)
	status, err = models.GetMsProduct(&product, strPro)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.Product.ProductKey = product.ProductKey
		responseData.Product.ProductCode = product.ProductCode
		responseData.Product.ProductName = product.ProductName
	}

	//check trans status
	var transStatus models.TrTransactionStatus
	strTrSt := strconv.FormatUint(transaction.TransStatusKey, 10)
	status, err = models.GetTrTransactionStatus(&transStatus, strTrSt)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.TransStatus.TransStatusKey = transStatus.TransStatusKey
		responseData.TransStatus.StatusCode = transStatus.StatusCode
		responseData.TransStatus.StatusDescription = transStatus.StatusDescription
	}

	//check trans type
	var transType models.TrTransactionType
	strTrTy := strconv.FormatUint(transaction.TransTypeKey, 10)
	status, err = models.GetMsTransactionType(&transType, strTrTy)
	if err != nil {
		if err != sql.ErrNoRows {
			return lib.CustomError(status)
		}
	} else {
		responseData.TransType.TransTypeKey = transType.TransTypeKey
		responseData.TransType.TypeCode = transType.TypeCode
		responseData.TransType.TypeDescription = transType.TypeDescription
	}

	//check bank
	var bank models.MsBank
	if transaction.TransBankKey != nil {
		strBank := strconv.FormatUint(*transaction.TransBankKey, 10)
		status, err = models.GetMsBank(&bank, strBank)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var tb models.TransBank
			tb.BankKey = bank.BankKey
			tb.BankCode = bank.BankCode
			tb.BankName = bank.BankName
			responseData.TransBank = &tb
		}
	}

	//check aca
	if transaction.AcaKey != nil {
		var aca models.TrAccountAgent
		strAca := strconv.FormatUint(*transaction.AcaKey, 10)
		status, err = models.GetTrAccountAgent(&aca, strAca)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var ac models.AcaTrans
			ac.AcaKey = aca.AcaKey
			ac.AccKey = aca.AccKey
			var agent models.MsAgent
			strAgent := strconv.FormatUint(aca.AgentKey, 10)
			status, err = models.GetMsAgent(&agent, strAgent)
			if err != nil {
				if err != sql.ErrNoRows {
					return lib.CustomError(status)
				}
			} else {
				ac.AgentKey = agent.AgentKey
				ac.AgentCode = agent.AgentCode
				ac.AgentName = agent.AgentName
			}

			responseData.Aca = &ac
		}
	}

	//check transaction confirmation
	if transaction.AcaKey != nil {
		var tc models.TrTransactionConfirmation
		strTrKey := strconv.FormatUint(transaction.TransactionKey, 10)
		status, err = models.GetTrTransactionConfirmation(&tc, strTrKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			var transTc models.TransactionConfirmation
			transTc.TcKey = tc.TcKey
			transTc.ConfirmDate = tc.ConfirmDate
			transTc.ConfirmedAmount = tc.ConfirmedAmount
			transTc.ConfirmedUnit = tc.ConfirmedUnit

			responseData.TransactionConfirmation = &transTc
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func TransactionApprovalCs(c echo.Context) error {
	errorAuth := initAuthCs()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	transStatusKeyDefault := "2"
	transStatusIds := []string{"1", "3", "4"}

	return ProsesApproval(transStatusKeyDefault, transStatusIds, c)
}

func TransactionApprovalCompliance(c echo.Context) error {
	errorAuth := initAuthKyc()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	transStatusKeyDefault := "4"
	transStatusIds := []string{"1", "3", "5"}

	return ProsesApproval(transStatusKeyDefault, transStatusIds, c)
}

func ProsesApproval(transStatusKeyDefault string, transStatusIds []string, c echo.Context) error {
	var err error
	var status int

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	params := make(map[string]string)

	transStatus := c.FormValue("trans_status_key")
	if transStatus == "" {
		log.Error("Missing required parameter: trans_status_key")
		return lib.CustomError(http.StatusBadRequest)
	} else {
		_, found := lib.Find(transStatusIds, transStatus)
		if !found {
			log.Error("Missing required parameter: trans_status_key")
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	n, err := strconv.ParseUint(transStatus, 10, 64)
	if err == nil && n > 0 {
		params["trans_status_key"] = transStatus
	} else {
		log.Error("Wrong input for parameter: trans_status_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_status_key", "Wrong input for parameter: trans_status_key")
	}

	notes := c.FormValue("notes")

	if lib.Profile.RoleKey == roleKeyKyc {
		trxrisklevel := c.FormValue("trx_risk_level")
		if trxrisklevel == "" {
			log.Error("Missing required parameter: trx_risk_level")
			return lib.CustomError(http.StatusBadRequest, "trx_risk_level can not be blank", "trx_risk_level can not be blank")
		} else {

			listLevelOption := []string{"114", "115"} //lookup group key 24
			_, found := lib.Find(listLevelOption, trxrisklevel)
			if !found {
				log.Error("Missing required parameter: trx_risk_level")
				return lib.CustomError(http.StatusBadRequest)
			}
		}

		n, err := strconv.ParseUint(trxrisklevel, 10, 64)
		if err == nil && n > 0 {
			params["trx_risk_level"] = trxrisklevel
		} else {
			log.Error("Wrong input for parameter: trx_risk_level")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trx_risk_level", "Wrong input for parameter: trx_risk_level")
		}
	}

	transactionkey := c.FormValue("transaction_key")
	if transactionkey == "" {
		log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err = strconv.ParseUint(transactionkey, 10, 64)
	if err == nil && n > 0 {
		params["transaction_key"] = transactionkey
	} else {
		log.Error("Wrong input for parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, transactionkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

	strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

	if strTransTypeKey == "3" {
		log.Error("Transaction not found")
		return lib.CustomError(http.StatusBadRequest)
	}

	if transStatusKeyDefault != strTransStatusKey {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	if lib.Profile.RoleKey == roleKeyCs {
		params["check1_notes"] = notes
		params["check1_date"] = time.Now().Format(dateLayout)
		params["check1_flag"] = "1"
		params["check1_references"] = strIDUserLogin
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		params["check2_notes"] = notes
		params["check2_date"] = time.Now().Format(dateLayout)
		params["check2_flag"] = "1"
		params["check2_references"] = strIDUserLogin
	}

	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	if strTransTypeKey == "4" {
		if transaction.ParentKey != nil {
			strParentKey := strconv.FormatUint(*transaction.ParentKey, 10)
			params["transaction_key"] = strParentKey
			log.Println(params)
			_, err = models.UpdateTrTransaction(params)
			if err != nil {
				log.Error("Error update tr transaction parent")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
			}
		}
	}

	log.Info("Success update transaksi")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func UpdateNavDate(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var err error
	var status int

	params := make(map[string]string)

	//date
	postnavdate := c.FormValue("nav_date")
	if postnavdate == "" {
		log.Error("Missing required parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}

	layoutISO := "2006-01-02"

	t, _ := time.Parse(layoutISO, postnavdate)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	w := lib.IsWeekend(t)
	if w {
		log.Error("Missing required parameter: nav_date cann't Weekend")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date cann't Weekend", "Missing required parameter: nav_date cann't Weekend")
	}

	paramHoliday := make(map[string]string)
	paramHoliday["holiday_date"] = postnavdate

	var holiday []models.MsHoliday
	status, err = models.GetAllMsHoliday(&holiday, paramHoliday)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	if len(holiday) > 0 {
		log.Error("nav_date is Bursa Holiday")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date is Bursa Holiday", "Missing required parameter: nav_date is Bursa Holiday")
	}

	params["nav_date"] = postnavdate

	transactionkey := c.FormValue("transaction_key")
	if transactionkey == "" {
		log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err := strconv.ParseUint(transactionkey, 10, 64)
	if err == nil && n > 0 {
		params["transaction_key"] = transactionkey
	} else {
		log.Error("Wrong input for parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, transactionkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strTransType := strconv.FormatUint(transaction.TransTypeKey, 10)

	if strTransType == "3" {
		log.Error("Transaction not found")
		return lib.CustomError(http.StatusBadRequest)
	}

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

	strStatusCutOff := "5"

	if strTransStatusKey != strStatusCutOff {
		log.Error("Data not found")
		return lib.CustomError(http.StatusUnauthorized, "Data not found", "Data not found")
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	params["settlement_date"] = postnavdate

	//set settlement_date by settlement period product
	if strTransType == "2" { //REDM
		//check product
		var product models.MsProduct
		strPro := strconv.FormatUint(transaction.ProductKey, 10)
		status, err = models.GetMsProduct(&product, strPro)
		if err != nil {
			if err != sql.ErrNoRows {
				return lib.CustomError(status)
			}
		} else {
			layoutISO := "2006-01-02"
			t, _ := time.Parse(layoutISO, postnavdate)

			if (product.SettlementPeriod != nil) && (*product.SettlementPeriod > 0) {
				params["settlement_date"] = SettDate(t, int(*product.SettlementPeriod))
			}

		}
	}

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	if strTransType == "4" {
		if transaction.ParentKey != nil {
			strParentKey := strconv.FormatUint(*transaction.ParentKey, 10)
			params["transaction_key"] = strParentKey
			log.Println(params)
			_, err = models.UpdateTrTransaction(params)
			if err != nil {
				log.Error("Error update tr transaction parent")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
			}
		}
	}

	log.Info("Success update transaksi")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func TransactionApprovalCutOff(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	//list id
	transIds := c.FormValue("trans_ids")
	if transIds == "" {
		log.Error("Missing required parameter: trans_ids")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_ids", "Missing required parameter: trans_ids")
	}

	s := strings.Split(transIds, ",")

	var transParamIds []string

	for _, value := range s {
		is := strings.TrimSpace(value)
		if is != "" {
			if _, ok := lib.Find(transParamIds, is); !ok {
				transParamIds = append(transParamIds, is)
			}
		}
	}

	var transactionList []models.TrTransaction
	if len(transParamIds) > 0 {
		status, err := models.GetTrTransactionIn(&transactionList, transParamIds, "transaction_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
		if len(transParamIds) != len(transactionList) {
			log.Error("Missing required parameter: trans_ids")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: Jumlah Data & Parameter berbeda", "Missing required parameter: Jumlah Data & Parameter berbeda")
		}
	} else {
		log.Error("Missing required parameter: trans_ids")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_ids", "Missing required parameter: trans_ids")
	}

	strStatusCutOff := "5"

	for _, tr := range transactionList {
		strTransStatusKey := strconv.FormatUint(tr.TransStatusKey, 10)
		if strTransStatusKey != strStatusCutOff {
			log.Error("Missing required parameter: trans_ids")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_ids ", "Missing required parameter: trans_ids")
		}

		strTransType := strconv.FormatUint(tr.TransTypeKey, 10)
		if strTransType == "4" {
			if tr.ParentKey != nil {
				if _, ok := lib.Find(transParamIds, strconv.FormatUint(*tr.ParentKey, 10)); !ok {
					transParamIds = append(transParamIds, strconv.FormatUint(*tr.ParentKey, 10))
				}
			}
		}
	}

	paramsUpdate := make(map[string]string)

	paramsUpdate["trans_status_key"] = "6"
	dateLayout := "2006-01-02 15:04:05"
	paramsUpdate["rec_modified_date"] = time.Now().Format(dateLayout)
	strKey := strconv.FormatUint(lib.Profile.UserID, 10)
	paramsUpdate["rec_modified_by"] = strKey

	_, err := models.UpdateTrTransactionByKeyIn(paramsUpdate, transParamIds, "transaction_key")
	if err != nil {
		log.Error("Error update oa request")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func SettDate(t time.Time, period int) string {
	var dateNext string

	layoutISO := "2006-01-02"
	dateNext = t.Format(layoutISO)

	for i := 0; i < period; i++ {

		t, _ := time.Parse(layoutISO, dateNext)
		dateAfter := t.AddDate(0, 0, 1)
		t = time.Date(dateAfter.Year(), dateAfter.Month(), dateAfter.Day(), 0, 0, 0, 0, time.UTC)
		dateAfter = SkipWeekend(dateAfter)

		strDate := dateAfter.Format(layoutISO)
		dateNext = CheckHolidayBursa(strDate)
	}
	return dateNext
}

func SkipWeekend(t time.Time) time.Time {
	dateAfter := t
	t = t.UTC()

	switch t.Weekday() {
	case time.Saturday:
		dateAfter = t.AddDate(0, 0, 2)
	case time.Sunday:
		dateAfter = t.AddDate(0, 0, 1)
	}

	return dateAfter
}

func CheckHolidayBursa(date string) string {
	dateStr := date
	layoutISO := "2006-01-02"
	paramHoliday := make(map[string]string)
	paramHoliday["holiday_date"] = date

	var holiday []models.MsHoliday
	_, err := models.GetAllMsHoliday(&holiday, paramHoliday)
	if err != nil {
		if err == sql.ErrNoRows {
			t, _ := time.Parse(layoutISO, date)
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			dateStr = date
		}
	}

	if len(holiday) > 0 {
		t, _ := time.Parse(layoutISO, dateStr)
		dateAfter := t.AddDate(0, 0, 1)
		strDate := dateAfter.Format(layoutISO)
		dateStr = CheckHolidayBursa(strDate)
	}

	w, _ := time.Parse(layoutISO, dateStr)
	w = time.Date(w.Year(), w.Month(), w.Day(), 0, 0, 0, 0, time.UTC)
	cek := lib.IsWeekend(w)
	if cek {
		dateSkip := SkipWeekend(w)

		dateStr = dateSkip.Format(layoutISO)
		dateStr = CheckHolidayBursa(dateStr)
	}
	return dateStr
}

func GetFormatExcelDownloadList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "7")

	var err error
	var status int

	params := make(map[string]string)

	//date
	postnavdate := c.QueryParam("nav_date")
	if postnavdate == "" {
		log.Error("Missing required parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}

	transactiontype := c.QueryParam("transaction_type")
	if transactiontype == "" {
		log.Error("Missing required parameter: transaction_type")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: transaction_type", "Missing required parameter: transaction_type")
	}

	rolestransactiontype := []string{"1", "2", "3", "4"}
	_, found := lib.Find(rolestransactiontype, transactiontype)
	if !found {
		return lib.CustomError(http.StatusUnauthorized, "Missing parameter: transaction_type", "Missing parameter: transaction_type")
	}

	params["rec_status"] = "1"
	params["nav_date"] = postnavdate
	params["trans_type_key"] = transactiontype

	var trTransaction []models.TrTransaction
	status, err = models.AdminGetAllTrTransaction(&trTransaction, 0, 0, true, params, transStatusKey, "trans_status_key", true)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get data")
	}
	if len(trTransaction) < 1 {
		log.Error("transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var transTypeIds []string
	var customerIds []string
	var productIds []string
	for _, tr := range trTransaction {
		if _, ok := lib.Find(customerIds, strconv.FormatUint(tr.CustomerKey, 10)); !ok {
			customerIds = append(customerIds, strconv.FormatUint(tr.CustomerKey, 10))
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(tr.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(tr.ProductKey, 10))
		}
		if _, ok := lib.Find(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10)); !ok {
			transTypeIds = append(transTypeIds, strconv.FormatUint(tr.TransTypeKey, 10))
		}
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
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	var custodianIds []string
	productData := make(map[uint64]models.MsProduct)
	for _, p := range msProduct {
		productData[p.ProductKey] = p

		if p.CustodianKey != nil {
			if _, ok := lib.Find(custodianIds, strconv.FormatUint(*p.CustodianKey, 10)); !ok {
				custodianIds = append(custodianIds, strconv.FormatUint(*p.CustodianKey, 10))
			}
		}
	}

	//mapping Trans type
	var transactionType []models.TrTransactionType
	if len(transTypeIds) > 0 {
		status, err = models.GetMsTransactionTypeIn(&transactionType, transTypeIds, "trans_type_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	transactionTypeData := make(map[uint64]models.TrTransactionType)
	for _, t := range transactionType {
		transactionTypeData[t.TransTypeKey] = t
	}

	//mapping ms custodian bank
	var custodianBank []models.MsCustodianBank
	if len(custodianIds) > 0 {
		status, err = models.GetMsCustodianBankIn(&custodianBank, custodianIds, "custodian_key")
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}
	custodianBankData := make(map[uint64]models.MsCustodianBank)
	for _, cb := range custodianBank {
		custodianBankData[cb.CustodianKey] = cb
	}

	//mapping tr nav
	var trNav []models.TrNav
	if len(productIds) > 0 {
		status, err = models.GetAllTrNavBetween(&trNav, postnavdate, postnavdate, productIds)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	trNavData := make(map[uint64]models.TrNav)
	for _, tr := range trNav {
		trNavData[tr.ProductKey] = tr
	}

	var responseData []models.DownloadFormatExcelList
	for _, tr := range trTransaction {
		var data models.DownloadFormatExcelList

		data.IDTransaction = tr.TransactionKey
		if n, ok := transactionTypeData[tr.TransTypeKey]; ok {
			data.IDCategory = *n.TypeCode
		}
		if n, ok := productData[tr.ProductKey]; ok {
			data.ProductName = n.ProductName

			if n.CustodianKey != nil {
				if cb, ok := custodianBankData[*n.CustodianKey]; ok {
					data.Keterangan = cb.CustodianCode
				}
			}
		}

		if n, ok := customerData[tr.CustomerKey]; ok {
			data.FullName = n.FullName
		}

		layout := "2006-01-02 15:04:05"
		newLayout := "01/02/2006"
		date, _ := time.Parse(layout, tr.NavDate)
		data.NavDate = date.Format(newLayout)
		date, _ = time.Parse(layout, tr.TransDate)
		data.TransactionDate = date.Format(newLayout)

		data.Units = tr.TransUnit
		data.NetAmount = tr.TransAmount

		data.NavValue = nil
		if nv, ok := trNavData[tr.ProductKey]; ok {
			data.NavValue = &nv.NavValue
		} else {
			data.Keterangan = "NAV VALUE NOT EXIST"
		}
		data.ApproveUnits = tr.TransUnit
		data.ApproveAmount = tr.TransAmount
		data.Result = ""

		responseData = append(responseData, data)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}

func UploadExcelConfirmation(c echo.Context) error {
	var err error

	var responseData []models.DownloadFormatExcelList

	err = os.MkdirAll(config.BasePath+"/transaksi_upload/confirmation/", 0755)
	if err != nil {
		log.Error(err.Error())
	} else {
		var file *multipart.FileHeader
		file, err = c.FormFile("excel")

		if file != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			// Get file extension
			extension := filepath.Ext(file.Filename)
			log.Println(extension)
			roles := []string{".xlsx", ".XLSX", ".xls", ".XLS"}
			_, found := lib.Find(roles, extension)
			if !found {
				return lib.CustomError(http.StatusUnauthorized, "Format file must xlsx/xls", "Format file must xlsx/xls")
			}

			// Generate filename
			var filename string
			filename = lib.RandStringBytesMaskImprSrc(20)
			log.Println("Generate filename:", filename)
			// Upload image and move to proper directory
			err = lib.UploadImage(file, config.BasePath+"/transaksi_upload/confirmation/"+filename+extension)
			if err != nil {
				log.Println(err)
				return lib.CustomError(http.StatusInternalServerError)
			}

			xlsx, err := excelize.OpenFile(config.BasePath + "/transaksi_upload/confirmation/" + filename + extension)
			if err != nil {
				log.Fatal("ERROR", err.Error())
			}

			sheet1Name := xlsx.GetSheetName(1)

			for i := 2; i < 1000; i++ {
				var data models.DownloadFormatExcelList

				iDTransaction := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("A%d", i))
				iDCategory := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("B%d", i))
				productName := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("C%d", i))
				fullName := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("D%d", i))
				navDate := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("E%d", i))
				transactionDate := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("F%d", i))
				units := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("G%d", i))
				netAmount := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("H%d", i))
				navValue := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("I%d", i))
				approveUnits := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("J%d", i))
				approveAmount := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("K%d", i))
				keterangan := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("L%d", i))
				result := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("M%d", i))

				log.Println(navDate)
				log.Println(transactionDate)

				if iDTransaction == "" {
					break
				}

				key, _ := strconv.ParseUint(iDTransaction, 10, 64)
				if key == 0 {
					return lib.CustomError(http.StatusNotFound)
				}

				data.IDTransaction = key
				data.IDCategory = iDCategory
				data.ProductName = productName
				data.FullName = fullName
				data.NavDate = navDate
				data.TransactionDate = transactionDate

				if units != "" {
					if unitsFloat, err := strconv.ParseFloat(units, 64); err == nil {
						data.Units = float32(unitsFloat)
					}
				}

				if netAmount != "" {
					if netAmountFloat, err := strconv.ParseFloat(netAmount, 64); err == nil {
						data.NetAmount = float32(netAmountFloat)
					}
				}

				if navValue != "" {
					if navValueFloat, err := strconv.ParseFloat(navValue, 64); err == nil {
						nav := float32(navValueFloat)
						data.NavValue = &nav
					}
				}

				var transUnitFifo float32

				if approveUnits != "" {
					if approveUnitsFloat, err := strconv.ParseFloat(approveUnits, 64); err == nil {
						data.ApproveUnits = float32(approveUnitsFloat)
						transUnitFifo = data.ApproveUnits
					}
				}

				if approveAmount != "" {
					if approveAmountFloat, err := strconv.ParseFloat(approveAmount, 64); err == nil {
						data.ApproveAmount = float32(approveAmountFloat)
					}
				}

				data.Keterangan = keterangan
				data.Result = result

				//cek transaction
				var transaction models.TrTransaction
				_, err := models.GetTrTransaction(&transaction, iDTransaction)
				if err != nil {
					if err == sql.ErrNoRows {
						data.Result = "Data Transaction Not Found"
					} else {
						data.Result = err.Error()
					}
					fmt.Printf("%v \n", data)
					responseData = append(responseData, data)
					continue
				}

				//cek trans status
				strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)
				if strTransStatusKey != "7" {
					data.Result = "Status Transaction Not in CONFIRMED"
					fmt.Printf("%v \n", data)
					responseData = append(responseData, data)
					continue
				}

				//cek transaction confirmation
				var transactionConf models.TrTransactionConfirmation
				_, err = models.GetTrTransactionConfirmationByTransactionKey(&transactionConf, iDTransaction)
				if err != nil {
					if err != sql.ErrNoRows {
						data.Result = err.Error()
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}
				} else {
					data.Result = "TC already exists"
					fmt.Printf("%v \n", data)
					responseData = append(responseData, data)
					continue
				}

				strProductKey := strconv.FormatUint(transaction.ProductKey, 10)
				strCustomerKey := strconv.FormatUint(transaction.CustomerKey, 10)

				var trNav []models.TrNav
				_, err = models.GetTrNavByProductKeyAndNavDate(&trNav, strProductKey, transaction.NavDate)
				if err != nil {
					if err != sql.ErrNoRows {
						data.Result = "NAV VALUE NOT EXIST"
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					} else {
						data.Result = err.Error()
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}
				}
				strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

				var trBalanceCustomer []models.TrBalanceCustomerProduk

				//redm cek balance / saldo aktif
				if (strTransTypeKey == "2") || (strTransTypeKey == "3") { // REDM
					_, err = models.GetLastBalanceCustomerByProductKey(&trBalanceCustomer, strCustomerKey, strProductKey)
					if err != nil {
						if err != sql.ErrNoRows {
							data.Result = "Balance is empty"
							fmt.Printf("%v \n", data)
							responseData = append(responseData, data)
							continue
						} else {
							data.Result = err.Error()
							fmt.Printf("%v \n", data)
							responseData = append(responseData, data)
							continue
						}
					}
				}

				//redm cek balance / saldo aktif di parent jika switch
				var transactionParent models.TrTransaction
				var trParentBalanceCustomer []models.TrBalanceCustomerProduk
				if strTransTypeKey == "4" { // SWITCH IN
					if transaction.ParentKey == nil {
						data.Result = "Parent Transaction is empty"
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}

					strTrParentKey := strconv.FormatUint(*transaction.ParentKey, 10)
					_, err := models.GetTrTransaction(&transactionParent, strTrParentKey)
					if err != nil {
						if err == sql.ErrNoRows {
							data.Result = "Data Parent Transaction Not Found"
						} else {
							data.Result = err.Error()
						}
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}

					//cek parent transaksi apakah sudah di posting
					strParentTransStatusKey := strconv.FormatUint(transactionParent.TransStatusKey, 10)
					if strParentTransStatusKey != "9" {
						data.Result = "Data Parent Transaction has not been posting"
						fmt.Printf("%v \n", data)
						responseData = append(responseData, data)
						continue
					}

					strParentProductKey := strconv.FormatUint(transactionParent.ProductKey, 10)
					strParentCustomerKey := strconv.FormatUint(transactionParent.CustomerKey, 10)
					_, err = models.GetLastBalanceCustomerByProductKey(&trParentBalanceCustomer, strParentCustomerKey, strParentProductKey)
					if err != nil {
						if err != sql.ErrNoRows {
							data.Result = "Parent Balance is empty"
							fmt.Printf("%v \n", data)
							responseData = append(responseData, data)
							continue
						} else {
							data.Result = err.Error()
							fmt.Printf("%v \n", data)
							responseData = append(responseData, data)
							continue
						}
					}
				}

				//data valid 1. create tr_transaction_confirmation, 2. update trans status, 3. create tr_transaction_fifo
				//1. create tr_transaction_confirmation
				dateLayout := "2006-01-02 15:04:05"
				params := make(map[string]string)
				params["confirm_date"] = time.Now().Format(dateLayout)
				params["transaction_key"] = iDTransaction
				params["confirmed_amount"] = approveAmount
				params["confirmed_unit"] = approveUnits
				params["confirm_result"] = "208"

				var approveUnitsFloat float32
				approveUnitsFloat = 0
				if approveUnits != "" {
					if appUnits, err := strconv.ParseFloat(approveUnits, 64); err == nil {
						approveUnitsFloat = float32(appUnits)
					}
				}
				if transaction.TransUnit > approveUnitsFloat {
					strTransUnit := fmt.Sprintf("%g", transaction.TransUnit-approveUnitsFloat)
					params["confirmed_unit_diff"] = strTransUnit
				} else if transaction.TransUnit < approveUnitsFloat {
					strTransUnit := fmt.Sprintf("%g", approveUnitsFloat-transaction.TransUnit)
					params["confirmed_unit_diff"] = strTransUnit
				} else {
					params["confirmed_unit_diff"] = "0"
				}

				var approveAmountFloat float32
				approveAmountFloat = 0
				if approveUnits != "" {
					if appAmount, err := strconv.ParseFloat(approveAmount, 64); err == nil {
						approveAmountFloat = float32(appAmount)
					}
				}
				if transaction.TransAmount > approveAmountFloat {
					strTransAmount := fmt.Sprintf("%g", transaction.TransAmount-approveAmountFloat)
					params["confirmed_amount_diff"] = strTransAmount
				} else if transaction.TransAmount < approveAmountFloat {
					strTransAmount := fmt.Sprintf("%g", approveAmountFloat-transaction.TransAmount)
					params["confirmed_amount_diff"] = strTransAmount
				} else {
					params["confirmed_amount_diff"] = "0"
				}

				params["rec_status"] = "1"
				params["rec_created_date"] = time.Now().Format(dateLayout)
				params["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)

				status, err, trConfirmationID := models.CreateTrTransactionConfirmation(params)
				if err != nil {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed input data")
				}

				// 2. update trans status
				paramsTrans := make(map[string]string)
				paramsTrans["trans_status_key"] = "8"
				paramsTrans["confirmed_date"] = time.Now().Format(dateLayout)
				paramsTrans["rec_modified_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
				paramsTrans["rec_modified_date"] = time.Now().Format(dateLayout)
				paramsTrans["transaction_key"] = iDTransaction
				_, err = models.UpdateTrTransaction(paramsTrans)
				if err != nil {
					log.Error("Error update tr transaction")
					return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
				}

				//3. create tr_transaction_fifo

				if (strTransTypeKey == "1") || (strTransTypeKey == "4") { // SUB / switchin
					paramsFifo := make(map[string]string)
					paramsFifo["trans_conf_sub_key"] = trConfirmationID
					if transaction.AcaKey != nil {
						strAcaKey := strconv.FormatUint(*transaction.AcaKey, 10)
						paramsFifo["sub_aca_key"] = strAcaKey
					}
					paramsFifo["holding_days"] = "0"
					paramsFifo["trans_unit"] = approveUnits
					paramsFifo["fee_nav_mode"] = "207"

					var transAmountFifo float32
					transAmountFifo = transUnitFifo * trNav[0].NavValue
					strTransAmountFifo := fmt.Sprintf("%g", transAmountFifo)
					paramsFifo["trans_amount"] = strTransAmountFifo

					// paramsFifo["trans_fee_amount"] = ""
					paramsFifo["trans_fee_tax"] = "0"
					// paramsFifo["trans_nett_amount"] = ""
					paramsFifo["rec_status"] = "1"
					paramsFifo["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
					paramsFifo["rec_created_date"] = time.Now().Format(dateLayout)
					_, err = models.CreateTrTransactionFifo(paramsFifo)
					if err != nil {
						log.Error("Error update tr transaction")
						return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed insert data")
					}
				}

				if (strTransTypeKey == "2") || (strTransTypeKey == "3") { // REDM / switchout
					sisaFifo := transUnitFifo
					for _, trBalance := range trBalanceCustomer {
						if sisaFifo > 0 {
							var balanceUsed float32

							if trBalance.BalanceUnit > sisaFifo {
								balanceUsed = sisaFifo
								sisaFifo = 0
							}

							if trBalance.BalanceUnit < sisaFifo {
								balanceUsed = trBalance.BalanceUnit
								sisaFifo = sisaFifo - trBalance.BalanceUnit
							}

							if trBalance.BalanceUnit == sisaFifo {
								balanceUsed = trBalance.BalanceUnit
								sisaFifo = 0
							}

							paramsFifo := make(map[string]string)
							paramsFifo["trans_conf_red_key"] = trConfirmationID
							strTcKeySub := strconv.FormatUint(trBalance.TcKey, 10)
							paramsFifo["trans_conf_sub_key"] = strTcKeySub
							if transaction.AcaKey != nil {
								strAcaKey := strconv.FormatUint(*transaction.AcaKey, 10)
								paramsFifo["sub_aca_key"] = strAcaKey
							}

							day1, _ := time.Parse(dateLayout, transaction.NavDate)
							day2, _ := time.Parse(dateLayout, trBalance.NavDate)

							days := day1.Sub(day2).Hours() / 24
							strDays := fmt.Sprintf("%g", days)

							paramsFifo["holding_days"] = strDays

							strUnitUsed := fmt.Sprintf("%g", balanceUsed)
							paramsFifo["trans_unit"] = strUnitUsed
							paramsFifo["fee_nav_mode"] = "207"

							var transAmountFifo float32
							transAmountFifo = transUnitFifo * trNav[0].NavValue
							strTransAmountFifo := fmt.Sprintf("%g", transAmountFifo)
							paramsFifo["trans_amount"] = strTransAmountFifo

							// paramsFifo["trans_fee_amount"] = ""
							paramsFifo["trans_fee_tax"] = "0"
							// paramsFifo["trans_nett_amount"] = ""
							paramsFifo["rec_status"] = "1"
							paramsFifo["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
							paramsFifo["rec_created_date"] = time.Now().Format(dateLayout)
							_, err = models.CreateTrTransactionFifo(paramsFifo)
							if err != nil {
								log.Error("Error update tr transaction")
								return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed insert data")
							}
						} else {
							break
						}
					}
				}

				data.Keterangan = ""
				data.Result = "SUCCESS"

				responseData = append(responseData, data)
			}
		} else {
			log.Error("File cann't be blank")
			return lib.CustomError(http.StatusNotFound, "File can not be blank", "File can not be blank")
		}
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)

}

func ProsesPosting(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var err error
	var status int

	params := make(map[string]string)

	transactionkey := c.FormValue("transaction_key")
	if transactionkey == "" {
		log.Error("Missing required parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest)
	}

	n, err := strconv.ParseUint(transactionkey, 10, 64)
	if err == nil && n > 0 {
		params["transaction_key"] = transactionkey
	} else {
		log.Error("Wrong input for parameter: transaction_key")
		return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: transaction_key", "Wrong input for parameter: transaction_key")
	}

	var transaction models.TrTransaction
	status, err = models.GetTrTransaction(&transaction, transactionkey)
	if err != nil {
		return lib.CustomError(status)
	}

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

	if strTransStatusKey != "8" {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusBadRequest)
	}

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	strTransTypeKey := strconv.FormatUint(transaction.TransTypeKey, 10)

	var transactionConf models.TrTransactionConfirmation
	strTransactionKey := strconv.FormatUint(transaction.TransactionKey, 10)
	_, err = models.GetTrTransactionConfirmationByTransactionKey(&transactionConf, strTransactionKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	var trBalanceCustomer []models.TrBalanceCustomerProduk
	strProductKey := strconv.FormatUint(transaction.ProductKey, 10)
	strCustomerKey := strconv.FormatUint(transaction.CustomerKey, 10)

	if strTransTypeKey == "2" { // REDM
		_, err = models.GetLastBalanceCustomerByProductKey(&trBalanceCustomer, strCustomerKey, strProductKey)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("Transaction have not balance")
				return lib.CustomError(http.StatusBadRequest)
			} else {
				log.Error(err.Error())
				return lib.CustomError(http.StatusBadRequest)
			}
		}
	}

	strTransUnit := fmt.Sprintf("%g", transactionConf.ConfirmedUnit)

	//create tr_balance
	if strTransTypeKey == "1" { // SUB
		paramsBalance := make(map[string]string)
		strAcaKey := strconv.FormatUint(*transaction.AcaKey, 10)
		paramsBalance["aca_key"] = strAcaKey
		strTransactionConf := strconv.FormatUint(transactionConf.TcKey, 10)
		paramsBalance["tc_key"] = strTransactionConf

		newlayout := "2006-01-02"
		t, _ := time.Parse(dateLayout, transactionConf.ConfirmDate)
		balanceDate := t.Format(newlayout)

		paramsBalance["balance_date"] = balanceDate + " 00:00:00"
		paramsBalance["balance_unit"] = strTransUnit
		paramsBalance["rec_order"] = "0"
		paramsBalance["rec_status"] = "1"
		paramsBalance["rec_created_date"] = time.Now().Format(dateLayout)
		paramsBalance["rec_created_by"] = strIDUserLogin
		status, err := models.CreateTrBalance(paramsBalance)
		if err != nil {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed input data")
		}
	}

	if strTransTypeKey == "2" { // REDM
		sisaFifo := transactionConf.ConfirmedUnit
		for _, trBalance := range trBalanceCustomer {
			if sisaFifo > 0 {
				var sisaBalance float32

				if trBalance.BalanceUnit > sisaFifo {
					sisaBalance = trBalance.BalanceUnit - sisaFifo
					sisaFifo = 0
				}

				if trBalance.BalanceUnit < sisaFifo {
					sisaBalance = 0
					sisaFifo = sisaFifo - trBalance.BalanceUnit
				}

				if trBalance.BalanceUnit == sisaFifo {
					sisaBalance = 0
					sisaFifo = 0
				}

				paramsBalance := make(map[string]string)
				strAcaKey := strconv.FormatUint(*transaction.AcaKey, 10)
				paramsBalance["aca_key"] = strAcaKey
				strTransactionSubs := strconv.FormatUint(trBalance.TcKey, 10)
				paramsBalance["tc_key"] = strTransactionSubs
				strTransactionRed := strconv.FormatUint(transactionConf.TcKey, 10)
				paramsBalance["tc_key_red"] = strTransactionRed

				newlayout := "2006-01-02"
				t, _ := time.Parse(dateLayout, transactionConf.ConfirmDate)
				balanceDate := t.Format(newlayout)

				strTransUnitSisa := fmt.Sprintf("%g", sisaBalance)

				paramsBalance["balance_date"] = balanceDate + " 00:00:00"
				paramsBalance["balance_unit"] = strTransUnitSisa

				var balance models.TrBalance
				status, err = models.GetLastTrBalanceByTcRed(&balance, strTransactionRed)
				if err != nil {
					paramsBalance["rec_order"] = "0"
				} else {
					if balance.RecOrder == nil {
						paramsBalance["rec_order"] = "0"
					} else {
						orderNext := int(*balance.RecOrder) + 1
						strOrderNext := strconv.FormatInt(int64(orderNext), 10)
						paramsBalance["rec_order"] = strOrderNext
					}
				}

				paramsBalance["rec_status"] = "1"
				paramsBalance["rec_created_date"] = time.Now().Format(dateLayout)
				paramsBalance["rec_created_by"] = strIDUserLogin
				status, err := models.CreateTrBalance(paramsBalance)
				if err != nil {
					log.Error(err.Error())
					return lib.CustomError(status, err.Error(), "Failed input data")
				}
			} else {
				break
			}
		}
	}

	if strTransTypeKey == "4" { // SWITCH
	}

	//update tr_transaction
	params["posted_units"] = strTransUnit
	params["trans_status_key"] = "9"
	params["posted_date"] = time.Now().Format(dateLayout)
	params["rec_modified_by"] = strIDUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
	}

	//create user message
	paramsUserMessage := make(map[string]string)
	paramsUserMessage["umessage_type"] = "245"

	var userLogin models.ScUserLogin
	_, err = models.GetScUserLoginByCustomerKey(&userLogin, strCustomerKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest)
	}

	strUserLoginKey := strconv.FormatUint(userLogin.UserLoginKey, 10)
	paramsUserMessage["umessage_recipient_key"] = strUserLoginKey
	paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_read"] = "0"
	paramsUserMessage["umessage_sender_key"] = strIDUserLogin
	paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["flag_sent"] = "1"
	if strTransTypeKey == "1" { // SUBS
		paramsUserMessage["umessage_subject"] = "Subscribe Produk"
	}
	if strTransTypeKey == "2" { // REDM
		paramsUserMessage["umessage_subject"] = "Redm Produk"
	}
	if strTransTypeKey == "4" { // SWITCH
		paramsUserMessage["umessage_subject"] = "Switch Produk"
	}
	paramsUserMessage["umessage_body"] = "Selamat !!! Transaksi anda sudah di approv."
	paramsUserMessage["umessage_category"] = "248"
	paramsUserMessage["flag_archieved"] = "0"
	paramsUserMessage["archieved_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_status"] = "1"
	paramsUserMessage["rec_created_date"] = time.Now().Format(dateLayout)
	paramsUserMessage["rec_created_by"] = strIDUserLogin

	log.Info("Success update transaksi")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}
