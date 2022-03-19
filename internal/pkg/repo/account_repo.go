package repo

import (
	"context"

	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
	"github.com/kostikan/bd_kursovaya/internal/pkg/sql"
)

// AccountRepo - account repo
type AccountRepo struct {
	db *sql.Balancer
}

// NewAccountRepo - returns new accountRepo
func NewAccountRepo(db *sql.Balancer) *AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

func (a *AccountRepo) CreateAccount(ctx context.Context, account model.Account) (id uint64, err error) {
	query := `INSERT INTO account(nickname,avatar,description) VALUES($1,$2,$3) RETURNING id`

	err = a.db.Write(ctx).Get(&id, query, account.Nickname, account.Avatar, account.Description)
	return
}

func (a *AccountRepo) Truncate(ctx context.Context) (err error) {
	query := `SELECT truncate_tables('bd_kursovaya');`

	_, err = a.db.Write(ctx).Exec(query)
	return
}
