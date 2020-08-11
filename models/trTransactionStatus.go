package models

type TrTransactionStatus struct{
	TransStatusKey         uint64    `db:"trans_status_key"        json:"trans_status_key"`
	StatusCode            *string    `db:"status_code"             json:"status_code"`
	StatusDescription     *string    `db:"status_description"      json:"status_description"`
	StatusOrder            uint64    `db:"status_order"            json:"status_order"`
	StatusPhase           *string    `db:"status_phase"            json:"status_phase"`
}