package models

import (
	"api/db"
	"net/http"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type AccountStatementCustomerProduct struct {
	TransactionKey uint64           `db:"transaction_key"         json:"transaction_key"`
	TransTypeKey   uint64           `db:"trans_type_key"          json:"trans_type_key"`
	CustomerKey    uint64           `db:"customer_key"            json:"customer_key"`
	AccKey         uint64           `db:"acc_key"                 json:"acc_key"`
	ProductKey     uint64           `db:"product_key"             json:"product_key"`
	ProductName    string           `db:"product_name"            json:"product_name"`
	ProductCode    string           `db:"product_code"            json:"product_code"`
	Trans          string           `db:"trans"                   json:"trans"`
	NavDate        string           `db:"nav_date"                json:"nav_date"`
	NavValue       decimal.Decimal  `db:"nav_value"               json:"nav_value"`
	AvgNav         *decimal.Decimal `db:"avg_nav"                 json:"avg_nav"`
	Amount         decimal.Decimal  `db:"confirmed_amount"        json:"confirmed_amount"`
	Unit           decimal.Decimal  `db:"confirmed_unit"          json:"confirmed_unit"`
	Fee            *decimal.Decimal `db:"fee"                     json:"fee"`
	Currency       string           `db:"currency"                json:"currency"`
	BankName       string           `db:"bank_name"               json:"bank_name"`
	AccountName    string           `db:"account_holder_name"     json:"account_holder_name"`
	AccountNo      string           `db:"account_no"              json:"account_no"`
}

func AdminGetAllAccountStatementCustomerProduct(c *[]AccountStatementCustomerProduct, customerKey string, dateFrom string, dateTo string) (int, error) {
	query := `SELECT 
				t.transaction_key,
				t.trans_type_key,
				t.customer_key,
				ta.acc_key,
				t.product_key,
				p.product_name_alt AS product_name,
				p.product_code,
				tt.type_code as trans,
				DATE_FORMAT(t.nav_date, '%d %M %Y') AS nav_date,
				nav.nav_value,
				tc.avg_nav,
				tc.confirmed_amount,
				tc.confirmed_unit,
				(t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) AS fee,
				cur.code AS currency,
				b.bank_name,
				ba.account_holder_name,
				ba.account_no 
			FROM tr_transaction AS t
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key AND c.rec_status =  1
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key AND tc.rec_status = 1
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			INNER JOIN tr_transaction_type AS tt ON tt.trans_type_key = t.trans_type_key
			INNER JOIN tr_nav AS nav ON nav.nav_date = t.nav_date AND nav.product_key = t.product_key
			INNER JOIN ms_currency AS cur ON cur.currency_key = p.currency_key
			INNER JOIN ms_product_bank_account AS pb ON pb.product_key = t.product_key AND pb.bank_account_purpose = '269'
			INNER JOIN ms_bank_account AS ba ON ba.bank_account_key = pb.bank_account_key 
			INNER JOIN ms_bank AS b ON b.bank_key = ba.bank_key 
			INNER JOIN tr_account as ta ON ta.customer_key = t.customer_key AND ta.product_key = t.product_key  
			WHERE t.customer_key = '` + customerKey + `' AND t.trans_status_key = 9 AND t.rec_status = 1
			AND (t.nav_date BETWEEN '` + dateFrom + `' AND '` + dateTo + `')
			GROUP BY t.transaction_key 
			ORDER BY t.product_key, t.nav_date ASC`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type AccountStatementCustomerAgent struct {
	TransactionKey uint64           `db:"transaction_key"         json:"transaction_key"`
	TransTypeKey   uint64           `db:"trans_type_key"          json:"trans_type_key"`
	AcaKey         uint64           `db:"aca_key"                 json:"aca_key"`
	SalesKey       uint64           `db:"sales_key"                 json:"sales_key"`
	SalesCode      string           `db:"sales_code"            json:"sales_code"`
	SalesName      string           `db:"sales_name"            json:"sales_name"`
	CustomerKey    uint64           `db:"customer_key"            json:"customer_key"`
	ProductKey     uint64           `db:"product_key"             json:"product_key"`
	ProductName    string           `db:"product_name"            json:"product_name"`
	ProductCode    string           `db:"product_code"            json:"product_code"`
	Trans          string           `db:"trans"                   json:"trans"`
	NavDate        string           `db:"nav_date"                json:"nav_date"`
	NavValue       decimal.Decimal  `db:"nav_value"               json:"nav_value"`
	AvgNav         *decimal.Decimal `db:"avg_nav"                 json:"avg_nav"`
	Amount         decimal.Decimal  `db:"confirmed_amount"        json:"confirmed_amount"`
	Unit           decimal.Decimal  `db:"confirmed_unit"          json:"confirmed_unit"`
	Fee            *decimal.Decimal `db:"fee"                     json:"fee"`
	Currency       string           `db:"currency"                json:"currency"`
}

func AdminGetAllAccountStatementCustomerAgent(c *[]AccountStatementCustomerAgent, customerKey string, dateFrom string, dateTo string) (int, error) {
	query := `SELECT 
				t.transaction_key,
				t.trans_type_key,
				t.aca_key,
				ag.agent_key AS sales_key,
				ag.agent_code AS sales_code,
				ag.agent_name AS sales_name, 
				t.customer_key,
				t.product_key,
				p.product_name_alt as product_name,
				p.product_code,
				tt.type_code as trans,
				DATE_FORMAT(t.nav_date, '%d %M %Y') AS nav_date,
				nav.nav_value,
				tc.avg_nav,
				tc.confirmed_amount,
				tc.confirmed_unit,
				(t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) AS fee,
				cur.code AS currency 
			FROM tr_transaction AS t
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key AND c.rec_status =  1
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key AND tc.rec_status = 1
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			INNER JOIN tr_transaction_type AS tt ON tt.trans_type_key = t.trans_type_key
			INNER JOIN tr_nav AS nav ON nav.nav_date = t.nav_date AND nav.product_key = t.product_key
			INNER JOIN ms_currency AS cur ON cur.currency_key = p.currency_key
			INNER JOIN tr_account_agent AS taa ON taa.aca_key = t.aca_key 
			INNER JOIN ms_agent AS ag ON ag.agent_key = taa.agent_key  
			WHERE t.customer_key = '` + customerKey + `' AND t.trans_status_key = 9 AND t.rec_status = 1
			AND (t.nav_date BETWEEN '` + dateFrom + `' AND '` + dateTo + `')
			GROUP BY t.transaction_key 
			ORDER BY t.product_key, ag.agent_key ASC`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
