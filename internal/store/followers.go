package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) Follow(ctx context.Context, FollowerID int64, UserID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, UserID, FollowerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrDuplicateKeyConflict
		}
	}
	return nil
}

func (s *FollowersStore) UnFollow(ctx context.Context, FollowerID int64, UserID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`
	_, err := s.db.ExecContext(ctx, query, UserID, FollowerID)
	if err != nil {
		return err
	}
	return nil
}
