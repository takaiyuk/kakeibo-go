package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

func TestNewIFTTTClient(t *testing.T) {
	apiKey := "key"
	i := pkg.NewIFTTTClient(apiKey)
	expected := &pkg.ExportedIFTTTClient{APIKey: apiKey}
	assert.Equal(t, expected, i)
}

// TODO: Do not monkey patch
//
// func TestIFTTTClient_Post(t *testing.T) {
// 	apiKey := "key"
// 	i := pkg.NewIFTTTClient(apiKey)
// 	monkey.PatchInstanceMethod(reflect.TypeOf(i), "Post", func(*pkg.ExportedIFTTTClient, string, string, io.Reader) error {
// 		return nil
// 	})
// 	err := i.Emit("event", "value1", "value2")
// 	assert.NoError(t, err)
// }
