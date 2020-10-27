package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"
)

type TrBalance struct {
	BalanceKey                uint64    `db:"balance_key"               json:"balance_key"`
	AcaKey                    uint64    `db:"aca_key"                   json:"aca_key"`
	TcKey                     uint64    `db:"tc_key"                    json:"tc_key"`
	BalanceDate               string    `db:"balance_date"              json:"balance_date"`
	BalanceUnit               float32   `db:"balance_unit"              json:"balance_unit"`
	TcKeyRed                  *uint64   `db:"tc_key_red"                json:"tc_key_red"`
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

func GetLastBalanceIn(c *[]TrBalance, acaKey []string,) (int, error) {
	inQuery := strings.Join(acaKey, ",")
	query2 := `SELECT 
				t1.balance_key, 
				t1.aca_key, 
				t1.tc_key, 
				t1.balance_date, 
				t1.balance_unit, 
				t1.tc_key_red FROM
				   tr_balance t1 JOIN 
				   (SELECT MAX(balance_key) balance_key FROM tr_balance GROUP BY tc_key) t2
			   ON t1.balance_key = t2.balance_key`
	query := query2 + " WHERE t1.aca_key IN(" + inQuery + ") GROUP BY tc_key"
	
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}