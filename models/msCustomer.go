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
	query := `SELECT ms_customer.* FROM ms_customer WHERE ms_customer.rec_status = 1 AND ms_customer.customer_key = ` + key
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
