package container

import (
	"fmt"
	"testing"
)

type itest interface {
	speak()
}

type test struct {
	Content string
}

func (t test) speak() {
	println(t.Content)
}

func TestNormal(t *testing.T) {
	a := &test{}
	SetT[itest](a)

	b := Get[itest]()

	fmt.Printf("%+v\n", b)
}
