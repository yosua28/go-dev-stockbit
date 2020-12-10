package models

import(
	"github.com/shopspring/decimal"
)

type MsAgentProduct struct {
	AgentProductKey      uint64    `db:"aglic_key"              json:"aglic_key"`
	BranchKey            uint64    `db:"branch_key"             json:"branch_key"`
	ProductKey           uint64    `db:"product_key"            json:"product_key"`
	EffDate              *string    `db:"eff_date"              json:"eff_date"`
	ProductNameSa        *string    `db:"product_name_sa"       json:"product_name_sa"`
	HoldingPeriodDays    uint64    `db:"holding_period_days"    json:"holding_period_days"`
	MgtFeeShareSa        decimal.Decimal   `db:"mgt_fee_share_sa"       json:"mgt_fee_share_sa"`
	SubFeeShareSa        decimal.Decimal   `db:"sub_fee_share_sa"       json:"sub_fee_share_sa"`
	RedFeeShareSa        decimal.Decimal   `db:"red_fee_share_sa"       json:"red_fee_share_sa"`
	SwtotFeeShareSa      decimal.Decimal   `db:"swtot_fee_share_sa"     json:"swtot_fee_share_sa"`
	SwtinFeeShareSa      decimal.Decimal   `db:"swtin_fee_share_sa"     json:"swtin_fee_share_sa"`
	OjkFeeShareSa        decimal.Decimal   `db:"ojk_fee_share_sa"       json:"ojk_fee_share_sa"`
	OtherFeeShareSa      decimal.Decimal   `db:"other_fee_share_sa"     json:"other_fee_share_sa"`
	FlagEnabled          uint8     `db:"flag_enabled"           json:"flag_enabled"`
	FlagSubscription     uint8     `db:"flag_subscription"      json:"flag_subscription"`
	FlagRedemption       uint8     `db:"flag_redemption"        json:"flag_redemption"`
	FlagSwitchOut        uint8     `db:"flag_switch_out"        json:"flag_switch_out"`
	FlagSwitchIn         uint8     `db:"flag_switch_in"         json:"flag_switch_in"`
	MaxSubFee            decimal.Decimal   `db:"max_sub_fee"            json:"max_sub_fee"`
	MaxRedFee            decimal.Decimal   `db:"max_red_fee"            json:"max_red_fee"`
	MaxSwiFee            decimal.Decimal   `db:"max_swi_fee"            json:"max_swi_fee"`
	MinSubAmount         decimal.Decimal   `db:"min_sub_amount"         json:"min_sub_amount"`
	MinRedAmount         decimal.Decimal   `db:"min_red_amount"         json:"min_red_amount"`
	MinRedUnit           decimal.Decimal   `db:"min_red_unit"           json:"min_red_unit"`
	MinUnitAfterRed      decimal.Decimal   `db:"min_unit_after_red"     json:"min_unit_after_red"`
	RecOrder             *uint64   `db:"rec_order"              json:"rec_order"`
	RecStatus            uint8     `db:"rec_status"             json:"rec_status"`
	RecCreatedDate       *string   `db:"rec_created_date"       json:"rec_created_date"`
	RecCreatedBy         *string   `db:"rec_created_by"         json:"rec_created_by"`
	RecModifiedDate      *string   `db:"rec_modified_date"      json:"rec_modified_date"`
	RecModifiedBy        *string   `db:"rec_modified_by"        json:"rec_modified_by"`
	RecImage1            *string   `db:"rec_image1"             json:"rec_image1"`
	RecImage2            *string   `db:"rec_image2"             json:"rec_image2"`
	RecApprovalStatus    *uint8    `db:"rec_approval_status"    json:"rec_approval_status"`
	RecApprovalStage     *uint64   `db:"rec_approval_stage"     json:"rec_approval_stage"`
	RecApprovedDate      *string   `db:"rec_approved_date"      json:"rec_approved_date"`
	RecApprovedBy        *string   `db:"rec_approved_by"        json:"rec_approved_by"`
	RecDeletedDate       *string   `db:"rec_deleted_date"       json:"rec_deleted_date"`
	RecDeletedBy         *string   `db:"rec_deleted_by"         json:"rec_deleted_by"`
	RecAttributeID1      *string   `db:"rec_attribute_id1"      json:"rec_attribute_id1"`
	RecAttributeID2      *string   `db:"rec_attribute_id2"      json:"rec_attribute_id2"`
	RecAttributeID3      *string   `db:"rec_attribute_id3"      json:"rec_attribute_id3"`
}