package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

func (i *Implementation) CreatePost(ctx context.Context, req *desc.CreatePostRequest) (*desc.CreatePostResponse, error) {
	id, err := i.facade.CreatePost(ctx, model.Post{
		AuthorID: req.GetAuthorId(),
		Title:    req.GetTitle(),
		Text:     req.GetText(),
	})
	if err != nil {
		return nil, err
	}

	return &desc.CreatePostResponse{
		Id: id,
	}, err
}
