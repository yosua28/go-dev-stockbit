package models

import (
	"api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsAgent struct {
	AgentKey          uint64  `db:"agent_key"            json:"agent_key"`
	AgentID           uint64  `db:"agent_id"             json:"agent_id"`
	AgentCode         string  `db:"agent_code"           json:"agent_code"`
	AgentName         string  `db:"agent_name"           json:"agent_name"`
	AgentEmail        *string `db:"agent_email"          json:"agent_email"`
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

type ListAgentAdmin struct {
	AgentKey      uint64  `db:"agent_key"        json:"agent_key"`
	BranchName    *string `db:"branch_name"      json:"branch_name"`
	AgentId       *string `db:"agent_id"         json:"agent_id"`
	AgentCode     *string `db:"agent_code"       json:"agent_code"`
	AgentName     *string `db:"agent_name"       json:"agent_name"`
	AgentCategory *string `db:"agent_category"   json:"agent_category"`
	AgentChannel  *string `db:"agent_channel"    json:"agent_channel"`
}

func AdminGetListAgent(c *[]ListAgentAdmin, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if searchLike != "" {
		condition += " AND"
		condition += " (b.branch_name like '%" + searchLike + "%' OR"
		condition += " a.agent_id like '%" + searchLike + "%' OR"
		condition += " a.agent_code like '%" + searchLike + "%' OR"
		condition += " a.agent_name like '%" + searchLike + "%' OR"
		condition += " cat.lkp_name like '%" + searchLike + "%' OR"
		condition += " cha.lkp_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				a.agent_key,
				b.branch_name, 
				a.agent_id,
				a.agent_code,
				a.agent_name,
				cat.lkp_name AS agent_category,
				cha.lkp_name AS agent_channel 
			FROM ms_agent AS a
			LEFT JOIN ms_agent_branch AS mab ON mab.agent_key = a.agent_key AND mab.rec_status = 1
			LEFT JOIN ms_branch AS b ON b.branch_key = mab.branch_key 
			LEFT JOIN gen_lookup AS cat ON cat.lookup_key = a.agent_category
			LEFT JOIN gen_lookup AS cha ON cha.lookup_key = a.agent_channel
			WHERE a.rec_status = 1 ` + condition

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	if !nolimit {
		limitOffset += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			limitOffset += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += orderCondition + limitOffset

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CountAdminGetListAgent(c *CountData, params map[string]string, searchLike string) (int, error) {
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if searchLike != "" {
		condition += " AND"
		condition += " (b.branch_name like '%" + searchLike + "%' OR"
		condition += " a.agent_id like '%" + searchLike + "%' OR"
		condition += " a.agent_code like '%" + searchLike + "%' OR"
		condition += " a.agent_name like '%" + searchLike + "%' OR"
		condition += " cat.lkp_name like '%" + searchLike + "%' OR"
		condition += " cha.lkp_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				count(a.agent_key) AS count_data 
			FROM ms_agent AS a
			LEFT JOIN ms_agent_branch AS mab ON mab.agent_key = a.agent_key AND mab.rec_status = 1
			LEFT JOIN ms_branch AS b ON b.branch_key = mab.branch_key 
			LEFT JOIN gen_lookup AS cat ON cat.lookup_key = a.agent_category
			LEFT JOIN gen_lookup AS cha ON cha.lookup_key = a.agent_channel
			WHERE a.rec_status = 1 ` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateMsAgent(params map[string]string) (int, error) {
	query := "UPDATE ms_agent SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "agent_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE agent_key = " + params["agent_key"]
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

func CreateMsAgent(params map[string]string) (int, error, string) {
	query := "INSERT INTO ms_agent"
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
		return http.StatusBadGateway, err, "0"
	}
	ret, err := tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func CountMsAgentValidateUnique(c *CountData, field string, value string, menuKey string) (int, error) {
	query := `SELECT 
				COUNT(agent_key) AS count_data 
			FROM ms_agent
			WHERE ` + field + ` = '` + value + `'`

	if menuKey != "" {
		query += " AND agent_key != '" + menuKey + "'"
	}

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type MsAgentBranchDetail struct {
	AgentKey       uint64  `db:"agent_key"            json:"agent_key"`
	BranchKey      *uint64 `db:"branch_key"           json:"branch_key"`
	BranchName     *string `db:"branch_name"          json:"branch_name"`
	AgentId        uint64  `db:"agent_id"             json:"agent_id"`
	AgentCode      string  `db:"agent_code"           json:"agent_code"`
	AgentName      string  `db:"agent_name"           json:"agent_name"`
	AgentEmail     *string `db:"agent_email"          json:"agent_email"`
	AgentShortName *string `db:"agent_short_name"     json:"agent_short_name"`
	AgentCategory  *uint64 `db:"agent_category"       json:"agent_category"`
	AgentChannel   *uint64 `db:"agent_channel"        json:"agent_channel"`
	Remarks        *string `db:"remarks"              json:"remarks"`
	ReferenceCode  *string `db:"reference_code"       json:"reference_code"`
	RecOrder       *uint64 `db:"rec_order"            json:"rec_order"`
}

func AdminGetDetailAgent(c *MsAgentBranchDetail, key string) (int, error) {
	query := `SELECT 
				a.agent_key,
				b.branch_key,
				b.branch_name, 
				a.agent_id,
				a.agent_code,
				a.agent_name,
				a.agent_email,
				a.agent_short_name,
				a.agent_category,
				a.agent_channel,
				a.remarks,
				a.reference_code,
				a.rec_order 
			FROM ms_agent AS a
			LEFT JOIN ms_agent_branch AS mab ON mab.agent_key = a.agent_key AND mab.rec_status = 1
			LEFT JOIN ms_branch AS b ON b.branch_key = mab.branch_key 
			LEFT JOIN gen_lookup AS cat ON cat.lookup_key = a.agent_category
			LEFT JOIN gen_lookup AS cha ON cha.lookup_key = a.agent_channel
			WHERE a.rec_status = 1 AND a.agent_key = '` + key + `'`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
