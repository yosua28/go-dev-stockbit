package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CmsQuizHeaderData struct {
	QuizHeaderKey          uint64                  `json:"quiz_header_key"`
	QuizName               string                  `json:"quiz_name"`
	QuizDesc               string                  `json:"quiz_desc"`
	Questions             *[]CmsQuizQuestionInfo   `json:"questions,omitempty"`
}

type CmsQuizHeader struct {
	QuizHeaderKey          uint64     `db:"quiz_header_key"         json:"quiz_header_key"`
	QuizTypeKey           *uint64     `db:"quiz_type_key"           json:"quiz_type_key"`
	QuizName              *string     `db:"quiz_name"               json:"quiz_name"`
	QuizDesc              *string     `db:"quiz_desc"               json:"quiz_desc"`
	QuizPublishedStart    *string     `db:"quiz_published_start"    json:"quiz_published_start"`
	QuizPublishedThru     *string     `db:"quiz_published_thru"     json:"quiz_published_thru"`
	QuizMinimumPass       *uint64     `db:"quiz_minimum_pass"       json:"quiz_minimum_pass"`
	RecOrder              *uint64     `db:"rec_order"               json:"rec_order"`
	RecStatus              uint8      `db:"rec_status"              json:"rec_status"`
	RecCreatedDate        *string     `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy          *string     `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate       *string     `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy         *string     `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1             *string     `db:"rec_image1"              json:"rec_image1"`
	RecImage2             *string     `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus     *uint8      `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage      *uint64     `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate       *string     `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy         *string     `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate        *string     `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy          *string     `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1       *string     `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2       *string     `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3       *string     `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

func GetAllCmsQuizHeader(c *[]CmsQuizHeader, params map[string]string) (int, error) {
	query := `SELECT
              cms_quiz_header.* FROM 
			  cms_quiz_header `
	var present bool
	var whereClause []string
	var condition string
	
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType"){
			whereClause = append(whereClause, "cms_quiz_header."+field + " = '" + value + "'")
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

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetCmsQuizHeader(c *CmsQuizHeader, key string) (int, error) {
	query := `SELECT cms_quiz_header.* FROM cms_quiz_header WHERE cms_quiz_header.quiz_header_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}