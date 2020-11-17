package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"
)

type TrBalance struct {
	BalanceKey        uint64   `db:"balance_key"               json:"balance_key"`
	AcaKey            uint64   `db:"aca_key"                   json:"aca_key"`
	TcKey             uint64   `db:"tc_key"                    json:"tc_key"`
	BalanceDate       string   `db:"balance_date"              json:"balance_date"`
	BalanceUnit       float32  `db:"balance_unit"              json:"balance_unit"`
	AvgNav            *float32 `db:"avg_nav"                   json:"avg_nav"`
	TcKeyRed          *uint64  `db:"tc_key_red"                json:"tc_key_red"`
	RecOrder          *uint64  `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8    `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string  `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string  `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string  `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string  `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string  `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string  `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8   `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64  `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string  `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string  `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string  `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string  `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string  `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string  `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string  `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type TrBalanceCustomerProduk struct {
	BalanceKey     uint64  `db:"balance_key"               json:"balance_key"`
	AcaKey         uint64  `db:"aca_key"                   json:"aca_key"`
	BalanceUnit    float32 `db:"balance_unit"              json:"balance_unit"`
	TcKey          uint64  `db:"tc_key"                    json:"tc_key"`
	TransactionKey uint64  `db:"transaction_key"           json:"transaction_key"`
	NavDate        string  `db:"nav_date"                  json:"nav_date"`
}

func GetLastBalanceIn(c *[]TrBalance, acaKey []string) (int, error) {
	inQuery := strings.Join(acaKey, ",")
	query2 := `SELECT 
				t1.balance_key, 
				t1.aca_key, 
				t1.tc_key, 
				t1.balance_date, 
				t1.balance_unit, 
				t1.avg_nav, 
				t1.tc_key_red FROM
				   tr_balance t1 JOIN 
				   (SELECT MAX(balance_key) balance_key FROM tr_balance GROUP BY tc_key) t2
			   ON t1.balance_key = t2.balance_key`
	query := query2 + " WHERE t1.aca_key IN(" + inQuery + ") GROUP BY tc_key ORDER BY t1.balance_key DESC"

	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateTrBalance(params map[string]string) (int, error) {
	query := "INSERT INTO tr_balance"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	log.Println(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetLastBalanceCustomerByProductKey(c *[]TrBalanceCustomerProduk, customerKey string, productKey string) (int, error) {
	query := `SELECT 
				tb.balance_key as balance_key, 
				tb.aca_key as aca_key, 
				tb.balance_unit as balance_unit, 
				tc.tc_key as tc_key, 
				tr.transaction_key as transaction_key, 
				tr.nav_date as nav_date 
				FROM tr_balance AS tb
				JOIN (SELECT MAX(balance_key) balance_key FROM tr_balance GROUP BY tc_key) AS t2 ON tb.balance_key = t2.balance_key 
				INNER JOIN tr_transaction_confirmation AS tc ON tb.tc_key = tc.tc_key
				INNER JOIN tr_transaction AS tr ON tc.transaction_key = tr.transaction_key
				WHERE tr.customer_key = ` + customerKey +
		` AND tr.product_key = ` + productKey +
		` AND tr.trans_status_key = 9 AND tr.rec_status = 1 AND tr.trans_type_key = 1 AND tb.balance_unit > 0 
				GROUP BY tb.tc_key  ORDER BY tc.tc_key ASC`

	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetLastTrBalanceByTcRed(c *TrBalance, tcKeyRed string) (int, error) {
	query := `SELECT * FROM tr_balance WHERE tc_key_red = ` + tcKeyRed + ` ORDER BY rec_order DESC LIMIT 1`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
