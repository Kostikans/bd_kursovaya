package repo

import (
	"context"

	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
	"github.com/kostikan/bd_kursovaya/internal/pkg/sql"
)

// PostRepo - post repo
type PostRepo struct {
	db *sql.Balancer
}

// NewPostRepo - returns new postRepo
func NewPostRepo(db *sql.Balancer) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

func (p *PostRepo) CreatePost(ctx context.Context, post model.Post) (id uint64, err error) {
	query := `INSERT INTO post(author_id,title, text) VALUES($1,$2,$3) RETURNING id`

	err = p.db.Write(ctx).Get(&id, query, post.AuthorID, post.Title, post.Text)
	return
}

func (p *PostRepo) CheckAuthorExist(ctx context.Context, post model.Post) (exist bool, err error) {
	query := `SELECT EXISTS(SELECT * FROM account WHERE id = $1)`

	err = p.db.Read(ctx).Get(&exist, query, post.AuthorID)
	return
}

func (p *PostRepo) CheckPostExist(ctx context.Context, postID uint64) (exist bool, err error) {
	query := `SELECT EXISTS(SELECT * FROM post WHERE id = $1)`

	err = p.db.Read(ctx).Get(&exist, query, postID)
	return
}

func (p *PostRepo) GetPostVote(ctx context.Context, post model.PostVote) (res model.PostVote, err error) {
	query := `
SELECT id,post_id,author_id,vote FROM post_vote WHERE author_id = $1 AND post_id = $2`

	err = p.db.Write(ctx).Get(&res, query, post.AuthorID, post.PostID)
	return
}

func (p *PostRepo) CreatePostVote(ctx context.Context, post model.PostVote) (id uint64, err error) {
	query := `
INSERT INTO post_vote(author_id, post_id,vote) VALUES($1,$2,$3)
 ON CONFLICT (post_id,author_id) DO UPDATE SET vote = $3 RETURNING id`

	err = p.db.Write(ctx).Get(&id, query, post.AuthorID, post.PostID, post.Vote)

	return
}

func (p *PostRepo) IncrementPostVote(ctx context.Context, postID uint64, likeCount int64, dislikeCount int64) (id uint64, err error) {

	query := `
INSERT INTO post_vote_agg(post_id,like_count,dislike_count) VALUES($1,$2,$3) 
ON CONFLICT (post_id) DO UPDATE SET like_count = excluded.like_count + $2,
                                    dislike_count = excluded.dislike_count + $3 RETURNING id`

	err = p.db.Write(ctx).Get(&id, query, postID, likeCount, dislikeCount)
	return
}

func (p *PostRepo) GetPosts(ctx context.Context, limit uint32, cursor uint64) (rows []model.ExtendedPost, next uint64, err error) {
	query := `
SELECT post.id,author_id,title,text,created_at,COALESCE(like_count,0) as like_count,COALESCE(dislike_count,0) as dislike_count,tags_id
from post LEFT JOIN post_vote_agg pva ON post.id = pva.post_id WHERE post.id > $1 ORDER BY post.id LIMIT $2
`
	if limit == 0 {
		limit = basePaginationLimit
	}

	err = p.db.Read(ctx).Select(&rows, query, cursor, limit+1)
	if err != nil {
		return
	}

	if len(rows) > int(limit) {
		rows = rows[:len(rows)-1]
		next = rows[len(rows)-1].ID
	}

	return
}
