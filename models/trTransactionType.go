package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"
)

type TrTransactionType struct {
	TransTypeKey      uint64  `db:"trans_type_key"          json:"trans_type_key"`
	TypeCode          *string `db:"type_code"               json:"type_code"`
	TypeDescription   *string `db:"type_description"        json:"type_description"`
	TypeOrder         uint64  `db:"type_order"              json:"type_order"`
	TypeDomain        *string `db:"type_domain"            json:"type_domain"`
	RecOrder          *uint64 `db:"rec_order"            json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"           json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"     json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"       json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"    json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"      json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"           json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"           json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"  json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"   json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"    json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"      json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"     json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"       json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"    json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"    json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"    json:"rec_attribute_id3"`
}

func GetMsTransactionTypeIn(c *[]TrTransactionType, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				tr_transaction_type.* FROM 
				tr_transaction_type `
	query := query2 + " WHERE tr_transaction_type.rec_status = 1 AND tr_transaction_type." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsTransactionType(c *TrTransactionType, key string) (int, error) {
	query := `SELECT tr_transaction_type.* FROM tr_transaction_type WHERE tr_transaction_type.rec_status = 1 AND tr_transaction_type.trans_type_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
