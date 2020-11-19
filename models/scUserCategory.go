package models

import (
	"api/db"
	"log"
	"net/http"
	"strconv"
)

type ScUserCategory struct {
	UserCategoryKey   uint64  `db:"user_category_key"              json:"user_category_key"`
	UcategoryCode     string  `db:"ucategory_code"                 json:"ucategory_code"`
	UcategoryName     string  `db:"ucategory_name"                 json:"ucategory_name"`
	UcategoryDesc     string  `db:"ucategory_desc"                 json:"ucategory_desc"`
	UcategoryClass    string  `db:"ucategory_class"                json:"ucategory_class"`
	RecOrder          *uint64 `db:"rec_order"                      json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                     json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"               json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"                 json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"              json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"                json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                     json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                     json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"            json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"             json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"              json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"                json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"               json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"                 json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"              json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"              json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"              json:"rec_attribute_id3"`
}

type ScUserCategoryInfo struct {
	UserCategoryKey uint64 `json:"user_category_key"`
	UcategoryCode   string `json:"ucategory_code"`
	UcategoryName   string `json:"ucategory_name"`
	UcategoryDesc   string `json:"ucategory_desc"`
}

func GetScUserCategory(c *ScUserCategory, key string) (int, error) {
	query := `SELECT sc_user_category.* FROM sc_user_category 
				WHERE sc_user_category.rec_status = 1 AND sc_user_category.user_category_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetAllScUserCategory(c *[]ScUserCategory, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              sc_user_category.* FROM 
			  sc_user_category `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "sc_user_category."+field+" = '"+value+"'")
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
