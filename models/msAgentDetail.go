package models

type MsAgentDetail struct {
	AgentKey        uint64    `db:"agent_key"            json:"agent_key"`
	MobileNo        uint64    `db:"mobile_no"            json:"mobile_no"`
	MobileNoAlt     *uint64   `db:"mobile_no_alt"        json:"mobile_no_alt"`
	EmailAddress    string    `db:"email_address"        json:"email_address"`
	BirtDate        *string   `db:"birth_date"           json:"birth_date"`
	BirthPlace      *string   `db:"birth_place"          json:"birth_place"`
	Gender          uint64    `db:"gender"               json:"gender"`
	Nationality     uint64    `db:"nationality"          json:"nationality"`
	CountryKey      *uint64   `db:"country_key"          json:"country_key"`
	IDType          uint64    `db:"IDType"               json:"IDType"`
	IDNo            *string   `db:"id_no"                json:"id_no"`
	Occupation      *string   `db:"occupation"           json:"occupation"`
	CompanyName     *string   `db:"company_name"         json:"company_name"`
	Position        *string   `db:"position"             json:"position"`
	PositionLevel   uint64    `db:"position_level"       json:"position_level"`
	JoinDate        string    `db:"join_date"            json:"join_date"`
	FlagResign      uint8     `db:"flag_resign"          json:"flag_resign"`
}