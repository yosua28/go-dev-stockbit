package models

import (
	"api/db"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsCountryList struct {
	CountryKey uint64 `json:"country_key"`
	CouCode    string `json:"cou_code"`
	CouName    string `json:"cou_name"`
}

type MsCountry struct {
	CountryKey        uint64  `db:"country_key"           json:"country_key"`
	CouCode           string  `db:"cou_code"              json:"cou_code"`
	CouName           string  `db:"cou_name"              json:"cou_name"`
	ShortName         *string `db:"short_name"            json:"short_name"`
	FlagBase          uint8   `db:"flag_base"             json:"flag_base"`
	CurrencyKey       *uint64 `db:"currency_key"          json:"currency_key"`
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

func GetAllMsCountry(c *[]MsCountry, params map[string]string) (int, error) {
	query := `SELECT
              ms_country.* FROM 
			  ms_country `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_country."+field+" = '"+value+"'")
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
		if orderBy == "cou_name" {
			condition += " ORDER BY FIELD(cou_name, 'Indonesia') DESC, cou_name ASC "
		} else {
			condition += " ORDER BY " + orderBy
			if orderType, present = params["orderType"]; present == true {
				condition += " " + orderType
			}
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

func GetMsCountryIn(c *[]MsCountry, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_country.* FROM 
				ms_country `
	query := query2 + " WHERE ms_country.rec_status = 1 AND ms_country." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsCountry(c *MsCountry, key string) (int, error) {
	query := `SELECT ms_country.* FROM ms_country WHERE ms_country.country_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
