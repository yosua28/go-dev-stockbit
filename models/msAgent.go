package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"
)

type MsAgent struct {
	AgentKey          uint64  `db:"agent_key"            json:"agent_key"`
	AgentID           uint64  `db:"agent_id"             json:"agent_id"`
	AgentCode         string  `db:"agent_code"           json:"agent_code"`
	AgentName         string  `db:"agent_name"           json:"agent_name"`
	AgentShotName     *string `db:"agent_short_name"     json:"agent_short_name"`
	AgentCategory     *uint64 `db:"agent_category"       json:"agent_category"`
	AgentChannel      *uint64 `db:"agent_channel"        json:"agent_channel"`
	ReferenceCode     *string `db:"reference_code"       json:"reference_code"`
	Remarks           *string `db:"remarks"              json:"remarks"`
	RecOrder          *uint64 `db:"rec_order"            json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"           json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"     json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"       json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"    json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"      json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"           json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"           json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"  json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"   json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"    json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"      json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"     json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"       json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"    json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"    json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"    json:"rec_attribute_id3"`
}

func GetMsAgentIn(c *[]MsAgent, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_agent.* FROM 
				ms_agent `
	query := query2 + " WHERE ms_agent.rec_status = 1 AND ms_agent." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsAgent(c *MsAgent, key string) (int, error) {
	query := `SELECT ms_agent.* FROM ms_agent WHERE ms_agent.rec_status = 1 AND ms_agent.agent_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsAgentByField(c *MsAgent, value string, field string) (int, error) {
	query := `SELECT ms_agent.* FROM ms_agent WHERE ms_agent.rec_status = 1 AND ms_agent.` + field + ` = ` + value
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
