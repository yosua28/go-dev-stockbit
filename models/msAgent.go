package models

import (
	"api/db"
	"net/http"
	"strings"
	"strconv"
	
	log "github.com/sirupsen/logrus"
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

type MsAgentDropdown struct {
	AgentKey  uint64 `db:"agent_key"            json:"agent_key"`
	AgentName string `db:"agent_name"           json:"agent_name"`
}

func GetAllMsAgent(c *[]MsAgent, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ms_agent.* FROM 
			  ms_agent`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_agent."+field+" = '"+value+"'")
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

func GetMsAgentIn(c *[]MsAgent, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_agent.* FROM 
				ms_agent `
	query := query2 + " WHERE ms_agent.rec_status = 1 AND ms_agent." + field + " IN(" + inQuery + ")"

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Info(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsAgent(c *MsAgent, key string) (int, error) {
	query := `SELECT ms_agent.* FROM ms_agent WHERE ms_agent.rec_status = 1 AND ms_agent.agent_key = ` + key
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Info(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsAgentByField(c *MsAgent, value string, field string) (int, error) {
	query := `SELECT ms_agent.* FROM ms_agent WHERE ms_agent.rec_status = 1 AND ms_agent.` + field + ` = ` + value
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Info(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsAgentDropdown(c *[]MsAgentDropdown) (int, error) {
	query := `SELECT 
				agent_key, 
 				CONCAT(agent_code, " - ", agent_name) AS agent_name 
			FROM ms_agent WHERE ms_agent.rec_status = 1`
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Info(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
