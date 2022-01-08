package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

func (i *Implementation) GetCommentList(ctx context.Context, req *desc.GetCommentListRequest) (*desc.GetCommentListResponse, error) {
	//id, err := i.facade.AddPostVote(ctx, model.PostVote{
	//	AuthorID: req.GetAuthorId(),
	//	PostID:   req.GetPostId(),
	//	Vote:     req.GetVote(),
	//})
	//if err != nil {
	//	return nil, err
	//}

	return &desc.GetCommentListResponse{}, nil
}
