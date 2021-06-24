package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ListMenuTypeDropdown struct {
	MenuTypeKey  uint64  `db:"menu_type_key"        json:"menu_type_key"`
	MenuTypeCode *string `db:"menu_type_code"       json:"menu_type_code"`
	MenuTypeName *string `db:"menu_type_name"       json:"menu_type_name"`
}

func AdminGetListMenuTypeDropdown(c *[]ListMenuTypeDropdown) (int, error) {
	query := `SELECT 
				menu_type_key, menu_type_code, menu_type_name
			FROM sc_menu_type
			WHERE rec_status = 1 ORDER BY rec_order ASC`
	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
