package facade

import (
	"context"

	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountRepo interface {
	CreateAccount(ctx context.Context, account model.Account) (id uint64, err error)
}

type CommentRepo interface {
	CreateComment(ctx context.Context, comment model.Comment) (id uint64, err error)
	CheckAuthorAndPostExist(ctx context.Context, comment model.Comment) (exist bool, err error)
}

type PostRepo interface {
	CreatePost(ctx context.Context, post model.Post) (id uint64, err error)
	CheckAuthorExist(ctx context.Context, post model.Post) (exist bool, err error)
	CheckPostExist(ctx context.Context, postID uint64) (exist bool, err error)
}

type TagRepo interface {
	CreateTag(ctx context.Context, tag model.Tag) (id uint64, err error)
	CheckAuthorExist(ctx context.Context, tag model.Tag) (exist bool, err error)
	BulkUpdatePostTags(ctx context.Context, ids []uint64, postID uint64) (err error)
	CheckPostAndTagsExist(ctx context.Context, ids []uint64, postID uint64) (exist bool, err error)
}

type Facade struct {
	accountRepo AccountRepo
	commentRepo CommentRepo
	postRepo    PostRepo
	tagRepo     TagRepo
}

type Opts struct {
	AccountRepo AccountRepo
	CommentRepo CommentRepo
	PostRepo    PostRepo
	TagRepo     TagRepo
}

func NewFacade(opts Opts) *Facade {
	return &Facade{
		accountRepo: opts.AccountRepo,
		commentRepo: opts.CommentRepo,
		postRepo:    opts.PostRepo,
		tagRepo:     opts.TagRepo,
	}
}

func (f *Facade) CreateAccount(ctx context.Context, account model.Account) (id uint64, err error) {
	return f.accountRepo.CreateAccount(ctx, account)
}

func (f *Facade) CreateComment(ctx context.Context, comment model.Comment) (id uint64, err error) {
	exist, err := f.commentRepo.CheckAuthorAndPostExist(ctx, comment)
	if err != nil {
		return
	}
	if !exist {
		return 0, status.Errorf(codes.NotFound, "author or post not found")
	}

	return f.commentRepo.CreateComment(ctx, comment)
}

func (f *Facade) CreatePost(ctx context.Context, post model.Post) (id uint64, err error) {
	exist, err := f.postRepo.CheckAuthorExist(ctx, post)
	if err != nil {
		return
	}
	if !exist {
		return 0, status.Errorf(codes.NotFound, "author not found")
	}

	return f.postRepo.CreatePost(ctx, post)
}

func (f *Facade) CreateTag(ctx context.Context, tag model.Tag) (id uint64, err error) {
	exist, err := f.tagRepo.CheckAuthorExist(ctx, tag)
	if err != nil {
		return
	}
	if !exist {
		return 0, status.Errorf(codes.NotFound, "author not found")
	}

	return f.tagRepo.CreateTag(ctx, tag)
}

func (f *Facade) UpdatePostTags(ctx context.Context, tagIDs []uint64, postID uint64) (err error) {
	exist, err := f.tagRepo.CheckPostAndTagsExist(ctx, tagIDs, postID)
	if err != nil {
		return
	}
	if !exist {
		return status.Errorf(codes.NotFound, "post or tags not found")
	}

	return f.tagRepo.BulkUpdatePostTags(ctx, tagIDs, postID)
}
