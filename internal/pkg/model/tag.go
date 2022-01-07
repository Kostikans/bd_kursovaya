package model

type Tag struct {
	ID       uint64 `db:"id"`
	AuthorID uint64 `db:"author_id"`
	Name     string `db:"name"`
}

type PostTag struct {
	PostID   uint64 `db:"post_id"`
	AuthorID string `db:"author_id"`
}
