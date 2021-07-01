package controllers

import (
	"api/lib"
	"api/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func AdminDetailAccountStatementCustomerProduct(c echo.Context) error {
	var err error
	decimal.MarshalJSONWithoutQuotes = true

	customerKey := c.QueryParam("customer_key")
	if customerKey == "" {
		log.Error("Missing required parameter: menu_key")
		return lib.CustomError(http.StatusBadRequest, "menu_key can not be blank", "menu_key can not be blank")
	} else {
		n, err := strconv.ParseUint(customerKey, 10, 64)
		if err != nil || n == 0 {
			log.Error("Wrong input for parameter: menu_key")
			return lib.CustomError(http.StatusBadRequest, "Wrong input for parameter: menu_key", "Wrong input for parameter: menu_key")
		}
	}

	var customer models.HeaderCustomerDetailAccountStatement
	_, err = models.GetHeaderCustomerDetailAccountStatement(&customer, customerKey)
	if err != nil {
		log.Error("Customer not found")
		return lib.CustomError(http.StatusBadRequest, "Customer not found", "Customer not found")
	}

	var datefrom string
	var dateto string

	datefrom = c.QueryParam("date_from")
	if datefrom == "" {
		log.Error("Missing required parameter: date_from")
		return lib.CustomError(http.StatusBadRequest, "date_from can not be blank", "date_from can not be blank")
	}

	dateto = c.QueryParam("date_to")
	if dateto == "" {
		log.Error("Missing required parameter: date_to")
		return lib.CustomError(http.StatusBadRequest, "date_to can not be blank", "date_to can not be blank")
	}
	layoutISO := "2006-01-02"

	from, _ := time.Parse(layoutISO, datefrom)
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)

	to, _ := time.Parse(layoutISO, dateto)
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)

	var dateawal string
	var dateakhir string

	if from.Before(to) {
		dateawal = datefrom
		dateakhir = dateto
	}

	if from.After(to) {
		dateawal = dateto
		dateakhir = datefrom
	}

	responseData := make(map[string]interface{})

	//get header
	layout := "2006-01-02"
	newLayout := "02 Jan 2006"
	header := make(map[string]interface{})
	dateParem, _ := time.Parse(layout, dateawal)
	header["date_from"] = dateParem.Format(newLayout)
	dateParem, _ = time.Parse(layout, dateakhir)
	header["date_to"] = dateParem.Format(newLayout)
	header["customer_key"] = customer.CustomerKey
	header["cif"] = customer.Cif
	header["sid"] = customer.Sid
	header["full_name"] = customer.FullName
	header["address"] = customer.Address

	responseData["header"] = header

	var transactions []models.AccountStatementCustomerProduct

	status, err := models.AdminGetAllAccountStatementCustomerProduct(&transactions, customerKey, dateawal, dateakhir)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data transaction")
		}
	}

	if len(transactions) > 0 {
		var datatrans []interface{}
		var productKey uint64
		var transGroupProduct []interface{}
		product := make(map[string]interface{})
		count := make(map[string]interface{})
		var totalSubs decimal.Decimal
		var totalRedm decimal.Decimal
		nol := decimal.NewFromInt(0)

		var accKeyLast uint64
		var productKeyLast uint64

		for idx, tr := range transactions {
			if idx != 0 {
				if productKey != tr.ProductKey {
					var balanceEnding models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalance(&balanceEnding, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

					endingbalance := make(map[string]interface{})
					if err != nil {
						dateParem, _ = time.Parse(layout, dateakhir)
						endingbalance["date"] = dateParem.Format(newLayout)
						endingbalance["description"] = "ENDING BALANCE"
						endingbalance["amount"] = nol
						endingbalance["nav_value"] = tr.NavValue.Truncate(2)
						endingbalance["unit"] = nol
						endingbalance["avg_nav"] = tr.AvgNav.Truncate(2)
						endingbalance["fee"] = nol
						transGroupProduct = append(transGroupProduct, endingbalance)
					} else {
						endingbalance["date"] = balanceEnding.Tanggal
						endingbalance["description"] = balanceEnding.Description
						endingbalance["amount"] = balanceEnding.Amount.Truncate(0)
						endingbalance["nav_value"] = balanceEnding.NavValue.Truncate(2)
						endingbalance["unit"] = balanceEnding.Unit.Truncate(2)
						endingbalance["avg_nav"] = balanceEnding.AvgNav.Truncate(2)
						endingbalance["fee"] = balanceEnding.Fee.Truncate(0)
						transGroupProduct = append(transGroupProduct, endingbalance)
					}

					row := make(map[string]interface{})
					row["count"] = count
					row["product"] = product
					row["list"] = transGroupProduct
					datatrans = append(datatrans, row)
					productKey = tr.ProductKey

					transGroupProduct = nil
					product = make(map[string]interface{})
					count = make(map[string]interface{})
					totalSubs = nol
					totalRedm = nol

					var balance models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalance(&balance, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

					beginning := make(map[string]interface{})
					if err != nil {
						dateParem, _ = time.Parse(layout, dateawal)
						beginning["date"] = dateParem.Format(newLayout)
						beginning["description"] = "BEGINNING BALANCE"
						beginning["amount"] = nol
						beginning["nav_value"] = tr.NavValue.Truncate(2)
						beginning["unit"] = nol
						beginning["avg_nav"] = tr.AvgNav.Truncate(2)
						beginning["fee"] = nol
						transGroupProduct = append(transGroupProduct, beginning)
					} else {
						beginning["date"] = balance.Tanggal
						beginning["description"] = balance.Description
						beginning["amount"] = balance.Amount.Truncate(0)
						beginning["nav_value"] = balance.NavValue.Truncate(2)
						beginning["unit"] = balance.Unit.Truncate(2)
						beginning["avg_nav"] = balance.AvgNav.Truncate(2)
						beginning["fee"] = balance.Fee.Truncate(0)
						transGroupProduct = append(transGroupProduct, beginning)
					}

					trans := make(map[string]interface{})
					trans["date"] = tr.NavDate
					trans["description"] = tr.Trans
					trans["amount"] = tr.Amount.Truncate(0)
					trans["nav_value"] = tr.NavValue.Truncate(2)
					trans["unit"] = tr.Unit.Truncate(2)
					trans["avg_nav"] = tr.AvgNav.Truncate(2)
					trans["fee"] = tr.Fee.Truncate(0)
					transGroupProduct = append(transGroupProduct, trans)

					if (tr.TransTypeKey == uint64(1)) || (tr.TransTypeKey == uint64(4)) {
						totalSubs = totalSubs.Add(tr.Amount).Truncate(0)
					}
					if (tr.TransTypeKey == uint64(2)) || (tr.TransTypeKey == uint64(3)) {
						totalRedm = totalRedm.Add(tr.Amount).Truncate(0)
					}
					count["subs"] = totalSubs
					count["redm"] = totalRedm

					product["product_id"] = tr.ProductKey
					product["product_name"] = tr.ProductName
					product["product_bank_name"] = tr.BankName
					product["product_bank_account_name"] = tr.AccountName
					product["product_bank_account_no"] = tr.AccountNo
					product["product_code"] = tr.ProductCode
					product["currency"] = tr.Currency

					if idx == (len(transactions) - 1) {
						var balanceEndingLast models.BeginningEndingBalance
						_, err = models.GetBeginningEndingBalance(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

						endingbalancelast := make(map[string]interface{})
						if err != nil {
							dateParem, _ = time.Parse(layout, dateakhir)
							endingbalancelast["date"] = dateParem.Format(newLayout)
							endingbalancelast["description"] = "ENDING BALANCE"
							endingbalancelast["amount"] = nol
							endingbalancelast["nav_value"] = tr.NavValue.Truncate(2)
							endingbalancelast["unit"] = nol
							endingbalancelast["avg_nav"] = tr.AvgNav.Truncate(2)
							endingbalancelast["fee"] = nol
							transGroupProduct = append(transGroupProduct, endingbalancelast)
						} else {
							endingbalancelast["date"] = balanceEndingLast.Tanggal
							endingbalancelast["description"] = balanceEndingLast.Description
							endingbalancelast["amount"] = balanceEndingLast.Amount.Truncate(0)
							endingbalancelast["nav_value"] = balanceEndingLast.NavValue.Truncate(2)
							endingbalancelast["unit"] = balanceEndingLast.Unit.Truncate(2)
							endingbalancelast["avg_nav"] = balanceEndingLast.AvgNav.Truncate(2)
							endingbalancelast["fee"] = balanceEndingLast.Fee.Truncate(0)
							transGroupProduct = append(transGroupProduct, endingbalancelast)
						}

						row := make(map[string]interface{})
						row["count"] = count
						row["product"] = product
						row["list"] = transGroupProduct
						datatrans = append(datatrans, row)
					}
				} else {
					trans := make(map[string]interface{})
					trans["date"] = tr.NavDate
					trans["description"] = tr.Trans
					trans["amount"] = tr.Amount.Truncate(0)
					trans["nav_value"] = tr.NavValue.Truncate(2)
					trans["unit"] = tr.Unit.Truncate(2)
					trans["avg_nav"] = tr.AvgNav.Truncate(2)
					trans["fee"] = tr.Fee.Truncate(0)
					transGroupProduct = append(transGroupProduct, trans)

					if (tr.TransTypeKey == 1) || (tr.TransTypeKey == 4) {
						totalSubs = totalSubs.Add(tr.Amount).Truncate(0)
					}
					if (tr.TransTypeKey == 2) || (tr.TransTypeKey == 3) {
						totalRedm = totalRedm.Add(tr.Amount).Truncate(0)
					}
					count["subs"] = totalSubs
					count["redm"] = totalRedm

					product["product_id"] = tr.ProductKey
					product["product_name"] = tr.ProductName
					product["product_bank_name"] = tr.BankName
					product["product_bank_account_name"] = tr.AccountName
					product["product_bank_account_no"] = tr.AccountNo
					product["product_code"] = tr.ProductCode
					product["currency"] = tr.Currency

					if idx == (len(transactions) - 1) {
						var balanceEndingLast models.BeginningEndingBalance
						_, err = models.GetBeginningEndingBalance(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

						endingbalancelast := make(map[string]interface{})
						if err != nil {
							dateParem, _ = time.Parse(layout, dateakhir)
							endingbalancelast["date"] = dateParem.Format(newLayout)
							endingbalancelast["description"] = "ENDING BALANCE"
							endingbalancelast["amount"] = nol
							endingbalancelast["nav_value"] = tr.NavValue.Truncate(2)
							endingbalancelast["unit"] = nol
							endingbalancelast["avg_nav"] = tr.AvgNav.Truncate(2)
							endingbalancelast["fee"] = nol
							transGroupProduct = append(transGroupProduct, endingbalancelast)
						} else {
							endingbalancelast["date"] = balanceEndingLast.Tanggal
							endingbalancelast["description"] = balanceEndingLast.Description
							endingbalancelast["amount"] = balanceEndingLast.Amount.Truncate(0)
							endingbalancelast["nav_value"] = balanceEndingLast.NavValue.Truncate(2)
							endingbalancelast["unit"] = balanceEndingLast.Unit.Truncate(2)
							endingbalancelast["avg_nav"] = balanceEndingLast.AvgNav.Truncate(2)
							endingbalancelast["fee"] = balanceEndingLast.Fee.Truncate(0)
							transGroupProduct = append(transGroupProduct, endingbalancelast)
						}

						row := make(map[string]interface{})
						row["count"] = count
						row["product"] = product
						row["list"] = transGroupProduct
						datatrans = append(datatrans, row)
					}

				}
			} else {

				var balance models.BeginningEndingBalance
				_, err = models.GetBeginningEndingBalance(&balance, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

				beginning := make(map[string]interface{})
				if err != nil {
					dateParem, _ = time.Parse(layout, dateawal)
					beginning["date"] = dateParem.Format(newLayout)
					beginning["description"] = "BEGINNING BALANCE"
					beginning["amount"] = nol
					beginning["nav_value"] = tr.NavValue.Truncate(2)
					beginning["unit"] = nol
					beginning["avg_nav"] = tr.AvgNav.Truncate(2)
					beginning["fee"] = nol
					transGroupProduct = append(transGroupProduct, beginning)
				} else {
					beginning["date"] = balance.Tanggal
					beginning["description"] = balance.Description
					beginning["amount"] = balance.Amount.Truncate(0)
					beginning["nav_value"] = balance.NavValue.Truncate(2)
					beginning["unit"] = balance.Unit.Truncate(2)
					beginning["avg_nav"] = balance.AvgNav.Truncate(2)
					beginning["fee"] = balance.Fee.Truncate(0)
					transGroupProduct = append(transGroupProduct, beginning)
				}

				trans := make(map[string]interface{})
				trans["date"] = tr.NavDate
				trans["description"] = tr.Trans
				trans["amount"] = tr.Amount.Truncate(0)
				trans["nav_value"] = tr.NavValue.Truncate(2)
				trans["unit"] = tr.Unit.Truncate(2)
				trans["avg_nav"] = tr.AvgNav.Truncate(2)
				trans["fee"] = tr.Fee.Truncate(0)
				transGroupProduct = append(transGroupProduct, trans)

				if (tr.TransTypeKey == 1) || (tr.TransTypeKey == 4) {
					totalSubs = totalSubs.Add(tr.Amount).Truncate(0)
				}
				if (tr.TransTypeKey == 2) || (tr.TransTypeKey == 3) {
					totalRedm = totalRedm.Add(tr.Amount).Truncate(0)
				}
				count["subs"] = totalSubs
				count["redm"] = totalRedm

				product["product_id"] = tr.ProductKey
				product["product_name"] = tr.ProductName
				product["product_bank_name"] = tr.BankName
				product["product_bank_account_name"] = tr.AccountName
				product["product_bank_account_no"] = tr.AccountNo
				product["product_code"] = tr.ProductCode
				product["currency"] = tr.Currency

				if idx == (len(transactions) - 1) {
					var balanceEndingLast models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalance(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

					endingbalancelast := make(map[string]interface{})
					if err != nil {
						dateParem, _ = time.Parse(layout, dateakhir)
						endingbalancelast["date"] = dateParem.Format(newLayout)
						endingbalancelast["description"] = "ENDING BALANCE"
						endingbalancelast["amount"] = nol
						endingbalancelast["nav_value"] = tr.NavValue.Truncate(2)
						endingbalancelast["unit"] = nol
						endingbalancelast["avg_nav"] = tr.AvgNav.Truncate(2)
						endingbalancelast["fee"] = nol
						transGroupProduct = append(transGroupProduct, endingbalancelast)
					} else {
						endingbalancelast["date"] = balanceEndingLast.Tanggal
						endingbalancelast["description"] = balanceEndingLast.Description
						endingbalancelast["amount"] = balanceEndingLast.Amount.Truncate(0)
						endingbalancelast["nav_value"] = balanceEndingLast.NavValue.Truncate(2)
						endingbalancelast["unit"] = balanceEndingLast.Unit.Truncate(2)
						endingbalancelast["avg_nav"] = balanceEndingLast.AvgNav.Truncate(2)
						endingbalancelast["fee"] = balanceEndingLast.Fee.Truncate(0)
						transGroupProduct = append(transGroupProduct, endingbalancelast)
					}

					row := make(map[string]interface{})
					row["count"] = count
					row["product"] = product
					row["list"] = transGroupProduct
					datatrans = append(datatrans, row)
				}
			}
			accKeyLast = tr.AccKey
			productKeyLast = tr.ProductKey
		}

		responseData["transaction"] = datatrans

	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = responseData

	return c.JSON(http.StatusOK, response)
}
