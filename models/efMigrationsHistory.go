package models

type EfMigrationsHistory struct{
	MigrationID        string    `db:"migrationid"          json:"migrationid"`
	ProductVersion     string    `db:"productversion"       json:"productversion"`
}