package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncr2(t *testing.T) {
	token := NewTempusToken()
	assert.True(t, IsTempusValid(token))
}
