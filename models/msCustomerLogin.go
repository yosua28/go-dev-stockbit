package models

type MsCustomerLogin struct {
	CustLoginKey           uint64     `db:"cust_login_key"          json:"cust_login_key"`
	CustomerKey            uint64     `db:"customer_key"            json:"customer_key"`
	LoginUserID            string     `db:"login_userid"            json:"login_userid"`
	LoginUserPassword      string     `db:"login_user_password"     json:"login_user_password"`
	LoginUserName          string     `db:"login_user_name"         json:"login_user_name"`
	LoginPinno             *string    `db:"login_pinno"             json:"login_pinno"`
	LoginActive            uint8      `db:"login_active"            json:"login_active"`
	EmailAddress           *string    `db:"email_address"           json:"email_address"`
	EmailDateVerified      *string    `db:"email_date_verified"     json:"email_date_verified"`
	MobileNo               *string    `db:"mobile_no"               json:"mobile_no"`
	MobilenoDateVerified   *string    `db:"mobileno_date_verified"  json:"mobileno_date_verified"`
	LoginCategory          uint64     `db:"login_category"          json:"login_category"`
	LoginCustomerType      uint64     `db:"login_customer_type"     json:"login_customer_type"`
	RecOrder               *uint64    `db:"rec_order"               json:"rec_order"`
	RecStatus              uint8      `db:"rec_status"              json:"rec_status"`
	RecCreatedDate         *string    `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy           *string    `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate        *string    `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy          *string    `db:"rec_modified_by"         json:"rec_modified_by"`
	RecImage1              *string    `db:"rec_image1"              json:"rec_image1"`
	RecImage2              *string    `db:"rec_image2"              json:"rec_image2"`
	RecApprovalStatus      *uint8     `db:"rec_approval_status"     json:"rec_approval_status"`
	RecApprovalStage       *uint64    `db:"rec_approval_stage"      json:"rec_approval_stage"`
	RecApprovedDate        *string    `db:"rec_approved_date"       json:"rec_approved_date"`
	RecApprovedBy          *string    `db:"rec_approved_by"         json:"rec_approved_by"`
	RecDeletedDate         *string    `db:"rec_deleted_date"        json:"rec_deleted_date"`
	RecDeletedBy           *string    `db:"rec_deleted_by"          json:"rec_deleted_by"`
	RecAttributeID1        *string    `db:"rec_attribute_id1"       json:"rec_attribute_id1"`
	RecAttributeID2        *string    `db:"rec_attribute_id2"       json:"rec_attribute_id2"`
	RecAttributeID3        *string    `db:"rec_attribute_id3"       json:"rec_attribute_id3"`
}