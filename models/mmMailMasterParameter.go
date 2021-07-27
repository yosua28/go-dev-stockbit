package models

import (
	"api/db"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type MmMailMasterParameter struct {
	MailParameterKey      uint64  `db:"mail_parameter_key"            json:"mail_parameter_key"`
	MailMasterKey         uint64  `db:"mail_master_key"               json:"mail_master_key"`
	MailParamCode         *string `db:"mail_param_code"               json:"mail_param_code"`
	MailParamName         *string `db:"mail_param_name"               json:"mail_param_name"`
	MailParamDescription  *string `db:"mail_param_description"        json:"mail_param_description"`
	MailParamDefaultValue *string `db:"mail_param_default_value"      json:"mail_param_default_value"`
	MailParamDataType     *string `db:"mail_param_data_type"          json:"mail_param_data_type"`
	MailParamColumnName   *string `db:"mail_param_column_name"        json:"mail_param_column_name"`
	RecOrder              *uint64 `db:"rec_order"                     json:"rec_order"`
	RecStatus             uint8   `db:"rec_status"                    json:"rec_status"`
	RecCreatedDate        *string `db:"rec_created_date"              json:"rec_created_date"`
	RecCreatedBy          *string `db:"rec_created_by"                json:"rec_created_by"`
	RecModifiedDate       *string `db:"rec_modified_date"             json:"rec_modified_date"`
	RecModifiedBy         *string `db:"rec_modified_by"               json:"rec_modified_by"`
	RecImage1             *string `db:"rec_image1"                    json:"rec_image1"`
	RecImage2             *string `db:"rec_image2"                    json:"rec_image2"`
	RecApprovalStatus     *uint8  `db:"rec_approval_status"           json:"rec_approval_status"`
	RecApprovalStage      *uint64 `db:"rec_approval_stage"            json:"rec_approval_stage"`
	RecApprovedDate       *string `db:"rec_approved_date"             json:"rec_approved_date"`
	RecApprovedBy         *string `db:"rec_approved_by"               json:"rec_approved_by"`
	RecDeletedDate        *string `db:"rec_deleted_date"              json:"rec_deleted_date"`
	RecDeletedBy          *string `db:"rec_deleted_by"                json:"rec_deleted_by"`
	RecAttributeID1       *string `db:"rec_attribute_id1"             json:"rec_attribute_id1"`
	RecAttributeID2       *string `db:"rec_attribute_id2"             json:"rec_attribute_id2"`
	RecAttributeID3       *string `db:"rec_attribute_id3"             json:"rec_attribute_id3"`
}

func GetMmMailMasterParameter(c *MmMailMaster, field string, key string) (int, error) {
	query := `SELECT mm_mail_master_parameter.* FROM mm_mail_master_parameter 
	WHERE mm_mail_master_parameter.rec_status = '1' 
	AND mm_mail_master_parameter.` + field + ` = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateMmMailMasterParameter(params map[string]string) (int, error) {
	query := "UPDATE mm_mail_master_parameter SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "mail_parameter_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE mail_parameter_key = " + params["mail_parameter_key"]
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

func CreateMultipleMmMailMasterParamenter(params []interface{}) (int, error) {

	q := `INSERT INTO mm_mail_master_parameter (
		mail_master_key, 
		mail_param_code,
		rec_status,
		rec_created_date,
		rec_created_by) VALUES `

	for i := 0; i < len(params); i++ {
		q += "(?)"
		if i < (len(params) - 1) {
			q += ","
		}
	}
	query, args, err := sqlx.In(q, params...)
	if err != nil {
		return http.StatusBadGateway, err
	}

	query = db.Db.Rebind(query)
	_, err = db.Db.Query(query, args...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}

func GetAllMmMailParametergent(c *[]MmMailMasterParameter, params map[string]string) (int, error) {
	query := `SELECT
			mm_mail_master_parameter.* FROM 
			mm_mail_master_parameter`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "mm_mail_master_parameter."+field+" = '"+value+"'")
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
		log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateDeleteAllParameter(field string, params map[string]string, mailParameterKeyValues []string) (int, error) {
	inQuery := strings.Join(mailParameterKeyValues, ",")
	query := "UPDATE mm_mail_master_parameter SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE " + field + " IN(" + inQuery + ")"
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
