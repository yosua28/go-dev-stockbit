package models

import (
	"api/db"
	"log"
	"strconv"
	"net/http"
	"strings"
)

type CmsPost struct{
	PostKey                   uint64    `db:"post_key"                  json:"post_key"`
	PostSubtypeKey            uint64    `db:"post_subtype_key"          json:"post_subtype_key"`
	PostTitle                 string    `db:"post_title"                json:"post_title"`
	PostSubTitle             *string    `db:"post_sub_title"            json:"post_sub_title"`
	PostContent              *string    `db:"post_content"              json:"post_content"`
	PostContentAuthor        *string    `db:"post_content_author"       json:"post_content_author"`
	PostContentSources       *string    `db:"post_content_sources"      json:"post_content_sources"`
	PostPublishStart          string    `db:"post_publish_start"        json:"post_publish_start"`
	PostPublishThru           string    `db:"post_publish_thru"         json:"post_publish_thru"`
	PostPageAllowed           uint8     `db:"post_page_allowed"         json:"post_page_allowed"`
	PostCommentAllowed        uint8     `db:"post_comment_allowed"      json:"post_comment_allowed"`
	PostCommentDisplayed      uint8     `db:"post_comment_displayed"    json:"post_comment_displayed"`
	PostFilesAllowed          uint8     `db:"post_files_allowed"        json:"post_files_allowed"`
	PostVideoAllowed          uint8     `db:"post_video_allowed"        json:"post_video_allowed"`
	PostVideoUrl             *string    `db:"post_video_url"            json:"post_video_url"`
	PostPinned                uint8     `db:"post_pinned"               json:"post_pinned"`
	PostOwnerKey             *uint64    `db:"post_owner_key"            json:"post_owner_key"`
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

type CmsPostData struct {
	PostKey                   uint64             `json:"post_key"`
	PostSubtype               CmsPostSubtypeInfo `json:"post_subtype"`
	PostTitle                 string             `json:"post_title"`
	PostSubTitle              string             `json:"post_sub_title"`
	PostContent               string             `json:"post_content"`
	PostContentAuthor         string             `json:"post_content_author"`
	PostContentSources        string             `json:"post_content_sources"`
	PostPublishStart          string             `json:"post_publish_start"`
	PostPublishThru           string             `json:"post_publish_thru"`
	PostPageAllowed           bool               `json:"post_page_allowed"`
	PostCommentAllowed        bool               `json:"post_comment_allowed"`
	PostCommentDisplayed      bool               `json:"post_comment_displayed"`
	PostFilesAllowed          bool               `json:"post_files_allowed"`
	PostVideoAllowed          bool               `json:"post_video_allowed"`
	PostVideoUrl              string             `json:"post_video_url"`
	PostPinned                bool               `json:"post_pinned"`
	RecImage1                 string             `json:"rec_image1"`
	RecImage2                 string             `json:"rec_image2"`
}

type CmsPostList struct {
	PostKey                   uint64             `json:"post_key"`
	PostSubtype               CmsPostSubtypeInfo `json:"post_subtype"`
	PostTitle                 string             `json:"post_title"`
	PostSubTitle              string             `json:"post_sub_title"`
	PostContentAuthor         string             `json:"post_content_author"`
	PostContentSources        string             `json:"post_content_sources"`
	PostPublishStart          string             `json:"post_publish_start"`
	PostPublishThru           string             `json:"post_publish_thru"`
	PostPinned                bool               `json:"post_pinned"`
	RecImage1                 string             `json:"rec_image1"`
}

func GetAllCmsPost(c *[]CmsPost, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              cms_post.* FROM 
			  cms_post WHERE 
			  cms_post.post_publish_start <= NOW() AND 
			  cms_post.post_publish_thru > NOW() AND 
			  cms_post.rec_status = 1 `
	var present bool
	var whereClause []string
	var condition string
	
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType"){
			whereClause = append(whereClause, "cms_post."+field + " = '" + value + "'")
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

func GetCmsPostIn(c *[]CmsPost, value []string, field string,limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				cms_post.* FROM 
				cms_post WHERE 
				cms_post.post_publish_start <= NOW() AND 
				cms_post.post_publish_thru > NOW() AND 
				cms_post.rec_status = 1 `
	query := query2 + " AND cms_post."+field+" IN(" + inQuery + ")"
	
	var present bool
	var whereClause []string
	var condition string
	
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType"){
			whereClause = append(whereClause, "cms_post."+field + " = " + value)
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

func GetCmsPost(c *CmsPost, key string) (int, error) {
	query := `SELECT cms_post.* FROM cms_post WHERE cms_post.post_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}