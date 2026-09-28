package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestNewConfigAppliesDefaultPort(t *testing.T) {
	v := viper.New()
	v.Set("host", "db")

	c, err := newConfig(v, 5432)
	if err != nil {
		t.Fatal(err)
	}
	if c.Port != 5432 {
		t.Errorf("port = %d, want 5432", c.Port)
	}
}

func TestNewConfigExplicitPortWins(t *testing.T) {
	v := viper.New()
	v.Set("host", "db")
	v.Set("port", "6543") // env vars arrive as strings

	c, err := newConfig(v, 5432)
	if err != nil {
		t.Fatal(err)
	}
	if c.Port != 6543 {
		t.Errorf("port = %d, want 6543", c.Port)
	}
}

func TestNewConfigRequiresHost(t *testing.T) {
	_, err := newConfig(viper.New(), 3306)
	if err == nil || !strings.Contains(err.Error(), "host") {
		t.Fatalf("err = %v, want error mentioning host", err)
	}
}

func TestNewConfigRejectsOutOfRangePort(t *testing.T) {
	for _, port := range []int{-1, 65536} {
		v := viper.New()
		v.Set("host", "db")
		v.Set("port", port)
		if _, err := newConfig(v, 3306); err == nil {
			t.Errorf("port %d: expected error", port)
		}
	}
}
