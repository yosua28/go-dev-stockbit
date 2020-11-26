package models

import (
	"api/db"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type TrTransactionSubRedmReportResponse struct {
	SidNo             *string `json:"sid_no"`
	FullName          string  `json:"full_name"`
	TransAmount       float32 `json:"trans_amount"`
	TransFeeAmount    float32 `json:"trans_fee_amount"`
	TotalAmount       float32 `json:"total_amount"`
	CustodianFullName *string `json:"custodian_full_name"`
	TypeDescription   *string `json:"type_description"`
}

type TrTransactionSubRedmReport struct {
	SidNo             *string `db:"sid_no"                    json:"sid_no"`
	FullName          string  `db:"full_name"                 json:"full_name"`
	TransAmount       float32 `db:"trans_amount"              json:"trans_amount"`
	TransFeeAmount    float32 `db:"trans_fee_amount"          json:"trans_fee_amount"`
	TotalAmount       float32 `db:"total_amount"              json:"total_amount"`
	CustodianFullName *string `db:"custodian_full_name"       json:"custodian_full_name"`
	TypeDescription   *string `db:"type_description"          json:"type_description"`
}

func DailySubRedmReport(c *[]TrTransactionSubRedmReport, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
				cus.sid_no AS sid_no, 
				cus.full_name AS full_name, 
				tr.trans_amount AS trans_amount, 
				tr.trans_fee_amount AS trans_fee_amount, 
				tr.total_amount AS total_amount, 
				cb.custodian_full_name AS custodian_full_name, 
				ty.type_description AS type_description  
			  FROM tr_transaction AS tr 
			  INNER JOIN ms_product AS p ON p.product_key = tr.product_key 
			  INNER JOIN tr_transaction_type AS ty ON ty.trans_type_key = tr.trans_type_key 
			  INNER JOIN ms_customer AS cus ON cus.customer_key = tr.customer_key 
			  LEFT JOIN ms_custodian_bank cb ON cb.custodian_key = p.custodian_key 
			  WHERE tr.rec_status = 1 AND tr.trans_status_key != 3`
	var present bool
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
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
