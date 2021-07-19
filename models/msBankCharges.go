package models

import (
	"api/db"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type MsBankCharges struct {
	BchargesKey       uint64           `db:"bcharges_key"           json:"bcharges_key"`
	BankNetworkType   uint64           `db:"bank_network_type"      json:"bank_network_type"`
	BankKey           *uint64          `db:"bank_key"               json:"bank_key"`
	CustodianKey      *uint64          `db:"custodian_key"          json:"custodian_key"`
	MinNominalTrx     *decimal.Decimal `db:"min_nominal_trx"        json:"min_nominal_trx"`
	ValueType         uint64           `db:"value_type"             json:"value_type"`
	ChargesValue      *decimal.Decimal `db:"charges_value"          json:"charges_value"`
	RecOrder          *uint64          `db:"rec_order"              json:"rec_order"`
	RecStatus         uint8            `db:"rec_status"             json:"rec_status"`
	RecCreatedDate    *string          `db:"rec_created_date"       json:"rec_created_date"`
	RecCreatedBy      *string          `db:"rec_created_by"         json:"rec_created_by"`
	RecModifiedDate   *string          `db:"rec_modified_date"      json:"rec_modified_date"`
	RecModifiedBy     *string          `db:"rec_modified_by"        json:"rec_modified_by"`
	RecImage1         *string          `db:"rec_image1"             json:"rec_image1"`
	RecImage2         *string          `db:"rec_image2"             json:"rec_image2"`
	RecApprovalStatus *uint8           `db:"rec_approval_status"    json:"rec_approval_status"`
	RecApprovalStage  *uint64          `db:"rec_approval_stage"     json:"rec_approval_stage"`
	RecApprovedDate   *string          `db:"rec_approved_date"      json:"rec_approved_date"`
	RecApprovedBy     *string          `db:"rec_approved_by"        json:"rec_approved_by"`
	RecDeletedDate    *string          `db:"rec_deleted_date"       json:"rec_deleted_date"`
	RecDeletedBy      *string          `db:"rec_deleted_by"         json:"rec_deleted_by"`
	RecAttributeID1   *string          `db:"rec_attribute_id1"      json:"rec_attribute_id1"`
	RecAttributeID2   *string          `db:"rec_attribute_id2"      json:"rec_attribute_id2"`
	RecAttributeID3   *string          `db:"rec_attribute_id3"      json:"rec_attribute_id3"`
}

func GetMsBankCharges(c *MsBankCharges, key string) (int, error) {
	query := `SELECT ms_bank_charges.* FROM ms_bank_charges WHERE ms_bank_charges.rec_status = '1' 
	AND ms_bank_charges.bcharges_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

type ListBankChargesAdmin struct {
	BchargesKey     uint64           `db:"bcharges_key"            json:"bcharges_key"`
	BankNetworkType string           `db:"bank_network_type"       json:"bank_network_type"`
	BankName        *string          `db:"bank_name"               json:"bank_name"`
	CustodianName   *string          `db:"custodian_name"          json:"custodian_name"`
	MinNominalTrx   *decimal.Decimal `db:"min_nominal_trx"         json:"min_nominal_trx"`
	ValueType       uint64           `db:"value_type"              json:"value_type"`
	ChargesValue    *decimal.Decimal `db:"charges_value"           json:"charges_value"`
}

func AdminGetListBankCharges(c *[]ListBankChargesAdmin, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
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

	if searchLike != "" {
		condition += " AND"
		condition += " (net.lkp_name like '%" + searchLike + "%' OR"
		condition += " b.bank_name like '%" + searchLike + "%' OR"
		condition += " cb.custodian_short_name like '%" + searchLike + "%' OR"
		condition += " bc.min_nominal_trx like '%" + searchLike + "%' OR"
		condition += " bc.value_type like '%" + searchLike + "%' OR"
		condition += " bc.charges_value like '%" + searchLike + "%')"
	}

	query := `SELECT
				bc.bcharges_key,
				net.lkp_name AS bank_network_type,
				b.bank_name,
				cb.custodian_short_name AS custodian_name,
				bc.min_nominal_trx, 
				bc.value_type,
				bc.charges_value 
			FROM ms_bank_charges AS bc
			INNER JOIN gen_lookup AS net ON net.lookup_key = bc.bank_network_type
			LEFT JOIN ms_bank AS b ON b.bank_key = bc.bank_key
			LEFT JOIN ms_custodian_bank AS cb ON cb.custodian_key = bc.custodian_key
			WHERE bc.rec_status = 1 ` + condition

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	if !nolimit {
		limitOffset += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			limitOffset += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += orderCondition + limitOffset

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CountAdminGetListBankCharges(c *CountData, params map[string]string, searchLike string) (int, error) {
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
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

	if searchLike != "" {
		condition += " AND"
		condition += " (net.lkp_name like '%" + searchLike + "%' OR"
		condition += " b.bank_name like '%" + searchLike + "%' OR"
		condition += " cb.custodian_short_name like '%" + searchLike + "%' OR"
		condition += " bc.min_nominal_trx like '%" + searchLike + "%' OR"
		condition += " bc.value_type like '%" + searchLike + "%' OR"
		condition += " bc.charges_value like '%" + searchLike + "%')"
	}

	query := `SELECT
				count(bc.bcharges_key) AS count_data 
			FROM ms_bank_charges AS bc
			INNER JOIN gen_lookup AS net ON net.lookup_key = bc.bank_network_type
			LEFT JOIN ms_bank AS b ON b.bank_key = bc.bank_key
			LEFT JOIN ms_custodian_bank AS cb ON cb.custodian_key = bc.custodian_key
			WHERE bc.rec_status = 1 ` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateMsBankCharges(params map[string]string) (int, error) {
	query := "INSERT INTO ms_bank_charges"
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

func UpdateMsBankCharges(params map[string]string) (int, error) {
	query := "UPDATE ms_bank_charges SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "bcharges_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE bcharges_key = " + params["bcharges_key"]
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}
