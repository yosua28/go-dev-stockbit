package models

import (
	"api/db"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type ScUserDept struct {
	UserDeptKey          uint64  `db:"user_dept_key"                  json:"user_dept_key"`
	UserDeptParent       *uint64 `db:"user_dept_parent"               json:"user_dept_parent"`
	UserDeptCode         string  `db:"user_dept_code"                 json:"user_dept_code"`
	UserDeptName         string  `db:"user_dept_name"                 json:"user_dept_name"`
	UserDeptDesc         *string `db:"user_dept_desc"                 json:"user_dept_desc"`
	UserDeptEmailAddress *string `db:"user_dept_email_address"        json:"user_dept_email_address"`
	RolePrivileges       *uint64 `db:"role_privileges"                json:"role_privileges"`
	BranchKey            *uint64 `db:"branch_key"                     json:"branch_key"`
	RecOrder             *uint64 `db:"rec_order"                      json:"rec_order"`
	RecStatus            uint8   `db:"rec_status"                     json:"rec_status"`
	RecCreatedDate       *string `db:"rec_created_date"               json:"rec_created_date"`
	RecCreatedBy         *string `db:"rec_created_by"                 json:"rec_created_by"`
	RecModifiedDate      *string `db:"rec_modified_date"              json:"rec_modified_date"`
	RecModifiedBy        *string `db:"rec_modified_by"                json:"rec_modified_by"`
	RecImage1            *string `db:"rec_image1"                     json:"rec_image1"`
	RecImage2            *string `db:"rec_image2"                     json:"rec_image2"`
	RecApprovalStatus    *uint8  `db:"rec_approval_status"            json:"rec_approval_status"`
	RecApprovalStage     *uint64 `db:"rec_approval_stage"             json:"rec_approval_stage"`
	RecApprovedDate      *string `db:"rec_approved_date"              json:"rec_approved_date"`
	RecApprovedBy        *string `db:"rec_approved_by"                json:"rec_approved_by"`
	RecDeletedDate       *string `db:"rec_deleted_date"               json:"rec_deleted_date"`
	RecDeletedBy         *string `db:"rec_deleted_by"                 json:"rec_deleted_by"`
	RecAttributeID1      *string `db:"rec_attribute_id1"              json:"rec_attribute_id1"`
	RecAttributeID2      *string `db:"rec_attribute_id2"              json:"rec_attribute_id2"`
	RecAttributeID3      *string `db:"rec_attribute_id3"              json:"rec_attribute_id3"`
}

type ScUserDeptInfo struct {
	UserDeptKey  uint64  `json:"user_dept_key"`
	UserDeptCode string  `json:"user_dept_code"`
	UserDeptName string  `json:"user_dept_name"`
	UserDeptDesc *string `json:"user_dept_desc"`
}

type ListUserDeptAdmin struct {
	UserDeptKey    uint64  `db:"user_dept_key"        json:"user_dept_key"`
	ParentDept     *string `db:"parent_dept"          json:"parent_dept"`
	UserDeptCode   *string `db:"user_dept_code"       json:"user_dept_code"`
	UserDeptName   *string `db:"user_dept_name"       json:"user_dept_name"`
	RolePrivilages *string `db:"role_privileges"      json:"role_privileges"`
	BranchName     *string `db:"branch_name"          json:"branch_name"`
}

func GetScUserDept(c *ScUserDept, key string) (int, error) {
	query := `SELECT sc_user_dept.* FROM sc_user_dept 
				WHERE sc_user_dept.rec_status = 1 AND sc_user_dept.user_dept_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetAllScUserDept(c *[]ScUserDept, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              sc_user_dept.* FROM 
			  sc_user_dept `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_user_dept."+field+" = '"+value+"'")
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
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetListScUserDept(c *[]ListUserDeptAdmin, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
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
		condition += " (ud.user_dept_key like '%" + searchLike + "%' OR"
		condition += " par.user_dept_name like '%" + searchLike + "%' OR"
		condition += " ud.user_dept_code like '%" + searchLike + "%' OR"
		condition += " ud.user_dept_name like '%" + searchLike + "%' OR"
		condition += " ud.user_dept_desc like '%" + searchLike + "%' OR"
		condition += " g.lkp_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				dat.user_dept_key,
				dat.parent_dept,
				dat.user_dept_code,
				dat.user_dept_name,
				dat.role_privileges,
				dat.branch_name 
			FROM
				(SELECT 
					ud.user_dept_key,
					par.user_dept_name AS parent_dept,
					ud.user_dept_code,
					ud.user_dept_name,
					ud.user_dept_desc,
					g.lkp_name AS role_privileges,
					b.branch_name,
					(CASE
						WHEN par.user_dept_key IS NULL THEN ud.rec_status
						ELSE par.rec_status
					END) AS parent_rec_status   
				FROM sc_user_dept AS ud 
					LEFT JOIN sc_user_dept AS par ON par.user_dept_key = ud.user_dept_parent 
					LEFT JOIN gen_lookup AS g ON g.lookup_key = ud.role_privileges AND g.rec_status = 1
					LEFT JOIN ms_branch AS b ON b.branch_key = ud.branch_key AND b.rec_status = 1
				WHERE ud.rec_status = 1  ` + condition + `
				) AS dat
			WHERE dat.parent_rec_status = 1`

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY dat." + orderBy
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

func CountAdminGetListScUserDept(c *CountData, params map[string]string, searchLike string) (int, error) {
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
		condition += " (ud.user_dept_key like '%" + searchLike + "%' OR"
		condition += " par.user_dept_name like '%" + searchLike + "%' OR"
		condition += " ud.user_dept_code like '%" + searchLike + "%' OR"
		condition += " ud.user_dept_name like '%" + searchLike + "%' OR"
		condition += " ud.user_dept_desc like '%" + searchLike + "%' OR"
		condition += " g.lkp_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				count(dat.user_dept_key) AS count_data 
			FROM
				(SELECT 
					ud.user_dept_key,
					(CASE
						WHEN par.user_dept_key IS NULL THEN ud.rec_status
						ELSE par.rec_status
					END) AS parent_rec_status   
				FROM sc_user_dept AS ud 
					LEFT JOIN sc_user_dept AS par ON par.user_dept_key = ud.user_dept_parent 
					LEFT JOIN gen_lookup AS g ON g.lookup_key = ud.role_privileges AND g.rec_status = 1
					LEFT JOIN ms_branch AS b ON b.branch_key = ud.branch_key AND b.rec_status = 1
				WHERE ud.rec_status = 1  ` + condition + `
				) AS dat
			WHERE dat.parent_rec_status = 1`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateScUserDept(params map[string]string) (int, error) {
	query := "UPDATE sc_user_dept SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "user_dept_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE user_dept_key = " + params["user_dept_key"]
	log.Println(query)

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
func CreateScUserDept(params map[string]string) (int, error) {
	query := "INSERT INTO sc_user_dept"
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

func CountScUserDeptValidateUnique(c *CountData, field string, value string, key string) (int, error) {
	query := `SELECT 
				COUNT(user_dept_key) AS count_data 
			FROM sc_user_dept
			WHERE rec_status = '1' AND ` + field + ` = '` + value + `'`

	if key != "" {
		query += " AND user_dept_key != '" + key + "'"
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
