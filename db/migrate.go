package db

import (
	"database/sql"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

func AddMigrations(db *sql.DB) {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{},
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Couldn't apply the migrations: %s", err)
	}

	log.Printf("Applied %d migrations!", n)
}
