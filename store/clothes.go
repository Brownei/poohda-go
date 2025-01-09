package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/lib/pq"
	"github.com/poohda-go/types"
)

type ClothesStore struct {
	db *sql.DB
}

func (s *ClothesStore) GetAllClothes() ([]types.Clothes, error) {
	clothings := []types.Clothes{}
	query := `SELECT cl.id, cl.name, cl.price, cl.description, cl.quantity, cl.category_id, array_agg(i.url) AS urls FROM "clothes" AS cl JOIN "image" AS "i" ON i.clothes_id = cl.id GROUP BY cl.id, cl.id, cl.name, cl.price, cl.description, cl.quantity, cl.category_id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

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
		); err != nil {
			return nil, err
		}

		clothings = append(clothings, clothing)
	}

	return clothings, nil
}

func (s *ClothesStore) GetOneClothes(ctx context.Context, name string) (*types.Clothes, error) {
	var clothing types.Clothes
	query := `SELECT id, name, price, description, quantity FROM "clothes" WHERE name=$1`

	if err := s.db.QueryRowContext(ctx, query, name).Scan(
		&clothing.Id,
		&clothing.Name,
		&clothing.Price,
		&clothing.Description,
		&clothing.Quantity,
	); err != nil {
		return nil, err
	}

	return &clothing, nil
}

func (s *ClothesStore) CreateNewClothes(ctx context.Context, payload types.ClothesDTO) (*types.Clothes, error) {
	var newClothing types.Clothes
	var newImagePictures types.ClothesPictures
	log.Print(payload)
	query := `INSERT INTO "clothes" (name, price, category_id, description, quantity) VALUES ($1, $2, $3, $4, $5) RETURNING id, name,  price, category_id, description, quantity`
	imageQuery := `INSERT INTO "image" (clothes_id, url) VALUES ($1, $2) RETURNING id, url`

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
