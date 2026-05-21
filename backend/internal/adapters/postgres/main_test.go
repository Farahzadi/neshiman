package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var globalPool *pgxpool.Pool

func TestMain(m *testing.M) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/neshiman_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to test database: %v\n", err)
		fmt.Fprintf(os.Stderr, "skipping integration tests (set DB_URL to a test database)\n")
		os.Exit(0)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping test database: %v\n", err)
		fmt.Fprintf(os.Stderr, "skipping integration tests\n")
		os.Exit(0)
	}

	globalPool = pool
	os.Exit(m.Run())
}

func testPool(tb testing.TB) *pgxpool.Pool {
	tb.Helper()
	if globalPool == nil {
		tb.Skip("no test database available (set DB_URL)")
	}
	return globalPool
}

func truncate(tb testing.TB, pool *pgxpool.Pool) {
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
		if _, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			tb.Fatalf("failed to truncate %s: %v", table, err)
		}
	}
}
