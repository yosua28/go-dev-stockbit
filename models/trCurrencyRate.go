package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type TrCurrencyRate struct {
	CurrRateKey               uint64    `db:"curr_rate_key"             json:"curr_rate_key"`
	RateDate                  string    `db:"rate_date"                 json:"rate_date"`
	RateType                  uint64    `db:"rate_type"                 json:"rate_type"`
	RateValue                 decimal.Decimal   `db:"rate_value"                json:"rate_value"`
	CurrencyKey               uint64    `db:"currency_key"              json:"currency_key"`
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

func GetLastCurrencyIn(c *[]TrCurrencyRate, key []string) (int, error) {
	inQuery := strings.Join(key, ",")
	query2 := `SELECT t1.curr_rate_key, t1.rate_value, t1.currency_key FROM
			   tr_currency_rate t1 JOIN (SELECT MAX(curr_rate_key) curr_rate_key FROM tr_currency_rate GROUP BY currency_key) t2
			   ON t1.curr_rate_key = t2.curr_rate_key`
	query := query2 + " WHERE t1.currency_key IN(" + inQuery + ") GROUP BY currency_key"

	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}