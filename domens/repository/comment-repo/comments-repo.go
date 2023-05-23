package comment_repo

import (
	"database/sql"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

type TestCommentRepo struct {
	Comments map[uuid.UUID]*models.Comment
	sync.Mutex
}

func NewCommentRepo() *TestCommentRepo {
	return &TestCommentRepo{
		Comments: make(map[uuid.UUID]*models.Comment),
	}
}

type PostgresCommentRepo struct {
	DB *sql.DB
}

func NewPostgresCommentRepo(DB *sql.DB) *PostgresCommentRepo {
	return &PostgresCommentRepo{
		DB: DB,
	}
}

func (p *PostgresCommentRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS comments(
        id UUID PRIMARY KEY,
    	userId UUID,
        gameId UUID,
		parentId UUID NULL,
		value VARCHAR NOT NULL,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (t *TestCommentRepo) CreateComment(comment models.Comment) error {
	t.Comments[comment.ID] = &comment
	return nil
}

func (p *PostgresCommentRepo) CreateComment(comment models.Comment) error {
	_, err := p.DB.Exec("INSERT INTO comments(id, userId, gameId, parentId, value, createdAt) values($1, $2, $3,  $4, $5, $6)", comment.ID, comment.UserId, comment.GameId, comment.ParentId, comment.Value, comment.CreatedAt)
	return err
}

func (t *TestCommentRepo) Delete(id uuid.UUID) error {
	delete(t.Comments, id)
	return nil
}

func (p *PostgresCommentRepo) Delete(id uuid.UUID) error {
	res, err := p.DB.Exec("DELETE FROM comments WHERE id = $1", id)
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
