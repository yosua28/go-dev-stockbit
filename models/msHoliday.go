package models

import (
	"api/db"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MsHoliday struct {
	HolidayKey        uint64  `db:"holiday_key"           json:"holiday_key"`
	StickMarketKey    uint64  `db:"stock_market"          json:"stock_market"`
	HolidayDate       string  `db:"holiday_date"          json:"holiday_date"`
	HolidayName       *string `db:"holiday_name"          json:"holiday_name"`
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

func GetAllMsHoliday(c *[]MsHoliday, params map[string]string) (int, error) {
	query := `SELECT
              ms_holiday.* FROM 
			  ms_holiday`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_holiday."+field+" = '"+value+"'")
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

func GetMsHoliday(c *MsHoliday, key string) (int, error) {
	query := `SELECT ms_holiday.* FROM ms_holiday WHERE ms_holiday.rec_status = '1' AND ms_holiday.holiday_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

type ListHoliday struct {
	HolidayKey  uint64  `db:"holiday_key"           json:"holiday_key"`
	StockMarket *string `db:"stock_market"          json:"stock_market"`
	HolidayDate string  `db:"holiday_date"          json:"holiday_date"`
	HolidayName *string `db:"holiday_name"          json:"holiday_name"`
}

func AdminGetListHoliday(c *[]ListHoliday, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
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
		condition += " (stock.lkp_name like '%" + searchLike + "%' OR"
		condition += " DATE_FORMAT(h.holiday_date, '%d %M %Y') like '%" + searchLike + "%' OR"
		condition += " h.holiday_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				h.holiday_key,
				stock.lkp_name AS stock_market,
				DATE_FORMAT(h.holiday_date, '%d %M %Y') AS holiday_date,
				h.holiday_name 
			FROM ms_holiday AS h
			LEFT JOIN gen_lookup AS stock ON stock.lookup_key = h.stock_market
			WHERE h.rec_status = 1` + condition

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

func CountAdminGetHoliday(c *CountData, params map[string]string, searchLike string) (int, error) {
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
		condition += " (stock.lkp_name like '%" + searchLike + "%' OR"
		condition += " DATE_FORMAT(h.holiday_date, '%d %M %Y') like '%" + searchLike + "%' OR"
		condition += " h.holiday_name like '%" + searchLike + "%')"
	}

	query := `SELECT
				count(h.holiday_key) AS count_data 
			FROM ms_holiday AS h
			LEFT JOIN gen_lookup AS stock ON stock.lookup_key = h.stock_market
			WHERE h.rec_status = 1` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateMsHoliday(params map[string]string) (int, error) {
	query := "INSERT INTO ms_holiday"
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

func UpdateMsHoliday(params map[string]string) (int, error) {
	query := "UPDATE ms_holiday SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "holiday_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE holiday_key = " + params["holiday_key"]
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

func CountMsHolidayValidateUnique(c *CountData, field string, value string, key string) (int, error) {
	query := `SELECT 
				COUNT(holiday_key) AS count_data 
			FROM ms_holiday
			WHERE rec_status = '1' AND ` + field + ` = '` + value + `'`

	if key != "" {
		query += " AND holiday_key != '" + key + "'"
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
