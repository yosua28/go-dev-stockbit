package models

import (
	"api/db"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type TrPromo struct {
	PromoKey              uint64          `db:"promo_key"               json:"promo_key"`
	PromoCode             *string         `db:"promo_code"              json:"promo_code"`
	PromoTitle            *string         `db:"promo_title"             json:"promo_title"`
	PromoCategory         uint64          `db:"promo_category"          json:"promo_category"`
	PromoNominal          decimal.Decimal `db:"promo_nominal"           json:"promo_nominal"`
	PromoMaxNominal       decimal.Decimal `db:"promo_max_nominal"       json:"promo_max_nominal"`
	PromoValuesType       uint64          `db:"promo_values_type"       json:"promo_values_type"`
	PromoMaxuser          uint64          `db:"promo_maxuser"           json:"promo_maxuser"`
	PromoStayPeriode      uint64          `db:"promo_stay_periode"      json:"promo_stay_periode"`
	PromoFlagUniqUser     uint8           `db:"promo_flag_uniq_user"    json:"promo_flag_uniq_user"`
	PromoValidDate1       string          `db:"promo_valid_date1"       json:"promo_valid_date1"`
	PromoValidDate2       string          `db:"promo_valid_date2"       json:"promo_valid_date2"`
	PromoNotifStart       string          `db:"promo_notif_start"       json:"promo_notif_start"`
	PromoNotifEnd         string          `db:"promo_notif_end"         json:"promo_notif_end"`
	PromoNotifType        uint64          `db:"promo_notif_type"        json:"promo_notif_type"`
	PromoNotifDescription string          `db:"promo_description"       json:"promo_description"`
	PromoTnc              string          `db:"promo_tnc"               json:"promo_tnc"`
	RecOrder              *uint64         `db:"rec_order"               json:"rec_order"`
	RecStatus             uint8           `db:"rec_status"              json:"rec_status"`
	RecCreatedDate        *string         `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy          *string         `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate       *string         `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy         *string         `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1             *string         `db:"rec_image1"              json:"rec_image1"`
	RecImage2             *string         `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus     *uint8          `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage      *uint64         `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate       *string         `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy         *string         `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate        *string         `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy          *string         `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1       *string         `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2       *string         `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3       *string         `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

type TrPromoData struct {
	PromoKey        uint64          `db:"promo_key"               json:"promo_key"`
	Category        string          `db:"category"                json:"category"`
	PromoCode       *string         `db:"promo_code"              json:"promo_code"`
	PromoTitle      *string         `db:"promo_title"             json:"promo_title"`
	PromoNominal    decimal.Decimal `db:"promo_nominal"           json:"promo_nominal"`
	PromoMaxNominal decimal.Decimal `db:"promo_max_nominal"       json:"promo_max_nominal"`
	ValueType       string          `db:"value_type"              json:"value_type"`
	PromoMaxuser    uint64          `db:"promo_maxuser"           json:"promo_maxuser"`
	StartValid      string          `db:"start_valid"             json:"start_valid"`
	EndValid        string          `db:"end_valid"               json:"end_valid"`
	PromoNotifStart string          `db:"promo_notif_start"       json:"promo_notif_start"`
	PromoNotifEnd   string          `db:"promo_notif_end"         json:"promo_notif_end"`
}

type TrPromoDetail struct {
	PromoKey              uint64               `json:"promo_key"`
	PromoCode             *string              `json:"promo_code"`
	PromoTitle            *string              `json:"promo_title"`
	PromoCategory         LookupTrans          `json:"promo_category"`
	PromoNominal          decimal.Decimal      `json:"promo_nominal"`
	PromoMaxNominal       decimal.Decimal      `json:"promo_max_nominal"`
	PromoValuesType       LookupTrans          `json:"promo_values_type"`
	PromoMaxuser          uint64               `json:"promo_maxuser"`
	PromoStayPeriode      uint64               `json:"promo_stay_periode"`
	PromoFlagUniqUser     uint8                `json:"promo_flag_uniq_user"`
	PromoValidDate1       string               `json:"promo_valid_date1"`
	PromoValidDate2       string               `json:"promo_valid_date2"`
	PromoNotifStart       string               `json:"promo_notif_start"`
	PromoNotifEnd         string               `json:"promo_notif_end"`
	PromoNotifType        LookupTrans          `json:"promo_notif_type"`
	PromoNotifDescription string               `json:"promo_description"`
	PromoTnc              string               `json:"promo_tnc"`
	Image                 *string              `json:"image"`
	PromoProduct          []TrPromoProductData `json:"promo_product"`
}

type TrPromoCron struct {
	PromoKey         uint64  `db:"promo_key"               json:"promo_key"`
	PromoCode        *string `db:"promo_code"              json:"promo_code"`
	PromoTitle       *string `db:"promo_title"             json:"promo_title"`
	PromoDescription string  `db:"promo_description"       json:"promo_description"`
}

type CheckPromo struct {
	PromoEnabled bool   `json:"promo_enabled"`
	Message      string `json:"message"`
}

func CreateTrPromo(params map[string]string) (int, error, string) {
	query := "INSERT INTO tr_promo"
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
	log.Println(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func AdminGetAllTrPromo(c *[]TrPromoData, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT 
				p.promo_key as promo_key,
				c.lkp_name AS category,
				p.promo_code as promo_code,
				p.promo_title as promo_title,
				p.promo_nominal as promo_nominal,
				p.promo_max_nominal as promo_max_nominal,
				v.lkp_name AS value_type,
				p.promo_maxuser as promo_maxuser,
				DATE_FORMAT(p.promo_valid_date1, '%d %M %Y') AS start_valid,
				DATE_FORMAT(p.promo_valid_date2, '%d %M %Y') AS end_valid,
				DATE_FORMAT(p.promo_notif_start, '%d %M %Y') as promo_notif_start,
				DATE_FORMAT(p.promo_notif_end, '%d %M %Y') as promo_notif_end 
			FROM tr_promo AS p 
			LEFT JOIN gen_lookup AS c ON c.lookup_key = p.promo_category 
			LEFT JOIN gen_lookup AS v ON v.lookup_key = p.promo_values_type 
			WHERE p.rec_status = 1`
	var present bool
	var whereClause []string
	var condition string
	dateFrom := ""
	dateTo := ""

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "start_valid" || field == "end_valid") {
			whereClause = append(whereClause, "p."+field+" = '"+value+"'")
		}
		if field == "start_valid" {
			dateFrom = value
		}
		if field == "end_valid" {
			dateTo = value
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

	if (dateFrom != "") && (dateTo != "") {
		query += " AND p.promo_valid_date2 >= '" + dateFrom + "'"
	}

	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY p." + orderBy
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

func AdminGetCountTrPromo(c *CountData, params map[string]string) (int, error) {
	query := `SELECT
              count(p.promo_key) as count_data
			  FROM tr_promo as p
			  WHERE p.rec_status = 1`

	var whereClause []string
	var condition string

	dateFrom := ""
	dateTo := ""

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType" || field == "start_valid" || field == "end_valid") {
			whereClause = append(whereClause, "p."+field+" = '"+value+"'")
		}
		if field == "start_valid" {
			dateFrom = value
		}
		if field == "end_valid" {
			dateTo = value
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

	if (dateFrom != "") && (dateTo != "") {
		query += " AND p.promo_valid_date2 >= '" + dateFrom + "'"
	}

	query += condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateTrPromo(params map[string]string) (int, error) {
	query := "UPDATE tr_promo SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "promo_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE promo_key = " + params["promo_key"]
	log.Println(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}
	var ret sql.Result
	ret, err = tx.Exec(query)
	row, _ := ret.RowsAffected()
	tx.Commit()
	if row > 0 {
	} else {
		return http.StatusNotFound, err
	}
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func GetTrPromo(c *TrPromo, field string, value string) (int, error) {
	query := `SELECT tr_promo.* FROM tr_promo WHERE tr_promo.rec_status = 1 AND tr_promo.` + field + ` = '` + value + `'`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetTrPromoValidasiDuplikat(c *TrPromo, field string, value string, promoKeyNot string) (int, error) {
	query := `SELECT tr_promo.* FROM tr_promo WHERE tr_promo.rec_status = 1 AND tr_promo.` + field + ` = '` + value + `' 
	AND tr_promo.promo_key != '` + promoKeyNot + `'`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func AdminGetAllPromoOnce(c *[]TrPromoCron) (int, error) {
	query := `SELECT 
				promo_key, promo_code, promo_title, promo_description 
			FROM tr_promo
			WHERE rec_status = 1 AND 
			promo_notif_type = '310' AND 
			DATE(NOW()) = DATE(promo_notif_start)`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetAllPromoOnceBeginOnceBefore(c *[]TrPromoCron) (int, error) {
	query := `SELECT 
				promo_key, promo_code, promo_title, promo_description  
			FROM tr_promo
			WHERE rec_status = 1 AND 
			promo_notif_type = '311' AND 
			((DATE(NOW()) = DATE(promo_notif_start)) OR (DATE(NOW()) = DATE(promo_notif_end)))`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetAllPromoOnceAday(c *[]TrPromoCron) (int, error) {
	query := `SELECT 
				promo_key, promo_code, promo_title, promo_description 
			FROM tr_promo
			WHERE rec_status = 1 AND 
			promo_notif_type = '312' AND 
			(DATE(NOW()) <= DATE(promo_notif_end)) AND (DATE(NOW()) >= DATE(promo_notif_start))`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetTrPromoProductActive(c *TrPromo, promoCode string, productKey string) (int, error) {
	query := `SELECT 
				t.* 
			FROM tr_promo AS t 
			INNER JOIN tr_promo_product AS pp ON pp.promo_key = t.promo_key 
			WHERE t.rec_status = 1 AND pp.rec_status = 1 AND t.promo_valid_date1 <= NOW() 
			AND t.promo_valid_date2 >= NOW() AND pp.product_key = '` + productKey + `' AND t.promo_code = '` + promoCode + `'`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
