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
				Down: []string{`DROP TABLE IF EXISTS "clothes"`},
			},

			{
				Id: "3",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "category" (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL UNIQUE, description VARCHAR(100) NOT NULL, is_featured BOOLEAN DEFAULT false)`,
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

			{
				Id: "7",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "clothes_sizes" (id SERIAL PRIMARY KEY, clothes_id INT, size VARCHAR(20));`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "clothes_sizes"`,
				},
			},

			{
				Id: "8",
				Up: []string{
					`ALTER TABLE "clothes_sizes" ADD CONSTRAINT "ClothesSizes_clothesId_fkey" FOREIGN KEY ("clothes_id") REFERENCES "clothes"("id")`,
				},
				Down: []string{
					`ALTER TABLE "clothes_sizes" DROP CONSTRAINT IF EXISTS "ClothesSizes_clothesId_fkey"`,
				},
			},

			{
				Id: "9",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "category_image" (id SERIAL PRIMARY KEY, category_id INT, url VARCHAR(100));`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "category_image"`,
				},
			},

			{
				Id: "10",
				Up: []string{
					`ALTER TABLE "category_image" ADD CONSTRAINT "CategoryImage_categoryId_fkey" FOREIGN KEY ("category_id") REFERENCES "category"("id")`,
				},
				Down: []string{
					`ALTER TABLE "category_image" DROP CONSTRAINT IF EXISTS "CategoryImage_clothesId_fkey"`,
				},
			},

			{
				Id: "11",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "orders" (id SERIAL PRIMARY KEY, name VARCHAR(255), quantity INT, address VARCHAR(255), price INT, is_delivered BOOLEAN DEFAULT FALSE, clothes_bought_id INT)`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "orders"`,
				},
			},

			{
				Id: "12",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "clothes_bought" (id SERIAL PRIMARY KEY, order_id INT, clothe_id INT, quantity INT)`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "clothes_bought"`,
				},
			},

			{
				Id: "13",
				Up: []string{
					`ALTER TABLE "orders" ADD CONSTRAINT "Orders_clothesBoughtId_fkey" FOREIGN KEY ("clothes_bought_id") REFERENCES "clothes_bought"("id")`,
				},
				Down: []string{
					`ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "Orders_clothesBoughtId_fkey"`,
				},
			},

			{
				Id: "14",
				Up: []string{
					`ALTER TABLE "clothes_bought" ADD CONSTRAINT "ClothesBought_clothesId_fkey" FOREIGN KEY ("clothe_id") REFERENCES "clothes"("id")`,
				},
				Down: []string{
					`ALTER TABLE "clothes_bought" DROP CONSTRAINT IF EXISTS "ClotheBought_clothesId_fkey"`,
				},
			},

			{
				Id: "15",
				Up: []string{
					`ALTER TABLE "clothes_bought" ADD CONSTRAINT "ClothesBought_orderId_fkey" FOREIGN KEY ("order_id") REFERENCES "orders"("id")`,
				},
				Down: []string{
					`ALTER TABLE "clothes_bought" DROP CONSTRAINT IF EXISTS "ClothesBought_orderId_fkey"`,
				},
			},
		},
	}

	// _, err := migrate.Exec(db, "postgres", migrations, migrate.Down)
	// if err != nil {
	// 	log.Fatalf("Couldn't reset changes: %s", err)
	// }

	up, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Couldn't apply the migrations: %s", err)
	}

	log.Printf("Applied %d migrations!", up)
}
