package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	t.Run("done RunCmd ping case", func(t *testing.T) {
		code := RunCmd([]string{"echo", "test"}, Environment{})
		require.Equal(t, 0, code)
	})

	t.Run("done RunCmd case", func(t *testing.T) {
		env := Environment{}
		os.Setenv("UNSET", "HAS_TO_BE_REMOVED")
		require.Equal(t, "HAS_TO_BE_REMOVED", os.Getenv("UNSET"))
		env["BAR"] = EnvValue{Value: "TEST_BAR", NeedRemove: false}
		env["UNSET"] = EnvValue{Value: "", NeedRemove: true}
		code := RunCmd([]string{"/bin/bash", "./testdata/echo.sh"}, env)
		require.Equal(t, 0, code)
		require.Equal(t, "", os.Getenv("UNSET"))
	})
}
