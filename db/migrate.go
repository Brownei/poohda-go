package db

import (
	"database/sql"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

func AddMigrations(db *sql.DB) {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "waitlist" (id BIGSERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, email VARCHAR(255) NOT NULL UNIQUE, number VARCHAR(15) NOT NULL)`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "waitlist"`,
				},
			},
		},
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Couldn't apply the migrations: %s", err)
	}

	log.Printf("Applied %d migrations!", n)
}
