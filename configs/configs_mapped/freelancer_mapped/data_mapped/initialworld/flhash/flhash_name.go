package flhash

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"strings"
	"unicode/utf16"
)

// NameHash implements the hashing algorithm used by the account folders and save files.
type NameHash struct {
	hasher
}

func (h *NameHash) Hash(name string) uint32 {
	codes := utf16.Encode([]rune(strings.ToLower(name)))
	bytes := make([]byte, 2*len(codes))
	for i, c := range codes {
		binary.LittleEndian.PutUint16(bytes[2*i:2*(i+1)], c)
	}
	return h.RawHash(bytes)
}

func (h *NameHash) SaveFile(name string) string {
	return fmt.Sprintf("%02x-%08x", len(name), h.Hash(name))
}

func NewNameHasher() *NameHash {
	h := NameHash{}
	h.table = *crc32.MakeTable(flNameHashPolynomial)
	return &h
}

var n = NewNameHasher()

func SaveFile(name string) string {
	return n.SaveFile(name)
}
