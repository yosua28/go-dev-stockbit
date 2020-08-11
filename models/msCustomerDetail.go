package models

type MsCustomerDetail struct {
	CustomerKey     uint64    `db:"customer_key"      json:"customer_key"`
	Nationality     uint64    `db:"nationality"       json:"nationality"`
	Gender          uint64    `db:"gender"            json:"gender"`
	IDType          uint64    `db:"id_type"           json:"id_type"`
	IDNumber        *string   `db:"id_number"         json:"id_number"`
	IDHolderName    *string   `db:"id_holder_name"    json:"id_holder_name"`
	FlagEmployee    uint8     `db:"flag_employee"     json:"flag_employee"`
	FlagGroup       uint8     `db:"flag_group"        json:"flag_group"`
}