package model

type Account struct {
	AccountID   uint64 `db:"id"`
	Nickname    string `db:"nickname"`
	Avatar      string `db:"avatar"`
	Description string `db:"description"`
}
