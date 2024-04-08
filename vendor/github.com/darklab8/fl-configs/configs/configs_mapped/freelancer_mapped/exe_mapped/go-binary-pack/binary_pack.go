// Copyright 2017 Roman Kachanovsky. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package binary_pack performs conversions between some Go values represented as byte slices.
	This can be used in handling binary data stored in files or from network connections,
	among other sources. It uses format slices of strings as compact descriptions of the layout
	of the Go structs.

	Format characters (some characters like H have been reserved for future implementation of unsigned numbers):
		? - bool, packed size 1 byte
		b - int8, packed size 1 bytes
		B - uint8, packed size 1 bytes
		h - int16, packed size 2 bytes
		H - uint16, packed size 2 bytes
		i, l - int32, packed size 4 bytes
		I, L - int32, packed size 4 bytes
		q - int64, packed size 8 bytes
		Q - uint64, packed size 8 bytes
		f - float32, packed size 4 bytes
		d - float64, packed size 8 bytes
		Ns - string, packed size N bytes, N is a number of runes to pack/unpack

*/
package binarypack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// BinaryPack presents a BinaryPack
type BinaryPack struct{}

// New create a new BinaryPack
func New() *BinaryPack {
	return &BinaryPack{}
}

// Pack returns a byte slice containing the values of msg slice packed according to the given format.
// The items of msg slice must match the values required by the format exactly.
func (bp *BinaryPack) Pack(format []string, msg []interface{}) ([]byte, error) {
	if len(format) > len(msg) {
		return nil, errors.New("Format is longer than values to pack")
	}

	res := []byte{}

	for i, f := range format {
		switch f {
		case "?":
			castedV, ok := msg[i].(bool)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (bool)")
			}
			res = append(res, boolToBytes(castedV)...)
		case "b", "B":
			castedV, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 2 bytes)")
			}
			res = append(res, intToBytes(castedV, 1)...)
		case "h", "H":
			castedV, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 2 bytes)")
			}
			res = append(res, intToBytes(castedV, 2)...)
		case "i", "I", "l", "L":
			castedV, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 4 bytes)")
			}
			res = append(res, intToBytes(castedV, 4)...)
		case "q", "Q":
			castedV, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 8 bytes)")
			}
			res = append(res, intToBytes(castedV, 8)...)
		case "f":
			castedV, ok := msg[i].(float32)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (float32)")
			}
			res = append(res, float32ToBytes(castedV, 4)...)
		case "d":
			castedV, ok := msg[i].(float64)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (float64)")
			}
			res = append(res, float64ToBytes(castedV, 8)...)
		default:
			if strings.Contains(f, "s") {
				castedV, ok := msg[i].(string)
				if !ok {
					return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (string)")
				}
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, []byte(fmt.Sprintf("%s%s",
					castedV, strings.Repeat("\x00", n-len(castedV))))...)
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return res, nil
}

// UnPack the byte slice (presumably packed by Pack(format, msg)) according to the given format.
// The result is a []interface{} slice even if it contains exactly one item.
// The byte slice must contain not less the amount of data required by the format
// (len(msg) must more or equal CalcSize(format)).
func (bp *BinaryPack) UnPack(format []string, msg []byte) ([]interface{}, error) {
	expectedByte, err := bp.CalcSize(format)

	if err != nil {
		return nil, err
	}

	if expectedByte > len(msg) {
		return nil, errors.New("Expected size is bigger than actual size of message")
	}

	res := []interface{}{}

	for _, f := range format {
		switch f {
		case "b":
			res = append(res, bytesToInt(msg[:1]))
			msg = msg[1:]
		case "B":
			res = append(res, bytesToUint(msg[:1]))
			msg = msg[1:]
		case "?":
			res = append(res, bytesToBool(msg[:1]))
			msg = msg[1:]
		case "h":
			res = append(res, bytesToInt(msg[:2]))
			msg = msg[2:]
		case "H":
			res = append(res, bytesToUint(msg[:2]))
			msg = msg[2:]
		case "i", "l":
			res = append(res, bytesToInt(msg[:4]))
			msg = msg[4:]
		case "I", "L":
			res = append(res, bytesToUint(msg[:4]))
			msg = msg[4:]
		case "q":
			res = append(res, bytesToInt(msg[:8]))
			msg = msg[8:]
		case "Q":
			res = append(res, bytesToUint(msg[:8]))
			msg = msg[8:]
		case "f":
			res = append(res, bytesToFloat32(msg[:4]))
			msg = msg[4:]
		case "d":
			res = append(res, bytesToFloat64(msg[:8]))
			msg = msg[8:]
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, string(msg[:n]))
				msg = msg[n:]
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return res, nil
}

// CalcSize Returns the size of the struct (and hence of the byte slice) corresponding to the given format.
func (bp *BinaryPack) CalcSize(format []string) (int, error) {
	var size int

	for _, f := range format {
		switch f {
		case "b", "B":
			size++
		case "?":
			size++
		case "h", "H":
			size += 2
		case "i", "I", "l", "L", "f":
			size += 4
		case "q", "Q", "d":
			size += 8
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				size = size + n
			} else {
				return 0, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return size, nil
}

func boolToBytes(x bool) []byte {
	if x {
		return intToBytes(1, 1)
	}
	return intToBytes(0, 1)
}

func bytesToBool(b []byte) bool {
	return bytesToInt(b) > 0
}

func intToBytes(n int, size int) []byte {
	buf := bytes.NewBuffer([]byte{})

	switch size {
	case 1:
		binary.Write(buf, binary.LittleEndian, int8(n))
	case 2:
		binary.Write(buf, binary.LittleEndian, int16(n))
	case 4:
		binary.Write(buf, binary.LittleEndian, int32(n))
	default:
		binary.Write(buf, binary.LittleEndian, int64(n))
	}

	return buf.Bytes()[0:size]
}

func bytesToInt(b []byte) int {
	buf := bytes.NewBuffer(b)

	switch len(b) {
	case 1:
		var x int8
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 2:
		var x int16
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 4:
		var x int32
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	default:
		var x int64
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	}
}

func bytesToUint(b []byte) int {
	buf := bytes.NewBuffer(b)

	switch len(b) {
	case 1:
		var x uint8
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 2:
		var x uint16
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 4:
		var x uint32
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	default:
		var x uint64
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	}
}

func float32ToBytes(n float32, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()[0:size]
}

func bytesToFloat32(b []byte) float32 {
	var x float32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &x)
	return x
}

func float64ToBytes(n float64, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()[0:size]
}

func bytesToFloat64(b []byte) float64 {
	var x float64
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &x)
	return x
}
