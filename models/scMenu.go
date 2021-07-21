package models

import (
	"api/db"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ScMenu struct {
	MenuKey           uint64  `db:"menu_key"                  json:"menu_key"`
	MenuParent        *uint64 `db:"menu_parent"               json:"menu_parent"`
	AppModuleKey      uint64  `db:"app_module_key"            json:"app_module_key"`
	MenuCode          string  `db:"menu_code"                 json:"menu_code"`
	MenuName          string  `db:"menu_name"                 json:"menu_name"`
	MenuPage          *string `db:"menu_page"                 json:"menu_page"`
	MenuURL           *string `db:"menu_url"                  json:"menu_url"`
	MenuTypeKey       uint64  `db:"menu_type_key"             json:"menu_type_key"`
	HasEndpoint       uint8   `db:"has_endpoint"              json:"has_endpoint"`
	MenuDesc          *string `db:"menu_desc"                 json:"menu_desc"`
	RecOrder          *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

type ScMenuDetail struct {
	MenuKey    uint64               `json:"menu_key"`
	ModuleName *string              `json:"module_name"`
	MenuName   string               `json:"menu_name"`
	MenuDesc   *string              `json:"menu_desc"`
	ChildMenu  *[]ScMenuDetailChild `json:"child_menu"`
}

type ScMenuDetailChild struct {
	MenuKey   uint64  `json:"menu_key"`
	MenuName  string  `json:"menu_name"`
	MenuDesc  *string `json:"menu_desc"`
	IsChecked bool    `json:"is_checked"`
}

type ScMenuListRoleLogin struct {
	MenuKey       uint64                 `json:"menu_key"`
	AppModuleName string                 `json:"app_module_name"`
	MenuCode      string                 `json:"menu_code"`
	MenuName      string                 `json:"menu_name"`
	MenuPage      *string                `json:"menu_page"`
	MenuURL       *string                `json:"menu_url"`
	MenuTypeKey   uint64                 `json:"menu_type_key"`
	HasEndpoint   uint8                  `json:"has_endpoint"`
	MenuDesc      *string                `json:"menu_desc"`
	IsChecked     *bool                  `json:"is_checked"`
	ChildMenu     *[]ScMenuListRoleLogin `json:"child"`
}

type ListMenuRoleManagement struct {
	MenuKey    uint64  `db:"menu_key"        json:"menu_key"`
	ModuleName *string `db:"module_name"     json:"module_name"`
	MenuParent *uint64 `db:"menu_parent"     json:"menu_parent"`
	MenuName   string  `db:"menu_name"       json:"menu_name"`
	MenuDesc   *string `db:"menu_desc"       json:"menu_desc"`
	Checked    string  `db:"checked"         json:"checked"`
}

type ListMenuRoleUser struct {
	MenuParent *uint64 `db:"menu_parent"     json:"menu_parent"`
	MenuPage   *string `db:"menu_page"       json:"menu_page"`
	MenuURL    *string `db:"menu_url"        json:"menu_url"`
	Icon       *string `db:"icon"            json:"icon"`
}

type ListParentMenuRoleUser struct {
	MenuKey   uint64  `db:"menu_key"           json:"menu_key"`
	ClassName *string `db:"class_name"         json:"class_name"`
	MenuPage  *string `db:"menu_page"          json:"menu_page"`
	Icon      *string `db:"icon"               json:"icon"`
}

type MenuUserRole struct {
	ClassName *string      `json:"_name"`
	Name      *string      `json:"name"`
	Icon      *string      `json:"icon"`
	To        *string      `json:"to,omitempty"`
	Items     *[]MenuChild `json:"items,omitempty"`
}

type MenuChild struct {
	Name *string `json:"name"`
	To   *string `json:"to"`
	Icon *string `json:"icon"`
}

func AdminGetListMenuRole(c *[]ListMenuRoleManagement, roleKey string, isParent bool) (int, error) {
	query := `SELECT 
				menu.menu_key AS menu_key, 
				app.app_module_name AS module_name, 
				menu.menu_parent AS menu_parent, 
				menu.menu_name AS menu_name, 
				menu.menu_desc AS menu_desc, 
				(CASE 
					WHEN ep.ep_auth_key IS NULL THEN '0' 
					ELSE '1' 
				END) AS checked 
			FROM sc_menu AS menu 
			LEFT JOIN (SELECT ee.* FROM sc_endpoint_auth AS ee WHERE ee.rec_status = 1 AND ee.role_key = '` + roleKey + `' GROUP BY ee.menu_key) AS ep ON ep.menu_key = menu.menu_key 
			LEFT JOIN sc_app_module AS app ON app.app_module_key = menu.app_module_key 
			WHERE menu.rec_status = 1 AND menu.app_module_key != 1`

	if isParent {
		query += " AND menu.menu_parent IS NULL"
	} else {
		query += " AND menu.menu_parent IS NOT NULL"
	}

	query += " ORDER BY menu.app_module_key ASC"

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetParentMenuListRoleLogin(c *[]ListParentMenuRoleUser, value []string) (int, error) {
	inQuery := strings.Join(value, ",")
	query := `SELECT 
				menu_key AS menu_key,
				rec_attribute_id2 AS class_name, 
				menu_page AS menu_page, 
				rec_attribute_id1 AS icon 
			 FROM sc_menu WHERE menu_key IN(` + inQuery + `) ORDER BY menu_key ASC`

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetMenuListRoleLogin(c *[]ListMenuRoleUser, roleKey string) (int, error) {
	query := `SELECT 
				m.menu_parent as menu_parent, 
				m.menu_page as menu_page, 
				m.menu_url as menu_url, 
				m.rec_attribute_id1 AS icon 
			FROM sc_endpoint_auth au 
			INNER JOIN sc_menu AS m ON m.menu_key = au.menu_key 
			WHERE au.role_key = ` + roleKey + ` AND au.rec_status = 1 AND m.rec_status = 1 
			GROUP BY au.menu_key ORDER BY m.rec_order`

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

type ListMenuAdmin struct {
	MenuKey       uint64  `db:"menu_key"        json:"menu_key"`
	ParentName    *string `db:"parent_name"     json:"parent_name"`
	MenuCode      *string `db:"menu_code"       json:"menu_code"`
	MenuName      *string `db:"menu_name"       json:"menu_name"`
	MenuPage      *string `db:"menu_page"       json:"menu_page"`
	MenuDesc      *string `db:"menu_desc"       json:"menu_desc"`
	MenuUrl       *string `db:"menu_url"        json:"menu_url"`
	AppModuleName *string `db:"app_module_name" json:"app_module_name"`
	MenuTypeName  *string `db:"menu_type_name"  json:"menu_type_name"`
}

func AdminGetListMenu(c *[]ListMenuAdmin, limit uint64, offset uint64, params map[string]string, searchLike string, nolimit bool) (int, error) {
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
		condition += " (m.menu_key like '%" + searchLike + "%' OR"
		condition += " par.menu_name like '%" + searchLike + "%' OR"
		condition += " m.menu_code like '%" + searchLike + "%' OR"
		condition += " m.menu_name like '%" + searchLike + "%' OR"
		condition += " m.menu_page like '%" + searchLike + "%' OR"
		condition += " m.menu_desc like '%" + searchLike + "%' OR"
		condition += " m.menu_url like '%" + searchLike + "%' OR"
		condition += " mo.app_module_name like '%" + searchLike + "%' OR"
		condition += " t.menu_type_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				dat.menu_key AS menu_key, 
				dat.parent_name AS parent_name, 
				dat.menu_code AS menu_code, 
				dat.menu_name AS menu_name, 
				dat.menu_page AS menu_page, 
				dat.menu_desc AS menu_desc, 
				dat.menu_url AS menu_url,
				dat.app_module_name AS app_module_name, 
				dat.menu_type_name AS menu_type_name   
			FROM
				(SELECT 
					m.menu_key,
					(CASE
						WHEN m.menu_parent IS NULL THEN ""
						ELSE par.menu_name
					END) AS parent_name,
					m.menu_code,
					m.menu_name,
					m.menu_page,
					m.menu_desc,
					m.menu_url,
					mo.app_module_name,
					t.menu_type_name,
					(CASE
						WHEN m.menu_parent IS NULL THEN m.rec_status
						ELSE par.rec_status
					END) AS parent_rec_status   
				FROM sc_menu AS m
				LEFT JOIN sc_menu AS par ON m.menu_parent = par.menu_key
				INNER JOIN sc_app_module AS mo ON mo.app_module_key = m.app_module_key AND mo.rec_status = 1
				INNER JOIN sc_menu_type AS t ON t.menu_type_key = m.menu_type_key AND t.rec_status = 1
				WHERE m.rec_status = 1 ` + condition + `
				) AS dat
			WHERE dat.parent_rec_status = 1`

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY dat." + orderBy
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

func CountAdminGetListMenu(c *CountData, params map[string]string, searchLike string) (int, error) {
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
		condition += " (m.menu_key like '%" + searchLike + "%' OR"
		condition += " par.menu_name like '%" + searchLike + "%' OR"
		condition += " m.menu_code like '%" + searchLike + "%' OR"
		condition += " m.menu_name like '%" + searchLike + "%' OR"
		condition += " m.menu_page like '%" + searchLike + "%' OR"
		condition += " m.menu_desc like '%" + searchLike + "%' OR"
		condition += " m.menu_url like '%" + searchLike + "%' OR"
		condition += " mo.app_module_name like '%" + searchLike + "%' OR"
		condition += " t.menu_type_name like '%" + searchLike + "%')"
	}

	query := `SELECT 
				count(dat.menu_key) AS count_data 
			FROM
				(SELECT 
					m.menu_key,
					(CASE
						WHEN m.menu_parent IS NULL THEN m.rec_status
						ELSE par.rec_status
					END) AS parent_rec_status   
				FROM sc_menu AS m
				LEFT JOIN sc_menu AS par ON m.menu_parent = par.menu_key
				INNER JOIN sc_app_module AS mo ON mo.app_module_key = m.app_module_key AND mo.rec_status = 1
				INNER JOIN sc_menu_type AS t ON t.menu_type_key = m.menu_type_key AND t.rec_status = 1
				WHERE m.rec_status = 1 ` + condition + `
				) AS dat
			WHERE dat.parent_rec_status = 1`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateScMenu(params map[string]string) (int, error) {
	query := "UPDATE sc_menu SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "menu_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE menu_key = " + params["menu_key"]
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

func CreateScMenu(params map[string]string) (int, error) {
	query := "INSERT INTO sc_menu"
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

func CountScMenuValidateUnique(c *CountData, field string, value string, key string) (int, error) {
	query := `SELECT 
				COUNT(menu_key) AS count_data 
			FROM sc_menu
			WHERE rec_status = '1' AND ` + field + ` = '` + value + `'`

	if key != "" {
		query += " AND menu_key != '" + key + "'"
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

func GetScMenu(c *ScMenu, key string) (int, error) {
	query := `SELECT sc_menu.* FROM sc_menu WHERE sc_menu.rec_status = 1 AND sc_menu.menu_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
