package models

import (
	"api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsBankList struct {
	BankKey      uint64 `json:"bank_key"`
	BankCode     string `json:"bank_code"`
	BankName     string `json:"bank_name"`
	BankFullname string `json:"bank_fullname"`
}

type MsBank struct {
	BankKey           uint64  `db:"bank_key"               json:"bank_key"`
	BankCode          string  `db:"bank_code"              json:"bank_code"`
	BankName          string  `db:"bank_name"              json:"bank_name"`
	BankFullname      *string `db:"bank_fullname"          json:"bank_fullname"`
	BiMemberCode      *string `db:"bi_member_code"         json:"bi_member_code"`
	SwiftCode         *string `db:"swift_code"             json:"swift_code"`
	FlagLocal         uint8   `db:"flag_local"             json:"flag_local"`
	FlagGoverment     uint8   `db:"flag_government"        json:"flag_government"`
	BankLogo          *string `db:"bank_logo"              json:"bank_logo"`
	BankWebUrl        *string `db:"bank_web_url"           json:"bank_web_url"`
	BankIbankUrl      *string `db:"bank_ibank_url"         json:"bank_ibank_url"`
	RecOrder          *uint64 `db:"rec_order"              json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"             json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"       json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"         json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"      json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"        json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"             json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"             json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"    json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"     json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"      json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"        json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"       json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"         json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"      json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"      json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"      json:"rec_attribute_id3"`
}

func GetAllMsBank(c *[]MsBank, params map[string]string) (int, error) {
	query := `SELECT
              ms_bank.* FROM 
			  ms_bank `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_bank."+field+" = '"+value+"'")
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

func GetMsBank(c *MsBank, key string) (int, error) {
	query := `SELECT ms_bank.* FROM ms_bank WHERE ms_bank.bank_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsBankIn(c *[]MsBank, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_bank.* FROM 
				ms_bank `
	query := query2 + " WHERE ms_bank." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type ListBankAdmin struct {
	BankKey        uint64  `db:"bank_key"        json:"bank_key"`
	BankCode       string  `db:"bank_code"       json:"bank_code"`
	BankName       string  `db:"bank_name"       json:"bank_name"`
	BankFullname   *string `db:"bank_fullname"   json:"bank_fullname"`
	BiMemberCode   *string `db:"bi_member_code"  json:"bi_member_code"`
	SwiftCode      *string `db:"swift_code"      json:"swift_code"`
	BankLocal      string  `db:"bank_local"      json:"bank_local"`
	BankGovernment string  `db:"bank_government" json:"bank_government"`
}

func AdminGetListBank(c *[]ListBankAdmin, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
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
		condition += " (bank_code like '%" + searchLike + "%' OR"
		condition += " bank_name like '%" + searchLike + "%' OR"
		condition += " bank_fullname like '%" + searchLike + "%' OR"
		condition += " bi_member_code like '%" + searchLike + "%' OR"
		condition += " swift_code like '%" + searchLike + "%')"
	}

	query := `SELECT
				bank_key,
				bank_code,
				bank_name,
				bank_fullname,
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
			FROM ms_bank 
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

func CountAdminGetListBank(c *CountData, params map[string]string, searchLike string) (int, error) {
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
		condition += " (bank_code like '%" + searchLike + "%' OR"
		condition += " bank_name like '%" + searchLike + "%' OR"
		condition += " bank_fullname like '%" + searchLike + "%' OR"
		condition += " bi_member_code like '%" + searchLike + "%' OR"
		condition += " swift_code like '%" + searchLike + "%')"
	}

	query := `SELECT
				count(bank_key) AS count_data 
			FROM ms_bank 
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

func UpdateMsBank(params map[string]string) (int, error) {
	query := "UPDATE ms_bank SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "bank_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE bank_key = " + params["bank_key"]
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

func CreateMsBank(params map[string]string) (int, error) {
	query := "INSERT INTO ms_bank"
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

func CountMsBankValidateUnique(c *CountData, field string, value string, key string) (int, error) {
	query := `SELECT 
				COUNT(bank_key) AS count_data 
			FROM ms_bank
			WHERE rec_status = '1' AND ` + field + ` = '` + value + `'`

	if key != "" {
		query += " AND bank_key != '" + key + "'"
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
