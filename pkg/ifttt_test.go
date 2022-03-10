package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

func TestNewIFTTT(t *testing.T) {
	apiKey := "key"
	i := pkg.ExportedNewIFTTT(apiKey)
	expected := &pkg.ExportedIFTTT{APIKey: apiKey}
	assert.Equal(t, expected, i)
}
