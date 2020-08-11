package models

type MsAgentAgreement struct {
	AgreementKey         uint64    `db:"agreement_key"        json:"agreement_key"`
	AgreementNo          *string   `db:"agreement_no"         json:"agreement_no"`
	AgreementDate        string    `db:"agreement_date"       json:"agreement_date"`
	AgreementSubject     string    `db:"agreement_subject"    json:"agreement_subject"`
	AgreementContent     *string   `db:"agreement_content"    json:"agreement_content"`
	SignedDate           *string   `db:"signed_date"          json:"signed_date"`
	SignCity             *string   `db:"sign_city"            json:"sign_city"`
	BranchKey            *uint64   `db:"branch_key"           json:"branch_key"`
	AgreementStatus      uint64    `db:"agreement_status"     json:"agreement_status"`
	RecOrder             *uint64   `db:"rec_order"            json:"rec_order"`
	RecStatus            uint8     `db:"rec_status"           json:"rec_status"`
	RecCreatedDate       *string   `db:"rec_created_date"     json:"rec_created_date"`
	RecCreatedBy         *string   `db:"rec_created_by"       json:"rec_created_by"`
	RecModifiedDate      *string   `db:"rec_modified_date"    json:"rec_modified_date"`
	RecModifiedBy        *string   `db:"rec_modified_by"      json:"rec_modified_by"`
	RecImage1            *string   `db:"rec_image1"           json:"rec_image1"`
	RecImage2            *string   `db:"rec_image2"           json:"rec_image2"`
	RecApprovalStatus    *uint8    `db:"rec_approval_status"  json:"rec_approval_status"`
	RecApprovalStage     *uint64   `db:"rec_approval_stage"   json:"rec_approval_stage"`
	RecApprovedDate      *string   `db:"rec_approved_date"    json:"rec_approved_date"`
	RecApprovedBy        *string   `db:"rec_approved_by"      json:"rec_approved_by"`
	RecDeletedDate       *string   `db:"rec_deleted_date"     json:"rec_deleted_date"`
	RecDeletedBy         *string   `db:"rec_deleted_by"       json:"rec_deleted_by"`
	RecAttributeID1      *string   `db:"rec_attribute_id1"    json:"rec_attribute_id1"`
	RecAttributeID2      *string   `db:"rec_attribute_id2"    json:"rec_attribute_id2"`
	RecAttributeID3      *string   `db:"rec_attribute_id3"    json:"rec_attribute_id3"`
}