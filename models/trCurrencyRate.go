package models

import (
	"api/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	log "github.com/sirupsen/logrus"
)

type TrCurrencyRate struct {
	CurrRateKey       uint64          `db:"curr_rate_key"             json:"curr_rate_key"`
	RateDate          string          `db:"rate_date"                 json:"rate_date"`
	RateType          uint64          `db:"rate_type"                 json:"rate_type"`
	RateValue         decimal.Decimal `db:"rate_value"                json:"rate_value"`
	CurrencyKey       uint64          `db:"currency_key"              json:"currency_key"`
	RecOrder          *uint64         `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8           `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string         `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string         `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string         `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string         `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string         `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string         `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8          `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64         `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string         `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string         `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string         `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string         `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string         `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string         `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string         `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func GetTrCurrencyRate(c *TrCurrencyRate, key string) (int, error) {
	query := `SELECT tr_currency_rate.* FROM tr_currency_rate WHERE tr_currency_rate.rec_status = 1 AND tr_currency_rate.curr_rate_key = ` + key
	log.Info(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Info(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
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

type ListCurrencyRate struct {
	CurrRateKey  uint64          `db:"curr_rate_key"            json:"curr_rate_key"`
	RateDate     string          `db:"rate_date"                json:"rate_date"`
	RateType     *string         `db:"rate_type"                json:"rate_type"`
	CurrencyCode *string         `db:"currency_code"            json:"currency_code"`
	CurrencyName *string         `db:"currency_name"            json:"currency_name"`
	RateValue    decimal.Decimal `db:"rate_value"               json:"rate_value"`
}

func AdminGetListCurrencyRate(c *[]ListCurrencyRate, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
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

	if searchLike != "" {
		condition += " AND"
		condition += " (DATE_FORMAT(cr.rate_date, '%d %M %Y') like '%" + searchLike + "%' OR"
		condition += " ty.lkp_name like '%" + searchLike + "%' OR"
		condition += " c.code like '%" + searchLike + "%' OR"
		condition += " c.name like '%" + searchLike + "%' OR"
		condition += " cr.rate_value like '%" + searchLike + "%')"
	}

	query := `SELECT 
				cr.curr_rate_key,
				DATE_FORMAT(cr.rate_date, '%d %M %Y') AS rate_date,
				ty.lkp_name AS rate_type,
				c.code AS currency_code,
				c.name AS currency_name,
				cr.rate_value AS rate_value  
			FROM tr_currency_rate AS cr
			LEFT JOIN gen_lookup AS ty ON ty.lookup_key = cr.rate_type
			LEFT JOIN ms_currency AS c ON c.currency_key = cr.currency_key
			WHERE cr.rec_status = 1` + condition

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	if !nolimit {
		limitOffset += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			limitOffset += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += orderCondition + limitOffset

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CountAdminGetCurrencyRate(c *CountData, params map[string]string, searchLike string) (int, error) {
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
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

	if searchLike != "" {
		condition += " AND"
		condition += " (DATE_FORMAT(cr.rate_date, '%d %M %Y') like '%" + searchLike + "%' OR"
		condition += " ty.lkp_name like '%" + searchLike + "%' OR"
		condition += " c.code like '%" + searchLike + "%' OR"
		condition += " c.name like '%" + searchLike + "%' OR"
		condition += " cr.rate_value like '%" + searchLike + "%')"
	}

	query := `SELECT
				count(cr.curr_rate_key) AS count_data 
			FROM tr_currency_rate AS cr
			LEFT JOIN gen_lookup AS ty ON ty.lookup_key = cr.rate_type
			LEFT JOIN ms_currency AS c ON c.currency_key = cr.currency_key
			WHERE cr.rec_status = 1` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateTrCurrenctyRate(params map[string]string) (int, error) {
	query := "INSERT INTO tr_currency_rate"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func UpdateTrCurrenctyRate(params map[string]string) (int, error) {
	query := "UPDATE tr_currency_rate SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "curr_rate_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE curr_rate_key = " + params["curr_rate_key"]
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

func CountTrCurrencyRateValidateUniqueDateRateCurrency(c *CountData, date string, rate string, currencyKey string, key string) (int, error) {
	query := `SELECT 
				COUNT(curr_rate_key) AS count_data 
			FROM tr_currency_rate
			WHERE rec_status = 1 AND rate_date = '` + date + `' AND rate_type = '` + rate + `' AND currency_key = '` + currencyKey + `'`

	if key != "" {
		query += " AND curr_rate_key != '" + key + "'"
	}

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
