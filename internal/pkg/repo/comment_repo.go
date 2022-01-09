package repo

import (
	"context"

	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
	"github.com/kostikan/bd_kursovaya/internal/pkg/sql"
)

// CommentRepo  - comment repo
type CommentRepo struct {
	db *sql.Balancer
}

// NewCommentRepo - returns new commentRepo
func NewCommentRepo(db *sql.Balancer) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

func (p *CommentRepo) CreateComment(ctx context.Context, comment model.Comment) (id uint64, err error) {
	query := `INSERT INTO comment(author_id, post_id, parent_id, text) VALUES($1,$2,$3,$4) RETURNING id`

	err = p.db.Write(ctx).Get(&id, query, comment.AuthorID, comment.PostID, comment.ParentID, comment.Text)
	return
}

func (p *CommentRepo) CheckAuthorAndPostExist(ctx context.Context, comment model.Comment) (exist bool, err error) {
	query := `
SELECT COALESCE((SELECT 1 FROM account where id = $1),0)::bool  AND
       COALESCE((SELECT 1 FROM post where id = $2),0)::bool`

	err = p.db.Read(ctx).Get(&exist, query, comment.AuthorID, comment.PostID)
	return
}

func (p *CommentRepo) GetCommentVote(ctx context.Context, comment model.CommentVote) (res model.CommentVote, err error) {
	query := `
SELECT id,post_id,author_id,comment_id,vote FROM comment_vote WHERE author_id = $1 AND post_id = $2 AND comment_id = $3`

	err = p.db.Read(ctx).Get(&res, query, comment.AuthorID, comment.PostID, comment.CommentID)
	return
}

func (p *CommentRepo) CreateCommentVote(ctx context.Context, comment model.CommentVote) (id uint64, err error) {
	query := `
INSERT INTO comment_vote(author_id, post_id,comment_id,vote) VALUES($1,$2,$3,$4)
 ON CONFLICT (post_id,author_id,comment_id) DO UPDATE SET vote = $4 RETURNING id`

	err = p.db.Write(ctx).Get(&id, query, comment.AuthorID, comment.PostID, comment.CommentID, comment.Vote)
	return
}

func (p *CommentRepo) IncrementCommentVote(ctx context.Context, commentID uint64, likeCount int64, dislikeCount int64) (id uint64, err error) {
	query := `
WITH curr AS (SELECT comment_id,like_count,dislike_count FROM comment_vote_agg WHERE comment_id = $1)
INSERT INTO comment_vote_agg(comment_id,like_count,dislike_count) VALUES($1,$2,$3) 
ON CONFLICT (comment_id) DO UPDATE SET like_count = (SELECT like_count FROM curr) + $2, dislike_count = (SELECT dislike_count FROM curr) + $3 RETURNING id`

	err = p.db.Write(ctx).Get(&id, query, commentID, likeCount, dislikeCount)
	return
}

func (p *CommentRepo) GetComments(ctx context.Context, postID uint64, limit uint32, cursor uint64) (rows []model.ExtendedComment, next uint64, err error) {
	query := `
WITH RECURSIVE comments_tree(id, post_id,author_id,created_at, parent_id, text, breadcrumbs) AS (
	SELECT id, post_id,author_id, created_at,parent_id,text, ARRAY [id] as breadcrumbs
	FROM comment
	WHERE parent_id = 0 AND post_id = $1
	UNION ALL
	SELECT comment.id, comment.post_id, comment.author_id, comment.created_at, comment.parent_id,comment.text,  com_tree.breadcrumbs || comment.id
	FROM comment
			 JOIN comments_tree com_tree ON comment.parent_id = com_tree.id
	WHERE NOT comment.id = ANY (com_tree.breadcrumbs) AND comment.post_id = $1
)
SELECT ct.id, ct.post_id,author_id,created_at, parent_id, text,breadcrumbs,COALESCE(like_count,0) as like_count,COALESCE(dislike_count,0) as dislike_count
FROM comments_tree ct LEFT JOIN post_vote_agg pva ON ct.id = pva.post_id WHERE ct.id > $2 ORDER BY ct.id LIMIT $3
`
	if limit == 0 {
		limit = basePaginationLimit
	}

	err = p.db.Read(ctx).Select(&rows, query, postID, cursor, limit+1)
	if err != nil {
		return
	}

	if len(rows) > int(limit) {
		rows = rows[:len(rows)-1]
		next = rows[len(rows)-1].ID
	}

	return
}
