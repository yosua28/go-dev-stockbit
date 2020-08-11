package models

type TrAccount struct {
	AccKey                   uint64    `db:"acc_key"                    json:"acc_key"`
	ProductKey               uint64    `db:"product_key"                json:"product_key"`
	CustomerKey              uint64    `db:"customer_key"               json:"customer_key"`
	AccountName              string    `db:"account_name"               json:"account_name"`
	AccountNo                string    `db:"account_no"                 json:"account_no"`
	IfuaNo                   *string   `db:"ifua_no"                    json:"ifua_no"`
	IfuaName                 *string   `db:"ifua_name"                  json:"ifua_name"`
	AccStatus                uint8     `db:"acc_status"                 json:"acc_status"`
	AccSuspendFlag           uint8     `db:"acc_suspend_flag"           json:"acc_suspend_flag"`
	AccSuspendModifiedFlag   uint8     `db:"acc_suspend_modified_flag"  json:"acc_suspend_modified_flag"`
	AccSuspendReason         *uint8    `db:"acc_suspend_reason"         json:"acc_suspend_reason"`
	RecOrder                 *uint64   `db:"rec_order"                  json:"rec_order"`
	RecStatus                uint8     `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate           *string   `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy             *string   `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate          *string   `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy            *string   `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1                *string   `db:"rec_image1"                 json:"rec_image1"`
	RecImage2                *string   `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus        *uint8    `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage         *uint64   `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate          *string   `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy            *string   `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate           *string   `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy             *string   `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1          *string   `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2          *string   `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3          *string   `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
}