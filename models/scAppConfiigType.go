package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ListDropdownScAppConfigType struct {
	ConfigTypeKey  uint64  `db:"config_type_key"        json:"config_type_key"`
	ConfigTypeCode *string `db:"config_type_code"      json:"config_type_code"`
	ConfigTypeName *string `db:"config_type_name"       json:"config_type_name"`
}

func AdminGetListDropdownScAppConfigType(c *[]ListDropdownScAppConfigType) (int, error) {
	query := `SELECT
				c.config_type_key,
				c.config_type_code,
				c.config_type_name 
			FROM sc_app_config_type AS c
			WHERE c.rec_status = 1 `
	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
