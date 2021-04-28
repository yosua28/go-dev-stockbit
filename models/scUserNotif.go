package models

import (
	"api/db"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ScUserNotif struct {
	NotifHdrKey       uint64  `db:"notif_hdr_key"             json:"notif_hdr_key"`
	NotifCategory     *uint64 `db:"notif_category"            json:"notif_category"`
	NotifDateSent     *string `db:"notif_date_sent"           json:"notif_date_sent"`
	UmessageSubject   *string `db:"umessage_subject"          json:"umessage_subject"`
	UmessageBody      *string `db:"umessage_body"             json:"umessage_body"`
	AlertNotifType    uint64  `db:"alert_notif_type"          json:"alert_notif_type"`
	RecOrder          *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func CreateScUserNotif(params map[string]string) (int, error) {
	query := "INSERT INTO sc_user_notif"
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

func GetScUserNotif(c *ScUserNotif, key string) (int, error) {
	query := `SELECT sc_user_notif.* FROM sc_user_notif WHERE sc_user_notif.notif_hdr_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateScUserNotif(params map[string]string, where map[string]string) (int, error) {
	query := "UPDATE sc_user_notif SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "notif_hdr_key" {

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
