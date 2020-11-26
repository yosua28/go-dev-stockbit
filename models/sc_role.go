package models

import (
	"api/db"
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

type ScRoleInfo struct {
	RoleKey         uint64 `json:"role_key"`
	RoleCategoryKey uint64 `json:"role_category_key"`
	RoleCode        string `json:"role_code"`
	RoleName        string `json:"role_name"`
	RoleDesc        string `json:"role_desc"`
}

type ScRole struct {
	RoleKey           uint64  `db:"role_key"                  json:"role_key"`
	RoleCategoryKey   *uint64 `db:"role_category_key"         json:"role_category_key"`
	RoleCode          *string `db:"role_code"                 json:"role_code"`
	RoleName          *string `db:"role_name"                 json:"role_name"`
	RoleDesc          *string `db:"role_desc"                 json:"role_desc"`
	RecOrder          *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type ScRoleInfoLogin struct {
	RoleKey  uint64  `json:"role_key"`
	RoleCode *string `json:"role_code"`
	RoleName *string `json:"role_name"`
	RoleDesc *string `json:"role_desc"`
}

type AdminRoleManagement struct {
	RoleKey          uint64  `db:"role_key"                  json:"role_key"`
	RoleCategoryCode *string `db:"role_category_code"        json:"role_category_code"`
	RoleCategoryName *string `db:"role_category_name"        json:"role_category_name"`
	RoleCategoryKey  *uint64 `db:"role_category_key"         json:"role_category_key"`
	RoleCode         *string `db:"role_code"                 json:"role_code"`
	RoleName         *string `db:"role_name"                 json:"role_name"`
	RoleDesc         *string `db:"role_desc"                 json:"role_desc"`
}

type AdminRoleManagementDetail struct {
	RoleKey      uint64              `json:"role_key"`
	RoleCategory *ScRoleCategoryInfo `json:"role_category"`
	RoleCode     *string             `json:"role_code"`
	RoleName     *string             `json:"role_name"`
	RoleDesc     *string             `json:"role_desc"`
}

func GetAllScRole(c *[]ScRole, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              sc_role.* FROM 
			  sc_role`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_role."+field+" = '"+value+"'")
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

func GetScRole(c *ScRole, key string) (int, error) {
	query := `SELECT sc_role.* FROM sc_role 
				WHERE sc_role.rec_status = 1 AND sc_role.role_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func AdminGetAllRoleManagement(c *[]AdminRoleManagement, limit uint64, offset uint64, params map[string]string, nolimit bool, searchLike *string) (int, error) {
	query := `SELECT
				role.role_key AS role_key, 
				cat.role_category_code AS role_category_code,
				cat.role_category_name AS role_category_name, 
				role.role_code AS role_code, 
				role.role_name AS role_name,
				role.role_desc AS role_desc 
			  FROM sc_role AS role 
			  INNER JOIN sc_role_category AS cat ON role.role_category_key = cat.role_category_key
			  WHERE role.rec_status = 1`
	var present bool
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

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " role.role_key LIKE '%" + *searchLike + "%' OR"
		condition += " cat.role_category_code LIKE '%" + *searchLike + "%' OR"
		condition += " cat.role_category_name LIKE '%" + *searchLike + "%' OR"
		condition += " role.role_code LIKE '%" + *searchLike + "%' OR"
		condition += " role.role_name LIKE '%" + *searchLike + "%' OR"
		condition += " role.role_desc LIKE '%" + *searchLike + "%' )"
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

func AdminCountDataRoleManagement(c *CountData, params map[string]string, searchLike *string) (int, error) {
	query := `SELECT
	            count(role.role_key) AS count_data 
			  FROM sc_role AS role 
			  INNER JOIN sc_role_category AS cat ON role.role_category_key = cat.role_category_key
			  WHERE role.rec_status = 1`
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

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " role.role_key LIKE '%" + *searchLike + "%' OR"
		condition += " cat.role_category_code LIKE '%" + *searchLike + "%' OR"
		condition += " cat.role_category_name LIKE '%" + *searchLike + "%' OR"
		condition += " role.role_code LIKE '%" + *searchLike + "%' OR"
		condition += " role.role_name LIKE '%" + *searchLike + "%' OR"
		condition += " role.role_desc LIKE '%" + *searchLike + "%' )"
	}

	query += condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateScRole(params map[string]string) (int, error) {
	query := "UPDATE sc_role SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "role_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE role_key = " + params["role_key"]
	log.Println(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func CreateScRole(params map[string]string) (int, error, string) {
	query := "INSERT INTO sc_role"
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
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func AdminGetValidateUniqueMsRole(c *CountData, paramsAnd map[string]string, updateKey *string) (int, error) {
	query := `SELECT
			  count(sc_role.role_key) as count_data 
			  FROM sc_role `

	var andWhereClause []string
	var condition string

	for fieldAnd, valueAnd := range paramsAnd {
		andWhereClause = append(andWhereClause, "sc_role."+fieldAnd+" = '"+valueAnd+"'")
	}

	// Combile where And clause
	if len(andWhereClause) > 0 {
		condition += " WHERE "

		for index, where := range andWhereClause {
			condition += where
			if (len(andWhereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if updateKey != nil {
		if len(andWhereClause) > 0 {
			condition += " AND "
		} else {
			condition += " WHERE "
		}

		condition += " sc_role.role_key != '" + *updateKey + "'"
	}

	query += condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
