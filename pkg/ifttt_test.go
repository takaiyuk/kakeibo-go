package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

func TestNewIFTTTClient(t *testing.T) {
	apiKey := "key"
	i := pkg.NewIFTTTClient(apiKey)
	expected := &pkg.IFTTTClient{APIKey: apiKey}
	assert.Equal(t, expected, i)
}
