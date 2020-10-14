package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type OaRequest struct {
	OaRequestKey      uint64  `db:"oa_request_key"             json:"oa_request_key"`
	OaRequestType     *uint64 `db:"oa_request_type"            json:"oa_request_type"`
	OaEntryStart      string  `db:"oa_entry_start"             json:"oa_entry_start"`
	OaEntryEnd        string  `db:"oa_entry_end"               json:"oa_entry_end"`
	Oastatus          *uint64 `db:"oa_status"                  json:"oa_status"`
	UserLoginKey      *uint64 `db:"user_login_key"             json:"user_login_key"`
	CustomerKey       *uint64 `db:"customer_key"               json:"customer_key"`
	SalesCode         *string `db:"sales_code"                 json:"sales_code"`
	Check1Date        *string `db:"check1_date"                json:"check1_date"`
	Check1Flag        *uint8  `db:"check1_flag"                json:"check1_flag"`
	Check1References  *string `db:"check1_references"          json:"check1_references"`
	Check1Notes       *string `db:"check1_notes"               json:"check1_notes"`
	Check2Date        *string `db:"check2_date"                json:"check2_date"`
	Check2Flag        *uint8  `db:"check2_flag"                json:"check2_flag"`
	Check2References  *string `db:"check2_references"          json:"check2_references"`
	Check2Notes       *string `db:"check2_notes"               json:"check2_notes"`
	OaRiskLevel       *uint64 `db:"oa_risk_level"              json:"oa_risk_level"`
	RecOrder          *uint64 `db:"rec_order"                  json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                 json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
}

type OaRequestListResponse struct {
	OaRequestKey uint64 `json:"oa_request_key"`
	OaEntryStart string `json:"oa_entry_start"`
	OaEntryEnd   string `json:"oa_entry_end"`
	Oastatus     string `json:"oa_status"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_mobile"`
	DateBirth    string `json:"date_birth"`
	FullName     string `json:"full_name"`
	IDCardNo     string `json:"idcard_no"`
}

type OaRequestCountData struct {
	CountData int `db:"count_data"             json:"count_data"`
}

type OaRequestDetailResponse struct {
	OaRequestKey        uint64               `json:"oa_request_key"`
	OaRequestType       *string              `json:"oa_request_type"`
	OaRiskLevel         *string              `json:"oa_risk_level"`
	OaEntryStart        string               `json:"oa_entry_start"`
	OaEntryEnd          string               `json:"oa_entry_end"`
	Oastatus            string               `json:"oa_status"`
	EmailAddress        string               `json:"email_address"`
	PhoneNumber         string               `json:"phone_mobile"`
	DateBirth           string               `json:"date_birth"`
	FullName            string               `json:"full_name"`
	IDCardNo            string               `json:"idcard_no"`
	Nationality         *string              `json:"nationality"`
	Gender              *string              `json:"gender"`
	PlaceBirth          string               `json:"place_birth"`
	MaritalStatus       *string              `json:"marital_status"`
	PhoneHome           string               `json:"phone_home"`
	Religion            *string              `json:"religion"`
	Education           *string              `json:"education"`
	PicKtp              *string              `json:"pic_ktp"`
	PicSelfie           *string              `json:"pic_selfie"`
	PicSelfieKtp        *string              `json:"pic_selfie_ktp"`
	OccupJob            *string              `json:"occup_job"`
	OccupCompany        *string              `json:"occup_company"`
	OccupPosition       *string              `json:"occup_position"`
	OccupPhone          *string              `json:"occup_phone"`
	OccupWebURL         *string              `json:"occup_web_url"`
	AnnualIncome        *string              `json:"annual_income"`
	SourceofFund        *string              `json:"sourceof_fund"`
	InvesmentObjectives *string              `json:"invesment_objectives"`
	Correspondence      *string              `json:"correspondence"`
	MotherMaidenName    string               `json:"mother_maiden_name"`
	BeneficialRelation  *string              `json:"beneficial_relation"`
	BeneficialFullName  *string              `json:"beneficial_full_name"`
	OccupBusinessFields *string              `json:"occup_business_fields"`
	IDcardAddress       Address              `json:"idcard_address"`
	DomicileAddress     Address              `json:"domicile_address"`
	OccupAddressKey     Address              `json:"occup_address_key"`
	BankAccount         BankAccount          `json:"bank_account"`
	Relation            Relation             `json:"relation"`
	Emergency           Emergency            `json:"emergency"`
	RiskProfile         []AdminOaRiskProfile `json:"risk_profile"`
	RiskProfileQuiz     []RiskProfileQuiz    `json:"risk_profile_quiz"`
}

type Address struct {
	Address    *string `json:"address"`
	Kabupaten  *string `json:"kabupaten"`
	Kecamatan  *string `json:"kecamatan"`
	PostalCode *string `json:"postal_code"`
}

type BankAccount struct {
	BankKey           *string `json:"bank_key"`
	AccountNo         *string `json:"account_no"`
	AccountHolderName *string `json:"account_holder_name"`
	BranchName        *string `json:"branch_name"`
	BankAccountType   *string `json:"bank_account_type"`
}

type Relation struct {
	RelationType           *string `json:"relation_type"`
	RelationFullName       *string `json:"relation_full_name"`
	RelationOccupation     *string `json:"relation_occupation"`
	RelationBusinessFields *string `json:"relation_business_fields"`
}

type Emergency struct {
	EmergencyFullName *string `json:"emergency_full_name"`
	EmergencyRelation *string `json:"emergency_relation"`
	EmergencyPhoneNo  *string `json:"emergency_phone_no"`
}

type RiskProfileQuiz struct {
	RiskProfileQuizKey  uint64               `json:"risk_profile_quiz_key"`
	QuizOptionUser      CmsQuizOptionsInfo   `json:"quiz_option_user"`
	QuizOptionScoreUser float32              `json:"quiz_option_score_user"`
	QuizQuestionKey     uint64               `json:"quiz_question_key"`
	HeaderQuizName      string               `json:"header_quiz_name"`
	QuizTitle           string               `json:"quiz_title"`
	Options             []CmsQuizOptionsInfo `json:"options"`
}

func CreateOaRequest(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_request"
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
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func GetAllOaRequest(c *[]OaRequest, limit uint64, offset uint64, nolimit bool, params map[string]string) (int, error) {
	query := `SELECT
              oa_request.*
			  FROM oa_request`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "oa_request."+field+" = '"+value+"'")
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

func GetOaRequest(c *OaRequest, key string) (int, error) {
	query := `SELECT oa_request.* FROM oa_request WHERE oa_request.rec_status = 1 AND oa_request.oa_request_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetCountOaRequest(c *OaRequestCountData, params map[string]string) (int, error) {
	query := `SELECT
              count(oa_request.oa_request_key) as count_data
			  FROM oa_request`

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "oa_request."+field+" = '"+value+"'")
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

func GetOaRequestsIn(c *[]OaRequest, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				oa_request.* FROM 
				oa_request `
	query := query2 + " WHERE oa_request.rec_status = 1 AND oa_request." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateOaRequest(params map[string]string) (int, error) {
	query := "UPDATE oa_request SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "oa_request_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE oa_request_key = " + params["oa_request_key"]
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
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
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetAllOaRequestApproval3(c *[]OaRequest, limit uint64, offset uint64,
	nolimit bool, params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := `SELECT
              oa_request.*
			  FROM oa_request`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "oa_request."+field+" = '"+value+"'")
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

	if len(valueIn) > 0 {
		if len(whereClause) < 1 {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " WHERE oa_request." + fieldIn + " IN(" + inQuery + ")"
			}
		} else {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " AND oa_request." + fieldIn + " IN(" + inQuery + ")"
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

func GetAllOaRequestDoTransaction(c *[]OaRequest, limit uint64, offset uint64, nolimit bool,
	params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := `SELECT oa_request.* 
				FROM 
				oa_request AS oa_request 
				INNER JOIN ms_customer AS cus ON oa_request.customer_key = cus.customer_key 
				INNER JOIN 
				( 
					SELECT customer_key 
					FROM tr_transaction 
					WHERE rec_status = 1 
					GROUP BY customer_key 
				) tr ON tr.customer_key = cus.customer_key `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "oa_request."+field+" = '"+value+"'")
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

	if len(valueIn) > 0 {
		if len(whereClause) < 1 {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " WHERE oa_request." + fieldIn + " IN(" + inQuery + ")"
			}
		} else {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " AND oa_request." + fieldIn + " IN(" + inQuery + ")"
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

func GetCountOaRequestDoTransaction(c *OaRequestCountData, params map[string]string,
	valueIn []string, fieldIn string) (int, error) {
	query := `SELECT count(oa_request.oa_request_key) as count_data
				FROM 
				oa_request AS oa_request 
				INNER JOIN ms_customer AS cus ON oa_request.customer_key = cus.customer_key 
				INNER JOIN 
				( 
					SELECT customer_key 
					FROM tr_transaction 
					WHERE rec_status = 1 
					GROUP BY customer_key 
				) tr ON tr.customer_key = cus.customer_key `

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "oa_request."+field+" = '"+value+"'")
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

	if len(valueIn) > 0 {
		if len(whereClause) < 1 {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " WHERE oa_request." + fieldIn + " IN(" + inQuery + ")"
			}
		} else {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " AND oa_request." + fieldIn + " IN(" + inQuery + ")"
			}
		}
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

func UpdateOaRequestByKeyIn(params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := "UPDATE oa_request SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}

	inQuery := strings.Join(valueIn, ",")
	query += " WHERE oa_request." + fieldIn + " IN(" + inQuery + ")"

	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
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
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
