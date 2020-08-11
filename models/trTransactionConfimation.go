package models

type TrTransactionConfirmation struct {
	TcKey                     uint64    `db:"tc_key"                    json:"tc_key"`
	ConfirmDate               string    `db:"confirm_date"              json:"confirm_date"`
	TransactionKey            uint64    `db:"transaction_key"           json:"transaction_key"`
	ConfirmedAmount           float32   `db:"confirmed_amount"          json:"confirmed_amount"`
	ConfirmedUnit             float32   `db:"confirmed_unit"            json:"confirmed_unit"`
	ConfirmResult             uint64    `db:"confirm_result"            json:"confirm_result"`
	ConfirmedAmountDiff       float32   `db:"confirmed_amount_diff"     json:"confirmed_amount_diff"`
	ConfirmedUnitDiff         float32   `db:"confirmed_unit_diff"       json:"confirmed_unit_diff"`
	ConfirmedRemarks         *string    `db:"confirmed_remarks"         json:"confirmed_remarks"`
	ConfirmedReferences      *string    `db:"confirmed_references"      json:"confirmed_references"`
	RecOrder                 *uint64    `db:"rec_order"                 json:"rec_order"`
	RecStatus                 uint8     `db:"rec_status"                json:"rec_status"`
	RecCreatedDate           *string    `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy             *string    `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate          *string    `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy            *string    `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1                *string    `db:"rec_image1"                json:"rec_image1"`
	RecImage2                *string    `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus        *uint8     `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage         *uint64    `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate          *string    `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy            *string    `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate           *string    `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy             *string    `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1          *string    `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2          *string    `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3          *string    `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}