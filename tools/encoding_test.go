package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	src := "123456jfdkalfj"

	d, err := Base64Decode(string(Base64Encode([]byte(src))))
	assert.Nil(t, err)
	assert.Equal(t, src, string(d))
}
