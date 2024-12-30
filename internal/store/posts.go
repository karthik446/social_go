package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
)

// Post This is the Model
type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type Feed struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tags := pq.Array(post.Tags)
	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, tags).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetById(ctx context.Context, postID int64) (*Post, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at, version FROM posts WHERE id = $1`
	var post Post
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, postID).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return &post, err
}

func (s *PostsStore) DeleteById(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	// Build dynamic query based on provided fields

	updates := make([]string, 0)
	args := make([]interface{}, 0)
	argPosition := 1

	// Always update updated_at
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argPosition))
	args = append(args, time.Now())
	argPosition++

	updates = append(updates, fmt.Sprintf("version = $%d", argPosition))
	args = append(args, post.Version+1) // Increment version
	argPosition++

	if post.Content != "" {
		updates = append(updates, fmt.Sprintf("content = $%d", argPosition))
		args = append(args, post.Content)
		argPosition++
	}

	if post.Title != "" {
		updates = append(updates, fmt.Sprintf("title = $%d", argPosition))
		args = append(args, post.Title)
		argPosition++
	}

	if len(post.Tags) > 0 {
		updates = append(updates, fmt.Sprintf("tags = $%d", argPosition))
		args = append(args, pq.Array(post.Tags))
		argPosition++
	}
	log.Println(updates)

	// If no fields to update except updated_at, return the existing post
	if len(updates) == 2 { // only updated_at
		return nil
	}

	// Add ID to args
	args = append(args, post.ID)
	args = append(args, post.Version) // Current version for comparison

	// Construct final query
	query := fmt.Sprintf(
		"UPDATE posts SET %s WHERE id = $%d AND version = $%d RETURNING id, user_id, content, title, tags, created_at, updated_at, version",
		strings.Join(updates, ", "),
		argPosition,
		argPosition+1,
	)
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	// Execute query and scan into new post
	row := s.db.QueryRowContext(ctx, query, args...)
	err := row.Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.Title,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return err
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]Feed, error) {
	query := `select p.id,
       				 p.user_id,
       				 p.title,
       				 p.content,
       				 p.created_at,
       				 p.version,
       				 p.tags,
       				 u.username,
       				 count(c.id) as comments_count
				from posts p
				         left join comments c on c.post_id = p.id
				         left join followers f on f.follower_id = p.user_id or p.user_id = $1
				         join users u on p.user_id = u.id
				where 
				f.user_id = $1 AND
				(p.title ilike '%' || $4 || '%' or p.content ilike '%' || $4 || '%') AND 
				($5::varchar[] is null OR array_length($5::varchar[], 1) is null OR p.tags && $5::varchar[])
				GROUP BY p.id, u.username, p.created_at
				order by p.created_at ` + fq.Sort + `
				limit $2 offset $3`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows: ", err)
		}
	}(rows)

	feeds := make([]Feed, 0)
	for rows.Next() {
		var f Feed
		f.User = User{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.Title, &f.Content, &f.CreatedAt, &f.Version, pq.Array(&f.Tags), &f.User.Username, &f.CommentsCount); err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}
