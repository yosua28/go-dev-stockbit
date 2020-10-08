package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ScUserMessage struct {
	UmessageKey          uint64  `db:"umessage_key"              json:"umessage_key"`
	UmessageType         *uint64 `db:"umessage_type"             json:"umessage_type"`
	UmessageRecipientKey uint64  `db:"umessage_recipient_key"    json:"umessage_recipient_key"`
	UmessageReceiptDate  *string `db:"umessage_receipt_date"     json:"umessage_receipt_date"`
	FlagRead             uint8   `db:"flag_read"                 json:"flag_read"`
	UmessageSenderKey    *uint64 `db:"umessage_sender_key"       json:"umessage_sender_key"`
	UmessageSentDate     *string `db:"umessage_sent_date"        json:"umessage_sent_date"`
	FlagSent             *uint8  `db:"flag_sent"                 json:"flag_sent"`
	UmessageSubject      *string `db:"umessage_subject"          json:"umessage_subject"`
	UmessageBody         *string `db:"umessage_body"             json:"umessage_body"`
	UparentKey           *uint64 `db:"uparent_key"               json:"uparent_key"`
	UmessageCategory     *uint64 `db:"umessage_category"         json:"umessage_category"`
	FlagArchieved        *uint8  `db:"flag_archieved"            json:"flag_archieved"`
	ArchievedDate        *string `db:"archieved_date"            json:"archieved_date"`
	RecOrder             *uint64 `db:"rec_order"                 json:"rec_order"`
	RecStatus            uint8   `db:"rec_status"                json:"rec_status"`
	RecCreatedDate       *string `db:"rec_created_date"          json:"rec_created_date"`
	RecCreatedBy         *string `db:"rec_created_by"            json:"rec_created_by"`
	RecModifiedDate      *string `db:"rec_modified_date"         json:"rec_modified_date"`
	RecModifiedBy        *string `db:"rec_modified_by"           json:"rec_modified_by"`
	RecImage1            *string `db:"rec_image1"                json:"rec_image1"`
	RecImage2            *string `db:"rec_image2"                json:"rec_image2"`
	RecApprovalStatus    *uint8  `db:"rec_approval_status"       json:"rec_approval_status"`
	RecApprovalStage     *uint64 `db:"rec_approval_stage"        json:"rec_approval_stage"`
	RecApprovedDate      *string `db:"rec_approved_date"         json:"rec_approved_date"`
	RecApprovedBy        *string `db:"rec_approved_by"           json:"rec_approved_by"`
	RecDeletedDate       *string `db:"rec_deleted_date"          json:"rec_deleted_date"`
	RecDeletedBy         *string `db:"rec_deleted_by"            json:"rec_deleted_by"`
	RecAttributeID1      *string `db:"rec_attribute_id1"         json:"rec_attribute_id1"`
	RecAttributeID2      *string `db:"rec_attribute_id2"         json:"rec_attribute_id2"`
	RecAttributeID3      *string `db:"rec_attribute_id3"         json:"rec_attribute_id3"`
}

func CreateScUserMessage(params map[string]string) (int, error) {
	query := "INSERT INTO sc_user_message"
	// Get params
	var fields, values string
	var bindvars []interface{}
	for key, value := range params {
		fields += key + ", "
		values += "?, "
		bindvars = append(bindvars, value)
	}
	fields = fields[:(len(fields) - 2)]
	values = values[:(len(values) - 2)]

	// Combine params to build query
	query += "(" + fields + ") VALUES(" + values + ")"
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}

	_, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
