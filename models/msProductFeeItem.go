package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsProductFeeItemInfo struct {
	ItemSeqno      uint64  `json:"item_seqno"`
	RowMax         uint8   `json:"row_max"`
	PrincipleLimit float64 `json:"principle_limit"`
	FeeValue       float64 `json:"fee_value"`
	ItemNotes      string  `json:"item_notes"`
}

type MsProductFeeItem struct {
	ProductFeeItemKey uint64  `db:"product_fee_item_key"  json:"product_fee_item_key"`
	ProductFeeKey     uint64  `db:"product_fee_key"       json:"product_fee_key"`
	ItemSeqno         uint64  `db:"item_seqno"            json:"item_seqno"`
	RowMax            uint8   `db:"row_max"               json:"row_max"`
	PrincipleLimit    float64 `db:"principle_limit"       json:"principle_limit"`
	FeeValue          float64 `db:"fee_value"             json:"fee_value"`
	ItemNotes         *string `db:"item_notes"            json:"item_notes"`
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

type MsProductFeeItemDetailList struct {
	ProductFeeItemKey uint64  `json:"product_fee_item_key"`
	PrincipleLimit    float64 `json:"principle_limit"`
	FeeValue          float64 `json:"fee_value"`
	ItemNotes         *string `json:"item_notes"`
}

func GetAllMsProductFeeItem(c *[]MsProductFeeItem, params map[string]string) (int, error) {
	query := `SELECT
              ms_product_fee_item.* FROM 
			  ms_product_fee_item `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_product_fee_item."+field+" = '"+value+"'")
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

func GetMsProductFeeItemIn(c *[]MsProductFeeItem, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_product_fee_item.* FROM 
				ms_product_fee_item WHERE
				ms_product_fee_item.rec_status = 1 `
	query := query2 + " AND ms_product_fee_item." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateMsProductFeeItemByField(params map[string]string, value string, field string) (int, error) {
	query := "UPDATE ms_product_fee_item SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"
		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE " + field + " = " + value
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

func CreateMsProductFeeItem(params map[string]string) (int, error) {
	query := "INSERT INTO ms_product_fee_item"
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

func GetLastMsProductFeeItemByFeeKey(c *MsProductFeeItem, feeKey string, fieldOrder string, orderByType string) (int, error) {
	query := `SELECT ms_product_fee_item.* 
			FROM ms_product_fee_item 
			WHERE ms_product_fee_item.rec_status = 1 AND 
			ms_product_fee_item.product_fee_key = ` + feeKey + ` ORDER BY ` + fieldOrder + ` ` + orderByType + ` LIMIT 1`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsProductFeeItem(c *MsProductFeeItem, key string) (int, error) {
	query := `SELECT ms_product_fee_item.* FROM ms_product_fee_item 
				WHERE ms_product_fee_item.rec_status = 1 AND ms_product_fee_item.product_fee_item_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsProductFeeItemCalculateFifoWithLimit(c *MsProductFeeItem, productKey string, pLimit string, feeType string) (int, error) {
	query := `SELECT i.* 
			FROM ms_product_fee_item AS i 
			JOIN ms_product_fee AS pf ON pf.fee_key = i.product_fee_key 
			WHERE i.rec_status = 1 AND i.principle_limit <= ` + pLimit + ` 
			AND pf.product_key = ` + productKey + ` AND pf.fee_type = ` + feeType + ` 
			ORDER BY i.principle_limit ASC LIMIT 1`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsProductFeeItemLastCalculateFifo(c *MsProductFeeItem, productKey string, feeType string) (int, error) {
	query := `SELECT i.* 
			FROM ms_product_fee_item AS i 
			JOIN ms_product_fee AS pf ON pf.fee_key = i.product_fee_key 
			WHERE i.rec_status = 1 AND pf.product_key = ` + productKey + ` AND pf.fee_type = ` + feeType + ` 
			ORDER BY i.principle_limit ASC LIMIT 1`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
