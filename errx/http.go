package errx

import "net/http"

func WrapByHttpError(err error, code int, detail ...map[string]string) error {
	e := &HttpError{
		Code:  code,
		error: err,
	}
	if len(detail) > 0 {
		e.Detail = detail[0]
	}
	return e
}

func NewHttpError(code int, msg string, detail ...map[string]string) error {
	e := &HttpError{
		Code:  code,
		error: New(msg, 1),
	}
	if detail != nil && len(detail) > 0 {
		e.Detail = detail[0]
	}
	return e
}

type HttpError struct {
	Code int
	error
	Detail map[string]string
}

func NewHttpErrorUnprocessableEntity(msg string, detail ...map[string]string) error {
	return WrapByHttpError(
		New(msg, 1),
		http.StatusUnprocessableEntity,
		detail...,
	)
}

func NewHttpErrorBadRequest(msg string) error {
	return WrapByHttpError(
		New(msg, 1),
		http.StatusBadRequest,
	)
}

func NewHttpErrorForbidden(msg string) error {
	return WrapByHttpError(
		New(msg, 1),
		http.StatusForbidden,
	)
}

func NewHttpErrorConflict(msg string) error {
	return WrapByHttpError(
		New(msg, 1),
		http.StatusConflict,
	)
}

func NewHttpErrorTeapot(msg string) error {
	return WrapByHttpError(
		New(msg, 1),
		http.StatusTeapot,
	)
}
