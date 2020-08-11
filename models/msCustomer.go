package models

type MsCustomer struct {
	CustomerKey              uint64     `db:"customer_key"                json:"customer_key"`
	IDCustomer               uint64     `db:"id_customer"                 json:"id_customer"`
	UnitHolderIDno           string     `db:"unit_holder_idno"            json:"unit_holder_idno"`
	FullName                 string     `db:"full_name"                   json:"full_name"`
	SidNo                    *string    `db:"sid_no"                      json:"sid_no"`
	InvestorType             string     `db:"investor_type"               json:"investor_type"`
	CustomerCategory         string     `db:"customer_category"           json:"customer_category"`
	ParticipantKey           *uint64    `db:"participant_key"             json:"participant_key"`
	CifSuspendFlag           uint8      `db:"cif_suspend_flag"            json:"cif_suspend_flag"`
	CifSuspendModifiedDate   *string    `db:"cif_suspend_modified_date"   json:"cif_suspend_modified_date"`
	CifSuspendReason         *string    `db:"cif_suspend_reason"          json:"cif_suspend_reason"`
	OpenaccBranchKey         *uint64    `db:"openacc_branch_key"          json:"openacc_branch_key"`
	OpenaccAgentKey          *uint64    `db:"openacc_agent_key"           json:"openacc_agent_key"`
	OpenaccDate              *string    `db:"openacc_date"                json:"openacc_date"`
	CloseaccBranchKey        *uint64    `db:"closeacc_branch_key"         json:"closeacc_branch_key"`
	CloseaccAgentKey         *uint64    `db:"closeacc_agent_key"          json:"closeacc_agent_key"`
	CloseaccDate             *string    `db:"closeacc_date"               json:"closeacc_date"`
	RecOrder                 *uint64    `db:"rec_order"                   json:"rec_order"`
	RecStatus                uint8      `db:"rec_status"                  json:"rec_status"`
	RecCreatedDate           *string    `db:"rec_created_date"            json:"rec_created_date"`
	RecCreatedBy             *string    `db:"rec_created_by"              json:"rec_created_by"`
	RecModifiedDate          *string    `db:"rec_modified_date"           json:"rec_modified_date"`
	RecModifiedBy            *string    `db:"rec_modified_by"             json:"rec_modified_by"`
	RecImage1                *string    `db:"rec_image1"                  json:"rec_image1"`
	RecImage2                *string    `db:"rec_image2"                  json:"rec_image2"`
	RecApprovalStatus        *uint8     `db:"rec_approval_status"         json:"rec_approval_status"`
	RecApprovalStage         *uint64    `db:"rec_approval_stage"          json:"rec_approval_stage"`
	RecApprovedDate          *string    `db:"rec_approved_date"           json:"rec_approved_date"`
	RecApprovedBy            *string    `db:"rec_approved_by"             json:"rec_approved_by"`
	RecDeletedDate           *string    `db:"rec_deleted_date"            json:"rec_deleted_date"`
	RecDeletedBy             *string    `db:"rec_deleted_by"              json:"rec_deleted_by"`
	RecAttributeID1          *string    `db:"rec_attribute_id1"           json:"rec_attribute_id1"`
	RecAttributeID2          *string    `db:"rec_attribute_id2"           json:"rec_attribute_id2"`
	RecAttributeID3          *string    `db:"rec_attribute_id3"           json:"rec_attribute_id3"`
}