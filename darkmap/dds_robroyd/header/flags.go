package header

// DDSf is the flag type for the DDSHeader.TextureFlags
type DDSf uint32

// flags are related to the DDSHeader.TextureFlags itself
const (
	DDSDCaps        DDSf = 0x1      // CapsHeader is set
	DDSDHeight      DDSf = 0x2      // DDSHeader.Height is set
	DDSDWidth       DDSf = 0x4      // DDSHeader.Width is set
	DDSDPitch       DDSf = 0x8      // if there is an uncompressed texture with pitch
	DDSDPixelFormat DDSf = 0x1000   // DDPFHeader is set
	DDSDMipMapCount DDSf = 0x20000  // DDSHeader.MipMapCount is set. If texture has mipmaps
	DDSDLinearSize  DDSf = 0x80000  // if there is a compressed texture with pitch
	DDSDDepth       DDSf = 0x800000 // DDSHeader.Depth is set. If texture has depth
)

// combined header flags related to the DDSHeader.TextureFlags
const (
	DDSDHeaderFlagsTexture    DDSf = DDSDCaps | DDSDHeight | DDSDWidth | DDSDPixelFormat
	DDSDHeaderFlagsMipMap     DDSf = DDSDMipMapCount
	DDSDHeaderFlagsVolume     DDSf = DDSDDepth
	DDSDHeaderFlagsPitch      DDSf = DDSDPitch
	DDSDHeaderFlagsLinearSize DDSf = DDSDLinearSize
)

// DDPFf is the flag type for the DDPFHeader.PixelFlags
type DDPFf uint32

// flags for the DDPFHeader.PixelFlags itself
const (
	DDPFAlphaPixels DDPFf = 0x1     // texture contains alpha data
	DDPFAlpha       DDPFf = 0x2     // texture contains only uncompressed alpha data
	DDPFFourCC      DDPFf = 0x4     // texture contains compressed RGB data
	DDPFRGB         DDPFf = 0x40    // texture contains uncompressed RGB data
	DDPFYUV         DDPFf = 0x200   // texture contains uncompressed YUV data
	DDPFLuminance   DDPFf = 0x20000 // texture contains a single channel uncompressed data
)

// DDSCf is the flag type for CapsHeader.Caps1
type DDSCf uint32

// flags for the first dword of the CapsHeader.Caps1
const (
	DDSCAPSComplex DDSCf = 0x8
	DDSCAPSMipmap  DDSCf = 0x400000
	DDSCAPSTexture DDSCf = 0x1000
)

// DDSDTc is the dimension flag for a DC10 texture
type DDSDTc uint32

// constants for the dimensions of a DX10 texture
const (
	DDSDT1D DDSDTc = iota + 2 // 1D texture
	DDSDT2D                   // 2D texture
	DDSDT3D                   // 3D texture
)
