package flhash

// Provided by marko_oktabyr (discord)

import (
	"fmt"
	"testing"
)

func TestHasher_Hash(t *testing.T) {
	tests := []struct {
		hash HashCode
		name string
	}{
		{hash: HashCode(2339324873), name: "pl_ge_fighter4"},
		{hash: HashCode(2194637003), name: "ge_fighter4_power01"},
		{hash: HashCode(2723858309), name: "ge_s_scanner_01"},
	}

	for _, test := range tests {
		code := HashNickname(test.name)
		if code != test.hash {
			t.Errorf("%s: expected %d but got %d", test.name, test.hash, code)
		}
	}
}

func TestNameHash_Hash(t *testing.T) {
	tests := []struct {
		flName string
		name   string
	}{
		{flName: "05-fda83b36", name: "Marko"},
		{flName: "04-09a4114e", name: "Test"},
		{flName: "23-ef57f351", name: "a3f82aa3-97a28ec6-6ef10b83-b3b5d470"},
	}

	for _, test := range tests {
		sName := SaveFile(test.name)
		if sName != test.flName {
			t.Errorf("%s: expected %s but got %s", test.name, test.flName, sName)
		}
	}
}

func TestHashSmth(t *testing.T) {
	hash := HashNickname("st02_03_base")
	fmt.Println(hash.ToIntStr())
	fmt.Println(hash.ToUintStr())
	fmt.Println(hash.ToHexStr())
}
