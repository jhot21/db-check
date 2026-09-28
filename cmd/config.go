package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// dbConfig holds the connection settings shared by all subcommands.
type dbConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string // postgres only
	TLS      string // mysql only
}

// newConfig reads and validates connection settings from v. defaultPort is
// used when no port was provided.
func newConfig(v *viper.Viper, defaultPort int) (dbConfig, error) {
	c := dbConfig{
		Host:     v.GetString("host"),
		Port:     v.GetInt("port"),
		Name:     v.GetString("name"),
		User:     v.GetString("user"),
		Password: v.GetString("password"),
		SSLMode:  v.GetString("sslmode"),
		TLS:      v.GetString("tls"),
	}

	if c.Host == "" {
		return dbConfig{}, errors.New("host is required (--host or HOST)")
	}
	if c.Port == 0 {
		c.Port = defaultPort
	}
	if c.Port < 1 || c.Port > 65535 {
		return dbConfig{}, fmt.Errorf("port %d is out of range (1-65535)", c.Port)
	}
	return c, nil
}
