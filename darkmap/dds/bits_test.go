/*
Copyright 2017 Luke Granger-Brown

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dds

import (
	"bytes"
	"fmt"
	"testing"
)

func TestReadDWORD(t *testing.T) {
	buf := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5}
	got, rest := readDWORD(buf)
	if got != 0xddccbbaa {
		t.Errorf("readDWORD(%x) = %x; want ddccbbaa", buf, got)
	}
	if !bytes.Equal(rest, buf[4:]) {
		t.Errorf("readDWORD(%x): rest is %x, want %x", buf, rest, buf[4:])
	}
}

func TestReadBits(t *testing.T) {
	for _, test := range []struct {
		buf  []byte
		n    uint32
		want uint32
	}{
		{
			[]byte{0xaa, 0xbb, 0xcc, 0xdd},
			32,
			0xddccbbaa,
		},
		{
			[]byte{0xaa, 0xbb, 0xcc, 0xdd},
			24,
			0xccbbaa,
		},
		{
			[]byte{0xaa, 0xbb, 0xcc, 0xdd},
			16,
			0xbbaa,
		},
		{
			[]byte{0xaa, 0xbb, 0xcc, 0xdd},
			8,
			0xaa,
		},
	} {
		got := readBits(test.buf, test.n)
		if got != test.want {
			t.Errorf("readBits(%x, %d) = %x; want %x", test.buf, test.n, got, test.want)
		}
	}
}

func TestReadBitsNotByteAligned(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("wanted test to panic, but it didn't")
		}
	}()
	readBits([]byte{0xaa, 0xbb, 0xcc, 0xdd}, 1)
}

func TestReadBitsMoreThan32(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("wanted test to panic, but it didn't")
		}
	}()
	readBits([]byte{0xaa, 0xbb, 0xcc, 0xdd}, 64)
}

func TestLowestSetBit(t *testing.T) {
	for _, test := range []struct {
		in   uint32
		want uint
	}{
		{0, 32},
		{0x1, 0},
		{0x2, 1},
		{0x3, 0},
		{0x4, 2},
		{0x5, 0},
		{0x6, 1},
		{0x7, 0},
		{0x10, 4},
		{0x100, 8},
		{0x1000, 12},
		{0x10000, 16},
		{0x100000, 20},
		{0x1000000, 24},
		{0x10000000, 28},
		{0x80000000, 31},
		{0x80000001, 0},
	} {
		t.Run(fmt.Sprintf("%d", test.in), func(t *testing.T) {
			got := lowestSetBit(test.in)
			if got != test.want {
				t.Errorf("lowestSetBit(%d) = %d; want %d", test.in, got, test.want)
			}
		})
	}
}
