package models

type MsCustodianBank struct {
	CustodianKey         uint64     `db:"custodian_key"         json:"custodian_key"`
	CustodianCode        string     `db:"custodian_code"        json:"custodian_code"`
	CustodianShortName   string     `db:"custodian_short_name"  json:"custodian_short_name"`
	CustodianFullName    *string    `db:"custodian_full_name"   json:"custodian_full_name"`
	BiMemberCode         *string    `db:"bi_member_code"        json:"bi_member_code"`
	SwiftCode            *string    `db:"swift_code"            json:"swift_code"`
	FlagLocal            uint8      `db:"flag_local"            json:"flag_local"`
	FlagGoverment        uint8      `db:"flag_government"       json:"flag_government"`
	BankLogo             *string    `db:"bank_logo"             json:"bank_logo"`
	CustodianProfile     *string    `db:"custodian_profile"     json:"custodian_profile"`
	RecOrder             *uint64    `db:"rec_order"             json:"rec_order"`
	RecStatus            uint8      `db:"rec_status"            json:"rec_status"`
	RecCreatedDate       *string    `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy         *string    `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate      *string    `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy        *string    `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1            *string    `db:"rec_image1"            json:"rec_image1"`
	RecImage2            *string    `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus    *uint8     `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage     *uint64    `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate      *string    `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy        *string    `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate       *string    `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy         *string    `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1      *string    `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2      *string    `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3      *string    `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}