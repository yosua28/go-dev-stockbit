package models

import (
	"api/db"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ListAppModuleDropdown struct {
	AppModuleKey  uint64  `db:"app_module_key"        json:"app_module_key"`
	AppModuleCode *string `db:"app_module_code"       json:"app_module_code"`
	AppModuleName *string `db:"app_module_name"       json:"app_module_name"`
}

func AdminGetListAppModuleDropdown(c *[]ListAppModuleDropdown) (int, error) {
	query := `SELECT 
				app_module_key, app_module_code, app_module_name
			FROM sc_app_module
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
