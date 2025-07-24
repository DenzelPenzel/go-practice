package utils

import (
	"log"
	"math"
	"os"
	"syscall"
)

//gcassert:inline
func Round(value float64) float64 {
	return math.Round(value*10.0) / 10.0
}

func Mmap(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return mmapFile(f)
}

func mmapFile(f *os.File) []byte {
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	sz := fi.Size()
	if int64(int(sz+4095)) != sz+4095 {
		log.Fatalf("%s: too large for mmap", f.Name())
	}
	n := int(sz)
	if n == 0 {
		return nil
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, (n+4095)&^4095, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("mmap %s: %v", f.Name(), err)
	}
	return data[:n]
}
