package repositories

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
	"testFIO/pkg/client/postgresql"
)

func TestPGStore(t *testing.T, cfg config.ServerConfig) (*PGSStore, func(...string)) {
	t.Helper()

	cfg.DatabaseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	logger := loggers.NewLogger()
	client, err := postgresql.NewClient(context.Background(), 5, &cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewPGSStore(client, &cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			_, err = s.client.Exec(context.Background(), fmt.Sprintf(
				"TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
