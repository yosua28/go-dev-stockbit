package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TrTransaction struct {
	TransactionKey    uint64   `db:"transaction_key"           json:"transaction_key"`
	IDTransaction     *uint64  `db:"id_transaction"            json:"id_transaction"`
	BranchKey         *uint64  `db:"branch_key"                json:"branch_key"`
	AgentKey          *uint64  `db:"agent_key"                 json:"agent_key"`
	CustomerKey       uint64   `db:"customer_key"              json:"customer_key"`
	ProductKey        uint64   `db:"product_key"               json:"product_key"`
	TransStatusKey    uint64   `db:"trans_status_key"          json:"trans_status_key"`
	TransDate         string   `db:"trans_date"                json:"trans_date"`
	TransTypeKey      uint64   `db:"trans_type_key"            json:"trans_type_key"`
	TrxCode           *uint64  `db:"trx_code"                  json:"trx_code"`
	NavDate           string   `db:"nav_date"                  json:"nav_date"`
	EntryMode         *uint64  `db:"entry_mode"                json:"entry_mode"`
	TransCalcMethod   *uint64  `db:"trans_calc_method"         json:"trans_calc_method"`
	TransAmount       float32  `db:"trans_amount"              json:"trans_amount"`
	TransUnit         float32  `db:"trans_unit"                json:"trans_unit"`
	TransUnitPercent  *float32 `db:"trans_unit_percent"        json:"trans_unit_percent"`
	FlagRedemtAll     *uint8   `db:"flag_redempt_all"          json:"flag_redempt_all"`
	FlagNewSub        *uint8   `db:"flag_newsub"               json:"flag_newsub"`
	TransFeePercent   float32  `db:"trans_fee_percent"         json:"trans_fee_percent"`
	TransFeeAmount    float32  `db:"trans_fee_amount"          json:"trans_fee_amount"`
	ChargesFeeAmount  float32  `db:"charges_fee_amount"        json:"charges_fee_amount"`
	ServicesFeeAmount float32  `db:"services_fee_amount"       json:"services_fee_amount"`
	TotalAmount       float32  `db:"total_amount"              json:"total_amount"`
	SettlementDate    *string  `db:"settlement_date"           json:"settlement_date"`
	TransBankAccNo    *string  `db:"trans_bank_accno"          json:"trans_bank_accno"`
	TransBankaccName  *string  `db:"trans_bankacc_name"        json:"trans_bankacc_name"`
	TransBankKey      *uint64  `db:"trans_bank_key"            json:"trans_bank_key"`
	TransRemarks      *string  `db:"trans_remarks"             json:"trans_remarks"`
	TransReferences   *string  `db:"trans_references"          json:"trans_references"`
	PromoCode         *string  `db:"promo_code"                json:"promo_code"`
	SalesCode         *string  `db:"sales_code"                json:"sales_code"`
	RiskWaiver        uint8    `db:"risk_waiver"               json:"risk_waiver"`
	AddtoAutoInvest   *uint8   `db:"addto_auto_invest"         json:"addto_auto_invest"`
	TransSource       *uint64  `db:"trans_source"              json:"trans_source"`
	FileUploadDate    *string  `db:"file_upload_date"          json:"file_upload_date"`
	PaymentMethod     *uint64  `db:"payment_method"            json:"payment_method"`
	Check1Date        *string  `db:"check1_date"               json:"check1_date"`
	Check1Flag        *uint8   `db:"check1_flag"               json:"check1_flag"`
	Check1References  *string  `db:"check1_references"         json:"check1_references"`
	Check1Notes       *string  `db:"check1_notes"              json:"check1_notes"`
	Check2Date        *string  `db:"check2_date"               json:"check2_date"`
	Check2Flag        *uint8   `db:"check2_flag"               json:"check2_flag"`
	Check2References  *string  `db:"check2_references"         json:"check2_references"`
	Check2Notes       *string  `db:"check2_notes"              json:"check2_notes"`
	TrxRiskLevel      *uint64  `db:"trx_risk_level"            json:"trx_risk_level"`
	ProceedDate       *string  `db:"proceed_date"              json:"proceed_date"`
	ProceedAmount     *float32 `db:"proceed_amount"            json:"proceed_amount"`
	SentDate          *string  `db:"sent_date"                 json:"sent_date"`
	SentReferences    *string  `db:"sent_references"           json:"sent_references"`
	ConfirmedDate     *string  `db:"confirmed_date"            json:"confirmed_date"`
	PostedDate        *string  `db:"posted_date"               json:"posted_date"`
	PostedUnits       *float32 `db:"posted_units"              json:"posted_units"`
	AcaKey            *uint64  `db:"aca_key"                   json:"aca_key"`
	SettledDate       *string  `db:"settled_date"              json:"settled_date"`
	RecOrder          *uint64  `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8    `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string  `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string  `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string  `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string  `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string  `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string  `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8   `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64  `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string  `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string  `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string  `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string  `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string  `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string  `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string  `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type AdminTrTransactionList struct {
	TransactionKey   uint64  `json:"transaction_key"`
	BranchName       string  `json:"branch_name"`
	AgentName        string  `json:"agent_name"`
	CustomerName     string  `json:"customer_name"`
	ProductName      string  `json:"product_name"`
	TransStatus      string  `json:"trans_status"`
	TransDate        string  `json:"trans_date"`
	TransType        string  `json:"trans_type"`
	NavDate          string  `json:"nav_date"`
	TransAmount      float32 `json:"trans_amount"`
	TransUnit        float32 `json:"trans_unit"`
	TotalAmount      float32 `json:"total_amount"`
	TransBankName    string  `json:"trans_bank_name"`
	TransBankAccNo   *string `json:"trans_bank_accno"`
	TransBankaccName *string `json:"trans_bankacc_name"`
}

type CountData struct {
	CountData int `db:"count_data"             json:"count_data"`
}

type AdminTransactionDetail struct {
	TransactionKey          uint64                   `json:"transaction_key"`
	Branch                  *BranchTrans             `json:"branch"`
	Agent                   *AgentTrans              `json:"agent"`
	Customer                CustomerTrans            `json:"customer"`
	Product                 ProductTrans             `json:"product"`
	TransStatus             TransStatus              `json:"trans_status"`
	TransDate               string                   `json:"trans_date"`
	TransType               TransType                `json:"trans_type"`
	TrxCode                 *LookupTrans             `json:"trx_code"`
	NavDate                 string                   `json:"nav_date"`
	EntryMode               *LookupTrans             `json:"entry_mode"`
	TransAmount             float32                  `json:"trans_amount"`
	TransUnit               float32                  `json:"trans_unit"`
	TransUnitPercent        *float32                 `json:"trans_unit_percent"`
	FlagRedemtAll           bool                     `json:"flag_redempt_all"`
	FlagNewSub              bool                     `json:"flag_newsub"`
	TransFeePercent         float32                  `json:"trans_fee_percent"`
	TransFeeAmount          float32                  `json:"trans_fee_amount"`
	ChargesFeeAmount        float32                  `json:"charges_fee_amount"`
	ServicesFeeAmount       float32                  `json:"services_fee_amount"`
	TotalAmount             float32                  `json:"total_amount"`
	SettlementDate          *string                  `json:"settlement_date"`
	TransBankAccNo          *string                  `json:"trans_bank_accno"`
	TransBankaccName        *string                  `json:"trans_bankacc_name"`
	TransBank               *TransBank               `json:"trans_bank"`
	TransRemarks            *string                  `json:"trans_remarks"`
	TransReferences         *string                  `json:"trans_references"`
	PromoCode               *string                  `json:"promo_code"`
	SalesCode               *string                  `json:"sales_code"`
	RiskWaiver              bool                     `json:"risk_waiver"`
	FileUploadDate          *string                  `json:"file_upload_date"`
	PaymentMethod           *LookupTrans             `json:"payment_method"`
	TrxRiskLevel            *LookupTrans             `json:"trx_risk_level"`
	ProceedDate             *string                  `json:"proceed_date"`
	ProceedAmount           *float32                 `json:"proceed_amount"`
	SentDate                *string                  `json:"sent_date"`
	SentReferences          *string                  `json:"sent_references"`
	ConfirmedDate           *string                  `json:"confirmed_date"`
	PostedDate              *string                  `json:"posted_date"`
	PostedUnits             *float32                 `json:"posted_units"`
	Aca                     *AcaTrans                `json:"aca"`
	SettledDate             *string                  `json:"settled_date"`
	RecImage1               *string                  `json:"rec_image1"`
	RecCreatedDate          *string                  `json:"rec_created_date"`
	RecCreatedBy            *string                  `json:"rec_created_by"`
	TransactionConfirmation *TransactionConfirmation `json:"transaction_confirmation"`
}

type BranchTrans struct {
	BranchKey  uint64 `json:"branch_key"`
	BranchCode string `json:"branch_code"`
	BranchName string `json:"branch_name"`
}

type AgentTrans struct {
	AgentKey  uint64 `json:"agent_key"`
	AgentCode string `json:"agent_code"`
	AgentName string `json:"agent_name"`
}

type CustomerTrans struct {
	CustomerKey    uint64  `json:"customer_key"`
	FullName       string  `json:"full_name"`
	SidNo          *string `json:"sid_no"`
	UnitHolderIDno string  `json:"unit_holder_idno"`
}

type ProductTrans struct {
	ProductKey  uint64 `json:"product_key"`
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
}

type TransStatus struct {
	TransStatusKey    uint64  `json:"trans_status_key"`
	StatusCode        *string `json:"status_code"`
	StatusDescription *string `json:"status_description"`
}

type TransType struct {
	TransTypeKey    uint64  `json:"trans_type_key"`
	TypeCode        *string `json:"type_code"`
	TypeDescription *string `json:"type_description"`
}

type TransBank struct {
	BankKey  uint64 `json:"bank_key"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

type LookupTrans struct {
	LookupKey   uint64  `json:"lookup_key"`
	LkpGroupKey uint64  `json:"lkp_group_key"`
	LkpCode     *string `json:"lkp_code"`
	LkpName     *string `json:"lkp_name"`
}
type AcaTrans struct {
	AcaKey    uint64 `json:"aca_key"`
	AccKey    uint64 `json:"acc_key"`
	AgentKey  uint64 `json:"agent_key"`
	AgentCode string `json:"agent_code"`
	AgentName string `json:"agent_name"`
}
type TransactionConfirmation struct {
	TcKey           uint64  `json:"tc_key"`
	ConfirmDate     string  `json:"confirm_date"`
	ConfirmedAmount float32 `json:"confirmed_amount"`
	ConfirmedUnit   float32 `json:"confirmed_unit"`
}

func AdminGetAllTrTransaction(c *[]TrTransaction, limit uint64, offset uint64, nolimit bool,
	params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := `SELECT
              tr_transaction.*
			  FROM tr_transaction`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_transaction."+field+" = '"+value+"'")
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
				condition += " WHERE tr_transaction." + fieldIn + " IN(" + inQuery + ")"
			}
		} else {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " AND tr_transaction." + fieldIn + " IN(" + inQuery + ")"
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

func AdminGetCountTrTransaction(c *CountData, params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := `SELECT
              count(tr_transaction.transaction_key) as count_data
			  FROM tr_transaction`

	var whereClause []string
	var condition string

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
				condition += " WHERE tr_transaction." + fieldIn + " IN(" + inQuery + ")"
			}
		} else {
			if len(valueIn) > 0 {
				inQuery := strings.Join(valueIn, ",")
				condition += " AND tr_transaction." + fieldIn + " IN(" + inQuery + ")"
			}
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

func GetAllTrTransaction(c *[]TrTransaction, params map[string]string) (int, error) {
	query := `SELECT
              tr_transaction.* FROM 
			  tr_transaction`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "tr_transaction."+field+" = '"+value+"'")
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

func UpdateTrTransaction(params map[string]string) (int, error) {
	query := "UPDATE tr_transaction SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "transaction_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE transaction_key = " + params["transaction_key"]
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	if row > 0 {
		tx.Commit()
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func CreateTrTransaction(params map[string]string) (int, error) {
	query := "INSERT INTO tr_transaction"
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

func GetTrTransaction(c *TrTransaction, key string) (int, error) {
	query := `SELECT tr_transaction.* FROM tr_transaction WHERE tr_transaction.rec_status = 1 AND tr_transaction.transaction_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func UpdateTrTransactionByKeyIn(params map[string]string, valueIn []string, fieldIn string) (int, error) {
	query := "UPDATE tr_transaction SET "
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
	query += " WHERE tr_transaction." + fieldIn + " IN(" + inQuery + ")"

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

func GetTrTransactionIn(c *[]TrTransaction, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_transaction.* FROM 
				tr_transaction `
	query := query2 + " WHERE tr_transaction.rec_status = 1 AND tr_transaction." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
