package cmd

import (
	"fmt"

	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(mysqlCmd)
}

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "check connection to mysql/mariadb database",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=10s",
			viper.GetString("user"),
			viper.GetString("password"),
			viper.GetString("host"),
			viper.GetString("port"),
			viper.GetString("name")))

		if err != nil {
			return fmt.Errorf("error connecting to database: %w", err)
		}

		defer db.Close()
		db.SetConnMaxLifetime(time.Minute * 1)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)

		stmtOut, err := db.Prepare("SELECT 1")
		if err != nil {
			return fmt.Errorf("error querying database: %w", err)
		}
		defer stmtOut.Close()

		var stuff string
		return stmtOut.QueryRow().Scan(&stuff)
	},
}
