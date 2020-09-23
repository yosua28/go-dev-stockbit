package models

import (
	"api/db"
	"net/http"
	"database/sql"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MsBankAccount struct {
	BankAccountKey       uint64    `db:"bank_account_key"       json:"bank_account_key"`
	BankKey              uint64    `db:"bank_key"               json:"bank_key"`
	AccountNo            string    `db:"account_no"             json:"account_no"`
	AccountHolderName    string    `db:"account_holder_name"    json:"account_holder_name"`
	BranchName           *string   `db:"branch_name"            json:"branch_name"`
	CurrencyKey          uint64    `db:"currency_key"           json:"currency_key"`
	BankAccountType      uint64    `db:"bank_account_type"      json:"bank_account_type"`
	RecDomain            uint64    `db:"rec_domain"             json:"rec_domain"`
	SwiftCode            *string   `db:"swift_code"             json:"swift_code"`
	RecOrder             *uint64   `db:"rec_order"              json:"rec_order"`
	RecStatus            uint8     `db:"rec_status"             json:"rec_status"`
	RecCreatedDate       *string   `db:"rec_created_date"       json:"rec_created_date"`
	RecCreatedBy         *string   `db:"rec_created_by"         json:"rec_created_by"`
	RecModifiedDate      *string   `db:"rec_modified_date"      json:"rec_modified_date"`
	RecModifiedBy        *string   `db:"rec_modified_by"        json:"rec_modified_by"`
	RecImage1            *string   `db:"rec_image1"             json:"rec_image1"`
	RecImage2            *string   `db:"rec_image2"             json:"rec_image2"`
	RecApprovalStatus    *uint8    `db:"rec_approval_status"    json:"rec_approval_status"`
	RecApprovalStage     *uint64   `db:"rec_approval_stage"     json:"rec_approval_stage"`
	RecApprovedDate      *string   `db:"rec_approved_date"      json:"rec_approved_date"`
	RecApprovedBy        *string   `db:"rec_approved_by"        json:"rec_approved_by"`
	RecDeletedDate       *string   `db:"rec_deleted_date"       json:"rec_deleted_date"`
	RecDeletedBy         *string   `db:"rec_deleted_by"         json:"rec_deleted_by"`
	RecAttributeID1      *string   `db:"rec_attribute_id1"      json:"rec_attribute_id1"`
	RecAttributeID2      *string   `db:"rec_attribute_id2"      json:"rec_attribute_id2"`
	RecAttributeID3      *string   `db:"rec_attribute_id3"      json:"rec_attribute_id3"`
}

func CreateMsBankAccount(params map[string]string) (int, error, string){
	query := "INSERT INTO ms_bank_account"
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
		query += "("+fields + ") VALUES(" + values + ")"
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