package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type Portofolio struct {
	Date          string
	Cif           string
	Sid           string
	Name          string
	Address       string
	City          string
	Country       string
	Datas         []ProductPortofolio
	Total         string
	TotalGainLoss string
}

type ProductPortofolio struct {
	ProductName string
	AvgNav      string
	Nav         string
	Unit        string
	CCY         string
	Amount      string
	GainLoss    string
	Kurs        string
	AmountIDR   string
	GainLossIDR string
}

type TrTransaction struct {
	TransactionKey    uint64           `db:"transaction_key"           json:"transaction_key"`
	ParentKey         *uint64          `db:"parent_key"                json:"parent_key"`
	IDTransaction     *uint64          `db:"id_transaction"            json:"id_transaction"`
	BranchKey         *uint64          `db:"branch_key"                json:"branch_key"`
	AgentKey          *uint64          `db:"agent_key"                 json:"agent_key"`
	CustomerKey       uint64           `db:"customer_key"              json:"customer_key"`
	ProductKey        uint64           `db:"product_key"               json:"product_key"`
	TransStatusKey    uint64           `db:"trans_status_key"          json:"trans_status_key"`
	TransDate         string           `db:"trans_date"                json:"trans_date"`
	TransTypeKey      uint64           `db:"trans_type_key"            json:"trans_type_key"`
	TrxCode           *uint64          `db:"trx_code"                  json:"trx_code"`
	NavDate           string           `db:"nav_date"                  json:"nav_date"`
	EntryMode         *uint64          `db:"entry_mode"                json:"entry_mode"`
	TransCalcMethod   *uint64          `db:"trans_calc_method"         json:"trans_calc_method"`
	TransAmount       decimal.Decimal  `db:"trans_amount"              json:"trans_amount"`
	TransUnit         decimal.Decimal  `db:"trans_unit"                json:"trans_unit"`
	TransUnitPercent  *decimal.Decimal `db:"trans_unit_percent"        json:"trans_unit_percent"`
	FlagRedemtAll     *uint8           `db:"flag_redempt_all"          json:"flag_redempt_all"`
	FlagNewSub        *uint8           `db:"flag_newsub"               json:"flag_newsub"`
	TransFeePercent   decimal.Decimal  `db:"trans_fee_percent"         json:"trans_fee_percent"`
	TransFeeAmount    decimal.Decimal  `db:"trans_fee_amount"          json:"trans_fee_amount"`
	ChargesFeeAmount  decimal.Decimal  `db:"charges_fee_amount"        json:"charges_fee_amount"`
	ServicesFeeAmount decimal.Decimal  `db:"services_fee_amount"       json:"services_fee_amount"`
	TotalAmount       decimal.Decimal  `db:"total_amount"              json:"total_amount"`
	SettlementDate    *string          `db:"settlement_date"           json:"settlement_date"`
	TransBankAccNo    *string          `db:"trans_bank_accno"          json:"trans_bank_accno"`
	TransBankaccName  *string          `db:"trans_bankacc_name"        json:"trans_bankacc_name"`
	TransBankKey      *uint64          `db:"trans_bank_key"            json:"trans_bank_key"`
	TransRemarks      *string          `db:"trans_remarks"             json:"trans_remarks"`
	TransReferences   *string          `db:"trans_references"          json:"trans_references"`
	PromoCode         *string          `db:"promo_code"                json:"promo_code"`
	SalesCode         *string          `db:"sales_code"                json:"sales_code"`
	RiskWaiver        uint8            `db:"risk_waiver"               json:"risk_waiver"`
	AddtoAutoInvest   *uint8           `db:"addto_auto_invest"         json:"addto_auto_invest"`
	TransSource       *uint64          `db:"trans_source"              json:"trans_source"`
	FileUploadDate    *string          `db:"file_upload_date"          json:"file_upload_date"`
	PaymentMethod     *uint64          `db:"payment_method"            json:"payment_method"`
	Check1Date        *string          `db:"check1_date"               json:"check1_date"`
	Check1Flag        *uint8           `db:"check1_flag"               json:"check1_flag"`
	Check1References  *string          `db:"check1_references"         json:"check1_references"`
	Check1Notes       *string          `db:"check1_notes"              json:"check1_notes"`
	Check2Date        *string          `db:"check2_date"               json:"check2_date"`
	Check2Flag        *uint8           `db:"check2_flag"               json:"check2_flag"`
	Check2References  *string          `db:"check2_references"         json:"check2_references"`
	Check2Notes       *string          `db:"check2_notes"              json:"check2_notes"`
	TrxRiskLevel      *uint64          `db:"trx_risk_level"            json:"trx_risk_level"`
	ProceedDate       *string          `db:"proceed_date"              json:"proceed_date"`
	ProceedAmount     *decimal.Decimal `db:"proceed_amount"            json:"proceed_amount"`
	SentDate          *string          `db:"sent_date"                 json:"sent_date"`
	SentReferences    *string          `db:"sent_references"           json:"sent_references"`
	ConfirmedDate     *string          `db:"confirmed_date"            json:"confirmed_date"`
	PostedDate        *string          `db:"posted_date"               json:"posted_date"`
	PostedUnits       *decimal.Decimal `db:"posted_units"              json:"posted_units"`
	AcaKey            *uint64          `db:"aca_key"                   json:"aca_key"`
	SettledDate       *string          `db:"settled_date"              json:"settled_date"`
	BatchKey          *uint64          `db:"batch_key"                 json:"batch_key"`
	RecOrder          *uint64          `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8            `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string          `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string          `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string          `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string          `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string          `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string          `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8           `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64          `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string          `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string          `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string          `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string          `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string          `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string          `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string          `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type TrTransactionList struct {
	TransactionKey uint64          `json:"transaction_key"`
	ProductName    string          `json:"product_name"`
	TransStatus    string          `json:"trans_status"`
	TransDate      string          `json:"trans_date"`
	TransType      string          `json:"trans_type"`
	NavDate        string          `json:"nav_date"`
	NavValue       decimal.Decimal `json:"nav_value"`
	TransAmount    decimal.Decimal `json:"trans_amount,omitempty"`
	TransUnit      decimal.Decimal `json:"trans_unit,omitempty"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	Uploaded       bool            `json:"uploaded"`
	DateUploaded   *string         `json:"date_uploaded"`
	BankName       *string         `json:"bank_name"`
	BankAccNo      *string         `json:"bank_accno"`
	BankAccName    *string         `json:"bankacc_name"`
	ProductOut     *string         `json:"product_name_out"`
	ProductIn      *string         `json:"product_name_in"`
}

type AdminTrTransactionList struct {
	TransactionKey   uint64          `json:"transaction_key"`
	BranchName       string          `json:"branch_name"`
	AgentName        string          `json:"agent_name"`
	CustomerName     string          `json:"customer_name"`
	ProductName      string          `json:"product_name"`
	TransStatus      string          `json:"trans_status"`
	TransDate        string          `json:"trans_date"`
	TransType        string          `json:"trans_type"`
	NavDate          string          `json:"nav_date"`
	TransAmount      decimal.Decimal `json:"trans_amount"`
	TransUnit        decimal.Decimal `json:"trans_unit"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
	TransBankName    string          `json:"trans_bank_name"`
	TransBankAccNo   *string         `json:"trans_bank_accno"`
	TransBankaccName *string         `json:"trans_bankacc_name"`
	ProductOut       *string         `json:"product_name_out"`
	ProductIn        *string         `json:"product_name_in"`
}

type AdminTrTransactionInquiryList struct {
	TransactionKey uint64          `json:"transaction_key"`
	CustomerKey    uint64          `json:"customer_key,omitempty"`
	ProductKey     uint64          `json:"product_key,omitempty"`
	BranchName     string          `json:"branch_name"`
	AgentName      string          `json:"agent_name"`
	CustomerName   string          `json:"customer_name"`
	ProductName    string          `json:"product_name"`
	TransStatus    string          `json:"trans_status"`
	TransDate      string          `json:"trans_date"`
	TransType      string          `json:"trans_type"`
	NavDate        string          `json:"nav_date"`
	TransAmount    decimal.Decimal `json:"trans_amount"`
	TransUnit      decimal.Decimal `json:"trans_unit"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
}

type CountData struct {
	CountData int `db:"count_data"             json:"count_data"`
}

type NavValue struct {
	NavValue *decimal.Decimal `db:"nav_value"              json:"nav_value"`
}

type AdminTransactionDetail struct {
	TransactionKey              uint64                               `json:"transaction_key"`
	Branch                      *BranchTrans                         `json:"branch"`
	Agent                       *AgentTrans                          `json:"agent"`
	Customer                    CustomerTrans                        `json:"customer"`
	Product                     ProductTrans                         `json:"product"`
	TransStatus                 TransStatus                          `json:"trans_status"`
	TransDate                   string                               `json:"trans_date"`
	TransType                   TransType                            `json:"trans_type"`
	TrxCode                     *LookupTrans                         `json:"trx_code"`
	NavDate                     string                               `json:"nav_date"`
	EntryMode                   *LookupTrans                         `json:"entry_mode"`
	TransAmount                 decimal.Decimal                      `json:"trans_amount"`
	TransUnit                   decimal.Decimal                      `json:"trans_unit"`
	TransUnitPercent            *decimal.Decimal                     `json:"trans_unit_percent"`
	FlagRedemtAll               bool                                 `json:"flag_redempt_all"`
	FlagNewSub                  bool                                 `json:"flag_newsub"`
	TransFeePercent             decimal.Decimal                      `json:"trans_fee_percent"`
	TransFeeAmount              decimal.Decimal                      `json:"trans_fee_amount"`
	ChargesFeeAmount            decimal.Decimal                      `json:"charges_fee_amount"`
	ServicesFeeAmount           decimal.Decimal                      `json:"services_fee_amount"`
	TotalAmount                 decimal.Decimal                      `json:"total_amount"`
	SettlementDate              *string                              `json:"settlement_date"`
	TransBankAccNo              *string                              `json:"trans_bank_accno"`
	TransBankaccName            *string                              `json:"trans_bankacc_name"`
	TransBank                   *TransBank                           `json:"trans_bank"`
	TransRemarks                *string                              `json:"trans_remarks"`
	TransReferences             *string                              `json:"trans_references"`
	PromoCode                   *string                              `json:"promo_code"`
	SalesCode                   *string                              `json:"sales_code"`
	RiskWaiver                  bool                                 `json:"risk_waiver"`
	FileUploadDate              *string                              `json:"file_upload_date"`
	PaymentMethod               *LookupTrans                         `json:"payment_method"`
	TransactionSettlement       *[]TransactionSettlement             `json:"transaction_settlement"`
	TrxRiskLevel                *LookupTrans                         `json:"trx_risk_level"`
	ProceedDate                 *string                              `json:"proceed_date"`
	ProceedAmount               *decimal.Decimal                     `json:"proceed_amount"`
	SentDate                    *string                              `json:"sent_date"`
	SentReferences              *string                              `json:"sent_references"`
	ConfirmedDate               *string                              `json:"confirmed_date"`
	PostedDate                  *string                              `json:"posted_date"`
	PostedUnits                 *decimal.Decimal                     `json:"posted_units"`
	Aca                         *AcaTrans                            `json:"aca"`
	SettledDate                 *string                              `json:"settled_date"`
	RecImage1                   *string                              `json:"rec_image1"`
	RecCreatedDate              *string                              `json:"rec_created_date"`
	RecCreatedBy                *string                              `json:"rec_created_by"`
	TransactionConfirmation     *TransactionConfirmation             `json:"transaction_confirmation"`
	ProductBankAccount          *MsProductBankAccountTransactionInfo `json:"product_bank_account"`
	CustomerBankAccount         *MsCustomerBankAccountInfo           `json:"customer_bank_account"`
	IsEnableUnposting           bool                                 `json:"is_enable_unposting"`
	MessageEnableUnposting      string                               `json:"message_enable_unposting"`
	TransactionConfirmationInfo *TrTransactionConfirmationInfo       `json:"transaction_confirmation_info"`
	Promo                       *TrPromoData                         `json:"promo"`
}

type DownloadFormatExcelList struct {
	IDTransaction   uint64           `json:"id_transaction"`
	IDCategory      string           `json:"id_category"`
	ProductName     string           `json:"product_name"`
	FullName        string           `json:"full_name"`
	NavDate         string           `json:"nav_date"`
	TransactionDate string           `json:"transaction_date"`
	Units           decimal.Decimal  `json:"units"`
	NetAmount       decimal.Decimal  `json:"net_amount"`
	NavValue        *decimal.Decimal `json:"nav_value"`
	ApproveUnits    decimal.Decimal  `json:"approve_units"`
	ApproveAmount   decimal.Decimal  `json:"approve_amount"`
	Keterangan      string           `json:"keterangan"`
	Result          string           `json:"result"`
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
	TcKey           uint64          `json:"tc_key"`
	ConfirmDate     string          `json:"confirm_date"`
	ConfirmedAmount decimal.Decimal `json:"confirmed_amount"`
	ConfirmedUnit   decimal.Decimal `json:"confirmed_unit"`
}

type ParamBatchTrTransaction struct {
	ProductCode    string  `db:"product_code"     json:"product_code"`
	TypeCode       string  `db:"type_code"        json:"type_code"`
	Bulan          string  `db:"bulan"            json:"bulan"`
	Tahun          string  `db:"tahun"            json:"tahun"`
	NavDate        string  `db:"nav_date"         json:"nav_date"`
	ProductKey     uint64  `db:"product_key"      json:"product_key"`
	TransTypeKey   uint64  `db:"trans_type_key"   json:"trans_type_key"`
	TransactionKey string  `db:"transaction_key"  json:"transaction_key"`
	TransDate      string  `db:"trans_date"       json:"trans_date"`
	Batch          *uint64 `db:"batch"            json:"batch"`
}

type ProductCheckAllowRedmSwtching struct {
	ProductKey uint64 `db:"product_key"             json:"product_key"`
}

type TransactionCustomerHistory struct {
	ProductName string  `db:"product_name"      json:"product_name"`
	FullName    string  `db:"full_name"         json:"full_name"`
	Cif         *string `db:"cif"               json:"cif"`
	Sid         *string `db:"sid"               json:"sid"`
	ParamDetail *string `db:"param_detail"      json:"param_detail"`
}

type TransactionConsumenProduct struct {
	TransactionKey  uint64          `db:"transaction_key"     json:"transaction_key"`
	TransTypeKey    uint64          `db:"trans_type_key"      json:"trans_type_key"`
	NavDate         string          `db:"nav_date"            json:"nav_date"`
	TypeDescription string          `db:"type_description"    json:"type_description"`
	NavValue        decimal.Decimal `db:"nav_value"            json:"nav_value"`
	Unit            decimal.Decimal `db:"unit"                json:"unit"`
	GrossAmount     decimal.Decimal `db:"gross_amount"        json:"gross_amount"`
	FeeAmount       decimal.Decimal `db:"fee_amount"          json:"fee_amount"`
	NetAmount       decimal.Decimal `db:"net_amount"          json:"net_amount"`
}

type DataDetailTransaksiCustomerProduct struct {
	DataTransaksi     DetailHeaderTransaksiCustomer `json:"data_transaksi"`
	DetailTransaksi   *[]TransactionConsumenProduct `json:"detail_transaksi"`
	CountSubscription decimal.Decimal               `json:"total_subscription"`
	CountRedemption   decimal.Decimal               `json:"total_redemption"`
	CountNetSub       decimal.Decimal               `json:"total_netsub"`
}

type DetailHeaderTransaksiCustomer struct {
	UnitHolder  string `db:"unit_holder"          json:"unit_holder"`
	FullName    string `db:"full_name"            json:"full_name"`
	Sid         string `db:"sid"                  json:"sid"`
	IfuaNo      string `db:"ifua_no"              json:"ifua_no"`
	NavDateFrom string `db:"nav_date_from"        json:"nav_date_from"`
	NavDateTo   string `db:"nav_date_to"          json:"nav_date_to"`
	ProductName string `db:"product_name"         json:"product_name"`
}

func AdminGetAllTrTransaction(c *[]TrTransaction, limit uint64, offset uint64, nolimit bool,
	params map[string]string, valueIn []string, fieldIn string, isAll bool) (int, error) {
	query := `SELECT
              t.*
			  FROM tr_transaction as t
			  inner join ms_customer as c on c.customer_key = t.customer_key
			  WHERE t.rec_status = 1 AND t.trans_status_key != 3`

	if isAll == false {
		query += " AND t.trans_type_key != 3"
	}
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			if field == "c.openacc_branch_key" {
				whereClause = append(whereClause, field+" = '"+value+"'")
			} else {
				whereClause = append(whereClause, "t."+field+" = '"+value+"'")
			}
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

	if len(valueIn) > 0 {
		inQuery := strings.Join(valueIn, ",")
		condition += " AND t." + fieldIn + " IN(" + inQuery + ")"
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
              count(t.transaction_key) as count_data
			  FROM tr_transaction as t
			  inner join ms_customer as c on c.customer_key = t.customer_key
			  WHERE t.trans_type_key != 3 `

	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			if field == "c.openacc_branch_key" {
				whereClause = append(whereClause, field+" = '"+value+"'")
			} else {
				whereClause = append(whereClause, "t."+field+" = '"+value+"'")
			}
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

	if len(valueIn) > 0 {
		inQuery := strings.Join(valueIn, ",")
		condition += " AND t." + fieldIn + " IN(" + inQuery + ")"
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

func CreateTrTransaction(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_transaction"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		if value == "NULL" {
			var s *string
			bindvars = append(bindvars, s)
		} else {
			bindvars = append(bindvars, value)
		}

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

func GetTrTransaction(c *TrTransaction, key string) (int, error) {
	query := `SELECT tr_transaction.* FROM tr_transaction WHERE tr_transaction.transaction_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionByField(c *TrTransaction, field string, value string) (int, error) {
	query := `SELECT tr_transaction.* FROM tr_transaction WHERE tr_transaction.rec_status = 1 AND tr_transaction.` + field + ` = '` + value + `'`
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

func GetAllTransactionByParamAndValueIn(c *[]TrTransaction, limit uint64, offset uint64,
	nolimit bool, params map[string]string, valueIn []string, fieldIn string) (int, error) {
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

func GetTrTransactionDateRange(c *[]TrTransaction, params map[string]string, start string, end string) (int, error) {
	query := `SELECT
              tr_transaction.* FROM 
			  tr_transaction`
	query += " WHERE tr_transaction.trans_date >= " + start + " AND tr_transaction.trans_date <= " + end
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

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrTransactionOnProcess(c *[]TrTransaction, params map[string]string) (int, error) {
	query := `SELECT
              tr_transaction.* FROM 
			  tr_transaction WHERE tr_transaction.trans_status_key < 9 AND DATE_ADD(tr_transaction.nav_date, INTERVAL 1 DAY) >= NOW()`
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

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func ParamBatchTrTransactionByKey(c *ParamBatchTrTransaction, transactionKey string) (int, error) {
	query := `SELECT
				p.product_code AS product_code,
				tt.type_code AS type_code,
				MONTH(t.nav_date) AS bulan,
				YEAR(t.nav_date) AS tahun,
				t.nav_date AS nav_date,
				p.product_key AS product_key,
				t.trans_type_key AS trans_type_key,
				t.transaction_key AS transaction_key,
				t.trans_date AS trans_date,
			    (SELECT batch_number AS bat FROM tr_transaction_batch ORDER BY batch_number DESC LIMIT 1) AS batch 
			FROM tr_transaction AS t
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			INNER JOIN tr_transaction_type AS tt ON tt.trans_type_key = t.trans_type_key
			WHERE t.transaction_key = ` + transactionKey + ` LIMIT 1`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CheckTrTransactionLastProductCustomer(c *TrTransaction, customerKey string, productKey string, transKey string) (int, error) {
	query := `SELECT
				tr_transaction.* FROM 
				tr_transaction `
	query += " WHERE tr_transaction.rec_status = 1"
	query += " AND trans_status_key = 9"
	query += " AND customer_key = " + customerKey
	query += " AND product_key = " + productKey
	query += " AND transaction_key > " + transKey + " LIMIT 1"

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CheckProductAllowRedmOrSwitching(c *[]ProductCheckAllowRedmSwtching, customerKey string, productKeyIn []string) (int, error) {

	inQuery := strings.Join(productKeyIn, ",")

	query := `SELECT product_key 
				FROM tr_transaction`
	query += " WHERE rec_status = 1"
	query += " AND trans_type_key IN (2,3)"
	query += " AND trans_status_key NOT IN (3,9)"
	query += " AND customer_key = " + customerKey
	query += " AND product_key IN(" + inQuery + ")"
	query += " GROUP BY product_key"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetTransactionCustomerHistory(c *[]TransactionCustomerHistory, limit uint64, offset uint64, params map[string]string, paramsLike map[string]string, nolimit bool) (int, error) {
	dateFrom := ""
	dateTo := ""

	var whereClause []string
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "dateFrom" || field == "dateTo") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
		if field == "dateFrom" {
			dateFrom = value
		}
		if field == "dateTo" {
			dateTo = value
		}
	}

	query := `SELECT 
				p.product_name_alt AS product_name,
				c.full_name AS full_name,
				c.unit_holder_idno AS cif,
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				TO_BASE64(CONCAT(c.customer_key, ",", p.product_key, ",", "` + dateFrom + `", ",", "` + dateTo + `")) AS param_detail 
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key 
			INNER JOIN ms_product AS p ON t.product_key = p.product_key
			WHERE t.trans_status_key = 9 AND t.rec_status = 1 AND tc.rec_status = 1`

	var present bool
	var condition string
	var conditionOrder string

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
	}

	if (dateFrom != "") && (dateTo != "") {
		query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"
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

	query += condition

	query += " GROUP BY t.customer_key, t.product_key"

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		conditionOrder += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			conditionOrder += " " + orderType
		}
	}
	query += conditionOrder

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

func AdminCountTransactionCustomerHistory(c *CountData, params map[string]string, paramsLike map[string]string) (int, error) {
	query := `SELECT 
				count(t.transaction_key) AS count_data 
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key
			INNER JOIN ms_customer AS c ON c.customer_key = t.customer_key 
			INNER JOIN ms_product AS p ON t.product_key = p.product_key
			WHERE t.trans_status_key = 9 AND t.rec_status = 1 AND tc.rec_status = 1`

	var whereClause []string
	var condition string
	dateFrom := ""
	dateTo := ""

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "dateFrom" || field == "dateTo") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
		if field == "dateFrom" {
			dateFrom = value
		}
		if field == "dateTo" {
			dateTo = value
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
	}

	if (dateFrom != "") && (dateTo != "") {
		query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"
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

	query += condition

	query += " GROUP BY t.customer_key, t.product_key"

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetTransactionConsumenProduct(c *[]TransactionConsumenProduct, params map[string]string, paramsLike map[string]string, dateFrom string, dateTo string) (int, error) {
	query := `SELECT 
				t.transaction_key as transaction_key,
				t.trans_type_key as trans_type_key,
				DATE_FORMAT(t.nav_date, '%d %M %Y') AS nav_date, 
				ty.type_description as type_description,
				nav.nav_value as nav_value,
				tc.confirmed_unit as unit,
				(CASE
					WHEN t.total_amount IS NULL OR t.total_amount = 0 THEN tc.confirmed_amount
					ELSE t.total_amount
				END) AS gross_amount,
				(t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) AS fee_amount, 
				(CASE
					WHEN t.total_amount IS NULL OR t.total_amount = 0 THEN tc.confirmed_amount - (t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount)
					ELSE (t.total_amount - (t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount))
				END) AS net_amount 
			FROM tr_transaction AS t 
			INNER JOIN tr_transaction_confirmation AS tc ON t.transaction_key = tc.transaction_key
			INNER JOIN tr_nav AS nav ON t.nav_date = nav.nav_date AND t.product_key = nav.product_key 
			INNER JOIN tr_transaction_type AS ty ON ty.trans_type_key = t.trans_type_key
			WHERE t.rec_status = 1 AND tc.rec_status = 1 AND t.trans_status_key = 9 `

	query += " AND (t.nav_date BETWEEN '" + dateFrom + "' AND '" + dateTo + "')"

	var present bool
	var condition string

	var whereClause []string
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

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetDetailHeaderTransaksiCustomer(c *DetailHeaderTransaksiCustomer, dateFrom string, dateTo string, params map[string]string) (int, error) {
	query := `SELECT
				c.unit_holder_idno AS unit_holder,
				c.full_name AS full_name,
				(CASE
					WHEN c.sid_no IS NULL THEN ""
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN a.ifua_no IS NULL THEN ""
					ELSE a.ifua_no
				END) AS ifua_no, 
				DATE_FORMAT("` + dateFrom + `", '%d %M %Y') AS nav_date_from,
				DATE_FORMAT("` + dateTo + `", '%d %M %Y') AS nav_date_to, 
				CONCAT(p.product_code, " - ", p.product_name_alt) AS product_name  
			FROM tr_account AS a 
			INNER JOIN ms_customer AS c ON c.customer_key = a.customer_key
			INNER JOIN ms_product AS p ON p.product_key = a.product_key
			INNER JOIN sc_user_login AS l ON l.customer_key = c.customer_key 
			INNER JOIN sc_user_dept AS d ON d.user_dept_key = l.user_dept_key 
			WHERE a.rec_status = 1 AND c.rec_status = 1 `

	var present bool
	var condition string

	var whereClause []string
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
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetAllTrTransactionPosting(c *[]TrTransaction, params map[string]string, valueIn []string, fieldIn string, isAll bool) (int, error) {
	query := `SELECT
              t.*
			  FROM tr_transaction as t
			  WHERE t.rec_status = 1 AND t.trans_status_key != 3`

	if isAll == false {
		query += " AND t.trans_type_key != 3"
	}
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "t."+field+" = '"+value+"'")
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

	if len(valueIn) > 0 {
		inQuery := strings.Join(valueIn, ",")
		condition += " AND t." + fieldIn + " IN(" + inQuery + ")"
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
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminLastAvgNav(c *NavValue, productKey string, customerKey string, date string) (int, error) {
	query := `SELECT 
				tc.avg_nav AS nav_value
			FROM tr_transaction AS t
			INNER JOIN tr_transaction_confirmation AS tc ON tc.transaction_key = t.transaction_key AND tc.rec_status = '1'
			WHERE t.product_key = '` + productKey + `' 
			AND t.customer_key = '` + customerKey + `' 
			AND t.nav_date <= '` + date + `' 
			AND t.trans_status_key = '9'
			ORDER BY tc.tc_key DESC LIMIT 1`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type DetailTransactionDataSentEmail struct {
	TransactionKey        uint64          `db:"transaction_key"           json:"transaction_key"`
	TransTypeKey          uint64          `db:"trans_type_key"            json:"trans_type_key"`
	FullName              string          `db:"full_name"                 json:"full_name"`
	Cif                   *string         `db:"cif"                       json:"cif"`
	TransDate             string          `db:"trans_date"                json:"trans_date"`
	TransTime             string          `db:"trans_time"                json:"trans_time"`
	ProductName           string          `db:"product_name"              json:"product_name"`
	CurrencySymbol        string          `db:"currency_symbol"           json:"currency_symbol"`
	EntryMode             *uint64         `db:"entry_mode"                json:"entry_mode"`
	TransAmount           decimal.Decimal `db:"trans_amount"              json:"trans_amount"`
	TransUnit             decimal.Decimal `db:"trans_unit"                json:"trans_unit"`
	Fee                   decimal.Decimal `db:"fee"                       json:"fee"`
	PaymentMethod         *uint64         `db:"payment_method"            json:"payment_method"`
	PaymentMethodName     *string         `db:"payment_method_name"       json:"payment_method_name"`
	RekBankCustodian      *string         `db:"rek_bank_custodian"        json:"rek_bank_custodian"`
	NoRekBankCustomer     *string         `db:"no_rek_bank_customer"      json:"no_rek_bank_customer"`
	NameRekBankCustomer   *string         `db:"name_rek_bank_customer"    json:"name_rek_bank_customer"`
	CabangRekBankCustomer *string         `db:"cabang_rek_bank_customer"  json:"cabang_rek_bank_customer"`
	BankRekBankCustomer   *string         `db:"bank_rek_bank_customer"    json:"bank_rek_bank_customer"`
	Sales                 *string         `db:"sales"                     json:"sales"`
	SalesEmail            *string         `db:"sales_email"               json:"sales_email"`
	BuktiTransafer        *string         `db:"bukti_transafer"           json:"bukti_transafer"`
	ProductTujuan         *string         `db:"product_tujuan"            json:"product_tujuan"`
	UserLoginKey          string          `db:"user_login_key"            json:"user_login_key"`
}

func AdminDetailTransactionDataSentEmail(c *DetailTransactionDataSentEmail, tansactionKey string) (int, error) {
	query := `SELECT 
				t.transaction_key,
				t.trans_type_key, 
				c.full_name AS full_name,
				c.unit_holder_idno AS cif,
				DATE_FORMAT(t.trans_date, '%d %M %Y') AS trans_date,
				CONCAT(DATE_FORMAT(t.trans_date, '%H:%i'), " WIB") AS trans_time,
				p.product_name_alt AS product_name,
				cu.symbol as currency_symbol,
				t.entry_mode,
				t.trans_amount,
				t.trans_unit,
				(t.trans_fee_amount + t.charges_fee_amount + t.services_fee_amount) AS fee,
				t.payment_method,
				mp.lkp_name AS payment_method_name,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE CONCAT(ba.account_no, " - ", ba.account_holder_name)
				END) AS rek_bank_custodian,
				t.rec_image1 AS bukti_transafer,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE ba_c.account_no
				END) AS no_rek_bank_customer,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE ba_c.account_holder_name
				END) AS name_rek_bank_customer,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE ba_c.branch_name
				END) AS cabang_rek_bank_customer,
				(CASE 
					WHEN tbank.trans_bankacc_key IS NULL THEN '-' 
					ELSE b.bank_name
				END) AS bank_rek_bank_customer,
				CONCAT(a.agent_code, " - ", a.agent_name) AS sales,
				a.agent_email AS sales_email,
				t.rec_image1 AS bukti_transafer,
				p_t.product_name_alt AS product_tujuan,
				ul.user_login_key  
			FROM tr_transaction AS t
			INNER JOIN ms_customer AS c ON t.customer_key = c.customer_key
			LEFT JOIN ms_agent AS a ON a.agent_key = c.openacc_agent_key
			INNER JOIN ms_product AS p ON p.product_key = t.product_key
			LEFT join ms_currency as cu on cu.currency_key = p.currency_key 
			LEFT JOIN ms_currency AS cur ON cur.currency_key = p.currency_key
			LEFT JOIN gen_lookup AS mp ON mp.lookup_key = t.payment_method 
			LEFT JOIN tr_transaction_bank_account AS tbank ON tbank.transaction_key = t.transaction_key
			LEFT JOIN ms_product_bank_account AS pbank ON pbank.prod_bankacc_key = tbank.prod_bankacc_key
			LEFT JOIN ms_customer_bank_account AS cbank ON cbank.cust_bankacc_key = tbank.cust_bankacc_key
			LEFT JOIN ms_bank_account AS ba ON ba.bank_account_key = pbank.bank_account_key 
			LEFT JOIN ms_bank_account AS ba_c ON ba_c.bank_account_key = cbank.bank_account_key 
			LEFT JOIN ms_bank AS b ON b.bank_key = ba_c.bank_key 
			LEFT JOIN tr_transaction AS t_ch ON t_ch.parent_key = t.transaction_key
			LEFT JOIN ms_product AS p_t ON p_t.product_key = t_ch.product_key 
			INNER JOIN sc_user_login as ul on ul.customer_key = t.customer_key 
			WHERE t.rec_status = 1 AND t.transaction_key = '` + tansactionKey + `'`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
