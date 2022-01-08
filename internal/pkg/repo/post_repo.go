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
	query := `SELECT EXISTS(SELECT * FROM account WHERE id = $1)`

	err = p.db.Get(&exist, query, post.AuthorID)
	return
}

func (p *PostRepo) CheckPostExist(ctx context.Context, postID uint64) (exist bool, err error) {
	query := `SELECT EXISTS(SELECT * FROM post WHERE id = $1)`

	err = p.db.Get(&exist, query, postID)
	return
}

func (p *PostRepo) GetPostVote(ctx context.Context, post model.PostVote) (res model.PostVote, err error) {
	query := `
SELECT id,post_id,author_id,vote FROM post_vote WHERE author_id = $1 AND post_id = $2`

	err = p.db.Get(&res, query, post.AuthorID, post.PostID)
	return
}

func (p *PostRepo) CreatePostVote(ctx context.Context, post model.PostVote) (id uint64, err error) {
	query := `
INSERT INTO post_vote(author_id, post_id,vote) VALUES($1,$2,$3)
 ON CONFLICT (post_id,author_id) DO UPDATE SET vote = $3 RETURNING id`

	err = p.db.Get(&id, query, post.AuthorID, post.PostID, post.Vote)
	return
}

func (p *PostRepo) IncrementPostVote(ctx context.Context, postID uint64, likeCount int64, dislikeCount int64) (id uint64, err error) {

	query := `
WITH curr AS (SELECT post_id,like_count,dislike_count FROM post_vote_agg WHERE post_id = $1)
INSERT INTO post_vote_agg(post_id,like_count,dislike_count) VALUES($1,$2,$3) 
ON CONFLICT (post_id) DO UPDATE SET like_count = (SELECT like_count FROM curr) + $2, dislike_count = (SELECT dislike_count FROM curr) + $3 RETURNING id`

	err = p.db.Get(&id, query, postID, likeCount, dislikeCount)
	return
}

func (p *PostRepo) GetPosts(ctx context.Context, limit uint32, cursor uint64) (rows []model.ExtendedPost, next uint64, err error) {
	query := `
SELECT post.id,author_id,title,text,created_at, from post INNER JOIN post_vote_agg pva ON post.id = pva.post_id
`

	err = p.db.Select(&rows, query, cursor, limit+1)
	if err != nil {
		return
	}

	if limit == 0 {
		limit = basePaginationLimit
	}

	if len(rows) > int(limit) {
		rows = rows[:len(rows)-1]
		next = rows[len(rows)-1].ID
	}

	return
}
