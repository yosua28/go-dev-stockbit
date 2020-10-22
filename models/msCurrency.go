package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"
)

type MsCurrencyInfo struct {
	CurrencyKey       uint64  `json:"currency_key"`
	Code              string  `json:"code"`
	Symbol            *string `json:"symbol"`
	Name              *string `json:"name"`
	FlagBase          uint8   `json:"flag_base"`
}

type MsCurrency struct {
	CurrencyKey       uint64  `db:"currency_key"          json:"currency_key"`
	Code              string  `db:"code"                  json:"code"`
	Symbol            *string `db:"symbol"                json:"symbol"`
	Name              *string `db:"name"                  json:"name"`
	FlagBase          uint8   `db:"flag_base"             json:"flag_base"`
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

func GetMsCurrencyIn(c *[]MsCurrency, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_currency.* FROM 
				ms_currency `
	query := query2 + " WHERE ms_currency." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsCurrency(c *MsCurrency, key string) (int, error) {
	query := `SELECT ms_currency.* FROM ms_currency WHERE ms_currency.currency_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
