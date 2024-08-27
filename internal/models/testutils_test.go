package models

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test_web:pass@localhost:5432/test_snippetbox?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	if err = pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, string(script))
	if err != nil {
		t.Fatal(err)
	}
	// Use the t.Cleanup() to register a function which will automatically be
	// called by Go when the current test (or sub-test) which calls newTestDB() has finished.
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = pool.Exec(ctx, string(script))
		if err != nil {
			t.Fatal(err)
		}

		pool.Close()
	})

	return pool
}
