package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/poohda-go/types"
)

type WaitlistStore struct {
	db *sql.DB
}

func (s *WaitlistStore) AddToWaitlist(ctx context.Context, payload types.SubscribePayload) error {
	var waitlist types.SubscribePayload
	var id int
	findUserQuery := `SELECT id FROM "waitlist" WHERE email=$1`
	if err := s.db.QueryRowContext(
		ctx,
		findUserQuery,
		payload.Email,
	).Scan(&id); err == nil {
		return fmt.Errorf("You have already joined the circle")
	} else {
		if err != sql.ErrNoRows {
			// Handle unexpected database errors
			return fmt.Errorf("Database error: %v", err)
		}
	}
	// fmt.Print(&id)

	query := `INSERT INTO "waitlist" (name, email, number) VALUES ($1, $2, $3) RETURNING name, email`

	if err := s.db.QueryRowContext(
		ctx,
		query,
		payload.Name,
		payload.Email,
		payload.Number,
	).Scan(
		&waitlist.Name,
		&waitlist.Email,
	); err != nil {
		return err
	}

	return nil
}

func (s *WaitlistStore) GetAllWaitlistParticipants() ([]types.Waitlist, error) {
	waitlst := []types.Waitlist{}
	query := `SELECT name, email, number FROM "waitlist"`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var individualInWaiting types.Waitlist

		if err := rows.Scan(
			&individualInWaiting.Name,
			&individualInWaiting.Email,
			&individualInWaiting.Number,
		); err != nil {
			return nil, err
		}

		waitlst = append(waitlst, individualInWaiting)
	}

	return waitlst, nil
}
