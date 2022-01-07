package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

func (i *Implementation) CreateTag(ctx context.Context, req *desc.CreateTagRequest) (*desc.CreateTagResponse, error) {
	id, err := i.facade.CreateTag(ctx, model.Tag{
		AuthorID: req.GetAuthorId(),
		Name:     req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	return &desc.CreateTagResponse{
		Id: id,
	}, err
}
