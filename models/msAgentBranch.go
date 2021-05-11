package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type MsAgentBranch struct {
	AgentBranchKey    uint64  `db:"agent_branch_key"     json:"agent_branch_key"`
	AgentKey          uint64  `db:"agent_key"            json:"agent_key"`
	BranchKey         uint64  `db:"branch_key"           json:"branch_key"`
	EffDate           *string `db:"eff_date"             json:"eff_date"`
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

type MsAgentLastBranch struct {
	AgentKey  uint64 `db:"agent_key"            json:"agent_key"`
	AgentCode string `db:"agent_code"            json:"agent_code"`
	Agentname string `db:"agent_name"           json:"agent_name"`
}

func GetMsAgentLastBranch(c *[]MsAgentLastBranch, branchKey string) (int, error) {
	query := `SELECT 
				a.agent_key AS agent_key,
				a.agent_code AS agent_code,
				CONCAT(a.agent_code, " - ", a.agent_name) AS agent_name 
			FROM ms_agent_branch AS mab 
			INNER JOIN ms_agent AS a ON mab.agent_key = a.agent_key
			JOIN (
				SELECT MAX(eff_date) eff_date, agent_key 
				FROM ms_agent_branch 
				GROUP BY agent_key 
				ORDER BY eff_date DESC
				) t2 ON mab.agent_key = t2.agent_key 
			WHERE mab.rec_status = 1 AND a.rec_status = 1 
			AND t2.eff_date = mab.eff_date AND mab.branch_key = '` + branchKey + `'  
			ORDER BY mab.agent_branch_key ASC`
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
