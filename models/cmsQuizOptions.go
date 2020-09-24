package models

import (
	"api/db"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CmsQuizOptionsInfo struct {
	QuizOptionKey          uint64     `json:"quiz_option_key"`
	QuizOptionLabel        string     `json:"quiz_option_label"`
	QuizOptionTitle        string     `json:"quiz_option_title"`
	QuizOptionScore        uint64     `json:"quiz_option_score"`
	QuizOptionDefault      uint8      `json:"quiz_option_default"`
}

type CmsQuizOptions struct {
	QuizOptionKey          uint64     `db:"quiz_option_key"         json:"quiz_option_key"`
	QuizQuestionKey        uint64     `db:"quiz_question_key"       json:"quiz_question_key"`
	QuizOptionLabel       *string     `db:"quiz_option_label"       json:"quiz_option_label"`
	QuizOptionTitle       *string     `db:"quiz_option_title"       json:"quiz_option_title"`
	QuizOptionScore       *uint64     `db:"quiz_option_score"       json:"quiz_option_score"`
	QuizOptionDefault     *uint8      `db:"quiz_option_default"     json:"quiz_option_default"`
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

func GetAllCmsQuizOptions(c *[]CmsQuizOptions, params map[string]string) (int, error) {
	query := `SELECT
              cms_quiz_options.* FROM 
			  cms_quiz_options `
	var present bool
	var whereClause []string
	var condition string
	
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType"){
			whereClause = append(whereClause, "cms_quiz_options."+field + " = '" + value + "'")
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

func GetCmsQuizOptionsIn(c *[]CmsQuizOptions, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				cms_quiz_options.* FROM 
				cms_quiz_options `
	query := query2 + " WHERE cms_quiz_options."+field+" IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}