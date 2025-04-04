package setup

import (
	"fmt"
	"os"
)

func getEnvs(cfg *Config) error {
	cfg.PgHostname.Name = "PG_HOSTNAME"
	cfg.PgPort.Name = "PG_PORT"
	cfg.PgDatabase.Name = "PG_DATABASE"
	cfg.PgUsername.Name = "PG_USERNAME"
	cfg.PgPassword.Name = "PG_PASSWORD"

	configs := [5]*ConfigItem{
		&cfg.PgHostname,
		&cfg.PgPort,
		&cfg.PgDatabase,
		&cfg.PgUsername,
		&cfg.PgPassword,
	}

	for _, c := range configs {
		c.Value = os.Getenv(c.Name)
		if c.Value == "" {
			return fmt.Errorf("environment variable %s is not defined or it is empty", c.Name)
		}
	}

	return nil
}
