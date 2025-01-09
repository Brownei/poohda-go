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
			{
				Id: "2",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "clothes" (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, description VARCHAR(100) NOT NULL, price INT, category_id INT NOT NULL, quantity INT, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
				},
				Down: []string{`"DROP TABLE IF EXISTS "clothes"`},
			},

			{
				Id: "3",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "category" (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL UNIQUE, picture VARCHAR(100))`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "category"`,
				},
			},

			{
				Id: "4",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "image" (id SERIAL PRIMARY KEY, clothes_id INT, url VARCHAR(100));`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "image"`,
				},
			},

			{
				Id: "5",
				Up: []string{
					`ALTER TABLE "clothes" ADD CONSTRAINT "Clothes_CategoryId_fkey" FOREIGN KEY ("category_id") REFERENCES "category"("id")`,
				},
				Down: []string{
					`ALTER TABLE "clothes" DROP CONSTRAINT IF EXISTS "Clothes_CategoryId_fkey"`,
				},
			},

			{
				Id: "6",
				Up: []string{
					`ALTER TABLE "image" ADD CONSTRAINT "Image_clothesId_fkey" FOREIGN KEY ("clothes_id") REFERENCES "clothes"("id")`,
				},
				Down: []string{
					`ALTER TABLE "image" DROP CONSTRAINT IF EXISTS "Image_clothesId_fkey"`,
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
