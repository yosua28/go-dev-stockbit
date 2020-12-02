package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type TrTransactionBacth struct {
	BatchKey          uint64  `db:"batch_key"             json:"batch_key"`
	BatchNumber       *uint64 `db:"batch_number"          json:"batch_number"`
	BatchDisplayNo    *string `db:"batch_display_no"      json:"batch_display_no"`
	NavDate           string  `db:"nav_date"              json:"nav_date"`
	ProductKey        uint64  `db:"product_key"           json:"product_key"`
	TransTypeKey      uint64  `db:"trans_type_key"        json:"trans_type_key"`
	RecOrder          *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type CheckTrTransactionBacth struct {
	BatchKey uint64 `db:"batch_key"             json:"batch_key"`
}

func CreateTrTransactionBacth(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_transaction_batch"
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
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func CheckTrTransactionBatchByTransKey(c *CheckTrTransactionBacth, transactionKey string) (int, error) {
	query := `SELECT
              b.batch_key AS batch_key 
			  FROM tr_transaction_batch AS b 
			  LEFT JOIN tr_transaction AS t ON t.nav_date = b.nav_date AND t.product_key = b.product_key AND t.trans_type_key = b.trans_type_key AND t.rec_status = 1 
			  LEFT JOIN ms_product AS p ON p.product_key = t.product_key
			  WHERE t.transaction_key = '` + transactionKey + `' LIMIT 1`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
