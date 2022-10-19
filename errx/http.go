package errx

import "net/http"

func WrapByHttpError(err error, code int, msg string, detail map[string]string) error {
	if err == nil {
		return nil
	}
	e := Wrap(err, msg)
	e.(*Error).Code = code
	e.(*Error).Detail = detail
	return e
}

func NewHttpError(code int, msg string, detail map[string]string) error {
	e := New(msg)
	e.(*Error).Detail = detail
	e.(*Error).Code = code
	return e
}

func NewHttpErrorUnprocessableEntity(msg string, detail map[string]string) error {
	return NewHttpError(
		http.StatusUnprocessableEntity,
		msg,
		detail,
	)
}

func NewHttpErrorBadRequest(msg string) error {
	return NewHttpError(
		http.StatusBadRequest,
		msg,
		nil,
	)
}

func NewHttpErrorForbidden(msg string) error {
	return NewHttpError(
		http.StatusForbidden,
		msg,
		nil,
	)
}

func NewHttpErrorConflict(msg string) error {
	return NewHttpError(
		http.StatusConflict,
		msg,
		nil,
	)
}

func NewHttpErrorTeapot(msg string) error {
	return NewHttpError(
		http.StatusTeapot,
		msg,
		nil,
	)
}

func NewHttpErrorUnauthorized(msg string) error {
	return NewHttpError(
		http.StatusUnauthorized,
		msg,
		nil,
	)
}

func NewHttpErrorNotFound(msg string) error {
	return NewHttpError(
		http.StatusNotFound,
		msg,
		nil,
	)
}

func NewHttpErrorInternalServer(msg string) error {
	return NewHttpError(
		http.StatusInternalServerError,
		msg,
		nil,
	)
}

func NewHttpErrorBadGateway(msg string) error {
	return NewHttpError(
		http.StatusBadGateway,
		msg,
		nil,
	)
}
