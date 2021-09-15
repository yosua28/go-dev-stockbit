package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ScLinkage struct {
	LinkedKey                uint64  `db:"linked_key"                    json:"linked_key"`
	UserLoginKey             uint64  `db:"user_login_key"                json:"user_login_key"`
	SettleChannel            uint64  `db:"settle_channel"                json:"settle_channel"`
	LinkedMobileno           *string `db:"linked_mobileno"               json:"linked_mobileno"`
	LinkedFullName           *string `db:"linked_full_name"              json:"linked_full_name"`
	LinkedMemberType         *string `db:"linked_member_type"            json:"linked_member_type"`
	LinkedToken              *string `db:"linked_token"                  json:"linked_token"`
	LinkedTokenExpired       *string `db:"linked_token_expired"          json:"linked_token_expired"`
	LinkedVerifiedOtp        *string `db:"linked_verified_otp"           json:"linked_verified_otp"`
	LinkedVerifiedDate       *string `db:"linked_verified_date"          json:"linked_verified_date"`
	LinkedVerifiedReferences *string `db:"linked_verified_references"    json:"linked_verified_references"`
	MotherMaidName           *string `db:"mother_maid_name"              json:"mother_maid_name"`
	IdType                   *string `db:"id_type"                       json:"id_type"`
	IdNumber                 *string `db:"id_number"                     json:"id_number"`
	BirthDate                *string `db:"birth_date"                    json:"birth_date"`
	BirthCity                *string `db:"birth_city"                    json:"birth_city"`
	Gender                   *string `db:"gender"                        json:"gender"`
	Address                  *string `db:"address"                       json:"address"`
	Nationality              *string `db:"nationality"                   json:"nationality"`
	Occupation               *string `db:"occupation"                    json:"occupation"`
	Remark                   *string `db:"remark"                        json:"remark"`
	DeviceId                 *string `db:"device_id"                     json:"device_id"`
	LinkedStatus             *string `db:"linked_status"                 json:"linked_status"`
	LinkedName               *string `db:"linked_name"                   json:"linked_name"`
	UserToken                *string `db:"user_token"                    json:"user_token"`
	RecOrder                 *uint8  `db:"rec_order"                     json:"rec_order"`
	RecStatus                uint8   `db:"rec_status"                    json:"rec_status"`
	RecCreatedDate           *string `db:"rec_created_date"              json:"rec_created_date"`
	RecCreatedBy             *string `db:"rec_created_by"                json:"rec_created_by"`
	RecModifiedDate          *string `db:"rec_modified_date"             json:"rec_modified_date"`
	RecModifiedBy            *string `db:"rec_modified_by"               json:"rec_modified_by"`
	RecImage1                *string `db:"rec_image1"                    json:"rec_image1"`
	RecImage2                *string `db:"rec_image2"                    json:"rec_image2"`
	RecApprovalStatus        *uint8  `db:"rec_approval_status"           json:"rec_approval_status"`
	RecApprovalStage         *uint64 `db:"rec_approval_stage"            json:"rec_approval_stage"`
	RecApprovedDate          *string `db:"rec_approved_date"             json:"rec_approved_date"`
	RecApprovedBy            *string `db:"rec_approved_by"               json:"rec_approved_by"`
	RecDeletedDate           *string `db:"rec_deleted_date"              json:"rec_deleted_date"`
	RecDeletedBy             *string `db:"rec_deleted_by"                json:"rec_deleted_by"`
	RecAttributeID1          *string `db:"rec_attribute_id1"             json:"rec_attribute_id1"`
	RecAttributeID2          *string `db:"rec_attribute_id2"             json:"rec_attribute_id2"`
	RecAttributeID3          *string `db:"rec_attribute_id3"             json:"rec_attribute_id3"`
}

func CreateScLinkage(params map[string]string) (int, error) {
	query := "INSERT INTO sc_linkage"
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

func GetLinkageByField(c *ScLinkage, value string, field string) (int, error) {
	query := `SELECT
				* FROM 
				sc_linkage where rec_status = '1' AND ` + field + ` = "` + value + `"  order by linked_key desc limit 1`
	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetLinkageByParams(c *ScLinkage, params map[string]string) (int, error) {
	query := `SELECT
			  sc_linkage.* FROM 
			  sc_linkage`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_linkage."+field+" = '"+value+"'")
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

	condition += " LIMIT 1"
	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateScLinkage(params map[string]string) (int, error) {
	query := "UPDATE sc_linkage SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "linked_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE linked_key = " + params["linked_key"]
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

func UnlinkedUser(params map[string]string, field string, value string) (int, error) {
	query := "UPDATE sc_linkage SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"
		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE rec_status = 1 AND " + field + " = '" + value + "'"
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
