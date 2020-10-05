package models

import (
	"api/db"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type GenLookupInfo struct {
	Name  string `json:"name"`
	Value uint64 `json:"value"`
}

type GenLookup struct {
	LookupKey         uint64  `db:"lookup_key"              json:"lookup_key"`
	LkpGroupKey       uint64  `db:"lkp_group_key"           json:"lkp_group_key"`
	LkpCode           *string `db:"lkp_code"                json:"lkp_code"`
	LkpName           *string `db:"lkp_name"                json:"lkp_name"`
	LkpDesc           *string `db:"lkp_desc"                json:"lkp_desc"`
	LkpVal1           *string `db:"lkp_val1"                json:"lkp_val1"`
	LkpVal2           *string `db:"lkp_val2"                json:"lkp_val2"`
	LkpVal3           *string `db:"lkp_val3"                json:"lkp_val3"`
	LkpText1          *string `db:"lkp_text1"               json:"lkp_text1"`
	LkpText2          *string `db:"lkp_text2"               json:"lkp_text2"`
	LkpText3          *string `db:"lkp_text3"               json:"lkp_text3"`
	RecOrder          *uint64 `db:"rec_order"               json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"              json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"              json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

func GetAllGenLookup(c *[]GenLookup, params map[string]string) (int, error) {
	query := `SELECT
              gen_lookup.* FROM 
			  gen_lookup `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "gen_lookup."+field+" = '"+value+"'")
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

func GetGenLookupIn(c *[]GenLookup, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				gen_lookup.* FROM 
				gen_lookup `
	query := query2 + " WHERE gen_lookup." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetGenLookup(c *GenLookup, key string) (int, error) {
	query := `SELECT gen_lookup.* FROM gen_lookup WHERE gen_lookup.lookup_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
