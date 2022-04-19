package pkg_test

import (
	"io"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

func TestNewIFTTTClient(t *testing.T) {
	apiKey := "key"
	i := pkg.NewIFTTTClient(apiKey)
	expected := &pkg.IFTTTClient{APIKey: apiKey}
	assert.Equal(t, expected, i)
}

func TestIFTTTClient_Post(t *testing.T) {
	apiKey := "key"
	i := pkg.NewIFTTTClient(apiKey)
	monkey.PatchInstanceMethod(reflect.TypeOf(i), "Post", func(*pkg.IFTTTClient, string, string, io.Reader) error {
		return nil
	})
	err := i.Emit("event", "value1", "value2")
	assert.NoError(t, err)
}
