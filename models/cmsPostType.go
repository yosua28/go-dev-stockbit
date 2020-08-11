package models

import (
	"api/db"
	"log"
	"strconv"
	"net/http"
)

type CmsPostTypeData struct{
	PostTypeKey               uint64        `json:"post_type_key"`
	PostTypeCode              string        `json:"post_type_code"`
	PostTypeName             *string        `json:"post_type_name"`
	PostTypeDesc             *string        `json:"post_type_desc"`
	PostTypeGroup            *string        `json:"post_type_group"`
	PostList                 []CmsPostData  `json:"post_list"`
}
type CmsPostType struct{
	PostTypeKey               uint64    `db:"post_type_key"             json:"post_type_key"`
	PostTypeCode              string    `db:"post_type_code"            json:"post_type_code"`
	PostTypeName             *string    `db:"post_type_name"            json:"post_type_name"`
	PostTypeDesc             *string    `db:"post_type_desc"            json:"post_type_desc"`
	PostTypeGroup            *string    `db:"post_type_group"           json:"post_type_group"`
	RecOrder                 *uint64    `db:"rec_order"                 json:"rec_order"`
	RecStatus                 uint8     `db:"rec_status"                json:"rec_status"`
	RecCreatedDate           *string    `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy             *string    `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate          *string    `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy            *string    `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1                *string    `db:"rec_image1"                json:"rec_image1"`
	RecImage2                *string    `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus        *uint8     `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage         *uint64    `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate          *string    `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy            *string    `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate           *string    `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy             *string    `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1          *string    `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2          *string    `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3          *string    `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func GetAllCmsPostType(c *[]CmsPostType, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              cms_post_type.*
			  FROM cms_post_type`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		whereClause = append(whereClause, field + " = " + value)
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

func GetCmsPostType(c *CmsPostType, field string, value string) (int, error) {
	query := "SELECT cms_post_type.* FROM cms_post_type WHERE cms_post_type."+field+" = ?"

	log.Println(query)
	err := db.Db.Get(c, query, value)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}