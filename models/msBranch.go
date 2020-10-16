package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"
)

type MsBranch struct {
	BranchKey          uint64  `db:"branch_key"            json:"branch_key"`
	ParticipantKey     *uint64 `db:"participant_key"       json:"participant_key"`
	BranchCode         string  `db:"branch_code"           json:"branch_code"`
	BranchName         string  `db:"branch_name"           json:"branch_name"`
	BranchCategory     *uint64 `db:"branch_category"       json:"branch_category"`
	CityKey            *uint64 `db:"city_key"              json:"city_key"`
	BranchAddress      *string `db:"branch_address"        json:"branch_address"`
	BranchEstablished  *string `db:"branch_established"    json:"branch_established"`
	BranchPicName      *string `db:"branch_pic_name"       json:"branch_pic_name"`
	BranchPicEmail     *string `db:"branch_pic_email"      json:"branch_pic_email"`
	BranchPicPhoneno   *string `db:"branch_pic_phoneno"    json:"branch_pic_phoneno"`
	BranchCostCenter   *string `db:"branch_cost_center"    json:"branch_cost_center"`
	BranchProfitCenter *string `db:"branch_profit_center"  json:"branch_profit_center"`
	RecOrder           *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus          uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate     *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy       *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate    *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy      *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1          *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2          *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus  *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage   *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate    *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy      *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate     *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy       *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1    *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2    *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3    *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

func GetMsBranchIn(c *[]MsBranch, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				ms_branch.* FROM 
				ms_branch `
	query := query2 + " WHERE ms_branch.rec_status = 1 AND ms_branch." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
