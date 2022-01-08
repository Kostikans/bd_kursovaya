package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

func (i *Implementation) CreatePostVote(ctx context.Context, req *desc.CreatePostVoteRequest) (*desc.CreatePostVoteResponse, error) {
	id, err := i.facade.AddPostVote(ctx, model.PostVote{
		AuthorID: req.GetAuthorId(),
		PostID:   req.GetPostId(),
		Vote:     req.GetVote(),
	})
	if err != nil {
		return nil, err
	}

	return &desc.CreatePostVoteResponse{
		Id: id,
	}, err
}
