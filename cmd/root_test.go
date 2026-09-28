package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// executeRoot runs the root command with a hermetic environment: the
// developer's real HOST/PORT/USER/... must not leak into viper's AutomaticEnv.
func executeRoot(t *testing.T, args ...string) (string, error) {
	t.Helper()
	for _, k := range []string{"HOST", "PORT", "NAME", "USER", "PASSWORD", "SSLMODE", "TLS"} {
		t.Setenv(k, "") // empty env values are treated as unset by viper
	}

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs(args)
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
	})
	err := rootCmd.Execute()
	return out.String(), err
}

func TestRootWithoutSubcommandFails(t *testing.T) {
	if _, err := executeRoot(t); err == nil {
		t.Fatal("bare invocation must fail so a misconfigured healthcheck is not reported healthy")
	}
}

func TestSubcommandsRejectPositionalArgs(t *testing.T) {
	for _, sub := range []string{"mysql", "postgres"} {
		if _, err := executeRoot(t, sub, "extra"); err == nil {
			t.Errorf("%s: expected error for positional arg", sub)
		}
	}
}

func TestSubcommandRequiresHost(t *testing.T) {
	for _, sub := range []string{"mysql", "postgres"} {
		_, err := executeRoot(t, sub)
		if err == nil || !strings.Contains(err.Error(), "host is required") {
			t.Errorf("%s: err = %v, want \"host is required\"", sub, err)
		}
	}
}

func TestNonCheckCommandsDoNotPrintSuccess(t *testing.T) {
	for _, args := range [][]string{{"help"}, {"--help"}, {"--version"}, {"completion", "bash"}} {
		out, err := executeRoot(t, args...)
		if err != nil {
			t.Errorf("%v: unexpected error: %v", args, err)
		}
		if strings.Contains(out, "Success") {
			t.Errorf("%v: printed Success", args)
		}
	}
}
