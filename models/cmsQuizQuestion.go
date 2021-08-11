package models

import (
	"api/db"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CmsQuizQuestionInfo struct {
	QuizQuestionKey   uint64                `json:"quiz_question_key"`
	QuizTitle         string                `json:"quiz_title"`
	FileImageAllowed  bool                  `json:"file_image_allowed"`
	QuizOptionType    string                `json:"quiz_opetion_type"`
	QuizOptionDefault string                `json:"quiz_opetion_default"`
	Options           *[]CmsQuizOptionsInfo `json:"options,omitempty"`
}

type CmsQuizQuestion struct {
	QuizQuestionKey   uint64  `db:"quiz_question_key"       json:"quiz_question_key"`
	QuizHeaderKey     uint64  `db:"quiz_header_key"         json:"quiz_header_key"`
	QuizTitle         *string `db:"quiz_title"              json:"quiz_title"`
	FileImageAllowed  uint8   `db:"file_image_allowed"      json:"file_image_allowed"`
	QuizOptionType    string  `db:"quiz_option_type"        json:"quiz_option_type"`
	QuizOptionDefault *string `db:"quiz_option_default"    json:"quiz_option_default"`
	RecOrder          *uint64 `db:"rec_order"               json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"              json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"              json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

type QuestionOptionQuiz struct {
	QuizQuestionKey string `db:"quiz_question_key"       json:"quiz_question_key"`
	QuizOptionKey   string `db:"quiz_option_key"         json:"quiz_option_key"`
	QuizOptionScore uint64 `db:"quiz_option_score"       json:"quiz_option_score"`
}

func GetAllCmsQuizQuestion(c *[]CmsQuizQuestion, params map[string]string) (int, error) {
	query := `SELECT
              cms_quiz_question.* FROM 
			  cms_quiz_question `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "cms_quiz_question."+field+" = '"+value+"'")
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

func AdminGetQuestionOptionQuiz(c *[]QuestionOptionQuiz, optionKey []string) (int, error) {
	inQuery := strings.Join(optionKey, ",")
	query := `SELECT 
				q.quiz_question_key AS quiz_question_key,
				o.quiz_option_key AS quiz_option_key,
				(CASE
					WHEN o.quiz_option_score IS NULL THEN 0
					ELSE o.quiz_option_score
				END) AS quiz_option_score 
			FROM cms_quiz_question AS q 
			INNER JOIN cms_quiz_options AS o ON o.quiz_question_key = q.quiz_question_key 
			AND q.quiz_header_key = 2 AND q.rec_status = 1 
			AND o.quiz_option_key IN(` + inQuery + `) ORDER BY q.rec_order ASC `
	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type QuestionOptionCustomer struct {
	QuizQuestionKey uint64  `db:"quiz_question_key"     json:"quiz_question_key"`
	QuizTitle       string  `db:"quiz_title"            json:"quiz_title"`
	QuizOptionKey   uint64  `db:"quiz_option_key"       json:"quiz_option_key"`
	QuizOptionTitle string  `db:"quiz_option_title"     json:"quiz_option_title"`
	UserAnswer      *string `db:"user_answer"           json:"user_answer"`
	IsCheck         string  `db:"is_check"              json:"is_check"`
}

func GetListQuestionOptionCustomer(c *[]QuestionOptionCustomer, oaReqKey string) (int, error) {
	query := `SELECT
				q.quiz_question_key,
				q.quiz_title,
				o.quiz_option_key,
				o.quiz_option_title,
				rp.quiz_option_key AS user_answer, 
				(CASE 
				WHEN o.quiz_option_key = rp.quiz_option_key THEN "1"
				ELSE "0"
				END) AS is_check 
			FROM cms_quiz_question AS q
			INNER JOIN cms_quiz_options AS o ON q.quiz_question_key = o.quiz_question_key
			LEFT JOIN oa_risk_profile_quiz AS rp ON rp.quiz_question_key = o.quiz_question_key AND rp.oa_request_key = "` + oaReqKey + `"
			WHERE q.quiz_header_key = 2 AND q.rec_status = 1 AND o.rec_status = 1
			ORDER BY q.quiz_question_key, o.quiz_option_key ASC`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
