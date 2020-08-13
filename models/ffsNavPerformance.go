package models

import (
	"api/db"
	"log"
	"strconv"
	"net/http"
)

type FfsNavPerformance struct {
	NavPerformanceKey    uint64    `db:"nav_performance_key"   json:"nav_performance_key"`
	ProductKey           uint64    `db:"product_key"           json:"product_key"`
	PeriodeKey           uint64    `db:"periode_key"           json:"periode_key"`
	NavDate              string    `db:"nav_date"              json:"nav_date"`
	NavD0                float32   `db:"nav_d0"                json:"nav_d0"`
	NavD1                float32   `db:"nav_d1"                json:"nav_d1"`
	NavM0                float32   `db:"nav_m0"                json:"nav_m0"`
	NavM1                float32   `db:"nav_m1"                json:"nav_m1"`
	NavM3                float32   `db:"nav_m3"                json:"nav_m3"`
	NavM6                float32   `db:"nav_m6"                json:"nav_m6"`
	NavYtd               float32   `db:"nav_ytd"               json:"nav_ytd"`
	Navy1                float32   `db:"nav_y1"                json:"nav_y1"`
	Navy3                float32   `db:"nav_y3"                json:"nav_y3"`
	Navy5                float32   `db:"nav_y5"                json:"nav_y5"`
	PerformD1            float32   `db:"perform_d1"            json:"perform_d1"`
	PerformMtd           float32   `db:"perform_mtd"           json:"perform_mtd"`
	PerformM1            float32   `db:"perform_m1"            json:"perform_m1"`
	PerformM3            float32   `db:"perform_m3"            json:"perform_m3"`
	PerformM6            float32   `db:"perform_m6"            json:"perform_m6"`
	PerformYtd           float32   `db:"perform_ytd"           json:"perform_ytd"`
	PerformY1            float32   `db:"perform_y1"            json:"perform_y1"`
	PerformY3            float32   `db:"perform_y3"            json:"perform_y3"`
	PerformY5            float32   `db:"perform_y5"            json:"perform_y5"`
	PerformCagr          float32   `db:"perform_cagr"          json:"perform_cagr"`
	PerformAll           float32   `db:"perform_all"           json:"perform_all"`
	RecOrder             *uint64   `db:"rec_order"             json:"rec_order"`
	RecStatus            uint8     `db:"rec_status"            json:"rec_status"`
	RecCreatedDate       *string   `db:"rec_created_date"      json:"rec_created_date"`
	RecCreatedBy         *string   `db:"rec_created_by"        json:"rec_created_by"`
	RecModifiedDate      *string   `db:"rec_modified_date"     json:"rec_modified_date"`
	RecModifiedBy        *string   `db:"rec_modified_by"       json:"rec_modified_by"`
	RecImage1            *string   `db:"rec_image1"            json:"rec_image1"`
	RecImage2            *string   `db:"rec_image2"            json:"rec_image2"`
	RecApprovalStatus    *uint8    `db:"rec_approval_status"   json:"rec_approval_status"`
	RecApprovalStage     *uint64   `db:"rec_approval_stage"    json:"rec_approval_stage"`
	RecApprovedDate      *string   `db:"rec_approved_date"     json:"rec_approved_date"`
	RecApprovedBy        *string   `db:"rec_approved_by"       json:"rec_approved_by"`
	RecDeletedDate       *string   `db:"rec_deleted_date"      json:"rec_deleted_date"`
	RecDeletedBy         *string   `db:"rec_deleted_by"        json:"rec_deleted_by"`
	RecAttributeID1      *string   `db:"rec_attribute_id1"     json:"rec_attribute_id1"`
	RecAttributeID2      *string   `db:"rec_attribute_id2"     json:"rec_attribute_id2"`
	RecAttributeID3      *string   `db:"rec_attribute_id3"     json:"rec_attribute_id3"`
}

func GetAllFfsNavPerformance(c *[]FfsNavPerformance, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ffs_nav_performance.* FROM 
			  ffs_nav_performance`
	var present bool
	var whereClause []string
	var condition string
	
	for field, value := range params {
		if !(field == "orderBy" || field == "orderType"){
			whereClause = append(whereClause, "ffs_nav_performance."+field + " = '" + value + "'")
		}
	} 

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " WHERE "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by
	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition

	// Query limit and offset
	if !nolimit {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			query += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}