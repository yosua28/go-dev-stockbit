package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type GenLookupGroupList struct {
	GroupName string           `json:"group_name"`
	Lookup    *[]GenLookupInfo `json:"lookup,omitempty"`
}

type GenLookupGroup struct {
	LkpGroupKey       uint64  `db:"lkp_group_key"           json:"lkp_group_key"`
	LkpGroupCode      string  `db:"lkp_group_code"          json:"lkp_group_code"`
	LkpGroupName      string  `db:"lkp_group_name"          json:"lkp_group_name"`
	LkpGroupDesc      *string `db:"lkp_group_desc"          json:"lkp_group_desc"`
	LkpGroupPurpose   string  `db:"lkp_group_purpose"       json:"lkp_group_purpose"`
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

func GetAllGenLookupGroup(c *[]GenLookupGroup, params map[string]string) (int, error) {
	query := `SELECT
              gen_lookup_group.* FROM 
			  gen_lookup_group `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "gen_lookup_group."+field+" = '"+value+"'")
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

type ListDropdownLookupGroup struct {
	LkpGroupKey  uint64  `db:"lkp_group_key"        json:"lkp_group_key"`
	LkpGroupCode *string `db:"lkp_group_code"      json:"lkp_group_code"`
	LkpGroupName *string `db:"lkp_group_name"       json:"lkp_group_name"`
}

func AdminGetListDropdownLookupGroup(c *[]ListDropdownLookupGroup) (int, error) {
	query := `SELECT
				c.lkp_group_key,
				c.lkp_group_code,
				c.lkp_group_name 
			FROM gen_lookup_group AS c
			WHERE c.rec_status = 1 `
	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
