package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

func (i *Implementation) CreateCommentVote(ctx context.Context, req *desc.CreateCommentVoteRequest) (*desc.CreateCommentVoteResponse, error) {
	return &desc.CreateCommentVoteResponse{}, nil
}
