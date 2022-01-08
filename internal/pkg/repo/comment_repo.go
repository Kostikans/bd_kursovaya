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

func (p *CommentRepo) GetCommentVote(ctx context.Context, comment model.CommentVote) (res model.CommentVote, err error) {
	query := `
SELECT id,post_id,author_id,comment_id,vote FROM comment_vote WHERE author_id = $1 AND post_id = $2 AND comment_id = $3`

	err = p.db.Get(&res, query, comment.AuthorID, comment.PostID, comment.CommentID)
	return
}

func (p *CommentRepo) CreateCommentVote(ctx context.Context, comment model.CommentVote) (id uint64, err error) {
	query := `
INSERT INTO comment_vote(author_id, post_id,comment_id,vote) VALUES($1,$2,$3,$4)
 ON CONFLICT (post_id,author_id,comment_id) DO UPDATE SET vote = $4 RETURNING id`

	err = p.db.Get(&id, query, comment.AuthorID, comment.PostID, comment.CommentID, comment.Vote)
	return
}

func (p *CommentRepo) IncrementCommentVote(ctx context.Context, commentID uint64, likeCount int64, dislikeCount int64) (id uint64, err error) {
	query := `
WITH curr AS (SELECT comment_id,like_count,dislike_count FROM comment_vote_agg WHERE comment_id = $1)
INSERT INTO comment_vote_agg(comment_id,like_count,dislike_count) VALUES($1,$2,$3) 
ON CONFLICT (comment_id) DO UPDATE SET like_count = (SELECT like_count FROM curr) + $2, dislike_count = (SELECT dislike_count FROM curr) + $3 RETURNING id`

	err = p.db.Get(&id, query, commentID, likeCount, dislikeCount)
	return
}
