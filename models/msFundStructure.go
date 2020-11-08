package models

import (
	"api/db"
	"log"
	"net/http"
)

type MsFundStructure struct {
	FundStructureKey  uint64  `db:"fund_structure_key"    json:"fund_structure_key"`
	FundStructureCode string  `db:"fund_structure_code"   json:"fund_structure_code"`
	FundStructureName string  `db:"fund_structure_name"   json:"fund_structure_name"`
	FundStructureDesc *string `db:"fund_structure_desc"   json:"fund_structure_desc"`
	RecOrder          *uint64 `db:"rec_order"             json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"            json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"            json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

type MsFundStructureInfo struct {
	FundStructureKey  uint64  `json:"fund_structure_key"`
	FundStructureCode string  `json:"fund_structure_code"`
	FundStructureName string  `json:"fund_structure_name"`
	FundStructureDesc *string `json:"fund_structure_desc"`
}

func GetMsFundStructure(c *MsFundStructure, key string) (int, error) {
	query := `SELECT ms_fund_structure.* FROM ms_fund_structure 
				WHERE ms_fund_structure.rec_status = 1 AND ms_fund_structure.fund_structure_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
