package models

import (
	"api/db"
	"log"
	"net/http"
	"strconv"
)

type CmsFinancialCalc struct {
	CalcKey           uint64   `db:"calc_key"                  json:"calc_key"`
	UserLoginKey      uint64   `db:"user_login_key"            json:"user_login_key"`
	InvestPurposeKey  uint64   `db:"invest_purpose_key"        json:"invest_purpose_key"`
	InvestTarget      *float64 `db:"invest_target"             json:"invest_target"`
	InvestReturn      *float64 `db:"invest_return"             json:"invest_return"`
	InvestPeriod      *float64 `db:"invest_period"             json:"invest_period"`
	ResultInvest      *float64 `db:"result_invest"             json:"result_invest"`
	ResultRemakrs     *string  `db:"result_remakrs"            json:"result_remakrs"`
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

type CmsFinancialCalcList struct {
	CalcKey          uint64   `db:"calc_key"                json:"calc_key"`
	UserLoginKey     uint64   `db:"user_login_key"          json:"user_login_key"`
	UserLoginName    string   `db:"user_login_name"         json:"user_login_name"`
	CustomerName     *string  `db:"customer_name"           json:"customer_name"`
	CustomerRole     *string  `db:"customer_role"           json:"customer_role"`
	InvestPurposeKey string   `db:"invest_purpose_key"      json:"invest_purpose_key"`
	PurposeCode      string   `db:"purpose_code"            json:"purpose_code"`
	PurposeName      *string  `db:"purpose_name"            json:"purpose_name"`
	PurposeDesc      *string  `db:"purpose_desc"            json:"purpose_desc"`
	InvestTarget     *float64 `db:"invest_target"           json:"invest_target"`
	InvestReturn     *float64 `db:"invest_return"           json:"invest_return"`
	InvestPeriod     *float64 `db:"invest_period"           json:"invest_period"`
	ResultInvest     *float64 `db:"result_invest"           json:"result_invest"`
	ResultRemakrs    *string  `db:"result_remakrs"          json:"result_remakrs"`
}

type CmsFinancialCalcDetail struct {
	CalcKey       uint64               `json:"invest_purpose_key"`
	UserLogin     UserLogin            `json:"user_login"`
	InvestPurpose CmsInvestPurposeList `json:"invest_purpose"`
	InvestTarget  *float64             `json:"invest_target"`
	InvestReturn  *float64             `json:"invest_return"`
	InvestPeriod  *float64             `json:"invest_period"`
	ResultInvest  *float64             `json:"result_invest"`
	ResultRemakrs *string              `json:"result_remakrs"`
}

type UserLogin struct {
	UserLoginKey  uint64  `json:"user_login_key"`
	UserLoginName string  `json:"user_login_name"`
	CustomerKey   *string `json:"customer_key"`
	CustomerName  *string `json:"customer_name"`
	RoleKey       *string `json:"role_key"`
	RoleName      *string `json:"role_name"`
}

type CmsFinancialCalcCount struct {
	CountData int `db:"count_data"             json:"count_data"`
}

func GetAllCmsFinancialCalc(c *[]CmsFinancialCalcList, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT 
				cfc.calc_key AS calc_key,
				cfc.user_login_key AS user_login_key,
				sul.ulogin_name AS user_login_name,
				cus.full_name AS customer_name,
				role.role_code AS customer_role,
				cfc.invest_purpose_key AS invest_purpose_key,
				cip.purpose_code AS purpose_code,
				cip.purpose_name AS purpose_name,
				cip.purpose_desc AS purpose_desc,
				cfc.invest_target AS invest_target,
				cfc.invest_return AS invest_return,
				cfc.invest_period AS invest_period,
				cfc.result_invest AS result_invest,
				cfc.result_remakrs AS result_remakrs
			FROM cms_financial_calc AS cfc 
			INNER JOIN sc_user_login AS sul ON sul.user_login_key = cfc.user_login_key
			INNER JOIN cms_invest_purpose cip ON cip.invest_purpose_key = cfc.invest_purpose_key
			LEFT JOIN ms_customer AS cus ON sul.customer_key = cus.customer_key
			LEFT JOIN sc_role AS role ON sul.role_key = role.role_key
			WHERE cfc.rec_status = 1 `
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

func GetCountCmsFinancialCalc(c *CmsFinancialCalcCount, params map[string]string) (int, error) {
	query := `SELECT count(cfc.calc_key) as count_data
			FROM cms_financial_calc AS cfc 
			INNER JOIN sc_user_login AS sul ON sul.user_login_key = cfc.user_login_key
			INNER JOIN cms_invest_purpose cip ON cip.invest_purpose_key = cfc.invest_purpose_key
			LEFT JOIN ms_customer AS cus ON sul.customer_key = cus.customer_key
			LEFT JOIN sc_role AS role ON sul.role_key = role.role_key
			WHERE cfc.rec_status = 1 `

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
