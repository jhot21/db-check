package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultMySQLPort = 3306

func init() {
	mysqlCmd.Flags().String("tls", "false", "TLS mode: true, false, skip-verify or preferred")
	viper.BindPFlag("tls", mysqlCmd.Flags().Lookup("tls"))
	rootCmd.AddCommand(mysqlCmd)
}

// mysqlConfig builds the driver config from discrete fields. It is used with
// mysql.NewConnector so no DSN string is ever formatted or parsed, which keeps
// special characters in credentials intact.
func mysqlConfig(c dbConfig) *mysql.Config {
	cfg := mysql.NewConfig()
	cfg.User = c.User
	cfg.Passwd = c.Password
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	cfg.DBName = c.Name
	cfg.Timeout = 10 * time.Second
	cfg.TLSConfig = c.TLS
	return cfg
}

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "check connection to mysql/mariadb database",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newConfig(viper.GetViper(), defaultMySQLPort)
		if err != nil {
			return err
		}

		connector, err := mysql.NewConnector(mysqlConfig(c))
		if err != nil {
			return fmt.Errorf("invalid connection settings: %w", err)
		}
		db := sql.OpenDB(connector)
		defer db.Close()

		ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			return err
		}
		reportSuccess(cmd)
		return nil
	},
}
