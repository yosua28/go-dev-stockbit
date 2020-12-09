package models

import (
	"api/db"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

type UdfValue struct {
	UdfValueKey       uint64  `db:"udf_value_key"         json:"udf_value_key"`
	UdfInfoKey        uint64  `db:"udf_info_key"          json:"udf_info_key"`
	RowDataKey        uint64  `db:"row_data_key"          json:"row_data_key"`
	UdfValues        *string  `db:"udf_value"             json:"udf_value"`
}


func GetUdfValueIn(c *[]UdfValue, value []string, field string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT
				udf_value.* FROM 
				udf_value `
	query := query2 + " WHERE udf_value." + field + " IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CreateMultipleUdfValue(params []interface{}) (int, error) {
	
		q := `INSERT INTO udf_value ( 
			udf_info_key,
			row_data_key,
			udf_values) VALUES `
	
		for i := 0; i < len(params); i++ {
			q += "(?)"
			if i < (len(params) - 1) {
				q += ","
			}
		}
		log.Info(q)
		query, args, err := sqlx.In(q, params...)
		if err != nil {
			return http.StatusBadGateway, err
		}
	
		query = db.Db.Rebind(query)
		_, err = db.Db.Query(query, args...)
		if err != nil {
			log.Error(err.Error())
			return http.StatusBadGateway, err
		}
		return http.StatusOK, nil
	}