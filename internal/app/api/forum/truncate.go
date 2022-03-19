package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

func (i *Implementation) Truncate(ctx context.Context, req *desc.TruncateRequest) (*desc.TruncateResponse, error) {
	err := i.facade.Truncate(ctx)
	if err != nil {
		return nil, err
	}
	return &desc.TruncateResponse{}, nil
}
