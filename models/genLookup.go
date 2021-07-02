package models

import (
	"api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type GenLookupInfo struct {
	Name  string `json:"name"`
	Value uint64 `json:"value"`
}

type GenLookup struct {
	LookupKey         uint64  `db:"lookup_key"              json:"lookup_key"`
	LkpGroupKey       uint64  `db:"lkp_group_key"           json:"lkp_group_key"`
	LkpCode           *string `db:"lkp_code"                json:"lkp_code"`
	LkpName           *string `db:"lkp_name"                json:"lkp_name"`
	LkpDesc           *string `db:"lkp_desc"                json:"lkp_desc"`
	LkpVal1           *string `db:"lkp_val1"                json:"lkp_val1"`
	LkpVal2           *string `db:"lkp_val2"                json:"lkp_val2"`
	LkpVal3           *string `db:"lkp_val3"                json:"lkp_val3"`
	LkpText1          *string `db:"lkp_text1"               json:"lkp_text1"`
	LkpText2          *string `db:"lkp_text2"               json:"lkp_text2"`
	LkpText3          *string `db:"lkp_text3"               json:"lkp_text3"`
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

type MetodePerhitungan struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func GetAllGenLookup(c *[]GenLookup, params map[string]string) (int, error) {
	query := `SELECT
              gen_lookup.* FROM 
			  gen_lookup `
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "gen_lookup."+field+" = '"+value+"'")
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

func GetGenLookupIn(c *[]GenLookup, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				gen_lookup.* FROM 
				gen_lookup `
	query := query2 + " WHERE gen_lookup." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetGenLookup(c *GenLookup, key string) (int, error) {
	query := `SELECT gen_lookup.* FROM gen_lookup WHERE gen_lookup.lookup_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

type ListLookup struct {
	LookupKey    uint64  `db:"lookup_key"        json:"lookup_key"`
	LkpGroupCode string  `db:"lkp_group_code"    json:"lkp_group_code"`
	LkpCode      string  `db:"lkp_code"          json:"lkp_code"`
	LkpName      string  `db:"lkp_name"          json:"lkp_name"`
	LkpDesc      *string `db:"lkp_desc"          json:"lkp_desc"`
}

func AdminGetLookup(c *[]ListLookup, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
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
		condition += " (gg.lkp_group_code like '%" + searchLike + "%' OR"
		condition += " g.lkp_code like '%" + searchLike + "%' OR"
		condition += " g.lkp_name like '%" + searchLike + "%' OR"
		condition += " g.lkp_desc  like '%" + searchLike + "%')"
	}

	query := `SELECT 
				g.lookup_key,
				gg.lkp_group_code,
				g.lkp_code,
				g.lkp_name,
				g.lkp_desc 
			FROM gen_lookup AS g
			INNER JOIN gen_lookup_group AS gg ON gg.lkp_group_key = g.lkp_group_key
			WHERE g.rec_status = 1 ` + condition

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

func CountAdminGetLookup(c *CountData, params map[string]string, searchLike string) (int, error) {
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
		condition += " (gg.lkp_group_code like '%" + searchLike + "%' OR"
		condition += " g.lkp_code like '%" + searchLike + "%' OR"
		condition += " g.lkp_name like '%" + searchLike + "%' OR"
		condition += " g.lkp_desc  like '%" + searchLike + "%')"
	}

	query := `SELECT 
				count(g.lookup_key) AS count_data 
			FROM gen_lookup AS g
			INNER JOIN gen_lookup_group AS gg ON gg.lkp_group_key = g.lkp_group_key
			WHERE g.rec_status = 1 ` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateLookup(params map[string]string) (int, error) {
	query := "UPDATE gen_lookup SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "lookup_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE lookup_key = " + params["lookup_key"]
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

func CreateLookup(params map[string]string) (int, error) {
	query := "INSERT INTO gen_lookup"
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

func GetLookup(c *GenLookup, key string) (int, error) {
	query := `SELECT * FROM gen_lookup WHERE rec_status = 1 AND lookup_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
