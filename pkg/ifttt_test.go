package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIFTTT(t *testing.T) {
	apiKey := "key"
	i := newIFTTT(apiKey)
	expected := &ifttt{apiKey: apiKey}
	assert.Equal(t, expected, i)
}
