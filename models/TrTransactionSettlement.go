package models

type TrTransactionSettlement struct{
	SettlementKey             uint64    `db:"settlement_key"            json:"settlement_key"`
	TransactionKey           *uint64    `db:"transaction_key"           json:"transaction_key"`
	SettlePurposed            string    `db:"settle_purposed"           json:"settle_purposed"`
	SettleDate                string    `db:"settle_date"               json:"settle_date"`
	SettleNominal             float32   `db:"settle_nominal"            json:"settle_nominal"`
	SettleStatus              uint64    `db:"settle_status"             json:"settle_status"`
	SettleRealizedDate        string    `db:"settle_realize_date"       json:"settle_realize_date"`
	SettleRemarks            *string    `db:"settle_remarks"            json:"settle_remarks"`
	SettleReference          *string    `db:"settle_reference"          json:"settle_reference"`
	SourceBankAccountKey     *uint64    `db:"source_bank_account_key"   json:"source_bank_account_key"`
	TargetBankAccountKey      uint64    `db:"target_bank_account_key"   json:"target_bank_account_key"`
	SettleNotes              *string    `db:"settle_notes"              json:"settle_notes"`
	SettleAcknowledgement    *string    `db:"settle_acknowledgement"    json:"settle_acknowledgement"`
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