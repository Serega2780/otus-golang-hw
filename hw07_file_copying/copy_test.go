package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	buf, buf2 := make([]byte, 1024*1024), make([]byte, 1024*1024)

	t.Run("offset0_limit0 case", func(t *testing.T) {
		err := Copy(`./testdata/input.txt`, `/tmp/random`, 0, 0)
		require.NoError(t, err, nil)

		f, _ := os.Open(`/tmp/random`)
		f2, _ := os.Open(`./testdata/out_offset0_limit0.txt`)
		defer f.Close()
		defer f2.Close()

		_, _ = f.ReadAt(buf, 0)
		_, _ = f2.ReadAt(buf2, 0)
		require.True(t, testEq(buf, buf2))
	})

	t.Run("offset0_limit10 case", func(t *testing.T) {
		err := Copy(`./testdata/input.txt`, `/tmp/random`, 0, 10)
		require.NoError(t, err, nil)

		f, _ := os.Open(`/tmp/random`)
		f2, _ := os.Open(`./testdata/out_offset0_limit10.txt`)
		defer f.Close()
		defer f2.Close()

		_, _ = f.ReadAt(buf, 0)
		_, _ = f2.ReadAt(buf2, 0)
		require.True(t, testEq(buf, buf2))
	})

	t.Run("offset0_limit1000 case", func(t *testing.T) {
		err := Copy(`./testdata/input.txt`, `/tmp/random`, 0, 1000)
		require.NoError(t, err, nil)

		f, _ := os.Open(`/tmp/random`)
		f2, _ := os.Open(`./testdata/out_offset0_limit1000.txt`)
		defer f.Close()
		defer f2.Close()

		_, _ = f.ReadAt(buf, 0)
		_, _ = f2.ReadAt(buf2, 0)
		require.True(t, testEq(buf, buf2))
	})

	t.Run("offset0_limit10000 case", func(t *testing.T) {
		err := Copy(`./testdata/input.txt`, `/tmp/random`, 0, 10000)
		require.NoError(t, err, nil)

		f, _ := os.Open(`/tmp/random`)
		f2, _ := os.Open(`./testdata/out_offset0_limit10000.txt`)
		defer f.Close()
		defer f2.Close()

		_, _ = f.ReadAt(buf, 0)
		_, _ = f2.ReadAt(buf2, 0)
		require.True(t, testEq(buf, buf2))
	})

	t.Run("offset100_limit1000 case", func(t *testing.T) {
		err := Copy(`./testdata/input.txt`, `/tmp/random`, 100, 1000)
		require.NoError(t, err, nil)

		f, _ := os.Open(`/tmp/random`)
		f2, _ := os.Open(`./testdata/out_offset100_limit1000.txt`)
		defer f.Close()
		defer f2.Close()

		_, _ = f.ReadAt(buf, 0)
		_, _ = f2.ReadAt(buf2, 0)
		require.True(t, testEq(buf, buf2))
	})

	t.Run("offset6000_limit1000 case", func(t *testing.T) {
		err := Copy(`./testdata/input.txt`, `/tmp/random`, 6000, 1000)
		require.NoError(t, err, nil)

		f, _ := os.Open(`/tmp/random`)
		f2, _ := os.Open(`./testdata/out_offset6000_limit1000.txt`)
		defer f.Close()
		defer f2.Close()

		_, _ = f.ReadAt(buf, 0)
		_, _ = f2.ReadAt(buf2, 0)
		require.True(t, testEq(buf, buf2))
	})

	t.Run("input_one_symbol case", func(t *testing.T) {
		err := Copy(`./testdata/input_one_symbol.txt`, `/tmp/random`, 0, 0)
		require.NoError(t, err, nil)

		f, _ := os.Open(`/tmp/random`)
		f2, _ := os.Open(`./testdata/input_one_symbol.txt`)
		defer f.Close()
		defer f2.Close()

		_, _ = f.ReadAt(buf, 0)
		_, _ = f2.ReadAt(buf2, 0)
		require.True(t, testEq(buf, buf2))
	})

	t.Run("offset10000_limit0 case", func(t *testing.T) {
		err := Copy(`./testdata/input.txt`, `/tmp/random`, 10000, 0)
		require.Error(t, err, ErrOffsetExceedsFileSize)
	})
}

func testEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
