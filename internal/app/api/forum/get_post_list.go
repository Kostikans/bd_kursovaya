package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) GetPosts(ctx context.Context, req *desc.GetPostListRequest) (*desc.GetPostListResponse, error) {
	posts, next, err := i.facade.GetPosts(ctx, req.GetLimit(), req.GetCursor())
	if err != nil {
		return nil, err
	}

	items := make([]*desc.Post, 0, len(posts))
	for _, post := range posts {
		descTags := make([]*desc.Tag, 0, len(post.Tags))
		for _, tag := range post.Tags {
			descTags = append(descTags, &desc.Tag{
				Id:       tag.ID,
				AuthorId: tag.AuthorID,
				Name:     tag.Name,
			})
		}

		items = append(items, &desc.Post{
			Title:        post.Title,
			AuthorId:     post.AuthorID,
			Id:           post.ID,
			Text:         post.Text,
			LikeCount:    post.LikeCount,
			DislikeCount: post.DislikeCount,
			CreatedAt:    timestamppb.New(post.CreatedAt),
			Tags:         descTags,
		})
	}

	return &desc.GetPostListResponse{
		Posts:   items,
		Next:    next,
		HasNext: next != 0,
	}, err
}
