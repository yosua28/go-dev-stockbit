package models

import (
	"api/db"
	"log"
	"net/http"
)

type MsCustomerBankAccount struct {
	CustBankaccKey    uint64  `db:"cust_bankacc_key"      json:"cust_bankacc_key"`
	CustomerKey       uint64  `db:"customer_key"          json:"customer_key"`
	AankAccountKey    uint64  `db:"bank_account_key"      json:"bank_account_key"`
	FlagPriority      uint8   `db:"flag_priority"         json:"flag_priority"`
	BankAccountName   string  `db:"bank_account_name"     json:"bank_account_name"`
	RecOrder          *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type MsCustomerBankAccountInfo struct {
	CustBankaccKey uint64  `db:"cust_bankacc_key"      json:"cust_bankacc_key"`
	BankName       string  `db:"bank_name"             json:"bank_name"`
	AccountNo      string  `db:"account_no"            json:"account_no"`
	AccountName    string  `db:"account_name"          json:"account_name"`
	BranchName     *string `db:"branch_name"           json:"branch_name"`
}

func CreateMsCustomerBankAccount(params map[string]string) (int, error) {
	query := "INSERT INTO ms_customer_bank_account"
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
	log.Println(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetAllMsCustomerBankAccountTransaction(c *[]MsCustomerBankAccountInfo, customerKey string) (int, error) {
	query2 := `SELECT 
				ba.cust_bankacc_key AS cust_bankacc_key, 
				bank.bank_fullname AS bank_name, 
				b.account_no AS account_no, 
				b.account_holder_name AS account_name,
				b.branch_name as branch_name 
			FROM ms_customer_bank_account AS ba 
			INNER JOIN ms_bank_account AS b ON b.bank_account_key = ba.bank_account_key
			INNER JOIN ms_bank AS bank ON bank.bank_key = b.bank_key`
	query := query2 + " WHERE ba.rec_status = 1 AND ba.customer_key = '" + customerKey + "'"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsCustomerBankAccountTransactionByKey(c *MsCustomerBankAccountInfo, custBankKey string) (int, error) {
	query2 := `SELECT 
				ba.cust_bankacc_key AS cust_bankacc_key, 
				bank.bank_fullname AS bank_name, 
				b.account_no AS account_no, 
				b.account_holder_name AS account_name,
				b.branch_name as branch_name  
			FROM ms_customer_bank_account AS ba 
			INNER JOIN ms_bank_account AS b ON b.bank_account_key = ba.bank_account_key
			INNER JOIN ms_bank AS bank ON bank.bank_key = b.bank_key`
	query := query2 + " WHERE ba.rec_status = 1 AND ba.cust_bankacc_key = '" + custBankKey + "'"

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
