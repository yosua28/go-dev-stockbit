package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ScLoginLog struct {
	LoginLogKey          uint64  `db:"login_log_key"             json:"login_log_key"`
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

func CreateScLoginLog(params map[string]string) (int, error) {
	query := "INSERT INTO sc_login_log"
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
