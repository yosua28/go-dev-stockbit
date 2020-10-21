package models

import (
	"api/db"
	"log"
	"net/http"
)

type MsAgentCustomer struct {
	AgentCustomerKey  uint64  `db:"agent_customer_key"   json:"agent_customer_key"`
	AgentKey          uint64  `db:"agent_key"            json:"agent_key"`
	CustomerKey       uint64  `db:"customer_key"         json:"customer_key"`
	EffDate           *string `db:"eff_date"             json:"eff_date"`
	RefCodeUsed       *string `db:"ref_code_used"        json:"ref_code_used"`
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

func GetLastAgenCunstomer(c *MsAgentCustomer, customerKey string) (int, error) {
	query2 := `SELECT t1.agent_customer_key, t1.agent_key, t1.customer_key, t1.eff_date, t1.ref_code_used FROM
			   ms_agent_customer t1 JOIN (SELECT MAX(agent_customer_key) agent_customer_key FROM ms_agent_customer GROUP BY customer_key) t2
			   ON t1.agent_customer_key = t2.agent_customer_key`
	query := query2 + " WHERE t1.customer_key =" + customerKey

	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateMsAgentCustomer(params map[string]string) (int, error) {
	query := "INSERT INTO ms_agent_customer"
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
