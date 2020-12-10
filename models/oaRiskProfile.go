package models

import (
	"api/db"
	_ "database/sql"
	"net/http"
	_ "strconv"

	log "github.com/sirupsen/logrus"
	"github.com/shopspring/decimal"
)

type OaRiskProfile struct {
	OaRiskProfileKey  uint64  `db:"oa_risk_profile_key"   json:"oa_risk_profile_key"`
	OaRequestKey      uint64  `db:"oa_request_key"        json:"oa_request_key"`
	RiskProfileKey    uint64  `db:"risk_profile_key"      json:"risk_profile_key"`
	ScoreResult       *uint64 `db:"score_result"          json:"score_result"`
	RecOrder          *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type AdminOaRiskProfile struct {
	OaRiskProfileKey *uint64 `db:"oa_risk_profile_key"        json:"oa_risk_profile_key"`
	RiskProfileKey   *uint64 `db:"risk_profile_key"           json:"risk_profile_key"`
	RiskCode         string  `db:"risk_code"                  json:"risk_code"`
	RiskName         *string `db:"risk_name"                  json:"risk_name"`
	RiskDesc         *string `db:"risk_desc"                  json:"risk_desc"`
	MinScore         decimal.Decimal `db:"min_score"                  json:"min_score"`
	MaxScore         decimal.Decimal `db:"max_score"                  json:"max_score"`
	ScoreResult      *uint64 `db:"score_result"               json:"score_result"`
}

func CreateOaRiskProfile(params map[string]string) (int, error) {
	query := "INSERT INTO oa_risk_profile"
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

func AdminGetOaRiskProfile(c *[]AdminOaRiskProfile, oaRequestKey string) (int, error) {
	query := `SELECT 
				oa.oa_risk_profile_key AS oa_risk_profile_key,
				oa.risk_profile_key AS risk_profile_key,
				ms.risk_code AS risk_code,
				ms.risk_name AS risk_name,
				ms.risk_desc AS risk_desc,
				ms.min_score AS min_score,
				ms.max_score AS max_score,
				oa.score_result AS score_result
			FROM oa_risk_profile AS oa 
			INNER JOIN ms_risk_profile AS ms ON ms.risk_profile_key = oa.risk_profile_key
			WHERE oa.rec_status = 1 AND oa.oa_request_key =  ` + oaRequestKey

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetOaRiskProfile(c *OaRiskProfile, key string, field string) (int, error) {
	query := "SELECT oa_risk_profile.* FROM oa_risk_profile WHERE oa_risk_profile." + field + " = " + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
