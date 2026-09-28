package cmd

import (
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

func TestMySQLConfigPreservesSpecialCharacters(t *testing.T) {
	c := dbConfig{
		Host:     "db.example.com",
		Port:     3307,
		Name:     "my/db",
		User:     "us:er@x",
		Password: "p@ss/word:1?#&=",
		TLS:      "false",
	}

	got := mysqlConfig(c)
	if got.User != c.User || got.Passwd != c.Password || got.DBName != c.Name {
		t.Errorf("credentials/db not preserved: user=%q password=%q db=%q", got.User, got.Passwd, got.DBName)
	}
	if got.Net != "tcp" || got.Addr != "db.example.com:3307" {
		t.Errorf("addr = %s://%s", got.Net, got.Addr)
	}
}

func TestMySQLConfigIPv6Host(t *testing.T) {
	got := mysqlConfig(dbConfig{Host: "::1", Port: 3306})
	if got.Addr != "[::1]:3306" {
		t.Errorf("addr = %s, want [::1]:3306", got.Addr)
	}
}

func TestMySQLConfigSetsTimeoutAndTLS(t *testing.T) {
	got := mysqlConfig(dbConfig{Host: "db", Port: 3306, TLS: "skip-verify"})
	if got.Timeout != 10*time.Second {
		t.Errorf("timeout = %v, want 10s", got.Timeout)
	}
	if got.TLSConfig != "skip-verify" {
		t.Errorf("tls = %q, want skip-verify", got.TLSConfig)
	}
}

func TestMySQLConfigAcceptedByConnector(t *testing.T) {
	for _, tls := range []string{"true", "false", "skip-verify", "preferred"} {
		if _, err := mysql.NewConnector(mysqlConfig(dbConfig{Host: "db", Port: 3306, TLS: tls})); err != nil {
			t.Errorf("tls=%s: %v", tls, err)
		}
	}
	if _, err := mysql.NewConnector(mysqlConfig(dbConfig{Host: "db", Port: 3306, TLS: "bogus"})); err == nil {
		t.Error("tls=bogus: expected error")
	}
}
