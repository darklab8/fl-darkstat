package header

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDDPFFlags(t *testing.T) {
	flag := Flags[DDPFf]{F: DDPFRGB | DDPFAlphaPixels}

	assert.True(t, flag.Has(DDPFAlphaPixels))
	assert.True(t, flag.Has(DDPFRGB))
	assert.False(t, flag.Has(DDPFFourCC))
	assert.True(t, flag.Has(DDPFRGB|DDPFAlphaPixels))
	assert.False(t, flag.Has(DDPFRGB|DDPFFourCC))

	assert.False(t, flag.Is(DDPFAlphaPixels))
	assert.False(t, flag.Is(DDPFRGB))
	assert.False(t, flag.Is(DDPFFourCC))
	assert.True(t, flag.Is(DDPFRGB|DDPFAlphaPixels))
	assert.False(t, flag.Is(DDPFRGB|DDPFFourCC))

	assert.False(t, flag.Not(DDPFAlphaPixels))
	assert.False(t, flag.Not(DDPFRGB))
	assert.True(t, flag.Not(DDPFFourCC))
	assert.False(t, flag.Not(DDPFRGB|DDPFAlphaPixels))
	assert.False(t, flag.Not(DDPFRGB|DDPFFourCC))
}
