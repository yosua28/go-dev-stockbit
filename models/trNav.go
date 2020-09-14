package models

import (
	"api/db"
	"log"
	"strconv"
	"net/http"
	"strings"
)

type TrNavInfo struct {
	NavDate                   string    `json:"nav_date"`
	NavValue                  float32   `json:"nav_value"`
}

type TrNav struct {
	NavKey                    uint64    `db:"nav_key"                   json:"nav_key"`
	ProductKey                uint64    `db:"product_key"               json:"product_key"`
	NavDate                   string    `db:"nav_date"                  json:"nav_date"`
	NavValue                  float32   `db:"nav_value"                 json:"nav_value"`
	OriginalValue             float32   `db:"original_value"            json:"original_value"`
	NavStatus                 uint64    `db:"nav_status"                json:"nav_status"`
	ProdAumTotal              float32   `db:"prod_aum_total"            json:"prod_aum_total"`
	ProdUnitTotal             float32   `db:"prod_unit_total"           json:"prod_unit_total"`
	PublishMode               uint64    `db:"publish_mode"              json:"publish_mode"`
	RecOrder                  *uint64   `db:"rec_order"                 json:"rec_order"`
	RecStatus                 uint8     `db:"rec_status"                json:"rec_status"`
	RecCreatedDate            *string   `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy              *string   `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate           *string   `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy             *string   `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1                 *string   `db:"rec_image1"                json:"rec_image1"`
	RecImage2                 *string   `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus         *uint8    `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage          *uint64   `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate           *string   `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy             *string   `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate            *string   `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy              *string   `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1           *string   `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2           *string   `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3           *string   `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func GetAllTrNav(c *[]TrNav, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              tr_nav.* FROM 
			  tr_nav`
	var present bool
	var whereClause []string
	var condition string
	
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType"){
			whereClause = append(whereClause, "tr_nav."+field + " = '" + value + "'")
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

func GetTrNav(c *TrNav, key string) (int, error) {
	query := `SELECT tr_nav.* FROM tr_nav WHERE tr_nav.product_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetLastNavIn(c *[]TrNav, productKey []string,) (int, error) {
	inQuery := strings.Join(productKey, ",")
	query2 := `SELECT MAX( nav_key ) as nav_key, product_key, nav_date, nav_value FROM
			   tr_nav`
	query := query2 + " WHERE tr_nav.product_key IN(" + inQuery + ") GROUP BY product_key"
	
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetAllTrNavBetween(c *[]TrNav, start string, end string, productKey []string) (int, error) {
	inQuery := strings.Join(productKey, ",")
	query := `SELECT
              tr_nav.* FROM 
			  tr_nav`
	query += " WHERE tr_nav.product_key IN(" + inQuery + ") AND tr_nav.nav_date BETWEEN '" + start + "' AND '" + end + "'"
	
	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
