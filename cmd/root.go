package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "db-check",
		Short: "A command line tool to check database connectivity",
		Long: `DB Check is a utility meant to be used as a dependency of containers that rely
on a remote database. This will allow you to wait to start those containers until the database is available.`,
	}
)

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
	// viper.ReadInConfig()
}
