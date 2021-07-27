package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MmMailSentLog struct {
	MailSentKey        uint64  `db:"mail_sent_key"                 json:"mail_sent_key"`
	MailMasterKey      uint64  `db:"mail_master_key"               json:"mail_master_key"`
	MailMasterType     *uint64 `db:"mail_master_type"              json:"mail_master_type"`
	MailMasterCategory *uint64 `db:"mail_master_category"          json:"mail_master_category"`
	MailTemplateName   string  `db:"mail_template_name"            json:"mail_template_name"`
	MailAccountKey     *uint64 `db:"mail_account_key"              json:"mail_account_key"`
	MailToGroupKey     *uint64 `db:"mail_to_group_key"             json:"mail_to_group_key"`
	MailToMailKey      *string `db:"mail_to_mail_key"              json:"mail_to_mail_key"`
	MailCcGroupKey     *uint64 `db:"mail_cc_group_key"             json:"mail_cc_group_key"`
	MailCcMailKey      *string `db:"mail_cc_mail_key"              json:"mail_cc_mail_key"`
	MailBccGroupKey    *uint64 `db:"mail_bcc_group_key"            json:"mail_bcc_group_key"`
	MailBccMailKey     *string `db:"mail_bcc_mail_key"             json:"mail_bcc_mail_key"`
	MailSubject        *string `db:"mail_subject"                  json:"mail_subject"`
	MailBody           *string `db:"mail_body"                     json:"mail_body"`
	MailImageFile      *string `db:"mail_image_file"               json:"mail_image_file"`
	MailHasAttachment  uint8   `db:"mail_has_attachment"           json:"mail_has_attachment"`
	MailFlagHtml       uint8   `db:"mail_flag_html"                json:"mail_flag_html"`
	JobExecuteDate     string  `db:"job_execute_date"              json:"job_execute_date"`
	MailJobItemKey     uint64  `db:"mail_job_item_key"             json:"mail_job_item_key"`
	JobIsExecute       uint8   `db:"job_is_execute"                json:"job_is_execute"`
	JobSentDate        string  `db:"job_sent_date"                 json:"job_sent_date"`
	JobSentCount       uint64  `db:"job_sent_count"                json:"job_sent_count"`
	JobErrorLog        *string `db:"job_error_log"                 json:"job_error_log"`
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

func GetMmMailSentLog(c *MmMailMaster, field string, key string) (int, error) {
	query := `SELECT mm_mail_sent_log.* FROM mm_mail_sent_log 
	WHERE mm_mail_sent_log.rec_status = '1' 
	AND mm_mail_sent_log.` + field + ` = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateMmMailSentLog(params map[string]string) (int, error, string) {
	query := "INSERT INTO mm_mail_sent_log"
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

func UpdateMmMailSentLog(params map[string]string) (int, error) {
	query := "UPDATE mm_mail_sent_log SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "mm_mail_sent_log" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE mm_mail_sent_log = " + params["mm_mail_sent_log"]
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
