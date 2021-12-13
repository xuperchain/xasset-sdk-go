package utils

import (
	"fmt"
	"testing"
)

func TestGenAssetId(t *testing.T) {
	id1 := GenAssetId(123456)
	id2 := GenAssetId(123456)
	id3 := GenAssetId(789)

	fmt.Println(id1, id2, id3)
}

func TestGenNonce(t *testing.T) {
	cnt := 1000000
	m := make(map[int64]int)
	for i := 0; i < cnt; i++ {
		n := GenNonce()
		if _, ok := m[n]; !ok {
			m[n] = 0
		}
		m[n] = m[n] + 1
	}

	dup := 0
	for _, v := range m {
		if v > 1 {
			dup += v
		}
	}

	fmt.Println(cnt, dup)
}

func TestGetFuncCall(t *testing.T) {
	file, fc := GetFuncCall(1)
	fmt.Println(file, fc)
}
