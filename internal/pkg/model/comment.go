package model

type Comment struct {
	CommentID uint64 `db:"id"`
	AuthorID  uint64 `db:"author_id"`
	PostID    uint64 `db:"post_id"`
	ParentID  uint64 `db:"parent_id"`
	Text      string `db:"text"`
}
