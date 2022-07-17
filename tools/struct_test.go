package tools

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	A string `from:"A"`
	B string `from:"B,string"`
	C struct {
		D string `from:"C.D"`
	}
}

func TestGetTags(t *testing.T) {
	m := GetTags(reflect.TypeOf(test{}), "from", true)
	assert.Equal(t, m["A"], "A")
	assert.Equal(t, m["B"], "B,string")
	assert.Equal(t, m["C"], "")
	assert.Equal(t, m["C.D"], "C.D")
}

func TestConvert(t *testing.T) {
	sTest := test{
		A: "1",
		B: "2",
		C: struct {
			D string `from:"C.D"`
		}{
			D: "3",
		},
	}
	dTest := test{}
	assert.NotEqual(t, sTest, dTest)
	Convert(sTest, &dTest)
	assert.Equal(t, sTest, dTest)

	type testDto struct {
		A string `from:"A"`
		B string `from:"B"`
		C string `from:"C.D"`
	}
	dto := testDto{}
	Convert(sTest, &dto)
	assert.Equal(t, dto, testDto{
		A: "1",
		B: "2",
		C: "3",
	})

	type test2 struct {
		A string
		B string
		C struct{ D string }
	}
	testDto2 := test2{}
	Convert(dto, &testDto2)
	assert.Equal(t, testDto2, test2{
		A: "1",
		B: "2",
		C: struct{ D string }{D: "3"},
	})
}

func TestModifyUnexportAttr(t *testing.T) {
	type testDto struct {
		A int    `from:"a"`
		B string `from:"b"`
		C int    `from:"c.d"`
	}
	type test2 struct {
		a int
		b string
		c struct{ d int }
	}

	t1 := testDto{
		A: 1,
		B: "2",
		C: 3,
	}
	t2 := test2{}
	Convert(t1, &t2)
	t2.c.d = 2
	assert.Equal(t, t1, testDto{
		A: 1,
		B: "2",
		C: 3,
	})
	assert.Equal(t, t2, test2{
		a: 1,
		b: "2",
		c: struct{ d int }{d: 2},
	})
}

func TestParseFromTag(t *testing.T) {
	tags := map[string]string{
		"a": "a,string",
		"b": "b",
	}
	m := parseFromTag(tags)
	assert.Equal(t, m["a"].Name, "a")
	assert.Equal(t, m["a"].Type, "string")
	assert.Equal(t, m["b"].Name, "b")
	assert.Equal(t, m["b"].Type, "")
}
