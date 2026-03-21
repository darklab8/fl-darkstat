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
	"io"
)

const (
	headerSize      = 124 // Size of DDS_HEADER structure
	pixelFormatSize = 32  // Size of DDS_PIXELFORMAT structure

	dCaps        = 0x1
	dHeight      = 0x2
	dWidth       = 0x4
	dPitch       = 0x8
	dPixelFormat = 0x1000
	dMipMapCount = 0x20000
	dLinearSize  = 0x80000
	dDepth       = 0x800000

	pfAlphaPixels = 0x1
	pfAlpha       = 0x2
	pfFourCC      = 0x4
	pfRGB         = 0x40
	pfYUV         = 0x200
	pfLuminance   = 0x20000

	headerFlagsTexture    = dCaps | dHeight | dWidth | dPixelFormat
	headerFlagsMipMap     = dMipMapCount
	headerFlagsVolume     = dDepth
	headerFlagsPitch      = dPitch
	headerFlagsLinearSize = dLinearSize
)

type pixelFormat struct {
	flags       uint32
	fourCC      uint32
	rgbBitCount uint32
	rBitMask    uint32
	gBitMask    uint32
	bBitMask    uint32
	aBitMask    uint32
}

type header struct {
	flags             uint32
	height            uint32
	width             uint32
	pitchOrLinearSize uint32
	depth             uint32
	mipMapCount       uint32
	pixelFormat       pixelFormat
	caps              [4]uint32
}

func readHeader(r io.Reader) (header, error) {
	var buf []byte

	// read the magic
	buf = make([]byte, 4)
	if n, err := r.Read(buf); n != 4 || err != nil {
		return header{}, fmt.Errorf("reading magic: %v", err)
	}
	if buf[0] != 'D' || buf[1] != 'D' || buf[2] != 'S' || buf[3] != ' ' {
		return header{}, fmt.Errorf("magic is incorrect, expected \"DDS \", got %v", buf)
	}

	// read the dds file header
	buf = make([]byte, 124)
	if n, err := r.Read(buf); n != 124 || err != nil {
		return header{}, fmt.Errorf("reading header: %v", err)
	}

	var t uint32
	if t, buf = readDWORD(buf); t != headerSize {
		return header{}, fmt.Errorf("DDS_HEADER reports wrong size, expected %d, got %d", t, headerSize)
	}

	var h header
	h.flags, buf = readDWORD(buf)
	h.height, buf = readDWORD(buf)
	h.width, buf = readDWORD(buf)
	h.pitchOrLinearSize, buf = readDWORD(buf)
	h.depth, buf = readDWORD(buf)
	h.mipMapCount, buf = readDWORD(buf)
	buf = buf[11*4:] // strip off reserved1
	if t, buf = readDWORD(buf); t != pixelFormatSize {
		return header{}, fmt.Errorf("DDS_PIXEL_FORMAT reports wrong size, expected %d, got %d", t, pixelFormatSize)
	}
	pf := h.pixelFormat
	pf.flags, buf = readDWORD(buf)
	pf.fourCC, buf = readDWORD(buf)
	pf.rgbBitCount, buf = readDWORD(buf)
	pf.rBitMask, buf = readDWORD(buf)
	pf.gBitMask, buf = readDWORD(buf)
	pf.bBitMask, buf = readDWORD(buf)
	pf.aBitMask, buf = readDWORD(buf)
	h.pixelFormat = pf
	for n := 0; n < 4; n++ {
		h.caps[n], buf = readDWORD(buf)
	}
	buf = buf[4:] // strip off reserved2
	if len(buf) > 0 {
		return header{}, fmt.Errorf("trailing garbage remains: %d bytes", len(buf))
	}

	// // check that flags is valid
	// if h.flags&headerFlagsTexture != headerFlagsTexture {
	// 	return header{}, fmt.Errorf("DDS_HEADER reports that one or more required fields are not set: flags was %x; should at least have %x set", h.flags, headerFlagsTexture)
	// }

	return h, nil
}
