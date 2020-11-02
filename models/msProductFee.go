package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type MsProductFeeInfo struct {
	FeeAnnotation           string                  `json:"fee_annotation"`
	FeeDesc                 string                  `json:"fee_desc"`
	FeeCode                 string                  `json:"fee_code"`
	FeeType                 uint64                  `json:"fee_type"`
	FlagShowOntnc           uint8                   `json:"flag_show_ontnc"`
	FeeItem                 []MsProductFeeItemInfo  `json:"fee_item"`
}

type MsProductFee struct {
	FeeKey               uint64     `db:"fee_key"               json:"fee_key"`
	ProductKey           uint64     `db:"product_key"           json:"product_key"`
	FeeType              *uint64    `db:"fee_type"              json:"fee_type"`
	FeeCode              *string    `db:"fee_code"              json:"fee_code"`
	FlagShowOntnc        *uint8     `db:"flag_show_ontnc"       json:"flag_show_ontnc"`
	FeeAnnotation        *string    `db:"fee_annotation"        json:"fee_annotation"`
	FeeDesc              *string    `db:"fee_desc"              json:"fee_desc"`
	FeeDateStart         *string    `db:"fee_date_start"        json:"fee_date_start"`
	FeeDateThru          *string    `db:"fee_date_thru"         json:"fee_date_thru"`
	FeeNominalType       *uint64    `db:"fee_nominal_type"      json:"fee_nominal_type"`
	EnabledMinAmount      uint8     `db:"enabled_min_amount"    json:"enabled_min_amount"`
	FeeMinAmount         *float64   `db:"fee_min_amount"        json:"fee_min_amount"`
	EnabledMaxAmount      uint8     `db:"enabled_max_amount"    json:"enabled_max_amount"`
	FeeMaxAmount         *float64   `db:"fee_max_amount"        json:"fee_max_amount"`
	FeeCalcMethod        *uint64    `db:"fee_calc_method"       json:"fee_calc_method"`
	CalculationBaseon    *uint64    `db:"calculation_baseon"    json:"calculation_baseon"`
	PeriodHold            uint64    `db:"period_hold"           json:"period_hold"`
	DaysInyear           *uint64    `db:"days_inyear"           json:"days_inyear"`
	RecOrder             *uint64    `db:"rec_order"             json:"rec_order"`
	RecStatus            uint8      `db:"rec_status"            json:"rec_status"`
	RecCreatedDate       *string    `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy         *string    `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate      *string    `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy        *string    `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1            *string    `db:"rec_image1"            json:"rec_image1"`
	RecImage2            *string    `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus    *uint8     `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage     *uint64    `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate      *string    `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy        *string    `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate       *string    `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy         *string    `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1      *string    `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2      *string    `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3      *string    `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
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
