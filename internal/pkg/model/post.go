package model

type Post struct {
	ID       uint64 `db:"id"`
	AuthorID uint64 `db:"author_id"`
	Title    string `db:"title"`
	Text     string `db:"text"`
}
