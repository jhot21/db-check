package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(postgresCmd)
}

var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "check connection to postgres database",
	RunE: func(cmd *cobra.Command, args []string) error {
		timeoutContext, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
		defer cancel()

		conn, err := pgx.Connect(timeoutContext, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=allow",
			viper.GetString("user"),
			viper.GetString("password"),
			viper.GetString("host"),
			viper.GetString("port"),
			viper.GetString("name")))
		if err != nil {
			return fmt.Errorf("error connecting to database: %w", err)
		}
		defer conn.Close(timeoutContext)

		var stuff string
		return conn.QueryRow(timeoutContext, "SELECT 1").Scan(&stuff)
	},
}
