package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// version is set at build time via -ldflags "-X github.com/jhot21/db-check/cmd.version=...".
var version = "dev"

var (
	rootCmd = &cobra.Command{
		Use:   "dbcheck",
		Short: "A command line tool to check database connectivity",
		Long: `DB Check is a utility meant to be used as a dependency of containers that rely
on a remote database. This will allow you to wait to start those containers until the database is available.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		// Without a subcommand there is nothing to check; fail so a
		// misconfigured healthcheck is not reported as healthy.
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Usage()
			return errors.New("a subcommand is required (mysql or postgres)")
		},
	}
)

// reportSuccess is called by check subcommands once the database responded.
// It lives here, not in a PersistentPostRun, so cobra's built-in help and
// completion commands never claim success.
func reportSuccess(cmd *cobra.Command) {
	fmt.Fprintln(cmd.OutOrStdout(), "Success")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("host", "", "database host")
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	rootCmd.PersistentFlags().IntP("port", "p", 0, "database port")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	rootCmd.PersistentFlags().StringP("name", "n", "", "database name")
	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))
	rootCmd.PersistentFlags().StringP("user", "u", "", "database username")
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	rootCmd.PersistentFlags().String("password", "", "database password")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
}

func initConfig() {
	viper.AutomaticEnv()
}
