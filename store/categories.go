package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/poohda-go/types"
)

type CategoriesStore struct {
	db *sql.DB
}

func (s *CategoriesStore) GetAllCategories() ([]types.Category, error) {
	categories := []types.Category{}
	query := `SELECT c.id, c.name, c.description, c.is_featured, array_agg(ci.url) AS pictures FROM "category" AS c JOIN "category_image" AS "ci" ON ci.category_id = c.id GROUP BY c.id, c.name, c.description`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category types.Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.IsFeatured,
			pq.Array(&category.Pictures),
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
	var newImagePictures types.Pictures
	query := `INSERT INTO "category" (name, description, is_featured) VALUES ($1, $2, $3) RETURNING id, name, description, is_featured`
	imageQuery := `INSERT INTO "category_image" (category_id, url) VALUES ($1, $2) RETURNING id, url`

	err := s.db.QueryRowContext(
		ctx,
		query,
		payload.Name,
		payload.Description,
		payload.IsFeatured,
	).Scan(
		&newCategory.Id,
		&newCategory.Name,
		&newCategory.Description,
		&newCategory.IsFeatured,
	)
	if err != nil {
		return nil, err
	}

	for _, picture := range payload.Pictures {
		err := s.db.QueryRowContext(
			ctx,
			imageQuery,
			newCategory.Id,
			picture,
		).Scan(
			&newImagePictures.Id,
			&newImagePictures.Url,
		)
		if err != nil {
			return nil, err
		}
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

	query := `SELECT cl.id, cl.name, cl.price, cl.description, cl.quantity FROM "clothes" AS cl JOIN "category" AS c ON cl.category_id = c.id WHERE c.name=$1`

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

func (s *CategoriesStore) GetOneCategory(ctx context.Context, id int) (*types.Category, error) {
	var category types.Category
	query := `SELECT c.id, c.name, c.description, c.is_featured, array_agg(ci.url) AS pictures FROM "category" AS c JOIN "category_image" AS "ci" ON c.id = ci.category_id WHERE c.id=$1 GROUP BY c.id, c.name, c.description `

	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&category.Id,
		&category.Name,
		&category.Description,
		&category.IsFeatured,
		pq.Array(&category.Pictures),
	); err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoriesStore) EditCategory(ctx context.Context, id int, payload types.CategoryDTO) (*types.Category, error) {
	var newCategory types.Category
	var newImagePictures types.Pictures
	query := `UPDATE "category" SET name=$1, description=$2, is_featured=$3 WHERE id=$4 RETURNING id, name, description, is_featured`
	imageQuery := `UPDATE "category_image" SET url=$1 WHERE category_id=$2 RETURNING id, url `

	err := s.db.QueryRowContext(
		ctx,
		query,
		payload.Name,
		payload.Description,
		payload.IsFeatured,
		id,
	).Scan(
		&newCategory.Id,
		&newCategory.Name,
		&newCategory.Description,
		&newCategory.IsFeatured,
	)
	if err != nil {
		return nil, err
	}

	for _, picture := range payload.Pictures {
		err := s.db.QueryRowContext(
			ctx,
			imageQuery,
			picture,
			id,
		).Scan(
			&newImagePictures.Id,
			&newImagePictures.Url,
		)
		if err != nil {
			return nil, err
		}
	}

	return &newCategory, nil
}

func (s *CategoriesStore) DeleteCategory(ctx context.Context, id int) (*types.Category, error) {
	return nil, nil
}
