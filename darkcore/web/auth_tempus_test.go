package web

import (
	"testing"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/stretchr/testify/assert"
)

func TestEncr2(t *testing.T) {
	token := NewTempusToken()
	assert.True(t, IsTempusValid(token, logus.Log))
}
