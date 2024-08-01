package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	t.Run("done ReadDir case", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		require.Nil(t, err)
		require.Equal(t, 5, len(env))
	})
}
