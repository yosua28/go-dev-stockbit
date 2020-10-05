package models

import (
	"api/db"
	"log"
	"strconv"
	"net/http"
)

type CmsInvestPartnerList struct{
	InvestPartnerKey          uint64    `json:"invest_partner_key"`
	PartnerCode               string    `json:"partner_code"`
	PartnerBusinessName       string    `json:"partner_business_name"`
	PartnerUrl  	          string    `json:"partner_url"`
	RecImage1                 string    `json:"rec_image1"`
}

type CmsInvestPartner struct{
	InvestPartnerKey          uint64    `db:"invest_partner_key"        json:"invest_partner_key"`
	InvestPurposeKey         *uint64    `db:"invest_purpose_key"        json:"invest_purpose_key"`
	PartnerCode               string    `db:"partner_code"              json:"partner_code"`
	PartnerBusinessName      *string    `db:"partner_business_name"     json:"partner_business_name"`
	PartnerDesc              *string    `db:"partner_desc"              json:"partner_desc"`
	PartnerPicname           *string    `db:"partner_picname"           json:"partner_picname"`
	PartnerMobileno          *string    `db:"partner_mobileno"          json:"partner_mobileno"`
	PartnerOfficeno          *string    `db:"partner_officeno"          json:"partner_officeno"`
	PartnerEmail   	         *string    `db:"partner_email"             json:"partner_email"`
	PartnerCity   	         *string    `db:"partner_city"              json:"partner_city"`
	PartnerAddress	         *string    `db:"partner_address"           json:"partner_address"`
	PartnerUrl  	         *string    `db:"partner_url"               json:"partner_url"`
	PartnerDateStarted       *string    `db:"partner_date_started"      json:"partner_date_started"`
	PartnerDateExpired       *string    `db:"partner_date_expired"      json:"partner_date_expired"`
	PartnerBannerHits        *string    `db:"partner_banner_hits"       json:"partner_banner_hits"`
	RecOrder                 *uint64    `db:"rec_order"                 json:"rec_order"`
	RecStatus                 uint8     `db:"rec_status"                json:"rec_status"`
	RecCreatedDate           *string    `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy             *string    `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate          *string    `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy            *string    `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1                *string    `db:"rec_image1"                json:"rec_image1"`
	RecImage2                *string    `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus        *uint8     `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage         *uint64    `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate          *string    `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy            *string    `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate           *string    `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy             *string    `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1          *string    `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2          *string    `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3          *string    `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func GetAllCmsInvestPartner(c *[]CmsInvestPartner, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              cms_invest_partner.* FROM 
			  cms_invest_partner WHERE 
			  cms_invest_partner.partner_date_started <= NOW() AND 
			  cms_invest_partner.partner_date_expired > NOW() `
	var present bool
	var whereClause []string
	var condition string
	
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType"){
			whereClause = append(whereClause, "cms_invest_partner."+field + " = '" + value +"'")
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