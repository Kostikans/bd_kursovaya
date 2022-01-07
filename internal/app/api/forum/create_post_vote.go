package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

func (i *Implementation) CreatePostVote(ctx context.Context, req *desc.CreatePostVoteRequest) (*desc.CreatePostVoteResponse, error) {
	return &desc.CreatePostVoteResponse{}, nil
}
