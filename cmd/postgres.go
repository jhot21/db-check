package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultPostgresPort = 5432

func init() {
	postgresCmd.Flags().String("sslmode", "allow", "SSL mode: disable, allow, prefer, require, verify-ca or verify-full")
	viper.BindPFlag("sslmode", postgresCmd.Flags().Lookup("sslmode"))
	rootCmd.AddCommand(postgresCmd)
}

// pgQuote quotes a value for a libpq keyword/value connection string.
func pgQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return "'" + s + "'"
}

// postgresConfig builds a pgx config using the keyword/value connection string
// form, so values containing URL-special characters are handled safely.
func postgresConfig(c dbConfig) (*pgx.ConnConfig, error) {
	kv := [][2]string{
		{"host", c.Host},
		{"port", strconv.Itoa(c.Port)},
		{"user", c.User},
		{"password", c.Password},
		{"dbname", c.Name},
		{"sslmode", c.SSLMode},
	}
	parts := make([]string, 0, len(kv))
	for _, p := range kv {
		if p[1] != "" {
			parts = append(parts, p[0]+"="+pgQuote(p[1]))
		}
	}
	return pgx.ParseConfig(strings.Join(parts, " "))
}

var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "check connection to postgres database",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newConfig(viper.GetViper(), defaultPostgresPort)
		if err != nil {
			return err
		}
		cfg, err := postgresConfig(c)
		if err != nil {
			return fmt.Errorf("invalid connection settings: %w", err)
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
		defer cancel()

		// Connect performs the full startup handshake (including authentication),
		// which is the connectivity check; no separate Ping is needed.
		conn, err := pgx.ConnectConfig(ctx, cfg)
		if err != nil {
			return fmt.Errorf("error connecting to database: %w", err)
		}
		// ctx may already be near its deadline; give Close its own.
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer closeCancel()
		defer conn.Close(closeCtx)

		reportSuccess(cmd)
		return nil
	},
}
