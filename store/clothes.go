package store

import (
	"context"
	"database/sql"

	"github.com/poohda-go/types"
)

type ClothesStore struct {
	db *sql.DB
}

func (s *ClothesStore) GetAllClothes() ([]types.Clothes, error) {
	clothings := []types.Clothes{}
	query := `SELECT id, name, price, description, quantity FROM "clothes"`

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
	query := `INSERT INTO "clothes" (name, price, category_id, description, quantity) VALUES ($1, $2, $3, $4, $5) RETURNING id, name,  price, category_id, description, quantity`

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

	return &newClothing, nil
}
