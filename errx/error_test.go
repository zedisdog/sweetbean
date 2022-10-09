package errx

import (
	"runtime/debug"
	"testing"
)

func TestNormal(t *testing.T) {
	println(string(debug.Stack()))
}
