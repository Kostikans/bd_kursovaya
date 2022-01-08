package repo

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	registerDBOnce      = &sync.Once{}
	basePaginationLimit = uint32(1000)
)

// NewTestDB - returns new test database connection
func NewTestDB(t testing.TB) *sql.DB {
	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		dsn = "adv-shopping-user:adv-shopping-api@localhost:5432/adv-shopping-api?sslmode=disable"
		dsn = "postgres://" + dsn // make secret search happy
	}

	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("txdb", dsn)
	require.NoError(t, err)

	return db
}

//func waitDBReady(t testing.TB, db sql.Balancer) {
//	beginTs := time.Now()
//	for {
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//		err := db.Ping(ctx)
//		cancel()
//		if err == nil {
//			return
//		}
//		if time.Since(beginTs) > time.Minute {
//			t.Fatalf("wait db ready, last error: %v", err)
//		}
//	}
//}
