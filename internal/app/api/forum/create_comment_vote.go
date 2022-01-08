package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

func (i *Implementation) CreateCommentVote(ctx context.Context, req *desc.CreateCommentVoteRequest) (*desc.CreateCommentVoteResponse, error) {
	id, err := i.facade.AddCommentVote(ctx, model.CommentVote{
		AuthorID:  req.GetAuthorId(),
		PostID:    req.GetPostId(),
		CommentID: req.GetCommentId(),
		Vote:      req.GetVote(),
	})
	if err != nil {
		return nil, err
	}

	return &desc.CreateCommentVoteResponse{
		Id: id,
	}, err
}
