package models

import (
	"api/db"
	"net/http"
	_ "database/sql"
	_ "strconv"

	log "github.com/sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

type ParamsRiskProfileQuiz struct {
	OaRequestKey           string     `db:"oa_request_key"`
	QuizQuestionKey        string     `db:"quiz_question_key"`
	QuizOptionKey          string     `db:"quiz_option_key"`
	QuizOptionScore        string     `db:"quiz_option_score"`
	RecStatus              string     `db:"rec_status"`
}

type OaRiskProfileQuiz struct {
	RiskProfileQuizKey     uint64     `db:"risk_profile_quiz_key"   json:"risk_profile_quiz_key"`
	OaRequestKey          *uint64     `db:"oa_request_key"          json:"oa_request_key"`
	QuizQuestionKey       *uint64     `db:"quiz_question_key"       json:"quiz_question_key"`
	QuizOptionKey         *uint64     `db:"quiz_option_key"         json:"quiz_option_key"`
	QuizOptionScore       *uint64     `db:"quiz_option_score"       json:"quiz_option_score"`
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

func CreateOaRiskProfileQuiz(params map[string]string) (int, error){
	query := "INSERT INTO oa_risk_profile_quiz"
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
	query += "("+fields + ") VALUES(" + values + ")"
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

func CreateMultipleOaRiskProfileQuiz(params []interface{}) (int, error){

	q := `INSERT INTO oa_risk_profile_quiz (
		oa_request_key, 
		quiz_question_key,
		quiz_option_key,
		quiz_option_score,
		rec_status) VALUES `

	for i := 0; i < len(params); i++ {
		q += "(?)"
		if i < (len(params)-1){
			q += ","
		}
	}
	log.Info(q)
	query, args, err := sqlx.In(q,params...)
	if err != nil {
		return http.StatusBadGateway, err
	}

	query = db.Db.Rebind(query)
	_, err = db.Db.Query(query, args...)
	if err != nil {
		log.Error(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}