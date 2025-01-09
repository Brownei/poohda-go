package store

import "database/sql"

type AuthStore struct {
	db *sql.DB
}
