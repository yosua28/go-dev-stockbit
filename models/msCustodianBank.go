package models

import (
	"api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsCustodianBankInfo struct {
	CustodianCode      string  `json:"custodian_code"`
	CustodianShortName string  `json:"custodian_short_name"`
	CustodianFullName  *string `json:"custodian_full_name"`
}

type MsCustodianBankInfoList struct {
	CustodianKey       uint64  `json:"custodian_key"`
	CustodianCode      string  `json:"custodian_code"`
	CustodianShortName string  `json:"custodian_short_name"`
	CustodianFullName  *string `json:"custodian_full_name"`
}

type MsCustodianBank struct {
	CustodianKey       uint64  `db:"custodian_key"         json:"custodian_key"`
	CustodianCode      string  `db:"custodian_code"        json:"custodian_code"`
	CustodianShortName string  `db:"custodian_short_name"  json:"custodian_short_name"`
	CustodianFullName  *string `db:"custodian_full_name"   json:"custodian_full_name"`
	BiMemberCode       *string `db:"bi_member_code"        json:"bi_member_code"`
	SwiftCode          *string `db:"swift_code"            json:"swift_code"`
	FlagLocal          uint8   `db:"flag_local"            json:"flag_local"`
	FlagGoverment      uint8   `db:"flag_government"       json:"flag_government"`
	BankWebUrl         *string `db:"bank_web_url"          json:"bank_web_url"`
	BankLogo           *string `db:"bank_logo"             json:"bank_logo"`
	CustodianProfile   *string `db:"custodian_profile"     json:"custodian_profile"`
	RecOrder           *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

func GetMsCustodianBank(c *MsCustodianBank, key string) (int, error) {
	query := `SELECT ms_custodian_bank.* FROM ms_custodian_bank WHERE ms_custodian_bank.rec_status = '1' 
	AND ms_custodian_bank.custodian_key = ` + key
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsCustodianBankIn(c *[]MsCustodianBank, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_custodian_bank.* FROM 
				ms_custodian_bank WHERE 
				ms_custodian_bank.rec_status = 1 `
	query := query2 + " AND ms_custodian_bank." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetListMsCustodianBank(c *[]MsCustodianBank, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query2 := `SELECT
				ms_custodian_bank.* FROM 
				ms_custodian_bank `
	query := query2 + " WHERE ms_custodian_bank.rec_status = 1 "

	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_custodian_bank."+field+" = '"+value+"'")
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

type ListCustodianBankAdmin struct {
	CustodianKey       uint64  `db:"custodian_key"              json:"custodian_key"`
	CustodianCode      string  `db:"custodian_code"             json:"custodian_code"`
	CustodianShortName string  `db:"custodian_short_name"       json:"custodian_short_name"`
	CustodianFullName  *string `db:"custodian_full_name"        json:"custodian_full_name"`
	BiMemberCode       *string `db:"bi_member_code"             json:"bi_member_code"`
	SwiftCode          *string `db:"swift_code"                 json:"swift_code"`
	BankLocal          string  `db:"bank_local"                 json:"bank_local"`
	BankGovernment     string  `db:"bank_government"            json:"bank_government"`
}

func AdminGetListCustodianBank(c *[]ListCustodianBankAdmin, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

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

	if searchLike != "" {
		condition += " AND"
		condition += " (custodian_code like '%" + searchLike + "%' OR"
		condition += " custodian_short_name like '%" + searchLike + "%' OR"
		condition += " custodian_full_name like '%" + searchLike + "%' OR"
		condition += " bi_member_code like '%" + searchLike + "%' OR"
		condition += " swift_code like '%" + searchLike + "%')"
	}

	query := `SELECT
				custodian_key,
				custodian_code,
				custodian_short_name,
				custodian_full_name,
				bi_member_code,
				swift_code,
				(CASE
					WHEN flag_local = "1" THEN "Ya"
					ELSE "Tidak"
				END) AS bank_local, 
				(CASE
					WHEN flag_government = "1" THEN "Ya"
					ELSE "Tidak"
				END) AS bank_government 
			FROM ms_custodian_bank 
			WHERE rec_status = 1 ` + condition

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	if !nolimit {
		limitOffset += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			limitOffset += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += orderCondition + limitOffset

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CountAdminGetListCustodianBank(c *CountData, params map[string]string, searchLike string) (int, error) {
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

	if searchLike != "" {
		condition += " AND"
		condition += " (custodian_code like '%" + searchLike + "%' OR"
		condition += " custodian_short_name like '%" + searchLike + "%' OR"
		condition += " custodian_full_name like '%" + searchLike + "%' OR"
		condition += " bi_member_code like '%" + searchLike + "%' OR"
		condition += " swift_code like '%" + searchLike + "%')"
	}

	query := `SELECT 
				count(custodian_key) AS count_data 
			FROM ms_custodian_bank 
			WHERE rec_status = 1 ` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateMsCustodianBank(params map[string]string) (int, error) {
	query := "UPDATE ms_custodian_bank SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "custodian_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE custodian_key = " + params["custodian_key"]
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

func CreateMsCustodianBank(params map[string]string) (int, error) {
	query := "INSERT INTO ms_custodian_bank"
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
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func CountMsCustodianBankValidateUnique(c *CountData, field string, value string, key string) (int, error) {
	query := `SELECT 
				COUNT(custodian_key) AS count_data 
			FROM ms_custodian_bank
			WHERE ` + field + ` = '` + value + `'`

	if key != "" {
		query += " AND custodian_key != '" + key + "'"
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
