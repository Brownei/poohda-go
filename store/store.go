package store

import (
	"context"
	"database/sql"

	"github.com/poohda-go/types"
)

type Store struct {
	Waitlist interface {
		AddToWaitlist(ctx context.Context, payload types.SubscribePayload) error
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Waitlist: &WaitlistStore{db},
	}
}
