package models

import (
	"api/db"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type TrTransactionConfirmation struct {
	TcKey               uint64           `db:"tc_key"                    json:"tc_key"`
	ConfirmDate         string           `db:"confirm_date"              json:"confirm_date"`
	TransactionKey      uint64           `db:"transaction_key"           json:"transaction_key"`
	ConfirmedAmount     decimal.Decimal  `db:"confirmed_amount"          json:"confirmed_amount"`
	ConfirmedUnit       decimal.Decimal  `db:"confirmed_unit"            json:"confirmed_unit"`
	AvgNav              *decimal.Decimal `db:"avg_nav"                   json:"avg_nav"`
	ConfirmResult       *uint64          `db:"confirm_result"            json:"confirm_result"`
	ConfirmedAmountDiff decimal.Decimal  `db:"confirmed_amount_diff"     json:"confirmed_amount_diff"`
	ConfirmedUnitDiff   decimal.Decimal  `db:"confirmed_unit_diff"       json:"confirmed_unit_diff"`
	ConfirmedRemarks    *string          `db:"confirmed_remarks"         json:"confirmed_remarks"`
	ConfirmedReferences *string          `db:"confirmed_references"      json:"confirmed_references"`
	RecOrder            *uint64          `db:"rec_order"                 json:"rec_order"`
	RecStatus           uint8            `db:"rec_status"                json:"rec_status"`
	RecCreatedDate      *string          `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy        *string          `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate     *string          `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy       *string          `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1           *string          `db:"rec_image1"                json:"rec_image1"`
	RecImage2           *string          `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus   *uint8           `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage    *uint64          `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate     *string          `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy       *string          `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate      *string          `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy        *string          `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1     *string          `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2     *string          `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3     *string          `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type TrTransactionConfirmationInfo struct {
	TcKey               uint64          `db:"tc_key"                    json:"tc_key"`
	ConfirmDate         string          `db:"confirm_date"              json:"confirm_date"`
	ConfirmedAmount     decimal.Decimal `db:"confirmed_amount"          json:"confirmed_amount"`
	ConfirmedUnit       decimal.Decimal `db:"confirmed_unit"            json:"confirmed_unit"`
	ConfirmedAmountDiff decimal.Decimal `db:"confirmed_amount_diff"     json:"confirmed_amount_diff"`
	ConfirmedUnitDiff   decimal.Decimal `db:"confirmed_unit_diff"       json:"confirmed_unit_diff"`
}

func GetTrTransactionConfirmation(c *TrTransactionConfirmation, key string) (int, error) {
	query := `SELECT tr_transaction_confirmation.* FROM tr_transaction_confirmation 
	          WHERE tr_transaction_confirmation.rec_status = 1 AND tr_transaction_confirmation.tc_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionConfirmationByTransactionKey(c *TrTransactionConfirmation, transactionKey string) (int, error) {
	query := `SELECT tr_transaction_confirmation.* FROM tr_transaction_confirmation 
	          WHERE tr_transaction_confirmation.rec_status = 1 AND tr_transaction_confirmation.transaction_key = ` + transactionKey
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateTrTransactionConfirmation(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_transaction_confirmation"
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
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func GetTrTransactionConfirmationIn(c *[]TrTransactionConfirmation, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
			tr_transaction_confirmation.* FROM 
			tr_transaction_confirmation `
	query := query2 + " WHERE tr_transaction_confirmation." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateTrTransactionConfirmation(params map[string]string) (int, error) {
	query := "UPDATE tr_transaction_confirmation SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "tc_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE tc_key = " + params["tc_key"]
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
