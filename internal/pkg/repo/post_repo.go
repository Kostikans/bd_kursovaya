package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

// PostRepo - post repo
type PostRepo struct {
	db *sqlx.DB
}

// NewPostRepo - returns new postRepo
func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

func (p *PostRepo) CreatePost(ctx context.Context, post model.Post) (id uint64, err error) {
	query := `INSERT INTO post(author_id,title, text) VALUES($1,$2,$3) RETURNING id`

	err = p.db.Get(&id, query, post.AuthorID, post.Title, post.Text)
	return
}

func (p *PostRepo) CheckAuthorExist(ctx context.Context, post model.Post) (exist bool, err error) {
	query := `SELECT EXISTS(SELECT * FROM account where id = $1)`

	err = p.db.Get(&exist, query, post.AuthorID)
	return
}

func (p *PostRepo) CheckPostExist(ctx context.Context, postID uint64) (exist bool, err error) {
	query := `SELECT EXISTS(SELECT * FROM post where id = $1)`

	err = p.db.Get(&exist, query, postID)
	return
}
