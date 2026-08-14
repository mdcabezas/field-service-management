package postgresshared

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func testPool() *pgxpool.Pool {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://fsm_admin:change_me_in_prod@localhost:5432/fsm_test_gas?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	return pool
}

func newPool(t *testing.T) *pgxpool.Pool {
	pool := testPool()
	t.Cleanup(func() { pool.Close() })
	return pool
}
