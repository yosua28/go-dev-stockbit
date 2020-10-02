package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type OaRequest struct {
	OaRequestKey      uint64  `db:"oa_request_key"             json:"oa_request_key"`
	OaRequestType     *uint64 `db:"oa_request_type"            json:"oa_request_type"`
	OaEntryStart      string  `db:"oa_entry_start"             json:"oa_entry_start"`
	OaEntryEnd        string  `db:"oa_entry_end"               json:"oa_entry_end"`
	Oastatus          *uint64 `db:"oa_status"                  json:"oa_status"`
	UserLoginKey      *uint64 `db:"user_login_key"             json:"user_login_key"`
	CustomerKey       *uint64 `db:"customer_key"               json:"customer_key"`
	SalesCode         *string `db:"sales_code"                 json:"sales_code"`
	Check1Date        *string `db:"check1_date"                json:"check1_date"`
	Check1Flag        *uint8  `db:"check1_flag"                json:"check1_flag"`
	Check1References  *string `db:"check1_references"          json:"check1_references"`
	Check1Notes       *string `db:"check1_notes"               json:"check1_notes"`
	Check2Date        *string `db:"check2_date"                json:"check2_date"`
	Check2Flag        *uint8  `db:"check2_flag"                json:"check2_flag"`
	Check2References  *string `db:"check2_references"          json:"check2_references"`
	Check2Notes       *string `db:"check2_notes"               json:"check2_notes"`
	OaRiskLevel       *uint64 `db:"oa_risk_level"              json:"oa_risk_level"`
	RecOrder          *uint64 `db:"rec_order"                  json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                 json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
}

type OaRequestDataResponse struct {
	OaRequestKey      uint64  `json:"oa_request_key"`
	OaRequestType     *string `json:"oa_request_type"`
	OaEntryStart      string  `json:"oa_entry_start"`
	OaEntryEnd        string  `json:"oa_entry_end"`
	Oastatus          string  `json:"oa_status"`
	UserLoginName     *string `json:"user_login_name"`
	UserLoginFullName *string `json:"user_login_full_name"`
	Customer          *string `json:"customer"`
	SalesCode         *string `json:"sales_code"`
	Check1Date        *string `json:"check1_date"`
	Check1Flag        *uint8  `json:"check1_flag"`
	Check1References  *string `json:"check1_references"`
	Check1Notes       *string `json:"check1_notes"`
	Check2Date        *string `json:"check2_date"`
	Check2Flag        *uint8  `json:"check2_flag"`
	Check2References  *string `json:"check2_references"`
	Check2Notes       *string `json:"check2_notes"`
	OaRiskLevel       *string `json:"oa_risk_level"`
	RecOrder          *uint64 `json:"rec_order"`
	RecStatus         uint8   `json:"rec_status"`
	RecCreatedDate    *string `json:"rec_created_date"`
	RecCreatedBy      *string `json:"rec_created_by"`
	RecModifiedDate   *string `json:"rec_modified_date"`
	RecModifiedBy     *string `json:"rec_modified_by"`
	RecImage1         *string `json:"rec_image1"`
	RecImage2         *string `json:"rec_image2"`
	RecApprovalStatus *uint8  `json:"rec_approval_status"`
	RecApprovalStage  *uint64 `json:"rec_approval_stage"`
	RecApprovedDate   *string `json:"rec_approved_date"`
	RecApprovedBy     *string `json:"rec_approved_by"`
	RecDeletedDate    *string `json:"rec_deleted_date"`
	RecDeletedBy      *string `json:"rec_deleted_by"`
	RecAttributeID1   *string `json:"rec_attribute_id1"`
	RecAttributeID2   *string `json:"rec_attribute_id2"`
	RecAttributeID3   *string `json:"rec_attribute_id3"`
}

type OaRequestListResponse struct {
	OaRequestKey      uint64  `json:"oa_request_key"`
	OaRequestType     *string `json:"oa_request_type"`
	OaEntryStart      string  `json:"oa_entry_start"`
	OaEntryEnd        string  `json:"oa_entry_end"`
	Oastatus          string  `json:"oa_status"`
	UserLoginName     *string `json:"user_login_name"`
	UserLoginFullName *string `json:"user_login_full_name"`
	Customer          *string `json:"customer"`
	SalesCode         *string `json:"sales_code"`
	Check1Date        *string `json:"check1_date"`
	Check1Flag        *uint8  `json:"check1_flag"`
	Check1References  *string `json:"check1_references"`
	Check1Notes       *string `json:"check1_notes"`
	Check2Date        *string `json:"check2_date"`
	Check2Flag        *uint8  `json:"check2_flag"`
	Check2References  *string `json:"check2_references"`
	Check2Notes       *string `json:"check2_notes"`
	OaRiskLevel       *string `json:"oa_risk_level"`
	RecOrder          *uint64 `json:"rec_order"`
	RecStatus         uint8   `json:"rec_status"`
}

type OaRequestCountData struct {
	CountData int `db:"count_data"             json:"count_data"`
}

func CreateOaRequest(params map[string]string) (int, error, string) {
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

func GetAllOaRequest(c *[]OaRequest, limit uint64, offset uint64, nolimit bool, params map[string]string, statusFilter uint64) (int, error) {
	query := `SELECT
              oa_request.*
			  FROM oa_request`
	var present bool
	// var whereClause []string
	var condition string

	// Check status by
	if statusFilter > 0 {
		condition += " WHERE oa_request.oa_status = " + strconv.FormatUint(statusFilter, 10)
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

func GetOaRequest(c *OaRequest, key string) (int, error) {
	query := `SELECT oa_request.* FROM oa_request WHERE oa_request.oa_request_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetCountOaRequest(c *OaRequestCountData, statusFilter uint64) (int, error) {
	query := `SELECT
              count(oa_request.oa_request_key) as count_data
			  FROM oa_request`

	if statusFilter > 0 {
		query += " WHERE oa_request.oa_status = " + strconv.FormatUint(statusFilter, 10)
	}

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
