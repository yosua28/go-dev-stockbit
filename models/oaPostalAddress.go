package models

import (
	"api/db"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type OaPostalAddress struct {
	PostalAddressKey  uint64  `db:"postal_address_key"         json:"postal_address_key"`
	AddressType       string  `db:"address_type"               json:"address_type"`
	KabupatenKey      uint64  `db:"kabupaten_key"              json:"kabupaten_key"`
	KecamatanKey      uint64  `db:"kecamatan_key"              json:"kecamatan_key"`
	AddressLine1      *string `db:"address_line1"              json:"address_line1"`
	AddressLine2      *string `db:"address_line2"              json:"address_line2"`
	AddressLine3      *string `db:"address_line3"              json:"address_line3"`
	PostalCode        *string `db:"postal_code"                json:"postal_code"`
	GeolocName        *string `db:"geoloc_name"                json:"geoloc_name"`
	GeolocLongitude   *string `db:"geoloc_longitude"           json:"geoloc_longitude"`
	GeolocLatitude    *string `db:"geoloc_latitude"            json:"geoloc_latitude"`
	RecOrder          *uint64 `db:"rec_order"                  json:"rec_order"`
	RecStatus         uint8   `db:"rec_status"                 json:"rec_status"`
	RecCreatedDate    *string `db:"rec_created_date"           json:"rec_created_date"`
	RecCreatedBy      *string `db:"rec_created_by"             json:"rec_created_by"`
	RecModifiedDate   *string `db:"rec_modified_date"          json:"rec_modified_date"`
	RecModifiedBy     *string `db:"rec_modified_by"            json:"rec_modified_by"`
	RecImage1         *string `db:"rec_image1"                 json:"rec_image1"`
	RecImage2         *string `db:"rec_image2"                 json:"rec_image2"`
	RecApprovalStatus *uint8  `db:"rec_approval_status"        json:"rec_approval_status"`
	RecApprovalStage  *uint64 `db:"rec_approval_stage"         json:"rec_approval_stage"`
	RecApprovedDate   *string `db:"rec_approved_date"          json:"rec_approved_date"`
	RecApprovedBy     *string `db:"rec_approved_by"            json:"rec_approved_by"`
	RecDeletedDate    *string `db:"rec_deleted_date"           json:"rec_deleted_date"`
	RecDeletedBy      *string `db:"rec_deleted_by"             json:"rec_deleted_by"`
	RecAttributeID1   *string `db:"rec_attribute_id1"          json:"rec_attribute_id1"`
	RecAttributeID2   *string `db:"rec_attribute_id2"          json:"rec_attribute_id2"`
	RecAttributeID3   *string `db:"rec_attribute_id3"          json:"rec_attribute_id3"`
}

func CreateOaPostalAddress(params map[string]string) (int, error, string) {
	query := "INSERT INTO oa_postal_address"
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
		return http.StatusBadGateway, err, "0"
	}
	var ret sql.Result
	ret, err = tx.Exec(query, bindvars...)
	tx.Commit()
	if err != nil {
		log.Error(err)
		return http.StatusBadRequest, err, "0"
	}
	lastID, _ := ret.LastInsertId()
	return http.StatusOK, nil, strconv.FormatInt(lastID, 10)
}

func GetOaPostalAddressIn(c *[]OaPostalAddress, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				oa_postal_address.* FROM 
				oa_postal_address `
	query := query2 + " WHERE oa_postal_address." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func GetOaPostalAddress(c *OaPostalAddress, key string) (int, error) {
	query := `SELECT oa_postal_address.* FROM oa_postal_address WHERE oa_postal_address.postal_address_key = ` + key
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
