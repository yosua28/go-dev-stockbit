package models

import (
	"api/db"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MsProductBankAccountInfo struct {
	BankAccountName    string      `json:"bank_account_name"`
	BankAccountPurpose uint64      `json:"bank_account_purpose"`
	BankAccount        BankAccount `json:"bank_account"`
}

type MsProductBankAccount struct {
	ProdBankaccKey     uint64  `db:"prod_bankacc_key"      json:"prod_bankacc_key"`
	ProductKey         *uint64 `db:"product_key"           json:"product_key"`
	BankAccountKey     *uint64 `db:"bank_account_key"      json:"bank_account_key"`
	BankAccountName    string  `db:"bank_account_name"     json:"bank_account_name"`
	BankAccountPurpose uint64  `db:"bank_account_purpose"  json:"bank_account_purpose"`
	RecOrder           *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type AdminMsProductBankAccountList struct {
	ProdBankaccKey     uint64  `db:"prod_bankacc_key"      json:"prod_bankacc_key"`
	ProductKey         *uint64 `db:"product_key"           json:"product_key"`
	ProductCode        string  `db:"product_code"          json:"product_code"`
	ProductNameAlt     string  `db:"product_name_alt"      json:"product_name_alt"`
	BankAccountName    string  `db:"bank_account_name"     json:"bank_account_name"`
	BankAccountPurpose *string `db:"bank_account_purpose"  json:"bank_account_purpose"`
	BankFullname       *string `db:"bank_fullname"         json:"bank_fullname"`
	AccountNo          string  `db:"account_no"            json:"account_no"`
	AccountHolderName  string  `db:"account_holder_name"   json:"account_holder_name"`
}

func GetAllMsProductBankAccount(c *[]MsProductBankAccount, params map[string]string) (int, error) {
	query := `SELECT
              ms_product_bank_account.* FROM 
			  ms_product_bank_account WHERE  
			  ms_product_bank_account.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product_bank_account."+field+" = '"+value+"'")
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

func AdminGetAllMsProductBankAccount(c *[]AdminMsProductBankAccountList, limit uint64, offset uint64, params map[string]string, nolimit bool, searchLike *string) (int, error) {
	query := `SELECT 
				pba.prod_bankacc_key AS prod_bankacc_key,
				p.product_key AS product_key,
				p.product_code AS product_code,
				p.product_name_alt AS product_name_alt,
				pba.bank_account_name AS bank_account_name,
				bank_account_purpose.lkp_name AS bank_account_purpose,
				bank.bank_fullname AS bank_fullname,
				ba.account_no AS account_no,
				ba.account_holder_name AS account_holder_name
			FROM ms_product_bank_account pba
			INNER JOIN ms_product p ON pba.product_key = p.product_key
			INNER JOIN ms_bank_account ba ON pba.bank_account_key = ba.bank_account_key
			LEFT JOIN gen_lookup AS bank_account_purpose ON bank_account_purpose.lookup_key = pba.bank_account_purpose
			LEFT JOIN ms_bank AS bank ON bank.bank_key = ba.bank_key
			WHERE pba.rec_status = 1`
	var present bool
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

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " pba.prod_bankacc_key LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_name_alt LIKE '%" + *searchLike + "%' OR"
		condition += " pba.bank_account_name LIKE '%" + *searchLike + "%' OR"
		condition += " bank_account_purpose.lkp_name LIKE '%" + *searchLike + "%' OR"
		condition += " bank.bank_fullname LIKE '%" + *searchLike + "%' OR"
		condition += " ba.account_no LIKE '%" + *searchLike + "%' OR"
		condition += " ba.account_holder_name LIKE '%" + *searchLike + "%')"
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
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminCountDataGetAllMsProductBankAccount(c *CountData, params map[string]string, searchLike *string) (int, error) {
	query := `SELECT 
				count(pba.prod_bankacc_key) AS count_data 
			FROM ms_product_bank_account pba 
			INNER JOIN ms_product p ON pba.product_key = p.product_key 
			INNER JOIN ms_bank_account ba ON pba.bank_account_key = ba.bank_account_key 
			LEFT JOIN gen_lookup AS bank_account_purpose ON bank_account_purpose.lookup_key = pba.bank_account_purpose 
			LEFT JOIN ms_bank AS bank ON bank.bank_key = ba.bank_key 
			WHERE pba.rec_status = 1`
	var present bool
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

	//search like all
	if searchLike != nil {
		condition += " AND ("
		condition += " pba.prod_bankacc_key LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_name_alt LIKE '%" + *searchLike + "%' OR"
		condition += " pba.bank_account_name LIKE '%" + *searchLike + "%' OR"
		condition += " bank_account_purpose.lkp_name LIKE '%" + *searchLike + "%' OR"
		condition += " bank.bank_fullname LIKE '%" + *searchLike + "%' OR"
		condition += " ba.account_no LIKE '%" + *searchLike + "%' OR"
		condition += " ba.account_holder_name LIKE '%" + *searchLike + "%')"
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
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
