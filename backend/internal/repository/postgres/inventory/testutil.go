package postgresinventory

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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

func strPtr(s string) *string        { return &s }
func intPtr(i int) *int              { return &i }
func floatPtr(f float64) *float64    { return &f }
func uuidPtr(u uuid.UUID) *uuid.UUID { return &u }
func timePtr(t time.Time) *time.Time { return &t }
