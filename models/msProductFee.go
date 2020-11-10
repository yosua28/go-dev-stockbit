package models

import (
	"api/db"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MsProductFeeInfo struct {
	FeeAnnotation string                 `json:"fee_annotation"`
	FeeDesc       string                 `json:"fee_desc"`
	FeeCode       string                 `json:"fee_code"`
	FeeType       uint64                 `json:"fee_type"`
	FlagShowOntnc uint8                  `json:"flag_show_ontnc"`
	FeeItem       []MsProductFeeItemInfo `json:"fee_item"`
}

type MsProductFee struct {
	FeeKey            uint64   `db:"fee_key"               json:"fee_key"`
	ProductKey        uint64   `db:"product_key"           json:"product_key"`
	FeeType           *uint64  `db:"fee_type"              json:"fee_type"`
	FeeCode           *string  `db:"fee_code"              json:"fee_code"`
	FlagShowOntnc     *uint8   `db:"flag_show_ontnc"       json:"flag_show_ontnc"`
	FeeAnnotation     *string  `db:"fee_annotation"        json:"fee_annotation"`
	FeeDesc           *string  `db:"fee_desc"              json:"fee_desc"`
	FeeDateStart      *string  `db:"fee_date_start"        json:"fee_date_start"`
	FeeDateThru       *string  `db:"fee_date_thru"         json:"fee_date_thru"`
	FeeNominalType    *uint64  `db:"fee_nominal_type"      json:"fee_nominal_type"`
	EnabledMinAmount  uint8    `db:"enabled_min_amount"    json:"enabled_min_amount"`
	FeeMinAmount      *float64 `db:"fee_min_amount"        json:"fee_min_amount"`
	EnabledMaxAmount  uint8    `db:"enabled_max_amount"    json:"enabled_max_amount"`
	FeeMaxAmount      *float64 `db:"fee_max_amount"        json:"fee_max_amount"`
	FeeCalcMethod     *uint64  `db:"fee_calc_method"       json:"fee_calc_method"`
	CalculationBaseon *uint64  `db:"calculation_baseon"    json:"calculation_baseon"`
	PeriodHold        uint64   `db:"period_hold"           json:"period_hold"`
	DaysInyear        *uint64  `db:"days_inyear"           json:"days_inyear"`
	RecOrder          *uint64  `db:"rec_order"             json:"rec_order"`
	RecStatus         uint8    `db:"rec_status"            json:"rec_status"`
	RecCreatedDate    *string  `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy      *string  `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate   *string  `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy     *string  `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1         *string  `db:"rec_image1"            json:"rec_image1"`
	RecImage2         *string  `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus *uint8   `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage  *uint64  `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate   *string  `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy     *string  `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate    *string  `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy      *string  `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1   *string  `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2   *string  `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3   *string  `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type MsProductFeeDetailAdmin struct {
	FeeKey            uint64                        `json:"fee_key"`
	Product           MsProductListDropdown         `json:"product"`
	FeeType           *LookupTrans                  `json:"fee_type"`
	FeeCode           *string                       `json:"fee_code"`
	FlagShowOntnc     bool                          `json:"flag_show_ontnc"`
	FeeAnnotation     *string                       `json:"fee_annotation"`
	FeeDesc           *string                       `json:"fee_desc"`
	FeeDateStart      *string                       `json:"fee_date_start"`
	FeeDateThru       *string                       `json:"fee_date_thru"`
	FeeNominalType    *LookupTrans                  `json:"fee_nominal_type"`
	EnabledMinAmount  bool                          `json:"enabled_min_amount"`
	FeeMinAmount      *float64                      `json:"fee_min_amount"`
	EnabledMaxAmount  bool                          `json:"enabled_max_amount"`
	FeeMaxAmount      *float64                      `json:"fee_max_amount"`
	FeeCalcMethod     *LookupTrans                  `json:"fee_calc_method"`
	CalculationBaseon *LookupTrans                  `json:"calculation_baseon"`
	PeriodHold        uint64                        `json:"period_hold"`
	DaysInyear        *LookupTrans                  `json:"days_inyear"`
	ProductFeeItems   *[]MsProductFeeItemDetailList `json:"product_fee_items"`
}

type AdminListMsProductFee struct {
	FeeKey       uint64  `db:"fee_key"               json:"fee_key"`
	FeeCode      *string `db:"fee_code"              json:"fee_code"`
	ProductKey   uint64  `db:"product_key"           json:"product_key"`
	ProductCode  string  `db:"product_code"          json:"product_code"`
	ProductName  string  `db:"product_name"          json:"product_name"`
	FeeTypeName  *string `db:"feetypename"           :"feetypename"`
	FeeDateStart *string `db:"fee_date_start"        json:"fee_date_start"`
	FeeDateThru  *string `db:"fee_date_thru"         json:"fee_date_thru"`
	PeriodHold   uint64  `db:"period_hold"           json:"period_hold"`
}

func GetAllMsProductFee(c *[]MsProductFee, params map[string]string) (int, error) {
	query := `SELECT
              ms_product_fee.* FROM 
			  ms_product_fee WHERE 
			  ms_product_fee.fee_date_start <= NOW() AND 
			  ms_product_fee.fee_date_thru > NOW() AND 
			  ms_product_fee.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product_fee."+field+" = '"+value+"'")
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

func AdminGetAllMsProductFee(c *[]AdminListMsProductFee, limit uint64, offset uint64, params map[string]string, nolimit bool, searchLike *string) (int, error) {
	query := `SELECT
				pf.fee_key AS fee_key, 
				pf.fee_code AS fee_code, 
				p.product_key AS product_key, 
				p.product_code AS product_code, 
				p.product_name AS product_name, 
				feetype.lkp_name AS feetypename, 
				DATE_FORMAT(pf.fee_date_start, '%d %M %Y') AS fee_date_start, 
				DATE_FORMAT(pf.fee_date_thru, '%d %M %Y') AS fee_date_thru,
				pf.period_hold AS period_hold 
			  FROM ms_product_fee AS pf
			  INNER JOIN ms_product AS p ON p.product_key = pf.product_key
			  LEFT JOIN gen_lookup AS feetype ON feetype.lookup_key = pf.fee_type
			  WHERE pf.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "pf."+field+" = '"+value+"'")
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
		condition += " pf.fee_key LIKE '%" + *searchLike + "%' OR"
		condition += " pf.fee_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_name LIKE '%" + *searchLike + "%' OR"
		condition += " feetype.lkp_name LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(pf.fee_date_start, '%d %M %Y') LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(pf.fee_date_thru, '%d %M %Y') LIKE '%" + *searchLike + "%' OR"
		condition += " pf.period_hold LIKE '%" + *searchLike + "%')"
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

func AdminCountDataGetAllMsProductFee(c *CountData, params map[string]string, searchLike *string) (int, error) {
	query := `SELECT
				count(pf.fee_key) AS count_data
			  FROM ms_product_fee AS pf
			  INNER JOIN ms_product AS p ON p.product_key = pf.product_key
			  LEFT JOIN gen_lookup AS feetype ON feetype.lookup_key = pf.fee_type
			  WHERE pf.rec_status = 1`
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "pf."+field+" = '"+value+"'")
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
		condition += " pf.fee_key LIKE '%" + *searchLike + "%' OR"
		condition += " pf.fee_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_code LIKE '%" + *searchLike + "%' OR"
		condition += " p.product_name LIKE '%" + *searchLike + "%' OR"
		condition += " feetype.lkp_name LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(pf.fee_date_start, '%d %M %Y') LIKE '%" + *searchLike + "%' OR"
		condition += " DATE_FORMAT(pf.fee_date_thru, '%d %M %Y') LIKE '%" + *searchLike + "%' OR"
		condition += " pf.period_hold LIKE '%" + *searchLike + "%')"
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

func GetMsProductFee(c *MsProductFee, key string) (int, error) {
	query := `SELECT ms_product_fee.* FROM ms_product_fee WHERE ms_product_fee.rec_status = 1 AND ms_product_fee.fee_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
