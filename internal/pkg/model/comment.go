package model

import (
	"time"

	"github.com/lib/pq"
)

type Comment struct {
	CommentID uint64 `db:"id"`
	AuthorID  uint64 `db:"author_id"`
	PostID    uint64 `db:"post_id"`
	ParentID  uint64 `db:"parent_id"`
	Text      string `db:"text"`
}

type CommentVote struct {
	ID        uint64 `db:"id"`
	CommentID uint64 `db:"comment_id"`
	AuthorID  uint64 `db:"author_id"`
	PostID    uint64 `db:"post_id"`
	Vote      bool   `db:"vote"`
}

type ExtendedComment struct {
	ID           uint64        `db:"id"`
	AuthorID     uint64        `db:"author_id"`
	PostID       uint64        `db:"post_id"`
	ParentID     uint64        `db:"parent_id"`
	Text         string        `db:"text"`
	LikeCount    uint64        `db:"like_count"`
	DislikeCount uint64        `db:"dislike_count"`
	CreatedAt    time.Time     `db:"created_at"`
	BreadCrumbs  pq.Int64Array `db:"breadcrumbs"`
}
