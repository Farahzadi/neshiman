package testhelper

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func MustNewPool(tb testing.TB) *pgxpool.Pool {
	tb.Helper()

	ctx := context.Background()
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/neshiman_test?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		tb.Fatalf("failed to connect to test database: %v", err)
	}
	tb.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		tb.Fatalf("failed to ping test database: %v", err)
	}

	return pool
}

func TruncateAll(tb testing.TB, pool *pgxpool.Pool) {
	tb.Helper()
	ctx := context.Background()

	tables := []string{
		"cross_team_requests",
		"reservations",
		"seats",
		"users",
		"rooms",
		"teams",
	}

	for _, table := range tables {
		_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			tb.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}
