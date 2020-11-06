package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"
)

type MsProductCategory struct {
	ProductCategoryKey uint64  `db:"product_category_key"  json:"product_category_key"`
	CategoryCode       *string `db:"category_code"         json:"category_code"`
	CategoryName       *string `db:"category_name"         json:"category_name"`
	CategoryDesc       *string `db:"category_desc"         json:"category_desc"`
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

type MsProductCategoryInfo struct {
	ProductCategoryKey uint64  `json:"product_category_key"`
	CategoryCode       *string `json:"category_code"`
	CategoryName       *string `json:"category_name"`
	CategoryDesc       *string `json:"category_desc"`
}

func GetMsProductCategoryIn(c *[]MsProductCategory, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_product_category.* FROM 
				ms_product_category `
	query := query2 + " WHERE ms_product_category.rec_status = 1 AND ms_product_category." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsProductCategory(c *MsProductCategory, key string) (int, error) {
	query := `SELECT ms_product_category.* FROM ms_product_category 
				WHERE ms_product_category.rec_status = 1 AND ms_product_category.product_category_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
