package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/poohda-go/types"
)

type CategoriesStore struct {
	db *sql.DB
}

func (s *CategoriesStore) GetAllCategories() ([]types.Category, error) {
	categories := []types.Category{}
	query := `SELECT id, name, picture FROM "category"`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category types.Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Picture,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (s *CategoriesStore) CreateNewCategory(ctx context.Context, payload types.CategoryDTO) (*types.Category, error) {
	var newCategory types.Category
	query := `INSERT INTO "category" (name, picture) VALUES ($1, $2) RETURNING id, name, picture`

	err := s.db.QueryRowContext(
		ctx,
		query,
		payload.Name,
		payload.Picture,
	).Scan(
		&newCategory.Id,
		&newCategory.Name,
		&newCategory.Picture,
	)
	if err != nil {
		return nil, err
	}

	return &newCategory, nil
}

func (s *CategoriesStore) GetAllClothesReferenceToACategory(ctx context.Context, categoryName string) ([]types.Clothes, error) {
	clothings := []types.Clothes{}
	findQuery := `SELECT id FROM "category" WHERE name=$1`
	if err := s.db.QueryRowContext(ctx, findQuery, categoryName).Scan(); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No category like this!")
		}

		return nil, err
	}

	query := `SELECT cl.id, cl.name, cl.rating, cl.price, cl.description, cl.is_featured, cl.quantity, cl.is_best_sales FROM "clothes" AS cl JOIN "category" AS c ON cl.category_id = c.id WHERE c.name=$1`

	rows, err := s.db.Query(query, categoryName)
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

func (s *CategoriesStore) GetOneCategory(ctx context.Context, name string) (*types.Category, error) {
	var category types.Category
	query := `SELECT id, name, picture FROM "category" WHERE name=$1`

	if err := s.db.QueryRowContext(ctx, query, name).Scan(
		&category.Id,
		&category.Name,
		&category.Picture,
	); err != nil {
		return nil, err
	}

	return &category, nil
}
