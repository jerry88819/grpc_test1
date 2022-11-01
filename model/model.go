package model

import "database/sql"

var (
    db              *sql.DB
)

func Init(_db *sql.DB) {
    db = _db
} // Init()