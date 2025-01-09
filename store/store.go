package store

import (
	"context"
	"database/sql"

	"github.com/poohda-go/types"
)

type Store struct {
	Auth interface {
		Login()
	}
	Waitlist interface {
		AddToWaitlist(ctx context.Context, payload types.SubscribePayload) error
		GetAllWaitlistParticipants() ([]types.Waitlist, error)
	}
	Categories interface {
		GetAllCategories() ([]types.Category, error)
		CreateNewCategory(context.Context, types.CategoryDTO) (*types.Category, error)
		GetOneCategory(context.Context, int) (*types.Category, error)
		GetAllClothesReferenceToACategory(context.Context, string) ([]types.Clothes, error)
	}
	Clothes interface {
		CreateNewClothes(ctx context.Context, payload types.ClothesDTO) (*types.Clothes, error)
		GetAllClothes() ([]types.Clothes, error)
		GetOneClothes(context.Context, string) (*types.Clothes, error)
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		// Auth:       &AuthStore{db},
		Waitlist:   &WaitlistStore{db},
		Categories: &CategoriesStore{db},
		Clothes:    &ClothesStore{db},
	}
}
