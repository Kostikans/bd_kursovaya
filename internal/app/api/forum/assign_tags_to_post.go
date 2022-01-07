package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

func (i *Implementation) AssignTagsToPost(ctx context.Context, req *desc.AssignTagsToPostRequest) (*desc.AssignTagsToPostResponse, error) {
	err := i.facade.UpdatePostTags(ctx, req.GetTagsId(), req.GetPostId())
	if err != nil {
		return nil, err
	}

	return &desc.AssignTagsToPostResponse{}, err
}
