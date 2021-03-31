package models

import (
	"api/db"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TrPromoProduct struct {
	PromoProductKey   uint64  `db:"promo_product_key"       json:"promo_product_key"`
	PromoKey          uint64  `db:"promo_key"               json:"promo_key"`
	ProductKey        uint64  `db:"product_key"             json:"product_key"`
	FlagAllowed       uint8   `db:"flag_allowed"            json:"flag_allowed"`
	RecOrder          *uint64 `db:"rec_order"               json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"              json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"              json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}

type TrPromoProductData struct {
	PromoProductKey uint64 `db:"promo_product_key"       json:"promo_product_key"`
	PromoKey        uint64 `db:"promo_key"               json:"promo_key"`
	ProductKey      uint64 `db:"product_key"             json:"product_key"`
	ProductName     string `db:"product_name"            json:"product_name"`
}

func CreateTrPromoProduct(params map[string]string) (int, error) {
	query := "INSERT INTO tr_promo_product"
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
		return http.StatusBadGateway, err
	}
	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func UpdateTrPromoProductByField(params map[string]string, field string, value string) (int, error) {
	query := "UPDATE tr_promo_product SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE " + field + " = '" + value + "'"
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

func AdminGetPromoProductInNotIn(c *[]TrPromoProduct, valueIn []string, fieldIn string, promoKey string, opr string) (int, error) {
	query := `SELECT
              *
			  FROM tr_promo_product
			  WHERE tr_promo_product.rec_status = 1 AND tr_promo_product.promo_key = '` + promoKey + `'`

	var condition string

	if len(valueIn) > 0 {
		inQuery := strings.Join(valueIn, ",")
		condition += " AND tr_promo_product." + fieldIn + " " + opr + "(" + inQuery + ")"
	}

	query += condition

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateTrPromoProductByFieldIn(params map[string]string, field string, valueIn []string) (int, error) {
	inQuery := strings.Join(valueIn, ",")
	query := "UPDATE tr_promo_product SET "
	// Get params
	i := 0
	for key, value := range params {
		query += key + " = '" + value + "'"

		if (len(params) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE " + field + " IN(" + inQuery + ")"
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

func AdminGetPromoProductByPromoKey(c *[]TrPromoProductData, promoKey string) (int, error) {
	query := `SELECT 
				a.promo_product_key AS promo_product_key,
				a.promo_key AS promo_key,
				a.product_key AS product_key,
				b.product_name_alt AS product_name 
			FROM tr_promo_product AS a
			INNER JOIN ms_product AS b ON a.product_key = b.product_key
			INNER JOIN tr_promo AS c ON c.promo_key = a.promo_key
			WHERE a.rec_status = 1 AND b.rec_status = 1 AND c.rec_status = 1 
			AND a.promo_key = '` + promoKey + `'`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateMultiplePromoProduct(params []interface{}) (int, error) {

	q := `INSERT INTO tr_promo_product (
		promo_key, 
		product_key,
		flag_allowed,
		rec_status,
		rec_created_date,
		rec_created_by) VALUES `

	for i := 0; i < len(params); i++ {
		q += "(?)"
		if i < (len(params) - 1) {
			q += ","
		}
	}
	query, args, err := sqlx.In(q, params...)
	if err != nil {
		return http.StatusBadGateway, err
	}

	query = db.Db.Rebind(query)
	_, err = db.Db.Query(query, args...)
	if err != nil {
		log.Println(err.Error())
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
}
