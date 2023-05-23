package comment_repo

import (
	"database/sql"
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrDeleteFailed error = errors.New("delete failed")
)

type TestCommentLikesRepo struct {
	CommentLikes map[uuid.UUID]*models.CommentLike
	sync.Mutex
}

func NewCommentLikesRepo() *TestCommentLikesRepo {
	return &TestCommentLikesRepo{
		CommentLikes: make(map[uuid.UUID]*models.CommentLike),
	}
}

type PostgresCommentLikesRepo struct {
	DB *sql.DB
}

func NewPostgresCommentLikesRepo(DB *sql.DB) *PostgresCommentLikesRepo {
	return &PostgresCommentLikesRepo{
		DB: DB,
	}
}

func (p *PostgresCommentLikesRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS commentLikes(
        id UUID PRIMARY KEY,
    	userId UUID,
        commentId UUID
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (t *TestCommentLikesRepo) CreateCommentLike(commentLike models.CommentLike) error {
	t.CommentLikes[commentLike.ID] = &commentLike
	return nil
}

func (p *PostgresCommentLikesRepo) CreateCommentLike(commentLike models.CommentLike) error {
	_, err := p.DB.Exec("INSERT INTO comments(id, userId, commentId) values($1, $2, $3)", commentLike.ID, commentLike.UserId, commentLike.CommentId)
	return err
}

func (t *TestCommentLikesRepo) Delete(id uuid.UUID) error {
	delete(t.CommentLikes, id)
	return nil
}

func (p *PostgresCommentLikesRepo) Delete(id uuid.UUID) error {
	res, err := p.DB.Exec("DELETE FROM commentLikes WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}
	return nil
}
