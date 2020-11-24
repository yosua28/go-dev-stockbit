package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type ScLoginSessionInfo struct {
	SessionID          string  `json:"session_id"`
	Email              *string `json:"email,omitempty"`
	Expired            *string `json:"expired,omitempty"`
	MustChangePassword bool    `json:"must_change_password"`
}

type ScLoginSession struct {
	LoginSessionKey      uint64  `db:"login_session_key"         json:"login_session_key"`
	SessionID            *string `db:"session_id"                json:"session_id"`
	LoginDate            string  `db:"login_date"                json:"login_date"`
	LogoutDate           *string `db:"logout_date"               json:"logout_date"`
	UserLoginKey         uint64  `db:"user_login_key"            json:"user_login_key"`
	TerminalName         *string `db:"terminal_name"             json:"terminal_name"`
	DeviceModel          *string `db:"device_model"              json:"device_model"`
	WorkstationName      *string `db:"workstation_name"          json:"workstation_name"`
	WorkstationIpaddress *string `db:"workstation_ipaddress"     json:"workstation_ipaddress"`
	ClientAgent          *string `db:"client_agent"              json:"client_agent"`
	AccessLocation       *string `db:"access_location"           json:"access_location"`
	AccessNotes          *string `db:"access_notes"              json:"access_notes"`
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

func GetAllScLoginSession(c *[]ScLoginSession, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              sc_login_session.* FROM 
			  sc_login_session`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_login_session."+field+" = '"+value+"'")
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

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetScLoginSession(c *ScLoginSession, key string) (int, error) {
	query := `SELECT sc_login_session.* WHERE sc_login_session.user_login_key = ` + key
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateScLoginSession(params map[string]string) (int, error) {
	query := "INSERT INTO sc_login_session"
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

func UpdateScLoginSession(params map[string]string) (int, error) {
	var isLogin bool
	isLogin = true
	query := "UPDATE sc_login_session SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "login_session_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}

		if (isLogin == true) && (key == "logout_date") {
			isLogin = false
		}
	}

	if isLogin == true {
		if i > 0 {
			query += ", logout_date = NULL "
		} else {
			query += " logout_date = NULL "
		}
	}

	query += " WHERE user_login_key = " + params["user_login_key"]
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
