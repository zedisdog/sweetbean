package errx

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type Error struct {
	file   string
	line   int
	Code   int
	Detail interface{}
	Msg    string
	err    error
}

func (e Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			detail, _ := json.Marshal(e.Detail)
			fmt.Printf("%s:%s\n\t%s:%d\n", e.Msg, detail, e.file, e.line)
			if e.err != nil {
				fmt.Printf("%+v", e.err)
			}
		case s.Flag('#'):
			fmt.Printf("%s(%s:%d)\n", e.Msg, e.file, e.line)
			if e.err != nil {
				fmt.Printf("%#v", e.err)
			}
		default:
			println(e.Error())
		}
	}
}

func (e Error) Error() (result string) {
	result = e.Msg
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
		Msg: msg,
	}, 0)
}

func NewWithSkip(msg string, delta int) error {
	return getPosition(&Error{
		Msg: msg,
	}, delta)
}

func WrapWithSkip(err error, msg string, delta int) error {
	return getPosition(&Error{
		Msg: msg,
		err: err,
	}, delta)
}

func Wrap(err error, msg string) error {
	return getPosition(&Error{
		Msg: msg,
		err: err,
	}, 0)
}

func getPosition(e *Error, delta int) *Error {
	skip := 2 + delta
	_, e.file, e.line, _ = runtime.Caller(skip)
	return e
}
