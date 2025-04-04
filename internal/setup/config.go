package setup

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigItem struct {
	Name  string
	Value string
}

type Config struct {
	PgHostname ConfigItem
	PgPort     ConfigItem
	PgDatabase ConfigItem
	PgUsername ConfigItem
	PgPassword ConfigItem
	PgPool     DBPool
	GinEngine  *gin.Engine
}

var cfg = Config{}

type PgxPoolWrapper struct {
	Pool *pgxpool.Pool
}

func (p *PgxPoolWrapper) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.Pool.Exec(ctx, sql, arguments...)
}

func (p *PgxPoolWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}

func (p *PgxPoolWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.Pool.Query(ctx, sql, args...)
}

func GetConfig() *Config {
	return &cfg
}

func InitializeConfig() error {
	cfg := GetConfig()
	cfg.PgHostname = ConfigItem{}
	cfg.PgPort = ConfigItem{}
	cfg.PgDatabase = ConfigItem{}
	cfg.PgUsername = ConfigItem{}
	cfg.PgPassword = ConfigItem{}

	getEnvsErr := getEnvs(cfg)
	if getEnvsErr != nil {
		return getEnvsErr
	}

	initDbErr := initializeDatabase(cfg)
	if initDbErr != nil {
		return initDbErr
	}

	return nil
}
