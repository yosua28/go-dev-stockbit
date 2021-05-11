package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type ScUserNotif struct {
	NotifHdrKey       uint64  `db:"notif_hdr_key"             json:"notif_hdr_key"`
	NotifCategory     *uint64 `db:"notif_category"            json:"notif_category"`
	NotifStart        *string `db:"notif_start"           json:"notif_start"`
	NotifEnd          *string `db:"notif_end"           json:"notif_end"`
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

type UserNotifField struct {
	NotifHdrKey       uint64  `db:"notif_hdr_key"            json:"notif_hdr_key"`
	NotifCategoryKey  *uint64 `db:"notif_category_key"       json:"notif_category_key"`
	NotifCategory     *string `db:"notif_category"           json:"notif_category"`
	NotifStart        *string `db:"notif_start"          json:"notif_start"`
	NotifEnd          *string `db:"notif_end"          json:"notif_end"`
	UmessageSubject   *string `db:"umessage_subject"         json:"umessage_subject"`
	UmessageBody      *string `db:"umessage_body"            json:"umessage_body"`
	AlertNotifTypeKey uint64  `db:"alert_notif_type_key"     json:"alert_notif_type_key"`
	AlertNotifType    string  `db:"alert_notif_type"         json:"alert_notif_type"`
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

func UpdateScUserNotif(params map[string]string) (int, error) {
	query := "UPDATE sc_user_notif SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "notif_hdr_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE notif_hdr_key = " + params["notif_hdr_key"]

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

func AdminGetAllUserNotif(c *[]UserNotifField, limit uint64, offset uint64, params map[string]string, paramsLike string, nolimit bool) (int, error) {
	query := `SELECT 
				s.notif_hdr_key AS notif_hdr_key,
				s.notif_category AS notif_category_key,
				cat.lkp_name AS notif_category,
				DATE_FORMAT(s.notif_start, '%d %M %Y') AS notif_start,
				DATE_FORMAT(s.notif_end, '%d %M %Y') AS notif_end,
				s.umessage_subject AS umessage_subject,
				s.umessage_body AS umessage_body,
				s.alert_notif_type AS alert_notif_type_key,
				ty.lkp_name AS alert_notif_type 
			FROM sc_user_notif AS s
			INNER JOIN gen_lookup AS cat ON cat.lookup_key = s.notif_category
			INNER JOIN gen_lookup AS ty ON ty.lookup_key = s.alert_notif_type
			WHERE s.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	if paramsLike != "" {
		condition += " AND ("
		condition += " s.notif_hdr_key LIKE '%" + paramsLike + "%' OR"
		condition += " cat.lkp_name LIKE '%" + paramsLike + "%' OR"
		condition += " DATE_FORMAT(s.notif_date_sent, '%d %M %Y') LIKE '%" + paramsLike + "%' OR"
		condition += " s.umessage_subject LIKE '%" + paramsLike + "%' OR"
		condition += " s.umessage_body LIKE '%" + paramsLike + "%' OR"
		condition += " ty.lkp_name LIKE '%" + paramsLike + "%')"
	}

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

func CountAdminGetAllUserNotif(c *CountData, params map[string]string, paramsLike string) (int, error) {
	query := `SELECT 
				count(s.notif_hdr_key) AS count_data 
			FROM sc_user_notif AS s
			INNER JOIN gen_lookup AS cat ON cat.lookup_key = s.notif_category
			INNER JOIN gen_lookup AS ty ON ty.lookup_key = s.alert_notif_type
			WHERE s.rec_status = 1`

	var whereClause []string
	var condition string

	if paramsLike != "" {
		condition += " AND ("
		condition += " s.notif_hdr_key LIKE '%" + paramsLike + "%' OR"
		condition += " cat.lkp_name LIKE '%" + paramsLike + "%' OR"
		condition += " DATE_FORMAT(s.notif_date_sent, '%d %M %Y') LIKE '%" + paramsLike + "%' OR"
		condition += " s.umessage_subject LIKE '%" + paramsLike + "%' OR"
		condition += " s.umessage_body LIKE '%" + paramsLike + "%' OR"
		condition += " ty.lkp_name LIKE '%" + paramsLike + "%')"
	}

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

func AdminGetDetailUserNotif(c *UserNotifField, key string) (int, error) {
	query := `SELECT 
				s.notif_hdr_key AS notif_hdr_key,
				s.notif_category AS notif_category_key,
				cat.lkp_name AS notif_category,
				DATE_FORMAT(s.notif_start, '%d %M %Y') AS notif_start,
				DATE_FORMAT(s.notif_end, '%d %M %Y') AS notif_end,
				s.umessage_subject AS umessage_subject,
				s.umessage_body AS umessage_body,
				s.alert_notif_type AS alert_notif_type_key,
				ty.lkp_name AS alert_notif_type 
			FROM sc_user_notif AS s
			INNER JOIN gen_lookup AS cat ON cat.lookup_key = s.notif_category
			INNER JOIN gen_lookup AS ty ON ty.lookup_key = s.alert_notif_type
			WHERE s.rec_status = 1 and s.notif_hdr_key = ` + key

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
