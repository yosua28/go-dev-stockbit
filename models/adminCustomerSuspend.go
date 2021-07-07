package models

import (
	"api/db"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type CustomerIndividuStatusSuspend struct {
	CustomerKey         uint64  `db:"customer_key"                json:"customer_key"`
	Cif                 string  `db:"cif"                         json:"cif"`
	FullName            string  `db:"full_name"                   json:"full_name"`
	DateBirth           string  `db:"date_birth"                  json:"date_birth"`
	IDcardNo            string  `db:"ktp"                         json:"ktp"`
	PhoneMobile         string  `db:"phone_mobile"                json:"phone_mobile"`
	Email               string  `db:"email"                       json:"email"`
	SidNo               string  `db:"sid"                         json:"sid"`
	CifSuspendFlag      string  `db:"cif_suspend_flag"            json:"cif_suspend_flag"`
	SuspendModifiedDate *string `db:"suspend_modified_date"       json:"suspend_modified_date"`
	SuspendReason       *string `db:"suspend_reason"              json:"suspend_reason"`
	OaStatus            string  `db:"oa_status"                   json:"oa_status"`
	BranchKey           *uint64 `db:"branch_key"                  json:"branch_key"`
	BranchName          *string `db:"branch_name"                 json:"branch_name"`
	AgentKey            *uint64 `db:"agent_key"                   json:"agent_key"`
	AgentName           *string `db:"agent_name"                  json:"agent_name"`
	CreatedDate         *string `db:"created_date"                json:"created_date"`
}

func AdminGetAllCustomerStatusSuspend(c *[]CustomerIndividuStatusSuspend, limit uint64, offset uint64, params map[string]string, paramsLike map[string]string, nolimit bool) (int, error) {
	var present bool
	var whereClause []string
	var condition string
	var limitOffset string
	var orderCondition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by

	query := `SELECT 
				c.customer_key AS customer_key,
				c.unit_holder_idno AS cif, 
				c.full_name AS full_name, 
				DATE_FORMAT(pd.date_birth, '%d %M %Y') AS date_birth, 
				pd.idcard_no AS ktp, 
				pd.email_address AS email, 
				pd.phone_mobile AS phone_mobile, 
				(CASE
					WHEN c.sid_no IS NULL THEN "-"
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN c.cif_suspend_flag = 0 THEN "Tidak"
					ELSE "Ya"
				END) AS cif_suspend_flag, 
				DATE_FORMAT(c.cif_suspend_modified_date, '%d %M %Y %H:%i') AS suspend_modified_date, 
				c.cif_suspend_reason AS suspend_reason, 
				r.oa_status AS oa_status, 
				c.openacc_branch_key AS branch_key, 
				c.openacc_agent_key AS agent_key, 
				br.branch_name AS branch_name, 
				ag.agent_name AS agent_name, 
				DATE_FORMAT(c.rec_created_date, '%d %M %Y %H:%i') AS created_date 
			FROM ms_customer AS c 
			INNER JOIN oa_request AS r ON r.customer_key = c.customer_key 
			INNER JOIN (
			SELECT MAX(oa_request_key) AS oa_request_key, customer_key FROM oa_request WHERE rec_status = 1 AND oa_status > 259 GROUP BY customer_key
			) AS t2 ON r.oa_request_key = t2.oa_request_key
			INNER JOIN oa_personal_data AS pd ON pd.oa_request_key = r.oa_request_key
			LEFT JOIN ms_branch AS br ON br.branch_key = c.openacc_branch_key 
			LEFT JOIN ms_agent AS ag ON ag.agent_key = c.openacc_agent_key 
			WHERE c.rec_status = 1 AND r.rec_status = 1 AND pd.rec_status = 1
			AND r.customer_key IS NOT NULL` + condition + ` 
			GROUP BY c.customer_key`

	var orderBy string
	var orderType string
	if orderBy, present = params["orderBy"]; present == true {
		orderCondition += " ORDER BY " + orderBy
		if orderType, present = params["orderType"]; present == true {
			orderCondition += " " + orderType
		}
	}

	if !nolimit {
		limitOffset += " LIMIT " + strconv.FormatUint(limit, 10)
		if offset > 0 {
			limitOffset += " OFFSET " + strconv.FormatUint(offset, 10)
		}
	}

	query += orderCondition + limitOffset

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func CountAdminGetAllCustomerStatusSuspend(c *CountData, params map[string]string, paramsLike map[string]string) (int, error) {
	var whereClause []string
	var condition string

	for field, value := range params {
		if !(field == "orderBy" || field == "orderType") {
			whereClause = append(whereClause, field+" = '"+value+"'")
		}
	}

	for fieldLike, valueLike := range paramsLike {
		whereClause = append(whereClause, fieldLike+" like '%"+valueLike+"%'")
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by

	query := `SELECT 
				COUNT(c.customer_key) AS count_data
			FROM ms_customer AS c 
			INNER JOIN oa_request AS r ON r.customer_key = c.customer_key 
			INNER JOIN (
			SELECT MAX(oa_request_key) AS oa_request_key, customer_key FROM oa_request WHERE rec_status = 1 AND oa_status > 259 GROUP BY customer_key
			) AS t2 ON r.oa_request_key = t2.oa_request_key
			INNER JOIN oa_personal_data AS pd ON pd.oa_request_key = r.oa_request_key
			LEFT JOIN ms_branch AS br ON br.branch_key = r.branch_key 
			LEFT JOIN ms_agent AS ag ON ag.agent_key = r.agent_key 
			WHERE c.rec_status = 1 AND r.rec_status = 1 AND pd.rec_status = 1
			AND r.customer_key IS NOT NULL ` + condition

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}

func AdminGetDetailCustomerStatusSuspend(c *CustomerIndividuStatusSuspend, params map[string]string) (int, error) {
	var whereClause []string
	var condition string

	for field, value := range params {
		whereClause = append(whereClause, field+" = '"+value+"'")
	}

	// Combile where clause
	if len(whereClause) > 0 {
		condition += " AND "
		for index, where := range whereClause {
			condition += where
			if (len(whereClause) - 1) > index {
				condition += " AND "
			}
		}
	}
	// Check order by

	query := `SELECT 
				c.customer_key AS customer_key,
				c.unit_holder_idno AS cif, 
				c.full_name AS full_name, 
				DATE_FORMAT(pd.date_birth, '%d %M %Y') AS date_birth, 
				pd.idcard_no AS ktp, 
				pd.email_address AS email, 
				pd.phone_mobile AS phone_mobile, 
				(CASE
					WHEN c.sid_no IS NULL THEN "-"
					ELSE c.sid_no
				END) AS sid,
				(CASE
					WHEN c.cif_suspend_flag = 0 THEN "Tidak"
					ELSE "Ya"
				END) AS cif_suspend_flag, 
				DATE_FORMAT(c.cif_suspend_modified_date, '%d %M %Y %H:%i') AS suspend_modified_date, 
				c.cif_suspend_reason AS suspend_reason, 
				r.oa_status AS oa_status, 
				c.openacc_branch_key AS branch_key, 
				c.openacc_agent_key AS agent_key, 
				br.branch_name AS branch_name, 
				ag.agent_name AS agent_name, 
				DATE_FORMAT(c.rec_created_date, '%d %M %Y %H:%i') AS created_date 
			FROM ms_customer AS c 
			INNER JOIN oa_request AS r ON r.customer_key = c.customer_key 
			INNER JOIN (
			SELECT MAX(oa_request_key) AS oa_request_key, customer_key FROM oa_request WHERE rec_status = 1 AND oa_status > 259 GROUP BY customer_key
			) AS t2 ON r.oa_request_key = t2.oa_request_key
			INNER JOIN oa_personal_data AS pd ON pd.oa_request_key = r.oa_request_key
			LEFT JOIN ms_branch AS br ON br.branch_key = c.openacc_branch_key 
			LEFT JOIN ms_agent AS ag ON ag.agent_key = c.openacc_agent_key 
			WHERE c.rec_status = 1 AND r.rec_status = 1 AND pd.rec_status = 1
			AND r.customer_key IS NOT NULL` + condition + ` 
			GROUP BY c.customer_key limit 1`

	// Main query
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
