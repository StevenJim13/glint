package glint

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func testArray(n int, b *testing.B) {
	testArray := []string{}
	for i := 0; i < n; i++ {
		testArray = append(testArray, strconv.Itoa(i))
	}
	exist := func(s []string, v string) bool {
		for index, _ := range s {
			if s[index] == v {
				return true
			}
		}
		return false
	}

	for i := 0; i < b.N; i++ {
		v := strconv.Itoa(i % n)
		exist(testArray, v)
	}
}

func testMap(n int, b *testing.B) {
	testMap := make(map[string]struct{})
	for i := 0; i < n; i++ {
		testMap[strconv.Itoa(i)] = struct{}{}
	}

	exist := func(s map[string]struct{}, v string) bool {
		_, ok := s[v]
		return ok
	}

	for i := 0; i < b.N; i++ {
		v := strconv.Itoa(i % n)
		exist(testMap, v)
	}
}

func BenchmarkArray6(b *testing.B) {
	testArray(6, b)
}

func BenchmarkMap6(b *testing.B) {
	testMap(6, b)
}

func BenchmarkArray10(b *testing.B) {
	testArray(10, b)
}

func BenchmarkMap10(b *testing.B) {
	testMap(10, b)
}

func TestMakeExclude(t *testing.T) {
	f, err := makeExcludeFunc(".*")
	require.NoError(t, err)
	fmt.Println(f("glint"))

}
