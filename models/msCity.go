package models

import (
	"api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MsCityList struct {
	CityKey    uint64 `json:"city_key"`
	ParentKey  uint64 `json:"parent_key"`
	CityCode   string `json:"city_code"`
	CityName   string `json:"city_name"`
	CityLevel  uint64 `json:"city_level"`
	PostalCode string `json:"postal_code"`
}

type MsCity struct {
	CityKey           uint64  `db:"city_key"              json:"city_key"`
	CountryKey        uint64  `db:"country_key"              json:"country_key"`
	ParentKey         *uint64 `db:"parent_key"            json:"parent_key"`
	CityCode          string  `db:"city_code"             json:"city_code"`
	CityName          string  `db:"city_name"             json:"city_name"`
	CityLevel         uint64  `db:"city_level"            json:"city_level"`
	PostalCode        *string `db:"postal_code"           json:"postal_code"`
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

func GetAllMsCity(c *[]MsCity, params map[string]string) (int, error) {
	query := `SELECT
              ms_city.* FROM 
			  ms_city `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_city."+field+" = '"+value+"'")
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

func GetMsCityIn(c *[]MsCity, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_city.* FROM 
				ms_city `
	query := query2 + " WHERE ms_city." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetMsCity(c *MsCity, key string) (int, error) {
	query := `SELECT ms_city.* FROM ms_city WHERE ms_city.rec_status = '1' AND ms_city.city_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func GetMsCityByParent(c *MsCity, key string) (int, error) {
	query := `SELECT * 
			FROM ms_city 
			WHERE city_key = (SELECT parent_key FROM ms_city WHERE city_key = '` + key + `')`
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

type ListCity struct {
	CurrRateKey uint64  `db:"city_key"            json:"city_key"`
	CouName     string  `db:"cou_name"            json:"cou_name"`
	CityParent  *string `db:"city_parent"         json:"city_parent"`
	CityName    string  `db:"city_name"           json:"city_name"`
	CityCode    string  `db:"city_code"           json:"city_code"`
	CityLevel   string  `db:"city_level"          json:"city_level"`
	PostalCode  *string `db:"postal_code"         json:"postal_code"`
}

func AdminGetListCity(c *[]ListCity, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
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
		condition += " (cou.cou_name like '%" + searchLike + "%' OR"
		condition += " par.city_name like '%" + searchLike + "%' OR"
		condition += " c.city_name like '%" + searchLike + "%' OR"
		condition += " c.city_code like '%" + searchLike + "%' OR"
		condition += " cl.lkp_name like '%" + searchLike + "%' OR"
		condition += " c.postal_code like '%" + searchLike + "%')"
	}

	query := `SELECT 
				c.city_key,
				cou.cou_name,
				par.city_name AS city_parent,
				c.city_name,
				c.city_code,
				cl.lkp_name AS city_level,
				c.postal_code 
			FROM ms_city AS c
			INNER JOIN ms_country AS cou ON cou.country_key = c.country_key
			LEFT JOIN ms_city AS par ON par.city_key = c.parent_key AND par.rec_status = '1'
			INNER JOIN gen_lookup AS cl ON cl.lkp_code = c.city_level AND cl.lkp_group_key = '47'
			WHERE c.rec_status = 1 AND c.city_level IN (1,2,3,4)` + condition

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

func CountAdminGetCity(c *CountData, params map[string]string, searchLike string) (int, error) {
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
		condition += " (cou.cou_name like '%" + searchLike + "%' OR"
		condition += " par.city_name like '%" + searchLike + "%' OR"
		condition += " c.city_name like '%" + searchLike + "%' OR"
		condition += " c.city_code like '%" + searchLike + "%' OR"
		condition += " cl.lkp_name like '%" + searchLike + "%' OR"
		condition += " c.postal_code like '%" + searchLike + "%')"
	}

	query := `SELECT
				count(c.city_key) AS count_data 
			FROM ms_city AS c
			INNER JOIN ms_country AS cou ON cou.country_key = c.country_key
			LEFT JOIN ms_city AS par ON par.city_key = c.parent_key AND par.rec_status = '1'
			INNER JOIN gen_lookup AS cl ON cl.lkp_code = c.city_level AND cl.lkp_group_key = '47'
			WHERE c.rec_status = 1 AND c.city_level IN (1,2,3,4)` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateMsCity(params map[string]string) (int, error) {
	query := "INSERT INTO ms_city"
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

func UpdateMsCity(params map[string]string) (int, error) {
	query := "UPDATE ms_city SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "city_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE city_key = " + params["city_key"]
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

func CountMsCityValidateUnique(c *CountData, field string, value string, key string) (int, error) {
	query := `SELECT 
				COUNT(city_key) AS count_data 
			FROM ms_city
			WHERE rec_status = '1' AND ` + field + ` = '` + value + `'`

	if key != "" {
		query += " AND city_key != '" + key + "'"
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

type ListParent struct {
	CurrRateKey uint64 `db:"parent_key"            json:"parent_key"`
	CityName    string `db:"city_name"             json:"city_name"`
}

func AdminGetListParent(c *[]ListParent) (int, error) {
	query := `SELECT
				city_key AS parent_key,
				(CASE
					WHEN city_level = "1" THEN CONCAT("- ", city_name)
					WHEN city_level = "2" THEN CONCAT("-- ", city_name)
					WHEN city_level = "3" THEN CONCAT("--- ", city_name)
				END) AS city_name 
			FROM ms_city 
			WHERE city_level IN (1,2,3) AND rec_status = 1
			ORDER BY city_level ASC`

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
