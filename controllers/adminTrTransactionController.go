package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	return getListAdmin(transStatusKey, c)
}

func GetTransactionCutOffList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "5")

	return getListAdmin(transStatusKey, c)
}

func GetTransactionCorrectionList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "6")

	return getListAdmin(transStatusKey, c)
}

func GetTransactionConfirmationList(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "7")

	return getListAdmin(transStatusKey, c)
}

func GetTransactionCorrectionAdminList(c echo.Context) error {
	errorAuth := initAuthBranchEntryHoEntry()
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var transStatusKey []string
	transStatusKey = append(transStatusKey, "1")

	return getListAdmin(transStatusKey, c)
}

func getListAdmin(transStatusKey []string, c echo.Context) error {

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
	status, err = models.AdminGetAllTrTransaction(&trTransaction, limit, offset, noLimit, params, transStatusKey, "trans_status_key")
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

	strTransStatusKey := strconv.FormatUint(transaction.TransStatusKey, 10)

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
		status := []string{"5", "6", "7"}
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

	if transStatusKeyDefault != strTransStatusKey {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}

	var roleKeyCs uint64
	roleKeyCs = 11
	var roleKeyKyc uint64
	roleKeyKyc = 12

	dateLayout := "2006-01-02 15:04:05"
	strIdUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	if lib.Profile.RoleKey == roleKeyCs {
		params["check1_notes"] = notes
		params["check1_date"] = time.Now().Format(dateLayout)
		params["check1_flag"] = "1"
		params["check1_references"] = strIdUserLogin
	}

	if lib.Profile.RoleKey == roleKeyKyc {
		params["check2_notes"] = notes
		params["check2_date"] = time.Now().Format(dateLayout)
		params["check2_flag"] = "1"
		params["check2_references"] = strIdUserLogin
	}

	params["rec_modified_by"] = strIdUserLogin
	params["rec_modified_date"] = time.Now().Format(dateLayout)

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
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

	_, err = models.UpdateTrTransaction(params)
	if err != nil {
		log.Error("Error update tr transaction")
		return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
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
