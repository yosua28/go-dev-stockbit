package models

import (
	"api/db"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type FfsPublish struct {
	FfsKey            uint64  `db:"ffs_key"              json:"ffs_key"`
	PeriodeKey        uint64  `db:"periode_key"          json:"periode_key"`
	ProductKey        uint64  `db:"product_key"          json:"product_key"`
	FfsLink           *string `db:"ffs_link"             json:"ffs_link"`
	DatePeriode       string  `db:"date_periode"         json:"date_periode"`
	DatePublished     *string `db:"date_published"       json:"date_published"`
	FfsNameDisplayed  *string `db:"ffs_name_displayed"   json:"ffs_name_displayed"`
	RecOrder          *uint64 `db:"rec_order"            json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"           json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"     json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"       json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"    json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"      json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"           json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"           json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"  json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"   json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"    json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"      json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"     json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"       json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"    json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"    json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"    json:"rec_attribute_id3"`
}

func GetLastFfsIn(c *[]FfsPublish, productKey []string) (int, error) {
	inQuery := strings.Join(productKey, ",")
	query2 := `SELECT MAX( ffs_key ) as ffs_key, product_key, ffs_link FROM
			   ffs_publish`
	query := query2 + " WHERE ffs_publish.product_key IN(" + inQuery + ") GROUP BY product_key"

	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetLastFfsByPeriodeIn(c *[]FfsPublish, productKey []string) (int, error) {
	inQuery := strings.Join(productKey, ",")
	query := `SELECT a.ffs_key, a.product_key, a.ffs_link 
				FROM ffs_publish AS a
				INNER JOIN (
				SELECT product_key, MAX(periode_key) AS periode_key 
				FROM ffs_publish
				WHERE rec_status = 1
				GROUP BY product_key
				) b ON (a.product_key = b.product_key AND a.periode_key=b.periode_key)
				WHERE a.product_key  IN(` + inQuery + `) GROUP BY product_key`

	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
