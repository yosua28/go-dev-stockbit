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
					_, err = models.GetBeginningEndingBalanceAcc(&balanceEnding, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

					endingbalance := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(productKeyLast, 10), customerKey, dateakhir)
						if err != nil {
							endingbalance["avg_nav"] = nol
						} else {
							endingbalance["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(productKeyLast, 10), dateakhir)
						if err != nil {
							endingbalance["nav_value"] = nol
						} else {
							endingbalance["nav_value"] = lastnav.NavValue.Truncate(2)
						}

						dateParem, _ = time.Parse(layout, dateakhir)
						endingbalance["date"] = dateParem.Format(newLayout)
						endingbalance["description"] = "ENDING BALANCE"
						endingbalance["amount"] = nol
						endingbalance["unit"] = nol
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
					_, err = models.GetBeginningEndingBalanceAcc(&balance, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

					beginning := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateawal)
						if err != nil {
							beginning["avg_nav"] = nol
						} else {
							beginning["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateawal)
						if err != nil {
							beginning["nav_value"] = nol
						} else {
							beginning["nav_value"] = lastnav.NavValue.Truncate(2)
						}
						dateParem, _ = time.Parse(layout, dateawal)
						beginning["date"] = dateParem.Format(newLayout)
						beginning["description"] = "BEGINNING BALANCE"
						beginning["amount"] = nol
						beginning["unit"] = nol
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
						_, err = models.GetBeginningEndingBalanceAcc(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

						endingbalancelast := make(map[string]interface{})
						if err != nil {

							var lastavgnav models.NavValue
							_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(productKeyLast, 10), customerKey, dateakhir)
							if err != nil {
								endingbalancelast["avg_nav"] = nol
							} else {
								endingbalancelast["avg_nav"] = lastavgnav.NavValue.Truncate(2)
							}

							var lastnav models.NavValue
							_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(productKeyLast, 10), dateakhir)
							if err != nil {
								endingbalancelast["nav_value"] = nol
							} else {
								endingbalancelast["nav_value"] = lastnav.NavValue.Truncate(2)
							}
							dateParem, _ = time.Parse(layout, dateakhir)
							endingbalancelast["date"] = dateParem.Format(newLayout)
							endingbalancelast["description"] = "ENDING BALANCE"
							endingbalancelast["amount"] = nol
							endingbalancelast["unit"] = nol
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
						_, err = models.GetBeginningEndingBalanceAcc(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(accKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

						endingbalancelast := make(map[string]interface{})
						if err != nil {
							var lastavgnav models.NavValue
							_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(productKeyLast, 10), customerKey, dateakhir)
							if err != nil {
								endingbalancelast["avg_nav"] = nol
							} else {
								endingbalancelast["avg_nav"] = lastavgnav.NavValue.Truncate(2)
							}

							var lastnav models.NavValue
							_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(productKeyLast, 10), dateakhir)
							if err != nil {
								endingbalancelast["nav_value"] = nol
							} else {
								endingbalancelast["nav_value"] = lastnav.NavValue.Truncate(2)
							}
							dateParem, _ = time.Parse(layout, dateakhir)
							endingbalancelast["date"] = dateParem.Format(newLayout)
							endingbalancelast["description"] = "ENDING BALANCE"
							endingbalancelast["amount"] = nol
							endingbalancelast["unit"] = nol
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
				_, err = models.GetBeginningEndingBalanceAcc(&balance, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

				beginning := make(map[string]interface{})
				if err != nil {
					var lastavgnav models.NavValue
					_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateawal)
					if err != nil {
						beginning["avg_nav"] = nol
					} else {
						beginning["avg_nav"] = lastavgnav.NavValue.Truncate(2)
					}

					var lastnav models.NavValue
					_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateawal)
					if err != nil {
						beginning["nav_value"] = nol
					} else {
						beginning["nav_value"] = lastnav.NavValue.Truncate(2)
					}
					dateParem, _ = time.Parse(layout, dateawal)
					beginning["date"] = dateParem.Format(newLayout)
					beginning["description"] = "BEGINNING BALANCE"
					beginning["amount"] = nol
					beginning["unit"] = nol
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
					_, err = models.GetBeginningEndingBalanceAcc(&balanceEndingLast, "ENDING BALANCE", dateakhir, strconv.FormatUint(tr.AccKey, 10), strconv.FormatUint(tr.ProductKey, 10))

					endingbalancelast := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateakhir)
						if err != nil {
							endingbalancelast["avg_nav"] = nol
						} else {
							endingbalancelast["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateakhir)
						if err != nil {
							endingbalancelast["nav_value"] = nol
						} else {
							endingbalancelast["nav_value"] = lastnav.NavValue.Truncate(2)
						}
						dateParem, _ = time.Parse(layout, dateakhir)
						endingbalancelast["date"] = dateParem.Format(newLayout)
						endingbalancelast["description"] = "ENDING BALANCE"
						endingbalancelast["amount"] = nol
						endingbalancelast["unit"] = nol
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

func AdminDetailAccountStatementCustomerAgent(c echo.Context) error {
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

	var transactions []models.AccountStatementCustomerAgent

	status, err := models.AdminGetAllAccountStatementCustomerAgent(&transactions, customerKey, dateawal, dateakhir)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err.Error())
			return lib.CustomError(status, err.Error(), "Failed get data transaction")
		}
	}

	if len(transactions) > 0 {
		var datatrans []interface{}
		nol := decimal.NewFromInt(0)
		var productKey uint64
		var salesKey uint64
		transGroupAgent := make(map[string]interface{})
		var transGroupProduct []interface{}
		var transGroupSales []interface{}
		product := make(map[string]interface{})
		sales := make(map[string]interface{})
		count := make(map[string]interface{})
		// product := make(map[string]interface{})
		// count := make(map[string]interface{})
		var totalSubs decimal.Decimal
		var totalRedm decimal.Decimal
		var totalNettsubs decimal.Decimal

		var acaKeyLast uint64
		var productKeyLast uint64

		for idx, tr := range transactions {
			if idx != 0 {
				if productKey != tr.ProductKey {
					//set ending
					var balanceEnding models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalanceAca(&balanceEnding, "ENDING BALANCE", dateakhir, strconv.FormatUint(acaKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

					endingbalance := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(productKeyLast, 10), customerKey, dateakhir)
						if err != nil {
							endingbalance["avg_nav"] = nol
						} else {
							endingbalance["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(productKeyLast, 10), dateakhir)
						if err != nil {
							endingbalance["nav_value"] = nol
						} else {
							endingbalance["nav_value"] = lastnav.NavValue.Truncate(2)
						}

						dateParem, _ = time.Parse(layout, dateakhir)
						endingbalance["date"] = dateParem.Format(newLayout)
						endingbalance["description"] = "ENDING BALANCE"
						endingbalance["amount"] = nol
						endingbalance["unit"] = nol
						endingbalance["fee"] = nol
						transGroupSales = append(transGroupSales, endingbalance)
					} else {
						endingbalance["date"] = balanceEnding.Tanggal
						endingbalance["description"] = balanceEnding.Description
						endingbalance["amount"] = balanceEnding.Amount.Truncate(0)
						endingbalance["nav_value"] = balanceEnding.NavValue.Truncate(2)
						endingbalance["unit"] = balanceEnding.Unit.Truncate(2)
						endingbalance["avg_nav"] = balanceEnding.AvgNav.Truncate(2)
						endingbalance["fee"] = balanceEnding.Fee.Truncate(0)
						transGroupSales = append(transGroupSales, endingbalance)
					}

					transGroupAgent["sales"] = sales
					transGroupAgent["transaksi"] = transGroupSales
					transGroupAgent["count"] = count

					transGroupProduct = append(transGroupProduct, transGroupAgent)
					row := make(map[string]interface{})
					row["product"] = product
					row["data"] = transGroupProduct
					datatrans = append(datatrans, row)

					//reset
					transGroupProduct = nil
					transGroupSales = nil
					product = make(map[string]interface{})
					transGroupSales = nil
					transGroupAgent = make(map[string]interface{})
					count = make(map[string]interface{})
					sales = make(map[string]interface{})
					totalSubs = nol
					totalRedm = nol
					totalNettsubs = nol
					salesKey = tr.SalesKey
					productKey = tr.ProductKey

					//set beginning
					var balanceBeginning models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalanceAca(&balanceBeginning, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AcaKey, 10), strconv.FormatUint(tr.ProductKey, 10))

					beginningbalance := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateawal)
						if err != nil {
							beginningbalance["avg_nav"] = nol
						} else {
							beginningbalance["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateawal)
						if err != nil {
							beginningbalance["nav_value"] = nol
						} else {
							beginningbalance["nav_value"] = lastnav.NavValue.Truncate(2)
						}

						dateParem, _ = time.Parse(layout, dateakhir)
						beginningbalance["date"] = dateParem.Format(newLayout)
						beginningbalance["description"] = "BEGINNING BALANCE"
						beginningbalance["amount"] = nol
						beginningbalance["unit"] = nol
						beginningbalance["fee"] = nol
						transGroupSales = append(transGroupSales, beginningbalance)
					} else {
						beginningbalance["date"] = balanceBeginning.Tanggal
						beginningbalance["description"] = balanceBeginning.Description
						beginningbalance["amount"] = balanceBeginning.Amount.Truncate(0)
						beginningbalance["nav_value"] = balanceBeginning.NavValue.Truncate(2)
						beginningbalance["unit"] = balanceBeginning.Unit.Truncate(2)
						beginningbalance["avg_nav"] = balanceBeginning.AvgNav.Truncate(2)
						beginningbalance["fee"] = balanceBeginning.Fee.Truncate(0)
						transGroupSales = append(transGroupSales, beginningbalance)
					}

					var balance models.SumBalanceUnit
					status, err = models.GetBalanceUnitByCustomerAndProduct(&balance, customerKey, strconv.FormatUint(tr.ProductKey, 10))
					if err != nil {
						product["unit"] = nol
						product["amount"] = nol
					} else {
						product["unit"] = balance.Unit.Truncate(2)
						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), time.Now().Format("2006-01-02"))
						if err != nil {
							product["amount"] = nol
						} else {
							product["amount"] = lastnav.NavValue.Mul(balance.Unit).Truncate(2)
							product["nav_value"] = lastnav.NavValue.Truncate(2)
						}

					}

					product["product_id"] = tr.ProductKey
					product["product_name"] = tr.ProductName
					product["currency"] = tr.Currency

					sales["sales_code"] = tr.SalesCode
					sales["sales_name"] = tr.SalesName

					if (tr.TransTypeKey == uint64(1)) || (tr.TransTypeKey == uint64(4)) {
						totalSubs = totalSubs.Add(tr.Amount).Truncate(0)
					}
					if (tr.TransTypeKey == uint64(2)) || (tr.TransTypeKey == uint64(3)) {
						totalRedm = totalRedm.Add(tr.Amount).Truncate(0)
					}
					totalNettsubs = totalNettsubs.Add(*tr.Fee).Truncate(0)

					count["subs"] = totalSubs
					count["redm"] = totalRedm
					count["nettsubs"] = totalNettsubs

					trans := make(map[string]interface{})
					trans["date"] = tr.NavDate
					trans["description"] = tr.Trans
					trans["amount"] = tr.Amount.Truncate(0)
					trans["nav_value"] = tr.NavValue.Truncate(2)
					trans["unit"] = tr.Unit.Truncate(2)
					trans["avg_nav"] = tr.AvgNav.Truncate(2)
					trans["fee"] = tr.Fee.Truncate(0)
					transGroupSales = append(transGroupSales, trans)

					if idx == (len(transactions) - 1) {
						//set ending
						var balanceBeginning models.BeginningEndingBalance
						_, err = models.GetBeginningEndingBalanceAca(&balanceBeginning, "ENDING BALANCE", dateakhir, strconv.FormatUint(tr.AcaKey, 10), strconv.FormatUint(tr.ProductKey, 10))

						beginningbalance := make(map[string]interface{})
						if err != nil {
							var lastavgnav models.NavValue
							_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateakhir)
							if err != nil {
								beginningbalance["avg_nav"] = nol
							} else {
								beginningbalance["avg_nav"] = lastavgnav.NavValue.Truncate(2)
							}

							var lastnav models.NavValue
							_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateakhir)
							if err != nil {
								beginningbalance["nav_value"] = nol
							} else {
								beginningbalance["nav_value"] = lastnav.NavValue.Truncate(2)
							}

							dateParem, _ = time.Parse(layout, dateakhir)
							beginningbalance["date"] = dateParem.Format(newLayout)
							beginningbalance["description"] = "ENDING BALANCE"
							beginningbalance["amount"] = nol
							beginningbalance["unit"] = nol
							beginningbalance["fee"] = nol
							transGroupSales = append(transGroupSales, beginningbalance)
						} else {
							beginningbalance["date"] = balanceBeginning.Tanggal
							beginningbalance["description"] = balanceBeginning.Description
							beginningbalance["amount"] = balanceBeginning.Amount.Truncate(0)
							beginningbalance["nav_value"] = balanceBeginning.NavValue.Truncate(2)
							beginningbalance["unit"] = balanceBeginning.Unit.Truncate(2)
							beginningbalance["avg_nav"] = balanceBeginning.AvgNav.Truncate(2)
							beginningbalance["fee"] = balanceBeginning.Fee.Truncate(0)
							transGroupSales = append(transGroupSales, beginningbalance)
						}

						transGroupAgent["sales"] = sales
						transGroupAgent["transaksi"] = transGroupSales
						transGroupAgent["count"] = count

						transGroupProduct = append(transGroupProduct, transGroupAgent)

						row := make(map[string]interface{})
						row["product"] = product
						row["data"] = transGroupProduct
						datatrans = append(datatrans, row)
					}
				} else {
					if salesKey != tr.SalesKey {
						// set ending
						var endingBel models.BeginningEndingBalance
						_, err = models.GetBeginningEndingBalanceAca(&endingBel, "ENDING BALANCE", dateakhir, strconv.FormatUint(acaKeyLast, 10), strconv.FormatUint(productKeyLast, 10))

						endingBall := make(map[string]interface{})
						if err != nil {
							var lastavgnav models.NavValue
							_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateakhir)
							if err != nil {
								endingBall["avg_nav"] = nol
							} else {
								endingBall["avg_nav"] = lastavgnav.NavValue.Truncate(2)
							}

							var lastnav models.NavValue
							_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateakhir)
							if err != nil {
								endingBall["nav_value"] = nol
							} else {
								endingBall["nav_value"] = lastnav.NavValue.Truncate(2)
							}

							dateParem, _ = time.Parse(layout, dateakhir)
							endingBall["date"] = dateParem.Format(newLayout)
							endingBall["description"] = "ENDING BALANCE"
							endingBall["amount"] = nol
							endingBall["unit"] = nol
							endingBall["fee"] = nol
							transGroupSales = append(transGroupSales, endingBall)
						} else {
							endingBall["date"] = endingBel.Tanggal
							endingBall["description"] = endingBel.Description
							endingBall["amount"] = endingBel.Amount.Truncate(0)
							endingBall["nav_value"] = endingBel.NavValue.Truncate(2)
							endingBall["unit"] = endingBel.Unit.Truncate(2)
							endingBall["avg_nav"] = endingBel.AvgNav.Truncate(2)
							endingBall["fee"] = endingBel.Fee.Truncate(0)
							transGroupSales = append(transGroupSales, endingBall)
						}

						transGroupAgent["sales"] = sales
						transGroupAgent["transaksi"] = transGroupSales
						transGroupAgent["count"] = count
						transGroupProduct = append(transGroupProduct, transGroupAgent)

						//reset
						transGroupSales = nil
						transGroupAgent = make(map[string]interface{})
						count = make(map[string]interface{})
						sales = make(map[string]interface{})
						totalSubs = nol
						totalRedm = nol
						totalNettsubs = nol
						salesKey = tr.SalesKey

						sales["sales_code"] = tr.SalesCode
						sales["sales_name"] = tr.SalesName

						if (tr.TransTypeKey == uint64(1)) || (tr.TransTypeKey == uint64(4)) {
							totalSubs = totalSubs.Add(tr.Amount).Truncate(0)
						}
						if (tr.TransTypeKey == uint64(2)) || (tr.TransTypeKey == uint64(3)) {
							totalRedm = totalRedm.Add(tr.Amount).Truncate(0)
						}
						totalNettsubs = totalNettsubs.Add(*tr.Fee).Truncate(0)

						count["subs"] = totalSubs
						count["redm"] = totalRedm
						count["nettsubs"] = totalNettsubs

						//set beginning
						var begginingBel models.BeginningEndingBalance
						_, err = models.GetBeginningEndingBalanceAca(&begginingBel, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AcaKey, 10), strconv.FormatUint(tr.ProductKey, 10))

						bBall := make(map[string]interface{})
						if err != nil {
							var lastavgnav models.NavValue
							_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateawal)
							if err != nil {
								bBall["avg_nav"] = nol
							} else {
								bBall["avg_nav"] = lastavgnav.NavValue.Truncate(2)
							}

							var lastnav models.NavValue
							_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateawal)
							if err != nil {
								bBall["nav_value"] = nol
							} else {
								bBall["nav_value"] = lastnav.NavValue.Truncate(2)
							}

							dateParem, _ = time.Parse(layout, dateakhir)
							bBall["date"] = dateParem.Format(newLayout)
							bBall["description"] = "BEGINNING BALANCE"
							bBall["amount"] = nol
							bBall["unit"] = nol
							bBall["fee"] = nol
							transGroupSales = append(transGroupSales, bBall)
						} else {
							bBall["date"] = begginingBel.Tanggal
							bBall["description"] = begginingBel.Description
							bBall["amount"] = begginingBel.Amount.Truncate(0)
							bBall["nav_value"] = begginingBel.NavValue.Truncate(2)
							bBall["unit"] = begginingBel.Unit.Truncate(2)
							bBall["avg_nav"] = begginingBel.AvgNav.Truncate(2)
							bBall["fee"] = begginingBel.Fee.Truncate(0)
							transGroupSales = append(transGroupSales, bBall)
						}

						trans := make(map[string]interface{})
						trans["date"] = tr.NavDate
						trans["description"] = tr.Trans
						trans["amount"] = tr.Amount.Truncate(0)
						trans["nav_value"] = tr.NavValue.Truncate(2)
						trans["unit"] = tr.Unit.Truncate(2)
						trans["avg_nav"] = tr.AvgNav.Truncate(2)
						trans["fee"] = tr.Fee.Truncate(0)
						transGroupSales = append(transGroupSales, trans)

						if idx == (len(transactions) - 1) {
							//set ending
							var eBel models.BeginningEndingBalance
							_, err = models.GetBeginningEndingBalanceAca(&eBel, "ENDING BALANCE", dateakhir, strconv.FormatUint(tr.AcaKey, 10), strconv.FormatUint(tr.ProductKey, 10))

							eball := make(map[string]interface{})
							if err != nil {
								var lastavgnav models.NavValue
								_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateakhir)
								if err != nil {
									eball["avg_nav"] = nol
								} else {
									eball["avg_nav"] = lastavgnav.NavValue.Truncate(2)
								}

								var lastnav models.NavValue
								_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateakhir)
								if err != nil {
									eball["nav_value"] = nol
								} else {
									eball["nav_value"] = lastnav.NavValue.Truncate(2)
								}

								dateParem, _ = time.Parse(layout, dateakhir)
								eball["date"] = dateParem.Format(newLayout)
								eball["description"] = "ENDING BALANCE"
								eball["amount"] = nol
								eball["unit"] = nol
								eball["fee"] = nol
								transGroupSales = append(transGroupSales, eball)
							} else {
								eball["date"] = eBel.Tanggal
								eball["description"] = eBel.Description
								eball["amount"] = eBel.Amount.Truncate(0)
								eball["nav_value"] = eBel.NavValue.Truncate(2)
								eball["unit"] = eBel.Unit.Truncate(2)
								eball["avg_nav"] = eBel.AvgNav.Truncate(2)
								eball["fee"] = eBel.Fee.Truncate(0)
								transGroupSales = append(transGroupSales, eball)
							}

							transGroupAgent["sales"] = sales
							transGroupAgent["transaksi"] = transGroupSales
							transGroupAgent["count"] = count

							transGroupProduct = append(transGroupProduct, transGroupAgent)

							row := make(map[string]interface{})
							row["product"] = product
							row["data"] = transGroupProduct
							datatrans = append(datatrans, row)
						}
					} else {
						if (tr.TransTypeKey == uint64(1)) || (tr.TransTypeKey == uint64(4)) {
							totalSubs = totalSubs.Add(tr.Amount).Truncate(0)
						}
						if (tr.TransTypeKey == uint64(2)) || (tr.TransTypeKey == uint64(3)) {
							totalRedm = totalRedm.Add(tr.Amount).Truncate(0)
						}
						totalNettsubs = totalNettsubs.Add(*tr.Fee).Truncate(0)

						count["subs"] = totalSubs
						count["redm"] = totalRedm
						count["nettsubs"] = totalNettsubs

						trans := make(map[string]interface{})
						trans["date"] = tr.NavDate
						trans["description"] = tr.Trans
						trans["amount"] = tr.Amount.Truncate(0)
						trans["nav_value"] = tr.NavValue.Truncate(2)
						trans["unit"] = tr.Unit.Truncate(2)
						trans["avg_nav"] = tr.AvgNav.Truncate(2)
						trans["fee"] = tr.Fee.Truncate(0)
						transGroupSales = append(transGroupSales, trans)

						if idx == (len(transactions) - 1) {
							//set ending
							var eBel models.BeginningEndingBalance
							_, err = models.GetBeginningEndingBalanceAca(&eBel, "ENDING BALANCE", dateakhir, strconv.FormatUint(tr.AcaKey, 10), strconv.FormatUint(tr.ProductKey, 10))

							eball := make(map[string]interface{})
							if err != nil {
								var lastavgnav models.NavValue
								_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateakhir)
								if err != nil {
									eball["avg_nav"] = nol
								} else {
									eball["avg_nav"] = lastavgnav.NavValue.Truncate(2)
								}

								var lastnav models.NavValue
								_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateakhir)
								if err != nil {
									eball["nav_value"] = nol
								} else {
									eball["nav_value"] = lastnav.NavValue.Truncate(2)
								}

								dateParem, _ = time.Parse(layout, dateakhir)
								eball["date"] = dateParem.Format(newLayout)
								eball["description"] = "ENDING BALANCE"
								eball["amount"] = nol
								eball["unit"] = nol
								eball["fee"] = nol
								transGroupSales = append(transGroupSales, eball)
							} else {
								eball["date"] = eBel.Tanggal
								eball["description"] = eBel.Description
								eball["amount"] = eBel.Amount.Truncate(0)
								eball["nav_value"] = eBel.NavValue.Truncate(2)
								eball["unit"] = eBel.Unit.Truncate(2)
								eball["avg_nav"] = eBel.AvgNav.Truncate(2)
								eball["fee"] = eBel.Fee.Truncate(0)
								transGroupSales = append(transGroupSales, eball)
							}

							transGroupAgent["sales"] = sales
							transGroupAgent["transaksi"] = transGroupSales
							transGroupAgent["count"] = count

							transGroupProduct = append(transGroupProduct, transGroupAgent)

							row := make(map[string]interface{})
							row["product"] = product
							row["data"] = transGroupProduct
							datatrans = append(datatrans, row)
						}
					}
				}
			} else {
				salesKey = tr.SalesKey
				//set beginning
				var balanceBeginning models.BeginningEndingBalance
				_, err = models.GetBeginningEndingBalanceAca(&balanceBeginning, "BEGINNING BALANCE", dateawal, strconv.FormatUint(tr.AcaKey, 10), strconv.FormatUint(tr.ProductKey, 10))

				beginningbalance := make(map[string]interface{})
				if err != nil {
					var lastavgnav models.NavValue
					_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateawal)
					if err != nil {
						beginningbalance["avg_nav"] = nol
					} else {
						beginningbalance["avg_nav"] = lastavgnav.NavValue.Truncate(2)
					}

					var lastnav models.NavValue
					_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateawal)
					if err != nil {
						beginningbalance["nav_value"] = nol
					} else {
						beginningbalance["nav_value"] = lastnav.NavValue.Truncate(2)
					}

					dateParem, _ = time.Parse(layout, dateakhir)
					beginningbalance["date"] = dateParem.Format(newLayout)
					beginningbalance["description"] = "BEGINNING BALANCE"
					beginningbalance["amount"] = nol
					beginningbalance["unit"] = nol
					beginningbalance["fee"] = nol
					transGroupSales = append(transGroupSales, beginningbalance)
				} else {
					beginningbalance["date"] = balanceBeginning.Tanggal
					beginningbalance["description"] = balanceBeginning.Description
					beginningbalance["amount"] = balanceBeginning.Amount.Truncate(0)
					beginningbalance["nav_value"] = balanceBeginning.NavValue.Truncate(2)
					beginningbalance["unit"] = balanceBeginning.Unit.Truncate(2)
					beginningbalance["avg_nav"] = balanceBeginning.AvgNav.Truncate(2)
					beginningbalance["fee"] = balanceBeginning.Fee.Truncate(0)
					transGroupSales = append(transGroupSales, beginningbalance)
				}

				var balance models.SumBalanceUnit
				status, err = models.GetBalanceUnitByCustomerAndProduct(&balance, customerKey, strconv.FormatUint(tr.ProductKey, 10))
				if err != nil {
					product["unit"] = nol
					product["amount"] = nol
				} else {
					product["unit"] = balance.Unit.Truncate(2)
					var lastnav models.NavValue
					_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), time.Now().Format("2006-01-02"))
					if err != nil {
						product["amount"] = nol
					} else {
						product["amount"] = lastnav.NavValue.Mul(balance.Unit).Truncate(2)
						product["nav_value"] = lastnav.NavValue.Truncate(2)
					}

				}

				product["product_id"] = tr.ProductKey
				product["product_name"] = tr.ProductName
				product["currency"] = tr.Currency

				sales["sales_code"] = tr.SalesCode
				sales["sales_name"] = tr.SalesName

				if (tr.TransTypeKey == uint64(1)) || (tr.TransTypeKey == uint64(4)) {
					totalSubs = totalSubs.Add(tr.Amount).Truncate(0)
				}
				if (tr.TransTypeKey == uint64(2)) || (tr.TransTypeKey == uint64(3)) {
					totalRedm = totalRedm.Add(tr.Amount).Truncate(0)
				}
				totalNettsubs = totalNettsubs.Add(*tr.Fee).Truncate(0)

				count["subs"] = totalSubs
				count["redm"] = totalRedm
				count["nettsubs"] = totalNettsubs

				trans := make(map[string]interface{})
				trans["date"] = tr.NavDate
				trans["description"] = tr.Trans
				trans["amount"] = tr.Amount.Truncate(0)
				trans["nav_value"] = tr.NavValue.Truncate(2)
				trans["unit"] = tr.Unit.Truncate(2)
				trans["avg_nav"] = tr.AvgNav.Truncate(2)
				trans["fee"] = tr.Fee.Truncate(0)
				transGroupSales = append(transGroupSales, trans)

				if idx == (len(transactions) - 1) {
					//set ending
					var eBel models.BeginningEndingBalance
					_, err = models.GetBeginningEndingBalanceAca(&eBel, "ENDING BALANCE", dateakhir, strconv.FormatUint(tr.AcaKey, 10), strconv.FormatUint(tr.ProductKey, 10))

					eball := make(map[string]interface{})
					if err != nil {
						var lastavgnav models.NavValue
						_, err = models.AdminLastAvgNav(&lastavgnav, strconv.FormatUint(tr.ProductKey, 10), customerKey, dateakhir)
						if err != nil {
							eball["avg_nav"] = nol
						} else {
							eball["avg_nav"] = lastavgnav.NavValue.Truncate(2)
						}

						var lastnav models.NavValue
						_, err = models.AdminLastNavValue(&lastnav, strconv.FormatUint(tr.ProductKey, 10), dateakhir)
						if err != nil {
							eball["nav_value"] = nol
						} else {
							eball["nav_value"] = lastnav.NavValue.Truncate(2)
						}

						dateParem, _ = time.Parse(layout, dateakhir)
						eball["date"] = dateParem.Format(newLayout)
						eball["description"] = "ENDING BALANCE"
						eball["amount"] = nol
						eball["unit"] = nol
						eball["fee"] = nol
						transGroupSales = append(transGroupSales, eball)
					} else {
						eball["date"] = eBel.Tanggal
						eball["description"] = eBel.Description
						eball["amount"] = eBel.Amount.Truncate(0)
						eball["nav_value"] = eBel.NavValue.Truncate(2)
						eball["unit"] = eBel.Unit.Truncate(2)
						eball["avg_nav"] = eBel.AvgNav.Truncate(2)
						eball["fee"] = eBel.Fee.Truncate(0)
						transGroupSales = append(transGroupSales, eball)
					}

					transGroupAgent["sales"] = sales
					transGroupAgent["transaksi"] = transGroupSales
					transGroupAgent["count"] = count

					transGroupProduct = append(transGroupProduct, transGroupAgent)

					row := make(map[string]interface{})
					row["product"] = product
					row["data"] = transGroupProduct
					datatrans = append(datatrans, row)
				}

			}
			acaKeyLast = tr.AcaKey
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
