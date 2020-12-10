package models

import (
	"api/db"
	"database/sql"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type TrTransactionFifo struct {
	TransFifoKey      uint64   `db:"trans_fifo_key"            json:"trans_fifo_key"`
	TransConfRedKey   *uint64  `db:"trans_conf_red_key"        json:"trans_conf_red_key"`
	TransConfSubKey   uint64   `db:"trans_conf_sub_key"        json:"trans_conf_sub_key"`
	SubAcaKey         *uint64  `db:"sub_aca_key"               json:"sub_aca_key"`
	HoldingDays       *uint64  `db:"holding_days"              json:"holding_days"`
	TransUnit         *decimal.Decimal `db:"trans_unit"                json:"trans_unit"`
	FeeNavMode        *uint64  `db:"fee_nav_mode"              json:"fee_nav_mode"`
	TransAmount       *decimal.Decimal `db:"trans_amount"              json:"trans_amount"`
	TransFeeAmount    *decimal.Decimal `db:"trans_fee_amount"          json:"trans_fee_amount"`
	TransFeeTax       *decimal.Decimal `db:"trans_fee_tax"             json:"trans_fee_tax"`
	TransNettAmount   *decimal.Decimal `db:"trans_nett_amount"         json:"trans_nett_amount"`
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

func CreateTrTransactionFifo(params map[string]string) (int, error) {
	query := "INSERT INTO tr_transaction_fifo"
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

func UpdateTrTransactionFifo(params map[string]string, value string, field string) (int, error) {
	query := "UPDATE tr_transaction_fifo SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE " + field + " = " + value
	// log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		// log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		// log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
