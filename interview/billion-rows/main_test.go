package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/DenzelPenzel/go-leetcode/interview/billion-rows/sol1"
	"github.com/DenzelPenzel/go-leetcode/interview/billion-rows/sol2"

	"github.com/stretchr/testify/assert"
)

func readFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func find(root, ext string) []string {
	var res []string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			res = append(res, path[:len(path)-len(ext)])
		}
		return nil
	})
	return res
}

func Test_TestSol1(t *testing.T) {
	fileNames := find("./test_cases", ".txt")
	for _, name := range fileNames {
		t.Run(name, func(t *testing.T) {
			got := sol1.Run(name + ".txt")
			want := readFile(name + ".out")
			assert.Equal(t, want, got)
		})
	}
}

func Test_TestSol2(t *testing.T) {
	fileNames := find("./test_cases", ".txt")
	for _, name := range fileNames {
		t.Run(name, func(t *testing.T) {
			got := sol2.Run(name + ".txt")
			want := readFile(name + ".out")
			assert.Equal(t, want, got)
		})
	}
}

// func Test_TestSol3(t *testing.T) {
// 	fileNames := find("./test_cases", ".txt")
// 	for _, name := range fileNames {
// 		t.Run(name, func(t *testing.T) {
// 			got := sol3.Run(name + ".txt")
// 			want := readFile(name + ".out")
// 			assert.Equal(t, want, got)
// 		})
// 	}
// }

// func Test_TestSol4(t *testing.T) {
// 	fileNames := find("./test_cases", ".txt")
// 	for _, name := range fileNames {
// 		t.Run(name, func(t *testing.T) {
// 			got := sol4.Run(name + ".txt")
// 			want := readFile(name + ".out")
// 			assert.Equal(t, want, got)
// 		})
// 	}
// }

// func Test_TestHash(t *testing.T) {
// 	type testCase struct {
// 		input []byte
// 		size  int
// 		want  []byte
// 	}
// 	tests := []testCase{
// 		{[]byte("a;bcdefg"), 1, []byte("a")},
// 		{[]byte("ab;cdefg"), 2, []byte("ab")},
// 		{[]byte("abc;defg"), 3, []byte("abc")},
// 		{[]byte("abcd;efg"), 4, []byte("abcd")},
// 		{[]byte("abcde;fg"), 5, []byte("abcde")},
// 		{[]byte("abcdef;g"), 6, []byte("abcdef")},
// 		{[]byte("abcdefg;"), 7, []byte("abcdefg")},
// 		{[]byte("abcdefgy"), 8, []byte("abcdefgy")},
// 		{[]byte("abcdefgpoi"), 8, []byte("abcdefgp")},
// 		{[]byte("abcdefgpoi"), 9, []byte("abcdefgpo")},
// 	}

// 	for i, tc := range tests {
// 		u := *(*uint64)(unsafe.Pointer(&tc.input[0]))
// 		r := sol3.MakeHashKey(u, tc.size)
// 		e := *(*sol3.Hash)(unsafe.Pointer(&tc.want[0]))

// 		if r != e {
// 			panic(fmt.Sprintf("ts id: %d, %q: want %d, got %d", i, tc.input, tc.want, r))
// 		}
// 	}
// }

// func Test_TestSemicol(t *testing.T) {
// 	type testCase struct {
// 		input []byte
// 		want  int
// 	}
// 	tests := []testCase{
// 		{[]byte("a;bcdefg"), 1},
// 		{[]byte("ab;cdefg"), 2},
// 		{[]byte("abc;defg"), 3},
// 		{[]byte("abcd;efg"), 4},
// 		{[]byte("abcde;fg"), 5},
// 		{[]byte("abcdef;g"), 6},
// 		{[]byte("abcdefgy"), -1},
// 		{[]byte("abcdefgpoi"), -1},
// 		{[]byte("a;;;;;;;"), 1},
// 		{[]byte("a;bc;def"), 1},
// 		{[]byte("a;dcdefg"), 1},
// 		{[]byte("abcdef;;"), 6},
// 		{[]byte(";abcdefg"), 0},
// 		{[]byte("abcdefg;"), 7},
// 		{[]byte("abcdefgh;"), -1},
// 		{[]byte("Zagreb;"), 6},
// 	}

// 	for _, tc := range tests {
// 		u := *(*uint64)(unsafe.Pointer(&tc.input[0]))
// 		r := sol3.FindSemicolon(u)
// 		if r != tc.want {
// 			panic(fmt.Sprintf("%q: want %d, got %d", tc.input, tc.want, r))
// 		}
// 	}
// }
