package models

import (
	"api/db"
	"database/sql"
	"net/http"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type ScUserMessageList struct {
	UmessageKey         uint64        `db:"umessage_key"              json:"umessage_key"`
	UmessageType        GenLookupInfo `db:"umessage_type"             json:"umessage_type"`
	UmessageReceiptDate *string       `db:"umessage_receipt_date"     json:"umessage_receipt_date"`
	FlagRead            uint8         `db:"flag_read"                 json:"flag_read"`
	UmessageSubject     *string       `db:"umessage_subject"          json:"umessage_subject"`
	UmessageCategory    GenLookupInfo `db:"umessage_category"         json:"umessage_category"`
}

type ScUserMessageData struct {
	UmessageKey         uint64        `db:"umessage_key"              json:"umessage_key"`
	UmessageType        GenLookupInfo `db:"umessage_type"             json:"umessage_type"`
	UmessageReceiptDate *string       `db:"umessage_receipt_date"     json:"umessage_receipt_date"`
	FlagRead            uint8         `db:"flag_read"                 json:"flag_read"`
	UmessageSubject     *string       `db:"umessage_subject"          json:"umessage_subject"`
	UmessageBody        *string       `db:"umessage_body"             json:"umessage_body"`
	UparentKey          *uint64       `db:"uparent_key"               json:"uparent_key"`
	UmessageCategory    GenLookupInfo `db:"umessage_category"         json:"umessage_category"`
}

type ScUserMessage struct {
	UmessageKey          uint64  `db:"umessage_key"              json:"umessage_key"`
	UmessageType         *uint64 `db:"umessage_type"             json:"umessage_type"`
	NotifHdrKey          *uint64 `db:"notif_hdr_key"             json:"notif_hdr_key"`
	UmessageRecipientKey uint64  `db:"umessage_recipient_key"    json:"umessage_recipient_key"`
	UmessageReceiptDate  *string `db:"umessage_receipt_date"     json:"umessage_receipt_date"`
	FlagRead             uint8   `db:"flag_read"                 json:"flag_read"`
	UmessageSenderKey    *uint64 `db:"umessage_sender_key"       json:"umessage_sender_key"`
	UmessageSentDate     *string `db:"umessage_sent_date"        json:"umessage_sent_date"`
	FlagSent             *uint8  `db:"flag_sent"                 json:"flag_sent"`
	UmessageSubject      *string `db:"umessage_subject"          json:"umessage_subject"`
	UmessageBody         *string `db:"umessage_body"             json:"umessage_body"`
	UparentKey           *uint64 `db:"uparent_key"               json:"uparent_key"`
	UmessageCategory     *uint64 `db:"umessage_category"         json:"umessage_category"`
	FlagArchieved        *uint8  `db:"flag_archieved"            json:"flag_archieved"`
	ArchievedDate        *string `db:"archieved_date"            json:"archieved_date"`
	RecOrder             *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus            uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate       *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy         *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate      *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy        *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1            *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2            *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus    *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage     *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate      *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy        *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate       *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy         *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1      *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2      *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3      *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func CreateScUserMessage(params map[string]string) (int, error) {
	query := "INSERT INTO sc_user_message"
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

func GetAllScUserMessage(c *[]ScUserMessage, params map[string]string) (int, error) {
	query := `SELECT
              sc_user_message.* FROM 
			  sc_user_message`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_user_message."+field+" = '"+value+"'")
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
		// condition += " , umessage_receipt_date DESC"
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

func GetScUserMessage(c *ScUserMessage, key string) (int, error) {
	query := `SELECT sc_user_message.* FROM sc_user_message WHERE sc_user_message.umessage_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateScUserMessage(params map[string]string, where map[string]string) (int, error) {
	query := "UPDATE sc_user_message SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "umessage_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 1) > i {
				query += ", "
			}
			i++
		}
	}

	var whereClause []string
	var condition string
	for field, value := range where {
		whereClause = append(whereClause, "sc_user_message."+field+" = '"+value+"'")
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
	query += condition
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetCountUserMessage(c *CountData, params map[string]string) (int, error) {
	query := `SELECT
              count(sc_user_message.umessage_key) as count_data FROM 
			  sc_user_message`

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_user_message."+field+" = '"+value+"'")
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

	query += condition

	// Main query
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateScUserMessageByField(params map[string]string, field string, value string) (int, error) {
	query := "UPDATE sc_user_message SET "
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
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func CreateMultipleUserMessage(params []interface{}) (int, error) {

	q := `INSERT INTO sc_user_message (
		umessage_type, 
		umessage_recipient_key,
		umessage_receipt_date,
		flag_read,
		flag_sent,
		umessage_subject,
		umessage_body,
		umessage_category,
		flag_archieved,
		archieved_date,
		rec_status,
		rec_created_date,
		rec_created_by) VALUES `

	for i := 0; i < len(params); i++ {
		q += "(?)"
		if i < (len(params) - 1) {
			q += ","
		}
	}
	log.Info(q)
	query, args, err := sqlx.In(q, params...)
	if err != nil {
		return http.StatusBadGateway, err
	}

	query = db.Db.Rebind(query)
	_, err = db.Db.Query(query, args...)
	if err != nil {
		log.Error(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
