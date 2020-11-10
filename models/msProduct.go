package models

import (
	"api/db"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type MsProductList struct {
	ProductKey     uint64                 `json:"product_key"`
	ProductID      uint64                 `json:"product_id"`
	ProductCode    string                 `json:"product_code"`
	ProductName    string                 `json:"product_name"`
	ProductNameAlt string                 `json:"product_name_alt"`
	MinSubAmount   float32                `json:"min_sub_amount"`
	RecImage1      string                 `json:"rec_image1"`
	FundType       *MsFundTypeInfo        `json:"fund_type,omitempty"`
	NavPerformance *FfsNavPerformanceInfo `json:"nav_performance,omitempty"`
	Nav            *TrNavInfo             `json:"nav,omitempty"`
	RiskProfile    *MsRiskProfileInfo     `json:"risk_profile,omitempty"`
}

type MsProductData struct {
	ProductKey       uint64                     `json:"product_key"`
	ProductID        uint64                     `json:"product_id"`
	ProductCode      string                     `json:"product_code"`
	ProductName      string                     `json:"product_name"`
	ProductNameAlt   string                     `json:"product_name_alt"`
	MinSubAmount     float32                    `json:"min_sub_amount"`
	MinRedAmount     float32                    `json:"min_red_amount"`
	MinRedUnit       float32                    `json:"min_red_unit"`
	MinUnitAfterRed  float32                    `json:"min_unit_after_red"`
	ProspectusLink   string                     `json:"prospectus_link"`
	FundFactSheet    string                     `json:"ffs_link"`
	RecImage1        string                     `json:"rec_image1"`
	FlagSubscription bool                       `json:"flag_subscription"`
	FlagRedemption   bool                       `json:"flag_redemption"`
	FlagSwitchOut    bool                       `json:"flag_switch_out"`
	FlagSwitchIn     bool                       `json:"flag_switch_in"`
	FeeService       string                     `json:"fee_service"`
	FeeTransfer      string                     `json:"fee_transfer"`
	InvestValue      string                     `json:"invest_value"`
	BalanceUnit      float32                    `json:"balance_unit"`
	Currency         MsCurrencyInfo             `json:"currency"`
	BankAcc          []MsProductBankAccountInfo `json:"bank_account"`
	ProductFee       []MsProductFeeInfo         `json:"product_fee"`
	NavPerformance   *FfsNavPerformanceInfo     `json:"nav_performance,omitempty"`
	Nav              *TrNavInfo                 `json:"nav,omitempty"`
	CustodianBank    *MsCustodianBankInfo       `json:"custodian_bank,omitempty"`
	RiskProfile      *MsRiskProfileInfo         `json:"risk_profile,omitempty"`
}

type MsProduct struct {
	ProductKey            uint64  `db:"product_key"             json:"product_key"`
	ProductID             uint64  `db:"product_id"              json:"product_id"`
	ProductCode           string  `db:"product_code"            json:"product_code"`
	ProductName           string  `db:"product_name"            json:"product_name"`
	ProductNameAlt        string  `db:"product_name_alt"        json:"product_name_alt"`
	CurrencyKey           *uint64 `db:"currency_key"            json:"currency_key"`
	ProductCategoryKey    *uint64 `db:"product_category_key"    json:"product_category_key"`
	ProductTypeKey        *uint64 `db:"product_type_key"        json:"product_type_key"`
	FundTypeKey           *uint64 `db:"fund_type_key"           json:"fund_type_key"`
	FundStructureKey      *uint64 `db:"fund_structure_key"      json:"fund_structure_key"`
	RiskProfileKey        *uint64 `db:"risk_profile_key"        json:"risk_profile_key"`
	ProductProfile        *string `db:"product_profile"         json:"product_profile"`
	InvestmentObjectives  *string `db:"investment_objectives"   json:"investment_objectives"`
	ProductPhase          *uint64 `db:"product_phase"           json:"product_phase"`
	NavValuationType      *uint64 `db:"nav_valuation_type"      json:"nav_valuation_type"`
	ProspectusLink        *string `db:"prospectus_link"         json:"prospectus_link"`
	LaunchDate            *string `db:"launch_date"             json:"launch_date"`
	InceptionDate         *string `db:"inception_date"          json:"inception_date"`
	IsinCode              *string `db:"isin_code"               json:"isin_code"`
	FlagSyariah           uint8   `db:"flag_syariah"            json:"flag_syariah"`
	MaxSubFee             float32 `db:"max_sub_fee"             json:"max_sub_fee"`
	MaxRedFee             float32 `db:"max_red_fee"             json:"max_red_fee"`
	MaxSwiFee             float32 `db:"max_swi_fee"             json:"max_swi_fee"`
	MinSubAmount          float32 `db:"min_sub_amount"          json:"min_sub_amount"`
	MinRedAmount          float32 `db:"min_red_amount"          json:"min_red_amount"`
	MinRedUnit            float32 `db:"min_red_unit"            json:"min_red_unit"`
	MinUnitAfterRed       float32 `db:"min_unit_after_red"      json:"min_unit_after_red"`
	ManagementFee         float32 `db:"management_fee"          json:"management_fee"`
	CustodianFee          float32 `db:"custodian_fee"           json:"custodian_fee"`
	CustodianKey          *uint64 `db:"custodian_key"           json:"custodian_key"`
	OjkFee                float32 `db:"ojk_fee"                 json:"ojk_fee"`
	ProductFeeAmount      float32 `db:"product_fee_amount"      json:"product_fee_amount"`
	OverwriteTransactFlag uint8   `db:"overwrite_transact_flag" json:"overwrite_transact_flag"`
	OverwriteFeeFlag      uint8   `db:"overwrite_fee_flag"      json:"overwrite_fee_flag"`
	OtherFeeAmount        float32 `db:"other_fee_amount"        json:"other_fee_amount"`
	SettlementPeriod      *uint64 `db:"settlement_period"       json:"settlement_period"`
	SinvestFundCode       *string `db:"sinvest_fund_code"       json:"sinvest_fund_code"`
	FlagEnabled           uint8   `db:"flag_enabled"            json:"flag_enabled"`
	FlagSubscription      uint8   `db:"flag_subscription"       json:"flag_subscription"`
	FlagRedemption        uint8   `db:"flag_redemption"         json:"flag_redemption"`
	FlagSwitchOut         uint8   `db:"flag_switch_out"         json:"flag_switch_out"`
	FlagSwitchIn          uint8   `db:"flag_switch_in"          json:"flag_switch_in"`
	RecOrder              *uint64 `db:"rec_order"               json:"rec_order"`
	RecStatus             uint8   `db:"rec_status"              json:"rec_status"`
	RecCreatedDate        *string `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy          *string `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate       *string `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy         *string `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1             *string `db:"rec_image1"              json:"rec_image1"`
	RecImage2             *string `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus     *uint8  `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage      *uint64 `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate       *string `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy         *string `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate        *string `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy          *string `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1       *string `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2       *string `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3       *string `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

type AdminMsProductList struct {
	ProductKey          uint64  `json:"product_key"`
	ProductID           uint64  `json:"product_id"`
	ProductCode         string  `json:"product_code"`
	ProductName         string  `json:"product_name"`
	ProductNameAlt      string  `json:"product_name_alt"`
	CurrencyName        *string `json:"currency_name"`
	ProductCategoryName *string `json:"product_category_name"`
	ProductTypeName     *string `json:"product_type_name"`
	RiskProfileName     *string `json:"risk_profile_name"`
	LaunchDate          *string `json:"launch_date"`
	InceptionDate       *string `json:"inception_date"`
	IsinCode            *string `json:"isin_code"`
	Syariah             string  `json:"syariah"`
	CustodianFullName   *string `json:"custodian_full_name"`
	SinvestFundCode     *string `json:"sinvest_fund_code"`
	Enabled             string  `json:"enabled"`
	Subscription        string  `json:"subscription"`
	Redemption          string  `json:"redemption"`
	SwitchOut           string  `json:"switch_out"`
	SwitchIn            string  `json:"switch_in"`
}

func CreateMsProduct(params map[string]string) (int, error) {
	query := "INSERT INTO ms_product"
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

type AdminMsProductDetail struct {
	ProductKey            uint64                   `json:"product_key"`
	ProductID             uint64                   `json:"product_id"`
	ProductCode           string                   `json:"product_code"`
	ProductName           string                   `json:"product_name"`
	ProductNameAlt        string                   `json:"product_name_alt"`
	Currency              *MsCurrencyInfo          `json:"currency"`
	ProductCategory       *MsProductCategoryInfo   `json:"product_category"`
	ProductType           *MsProductTypeInfo       `json:"product_type"`
	FundType              *MsFundTypeInfo          `json:"fund_type"`
	FundStructure         *MsFundStructureInfo     `json:"fund_structure"`
	RiskProfile           *MsRiskProfileInfoAdmin  `json:"risk_profile"`
	ProductProfile        *string                  `json:"product_profile"`
	InvestmentObjectives  *string                  `json:"investment_objectives"`
	ProductPhase          *LookupTrans             `json:"product_phase"`
	NavValuationType      *LookupTrans             `json:"nav_valuation_type"`
	ProspectusLink        *string                  `json:"prospectus_link"`
	LaunchDate            *string                  `json:"launch_date"`
	InceptionDate         *string                  `json:"inception_date"`
	IsinCode              *string                  `json:"isin_code"`
	FlagSyariah           bool                     `json:"flag_syariah"`
	MaxSubFee             float32                  `json:"max_sub_fee"`
	MaxRedFee             float32                  `json:"max_red_fee"`
	MaxSwiFee             float32                  `json:"max_swi_fee"`
	MinSubAmount          float32                  `json:"min_sub_amount"`
	MinRedAmount          float32                  `json:"min_red_amount"`
	MinRedUnit            float32                  `json:"min_red_unit"`
	MinUnitAfterRed       float32                  `json:"min_unit_after_red"`
	ManagementFee         float32                  `json:"management_fee"`
	CustodianFee          float32                  `json:"custodian_fee"`
	Custodian             *MsCustodianBankInfoList `json:"custodian"`
	OjkFee                float32                  `json:"ojk_fee"`
	ProductFeeAmount      float32                  `json:"product_fee_amount"`
	OverwriteTransactFlag bool                     `json:"overwrite_transact_flag"`
	OverwriteFeeFlag      bool                     `json:"overwrite_fee_flag"`
	OtherFeeAmount        float32                  `json:"other_fee_amount"`
	SettlementPeriod      *uint64                  `json:"settlement_period"`
	SinvestFundCode       *string                  `json:"sinvest_fund_code"`
	FlagEnabled           bool                     `json:"flag_enabled"`
	FlagSubscription      bool                     `json:"flag_subscription"`
	FlagRedemption        bool                     `json:"flag_redemption"`
	FlagSwitchOut         bool                     `json:"flag_switch_out"`
	FlagSwitchIn          bool                     `json:"flag_switch_in"`
}

func GetAllMsProduct(c *[]MsProduct, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ms_product.* FROM 
			  ms_product `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product."+field+" = '"+value+"'")
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

func GetMsProduct(c *MsProduct, key string) (int, error) {
	query := `SELECT ms_product.* FROM ms_product WHERE ms_product.product_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsProductIn(c *[]MsProduct, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_product.* FROM 
				ms_product WHERE 
				ms_product.rec_status = 1 `
	query := query2 + " AND ms_product." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetAllMsProductWithLike(c *[]MsProduct, limit uint64, offset uint64, params map[string]string, paramsLike map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ms_product.* FROM 
			  ms_product `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product."+field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, "ms_product."+fieldLike+" like '%"+valueLike+"%'")
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

func AdminGetCountMsProductWithLike(c *CountData, params map[string]string, paramsLike map[string]string) (int, error) {
	query := `SELECT
			  count(ms_product.product_key) as count_data 
			  FROM ms_product `
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product."+field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, "ms_product."+fieldLike+" like '%"+valueLike+"%'")
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

func UpdateMsProduct(params map[string]string) (int, error) {
	query := "UPDATE ms_product SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "product_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE product_key = " + params["product_key"]
	log.Println(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func AdminGetValidateUniqueDataInsertUpdate(c *CountData, paramsOr map[string]string, paramsAnd map[string]string, updateKey *string) (int, error) {
	query := `SELECT
			  count(ms_product.product_key) as count_data 
			  FROM ms_product `
	var orWhereClause []string
	var andWhereClause []string
	var condition string

	for fieldOr, valueOr := range paramsOr {
		orWhereClause = append(orWhereClause, "ms_product."+fieldOr+" = '"+valueOr+"'")
	}

	for fieldAnd, valueAnd := range paramsAnd {
		andWhereClause = append(andWhereClause, "ms_product."+fieldAnd+" like '"+valueAnd+"'")
	}

	// Combile where Or clause
	if len(orWhereClause) > 0 {
		condition += " WHERE ("
		for index, where := range orWhereClause {
			condition += where
			if (len(orWhereClause) - 1) > index {
				condition += " OR "
			} else {
				condition += ") "
			}
		}
	}

	// Combile where And clause
	if len(andWhereClause) > 0 {
		if len(orWhereClause) > 0 {
			condition += " AND "
		} else {
			condition += " WHERE "
		}

		for index, where := range andWhereClause {
			condition += where
			if (len(andWhereClause) - 1) > index {
				condition += " AND "
			}
		}
	}

	if updateKey != nil {
		if len(orWhereClause) > 0 {
			condition += " AND "
		} else if len(andWhereClause) > 0 {
			condition += " AND "
		} else {
			condition += " WHERE "
		}

		condition += " ms_product.product_key != '" + *updateKey + "'"
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
