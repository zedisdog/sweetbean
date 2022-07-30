package errx

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type Error struct {
	File   string
	Line   int
	detail []interface{}
	msg    string
	err    error
}

func (e Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			detail, _ := json.Marshal(e.detail)
			fmt.Printf("%s:%s\n\t%s:%d\n", e.msg, detail, e.File, e.Line)
			if e.err != nil {
				fmt.Printf("%+v", e.err)
			}
		case s.Flag('#'):
			fmt.Printf("%s(%s:%d)\n", e.msg, e.File, e.Line)
			if e.err != nil {
				fmt.Printf("%#v", e.err)
			}
		default:
			println(e.Error())
		}
	}
}

func (e Error) Error() (result string) {
	result = e.msg
	if e.err != nil {
		result += "->" + e.err.Error()
	}
	return
}

func (e Error) Unwrap() error {
	return e.err
}

func New(msg string) *Error {
	return getPosition(&Error{
		msg: msg,
	}, 0)
}

func NewWithSkip(msg string, delta int) error {
	return getPosition(&Error{
		msg: msg,
	}, delta)
}

func WrapWithSkip(err error, msg string, delta int) error {
	return getPosition(&Error{
		msg: msg,
		err: err,
	}, delta)
}

func Wrap(err error, msg string) error {
	return getPosition(&Error{
		msg: msg,
		err: err,
	}, 0)
}

func getPosition(e *Error, delta int) *Error {
	skip := 2 + delta
	_, e.File, e.Line, _ = runtime.Caller(skip)
	return e
}
