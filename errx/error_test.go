package errx

import (
	"fmt"
	"testing"
)

func TestNormal(t *testing.T) {
	err := Wrap(New("test"), "testtt")
	fmt.Printf("%+v", err)
	println(err.Error())
}
