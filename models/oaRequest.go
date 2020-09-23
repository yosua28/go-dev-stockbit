package models

import (
	"api/db"
	"net/http"
	"database/sql"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type OaRequest struct {
	OaRequestKey           uint64    `db:"oa_request_key"             json:"oa_request_key"`
	OaRequestType          string    `db:"oa_request_type"            json:"oa_request_type"`
	OaEntryStart           string    `db:"oa_entry_start"             json:"oa_entry_start"`
	OaEntryEnd             string    `db:"oa_entry_end"               json:"oa_entry_end"`
	Oastatus               string    `db:"oa_status"                  json:"oa_status"`
	UserLoginKey          *uint64    `db:"user_login_key"             json:"user_login_key"`
	CustomerKey           *uint64    `db:"customer_key"               json:"customer_key"`
	SalesCode             *string    `db:"sales_code"                 json:"sales_code"`
	Check1Date            *string    `db:"check1_date"                json:"check1_date"`
	Check1Flag            *uint8     `db:"check1_flag"                json:"check1_flag"`
	Check1References      *string    `db:"check1_references"          json:"check1_references"`
	Check1Notes           *string    `db:"check1_notes"               json:"check1_notes"`
	Check2Date            *string    `db:"check2_date"                json:"check2_date"`
	Check2Flag            *uint8     `db:"check2_flag"                json:"check2_flag"`
	Check2References      *string    `db:"check2_references"          json:"check2_references"`
	Check2Notes           *string    `db:"check2_notes"               json:"check2_notes"`
	TrxRiskLevel          *string    `db:"trx_risk_level"             json:"trx_risk_level"`
	RecOrder              *uint64    `db:"rec_order"                  json:"rec_order"`
	RecStatus              uint8     `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate        *string    `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy          *string    `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate       *string    `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy         *string    `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1             *string    `db:"rec_image1"                 json:"rec_image1"`
	RecImage2             *string    `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus     *uint8     `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage      *uint64    `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate       *string    `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy         *string    `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate        *string    `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy          *string    `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1       *string    `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2       *string    `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3       *string    `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
}

func CreateOaRequest(params map[string]string) (int, error, string){
	query := "INSERT INTO oa_request"
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