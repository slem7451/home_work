package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		return
	}
	defer r.Close()

	if len(r.File) != 1 {
		return
	}

	data, err := r.File[0].Open()
	if err != nil {
		return
	}

	
	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
}