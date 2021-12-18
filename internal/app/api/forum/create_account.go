package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
)

// CreateAccount - create account
func (i *Implementation) CreateAccount(ctx context.Context, req *desc.CreateAccountRequest) (*desc.CreateAccountResponse, error) {
	return &desc.CreateAccountResponse{
		Id: 3,
	}, nil
}

