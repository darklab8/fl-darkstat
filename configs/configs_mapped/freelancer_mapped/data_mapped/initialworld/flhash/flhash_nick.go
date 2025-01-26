package flhash

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
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

func (h HashCode) ToIntStr() string {
	return fmt.Sprintf("%d", int32(h))
}

func (h HashCode) ToUintStr() string {
	return strconv.FormatUint(uint64(int(h)), 10)
}

func (h HashCode) ToHexStr() string {
	return fmt.Sprintf("%x", int(h))
}

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
