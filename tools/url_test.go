package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeQuery(t *testing.T) {
	encoded := EncodeQuery("mysql://root:toor@tcp(mysql)/main?collation=utf8mb4_general_ci&loc=Asia/Shanghai&parseTime=true")
	assert.Equal(t, "mysql://root:toor@tcp(mysql)/main?collation=utf8mb4_general_ci&loc=Asia%2FShanghai&parseTime=true", encoded)
	encoded = EncodeQuery("mysql://root:toor@tcp(mysql:1234)/main?collation=utf8mb4_general_ci&loc=Asia/Shanghai&parseTime=true")
	assert.Equal(t, "mysql://root:toor@tcp(mysql:1234)/main?collation=utf8mb4_general_ci&loc=Asia%2FShanghai&parseTime=true", encoded)
	encoded = EncodeQuery("mysql://root:toor@tcp(mysql:1234)/main?")
	assert.Equal(t, "mysql://root:toor@tcp(mysql:1234)/main?", encoded)
	encoded = EncodeQuery("mysql://root:toor@tcp(mysql:1234)/main")
	assert.Equal(t, "mysql://root:toor@tcp(mysql:1234)/main", encoded)
}
