package models

import (
	"api/db"
	"log"
	"net/http"
	"strconv"
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
