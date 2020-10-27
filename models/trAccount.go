package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TrAccount struct {
	AccKey                 uint64  `db:"acc_key"                    json:"acc_key"`
	ProductKey             uint64  `db:"product_key"                json:"product_key"`
	CustomerKey            uint64  `db:"customer_key"               json:"customer_key"`
	AccountName            *string `db:"account_name"               json:"account_name"`
	AccountNo              *string `db:"account_no"                 json:"account_no"`
	IfuaNo                 *string `db:"ifua_no"                    json:"ifua_no"`
	IfuaName               *string `db:"ifua_name"                  json:"ifua_name"`
	AccStatus              *uint64 `db:"acc_status"                 json:"acc_status"`
	SubSuspendFlag         *uint8  `db:"sub_suspend_flag"           json:"sub_suspend_flag"`
	SubSuspendModifiedDate *string `db:"sub_suspend_modified_date"  json:"sub_suspend_modified_date"`
	SubSuspendReason       *string `db:"sub_suspend_reason"         json:"sub_suspend_reason"`
	SubSuspendReference    *string `db:"sub_suspend_reference"      json:"sub_suspend_reference"`
	RedSuspendFlag         *uint8  `db:"red_suspend_flag"           json:"red_suspend_flag"`
	RedSuspendModifiedDate *string `db:"red_suspend_modified_date"  json:"red_suspend_modified_date"`
	RedSuspendReason       *string `db:"red_suspend_reason"         json:"red_suspend_reason"`
	RedSuspendReference    *string `db:"red_suspend_reference"      json:"red_suspend_reference"`
	RecOrder               *uint64 `db:"rec_order"                  json:"rec_order"`
	RecStatus              uint8   `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate         *string `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy           *string `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate        *string `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy          *string `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1              *string `db:"rec_image1"                 json:"rec_image1"`
	RecImage2              *string `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus      *uint8  `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage       *uint64 `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate        *string `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy          *string `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate         *string `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy           *string `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1        *string `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2        *string `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3        *string `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
}

func CreateTrAccount(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_account"
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

func GetAllTrAccount(c *[]TrAccount, params map[string]string) (int, error) {
	query := `SELECT
              tr_account.* FROM 
			  tr_account`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_account."+field+" = '"+value+"'")
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

func GetTrAccountIn(c *[]TrAccount, value []string, field string, groupBy *string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_account.* FROM 
				tr_account WHERE 
				tr_account.rec_status = 1 `
	query := query2 + " AND tr_account." + field + " IN(" + inQuery + ")"

	if groupBy != nil {
		query += " AND tr_account.ifua_no is not null"
		query += " GROUP BY tr_account." + *groupBy
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
