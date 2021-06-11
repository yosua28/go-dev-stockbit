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

type TrAccountAdmin struct {
	AccKey           uint64  `db:"acc_key"                    json:"acc_key"`
	ProductKey       uint64  `db:"product_key"                json:"product_key"`
	CustomerKey      uint64  `db:"customer_key"               json:"customer_key"`
	IfuaNo           *string `db:"ifua_no"                    json:"ifua_no"`
	ProductName      *string `db:"product_name"               json:"product_name"`
	Cif              *string `db:"cif"                        json:"cif"`
	FullName         *string `db:"full_name"                  json:"full_name"`
	Sid              *string `db:"sid"                        sjson:"sid"`
	AccStatus        *uint64 `db:"acc_status"                 json:"acc_status"`
	SubSuspendFlag   *uint8  `db:"sub_suspend_flag"           json:"sub_suspend_flag"`
	SubSuspendReason *string `db:"sub_suspend_reason"         json:"sub_suspend_reason"`
	RedSuspendFlag   *uint8  `db:"red_suspend_flag"           json:"red_suspend_flag"`
	RedSuspendReason *string `db:"red_suspend_reason"         json:"red_suspend_reason"`
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

func UpdateTrAccountUploadSinvest(params map[string]string, value string, field string) (int, error) {
	query := "UPDATE tr_account SET "
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
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func AdminGetAllTrAccount(c *[]TrAccountAdmin, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
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

	query := `SELECT 
				a.acc_key AS acc_key,
				a.product_key AS product_key,
				a.customer_key AS customer_key,
				a.ifua_no AS ifua_no,
				p.product_name_alt AS product_name,
				c.unit_holder_idno AS cif,
				c.full_name AS full_name,
				c.sid_no AS sid,
				(CASE
					WHEN a.sub_suspend_flag = 1 THEN 1
					ELSE 0
				END) AS sub_suspend_flag, 
				(CASE
					WHEN a.red_suspend_flag = 1 THEN 1
					ELSE 0
				END) AS red_suspend_flag, 
				a.sub_suspend_reason AS sub_suspend_reason,
				a.red_suspend_reason AS red_suspend_reason 
				FROM tr_account AS a 
			INNER JOIN ms_product AS p ON a.product_key = p.product_key
			INNER JOIN ms_customer AS c ON c.customer_key = a.customer_key
			WHERE a.rec_status = 1 AND c.rec_status = 1 AND p.rec_status = 1`

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}

	if !nolimit {
		condition += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			condition += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += condition

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CountAdminGetAllTrAccount(c *CountData, params map[string]string) (int, error) {
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

	query := `SELECT 
				COUNT(a.acc_key) AS count_data 
			FROM tr_account AS a
			INNER JOIN ms_product AS p ON a.product_key = p.product_key
			INNER JOIN ms_customer AS c ON c.customer_key = a.customer_key
			WHERE a.rec_status = 1 AND c.rec_status = 1 AND p.rec_status = 1`

	query += condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetDetailTrAccount(c *TrAccountAdmin, accKey string) (int, error) {
	query := `SELECT 
				a.acc_key AS acc_key,
				a.product_key AS product_key,
				a.customer_key AS customer_key,
				a.ifua_no AS ifua_no,
				p.product_name_alt AS product_name,
				c.unit_holder_idno AS cif,
				c.full_name AS full_name,
				c.sid_no AS sid,
				(CASE
					WHEN a.sub_suspend_flag = 1 THEN 1
					ELSE 0
				END) AS sub_suspend_flag, 
				(CASE
					WHEN a.red_suspend_flag = 1 THEN 1
					ELSE 0
				END) AS red_suspend_flag, 
				a.sub_suspend_reason AS sub_suspend_reason,
				a.red_suspend_reason AS red_suspend_reason 
				FROM tr_account AS a 
			INNER JOIN ms_product AS p ON a.product_key = p.product_key
			INNER JOIN ms_customer AS c ON c.customer_key = a.customer_key
			WHERE a.rec_status = 1 AND c.rec_status = 1 AND p.rec_status = 1 AND a.acc_key = '` + accKey + `'`

	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateTrAccountByProductAndCustomer(params map[string]string, productKey string, customerKey string) (int, error) {
	query := "UPDATE tr_account SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE product_key = " + productKey + " AND customer_key = " + customerKey
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetCustomerAccountByProduct(c *[]CustomerDropdown, productKey string) (int, error) {

	query := `SELECT 
				c.customer_key as customer_key,
				CONCAT(c.unit_holder_idno, " - ", c.full_name) AS name,
				c.openacc_branch_key as branch_key,
				c.openacc_agent_key as agent_key 
			FROM tr_account as a 
			INNER JOIN ms_customer AS c ON c.customer_key = a.customer_key
			INNER JOIN sc_user_login AS l ON l.customer_key = c.customer_key 
			WHERE c.rec_status = 1 AND l.rec_status = 1 AND a.rec_status = 1 AND a.product_key = '` + productKey + `' 
			GROUP BY a.customer_key ORDER BY c.full_name ASC `

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
