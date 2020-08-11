package models

type TrTransactionType struct{
	TransTypeKey           uint64    `db:"trans_type_key"          json:"trans_type_key"`
	TypeCode              *string    `db:"type_code"               json:"type_code"`
	TypeDescription       *string    `db:"type_description"        json:"type_description"`
	TypeOrder              uint64    `db:"type_order"              json:"type_order"`
	TypeDomain            *string    `db:"type_domain"            json:"type_domain"`
}