package models

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
