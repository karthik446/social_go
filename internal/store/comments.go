package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id FROM comments c
    JOIN users u ON c.user_id = u.id
    WHERE c.post_id = $1
    ORDER BY c.created_at DESC`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (s *CommentsStore) GetById(ctx context.Context, commentID int64) (*Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	var c Comment
	c.User = User{}
	err := s.db.QueryRowContext(ctx, query, commentID).Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CommentsStore) Create(ctx context.Context, c *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, c.PostID, c.UserID, c.Content).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentsStore) DeleteById(ctx context.Context, commentID int64) error {
	query := `DELETE FROM comments WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentsStore) Update(ctx context.Context, c *Comment) error {
	query := `UPDATE comments SET content = $1 WHERE id = $2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, c.Content, c.ID)
	if err != nil {
		return err
	}
	return nil
}
