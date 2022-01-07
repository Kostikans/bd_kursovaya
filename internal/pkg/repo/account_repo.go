package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
)

// AccountRepo - account repo
type AccountRepo struct {
	db *sqlx.DB
}

// NewAccountRepo - returns new accountRepo
func NewAccountRepo(db *sqlx.DB) *AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

func (a *AccountRepo) CreateAccount(ctx context.Context, account model.Account) (id uint64, err error) {
	query := `INSERT INTO account(nickname,avatar,description) VALUES($1,$2,$3) RETURNING id`

	err = a.db.Get(&id, query, account.Nickname, account.Avatar, account.Description)
	return
}
