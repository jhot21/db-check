package cmd

import "testing"

func TestPostgresConfigRoundTripsSpecialCharacters(t *testing.T) {
	c := dbConfig{
		Host:     "db.example.com",
		Port:     5433,
		Name:     "my db",
		User:     "us:er@x",
		Password: `p@ss/word:1?#'\ "x"`,
		SSLMode:  "disable",
	}

	got, err := postgresConfig(c)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("nil config")
	}
	if got.Host != c.Host || got.Port != uint16(c.Port) {
		t.Errorf("addr = %s:%d", got.Host, got.Port)
	}
	if got.User != c.User || got.Password != c.Password || got.Database != c.Name {
		t.Errorf("credentials/db not preserved: user=%q password=%q db=%q", got.User, got.Password, got.Database)
	}
}

func TestPostgresConfigSSLMode(t *testing.T) {
	base := dbConfig{Host: "db", Port: 5432, User: "u", Name: "n"}

	disable := base
	disable.SSLMode = "disable"
	got, err := postgresConfig(disable)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("nil config")
	}
	if got.TLSConfig != nil {
		t.Error("sslmode=disable should not configure TLS")
	}

	require := base
	require.SSLMode = "require"
	got, err = postgresConfig(require)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("nil config")
	}
	if got.TLSConfig == nil {
		t.Error("sslmode=require should configure TLS")
	}
}

func TestPostgresConfigRejectsInvalidSSLMode(t *testing.T) {
	_, err := postgresConfig(dbConfig{Host: "db", Port: 5432, SSLMode: "bogus"})
	if err == nil {
		t.Fatal("expected error for invalid sslmode")
	}
}
