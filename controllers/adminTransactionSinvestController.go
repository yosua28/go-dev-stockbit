package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/shopspring/decimal"
)

func DownloadTransactionFormatSinvest(c echo.Context) error {
	errorAuth := initAuthFundAdmin()
	zero := decimal.NewFromInt(0)
	if errorAuth != nil {
		log.Error("User Autorizer")
		return lib.CustomError(http.StatusUnauthorized, "User Not Allowed to access this page", "User Not Allowed to access this page")
	}
	var err error
	var status int
	var offset uint64

	params := make(map[string]string)

	//date
	postnavdate := c.QueryParam("nav_date")
	if postnavdate == "" {
		log.Error("Missing required parameter: nav_date")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: nav_date", "Missing required parameter: nav_date")
	}

	params["rec_status"] = "1"
	params["nav_date"] = postnavdate
	params["trans_status_key"] = "6"

	var transTypeKeySubRed []string
	transTypeKeySubRed = append(transTypeKeySubRed, "1")
	transTypeKeySubRed = append(transTypeKeySubRed, "2")

	var transTypeKeySwitch []string
	transTypeKeySwitch = append(transTypeKeySwitch, "4")

	var transSubRed []models.TrTransaction
	status, err = models.GetAllTransactionByParamAndValueIn(&transSubRed, config.LimitQuery, offset, true, params, transTypeKeySubRed, "trans_type_key")
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	var transSwitch []models.TrTransaction
	status, err = models.GetAllTransactionByParamAndValueIn(&transSwitch, config.LimitQuery, offset, true, params, transTypeKeySwitch, "trans_type_key")
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data")
		}
	}

	if (len(transSubRed)) == 0 && (len(transSwitch)) == 0 {
		log.Error("transaction not found")
		return lib.CustomError(http.StatusNotFound, "Transaction not found", "Transaction not found")
	}

	var productIds []string
	var customerIds []string
	var transactionIds []string
	for _, trSubRed := range transSubRed {
		if _, ok := lib.Find(transactionIds, strconv.FormatUint(trSubRed.TransactionKey, 10)); !ok {
			transactionIds = append(transactionIds, strconv.FormatUint(trSubRed.TransactionKey, 10))
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(trSubRed.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(trSubRed.ProductKey, 10))
		}
		if _, ok := lib.Find(customerIds, strconv.FormatUint(trSubRed.CustomerKey, 10)); !ok {
			customerIds = append(customerIds, strconv.FormatUint(trSubRed.CustomerKey, 10))
		}
	}

	var parentIds []string
	for _, trSw := range transSwitch {
		if _, ok := lib.Find(transactionIds, strconv.FormatUint(trSw.TransactionKey, 10)); !ok {
			transactionIds = append(transactionIds, strconv.FormatUint(trSw.TransactionKey, 10))
		}
		strTransType := strconv.FormatUint(trSw.TransTypeKey, 10)
		if strTransType == "4" {
			if trSw.ParentKey != nil {
				if _, ok := lib.Find(transactionIds, strconv.FormatUint(*trSw.ParentKey, 10)); !ok {
					transactionIds = append(transactionIds, strconv.FormatUint(*trSw.ParentKey, 10))
				}
			}
		}
		if _, ok := lib.Find(productIds, strconv.FormatUint(trSw.ProductKey, 10)); !ok {
			productIds = append(productIds, strconv.FormatUint(trSw.ProductKey, 10))
		}
		if _, ok := lib.Find(customerIds, strconv.FormatUint(trSw.CustomerKey, 10)); !ok {
			customerIds = append(customerIds, strconv.FormatUint(trSw.CustomerKey, 10))
		}

		if trSw.ParentKey != nil {
			if _, ok := lib.Find(parentIds, strconv.FormatUint(*trSw.ParentKey, 10)); !ok {
				parentIds = append(parentIds, strconv.FormatUint(*trSw.ParentKey, 10))
			}
		}
	}

	//mapping product
	var productList []models.MsProduct
	if len(productIds) > 0 {
		status, err = models.GetMsProductIn(&productList, productIds, "product_key")
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	productData := make(map[uint64]models.MsProduct)
	for _, pro := range productList {
		productData[pro.ProductKey] = pro
	}

	//mapping tr account
	var accountList []models.TrAccount
	if len(customerIds) > 0 {
		groupBy := "customer_key"
		status, err = models.GetTrAccountIn(&accountList, customerIds, "customer_key", &groupBy)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	accountData := make(map[uint64]models.TrAccount)
	for _, acc := range accountList {
		accountData[acc.CustomerKey] = acc
	}

	//mapping bank transaction
	var bankTrans []models.AdminTransactionBankInfo
	if len(customerIds) > 0 {
		status, err = models.GetTransactionBankInfoCustomerIn(&bankTrans, customerIds)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	bankTransData := make(map[uint64]models.AdminTransactionBankInfo)
	for _, bt := range bankTrans {
		bankTransData[bt.CustomerKey] = bt
	}

	//mapping parent transaction
	var parentTransaction []models.DataTransactionParent
	if len(parentIds) > 0 {
		status, err = models.GetDataParentTransactionSwitch(&parentTransaction, parentIds)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error(err.Error())
				return lib.CustomError(status, err.Error(), "Failed get data")
			}
		}
	}
	parentTrData := make(map[uint64]models.DataTransactionParent)
	for _, pt := range parentTransaction {
		parentTrData[pt.TransactionKey] = pt
	}

	var responseData models.TransactionFormatSinvest

	var dataRowSubRedArray []models.OaRequestCsvFormatFiksTxt

	if len(transSubRed) > 0 {
		txtHeaderSubRed := "transaction_date|transaction_type|sa_code|investor_fund_unit|fund_code|amount_nominal|amount_unit|amount_all_unit|fee_nominal|fee_unit|fee_percent|redm_payment_ac_sequence_code|redm_payment_bank_bic_code|redm_payment_bank_bi_member_code|redm_payment_ac_code|payment_date|transfer_type|sa_reference_no"
		var dataRowSubRed models.OaRequestCsvFormatFiksTxt
		dataRowSubRed.DataRow = txtHeaderSubRed
		dataRowSubRedArray = append(dataRowSubRedArray, dataRowSubRed)

		for _, trSubRed := range transSubRed {
			var txt models.OaRequestCsvFormatFiksTxt
			var subred models.SubscriptionRedeemption
			layout := "2006-01-02 15:04:05"
			newLayout := "20060102"
			date, _ := time.Parse(layout, trSubRed.NavDate)
			subred.TransactionDate = date.Format(newLayout)

			subred.TransactionType = strconv.FormatUint(trSubRed.TransTypeKey, 10)
			subred.SACode = "EP002"

			subred.InvestorFundUnitACNo = ""
			if co, ok := accountData[trSubRed.CustomerKey]; ok {
				subred.InvestorFundUnitACNo = *co.IfuaNo
			}

			subred.FundCode = ""
			if pro, ok := productData[trSubRed.ProductKey]; ok {
				subred.FundCode = *pro.SinvestFundCode
			}

			strTransTypeKey := strconv.FormatUint(trSubRed.TransTypeKey, 10)

			subred.AmountNominal = ""
			if strTransTypeKey == "1" { //SUB
				if trSubRed.TransAmount.Cmp(zero) == 1 {
					strTransAmount := fmt.Sprintf("%.2f", trSubRed.TransAmount)
					subred.AmountNominal = strTransAmount
				} else {
					subred.AmountNominal = "0"
				}
			} else {
				if trSubRed.TransAmount.Cmp(zero) == 1 {
					strTransAmount := fmt.Sprintf("%.2f", trSubRed.TransAmount)
					subred.AmountNominal = strTransAmount
				}
			}

			subred.AmountUnit = ""
			subred.AmountAllUnit = ""
			if strTransTypeKey == "2" { //REDM
				if trSubRed.TransUnit.Cmp(zero) == 1 {
					strTransUnit := fmt.Sprintf("%.2f", trSubRed.TransUnit)
					subred.AmountUnit = strTransUnit
				} else {
					subred.AmountUnit = "0"
				}

				if trSubRed.FlagRedemtAll != nil {
					if int(*trSubRed.FlagRedemtAll) > 0 {
						subred.AmountAllUnit = "Y"
					}
				}
			}

			subred.FeeNominal = ""
			if trSubRed.TransFeeAmount.Cmp(zero) == 1 {
				strFeeAmount := fmt.Sprintf("%.2f", trSubRed.TransFeeAmount)
				subred.FeeNominal = strFeeAmount
			}

			subred.FeeUnit = ""

			subred.FeePercent = ""
			if trSubRed.TransFeePercent.Cmp(zero) == 1 {
				strFeePercent := fmt.Sprintf("%.2f", trSubRed.TransFeePercent)
				subred.FeePercent = strFeePercent
			}

			subred.RedmPaymentACSequenceCode = "1"

			subred.RedmPaymentBankBICCode = ""
			subred.RedmPaymentBankBIMemberCode = ""
			subred.RedmPaymentACCode = ""
			subred.TransferType = ""
			if strTransTypeKey == "2" { //REDM
				if ba, ok := bankTransData[trSubRed.CustomerKey]; ok {
					if ba.SwiftCode != nil {
						subred.RedmPaymentBankBICCode = *ba.SwiftCode
					}
					if ba.BiMemberCode != nil {
						subred.RedmPaymentBankBIMemberCode = *ba.BiMemberCode
					}
					subred.RedmPaymentACCode = ba.CustomerAccountNo
				}
				subred.TransferType = "3"
			}

			if trSubRed.SettlementDate != nil {
				date, _ = time.Parse(layout, *trSubRed.SettlementDate)
				subred.PaymentDate = date.Format(newLayout)
			} else {
				date, _ = time.Parse(layout, trSubRed.NavDate)
				subred.PaymentDate = date.Format(newLayout)
			}
			subred.SaReferenceNo = strconv.FormatUint(trSubRed.TransactionKey, 10)

			txtData := subred.TransactionDate + "|" +
				subred.TransactionType + "|" +
				subred.SACode + "|" +
				subred.InvestorFundUnitACNo + "|" +
				subred.FundCode + "|" +
				subred.AmountNominal + "|" +
				subred.AmountUnit + "|" +
				subred.AmountAllUnit + "|" +
				subred.FeeNominal + "|" +
				subred.FeeUnit + "|" +
				subred.FeePercent + "|" +
				subred.RedmPaymentACSequenceCode + "|" +
				subred.RedmPaymentBankBICCode + "|" +
				subred.RedmPaymentBankBIMemberCode + "|" +
				subred.RedmPaymentACCode + "|" +
				subred.PaymentDate + "|" +
				subred.TransferType + "|" +
				subred.SaReferenceNo

			txt.DataRow = txtData

			dataRowSubRedArray = append(dataRowSubRedArray, txt)
		}
	}

	var dataRowSwitchArray []models.OaRequestCsvFormatFiksTxt

	if len(transSwitch) > 0 {
		txtHeaderSubRed := "transaction_date|transaction_type|sa_code|investor_fund_unit|switch_out_fund_code|switch_out_amount_nominal|switch_out_amount_unit|switch_out_amount_all|switching_fee_charge_fund|fee_nominal|fee_unit|fee_percent|switch_in_fund_code|payment_date|transfer_type|sa_reference_no"
		var dataRowSwitch models.OaRequestCsvFormatFiksTxt
		dataRowSwitch.DataRow = txtHeaderSubRed
		dataRowSwitchArray = append(dataRowSwitchArray, dataRowSwitch)

		for _, trSwitch := range transSwitch {
			var txtSw models.OaRequestCsvFormatFiksTxt
			var swc models.SwitchTransaction

			layout := "2006-01-02 15:04:05"
			newLayout := "20060102"
			date, _ := time.Parse(layout, trSwitch.NavDate)
			swc.TransactionDate = date.Format(newLayout)

			swc.TransactionType = strconv.FormatUint(trSwitch.TransTypeKey, 10)
			swc.SACode = "EP002"

			swc.InvestorFundUnitACNo = ""
			if co, ok := accountData[trSwitch.CustomerKey]; ok {
				swc.InvestorFundUnitACNo = *co.IfuaNo
			}

			swc.SwitchOutFundCode = ""
			swc.SwitchOutAmountNominal = ""
			swc.SwitchOutAmountUnit = ""
			swc.SwitchOutAmountAll = ""
			if trSwitch.ParentKey != nil {
				if pt, ok := parentTrData[*trSwitch.ParentKey]; ok {
					if pt.SinvestFundCode != nil {
						swc.SwitchOutFundCode = *pt.SinvestFundCode
					}

					if pt.TransAmount.Cmp(zero) == 1 {
						strTransAmount := fmt.Sprintf("%.2f", pt.TransAmount)
						swc.SwitchOutAmountNominal = strTransAmount
					}
					if pt.TransUnit.Cmp(zero) == 1 {
						strTransUnit := fmt.Sprintf("%.2f", pt.TransUnit)
						swc.SwitchOutAmountUnit = strTransUnit
					}
					if pt.FlagRedemtAll != nil {
						if int(*pt.FlagRedemtAll) > 0 {
							swc.SwitchOutAmountAll = "Y"
						}
					}
				}
			}

			swc.SwitchingFeeChargeFund = "1"
			if trSwitch.ChargesFeeAmount.Cmp(zero) == 1 {
				swc.SwitchingFeeChargeFund = "2"
			}

			swc.FeeNominal = ""
			if trSwitch.TransFeeAmount.Cmp(zero) == 1 {
				strFeeAmount := fmt.Sprintf("%.2f", trSwitch.TransFeeAmount)
				swc.FeeNominal = strFeeAmount
			}

			swc.FeeUnit = ""

			swc.FeePercent = ""
			if trSwitch.TransFeePercent.Cmp(zero) == 1 {
				strFeePercent := fmt.Sprintf("%.2f", trSwitch.TransFeePercent)
				swc.FeePercent = strFeePercent
			}

			swc.SwitchInFundCode = ""
			if pro, ok := productData[trSwitch.ProductKey]; ok {
				swc.SwitchInFundCode = *pro.SinvestFundCode
			}

			if trSwitch.SettlementDate != nil {
				date, _ = time.Parse(layout, *trSwitch.SettlementDate)
				swc.PaymentDate = date.Format(newLayout)
			} else {
				date, _ = time.Parse(layout, trSwitch.NavDate)
				swc.PaymentDate = date.Format(newLayout)
			}

			swc.TransferType = "3"
			swc.SaReferenceNo = strconv.FormatUint(trSwitch.TransactionKey, 10)

			txtData := swc.TransactionDate + "|" +
				swc.TransactionType + "|" +
				swc.SACode + "|" +
				swc.InvestorFundUnitACNo + "|" +
				swc.SwitchOutFundCode + "|" +
				swc.SwitchOutAmountNominal + "|" +
				swc.SwitchOutAmountUnit + "|" +
				swc.SwitchOutAmountAll + "|" +
				swc.SwitchingFeeChargeFund + "|" +
				swc.FeeNominal + "|" +
				swc.FeeUnit + "|" +
				swc.FeePercent + "|" +
				swc.SwitchInFundCode + "|" +
				swc.PaymentDate + "|" +
				swc.TransferType + "|" +
				swc.SaReferenceNo

			txtSw.DataRow = txtData

			dataRowSwitchArray = append(dataRowSwitchArray, txtSw)
		}
	}

	responseData.SubscriptionRedeemption = &dataRowSubRedArray
	responseData.SwitchTransaction = &dataRowSwitchArray

	if len(transactionIds) > 0 {
		for _, trID := range transactionIds {
			paramsUpdate := make(map[string]string)
			transBatchKey := CheckTransactionBatch(trID)
			paramsUpdate["trans_status_key"] = "7"
			dateLayout := "2006-01-02 15:04:05"
			paramsUpdate["rec_modified_date"] = time.Now().Format(dateLayout)
			strKey := strconv.FormatUint(lib.Profile.UserID, 10)
			paramsUpdate["rec_modified_by"] = strKey
			paramsUpdate["batch_key"] = transBatchKey
			paramsUpdate["transaction_key"] = trID

			_, err = models.UpdateTrTransaction(paramsUpdate)
			if err != nil {
				log.Error("Error update Transaction")
				return lib.CustomError(http.StatusInternalServerError, err.Error(), "Failed update data")
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

func CheckTransactionBatch(transKey string) string {
	transBatchKet := ""

	var trBatch models.CheckTrTransactionBacth
	_, err := models.CheckTrTransactionBatchByTransKey(&trBatch, transKey)
	if err == nil {
		transBatchKet = strconv.FormatUint(trBatch.BatchKey, 10)
	} else {
		//insert into batch
		var paramBatch models.ParamBatchTrTransaction
		_, err := models.ParamBatchTrTransactionByKey(&paramBatch, transKey)
		if err == nil {
			dateLayout := "2006-01-02 15:04:05"
			paramsBatch := make(map[string]string)
			if paramBatch.Batch == nil {
				paramsBatch["batch_number"] = "1"
			} else {
				paramsBatch["batch_number"] = strconv.FormatUint(*paramBatch.Batch+1, 10)
			}
			bulanRomawi := ""
			if paramBatch.Bulan == "1" {
				bulanRomawi = "I"
			}
			if paramBatch.Bulan == "2" {
				bulanRomawi = "II"
			}
			if paramBatch.Bulan == "3" {
				bulanRomawi = "III"
			}
			if paramBatch.Bulan == "4" {
				bulanRomawi = "IV"
			}
			if paramBatch.Bulan == "5" {
				bulanRomawi = "V"
			}
			if paramBatch.Bulan == "6" {
				bulanRomawi = "VI"
			}
			if paramBatch.Bulan == "7" {
				bulanRomawi = "VII"
			}
			if paramBatch.Bulan == "8" {
				bulanRomawi = "VIII"
			}
			if paramBatch.Bulan == "9" {
				bulanRomawi = "IX"
			}
			if paramBatch.Bulan == "10" {
				bulanRomawi = "X"
			}
			if paramBatch.Bulan == "11" {
				bulanRomawi = "XI"
			}
			if paramBatch.Bulan == "12" {
				bulanRomawi = "XII"
			}

			batchDisplayNo := paramBatch.ProductCode + "/" + paramBatch.TypeCode + "/" + bulanRomawi + "/" + paramBatch.Tahun

			paramsBatch["batch_display_no"] = batchDisplayNo
			paramsBatch["nav_date"] = paramBatch.NavDate
			paramsBatch["product_key"] = strconv.FormatUint(paramBatch.ProductKey, 10)
			paramsBatch["trans_type_key"] = strconv.FormatUint(paramBatch.TransTypeKey, 10)
			paramsBatch["rec_order"] = "0"
			paramsBatch["rec_status"] = "1"
			paramsBatch["rec_created_date"] = time.Now().Format(dateLayout)
			strKey := strconv.FormatUint(lib.Profile.UserID, 10)
			paramsBatch["rec_created_by"] = strKey

			_, err, lastID := models.CreateTrTransactionBacth(paramsBatch)
			if err != nil {
				log.Error(err.Error())
				log.Error("Error create batch")
			}

			transBatchKet = lastID
		}
	}

	return transBatchKet
}
