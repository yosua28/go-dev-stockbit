package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsCustomer struct {
	CustomerKey            uint64  `db:"customer_key"                json:"customer_key"`
	IDCustomer             uint64  `db:"id_customer"                 json:"id_customer"`
	UnitHolderIDno         string  `db:"unit_holder_idno"            json:"unit_holder_idno"`
	FullName               string  `db:"full_name"                   json:"full_name"`
	SidNo                  *string `db:"sid_no"                      json:"sid_no"`
	InvestorType           string  `db:"investor_type"               json:"investor_type"`
	CustomerCategory       string  `db:"customer_category"           json:"customer_category"`
	ParticipantKey         *uint64 `db:"participant_key"             json:"participant_key"`
	CifSuspendFlag         uint8   `db:"cif_suspend_flag"            json:"cif_suspend_flag"`
	CifSuspendModifiedDate *string `db:"cif_suspend_modified_date"   json:"cif_suspend_modified_date"`
	CifSuspendReason       *string `db:"cif_suspend_reason"          json:"cif_suspend_reason"`
	OpenaccBranchKey       *uint64 `db:"openacc_branch_key"          json:"openacc_branch_key"`
	OpenaccAgentKey        *uint64 `db:"openacc_agent_key"           json:"openacc_agent_key"`
	OpenaccDate            *string `db:"openacc_date"                json:"openacc_date"`
	CloseaccBranchKey      *uint64 `db:"closeacc_branch_key"         json:"closeacc_branch_key"`
	CloseaccAgentKey       *uint64 `db:"closeacc_agent_key"          json:"closeacc_agent_key"`
	CloseaccDate           *string `db:"closeacc_date"               json:"closeacc_date"`
	FlagEmployee           uint8   `db:"flag_employee"               json:"flag_employee"`
	FlagGroup              uint8   `db:"flag_group"                  json:"flag_group"`
	EmployeeNumber         *string `db:"employee_number"             json:"employee_number"`
	EmployeeEmail          *string `db:"employee_email"              json:"employee_email"`
	EmployeeNotes          *string `db:"employee_notes"              json:"employee_notes"`
	ParentKey              *uint64 `db:"parent_key"                  json:"parent_key"`
	MargingFlag            uint8   `db:"merging_flag"                json:"merging_flag"`
	FirstName              *string `db:"first_name"                  json:"first_name"`
	MiddleName             *string `db:"middle_name"                 json:"middle_name"`
	LastName               *string `db:"last_name"                   json:"last_name"`
	ClientCode             *string `db:"client_code"                 json:"client_code"`
	TinNumber              *string `db:"tin_number"                  json:"tin_number"`
	TinIssuanceDate        *string `db:"tin_issuance_date"           json:"tin_issuance_date"`
	TinIssuanceCountry     *uint64 `db:"tin_issuance_country"        json:"tin_issuance_country"`
	FatcaStatus            *uint64 `db:"fatca_status"                json:"fatca_status"`
	RecOrder               *uint64 `db:"rec_order"                   json:"rec_order"`
	RecStatus              uint8   `db:"rec_status"                  json:"rec_status"`
	RecCreatedDate         *string `db:"rec_created_date"            json:"rec_created_date"`
	RecCreatedBy           *string `db:"rec_created_by"              json:"rec_created_by"`
	RecModifiedDate        *string `db:"rec_modified_date"           json:"rec_modified_date"`
	RecModifiedBy          *string `db:"rec_modified_by"             json:"rec_modified_by"`
	RecImage1              *string `db:"rec_image1"                  json:"rec_image1"`
	RecImage2              *string `db:"rec_image2"                  json:"rec_image2"`
	RecApprovalStatus      *uint8  `db:"rec_approval_status"         json:"rec_approval_status"`
	RecApprovalStage       *uint64 `db:"rec_approval_stage"          json:"rec_approval_stage"`
	RecApprovedDate        *string `db:"rec_approved_date"           json:"rec_approved_date"`
	RecApprovedBy          *string `db:"rec_approved_by"             json:"rec_approved_by"`
	RecDeletedDate         *string `db:"rec_deleted_date"            json:"rec_deleted_date"`
	RecDeletedBy           *string `db:"rec_deleted_by"              json:"rec_deleted_by"`
	RecAttributeID1        *string `db:"rec_attribute_id1"           json:"rec_attribute_id1"`
	RecAttributeID2        *string `db:"rec_attribute_id2"           json:"rec_attribute_id2"`
	RecAttributeID3        *string `db:"rec_attribute_id3"           json:"rec_attribute_id3"`
}

type CustomerIndividuInquiry struct {
	CustomerKey      uint64 `db:"customer_key"                json:"customer_key"`
	Cif              string `db:"cif"                         json:"cif"`
	FullName         string `db:"full_name"                   json:"full_name"`
	DateBirth        string `db:"date_birth"                  json:"date_birth"`
	IDcardNo         string `db:"ktp"                         json:"ktp"`
	PhoneMobile      string `db:"phone_mobile"                json:"phone_mobile"`
	SidNo            string `db:"sid"                         json:"sid"`
	CifSuspendFlag   string `db:"cif_suspend_flag"            json:"cif_suspend_flag"`
	MotherMaidenName string `db:"mother_maiden_name"          json:"mother_maiden_name"`
}

type CustomerInstituionInquiry struct {
	CustomerKey    uint64 `db:"customer_key"                json:"customer_key"`
	Cif            string `db:"cif"                         json:"cif"`
	FullName       string `db:"full_name"                   json:"full_name"`
	Npwp           string `db:"npwp"                        json:"npwp"`
	Institusion    string `db:"institution"                 json:"institution"`
	SidNo          string `db:"sid"                         json:"sid"`
	CifSuspendFlag string `db:"cif_suspend_flag"            json:"cif_suspend_flag"`
}

type DetailCustomerIndividuInquiry struct {
	Header       CustomerIndividuInquiry `json:"header"`
	PersonalData *[]OaCustomer           `json:"personal_data"`
}

type DetailCustomerInstitutionInquiry struct {
	Header       CustomerInstituionInquiry `json:"header"`
	PersonalData *[]OaCustomer             `json:"personal_data"`
}

type DetailCustomerInquiryResponse struct {
	Header       DetailHeaderCustomerInquiry `json:"header"`
	PersonalData *[]OaCustomer               `json:"personal_data"`
}

type DetailCustomerInquiry struct {
	InvestorType     string  `db:"investor_type"               json:"investor_type"`
	CustomerKey      uint64  `db:"customer_key"                json:"customer_key"`
	Cif              string  `db:"cif"                         json:"cif"`
	FullName         string  `db:"full_name"                   json:"full_name"`
	DateBirth        *string `db:"date_birth"                  json:"date_birth"`
	IDcardNo         *string `db:"ktp"                         json:"ktp"`
	PhoneMobile      *string `db:"phone_mobile"                json:"phone_mobile"`
	SidNo            string  `db:"sid"                         json:"sid"`
	CifSuspendFlag   string  `db:"cif_suspend_flag"            json:"cif_suspend_flag"`
	MotherMaidenName *string `db:"mother_maiden_name"          json:"mother_maiden_name"`
	Npwp             *string `db:"npwp"                        json:"npwp"`
	Institusion      *string `db:"institution"                 json:"institution"`
}

type DetailHeaderCustomerInquiry struct {
	InvestorType     string  `json:"investor_type"`
	CustomerKey      uint64  `json:"customer_key"`
	Cif              string  `json:"cif"`
	FullName         string  `json:"full_name"`
	DateBirth        *string `json:"date_birth,omitempty"`
	IDcardNo         *string `json:"ktp,omitempty"`
	PhoneMobile      *string `json:"phone_mobile,omitempty"`
	SidNo            string  `json:"sid"`
	CifSuspendFlag   string  `json:"cif_suspend_flag"`
	MotherMaidenName *string `json:"mother_maiden_name,omitempty"`
	Npwp             *string `json:"npwp,omitempty"`
	Institusion      *string `json:"institution,omitempty"`
}

func GetMsCustomerIn(c *[]MsCustomer, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_customer.* FROM 
				ms_customer `
	query := query2 + " WHERE ms_customer." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsCustomer(c *MsCustomer, key string) (int, error) {
	query := `SELECT ms_customer.* FROM ms_customer WHERE ms_customer.customer_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func CreateMsCustomer(params map[string]string) (int, error, string) {
	query := "INSERT INTO ms_customer"
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

func UpdateMsCustomer(params map[string]string) (int, error) {
	query := "UPDATE ms_customer SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "customer_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE customer_key = " + params["customer_key"]
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

func GetMsCustomerByClientCode(c *MsCustomer, clientCode string) (int, error) {
	query := `SELECT ms_customer.* FROM ms_customer WHERE ms_customer.rec_status = 1 AND ms_customer.client_code = ` + clientCode
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetLastUnitHolder(c *MsCustomer, value string) (int, error) {
	query := `SELECT ms_customer.* FROM ms_customer 
	WHERE ms_customer.unit_holder_idno LIKE '` + value + `%' AND ms_customer.rec_status = 1
	ORDER BY unit_holder_idno DESC LIMIT 1`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func AdminGetAllCustomerIndividuInquery(c *[]CustomerIndividuInquiry, limit uint64, offset uint64, params map[string]string, paramsLike map[string]string, nolimit bool) (int, error) {
	query := `SELECT 
				c.customer_key AS customer_key, 
				c.unit_holder_idno AS cif, 
				c.full_name AS full_name, 
				DATE_FORMAT(pd.date_birth, '%d %M %Y') AS date_birth, 
				pd.idcard_no AS ktp, 
				pd.phone_mobile AS phone_mobile, 
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN c.cif_suspend_flag = 0 THEN "Tidak"
					ELSE "Ya"
				END) AS cif_suspend_flag, 
				pd.mother_maiden_name AS mother_maiden_name  
			FROM ms_customer AS c 
			INNER JOIN (SELECT MAX(oa_request_key) AS oa_request_key, customer_key FROM oa_request WHERE rec_status = 1 GROUP BY customer_key) 
			AS t2 ON c.customer_key = t2.customer_key
			INNER JOIN oa_request AS r ON c.customer_key = r.customer_key AND r.oa_request_key = t2.oa_request_key
			INNER JOIN oa_personal_data AS pd ON pd.oa_request_key = r.oa_request_key
			WHERE c.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
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

func CountAdminGetAllCustomerIndividuInquery(c *CountData, params map[string]string, paramsLike map[string]string) (int, error) {
	query := `SELECT 
				count(c.customer_key) AS count_data 
			FROM ms_customer AS c
			INNER JOIN oa_request AS r ON c.customer_key = r.customer_key
			INNER JOIN oa_personal_data AS pd ON pd.oa_request_key = r.oa_request_key
			WHERE c.rec_status = 1`

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
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

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetAllCustomerInstitutionInquery(c *[]CustomerInstituionInquiry, limit uint64, offset uint64, params map[string]string, paramsLike map[string]string, nolimit bool) (int, error) {
	query := `SELECT 
				c.customer_key AS customer_key, 
				c.unit_holder_idno AS cif, 
				c.full_name AS full_name, 
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN c.cif_suspend_flag = 0 THEN "Tidak"
					ELSE "Ya"
				END) AS cif_suspend_flag, 
				pd.npwp_no AS npwp, 
				pd.insti_full_name AS institution 
			FROM ms_customer AS c 
			INNER JOIN (SELECT MAX(oa_request_key) AS oa_request_key, customer_key FROM oa_request WHERE rec_status = 1 GROUP BY customer_key) 
			AS t2 ON c.customer_key = t2.customer_key
			INNER JOIN oa_request AS r ON c.customer_key = r.customer_key AND r.oa_request_key = t2.oa_request_key
			INNER JOIN oa_institution_data AS pd ON pd.oa_request_key = r.oa_request_key
			WHERE c.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
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

func CountAdminGetAllCustomerInstitutionInquery(c *CountData, params map[string]string, paramsLike map[string]string) (int, error) {
	query := `SELECT 
				count(c.customer_key) AS count_data 
			FROM ms_customer AS c
			INNER JOIN oa_request AS r ON c.customer_key = r.customer_key
			INNER JOIN oa_institution_data AS pd ON pd.oa_request_key = r.oa_request_key
			WHERE c.rec_status = 1`

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
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

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetHeaderCustomerIndividu(c *CustomerIndividuInquiry, customerKey string) (int, error) {
	query := `SELECT 
				c.customer_key AS customer_key, 
				c.unit_holder_idno AS cif, 
				c.full_name AS full_name, 
				DATE_FORMAT(pd.date_birth, '%d %M %Y') AS date_birth, 
				pd.idcard_no AS ktp, 
				pd.phone_mobile AS phone_mobile, 
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN c.cif_suspend_flag = 0 THEN "Tidak"
					ELSE "Ya"
				END) AS cif_suspend_flag, 
				pd.mother_maiden_name AS mother_maiden_name  
			FROM ms_customer AS c 
			INNER JOIN (SELECT MAX(oa_request_key) AS oa_request_key, customer_key FROM oa_request WHERE rec_status = 1 GROUP BY customer_key) 
			AS t2 ON c.customer_key = t2.customer_key
			INNER JOIN oa_request AS r ON c.customer_key = r.customer_key AND r.oa_request_key = t2.oa_request_key
			INNER JOIN oa_personal_data AS pd ON pd.oa_request_key = r.oa_request_key
			WHERE c.rec_status = 1 AND c.investor_type = 263 AND c.customer_key = ` + customerKey

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetHeaderCustomerInstitution(c *CustomerInstituionInquiry, customerKey string) (int, error) {
	query := `SELECT 
				c.customer_key AS customer_key, 
				c.unit_holder_idno AS cif, 
				c.full_name AS full_name, 
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN c.cif_suspend_flag = 0 THEN "Tidak"
					ELSE "Ya"
				END) AS cif_suspend_flag, 
				pd.npwp_no AS npwp, 
				pd.insti_full_name AS institution 
			FROM ms_customer AS c 
			INNER JOIN (SELECT MAX(oa_request_key) AS oa_request_key, customer_key FROM oa_request WHERE rec_status = 1 GROUP BY customer_key) 
			AS t2 ON c.customer_key = t2.customer_key
			INNER JOIN oa_request AS r ON c.customer_key = r.customer_key AND r.oa_request_key = t2.oa_request_key
			INNER JOIN oa_institution_data AS pd ON pd.oa_request_key = r.oa_request_key
			WHERE c.rec_status = 1 AND c.investor_type = 264 AND c.customer_key = ` + customerKey
	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetHeaderDetailCustomer(c *DetailCustomerInquiry, customerKey string) (int, error) {
	query := `SELECT 
				c.investor_type AS investor_type,
				c.customer_key AS customer_key, 
				c.unit_holder_idno AS cif, 
				c.full_name AS full_name, 
				DATE_FORMAT(pd.date_birth, '%d %M %Y') AS date_birth, 
				pd.idcard_no AS ktp, 
				pd.phone_mobile AS phone_mobile, 
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN c.cif_suspend_flag = 0 THEN "Tidak"
					ELSE "Ya"
				END) AS cif_suspend_flag, 
				pd.mother_maiden_name AS mother_maiden_name,
				id.npwp_no AS npwp, 
				id.insti_full_name AS institution   
			FROM ms_customer AS c
			INNER JOIN (SELECT MAX(oa_request_key) AS oa_request_key, customer_key FROM oa_request WHERE rec_status = 1 GROUP BY customer_key) 
			AS t2 ON c.customer_key = t2.customer_key
			INNER JOIN oa_request AS r ON c.customer_key = r.customer_key AND r.oa_request_key = t2.oa_request_key
			LEFT JOIN oa_personal_data AS pd ON pd.oa_request_key = r.oa_request_key
			LEFT JOIN oa_institution_data AS id ON id.oa_request_key = r.oa_request_key
			WHERE c.rec_status = 1 AND c.customer_key = ` + customerKey
	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
