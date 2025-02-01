package store

import (
	"context"
	"database/sql"
	// "log"

	"github.com/lib/pq"
	"github.com/poohda-go/types"
)

type ClothesStore struct {
	db *sql.DB
}

func (s *ClothesStore) GetAllClothes() ([]types.Clothes, error) {
	clothings := []types.Clothes{}
	query := `SELECT cl.id, cl.name, cl.price, cl.description, cl.quantity, cl.category_id, array_agg(DISTINCT i.url) AS urls, array_agg(DISTINCT s.size) AS sizes FROM "clothes" AS cl JOIN "image" AS "i" ON i.clothes_id = cl.id LEFT JOIN "clothes_sizes" AS "s" ON cl.id = s.clothes_id GROUP BY  cl.id, cl.name, cl.price, cl.description, cl.quantity, cl.category_id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var clothing types.Clothes

		if err := rows.Scan(
			&clothing.Id,
			&clothing.Name,
			&clothing.Price,
			&clothing.Description,
			&clothing.Quantity,
			&clothing.CategoryId,
			pq.Array(&clothing.Pictures),
			pq.Array(&clothing.Sizes),
		); err != nil {
			return nil, err
		}

		clothings = append(clothings, clothing)
	}

	return clothings, nil
}

func (s *ClothesStore) GetOneClothes(ctx context.Context, id int) (*types.Clothes, error) {
	var clothing types.Clothes
	query := `SELECT cl.id, cl.name, cl.price, cl.description, cl.quantity, array_agg(DISTINCT i.url) AS "pictures", array_agg(DISTINCT s.size) AS "sizes" FROM "clothes" AS cl JOIN "image" AS i ON cl.id = i.clothes_id LEFT JOIN "clothes_sizes" as "s" ON s.clothes_id = cl.id WHERE cl.id=$1 GROUP BY cl.id, cl.name, cl.price, cl.description, cl.quantity;
`

	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&clothing.Id,
		&clothing.Name,
		&clothing.Price,
		&clothing.Description,
		&clothing.Quantity,
		pq.Array(&clothing.Pictures),
		pq.Array(&clothing.Sizes),
	); err != nil {
		return nil, err
	}

	return &clothing, nil
}

func (s *ClothesStore) CreateNewClothes(ctx context.Context, payload types.ClothesDTO) (*types.Clothes, error) {
	var newClothing types.Clothes
	var newImagePictures types.Pictures
	var newImageSize types.Sizes

	// log.Print(payload)
	query := `INSERT INTO "clothes" (name, price, category_id, description, quantity) VALUES ($1, $2, $3, $4, $5) RETURNING id, name,  price, category_id, description, quantity`
	imageQuery := `INSERT INTO "image" (clothes_id, url) VALUES ($1, $2) RETURNING id, url`
	sizeQuery := `INSERT INTO "clothes_sizes" (clothes_id, size) VALUES ($1, $2) RETURNING id, size`

	err := s.db.QueryRowContext(
		ctx,
		query,
		payload.Name,
		payload.Price,
		payload.CategoryId,
		payload.Description,
		payload.Quantity,
	).Scan(
		&newClothing.Id,
		&newClothing.Name,
		&newClothing.Price,
		&newClothing.CategoryId,
		&newClothing.Description,
		&newClothing.Quantity,
	)
	if err != nil {
		return nil, err
	}

	for _, size := range payload.Sizes {
		err := s.db.QueryRowContext(
			ctx,
			sizeQuery,
			newClothing.Id,
			size,
		).Scan(
			&newImageSize.Id,
			&newImageSize.Size,
		)
		if err != nil {
			return nil, err
		}
	}

	for _, picture := range payload.Pictures {
		err := s.db.QueryRowContext(
			ctx,
			imageQuery,
			newClothing.Id,
			picture,
		).Scan(
			&newImagePictures.Id,
			&newImagePictures.Url,
		)
		if err != nil {
			return nil, err
		}
	}

	return &newClothing, nil
}

func (s *ClothesStore) GetClothesThroughName(ctx context.Context, searchString string) ([]types.Clothes, error) {
	clothes := []types.Clothes{}
	query := `SELECT cl.id, cl.name, cl.price, cl.description, cl.quantity, array_agg(DISTINCT i.url) AS "pictures", array_agg(DISTINCT s.size) AS "sizes" FROM "clothes" AS cl JOIN "image" AS i ON cl.id = i.clothes_id LEFT JOIN "clothes_sizes" as "s" ON s.clothes_id = cl.id WHERE cl.name ILIKE $1 GROUP BY cl.id, cl.name, cl.price, cl.description, cl.quantity;`

	rows, err := s.db.QueryContext(ctx, query, "%"+searchString+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var clothing types.Clothes

		if err := rows.Scan(
			&clothing.Id,
			&clothing.Name,
			&clothing.Price,
			&clothing.Description,
			&clothing.Quantity,
			pq.Array(&clothing.Pictures),
			pq.Array(&clothing.Sizes),
		); err != nil {
			return nil, err
		}

		clothes = append(clothes, clothing)
	}

	return clothes, nil
}

func (s *ClothesStore) EditClothes(ctx context.Context, id int) (*types.Clothes, error) {
	// query := `UPDATE "clothes" SET name=$1, price=$2, category_id=$3, description=$4, quantity=$5 WHERE id=$6 RETURNING id, price, category_id, description, quantity`
	// imageQuery := `UPDATE "image" SET url`
	return nil, nil
}

func (s *ClothesStore) DeleteClothes(ctx context.Context, id int) (*types.Clothes, error) {
	tx, err := s.db.BeginTx(ctx, nil) // Start a transaction
	if err != nil {
		return nil, err
	}

	// 1️⃣ Get the clothes details before deleting
	var clothing types.Clothes
	images := []string{}
	query := `SELECT id, name, price, category_id, description, quantity FROM "clothes" WHERE id=$1`
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&clothing.Id,
		&clothing.Name,
		&clothing.Price,
		&clothing.CategoryId,
		&clothing.Description,
		&clothing.Quantity,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2️⃣ Get related images
	imageQuery := `SELECT url FROM "image" WHERE clothes_id=$1`
	rows, err := tx.QueryContext(ctx, imageQuery, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img string
		if err := rows.Scan(&img); err != nil {
			tx.Rollback()
			return nil, err
		}
		images = append(images, img)
	}
	clothing.Pictures = images

	// 3️⃣ Get related sizes
	sizeQuery := `SELECT size FROM "clothes_sizes" WHERE clothes_id=$1`
	sizeRows, err := tx.QueryContext(ctx, sizeQuery, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer sizeRows.Close()

	var sizes []string
	for sizeRows.Next() {
		var size string
		if err := sizeRows.Scan(&size); err != nil {
			tx.Rollback()
			return nil, err
		}
		sizes = append(sizes, size)
	}
	clothing.Sizes = sizes

	// 4️⃣ Delete related images
	_, err = tx.ExecContext(ctx, `DELETE FROM "image" WHERE clothes_id=$1`, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 5️⃣ Delete related sizes
	_, err = tx.ExecContext(ctx, `DELETE FROM "clothes_sizes" WHERE clothes_id=$1`, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 6️⃣ Delete the clothes entry
	_, err = tx.ExecContext(ctx, `DELETE FROM "clothes" WHERE id=$1`, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// ✅ Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &clothing, nil // Return deleted item details
}
