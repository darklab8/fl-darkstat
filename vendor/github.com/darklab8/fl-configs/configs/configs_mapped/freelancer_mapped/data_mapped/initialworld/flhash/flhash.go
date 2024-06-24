package flhash

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"strings"
	"unicode/utf16"
)

const (
	logicalBits                 = 30
	physicalBits                = 32
	flHashPolynomial     uint32 = 0xA001 << (logicalBits - 16)
	flNameHashPolynomial uint32 = 0x50008 << (physicalBits - 20)
)

type hasher struct {
	table [256]uint32
}

// Function for calculating the Freelancer data nickname hash.
// Algorithm from flhash.exe by sherlog@t-online.de (2003-06-11)
func (h *hasher) RawHash(data []byte) uint32 {
	var hash uint32
	for _, b := range data {
		hash = (hash >> 8) ^ h.table[byte(hash)^b]
	}
	hash = (hash >> 24) | ((hash >> 8) & 0x0000FF00) | ((hash << 8) & 0x00FF0000) | (hash << 24)
	return hash
}

// NicknameHasher implements the hashing algorithm used by item, base, etc. nicknames
type NicknameHasher struct {
	hasher
}

type HashCode int

func (h *NicknameHasher) Hash(name string) HashCode {
	bytes := []byte(strings.ToLower(name))
	hash := h.RawHash(bytes)
	hash = (hash >> (physicalBits - logicalBits)) | 0x80000000
	return HashCode(hash)
}

func NewHasher() *NicknameHasher {
	h := NicknameHasher{}
	h.table = *crc32.MakeTable(flHashPolynomial)
	return &h
}

var nick = NewHasher()

func HashNickname(name string) HashCode {
	return nick.Hash(name)
}

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
