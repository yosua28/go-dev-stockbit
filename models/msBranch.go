package models

import (
	"api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsBranch struct {
	BranchKey          uint64  `db:"branch_key"            json:"branch_key"`
	ParticipantKey     *uint64 `db:"participant_key"       json:"participant_key"`
	BranchCode         string  `db:"branch_code"           json:"branch_code"`
	BranchName         string  `db:"branch_name"           json:"branch_name"`
	BranchCategory     *uint64 `db:"branch_category"       json:"branch_category"`
	CityKey            *uint64 `db:"city_key"              json:"city_key"`
	BranchAddress      *string `db:"branch_address"        json:"branch_address"`
	BranchEstablished  *string `db:"branch_established"    json:"branch_established"`
	BranchPicName      *string `db:"branch_pic_name"       json:"branch_pic_name"`
	BranchPicEmail     *string `db:"branch_pic_email"      json:"branch_pic_email"`
	BranchPicPhoneno   *string `db:"branch_pic_phoneno"    json:"branch_pic_phoneno"`
	BranchCostCenter   *string `db:"branch_cost_center"    json:"branch_cost_center"`
	BranchProfitCenter *string `db:"branch_profit_center"  json:"branch_profit_center"`
	RecOrder           *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}
type MsBranchDropdown struct {
	BranchKey  uint64 `db:"branch_key"            json:"branch_key"`
	BranchName string `db:"branch_name"           json:"branch_name"`
}

func GetMsBranchIn(c *[]MsBranch, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_branch.* FROM 
				ms_branch `
	query := query2 + " WHERE ms_branch.rec_status = 1 AND ms_branch." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsBranch(c *MsBranch, key string) (int, error) {
	query := `SELECT ms_branch.* FROM ms_branch WHERE ms_branch.rec_status = 1 AND ms_branch.branch_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsBranchDropdown(c *[]MsBranchDropdown) (int, error) {
	query := `SELECT 
				branch_key, 
 				CONCAT(branch_code, " - ", branch_name) AS branch_name 
			FROM ms_branch WHERE ms_branch.rec_status = 1`
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

type ListMsBranch struct {
	AppConfigKey      uint64  `db:"branch_key"            json:"branch_key"`
	ConfigTypeCode    *string `db:"participant_name"      json:"participant_name"`
	AppConfigCode     string  `db:"branch_code"           json:"branch_code"`
	AppConfigName     string  `db:"branch_name"           json:"branch_name"`
	AppConfigDesc     *string `db:"branch_category"       json:"branch_category"`
	AppConfigDatatype *string `db:"city_name"             json:"city_name"`
}

func AdminGetListMsBranch(c *[]ListMsBranch, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
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
		condition += " (p.participant_name like '%" + searchLike + "%' OR"
		condition += " b.branch_code like '%" + searchLike + "%' OR"
		condition += " b.branch_name like '%" + searchLike + "%' OR"
		condition += " cat.lkp_name like '%" + searchLike + "%' OR"
		condition += " c.city_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				b.branch_key,
				p.participant_name,
				b.branch_code,
				b.branch_name,
				cat.lkp_name AS branch_category,
				c.city_name 
			FROM ms_branch AS b
			LEFT JOIN ms_participant AS p ON p.participant_key = b.participant_key
			LEFT JOIN ms_city AS c ON c.city_key = b.city_key
			LEFT JOIN gen_lookup AS cat ON cat.lookup_key = b.branch_category
			WHERE b.rec_status = 1 ` + condition

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

func CountAdminGetListMsBranch(c *CountData, params map[string]string, searchLike string) (int, error) {
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
		condition += " (p.participant_name like '%" + searchLike + "%' OR"
		condition += " b.branch_code like '%" + searchLike + "%' OR"
		condition += " b.branch_name like '%" + searchLike + "%' OR"
		condition += " cat.lkp_name like '%" + searchLike + "%' OR"
		condition += " c.city_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				count(b.branch_key) as count_data 
			FROM ms_branch AS b
			LEFT JOIN ms_participant AS p ON p.participant_key = b.participant_key
			LEFT JOIN ms_city AS c ON c.city_key = b.city_key
			LEFT JOIN gen_lookup AS cat ON cat.lookup_key = b.branch_category
			WHERE b.rec_status = 1 ` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateMsBranch(params map[string]string) (int, error) {
	query := "UPDATE ms_branch SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "branch_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE branch_key = " + params["branch_key"]
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

func CreateMsBranch(params map[string]string) (int, error) {
	query := "INSERT INTO ms_branch"
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

func CountMsBranchValidateUnique(c *CountData, field string, value string, key string) (int, error) {
	query := `SELECT 
				COUNT(branch_key) AS count_data 
			FROM ms_branch
			WHERE rec_status = '1' AND ` + field + ` = '` + value + `'`

	if key != "" {
		query += " AND branch_key != '" + key + "'"
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
