package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) GetComments(ctx context.Context, req *desc.GetCommentListRequest) (*desc.GetCommentListResponse, error) {
	comments, next, err := i.facade.GetComments(ctx, req.GetPostId(), req.GetLimit(), req.GetCursor())
	if err != nil {
		return nil, err
	}

	items := make([]*desc.Comment, 0, len(comments))
	for _, comment := range comments {
		items = append(items, &desc.Comment{
			PostId:       comment.PostID,
			AuthorId:     comment.AuthorID,
			ParentId:     comment.ParentID,
			Id:           comment.ID,
			Text:         comment.Text,
			LikeCount:    comment.LikeCount,
			DislikeCount: comment.DislikeCount,
			CreatedAt:    timestamppb.New(comment.CreatedAt),
			Depth:        uint64(len(comment.BreadCrumbs)),
		})
	}

	return &desc.GetCommentListResponse{
		Comments: items,
		Next:     next,
		HasNext:  next != 0,
	}, err
}
