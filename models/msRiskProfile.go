package models

import (
	"api/db"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsRiskProfileInfo struct {
	RiskProfileKey       uint64     `json:"risk_profile_key"`
	RiskCode             string     `json:"risk_code"`
	RiskName            *string     `json:"risk_name"`
	RiskDesc            *string     `json:"risk_desc"`
	Score                uint64     `json:"score"`
}

type MsRiskProfile struct {
	RiskProfileKey       uint64     `db:"risk_profile_key"      json:"risk_profile_key"`
	RiskCode             string     `db:"risk_code"             json:"risk_code"`
	RiskName             *string    `db:"risk_name"             json:"risk_name"`
	RiskDesc             *string    `db:"risk_desc"             json:"risk_desc"`
	MinScore             float32    `db:"min_score"             json:"min_score"`
	MaxScore             float32    `db:"max_score"             json:"max_score"`
	MaxFlag              uint8      `db:"max_flag"              json:"max_flag"`
	RecOrder             *uint64    `db:"rec_order"             json:"rec_order"`
	RecStatus            uint8      `db:"rec_status"            json:"rec_status"`
	RecCreatedDate       *string    `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy         *string    `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate      *string    `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy        *string    `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1            *string    `db:"rec_image1"            json:"rec_image1"`
	RecImage2            *string    `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus    *uint8     `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage     *uint64    `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate      *string    `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy        *string    `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate       *string    `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy         *string    `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1      *string    `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2      *string    `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3      *string    `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

func GetMsRiskProfile(c *MsRiskProfile, key string) (int, error) {
	query := `SELECT ms_risk_profile.* FROM ms_risk_profile WHERE ms_risk_profile.risk_profile_key = ` + key
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsRiskProfileIn(c *[]MsRiskProfile, value []string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_risk_profile.* FROM 
				ms_risk_profile `
	query := query2 + " WHERE ms_risk_profile.risk_profile_key IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsRiskProfileScore(c *MsRiskProfile, score string) (int, error) {
	query := "SELECT ms_risk_profile.* FROM ms_risk_profile WHERE ms_risk_profile.min_score <= " + score + " AND ms_risk_profile.max_score >= " + score
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}