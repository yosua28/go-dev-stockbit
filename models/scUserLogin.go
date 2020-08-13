package models

import (
	"api/db"
	"log"
	"net/http"
	"strconv"
)

type ScUserLogin struct {
	UserLoginKey            uint64    `db:"user_login_key"            json:"user_login_key"`
	UserCategoryKey         uint64    `db:"user_category_key"         json:"user_category_key"`
	UserDeptKey            *uint64    `db:"user_dept_key"             json:"user_dept_key"`
	UserDeptKey1           *uint64    `db:"user_dept_key1"            json:"user_dept_key1"`
	UloginName              string    `db:"ulogin_name"               json:"ulogin_name"`
	UloginFullName          string    `db:"ulogin_full_name"          json:"ulogin_full_name"`
	UloginPassword          string    `db:"ulogin_password"           json:"ulogin_password"`
	UloginEmail             string    `db:"ulogin_email"              json:"ulogin_email"`
	LastVerifiedEmail      *string    `db:"last_verified_email"       json:"last_verified_email"`
	VerifiedEmail          *uint8     `db:"verified_email"            json:"verified_email"`
	UloginPin              *string    `db:"ulogin_pin"                json:"ulogin_pin"`
	LastChangedPin         *string    `db:"last_changed_pin"          json:"last_changed_pin"`
	UloginMobileno         *string    `db:"ulogin_mobileno"           json:"ulogin_mobileno"`
	LastVerifiedMobileno   *string    `db:"last_verified_mobileno"    json:"last_verified_mobileno"`
	VerifiedMobileno        uint8     `db:"verified_mobileno"         json:"verified_mobileno"`
	UloginMustChangepwd     string    `db:"ulogin_must_changepwd"     json:"ulogin_must_changepwd"`
	LastPasswordChanged    *string    `db:"last_password_changed"     json:"last_password_changed"`
	UloginLocked            uint8     `db:"ulogin_locked"             json:"ulogin_locked"`
	UloginEnabled           uint8     `db:"ulogin_enabled"            json:"ulogin_enabled"`
	UloginFailedCount       uint64    `db:"ulogin_failed_count"       json:"ulogin_failed_count"`
	LastAccess             *string    `db:"last_access"               json:"last_access"`
	AcceptLoginTnc          uint8     `db:"accept_login_tnc"          json:"accept_login_tnc"`
	AllowedSharingLogin     uint8     `db:"allowed_sharing_login"     json:"allowed_sharing_login"`
	RecOrder               *uint64    `db:"rec_order"                 json:"rec_order"`
	RecStatus               uint8     `db:"rec_status"                json:"rec_status"`
	RecCreatedDate         *string    `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy           *string    `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate        *string    `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy          *string    `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1              *string    `db:"rec_image1"                json:"rec_image1"`
	RecImage2              *string    `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus      *uint8     `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage       *uint64    `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate        *string    `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy          *string    `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate         *string    `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy           *string    `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1        *string    `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2        *string    `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3        *string    `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func GetAllScUserLogin(c *[]ScUserLogin, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              sc_user_login.* FROM 
			  sc_user_login`
	var present bool
	var whereClause []string
	var condition string
	
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType"){
			whereClause = append(whereClause, "sc_user_login."+field + " = '" + value + "'")
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
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetScUserLogin(c *ScUserLogin, email string) (int, error) {
	query := `SELECT sc_user_login.* WHERE sc_user_login.ulogin_email = ` + email
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateScUserLogin(params map[string]string) (int, error){
	query := "INSERT INTO sc_user_login"
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
		query += "("+fields + ") VALUES(" + values + ")"
		log.Println(query)
	
		tx, err := db.Db.Begin()
		if err != nil {
			log.Println(err)
			return http.StatusBadGateway, err
		}
		_, err = tx.Exec(query, bindvars...)
		tx.Commit()
		if err != nil {
			log.Println(err)
			return http.StatusBadRequest, err
		}
		return http.StatusOK, nil
}