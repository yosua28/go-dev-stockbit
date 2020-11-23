package models

type ScEndpoint struct {
	EndpointKey         uint64  `db:"endpoint_key"              json:"endpoint_key"`
	MenuKey             *uint64 `db:"menu_key"                  json:"menu_key"`
	EndpointCategoryKey uint64  `db:"endpoint_category_key"     json:"endpoint_category_key"`
	EndpointCode        string  `db:"endpoint_code"             json:"endpoint_code"`
	EndpointName        *string `db:"endpoint_name"             json:"endpoint_name"`
	EndpointVerb        string  `db:"endpoint_verb"             json:"endpoint_verb"`
	EndpointURI         *string `db:"endpoint_uri"              json:"endpoint_uri"`
	EndpointRoute       *string `db:"endpoint_route"            json:"endpoint_route"`
	EndpointController  *string `db:"endpoint_controller"       json:"endpoint_controller"`
	EndpointVersion     uint64  `db:"endpoint_version"          json:"endpoint_version"`
	PrivilegesKey       string  `db:"privileges_key"            json:"privileges_key"`
	EndpointDesc        *string `db:"endpoint_desc"             json:"endpoint_desc"`
	RecOrder            *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus           uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate      *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy        *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate     *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy       *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1           *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2           *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus   *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage    *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate     *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy       *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate      *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy        *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1     *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2     *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3     *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}
