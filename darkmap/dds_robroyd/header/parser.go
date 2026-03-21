package header

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	sizeDDTF   = 128    // Size of the whole texture file header. is 128
	sizeDDSD   = 124    // Size of the serialized DDSHeader. is 124
	sizeDDPF   = 32     // Size of the serialized DDPFHeader. is 32
	sizeDX10   = 20     // Size of the serialized optional DX10Header. is 20
	FourCCDX10 = "DX10" // the fourCC string for the presence of the extra DX10 header
)

// deserializer is used to parse all header bytes into a structure
type deserializer struct {
	MagicNumber     uint32     // magic number must be "DDS "
	HeaderSize      uint32     // header size. must be 124
	DDSHeader                  // the texture header
	_               [11]uint32 // reserved1
	PixelFormatSize uint32     // pixel format size. must be 32
	DDPFHeader                 // the pixel format header
	CapsHeader                 // header for more complex textures with mipmaps or cube-maps
	_               [1]uint32  // reserved2
}

// Read tries to take 128 Bytes from the reader and then tries to create a Header from it.
func Read(r io.Reader) (*Header, error) {
	return new(deserializer).read(r)
}

// read tries to take sizeDDTF Bytes from the reader and then tries to create a Header from it.
// If it finds the FourCCDX10 header on the DDPFHeader.FourCC it will try to parse the DX10Header.
// Calls verification on a successful parsed Header, which might return an error in the case of a
// wrongly configured header.
func (d *deserializer) read(r io.Reader) (*Header, error) {
	if err := d.readChunk(r, sizeDDTF, d); err != nil {
		return nil, err
	} else if err = d.verify(); err != nil {
		return nil, err
	}

	header := &Header{
		DDSHeader:    d.DDSHeader,
		DDPFHeader:   d.DDPFHeader,
		CapsHeader:   d.CapsHeader,
		FourCCString: d.toString(d.FourCC),
	}

	if header.FourCCString == FourCCDX10 {
		if err := d.readChunk(r, sizeDX10, &header.DX10Header); err != nil {
			return nil, err
		}
	}
	return header, nil
}

// readChunk reads in a portion of the stream and tries to deserialize it to the given target
func (*deserializer) readChunk(r io.Reader, size int, target any) error {
	buf := make([]byte, size, size)
	if n, err := r.Read(buf); err != nil {
		return err
	} else if n != size {
		return fmt.Errorf("bytes in header. expected %d, only got : %d", size, n)
	}
	return binary.Read(bytes.NewReader(buf), binary.LittleEndian, target)
}

// verify makes some semantic checks for validity
func (d *deserializer) verify() error {
	if mn := d.toString(d.MagicNumber); mn != "DDS " {
		return fmt.Errorf("magic is incorrect, expected \"DDS \", got %v", mn)
	}
	if d.HeaderSize != sizeDDSD {
		return fmt.Errorf("DDS_HEADER reports wrong size, expected %d, got %d", sizeDDSD, d.HeaderSize)
	}
	if d.PixelFormatSize != sizeDDPF {
		return fmt.Errorf("DDS_PIXEL_FORMAT reports wrong size, expected %d, got %d", sizeDDPF, d.PixelFormatSize)
	}

	// check that it's actually a texture per requirements
	if !d.TextureFlags.Has(DDSDHeaderFlagsTexture) {
		return fmt.Errorf("DDS_HEADER reports that one or more required fields are not set: flags was %x; should at least have %x set", d.TextureFlags, DDSDHeaderFlagsTexture)
	}
	return nil
}

func (*deserializer) toString(i uint32) string {
	return string(binary.LittleEndian.AppendUint32(nil, i))
}
