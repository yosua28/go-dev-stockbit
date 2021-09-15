package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type TrTransactionSettlement struct {
	SettlementKey        uint64          `db:"settlement_key"            json:"settlement_key"`
	TransactionKey       *uint64         `db:"transaction_key"           json:"transaction_key"`
	SettlePurposed       uint64          `db:"settle_purposed"           json:"settle_purposed"`
	SettleDate           string          `db:"settle_date"               json:"settle_date"`
	SettleNominal        decimal.Decimal `db:"settle_nominal"            json:"settle_nominal"`
	ClientSubaccountNo   string          `db:"client_subaccount_no"      json:"client_subaccount_no"`
	SettleStatus         uint64          `db:"settled_status"            json:"settled_status"`
	SettleRealizedDate   *string         `db:"settle_realized_date"      json:"settle_realized_date"`
	SettleRemarks        *string         `db:"settle_remarks"            json:"settle_remarks"`
	SettleReference      *string         `db:"settle_references"         json:"settle_references"`
	SourceBankAccountKey *uint64         `db:"source_bank_account_key"   json:"source_bank_account_key"`
	TargetBankAccountKey uint64          `db:"target_bank_account_key"   json:"target_bank_account_key"`
	SettleNotes          *string         `db:"settle_notes"              json:"settle_notes"`
	SettleChannel        uint64          `db:"settle_channel"            json:"settle_channel"`
	SettlePaymentMethod  uint64          `db:"settle_payment_method"     json:"settle_payment_method"`
	SettleResponse       *string         `db:"settle_response"           json:"settle_response"`
	RecOrder             *uint64         `db:"rec_order"                 json:"rec_order"`
	RecStatus            uint8           `db:"rec_status"                json:"rec_status"`
	RecCreatedDate       *string         `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy         *string         `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate      *string         `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy        *string         `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1            *string         `db:"rec_image1"                json:"rec_image1"`
	RecImage2            *string         `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus    *uint8          `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage     *uint64         `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate      *string         `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy        *string         `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate       *string         `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy         *string         `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1      *string         `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2      *string         `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3      *string         `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type TransactionSettlement struct {
	SettlementKey       uint64          `json:"settlement_key"`
	SettlePurposed      string          `json:"settle_purposed"`
	SettleDate          string          `json:"settle_date"`
	SettleNominal       decimal.Decimal `json:"settle_nominal"`
	SettleStatus        string          `json:"settle_status"`
	SettleRealizedDate  string          `json:"settle_realized_date"`
	SettleRemarks       *string         `json:"settle_remarks"`
	SettleReference     *string         `json:"settle_references"`
	SettleChannel       string          `json:"settle_channel"`
	SettlePaymentMethod string          `json:"settle_payment_method"`
}

func GetAllTrTransactionSettlement(c *[]TrTransactionSettlement, params map[string]string) (int, error) {
	query := `SELECT
              tr_transaction_settlement.* FROM 
			  tr_transaction_settlement `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_transaction_settlement."+field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
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

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionSettlementIn(c *[]TrTransactionSettlement, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_transaction_settlement.* FROM 
				tr_transaction_settlement `
	query := query2 + " WHERE tr_transaction_settlement." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionSettlement(c *TrTransactionSettlement, field string, key string) (int, error) {
	query := `SELECT tr_transaction_settlement.* FROM tr_transaction_settlement WHERE 
	tr_transaction_settlement.rec_status = "1" AND tr_transaction_settlement.` + field + ` = "` + key + `"`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateTrTransactionSettlement(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_transaction_settlement"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		if value == "NULL" {
			var s *string
			bindvars = append(bindvars, s)
		} else {
			bindvars = append(bindvars, value)
		}

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

func UpdateTrTransactionSettlement(params map[string]string) (int, error) {
	query := "UPDATE tr_transaction_settlement SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "settlement_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE settlement_key = " + params["settlement_key"]
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
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
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
