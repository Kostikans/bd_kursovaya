package facade

import (
	"context"
	"database/sql"

	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	averageTagsCount = 10
)

type TxManager interface {
	RunTX(ctx context.Context, name string, op func(ctx context.Context) error) error
}

type AccountRepo interface {
	CreateAccount(ctx context.Context, account model.Account) (id uint64, err error)
	Truncate(ctx context.Context) (err error)
}

type CommentRepo interface {
	CreateComment(ctx context.Context, comment model.Comment) (id uint64, err error)
	CheckAuthorAndPostExist(ctx context.Context, comment model.Comment) (exist bool, err error)

	GetCommentVote(ctx context.Context, comment model.CommentVote) (res model.CommentVote, err error)
	CreateCommentVote(ctx context.Context, comment model.CommentVote) (id uint64, err error)
	IncrementCommentVote(ctx context.Context, commentID uint64, likeCount int64, dislikeCount int64) (id uint64, err error)

	GetComments(ctx context.Context, postID uint64, limit uint32, cursor uint64) (rows []model.ExtendedComment, next uint64, err error)
	CreateCommentPartition(ctx context.Context, comment model.Comment) (err error)
}

type PostRepo interface {
	CreatePost(ctx context.Context, post model.Post) (id uint64, err error)
	CheckAuthorExist(ctx context.Context, post model.Post) (exist bool, err error)
	CheckPostExist(ctx context.Context, postID uint64) (exist bool, err error)

	GetPostVote(ctx context.Context, post model.PostVote) (res model.PostVote, err error)
	CreatePostVote(ctx context.Context, post model.PostVote) (id uint64, err error)
	IncrementPostVote(ctx context.Context, postID uint64, likeCount int64, dislikeCount int64) (id uint64, err error)

	GetPosts(ctx context.Context, limit uint32, cursor uint64) (rows []model.ExtendedPost, next uint64, err error)
}

type TagRepo interface {
	CreateTag(ctx context.Context, tag model.Tag) (id uint64, err error)
	CheckAuthorExist(ctx context.Context, tag model.Tag) (exist bool, err error)
	BulkUpdatePostTags(ctx context.Context, ids []uint64, postID uint64) (err error)
	CheckPostAndTagsExist(ctx context.Context, ids []uint64, postID uint64) (exist bool, err error)
	UpdatePostTags(ctx context.Context, ids []uint64, postID uint64) (err error)
	GetTagsByIDs(ctx context.Context, ids []uint64) (rows []model.Tag, err error)
}

type Facade struct {
	accountRepo AccountRepo
	commentRepo CommentRepo
	postRepo    PostRepo
	tagRepo     TagRepo
	txManager   TxManager
}

type Opts struct {
	AccountRepo AccountRepo
	CommentRepo CommentRepo
	PostRepo    PostRepo
	TagRepo     TagRepo
	TxManager   TxManager
}

func NewFacade(opts Opts) *Facade {
	return &Facade{
		accountRepo: opts.AccountRepo,
		commentRepo: opts.CommentRepo,
		postRepo:    opts.PostRepo,
		tagRepo:     opts.TagRepo,
		txManager:   opts.TxManager,
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

	//err = f.commentRepo.CreateCommentPartition(ctx, comment)
	//fmt.Println(err)
	//if err != nil {
	//	return
	//}

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
	err = f.txManager.RunTX(ctx, "update post tags", func(ctx context.Context) (err error) {
		exist, err := f.tagRepo.CheckPostAndTagsExist(ctx, tagIDs, postID)
		if err != nil {
			return
		}
		if !exist {
			return status.Errorf(codes.NotFound, "post or tags not found")
		}

		err = f.tagRepo.BulkUpdatePostTags(ctx, tagIDs, postID)
		if err != nil {
			return
		}

		return f.tagRepo.UpdatePostTags(ctx, tagIDs, postID)
	})

	return
}

func (f *Facade) AddPostVote(ctx context.Context, postVote model.PostVote) (id uint64, err error) {
	err = f.txManager.RunTX(ctx, "add post vote", func(ctx context.Context) (err error) {
		var (
			hasPrev      = true
			prevVote     model.PostVote
			likeCount    = int64(0)
			dislikeCount = int64(0)
		)

		prevVote, err = f.postRepo.GetPostVote(ctx, postVote)
		if err != nil {
			if err == sql.ErrNoRows {
				hasPrev = false
			} else {
				return err
			}
		}

		if postVote.Vote {
			likeCount++
		} else {
			dislikeCount++
		}

		if hasPrev {
			if prevVote.Vote != postVote.Vote {
				if prevVote.Vote && !postVote.Vote {
					likeCount = -1
					dislikeCount = 1
				} else {
					likeCount = 1
					dislikeCount = -1
				}
			} else {
				likeCount = 0
				dislikeCount = 0
			}
		}

		_, err = f.postRepo.CreatePostVote(ctx, postVote)
		if err != nil {
			return err
		}

		_, err = f.postRepo.IncrementPostVote(ctx, postVote.PostID, likeCount, dislikeCount)
		return err
	})

	return 0, err
}

func (f *Facade) AddCommentVote(ctx context.Context, commentVote model.CommentVote) (id uint64, err error) {
	err = f.txManager.RunTX(ctx, "add post vote", func(ctx context.Context) (err error) {
		var (
			hasPrev      = true
			prevVote     model.CommentVote
			likeCount    = int64(0)
			dislikeCount = int64(0)
		)

		prevVote, err = f.commentRepo.GetCommentVote(ctx, commentVote)
		if err != nil {
			if err == sql.ErrNoRows {
				hasPrev = false
			} else {
				return err
			}
		}

		if commentVote.Vote {
			likeCount++
		} else {
			dislikeCount++
		}

		if hasPrev {
			if prevVote.Vote != commentVote.Vote {
				if prevVote.Vote && !commentVote.Vote {
					likeCount = -1
					dislikeCount = 1
				} else {
					likeCount = 1
					dislikeCount = -1
				}
			} else {
				likeCount = 0
				dislikeCount = 0
			}
		}

		_, err = f.commentRepo.CreateCommentVote(ctx, commentVote)
		if err != nil {
			return err
		}

		_, err = f.commentRepo.IncrementCommentVote(ctx, commentVote.CommentID, likeCount, dislikeCount)
		return err
	})

	return 0, err
}

func (f *Facade) GetPosts(ctx context.Context, limit uint32, cursor uint64) (res []model.ExtendedPost, next uint64, err error) {
	res, next, err = f.postRepo.GetPosts(ctx, limit, cursor)
	if err != nil {
		return res, next, err
	}

	tagsIDs := make([]uint64, 0, len(res)*averageTagsCount)
	for _, post := range res {
		for _, tagID := range post.TagsIDs {
			tagsIDs = append(tagsIDs, uint64(tagID))
		}
	}

	tags, err := f.tagRepo.GetTagsByIDs(ctx, tagsIDs)

	tagMap := make(map[uint64]model.Tag, len(tags))
	for _, tag := range tags {
		tagMap[tag.ID] = tag
	}

	for i := range res {
		res[i].Tags = make([]model.Tag, 0, len(res[i].TagsIDs))
		for _, tagID := range res[i].TagsIDs {
			tag, ok := tagMap[uint64(tagID)]
			if !ok {
				continue
			}

			res[i].Tags = append(res[i].Tags, tag)
		}
	}

	return
}

func (f *Facade) GetComments(ctx context.Context, postID uint64, limit uint32, cursor uint64) (res []model.ExtendedComment, next uint64, err error) {
	res, next, err = f.commentRepo.GetComments(ctx, postID, limit, cursor)
	if err != nil {
		return res, next, err
	}

	return
}

func (f *Facade) Truncate(ctx context.Context) (err error) {
	return f.accountRepo.Truncate(ctx)
}
