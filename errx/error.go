package errx

import (
	"errors"
	"fmt"
	"runtime"
)

type Error struct {
	File string
	Line int
	error
}

func (e Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			fmt.Printf("%s:%d %s\n", e.File, e.Line, e.error.Error())
		default:
			println(e.error.Error())
		}
	}
	println(e.error.Error())
}

func New(msg string, delta ...int) error {
	e := &Error{
		File:  "???",
		error: errors.New(msg),
	}
	skip := 1
	if len(delta) > 0 {
		skip += delta[0]
	}
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		e.File = file
		e.Line = line
	}
	return e
}

func Wrap(err error, msg string, delta ...int) error {
	e := &Error{
		File:  "???",
		error: fmt.Errorf("%s:%w", msg, err),
	}
	skip := 1
	if len(delta) > 0 {
		skip += delta[0]
	}
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		e.File = file
		e.Line = line
	}
	return e
}
