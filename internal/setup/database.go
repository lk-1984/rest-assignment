package setup

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func initializeDatabase(cfg *Config) error {
	pgConnectionurl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.PgUsername.Value,
		cfg.PgPassword.Value,
		cfg.PgHostname.Value,
		cfg.PgPort.Value,
		cfg.PgDatabase.Value,
	)

	fmt.Println(pgConnectionurl)

	pool, err := pgxpool.New(context.Background(), pgConnectionurl)
	if err != nil {
		return err
	}
	cfg.PgPool = &PgxPoolWrapper{Pool: pool}
	return nil
}
