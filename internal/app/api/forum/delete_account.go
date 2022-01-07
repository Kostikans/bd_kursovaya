package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

func (i *Implementation) DeleteAccount(ctx context.Context, req *desc.DeleteAccountRequest) (*desc.DeleteAccountResponse, error) {
	return &desc.DeleteAccountResponse{}, nil
}
