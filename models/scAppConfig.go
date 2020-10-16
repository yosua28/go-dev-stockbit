package models

import (
	"api/db"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ScAppConfig struct {
	AppConfigKey      uint64  `db:"app_config_key"             json:"app_config_key"`
	ConfigTypeKey     uint64  `db:"config_type_key"            json:"config_type_key"`
	AppConfigCode     *string `db:"app_config_code"            json:"app_config_code"`
	AppConfigName     *string `db:"app_config_name"            json:"app_config_name"`
	AppConfigDesc     *string `db:"app_config_desc"            json:"app_config_desc"`
	AppConfigDatatype *string `db:"app_config_datatype"        json:"app_config_datatype"`
	AppConfigValue    *string `db:"app_config_value"           json:"app_config_value"`
	RecOrder          *uint64 `db:"rec_order"                  json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                 json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
}

func GetScAppConfigByCode(c *ScAppConfig, code string) (int, error) {
	query := `SELECT sc_app_config.* FROM sc_app_config WHERE sc_app_config.rec_status = 1 
			AND sc_app_config.app_config_code = "` + code + `"`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateMsCustomerByConfigCode(params map[string]string, code string) (int, error) {
	query := "UPDATE sc_app_config SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE app_config_code = '" + code + "'"
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetAllScAppConfig(c *[]ScAppConfig, params map[string]string) (int, error) {
	query := `SELECT
              sc_app_config.* FROM 
			  sc_app_config`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_app_config."+field+" = '"+value+"'")
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
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
