package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

func (i *Implementation) CreateComment(ctx context.Context, req *desc.CreateCommentRequest) (*desc.CreateCommentResponse, error) {
	id, err := i.facade.CreateComment(ctx, model.Comment{
		AuthorID: req.GetAuthorId(),
		PostID:   req.GetPostId(),
		ParentID: req.GetParentId(),
		Text:     req.GetText(),
	})
	if err != nil {
		return nil, err
	}

	return &desc.CreateCommentResponse{
		Id: id,
	}, err
}
