package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MmMailMaster struct {
	MailMasterKey      uint64  `db:"mail_master_key"               json:"mail_master_key"`
	MailMasterType     uint64  `db:"mail_master_type"              json:"mail_master_type"`
	MailMasterCategory *uint64 `db:"mail_master_category"          json:"mail_master_category"`
	MailTemplateName   string  `db:"mail_template_name"            json:"mail_template_name"`
	MailTemplateDesc   *string `db:"mail_template_desc"            json:"mail_template_desc"`
	MailAccountKey     *uint64 `db:"mail_account_key"              json:"mail_account_key"`
	MailToGroupKey     *uint64 `db:"mail_to_group_key"             json:"mail_to_group_key"`
	MailToEmailParam   *string `db:"mail_to_email_param"           json:"mail_to_email_param"`
	MailCcGroupKey     *uint64 `db:"mail_cc_group_key"             json:"mail_cc_group_key"`
	MailCcMailKey      *uint64 `db:"mail_cc_mail_key"              json:"mail_cc_mail_key"`
	MailCcEmailParam   *string `db:"mail_cc_email_param"           json:"mail_cc_email_param"`
	MailBccGroupKey    *uint64 `db:"mail_bcc_group_key"            json:"mail_bcc_group_key"`
	MailBccMailKey     *uint64 `db:"mail_bcc_mail_key"             json:"mail_bcc_mail_key"`
	MailBccEmailParam  *string `db:"mail_bcc_email_param"          json:"mail_bcc_email_param"`
	MailSubject        *string `db:"mail_subject"                  json:"mail_subject"`
	MailBody           *string `db:"mail_body"                     json:"mail_body"`
	MailImageFile      *string `db:"mail_image_file"               json:"mail_image_file"`
	MailHasAttachment  uint8   `db:"mail_has_attachment"           json:"mail_has_attachment"`
	MailFlagHtml       uint8   `db:"mail_flag_html"                json:"mail_flag_html"`
	RecOrder           *uint64 `db:"rec_order"                     json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"                    json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"              json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"                json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"             json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"               json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"                    json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"                    json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"           json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"            json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"             json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"               json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"              json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"                json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"             json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"             json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"             json:"rec_attribute_id3"`
}

func GetMmMailMaster(c *MmMailMaster, field string, key string) (int, error) {
	query := `SELECT mm_mail_master.* FROM mm_mail_master 
	WHERE mm_mail_master.rec_status = '1' 
	AND mm_mail_master.` + field + ` = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateMmMailMaster(params map[string]string) (int, error, string) {
	query := "INSERT INTO mm_mail_master"
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

func UpdateMmMailMaster(params map[string]string) (int, error) {
	query := "UPDATE mm_mail_master SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "mail_master_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE mail_master_key = " + params["mail_master_key"]
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

func CountMmMailMasterValidateUnique(c *CountData, field string, value string, key string) (int, error) {
	query := `SELECT 
				COUNT(mail_master_key) AS count_data 
			FROM mm_mail_master
			WHERE rec_status = '1' AND ` + field + ` = '` + value + `'`

	if key != "" {
		query += " AND mail_master_key != '" + key + "'"
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

type ListMmMailMaster struct {
	MailMasterKey      uint64  `db:"mail_master_key"            json:"mail_master_key"`
	MailMasterType     string  `db:"mail_master_type"           json:"mail_master_type"`
	MailMasterCategory *string `db:"mail_master_category"       json:"mail_master_category"`
	TemplateName       string  `db:"template_name"              json:"template_name"`
	Description        *string `db:"description"                json:"description"`
}

func AdminGetListMmMailMaster(c *[]ListMmMailMaster, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
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
		condition += " (ty.lkp_name like '%" + searchLike + "%' OR"
		condition += " ct.lkp_name like '%" + searchLike + "%' OR"
		condition += " m.mail_template_name like '%" + searchLike + "%' OR"
		condition += " m.mail_template_desc like '%" + searchLike + "%')"
	}

	query := `SELECT 
				m.mail_master_key,
				ty.lkp_name AS mail_master_type,
				ct.lkp_name AS mail_master_category,
				m.mail_template_name AS template_name,
				m.mail_template_desc AS description 
			FROM mm_mail_master AS m
			INNER JOIN gen_lookup AS ty ON ty.lookup_key = m.mail_master_type
			LEFT JOIN gen_lookup AS ct ON ct.lookup_key = m.mail_master_category
			WHERE m.rec_status = 1` + condition

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

func CountAdminGetMmMailMaster(c *CountData, params map[string]string, searchLike string) (int, error) {
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
		condition += " (ty.lkp_name like '%" + searchLike + "%' OR"
		condition += " ct.lkp_name like '%" + searchLike + "%' OR"
		condition += " m.mail_template_name like '%" + searchLike + "%' OR"
		condition += " m.mail_template_desc like '%" + searchLike + "%')"
	}

	query := `SELECT
				count(m.mail_master_key) AS count_data 
			FROM mm_mail_master AS m
			INNER JOIN gen_lookup AS ty ON ty.lookup_key = m.mail_master_type
			LEFT JOIN gen_lookup AS ct ON ct.lookup_key = m.mail_master_category
			WHERE m.rec_status = 1` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
