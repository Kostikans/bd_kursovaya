package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

type Facade interface {
	CreateAccount(ctx context.Context, account model.Account) (id uint64, err error)
	CreateComment(ctx context.Context, comment model.Comment) (id uint64, err error)
	CreatePost(ctx context.Context, post model.Post) (id uint64, err error)
	CreateTag(ctx context.Context, tag model.Tag) (id uint64, err error)
	UpdatePostTags(ctx context.Context, tagIDs []uint64, postID uint64) (err error)
}

type Implementation struct {
	desc.UnimplementedForumServer
	facade Facade
}

// Opts - configuration service dependencies
type Opts struct {
	Facade Facade
}

// NewForum return new instance of Implementation.
func NewForum(opts Opts) *Implementation {
	return &Implementation{
		facade: opts.Facade,
	}
}
