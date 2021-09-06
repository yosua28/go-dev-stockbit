package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/labstack/echo"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func CreateTransaction(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	params := make(map[string]string)
	paramsTransaction := make(map[string]string)
	zero := decimal.NewFromInt(0)
	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}

	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	params["customer_key"] = customerKey
	params["rec_status"] = "1"

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

	productIDs := []string{productKeyStr}

	var navDB []models.TrNav
	status, err = models.GetLastNavIn(&navDB, productIDs)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed get nav data")
	}

	paramsAcc := make(map[string]string)
	paramsAcc["customer_key"] = strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	paramsAcc["product_key"] = productKeyStr
	paramsAcc["rec_status"] = "1"
	var balanceUnit decimal.Decimal
	var investValue decimal.Decimal
	var accDB []models.TrAccount
	status, err = models.GetAllTrAccount(&accDB, paramsAcc)
	if err != nil {
		log.Error(err.Error())
	}

	var accIDs []string
	accProduct := make(map[uint64]uint64)
	acaProduct := make(map[uint64]uint64)
	var acaDB []models.TrAccountAgent
	if len(accDB) > 0 {
		for _, acc := range accDB {
			accIDs = append(accIDs, strconv.FormatUint(acc.AccKey, 10))
			accProduct[acc.AccKey] = acc.ProductKey
		}
		status, err = models.GetTrAccountAgentIn(&acaDB, accIDs, "acc_key")
		if err != nil {
			log.Error(err.Error())
		}
		if len(acaDB) > 0 {
			var acaIDs []string
			for _, aca := range acaDB {
				acaIDs = append(acaIDs, strconv.FormatUint(aca.AcaKey, 10))
				acaProduct[aca.AcaKey] = aca.AccKey
			}
			var balanceDB []models.TrBalance
			status, err = models.GetLastBalanceIn(&balanceDB, acaIDs)
			if err != nil {
				log.Error(err.Error())
			}
			if len(balanceDB) > 0 {
				for _, balance := range balanceDB {
					balanceUnit = balanceUnit.Add(balance.BalanceUnit)
				}
				investValue = navDB[0].NavValue.Mul(balanceUnit)
			}
		}
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
	typeStr := ""
	parentKeyStr := "NULL"
	if typeKeyStr == "4" {
		typeStr = "switching"
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
	var unitValue decimal.Decimal
	if transUnitStr != "" {
		unitValue, err = decimal.NewFromString(transUnitStr)
		if err != nil {
			log.Error("Wrong input for parameter: trans_unit")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_unit", "Wrong input for parameter: trans_unit")
		}

		if typeKeyStr == "2" {
			typeStr = "redemption"
			if unitValue.Cmp(balanceUnit) == 1 {
				log.Error("unit redemp > balance")
				return lib.CustomError(http.StatusBadRequest, "unit redemp > balance", "Jumlah unit yang di redem melebihi balance yang ada")
			}
			if unitValue.Cmp(product.MinRedUnit) == -1 {
				log.Error("red unit < minimum red")
				return lib.CustomError(http.StatusBadRequest, "red unit < minum red", "Minimum redemption untuk product ini adalah: "+product.MinRedUnit.Truncate(2).String()+"unit")
			}
			if flagRedemAll == false {
				if balanceUnit.Sub(unitValue).Cmp(product.MinUnitAfterRed) == -1 {
					log.Error("unit after redemption < minimum unit after red")
					return lib.CustomError(http.StatusBadRequest, "unit after redemption < minimum unit after red", "Minumum unit setelah redemption untuk product ini adalah: "+product.MinUnitAfterRed.Truncate(2).String()+"unit")
				}
			}
		}
	} else {
		log.Error("Missing required parameter: trans_unit")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_unit", "Missing required parameter: trans_unit")
	}

	transAmountStr := c.FormValue("trans_amount")
	if transAmountStr != "" {
		value, err := decimal.NewFromString(transAmountStr)
		if err != nil {
			log.Error("Wrong input for parameter: trans_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_amount", "Wrong input for parameter: trans_amount")
		}
		if typeKeyStr == "1" {
			typeStr = "subscription"
			if balanceUnit.Cmp(zero) == 1 {
				typeStr = "topup"
			}
			if value.Cmp(product.MinSubAmount) == -1 {
				log.Error("sub amount < minimum sub")
				return lib.CustomError(http.StatusBadRequest, "sub amount < minum sub", "Minumum subscription untuk product ini adalah: "+product.MinSubAmount.Truncate(2).String())
			}
		} else if typeKeyStr == "2" {
			typeStr = "redemption"
			if value.Cmp(investValue) == 1 {
				log.Error("Amount redemp > invest value")
				return lib.CustomError(http.StatusBadRequest, "amount redemp > invest value", "Jumlah redem melebihi total invest value untuk product ini")
			}
			if unitValue == zero && value.Cmp(product.MinRedAmount) == -1 {
				log.Error("red amount < minimum red")
				return lib.CustomError(http.StatusBadRequest, "red amount < minimum red", "Minumum redemption untuk product ini adalah: "+product.MinRedAmount.Truncate(2).String())
			}
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

	paymentChannel := c.FormValue("payment_channel")
	if paymentChannel != "" {
		_, err := strconv.ParseUint(paymentChannel, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: payment_channel")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: payment_channel", "Wrong input for parameter: payment_channel")
		}
	} else {
		log.Error("Missing required parameter: payment_channel")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: payment_channel", "Missing required parameter: payment_channel")
	}
	paymentMethod := c.FormValue("payment_method")
	if paymentMethod != "" {
		_, err := strconv.ParseUint(paymentMethod, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: payment_method")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: payment_method", "Wrong input for parameter: payment_method")
		}
	} else {
		log.Error("Missing required parameter: payment_method")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: payment_method", "Missing required parameter: payment_method")
	}

	var promoKey *string
	promoCode := c.FormValue("promo_code")
	if promoCode != "" {
		err, enable, text, promoKeyRes := validatePromo(promoCode, customerKey, productKeyStr)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
		} else {
			if enable == false {
				return lib.CustomError(http.StatusBadRequest, text, text)
			} else {
				promoKey = promoKeyRes
			}
		}
	}

	transCalcMethod := c.FormValue("trans_calc_method")
	if transCalcMethod != "" {
		transCalcMethodKey, _ := strconv.ParseUint(transCalcMethod, 10, 64)
		if transCalcMethodKey == 0 {
			log.Error("Wrong value for parameter: trans_calc_method")
			return lib.CustomError(http.StatusBadRequest, "Wrong value for parameter: trans_calc_method", "Missing required parameter: trans_calc_method")
		}
	}

	var accKey string
	var trAccountDB []models.TrAccount
	status, err = models.GetAllTrAccount(&trAccountDB, params)
	if len(trAccountDB) > 0 {
		accKey = strconv.FormatUint(trAccountDB[0].AccKey, 10)
		if (typeKeyStr == "1" || typeKeyStr == "4") && (trAccountDB[0].SubSuspendFlag != nil && *trAccountDB[0].SubSuspendFlag == 1) ||
			(typeKeyStr == "2" || typeKeyStr == "3") && (trAccountDB[0].RedSuspendFlag != nil && *trAccountDB[0].RedSuspendFlag == 1) {
			log.Error("Account suspended")
			return lib.CustomError(status, "Account suspended", "Account suspended")
		}
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
	params["trans_calc_method"] = transCalcMethod

	layout := "2006-01-02"
	now := time.Now()
	nowString := now.Format(layout)
	t, _ := time.Parse(layout, now.AddDate(0, 0, -1).Format(layout))
	dateBursa := SettDate(t, int(1))
	navDate := ""
	if nowString == dateBursa && (now.Hour() == 12 && now.Minute() > 0) || now.Hour() > 12 {
		t, _ := time.Parse(layout, dateBursa)
		navDate = SettDate(t, int(1)) + " 00:00:00"
	} else {
		navDate = dateBursa + " 00:00:00"
	}
	params["nav_date"] = navDate
	params["entry_mode"] = "140"
	if transAmountStr == "0" {
		params["entry_mode"] = "139"
	}

	params["trans_amount"] = transAmountStr
	params["trans_unit"] = transUnitStr
	params["trans_type_key"] = typeKeyStr
	if flagRedemAll {
		params["flag_redempt_all"] = redemAllStr
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
	params["payment_method"] = paymentChannel
	params["aca_key"] = acaKey
	params["rec_status"] = "1"

	settlementParams := make(map[string]string)
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
					settlementParams["rec_image1"] = filename + extension
					dateLayout := "2006-01-02 15:04:05"
					params["file_upload_date"] = time.Now().Format(dateLayout)
					settlementParams["settle_realized_date"] = time.Now().Format(dateLayout)
				}
			}
		}
	}

	productBankAccountKey := c.FormValue("product_bankacc_key")
	if typeKeyStr == "1" {
		if productBankAccountKey != "" {
			productBankAccount, _ := strconv.ParseUint(productBankAccountKey, 10, 64)
			if productBankAccount > 0 {
				paramsTransaction["prod_bankacc_key"] = productBankAccountKey
			} else {
				log.Error("Wrong input for parameter: product_bankacc_key")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_bankacc_key", "Wrong input for parameter: product_bankacc_key")
			}
		} else {
			log.Error("Missing required parameter: product_bankacc_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_bankacc_key", "Missing required parameter: product_bankacc_key")
		}
	} else {
		purpose := "270"
		if typeKeyStr == "4" {
			purpose = "269"
		}
		paramsProBankAcc := make(map[string]string)
		paramsProBankAcc["bank_account_purpose"] = purpose
		paramsProBankAcc["product_key"] = productKeyStr
		var productBankDB []models.MsProductBankAccount
		status, err = models.GetAllMsProductBankAccount(&productBankDB, paramsProBankAcc)
		if err != nil {
			log.Error(err.Error())
			paramsTransaction["prod_bankacc_key"] = "1"
		} else {
			paramsTransaction["prod_bankacc_key"] = strconv.FormatUint(productBankDB[0].ProdBankaccKey, 10)
		}
	}

	if promoKey != nil {
		params["promo_code"] = promoCode
	}
	params["rec_attribute_id3"] = c.Request().UserAgent()
	status, err, transactionID := models.CreateTrTransaction(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "Failed input data")
	}

	//save to promo used
	if promoKey != nil {
		paramsPromoUsed := make(map[string]string)
		paramsPromoUsed["used_date"] = time.Now().Format(dateLayout)
		paramsPromoUsed["promo_key"] = *promoKey
		paramsPromoUsed["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
		paramsPromoUsed["customer_key"] = customerKey
		paramsPromoUsed["transaction_key"] = transactionID
		paramsPromoUsed["used_status"] = "317"
		paramsPromoUsed["rec_status"] = "1"
		paramsPromoUsed["rec_created_date"] = time.Now().Format(dateLayout)
		paramsPromoUsed["rec_created_by"] = strconv.FormatUint(lib.Profile.UserID, 10)
		_, err := models.CreateTrPromoUsed(paramsPromoUsed)
		if err != nil {
			log.Error(err.Error())
		}
	}

	paramsTransaction["transaction_key"] = transactionID
	customerBankAccountKey := c.FormValue("customer_bankacc_key")
	if typeKeyStr == "2" {
		if customerBankAccountKey != "" {
			productBankAccount, _ := strconv.ParseUint(customerBankAccountKey, 10, 64)
			if productBankAccount > 0 {
				paramsTransaction["cust_bankacc_key"] = customerBankAccountKey
			} else {
				log.Error("Wrong input for parameter: customer_bankacc_key")
				return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: customer_bankacc_key", "Wrong input for parameter: customer_bankacc_key")
			}
		} else {
			log.Error("Missing required parameter: customer_bankacc_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: customer_bankacc_key", "Missing required parameter: customer_bankacc_key")
		}
	} else {
		var customerBankDB []models.MsCustomerBankAccount
		paramCustomerBank := make(map[string]string)
		paramCustomerBank["customer_key"] = customerKey
		paramCustomerBank["orderBy"] = "cust_bankacc_key"
		paramCustomerBank["orderType"] = "DESC"
		status, err = models.GetAllMsCustomerBankAccount(&customerBankDB, paramCustomerBank)
		if err != nil {
			log.Error(err.Error())
			paramsTransaction["cust_bankacc_key"] = "1"
		} else {
			paramsTransaction["cust_bankacc_key"] = strconv.FormatUint(customerBankDB[0].CustBankaccKey, 10)
		}
	}

	if paymentChannel == "299" {
		orderID := c.FormValue("order_id")
		if orderID == "" {
			log.Error("Missing required parameter: order_id")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: order_id", "Missing required parameter: order_id")
		}
		referenceNumber := c.FormValue("reference_number")
		if referenceNumber == "" {
			log.Error("Missing required parameter: reference_number")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: reference_number", "Missing required parameter: reference_number")
		}
		otp := c.FormValue("otp")
		if otp == "" {
			log.Error("Missing required parameter: otp")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: otp", "Missing required parameter: otp")
		}
		params := make(map[string]string)
		params["order_id"] = orderID
		params["phone"] = lib.Profile.PhoneNumber
		params["otp_private_code"] = otp
		params["payment_method"] = "SPINPAY"
		//params["payment_data"] = "Transaction " + params["reference_code"]

		status, res, err := lib.Spin(orderID, "PAY_ORDER", params)
		if err != nil {
			log.Error(status, err.Error())
			return lib.CustomError(status, err.Error(), "Pay order failed")
		}
		if status != 200 {
			log.Error(status, "Pay order failed")
			return lib.CustomError(status, "Pay order failed", "Pay order failed")
		}
		var sec map[string]interface{}
		if err = json.Unmarshal([]byte(res), &sec); err != nil {
			log.Error("Error4", err.Error())
			return lib.CustomError(http.StatusBadGateway, err.Error(), "Parsing data failed")
		}

		data := sec["message_action"].(string)
		if data == "PAYMENT_SUCCESS" {
			settlementParams["settle_realized_date"] = time.Now().Format(dateLayout)
			settlementParams["settle_response"] = res
			settlementParams["settle_remarks"] = orderID
			settlementParams["settle_references"] = referenceNumber
		}
	}

	paramsTransaction["rec_status"] = "1"
	status, err = models.CreateTrTransactionBankAccount(paramsTransaction)
	if err != nil {
		log.Error(err.Error())
	}

	if typeKeyStr == "1" {

		settlementParams["transaction_key"] = transactionID
		settlementParams["settle_purposed"] = "297"
		settlementParams["settle_date"] = navDate
		settlementParams["settle_nominal"] = totalAmountStr
		settlementParams["client_subaccount_no"] = ""
		if paymentChannel != "323" {
			subAcc := c.FormValue("sub_acc")
			if subAcc == "" {
				log.Error("Missing required parameter: sub_acc")
				return lib.CustomError(http.StatusBadRequest, "Missing required parameter: sub_acc", "Missing required parameter: sub_acc")
			} else {
				settlementParams["client_subaccount_no"] = subAcc
			}
		}
		settlementParams["settled_status"] = "243"
		settlementParams["source_bank_account_key"] = paramsTransaction["cust_bankacc_key"]
		settlementParams["target_bank_account_key"] = productBankAccountKey
		settlementParams["settle_channel"] = paymentChannel
		settlementParams["settle_payment_method"] = paymentMethod
		settlementParams["rec_status"] = "1"

		_, err, _ = models.CreateTrTransactionSettlement(settlementParams)
		if err != nil {
			log.Error(err.Error())
		}
	}

	if typeKeyStr != "3" {
		params["product_name"] = product.ProductNameAlt
		params["currency"] = strconv.FormatUint(*product.CurrencyKey, 10)
		params["parrent"] = parentKeyStr
		params["customer_bank_account_key"] = customerBankAccountKey
		err = mailTransaction(typeStr, params)
	}

	//insert message notif in app
	if typeKeyStr != "3" {
		strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)
		paramsUserMessage := make(map[string]string)
		paramsUserMessage["umessage_type"] = "245"
		paramsUserMessage["umessage_recipient_key"] = strIDUserLogin
		paramsUserMessage["umessage_receipt_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["flag_read"] = "0"
		paramsUserMessage["umessage_sent_date"] = time.Now().Format(dateLayout)
		paramsUserMessage["flag_sent"] = "1"
		var subject string
		var body string
		if typeKeyStr == "1" { // SUBS
			if params["flag_newsub"] == "1" {
				subject = "Subscription sedang Diproses"
				body = "Terima kasih telah melakukan subscription. Kami sedang memproses transaksi kamu."
			} else {
				subject = "Top Up sedang Diproses"
				body = "Terima kasih telah melakukan transaksi top up. Kami sedang memproses transaksi kamu."
			}
		}

		if typeKeyStr == "2" { // REDM
			subject = "Redemption sedang Diproses"
			body = "Redemption kamu telah kami terima. Kami akan memproses transaksi kamu."
		}
		if typeKeyStr == "4" || typeKeyStr == "3" { // SWITCH
			subject = "Switching sedang Diproses"
			body = "Switching kamu telah kami terima. Kami sedang memproses transaksi kamu."
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
		lib.CreateNotifCustomerFromAdminByUserLoginKey(strIDUserLogin, subject, body, "TRANSACTION")
	}

	//sent email to BO role 11 & Sales
	if typeKeyStr != "3" {
		if typeKeyStr == "1" { //subs
			if _, ok := params["rec_image1"]; ok { //jika upload image
				SentEmailTransactionToBackOfficeAndSales(transactionID, "11")
			}
		} else if typeKeyStr == "2" { //redm
			SentEmailTransactionToBackOfficeAndSales(transactionID, "11")
		} else if typeKeyStr == "4" { //switching
			SentEmailTransactionToBackOfficeAndSales(parentKeyStr, "11")
		}
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

func Subscription(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	params := make(map[string]string)

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}

	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)
	var customer models.MsCustomer
	_, err = models.GetMsCustomer(&customer, customerKey)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Data customer tidak ditemukan")
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

	transCalcMethod := c.FormValue("trans_calc_method")
	if transCalcMethod != "" {
		transCalcMethodKey, err := strconv.ParseUint(transCalcMethod, 10, 64)
		if err == nil && transCalcMethodKey > 0 {
			params["trans_calc_method"] = transCalcMethod
		} else {
			log.Error("Missing required parameter: trans_calc_method")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_calc_method", "Missing required parameter: trans_calc_method")
		}
	} else {
		log.Error("Missing required parameter: trans_calc_method")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_calc_method", "Missing required parameter: trans_calc_method")
	}

	transFeePercentStr := c.FormValue("trans_fee_percent")
	if transFeePercentStr != "" {
		_, err := decimal.NewFromString(transFeePercentStr)
		if err != nil {
			log.Error("Wrong input for parameter: trans_fee_percent")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_fee_percent", "Wrong input for parameter: trans_fee_percent")
		}
		params["trans_fee_percent"] = transFeePercentStr
	} else {
		log.Error("Missing required parameter: trans_fee_percent")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_fee_percent", "Missing required parameter: trans_fee_percent")
	}

	transFeeAmountStr := c.FormValue("trans_fee_amount")
	if transFeeAmountStr != "" {
		_, err := decimal.NewFromString(transFeeAmountStr)
		if err != nil {
			log.Error("Wrong input for parameter: trans_fee_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: trans_fee_amount", "Wrong input for parameter: trans_fee_amount")
		}
		params["trans_fee_amount"] = transFeeAmountStr
	} else {
		log.Error("Missing required parameter: trans_fee_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: trans_fee_amount", "Missing required parameter: trans_fee_amount")
	}

	chargesFeeAmountStr := c.FormValue("charges_fee_amount")
	if chargesFeeAmountStr != "" {
		_, err := decimal.NewFromString(chargesFeeAmountStr)
		if err != nil {
			log.Error("Wrong input for parameter: charges_fee_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: charges_fee_amount", "Wrong input for parameter: charges_fee_amount")
		}
		params["charges_fee_amount"] = chargesFeeAmountStr
	} else {
		log.Error("Missing required parameter: charges_fee_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: charges_fee_amount", "Missing required parameter: charges_fee_amount")
	}

	servicesFeeAmountStr := c.FormValue("services_fee_amount")
	if servicesFeeAmountStr != "" {
		_, err := decimal.NewFromString(servicesFeeAmountStr)
		if err != nil {
			log.Error("Wrong input for parameter: services_fee_amount")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: services_fee_amount", "Wrong input for parameter: services_fee_amount")
		}
		params["services_fee_amount"] = servicesFeeAmountStr
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
		params["total_amount"] = totalAmountStr
	} else {
		log.Error("Missing required parameter: total_amount")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: total_amount", "Missing required parameter: total_amount")
	}

	var promoKey *string
	promoCode := c.FormValue("promo_code")
	if promoCode != "" {
		err, enable, text, promoKeyRes := validatePromo(promoCode, customerKey, productKeyStr)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
		} else {
			if enable == false {
				return lib.CustomError(http.StatusBadRequest, text, text)
			} else {
				promoKey = promoKeyRes
			}
		}
	}

	paymentChannel := c.FormValue("payment_channel")
	if paymentChannel != "" {
		_, err := strconv.ParseUint(paymentChannel, 10, 64)
		if err != nil {
			log.Error("Wrong input for parameter: payment_channel")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: payment_channel", "Wrong input for parameter: payment_channel")
		}
	} else {
		log.Error("Missing required parameter: payment_channel")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: payment_channel", "Missing required parameter: payment_channel")
	}

	paymentMethod := c.FormValue("payment_method")
	if paymentMethod != "" {
		paymentKey, err := strconv.ParseUint(paymentMethod, 10, 64)
		if err == nil && paymentKey > 0 {
			params["payment_method"] = paymentMethod
		} else {
			log.Error("Missing required parameter: payment_method")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: payment_method", "Missing required parameter: payment_method")
		}
	} else {
		log.Error("Missing required parameter: payment_method")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: payment_method", "Missing required parameter: payment_method")
	}

	bankStr := c.FormValue("product_bankacc_key")
	if bankStr != "" {
		bankKey, err := strconv.ParseUint(bankStr, 10, 64)
		if err != nil || bankKey == 0 {
			log.Error("Missing required parameter: product_bankacc_key")
			return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_bankacc_key", "Missing required parameter: product_bankacc_key")
		}
	} else {
		log.Error("Missing required parameter: product_bankacc_key")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_bankacc_key", "Missing required parameter: product_bankacc_key")
	}

	transRemark := c.FormValue("trans_remarks")
	params["trans_remarks"] = transRemark

	dateLayout := "2006-01-02 15:04:05"
	strIDUserLogin := strconv.FormatUint(lib.Profile.UserID, 10)

	//cek tr_account / save
	var accKey string
	paramsAcc := make(map[string]string)
	paramsAcc["customer_key"] = customerKey
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
	
	var agentCustomerDB models.MsAgentCustomer
	status, err = models.GetLastAgenCunstomer(&agentCustomerDB, customerKey)
	if err != nil {
		log.Error("Failed get data agent: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	agentkey := strconv.FormatUint(agentCustomerDB.AgentKey, 10)

	//cek tr_account_agent / save
	paramsAccAgent := make(map[string]string)
	paramsAccAgent["acc_key"] = accKey
	paramsAccAgent["agent_key"] = agentkey
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
		paramsCreateAccAgent["agent_key"] = agentkey
		paramsCreateAccAgent["branch_key"] = "1"
		paramsCreateAccAgent["rec_status"] = "1"
		status, err, acaKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
		if err != nil {
			log.Error("Failed create account agent data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}
	//save tr_transaction

	params["branch_key"] = "1"
	params["customer_key"] = customerKey
	params["agent_key"] = agentkey
	params["trans_status_key"] = "2"
	params["trans_date"] = time.Now().Format(dateLayout)
	params["trans_type_key"] = "1"
	params["trx_code"] = "137"
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
	params["aca_key"] = acaKey
	params["trans_source"] = "141"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strIDUserLogin
	params["risk_waiver"] = "0"

	var riskProfil models.RiskProfilCustomer
	status, err = models.GetRiskProfilCustomer(&riskProfil, customerKey)
	if err != nil {
		if product.RiskProfileKey != nil {
			if riskProfil.RiskProfileKey < *product.RiskProfileKey {
				params["risk_waiver"] = "1"
			}
		}
	}

	var userData models.ScUserLogin
	status, err = models.GetScUserLoginByCustomerKey(&userData, customerKey)
	if err != nil {
		return lib.CustomError(status)
	}

	settlementParams := make(map[string]string)
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
				settlementParams["rec_image1"] = filename + extension
				dateLayout := "2006-01-02 15:04:05"
				params["file_upload_date"] = time.Now().Format(dateLayout)
				settlementParams["settle_realized_date"] = time.Now().Format(dateLayout)
			}
		}
	}

	if promoKey != nil {
		params["promo_code"] = promoCode
	}

	status, err, transactionID := models.CreateTrTransaction(params)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}

	//save to promo used
	if promoKey != nil {
		paramsPromoUsed := make(map[string]string)
		paramsPromoUsed["used_date"] = time.Now().Format(dateLayout)
		paramsPromoUsed["promo_key"] = *promoKey
		paramsPromoUsed["user_login_key"] = strIDUserLogin
		paramsPromoUsed["customer_key"] = customerKey
		paramsPromoUsed["transaction_key"] = transactionID
		paramsPromoUsed["used_status"] = "317"
		paramsPromoUsed["rec_status"] = "1"
		paramsPromoUsed["rec_created_date"] = time.Now().Format(dateLayout)
		paramsPromoUsed["rec_created_by"] = strIDUserLogin
		_, err := models.CreateTrPromoUsed(paramsPromoUsed)
		if err != nil {
			log.Error(err.Error())
		}
	}

	//save tr_transaction_bank_account
	paramsBankTransaction := make(map[string]string)
	paramsBankTransaction["transaction_key"] = transactionID
	paramsBankTransaction["prod_bankacc_key"] = bankStr
	paramsBankTransaction["rec_status"] = "1"

	var customerBankDB []models.MsCustomerBankAccount
	paramCustomerBank := make(map[string]string)
	paramCustomerBank["customer_key"] = customerKey
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

	//create to tr_transaction_settlement
	var subAcc string
	var responseFM interface{}
	if paymentMethod == "287" {
		subAcc = "7029123" + lib.Profile.PhoneNumber
		fmParams := make(map[string]string)
		fmParams["merchant_code"] = config.MerchantCode
		if customer.FirstName != nil && *customer.FirstName != "" {
			fmParams["first_name"] = *customer.FirstName
		} else {
			fmParams["first_name"] = customer.FullName
		}
		if customer.LastName != nil && *customer.LastName != "" {
			fmParams["last_name"] = *customer.LastName
		} else {
			fmParams["last_name"] = ""
		}
		fmParams["email"] = lib.Profile.Email
		fmParams["phone"] = lib.Profile.PhoneNumber
		fmParams["order_id"] = transactionID
		fmParams["no_reference"] = lib.Profile.PhoneNumber
		fmParams["amount"] = totalAmountStr
		fmParams["currency"] = "IDR"
		fmParams["item_details"] = typ + " product " + productKeyStr 
		fmParams["datetime_request"] = time.Now().Format(dateLayout) 
		fmParams["payment_method"] = "va_mandiri" 
		fmParams["time_limit"] = "1440" 
		fmParams["notif_url"] = config.BaseUrl + "/fmnotif" 
		fmParams["thanks_url"] = config.BaseUrl + "/fmthankyou" 
		_, responseFM, err = lib.FMPostPaymentData(fmParams)
		if err != nil {
			log.Error("Error POST payment data to FM: ",err.Error())
		}
	}
	settlementParams["transaction_key"] = transactionID
	settlementParams["settle_purposed"] = "297"
	settlementParams["settle_date"] = dateBursa + " 00:00:00"
	settlementParams["settle_nominal"] = totalAmountStr
	settlementParams["client_subaccount_no"] = subAcc
	settlementParams["settled_status"] = "243"
	settlementParams["target_bank_account_key"] = bankStr
	settlementParams["settle_channel"] = paymentChannel
	settlementParams["settle_payment_method"] = paymentMethod
	settlementParams["rec_status"] = "1"
	settlementParams["rec_created_date"] = time.Now().Format(dateLayout)
	settlementParams["rec_created_by"] = strIDUserLogin

	_, err, _ = models.CreateTrTransactionSettlement(settlementParams)
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
	lib.CreateNotifCustomerFromAdminByCustomerId(customerKey, subject, body, "TRANSACTION")

	//send email
	params["product_name"] = product.ProductNameAlt
	params["currency"] = strconv.FormatUint(*product.CurrencyKey, 10)
	params["parrent"] = transactionID
	err = mailSubscription(typ, params)
	responseData := make(map[string]interface{})
	responseData["transaction_key"] = transactionID
	responseData["response_fm"] = responseFM
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

func Redemption(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	params := make(map[string]string)

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}

	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)

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
	status, err = models.CheckProductAllowRedmOrSwitching(&productNotAllow, customerKey, productIds)
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
	status, err = models.GetBalanceUnitByCustomerAndProduct(&balance, customerKey, productKeyStr)
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

	metodePerhitungan := c.FormValue("calc_method")
	if metodePerhitungan != "" {
		if metodePerhitungan == "all" { //all unit
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
				return lib.CustomError(http.StatusBadRequest, "red unit < minimum red unit", "Minumum redemption unit untuk product ini adalah: "+product.MinRedUnit.String())
			}

			if value.Cmp(unitTersedia) == 1 {
				log.Error("red unit > unit tersedia")
				return lib.CustomError(http.StatusBadRequest, "red unit > unit tersedia", "Redemption unit tidak boleh lebih besar dari unit tersedia. Unit tersedia saat ini adalah: "+balance.Unit.String())
			}

			params["trans_unit"] = transUnitStr
			params["total_amount"] = "0"
		} else if metodePerhitungan == "unit" { //unit penyertaan
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

			if sisaUnitAfterRed != zero && sisaUnitAfterRed.Cmp(minSisa) == -1 {
				log.Error("Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All")
				return lib.CustomError(http.StatusBadRequest, "Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All. Sisa unit harus minimal : "+minSisa.String(), "Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All. Sisa unit harus minimal : "+minSisa.String())
			}

			params["trans_unit"] = transUnitStr
			params["total_amount"] = "0"
		} else if metodePerhitungan == "amount" { //Nominal
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

			if sisaUnitAfterRed != zero && sisaUnitAfterRed.Cmp(minSisa) == -1 {
				log.Error("Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All")
				return lib.CustomError(http.StatusBadRequest, "Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All. Sisa unit harus minimal : "+minSisa.String(), "Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All. Sisa unit harus minimal : "+minSisa.String())
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
	paramsAcc["customer_key"] = customerKey
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
	var agentCustomerDB models.MsAgentCustomer
	status, err = models.GetLastAgenCunstomer(&agentCustomerDB, customerKey)
	if err != nil {
		log.Error("Failed get data agent: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	agentkey := strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	//cek tr_account_agent / save
	paramsAccAgent := make(map[string]string)
	paramsAccAgent["acc_key"] = accKey
	paramsAccAgent["agent_key"] = agentkey
	paramsAccAgent["rec_status"] = "1"

	// var acaKey string
	// var accountAgentDB []models.TrAccountAgent
	// status, err = models.GetAllTrAccountAgent(&accountAgentDB, paramsAccAgent)
	// if len(accountAgentDB) > 0 {
	// 	acaKey = strconv.FormatUint(accountAgentDB[0].AcaKey, 10)
	// } else {
	// 	paramsCreateAccAgent := make(map[string]string)
	// 	paramsCreateAccAgent["acc_key"] = accKey
	// 	paramsCreateAccAgent["eff_date"] = time.Now().Format(dateLayout)
	// 	paramsCreateAccAgent["rec_created_date"] = time.Now().Format(dateLayout)
	// 	paramsCreateAccAgent["rec_created_by"] = strIDUserLogin
	// 	paramsCreateAccAgent["agent_key"] = agentkey
	// 	paramsCreateAccAgent["branch_key"] = "1"
	// 	paramsCreateAccAgent["rec_status"] = "1"
	// 	status, err, acaKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
	// 	if err != nil {
	// 		log.Error("Failed create account agent data: " + err.Error())
	// 		return lib.CustomError(status, err.Error(), "failed input data")
	// 	}
	// }
	//save tr_transaction
	params["branch_key"] = "1"

	params["agent_key"] = agentkey
	params["trans_status_key"] = "2"
	params["trans_date"] = time.Now().Format(dateLayout)
	params["trans_type_key"] = "2"
	params["trx_code"] = "138"
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
	if metodePerhitungan == "3" { //amount
		params["entry_mode"] = "139"
	} else {
		params["entry_mode"] = "140"
	}
	params["trans_fee_percent"] = "0"
	params["trans_fee_amount"] = "0"
	params["charges_fee_amount"] = "0"
	params["services_fee_amount"] = "0"
	//params["aca_key"] = acaKey
	params["trans_source"] = "141"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strIDUserLogin
	params["risk_waiver"] = "0"

	var userData models.ScUserLogin
	status, err = models.GetScUserLoginByCustomerKey(&userData, customerKey)
	if err != nil {
		return lib.CustomError(status)
	}
	params["rec_attribute_id3"] = c.Request().UserAgent()

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
	lib.CreateNotifCustomerFromAdminByCustomerId(customerKey, subject, body, "TRANSACTION")

	//send email to BO role 11 & Sales
	SentEmailTransactionToBackOfficeAndSales(transactionID, "11")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func Switching(c echo.Context) error {
	var err error
	var status int
	decimal.MarshalJSONWithoutQuotes = true
	params := make(map[string]string)
	paramsSwIn := make(map[string]string)

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}

	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)

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
	status, err = models.CheckProductAllowRedmOrSwitching(&productNotAllow, customerKey, productIds)
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
		} else {
			log.Error("Wrong input for parameter: product_to")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: product_to", "Wrong input for parameter: product_to")
		}
	} else {
		log.Error("Missing required parameter: product_to")
		return lib.CustomError(http.StatusBadRequest, "Missing required parameter: product_to", "Missing required parameter: product_to")
	}

	status, err = models.GetMsProduct(&productTo, productToKeyStr)
	if err != nil {
		log.Error(err.Error())
		return lib.CustomError(http.StatusBadRequest, err.Error(), "Product Tujuan tidak ditemukan")
	}

	zero := decimal.NewFromInt(0)
	var balance models.SumBalanceUnit
	status, err = models.GetBalanceUnitByCustomerAndProduct(&balance, customerKey, productKeyStr)
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
	metodePerhitungan := c.FormValue("calc_method")
	if metodePerhitungan != "" {
		if metodePerhitungan == "all" { //all unit
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

			valueSwitchToAmount := value.Mul(navDB[0].NavValue).Truncate(0)
			if valueSwitchToAmount.Cmp(productTo.MinSubAmount) == -1 {
				log.Error("Min. Product Switch In Amount < Switching unit * Last NAB")
				return lib.CustomError(http.StatusBadRequest, "Min. Product Switch In Amount < Switching unit * Last NAB", "Min. Product Switch In Amount < Switching unit * Last NAB. Min SProduct Switch In Amount : "+productTo.MinSubAmount.String())
			}

			params["trans_unit"] = transUnitStr
			params["total_amount"] = "0"
		} else if metodePerhitungan == "unit" { //unit penyertaan
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

			if sisaUnitAfterRed != zero && sisaUnitAfterRed.Cmp(minSisa) == -1 {
				log.Error("Sisa unit setelah switching kurang dari minimal unit, Silakan switch All")
				return lib.CustomError(http.StatusBadRequest, "Sisa unit setelah switching kurang dari minimal unit, Silakan switching All. Sisa unit harus minimal : "+minSisa.String(), "Sisa unit setelah switching kurang dari minimal unit, Silakan switching All. Sisa unit harus minimal : "+minSisa.String())
			}

			valueSwitchToAmount := value.Mul(navDB[0].NavValue).Truncate(0)
			if valueSwitchToAmount.Cmp(productTo.MinSubAmount) == -1 {
				log.Error("Min. Product Switch In Amount < Switching unit * Last NAB")
				return lib.CustomError(http.StatusBadRequest, "Min. Product Switch In Amount < Switching unit * Last NAB", "Min. Product Switch In Amount < Switching unit * Last NAB. Min SProduct Switch In Amount : "+productTo.MinSubAmount.String())
			}

			params["trans_unit"] = transUnitStr
			params["total_amount"] = "0"
		} else if metodePerhitungan == "amount" { //Nominal
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

			if sisaUnitAfterRed != zero && sisaUnitAfterRed.Cmp(minSisa) == -1 {
				log.Error("Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All")
				return lib.CustomError(http.StatusBadRequest, "Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All. Sisa unit harus minimal : "+minSisa.String(), "Sisa unit setelah redemption kurang dari minimal unit, Silakan redemption All. Sisa unit harus minimal : "+minSisa.String())
			}

			valueSwitchToAmount := value.Truncate(0)
			if valueSwitchToAmount.Cmp(productTo.MinSubAmount) == -1 {
				log.Error("Min. Product Switch In Amount < Switching unit * Last NAB")
				return lib.CustomError(http.StatusBadRequest, "Min. Product Switch In Amount < Switching unit * Last NAB", "Min. Product Switch In Amount < Switching unit * Last NAB. Min SProduct Switch In Amount : "+productTo.MinSubAmount.String())
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
	paramsAcc["customer_key"] = customerKey
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
	var agentCustomerDB models.MsAgentCustomer
	status, err = models.GetLastAgenCunstomer(&agentCustomerDB, customerKey)
	if err != nil {
		log.Error("Failed get data agent: " + err.Error())
		return lib.CustomError(status, err.Error(), "failed input data")
	}
	agentkey := strconv.FormatUint(agentCustomerDB.AgentKey, 10)
	//cek tr_account_agent / save
	paramsAccAgent := make(map[string]string)
	paramsAccAgent["acc_key"] = accKey
	paramsAccAgent["agent_key"] = agentkey
	paramsAccAgent["rec_status"] = "1"

	// var acaKey string
	// var accountAgentDB []models.TrAccountAgent
	// status, err = models.GetAllTrAccountAgent(&accountAgentDB, paramsAccAgent)
	// if len(accountAgentDB) > 0 {
	// 	acaKey = strconv.FormatUint(accountAgentDB[0].AcaKey, 10)
	// } else {
	// 	paramsCreateAccAgent := make(map[string]string)
	// 	paramsCreateAccAgent["acc_key"] = accKey
	// 	paramsCreateAccAgent["eff_date"] = time.Now().Format(dateLayout)
	// 	paramsCreateAccAgent["rec_created_date"] = time.Now().Format(dateLayout)
	// 	paramsCreateAccAgent["rec_created_by"] = strIDUserLogin
	// 	paramsCreateAccAgent["agent_key"] = agentkey
	// 	paramsCreateAccAgent["branch_key"] = "1"
	// 	paramsCreateAccAgent["rec_status"] = "1"
	// 	status, err, acaKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
	// 	if err != nil {
	// 		log.Error("Failed create account agent data: " + err.Error())
	// 		return lib.CustomError(status, err.Error(), "failed input data")
	// 	}
	// }

	params["branch_key"] = "1"
	paramsSwIn["branch_key"] = "1"

	params["agent_key"] = agentkey
	params["trans_status_key"] = "2"
	params["trans_date"] = time.Now().Format(dateLayout)
	params["trans_type_key"] = "3"
	params["trx_code"] = "138"
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
	if metodePerhitungan == "3" { //Nominal
		params["entry_mode"] = "139"
	} else {
		params["entry_mode"] = "140"
	}
	params["trans_fee_percent"] = "0"
	params["trans_fee_amount"] = "0"
	params["charges_fee_amount"] = "0"
	params["services_fee_amount"] = "0"

	//params["aca_key"] = acaKey
	params["trans_source"] = "141"
	params["rec_status"] = "1"
	params["rec_created_date"] = time.Now().Format(dateLayout)
	params["rec_created_by"] = strIDUserLogin
	params["risk_waiver"] = "0"
	params["flag_newsub"] = "0"

	var userData models.ScUserLogin
	status, err = models.GetScUserLoginByCustomerKey(&userData, customerKey)
	if err != nil {
		return lib.CustomError(status)
	}
	params["rec_attribute_id3"] = c.Request().UserAgent()

	status, err, transactionID := models.CreateTrTransaction(params)

	//SAVE PRODUCT TO
	//cek tr_account / save
	var accNewKey string
	paramsNewProdAcc := make(map[string]string)
	paramsNewProdAcc["customer_key"] = customerKey
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
	paramsNewProdAccAgent["agent_key"] = agentkey
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
		paramsCreateAccAgent["agent_key"] = agentkey
		paramsCreateAccAgent["branch_key"] = "1"
		paramsCreateAccAgent["rec_status"] = "1"
		status, err, acaNewProdKey = models.CreateTrAccountAgent(paramsCreateAccAgent)
		if err != nil {
			log.Error("Failed create account agent data: " + err.Error())
			return lib.CustomError(status, err.Error(), "failed input data")
		}
	}

	paramsSwIn["agent_key"] = agentkey
	paramsSwIn["trans_status_key"] = "2"
	paramsSwIn["trans_date"] = time.Now().Format(dateLayout)
	paramsSwIn["trans_type_key"] = "4"
	paramsSwIn["trx_code"] = "137"
	paramsSwIn["payment_method"] = "284"
	paramsSwIn["trans_calc_method"] = "288"

	paramsSwIn["entry_mode"] = "140"
	paramsSwIn["trans_fee_percent"] = "0"
	paramsSwIn["trans_fee_amount"] = "0"
	paramsSwIn["charges_fee_amount"] = "0"
	paramsSwIn["services_fee_amount"] = "0"

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
	paramsSwIn["rec_attribute_id3"] = c.Request().UserAgent()

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
	lib.CreateNotifCustomerFromAdminByCustomerId(customerKey, subject, body, "TRANSACTION")

	//send email to role 11 & sales
	SentEmailTransactionToBackOfficeAndSales(transactionID, "11")

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = ""
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

	//sent email to BO role 11 & Sales
	SentEmailTransactionToBackOfficeAndSales(transactionKeyStr, "11")

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
	decimal.MarshalJSONWithoutQuotes = true
	zero := decimal.NewFromInt(0)

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
	var transKeyAll []string
	var statusIDs []string
	var typeIDs []string
	var transKeyParent []string
	var transSwInParentKey []string
	var bankIDs []string
	var navDates []string
	switchout := make(map[uint64]models.TrTransaction)
	for _, transaction := range transactionDB {
		if transaction.TransTypeKey == 3 {
			switchout[transaction.TransactionKey] = transaction
			transSwInParentKey = append(transSwInParentKey, strconv.FormatUint(transaction.TransactionKey, 10))
		}

		if transaction.TransTypeKey == 4 {
			if transaction.ParentKey != nil {
				if _, ok := lib.Find(transKeyParent, strconv.FormatUint(*transaction.ParentKey, 10)); !ok {
					transKeyParent = append(transKeyParent, strconv.FormatUint(*transaction.ParentKey, 10))
				}
			}
		}

		if _, ok := lib.Find(transKeyAll, strconv.FormatUint(transaction.TransactionKey, 10)); !ok {
			transKeyAll = append(transKeyAll, strconv.FormatUint(transaction.TransactionKey, 10))
		}
		productIDs = append(productIDs, strconv.FormatUint(transaction.ProductKey, 10))
		statusIDs = append(statusIDs, strconv.FormatUint(transaction.TransStatusKey, 10))
		typeIDs = append(typeIDs, strconv.FormatUint(transaction.TransTypeKey, 10))
		if transaction.TransBankKey != nil {
			bankIDs = append(bankIDs, strconv.FormatUint(*transaction.TransBankKey, 10))
		}
		navDates = append(navDates, "'"+transaction.NavDate+"'")
	}

	var transactionParentList []models.TrTransaction
	if len(transKeyParent) > 0 {
		_, err := models.GetTrTransactionIn(&transactionParentList, transKeyParent, "transaction_key")
		if err == nil {
			for _, tr := range transactionParentList {
				if tr.TransTypeKey == 3 {
					switchout[tr.TransactionKey] = tr
				}
				productIDs = append(productIDs, strconv.FormatUint(tr.ProductKey, 10))
			}
		}
	}

	var transactionSwInList []models.TrTransaction
	parentTrans := make(map[uint64]models.TrTransaction)
	if len(transSwInParentKey) > 0 {
		_, err := models.GetTrTransactionIn(&transactionSwInList, transSwInParentKey, "parent_key")
		if err == nil {
			for _, trswin := range transactionSwInList {
				if trswin.ParentKey != nil {
					parentTrans[*trswin.ParentKey] = trswin
				}
				if trswin.TransTypeKey == 4 {
					if _, ok := lib.Find(transKeyAll, strconv.FormatUint(trswin.TransactionKey, 10)); !ok {
						transactionDB = append(transactionDB, trswin)
					}
				}
				productIDs = append(productIDs, strconv.FormatUint(trswin.ProductKey, 10))
			}
		}
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

	var responseData []models.TrTransactionList
	for _, transaction := range transactionDB {
		if trStatus == "process" {
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
					data.TransType = *typ.TypeDescription
				}

				strProductKey := strconv.FormatUint(transaction.ProductKey, 10)
				var navTrans models.TrNav
				status, err = models.GetNavByProductKeyAndNavDate(&navTrans, strProductKey, transaction.NavDate)
				if err == nil {
					data.NavValue = navTrans.NavValue
				} else {
					data.NavValue = zero
				}
				data.TransDate = transaction.TransDate
				data.NavDate = transaction.NavDate

				//cek transaction confirmation
				var transactionConf models.TrTransactionConfirmation
				strTrKey := strconv.FormatUint(transaction.TransactionKey, 10)
				_, err = models.GetTrTransactionConfirmationByTransactionKey(&transactionConf, strTrKey)
				if err != nil {
					if transaction.TransTypeKey == 4 {
						if transaction.ParentKey != nil {
							if swot, ok := switchout[*transaction.ParentKey]; ok {
								log.Println("HAHAHHAA")
								data.TransUnit = swot.TransUnit.Truncate(2)
								data.TransAmount = swot.TransAmount.Truncate(0)
							}
						} else {
							log.Println("qqqqqqqqqqqq")
							data.TransUnit = transactionConf.ConfirmedUnit.Truncate(2)
							data.TransAmount = transactionConf.ConfirmedAmount.Truncate(0)
						}
					} else {
						log.Println("rrrrrrrrrr")
						data.TransUnit = transaction.TransUnit.Truncate(2)
						data.TransAmount = transaction.TransAmount.Truncate(0)
					}
				} else {
					data.TransUnit = transactionConf.ConfirmedUnit.Truncate(2)
					data.TransAmount = transactionConf.ConfirmedAmount.Truncate(0)
				}

				data.TotalAmount = transaction.TotalAmount
				if transaction.FileUploadDate != nil {
					data.Uploaded = true
					data.DateUploaded = transaction.FileUploadDate
				} else {
					if transaction.TransTypeKey == 1 {
						data.Uploaded = false
					} else {
						data.Uploaded = true
					}
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
						data.ProductIn = &product.ProductNameAlt
					}
					if transaction.ParentKey != nil {
						if swot, ok := switchout[*transaction.ParentKey]; ok {
							if product, ok := pData[swot.ProductKey]; ok {
								data.ProductOut = &product.ProductNameAlt
							}
						}
					}
				}

				responseData = append(responseData, data)
			}
		} else {
			var data models.TrTransactionList

			data.TransactionKey = transaction.TransactionKey
			if product, ok := pData[transaction.ProductKey]; ok {
				data.ProductName = product.ProductNameAlt
			}
			if status, ok := sData[transaction.TransStatusKey]; ok {
				data.TransStatus = *status.StatusCode
			}
			if typ, ok := tData[transaction.TransTypeKey]; ok {
				data.TransType = *typ.TypeDescription
			}

			strProductKey := strconv.FormatUint(transaction.ProductKey, 10)
			var navTrans models.TrNav
			status, err = models.GetNavByProductKeyAndNavDate(&navTrans, strProductKey, transaction.NavDate)
			if err == nil {
				data.NavValue = navTrans.NavValue
			} else {
				data.NavValue = zero
			}
			data.TransDate = transaction.TransDate
			data.NavDate = transaction.NavDate

			//cek transaction confirmation
			var transactionConf models.TrTransactionConfirmation
			strTrKey := strconv.FormatUint(transaction.TransactionKey, 10)
			_, err = models.GetTrTransactionConfirmationByTransactionKey(&transactionConf, strTrKey)
			if err != nil {
				data.TransUnit = transaction.TransUnit.Truncate(2)
				data.TransAmount = transaction.TransAmount.Truncate(0)
			} else {
				data.TransUnit = transactionConf.ConfirmedUnit.Truncate(2)
				data.TransAmount = transactionConf.ConfirmedAmount.Truncate(0)
			}

			data.TotalAmount = transaction.TotalAmount
			if transaction.FileUploadDate != nil {
				data.Uploaded = true
				data.DateUploaded = transaction.FileUploadDate
			} else {
				if transaction.TransTypeKey == 1 {
					data.Uploaded = false
				} else {
					data.Uploaded = true
				}
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
					data.ProductIn = &product.ProductNameAlt
				}
				if transaction.ParentKey != nil {
					if swot, ok := switchout[*transaction.ParentKey]; ok {
						if product, ok := pData[swot.ProductKey]; ok {
							data.ProductOut = &product.ProductNameAlt
						}
					}
				}
			}
			if transaction.TransTypeKey == 3 {
				if product, ok := pData[transaction.ProductKey]; ok {
					data.ProductOut = &product.ProductNameAlt
				}
				if trParent, ok := parentTrans[transaction.TransactionKey]; ok {
					if product, ok := pData[trParent.ProductKey]; ok {
						data.ProductIn = &product.ProductNameAlt
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

func ValidatePromoTransaction(c echo.Context) error {

	if lib.Profile.CustomerKey == nil || *lib.Profile.CustomerKey == 0 {
		log.Error("No customer found")
		return lib.CustomError(http.StatusBadRequest, "No customer found", "No customer found, please open account first")
	}

	customerKey := strconv.FormatUint(*lib.Profile.CustomerKey, 10)

	productKeyStr := c.FormValue("product_key")
	if productKeyStr != "" {
		productKey, err := strconv.ParseUint(productKeyStr, 10, 64)
		if err == nil && productKey > 0 {
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

	responseData := make(map[string]string)
	promoCode := c.FormValue("promo_code")
	if promoCode != "" {
		err, enable, text, _ := validatePromo(promoCode, customerKey, productKeyStr)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest, err.Error(), "Failed get data")
		} else {
			if enable == false {
				return lib.CustomError(http.StatusBadRequest, text, text)
			}
		}
		responseData["message"] = text
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
	if err := t.Execute(&tpl, struct{ FileUrl string }{FileUrl: config.FileUrl + "/images/mail"}); err != nil {
		log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", lib.Profile.Email)
	mailer.SetHeader("Subject", "[MotionFunds] Laporan Akun")
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
	decimal.MarshalJSONWithoutQuotes = true
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
			fee := transaction.TransFeeAmount
			data["Fee"] = fee
			if tc, ok := tcData[transaction.TransactionKey]; ok {
				data["Unit"] = tc.ConfirmedUnit.Truncate(2)
				data["Amount"] = tc.ConfirmedAmount.Truncate(0)
				data["Total"] = tc.ConfirmedAmount.Add(fee).Truncate(0)
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
	if err := t.Execute(&tpl, struct{ FileUrl string }{FileUrl: config.FileUrl + "/images/mail"}); err != nil {
		log.Println(err)
	}

	result := tpl.String()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", lib.Profile.Email)
	mailer.SetHeader("Subject", "[MotionFunds] Histori Transaksi")
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

func mailTransaction(typ string, params map[string]string) error {
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
		if _, ok := params["rec_image1"]; ok {
			mailTemp = "index-" + typ + "-complete.html"
			s := "Subscription"
			if typ == "topup" {
				s = "Top Up"
			}
			subject = s + " Kamu sedang Diproses"
		} else {
			mailTemp = "index-" + typ + "-uncomplete.html"
			subject = "Ayo Upload Bukti Transfer Kamu"
		}
	} else if typ == "redemption" {
		mailTemp = "index-" + typ + ".html"
		subject = "Redemption Kamu sedang Diproses"
		// var requestDB []models.OaRequest
		// paramRequest := make(map[string]string)
		// paramRequest["user_login_key"] = strconv.FormatUint(lib.Profile.UserID, 10)
		// paramRequest["orderBy"] = "oa_request_key"
		// paramRequest["orderType"] = "DESC"
		// _, err = models.GetAllOaRequest(&requestDB, 1, 0, false, paramRequest)
		// if err != nil {
		// 	log.Error("Failed send mail: " + err.Error())
		// 	return err
		// }
		// if len(requestDB) < 1 {
		// 	log.Error("Failed send mail: no bank data")
		// 	return nil
		// }
		// var personalData models.OaPersonalData
		// requestKey := strconv.FormatUint(requestDB[0].OaRequestKey, 10)
		// _, err = models.GetOaPersonalData(&personalData, requestKey, "oa_request_key")
		// if err != nil {
		// 	log.Error("Failed send mail: " + err.Error())
		// 	return err
		// }
		var cusBankAcc models.MsCustomerBankAccount
		_, err = models.GetMsCustomerBankAccount(&cusBankAcc, params["customer_bank_account_key"])
		if err != nil {
			log.Error("Failed send mail: " + err.Error())
			return err
		}
		var bankAccount models.MsBankAccount
		bankAccountKey := strconv.FormatUint(cusBankAcc.BankAccountKey, 10)
		_, err = models.GetBankAccount(&bankAccount, bankAccountKey)
		if err != nil {
			log.Error("Failed send mail: " + err.Error())
			return err
		}
		var bank models.MsBank
		bankKey := strconv.FormatUint(bankAccount.BankKey, 10)
		_, err = models.GetMsBank(&bank, bankKey)
		if err != nil {
			log.Error("Failed send mail: " + err.Error())
			return err
		}

		mailParam["BankName"] = bank.BankName
		mailParam["BankAccNo"] = bankAccount.AccountNo
		mailParam["AccHolderName"] = bankAccount.AccountHolderName
		mailParam["Branch"] = *bankAccount.BranchName

	} else if typ == "switching" {
		mailTemp = "index-" + typ + ".html"
		subject = "Switching Kamu sedang Diproses"

		var transaction models.TrTransaction
		_, err = models.GetTrTransaction(&transaction, params["parrent"])
		if err != nil {
			log.Error("Failed send mail: " + err.Error())
			return err
		}

		var product models.MsProduct
		_, err = models.GetMsProduct(&product, strconv.FormatUint(transaction.ProductKey, 10))
		if err != nil {
			log.Error("Failed send mail: " + err.Error())
			return err
		}
		mailParam["ProductOut"] = product.ProductNameAlt
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
	_, err = models.GetMsCustomer(&customer, strconv.FormatUint(*lib.Profile.CustomerKey, 10))
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
	mailer.SetHeader("To", lib.Profile.Email)
	mailer.SetHeader("Subject", "[MotionFunds] "+subject)
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
func mailSubscription(typ string, params map[string]string) error {
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
	if _, ok := params["rec_image1"]; ok {
		mailTemp = "index-" + typ + "-complete.html"
		s := "Subscription"
		if typ == "topup" {
			s = "Top Up"
		}
		subject = s + " Kamu sedang Diproses"
	} else {
		mailTemp = "index-" + typ + "-uncomplete.html"
		subject = "Ayo Upload Bukti Transfer Kamu"
	}
	val, _ = decimal.NewFromString(params["trans_amount"])
	mailParam["Amount"] = ac0.FormatMoneyDecimal(val)
	var customer models.MsCustomer
	_, err = models.GetMsCustomer(&customer, strconv.FormatUint(*lib.Profile.CustomerKey, 10))
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
	mailer.SetHeader("To", lib.Profile.Email)
	mailer.SetHeader("Subject", "[MotionFunds] "+subject)
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
