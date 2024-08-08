package hw10programoptimization

import (
	"archive/zip"
	"io"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	data, zrc := getReader()
	defer zrc.Close()

	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
}

func getReader() (io.ReadCloser, *zip.ReadCloser) {
	r, _ := zip.OpenReader("testdata/users.dat.zip")
	data, _ := r.File[0].Open()
	return data, r
}
