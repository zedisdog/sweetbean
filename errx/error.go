package errx

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
)

type Error struct {
	Code   int
	Stack  []byte
	Detail interface{}
	Msg    string
	err    error
}

func (e Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			println(e.BuildDetail())
		case s.Flag('#'):
			fallthrough
		default:
			println(e.Error())
		}
	}
}

func (e Error) BuildDetail() string {
	detail, _ := json.Marshal(e.Detail)
	stack := string(e.Stack)
	if e.err != nil {
		return fmt.Sprintf("\n>>>>>>>>\n%s: %s\ndetail: %s\n%s\n<<<<<<<<<\n", e.Msg, e.err.Error(), detail, stack)
	} else {
		return fmt.Sprintf("\n>>>>>>>>\n%s\ndetail: %s\n%s\n<<<<<<<<<\n", e.Msg, detail, stack)
	}
}

func (e Error) Error() (result string) {
	result = e.Msg
	if e.err != nil {
		result += ": " + e.err.Error()
	}
	return
}

func (e Error) Unwrap() error {
	return e.err
}

func New(msg string) error {
	return &Error{
		Msg:   msg,
		Stack: debug.Stack(),
	}
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	e := &Error{
		Msg:   msg,
		err:   err,
		Stack: debug.Stack(),
	}

	if ee, ok := interface{}(err).(*Error); ok {
		e.Code = ee.Code
	}

	return e
}
