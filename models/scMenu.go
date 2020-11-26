package models

import (
	"api/db"
	"net/http"

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
