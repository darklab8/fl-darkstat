package header

type (
	// Header holds the combined used header information according to the specifications.
	Header struct {
		DDSHeader
		DDPFHeader
		CapsHeader
		DX10Header
		FourCCString string // the string representation of the DDPFHeader.FourCC
	}

	// DDSHeader is the definition header for the dds texture file
	DDSHeader struct {
		TextureFlags      Flags[DDSf] // flag-set
		Height            uint32      // width of the texture in pixels
		Width             uint32      // height of the texture in pixels
		PitchOrLinearSize uint32      // size of a scanline for uncompressed or the first texture for compressed
		Depth             uint32      // depth for a volume texture in pixels
		MipMapCount       uint32      // number of mip-map levels
	}

	// DDPFHeader holds the information how the textures stores its pixel-data
	DDPFHeader struct {
		PixelFlags  Flags[DDPFf] // flag-set
		FourCC      uint32       // 4 character-code for the used texture compression (compressed)
		RgbBitCount uint32       // bit-size for every color/pixel (uncompressed)
		RBitMask    uint32       // mask for reading the R in RGB or Y in YUV or the data for Luminance (uncompressed)
		GBitMask    uint32       // mask for reading the G in RGB or U in YUV (uncompressed)
		BBitMask    uint32       // mask for reading the B in RGB or V in YUV (uncompressed)
		ABitMask    uint32       // mask for reading the alpha data (if DDPFAlphaPixels) (uncompressed)
	}

	// CapsHeader wraps the cube-map specific flag set
	CapsHeader struct {
		Caps1 Flags[DDSCf] // Specifies the complexity of the surfaces stored
		Caps2 uint32       // Additional detail about the surfaces stored
		Caps3 uint32       // unused
		Caps4 uint32       // unused
	}

	// DX10Header is an extension fo the default header in case the FourCC is set to "DX10"
	DX10Header struct {
		DxgiFormat        uint32        // the pixel format as gigantic enum. replaces the DDPFHeader definitions
		ResourceDimension Flags[DDSDTc] // dimension of the texture: 1D, 2D or 3D
		MiscFlag          uint32        // more obscure settings regarding cube maps
		ArraySize         uint32        // the number of elements in the array (amount of textures inside)
		MiscFlags2        uint32        // bits regarding more precise description of alpha values
	}

	// Flags is a convenience type for flag-bit checking operations
	Flags[f ~uint32] struct{ F f }
)

// Has returns if the flags contain all given bits
func (d Flags[f]) Has(v f) bool {
	return d.F&v == v
}

// Is returns if the flags only contains the given bits
func (d Flags[f]) Is(v f) bool {
	return d.F^v == 0
}

// Not return if the flags do not contain any given bits at all.
func (d Flags[f]) Not(v f) bool {
	return d.F&^v == d.F
}
