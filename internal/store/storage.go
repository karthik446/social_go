package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound             = errors.New("resource not found")
	ErrDuplicateKeyConflict = errors.New("duplicate key value violates unique constraint")
	QueryTimeOutDuration    = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		DeleteById(context.Context, int64) error
		GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]Feed, error)
	}

	Users interface {
		Create(context.Context, *User) error
		GetById(ctx context.Context, userID int64) (*User, error)
	}
	Comments interface {
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
		Create(ctx context.Context, c *Comment) error
		DeleteById(ctx context.Context, commentID int64) error
		GetById(ctx context.Context, commentID int64) (*Comment, error)
		Update(ctx context.Context, c *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, FollowerID int64, UserID int64) error
		UnFollow(ctx context.Context, FollowerID int64, UserID int64) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowersStore{db},
	}
}
