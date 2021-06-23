package models

import (
	"api/db"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MsFile struct {
	FileKey            uint64       `db:"file_key"                json:"file_key"`
	RefFkKey           uint64       `db:"ref_fk_key"              json:"ref_fk_key"`
	RefFkDomain        string       `db:"ref_fk_domain"           json:"ref_fk_domain"`
	FileName           string       `db:"file_name"               json:"file_name"`
	FileExt            string       `db:"file_ext"                json:"file_ext"`
	BlobMode           uint8        `db:"blob_mode"               json:"blob_mode"`
	FilePath           *string      `db:"file_path"               json:"file_path"`
	FileUrl            *string      `db:"file_url"                json:"file_url"`
	FileNotes          *string      `db:"file_notes"              json:"file_notes"`
	FileObj            *interface{} `db:"file_obj"                json:"properties"`
	RecCreatedDate     *string      `db:"rec_created_date"        json:"rec_created_date"`
	RecCreatedBy       *string      `db:"rec_created_by"          json:"rec_created_by"`
	RecModifiedDate    *string      `db:"rec_modified_date"       json:"rec_modified_date"`
	RecModifiedBy      *string      `db:"rec_modified_by"         json:"rec_modified_by"`
	RecStatus          uint8        `db:"rec_status"              json:"rec_status"`
	CmsPostPostKey     uint64       `db:"cms_postpost_key"        json:"cms_postpost_key"`
	CmsQuizQuestionKey uint64       `db:"cms_quiz_question_key"   json:"cms_quiz_question_key"`
}

type MsFileDetail struct {
	FileKey        uint64  `db:"file_key"                json:"file_key"`
	RefFkKey       uint64  `db:"ref_fk_key"              json:"ref_fk_key"`
	FileName       string  `db:"file_name"               json:"file_name"`
	Path           string  `db:"path"                    json:"path"`
	FileExt        string  `db:"file_ext"                json:"file_ext"`
	FileNotes      *string `db:"file_notes"              json:"file_notes"`
	RecCreatedDate *string `db:"rec_created_date"        json:"rec_created_date"`
}

type CustomerDocumentDetail struct {
	Customer CustomerIndividuStatusSuspend `json:"customer"`
	Document []MsFileDetail                `json:"document"`
}

func GetAllMsFile(c *[]MsFile, limit uint64, offset uint64, params map[string]string, nolimit bool) (int, error) {
	query := `SELECT
              ms_file.* FROM 
			  ms_file`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_file."+field+" = '"+value+"'")
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
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func UpdateMsFile(params map[string]string) (int, error) {
	query := "UPDATE ms_file SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "file_key" {

			query += key + " = '" + value + "'"

			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE file_key = " + params["file_key"]
	log.Info(query)

	tx, err := db.Db.Begin()
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}
	// var ret sql.Result
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		log.Error(err)
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return http.StatusOK, nil
}

func CreateMsFile(params map[string]string) (int, error) {
	query := "INSERT INTO ms_file"
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

func GetALlDetailMsFile(c *[]MsFileDetail, params map[string]string) (int, error) {
	query := `SELECT 
				ms_file.file_key,
				ms_file.ref_fk_key,
				ms_file.file_name,
				ms_file.file_ext,
				ms_file.file_notes,
				DATE_FORMAT(ms_file.rec_created_date, '%d %M %Y %H:%i') AS rec_created_date 
			FROM ms_file AS ms_file`
	var present bool
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, "ms_file."+field+" = '"+value+"'")
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

	// Main query
	log.Info(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Error(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
