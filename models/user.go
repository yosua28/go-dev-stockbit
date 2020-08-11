package models

type User struct{
	ID                  uint64    `db:"id"               json:"id"`
	Username           *string    `db:"username"         json:"username"`
	Password           *string    `db:"password"         json:"password"`
	Role               *string    `db:"role"             json:"role"`
}