package models

type TrAccountAgent struct {
	AcaKey                    uint64    `db:"aca_key"                   json:"aca_key"`
	AccKey                    uint64    `db:"acc_key"                   json:"acc_key"`
	EffDate                   string    `db:"eff_date"                  json:"eff_date"`
	BranchKey                 *uint64   `db:"branch_key"                json:"branch_key"`
	AgentKey                  *uint64   `db:"agent_key"                 json:"agent_key"`
	RecOrder                  *uint64   `db:"rec_order"                 json:"rec_order"`
	RecStatus                 uint8     `db:"rec_status"                json:"rec_status"`
	RecCreatedDate            *string   `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy              *string   `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate           *string   `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy             *string   `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1                 *string   `db:"rec_image1"                json:"rec_image1"`
	RecImage2                 *string   `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus         *uint8    `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage          *uint64   `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate           *string   `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy             *string   `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate            *string   `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy              *string   `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1           *string   `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2           *string   `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3           *string   `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}