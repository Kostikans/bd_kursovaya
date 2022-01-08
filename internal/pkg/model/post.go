package model

import "time"

type Post struct {
	ID       uint64 `db:"id"`
	AuthorID uint64 `db:"author_id"`
	Title    string `db:"title"`
	Text     string `db:"text"`
}

type PostVote struct {
	ID       uint64 `db:"id"`
	AuthorID uint64 `db:"author_id"`
	PostID   uint64 `db:"post_id"`
	Vote     bool   `db:"vote"`
}

type ExtendedPost struct {
	ID           uint64    `db:"id"`
	AuthorID     uint64    `db:"author_id"`
	Title        string    `db:"title"`
	Text         string    `db:"text"`
	Tags         []Tag     `db:"tags"`
	LikeCount    uint64    `db:"like_count"`
	DislikeCount uint64    `db:"dislike_count"`
	CreatedAt    time.Time `db:"created_at"`
}
