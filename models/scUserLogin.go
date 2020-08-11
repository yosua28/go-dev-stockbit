package models

type ScUserLogin struct {
	UloginKey               uint64    `db:"ulogin_key"                json:"ulogin_key"`
	UloginName              string    `db:"ulogin_name"               json:"ulogin_name"`
	UloginFullName          string    `db:"ulogin_full_name"          json:"ulogin_full_name"`
	UloginPassword          string    `db:"ulogin_password"           json:"ulogin_password"`
	UloginEmail             string    `db:"ulogin_email"              json:"ulogin_email"`
	UloginPin               *string   `db:"ulogin_pin"                json:"ulogin_pin"`
	UloginMobileno          *string   `db:"ulogin_mobileno"           json:"ulogin_mobileno"`
	UloginRole              *string   `db:"ulogin_role"               json:"ulogin_role"`
	UloginLocked            uint8     `db:"ulogin_locked"             json:"ulogin_locked"`
	UloginEnable            uint8     `db:"ulogin_enable"             json:"ulogin_enable"`
	UloginFailedCount       uint64    `db:"ulogin_failed_count"       json:"ulogin_failed_count"`
	UloginLastAccess        string    `db:"ulogin_last_access"        json:"ulogin_last_access"`
	UloginEmailVerified     string    `db:"ulogin_email_verified"     json:"ulogin_email_verified"`
	UloginMobilenoVerified  string    `db:"ulogin_mobileno_verified"  json:"ulogin_mobileno_verified"`
}