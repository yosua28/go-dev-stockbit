package models

import (
	"api/config"
	"api/db"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type ResponseDailySubscriptionBatchForm struct {
	Header HeaderDailySubsRedmBatchForm      `json:"header"`
	Data   *[]ResponseDailySubsRedmBatchForm `json:"data_list"`
	Count  CountNominal                      `json:"count"`
}

type HeaderDailySubsRedmBatchForm struct {
	Logo            string `db:"logo"                json:"logo"`
	ProductName     string `db:"product_name"        json:"product_name"`
	NoRekProduct    string `db:"no_rek_product"      json:"no_rek_product"`
	TradeDate       string `db:"trade_date"          json:"trade_date"`
	ReferenceNo     string `db:"reference_no"        json:"reference_no"`
	BankNameProduct string `db:"bank_name_product"   json:"bank_name_product"`
}
type CountNominal struct {
	CountUnit       *decimal.Decimal `json:"count_unit,omitempty"`
	CountAmount     decimal.Decimal  `json:"count_amount"`
	CountFeeAmount  decimal.Decimal  `json:"count_fee_amount"`
	CountNettAmount decimal.Decimal  `json:"count_nett_amount"`
}

type DailySubsRedmBatchForm struct {
	ProductKey      string           `db:"product_key"           json:"product_key"`
	CustomerKey     string           `db:"customer_key"          json:"customer_key"`
	Sid             *string          `db:"sid"                   json:"sid"`
	IfuaNo          *string          `db:"ifua_no"               json:"ifua_no"`
	AccountNo       *string          `db:"account_no"            json:"account_no"`
	UnitHolderIDNo  *string          `db:"unit_holder_idno"      json:"unit_holder_idno"`
	FullName        string           `db:"full_name"             json:"full_name"`
	Amount          decimal.Decimal  `db:"amount"                json:"amount"`
	FeeAmount       decimal.Decimal  `db:"fee_amount"            json:"fee_amount"`
	NettAmount      decimal.Decimal  `db:"nett_amount"           json:"nett_amount"`
	BankFullName    *string          `db:"bank_fullname"         json:"bank_fullname"`
	NoRekening      *string          `db:"no_rekening"           json:"no_rekening"`
	TypeDescription string           `db:"type_description"      json:"type_description"`
	Unit            *decimal.Decimal `db:"unit"                  json:"unit"`
	PaymentDate     *string          `db:"payment_date"          json:"payment_date"`
}

type ResponseDailySubsRedmBatchForm struct {
	Sid             *string          `json:"sid"`
	IfuaNo          *string          `json:"ifua_no"`
	AccountNo       *string          `json:"account_no"`
	UnitHolderIDNo  *string          `json:"unit_holder_idno"`
	FullName        string           `json:"full_name"`
	Unit            *decimal.Decimal `json:"unit,omitempty"`
	Amount          decimal.Decimal  `json:"amount"`
	FeeAmount       decimal.Decimal  `json:"fee_amount"`
	NettAmount      decimal.Decimal  `json:"nett_amount"`
	BankFullName    *string          `json:"bank_fullname"`
	NoRekening      *string          `json:"no_rekening"`
	TypeDescription string           `json:"type_description"`
	PaymentDate     *string          `json:"payment_date,omitempty"`
	Notes1          *string          `json:"notes_1,omitempty"`
	Notes2          *string          `json:"notes_2,omitempty"`
	Notes3          *string          `json:"notes_3,omitempty"`
}

type NotesRedemption struct {
	Note1  *string         `db:"note1"           json:"note1"`
	Unit   decimal.Decimal `db:"unit"            json:"unit"`
	Amount decimal.Decimal `db:"amount"          json:"amount"`
	Note3  *string         `db:"note3"           json:"note3"`
}

type DailyTransactionReportField struct {
	ClientName         string          `db:"client_name"           json:"client_name"`
	Product            string          `db:"product"               json:"product"`
	SubscriptionAmount decimal.Decimal `db:"subscription_amount"   json:"subscription_amount"`
	SubscriptionFee    decimal.Decimal `db:"subscription_fee"      json:"subscription_fee"`
	RedemptionAmount   decimal.Decimal `db:"redemption_amount"     json:"redemption_amount"`
	RedemptionFee      decimal.Decimal `db:"redemption_fee"        json:"redemption_fee"`
	Category           *string         `db:"category"              json:"category"`
	Division           *string         `db:"division"              json:"division"`
	Branch             *string         `db:"branch"                json:"branch"`
	Sales              *string         `db:"sales"                 json:"sales"`
}

type DailyTransactionReportTotalField struct {
	TotalSubscriptionAmount decimal.Decimal `db:"total_subscription_amount"   json:"total_subscription_amount"`
	TotalSubscriptionFee    decimal.Decimal `db:"total_subscription_fee"      json:"total_subscription_fee"`
	TotalRedemptionAmount   decimal.Decimal `db:"total_redemption_amount"     json:"total_redemption_amount"`
	TotalRedemptionFee      decimal.Decimal `db:"total_redemption_fee"        json:"total_redemption_fee"`
}

type DailyTransactionReportResponse struct {
	Total                  *DailyTransactionReportTotalField `json:"total"`
	DailyTransactionReport *[]DailyTransactionReportField    `json:"data"`
}

func AdminGetHeaderDailySubsRedmBatchForm(c *HeaderDailySubsRedmBatchForm, params map[string]string) (int, error) {
	query := `SELECT 
			concat("` + config.BaseUrl + `", "/images/mail/report_logo_mnc.jpg") AS logo,
			p.product_name_alt AS product_name,
			DATE_FORMAT(t.nav_date, '%d %M %Y') AS trade_date,
			batch.batch_display_no AS reference_no,
			b.bank_fullname AS bank_name_product,
			mba.account_no AS no_rek_product 
		FROM tr_transaction AS t 
		INNER JOIN ms_product AS p ON p.product_key = t.product_key 
		INNER JOIN tr_transaction_batch AS batch ON batch.batch_key = t.batch_key 
		INNER JOIN tr_transaction_bank_account AS ba ON ba.transaction_key = t.transaction_key 
		INNER JOIN ms_product_bank_account AS pba ON pba.prod_bankacc_key = ba.prod_bankacc_key 
		INNER JOIN ms_bank_account AS mba ON mba.bank_account_key = pba.bank_account_key 
		INNER JOIN ms_bank AS b ON b.bank_key = mba.bank_key 
		WHERE t.rec_status = 1 AND t.trans_status_key >= 6`

	var condition string

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	query += " LIMIT 1"

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetDailySubsRedmBatchForm(c *[]DailySubsRedmBatchForm, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT 
				t.product_key AS product_key,
				t.customer_key AS customer_key,
				c.sid_no AS sid,
				acc.ifua_no AS ifua_no,
				acc.account_no AS account_no,
				c.unit_holder_idno AS unit_holder_idno,
				c.full_name AS full_name,
				(t.trans_amount + t.trans_fee_amount) AS amount,
				t.trans_fee_amount AS fee_amount,
				t.total_amount AS nett_amount,
				b.bank_fullname AS bank_fullname,
				mba.account_no AS no_rekening,
				ttt.type_description AS type_description,
				(CASE
					WHEN t.trans_unit IS NULL OR t.trans_unit = 0 THEN 
						(CASE
							WHEN tc.tc_key IS NOT NULL THEN tc.confirmed_unit
							ELSE 0
						END)
					ELSE t.trans_unit
				END) AS unit,
				DATE_FORMAT(t.settlement_date, '%d %M %Y') AS payment_date     
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_type AS ttt ON ttt.trans_type_key = t.trans_type_key 
			LEFT JOIN tr_transaction_confirmation AS tc ON t.transaction_key = tc.transaction_key 
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key 
			INNER JOIN tr_account_agent AS aa ON aa.aca_key = t.aca_key 
			INNER JOIN tr_account AS acc ON acc.acc_key = aa.acc_key 
			LEFT JOIN tr_transaction_bank_account AS ba ON ba.transaction_key = t.transaction_key 
			LEFT JOIN ms_customer_bank_account AS cba ON cba.cust_bankacc_key = ba.cust_bankacc_key 
			LEFT JOIN ms_bank_account AS mba ON mba.bank_account_key = cba.bank_account_key 
			LEFT JOIN ms_bank AS b ON b.bank_key = mba.bank_key
			WHERE t.rec_status = 1 AND t.trans_status_key >= 6`

	var present bool
	var condition string

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminCountDailySubsRedmBatchForm(c *CountData, params map[string]string) (int, error) {
	query := `SELECT 
				count(t.transaction_key) AS count_data 
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_type AS ttt ON ttt.trans_type_key = t.trans_type_key 
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key 
			INNER JOIN tr_account_agent AS aa ON aa.aca_key = t.aca_key 
			INNER JOIN tr_account AS acc ON acc.acc_key = aa.acc_key 
			LEFT JOIN tr_transaction_bank_account AS ba ON ba.transaction_key = t.transaction_key 
			LEFT JOIN ms_customer_bank_account AS cba ON cba.cust_bankacc_key = ba.cust_bankacc_key 
			LEFT JOIN ms_bank_account AS mba ON mba.bank_account_key = cba.bank_account_key 
			LEFT JOIN ms_bank AS b ON b.bank_key = mba.bank_key
			WHERE t.rec_status = 1 AND t.trans_status_key >= 6 `

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type BankProductTransactionReport struct {
	ProductKey  string `db:"prod_bankacc_key"    json:"prod_bankacc_key"`
	CustomerKey string `db:"bank_name"           json:"bank_name"`
}

func AdminGetBankProductTransactionReport(c *[]BankProductTransactionReport, params map[string]string) (int, error) {
	query := `SELECT 
				ba.prod_bankacc_key AS prod_bankacc_key, 
				CONCAT(b.bank_fullname, ' - ', '(', mba.account_no, ' / ', mba.account_holder_name, ')') AS bank_name 
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_type AS ttt ON ttt.trans_type_key = t.trans_type_key 
			INNER JOIN tr_transaction_bank_account AS ba ON ba.transaction_key = t.transaction_key 
			INNER JOIN ms_product_bank_account AS pba ON pba.prod_bankacc_key = ba.prod_bankacc_key 
			INNER JOIN ms_bank_account AS mba ON mba.bank_account_key = pba.bank_account_key 
			INNER JOIN ms_bank AS b ON b.bank_key = mba.bank_key 
			WHERE t.rec_status = 1 AND t.trans_status_key != 3 `

	var condition string

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	query += " GROUP BY ba.prod_bankacc_key"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetNotesRedemption(c *NotesRedemption, customerKey string, productKey string, navDate string) (int, error) {
	query := `SELECT 
				CONCAT(ty.type_code,' - ', DATE_FORMAT(t.nav_date, '%d %M %Y'), ' - ', p.product_code) AS note1,
				(CASE
					WHEN t.trans_unit IS NULL OR t.trans_unit = 0 THEN 
						(CASE
							WHEN tc.tc_key IS NOT NULL THEN tc.confirmed_unit
							ELSE 0
						END)
					ELSE t.trans_unit
				END) AS unit,
				(CASE
					WHEN t.trans_amount IS NULL OR t.trans_amount = 0 THEN 
						(CASE
							WHEN tc.tc_key IS NOT NULL THEN tc.confirmed_amount
							ELSE 0
						END)
					ELSE t.trans_amount
				END) AS amount,
				CONCAT('SettDate : ', DATE_FORMAT(t.settlement_date, '%d %M %Y')) AS note3 
			FROM tr_transaction AS t 
			LEFT JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key
			LEFT JOIN tr_transaction_type AS ty ON ty.trans_type_key = t.trans_type_key 
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			WHERE t.rec_status = 1 AND t.trans_type_key = 3 AND t.trans_status_key != 3 AND 
			t.transaction_key = 
			(
				SELECT parent_key FROM tr_transaction 
				WHERE rec_status = 1 
				AND nav_date < '` + navDate + `' AND trans_type_key = 4 AND customer_key = '` + customerKey + `' AND product_key = '` + productKey + `' 
				ORDER BY transaction_key DESC LIMIT 1
			)
			AND DATE_ADD(t.nav_date, INTERVAL 
				(SELECT pr.settlement_period
				FROM tr_transaction AS tr
				INNER JOIN ms_product AS pr ON pr.product_key = tr.product_key
				WHERE tr.rec_status = 1 AND tr.nav_date < '` + navDate + `' AND tr.trans_type_key = 4 
				AND tr.customer_key = '` + customerKey + `' AND tr.product_key = '` + productKey + `' 
				ORDER BY tr.transaction_key DESC LIMIT 1) DAY) > '` + navDate + `' 
			ORDER BY t.transaction_key DESC LIMIT 1`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func DailyTransactionReport(c *[]DailyTransactionReportField, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	dateFrom := ""
	dateTo := ""

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "dateFrom" || field == "dateTo") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
		if field == "dateFrom" {
			dateFrom = value
		}
		if field == "dateTo" {
			dateTo = value
		}
	}

	query := `SELECT 
				c.full_name AS client_name,
				CONCAT(p.product_code, ' - ', p.product_name_alt) AS product,
				SUM(CASE WHEN t.trans_type_key = 1 
						THEN trans_amount ELSE 0 END) subscription_amount,
				SUM(CASE WHEN t.trans_type_key = 1 
						THEN (t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) ELSE 0 END) subscription_fee,
				SUM(CASE WHEN t.trans_type_key = 2 
						THEN trans_amount ELSE 0 END) redemption_amount,
				SUM(CASE WHEN t.trans_type_key = 2 
						THEN (t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) ELSE 0 END) redemption_fee,
				ct.lkp_name AS category,
				division.lkp_name AS division,
				b.branch_name AS branch,
				a.agent_name AS sales 
			FROM tr_transaction t 
			INNER JOIN ms_customer AS c ON t.customer_key = c.customer_key 
			INNER JOIN ms_product AS p ON t.product_key = p.product_key 
			INNER JOIN gen_lookup AS ct ON ct.lookup_key = c.investor_type
			INNER JOIN gen_lookup AS division ON division.lookup_key = c.customer_category
			INNER JOIN tr_account_agent AS taa ON taa.aca_key = t.aca_key 
			INNER JOIN ms_agent AS a ON taa.agent_key = a.agent_key 
			INNER JOIN ms_branch AS b ON b.branch_key = t.branch_key 
			WHERE t.rec_status = 1 AND t.trans_type_key IN (1,2) AND t.trans_status_key >= 6 `

	var present bool
	var condition string
	var conditionOrder string

	if (dateFrom != "") && (dateTo != "") {
		query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	query += " GROUP BY t.customer_key, t.product_key"

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		conditionOrder += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			conditionOrder += " " + orderType
		}
	}
	query += conditionOrder

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func DailyTransactionReportTotal(c *DailyTransactionReportTotalField, params map[string]string) (int, error) {
	dateFrom := ""
	dateTo := ""

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "dateFrom" || field == "dateTo") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
		if field == "dateFrom" {
			dateFrom = value
		}
		if field == "dateTo" {
			dateTo = value
		}
	}

	query := `SELECT 
				SUM(CASE WHEN t.trans_type_key = 1 
						THEN trans_amount ELSE 0 END) total_subscription_amount,
				SUM(CASE WHEN t.trans_type_key = 1 
						THEN (t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) ELSE 0 END) total_subscription_fee,
				SUM(CASE WHEN t.trans_type_key = 2 
						THEN trans_amount ELSE 0 END) total_redemption_amount,
				SUM(CASE WHEN t.trans_type_key = 2 
						THEN (t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) ELSE 0 END) total_redemption_fee 
			FROM tr_transaction t 
			INNER JOIN ms_customer AS c ON t.customer_key = c.customer_key 
			INNER JOIN ms_product AS p ON t.product_key = p.product_key 
			INNER JOIN gen_lookup AS ct ON ct.lookup_key = c.investor_type
			INNER JOIN gen_lookup AS division ON division.lookup_key = c.customer_category
			INNER JOIN tr_account_agent AS taa ON taa.aca_key = t.aca_key 
			INNER JOIN ms_agent AS a ON taa.agent_key = a.agent_key 
			INNER JOIN ms_branch AS b ON b.branch_key = t.branch_key 
			WHERE t.rec_status = 1 AND t.trans_type_key IN (1,2) AND t.trans_status_key >= 6 `

	var condition string

	if (dateFrom != "") && (dateTo != "") {
		query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func DailyTransactionReportCountRow(c *CountData, params map[string]string) (int, error) {
	dateFrom := ""
	dateTo := ""

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "dateFrom" || field == "dateTo") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
		if field == "dateFrom" {
			dateFrom = value
		}
		if field == "dateTo" {
			dateTo = value
		}
	}

	query := `SELECT 
				COUNT(trans.jml_row) AS count_data 
			  FROM
			  (
				SELECT
					COUNT(*) AS jml_row
				FROM tr_transaction t 
				INNER JOIN ms_customer AS c ON t.customer_key = c.customer_key 
				INNER JOIN ms_product AS p ON t.product_key = p.product_key 
				INNER JOIN gen_lookup AS ct ON ct.lookup_key = c.investor_type
				INNER JOIN gen_lookup AS division ON division.lookup_key = c.customer_category
				INNER JOIN tr_account_agent AS taa ON taa.aca_key = t.aca_key 
				INNER JOIN ms_agent AS a ON taa.agent_key = a.agent_key 
				INNER JOIN ms_branch AS b ON b.branch_key = t.branch_key 
				WHERE t.rec_status = 1 AND t.trans_type_key IN (1,2) AND t.trans_status_key >= 6 `

	var condition string

	if (dateFrom != "") && (dateTo != "") {
		query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	query += condition

	query += " GROUP BY t.customer_key, t.product_key) trans"

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
