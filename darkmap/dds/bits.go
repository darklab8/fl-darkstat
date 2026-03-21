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
	"fmt"
)

func readDWORD(buf []byte) (uint32, []byte) {
	return uint32(buf[3])<<24 | uint32(buf[2])<<16 | uint32(buf[1])<<8 | uint32(buf[0]), buf[4:]
}

func readBits(buf []byte, n uint32) uint32 {
	if n&0x7 != 0 || n > 32 {
		panic(fmt.Sprintf("can't read %d bits", n))
	}
	var u uint32
	x := n / 8
	for y := uint32(0); y < x; y++ {
		u |= uint32(buf[y]) << uint(y*8)
	}
	return u
}

var (
	ctz32DeBruijn = [32]byte{
		0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,
		31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9,
	}
)

func lowestSetBit(n uint32) uint {
	// TODO(lukegb): switch to using math/bits when available
	if n == 0 {
		return 32
	}
	return uint(ctz32DeBruijn[(n&-n)*0x077CB531>>(32-5)])
}
