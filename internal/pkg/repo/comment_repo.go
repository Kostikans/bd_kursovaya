package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

// CommentRepo  - comment repo
type CommentRepo struct {
	db *sqlx.DB
}

// NewCommentRepo - returns new commentRepo
func NewCommentRepo(db *sqlx.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

func (p *CommentRepo) CreateComment(ctx context.Context, comment model.Comment) (id uint64, err error) {
	query := `INSERT INTO comment(author_id, post_id, parent_id, text) VALUES($1,$2,$3,$4) RETURNING id`

	err = p.db.Get(&id, query, comment.AuthorID, comment.PostID, comment.ParentID, comment.Text)
	return
}

func (p *CommentRepo) CheckAuthorAndPostExist(ctx context.Context, comment model.Comment) (exist bool, err error) {
	query := `
SELECT COALESCE((SELECT 1 FROM account where id = $1),0)::bool  AND
       COALESCE((SELECT 1 FROM post where id = $2),0)::bool`

	err = p.db.Get(&exist, query, comment.AuthorID, comment.PostID)
	return
}
