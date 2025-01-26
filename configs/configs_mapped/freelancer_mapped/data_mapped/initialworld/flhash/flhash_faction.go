package flhash

import (
	"strings"
)

const (
	FLFACTIONHASH_POLYNOMIAL = 0x1021
)

type faction_hasher struct {
	table [256]uint32
}

// Function for calculating the Freelancer data nickname hash.
// Algorithm from flhash.exe by sherlog@t-online.de (2003-06-11)
func (h *faction_hasher) RawHash(data []byte) uint32 {
	var hash uint32 = 0xFFFF
	for _, b := range data {
		hash = (hash >> 8) ^ h.table[byte(hash&0xFF)^b]
	}
	return hash
}

// NicknameHasher implements the hashing algorithm used by item, base, etc. nicknames
type FactionNickHasher struct {
	faction_hasher
}

func (h *FactionNickHasher) Hash(name string) HashCode {
	bytes := []byte(strings.ToLower(name))
	hash := h.RawHash(bytes)
	// hash = (hash >> (physicalBits - logicalBits)) | 0x80000000
	return HashCode(hash)
}

func NewFactionHasher() *FactionNickHasher {
	h := FactionNickHasher{}
	NotSimplePopulateTable(FLFACTIONHASH_POLYNOMIAL, &h.table)
	return &h
}

var faction = NewFactionHasher()

func HashFaction(name string) HashCode {
	return faction.Hash(name)
}

// Copied foundation from golang crc32 that had SimplePopulate
// And then modified to how it should at Python version given by Alex
// Matches also Laz version in C# i think
func NotSimplePopulateTable(poly uint32, t *[256]uint32) {
	for i := 0; i < 256; i++ {
		crc := uint32(i) << 8
		for j := 0; j < 8; j++ {
			crc <<= 1

			if crc&0x10000 != 0 {
				crc = (crc ^ poly) & 0xFFFF
			}
		}
		t[i] = crc
	}
}
