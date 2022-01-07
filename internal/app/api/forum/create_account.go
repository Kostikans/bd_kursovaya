package forum

import (
	"context"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

// CreateAccount - create account
func (i *Implementation) CreateAccount(ctx context.Context, req *desc.CreateAccountRequest) (*desc.CreateAccountResponse, error) {
	id, err := i.facade.CreateAccount(ctx, model.Account{
		Nickname:    req.GetAccount().GetNickname(),
		Avatar:      req.GetAccount().GetAvatar(),
		Description: req.GetAccount().GetDescription(),
	})
	if err != nil {
		return nil, err
	}

	return &desc.CreateAccountResponse{
		Id: id,
	}, err
}
