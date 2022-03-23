package tools

import (
	"context"
	"errors"
)

//TimeOutFunc receive a chan of interface{}.
//error or result will be transmitted through channel.
type TimeOutFunc func(chan<- interface{})

func SetTimeout(ctx context.Context, f TimeOutFunc) (interface{}, error) {
	c := make(chan interface{})
	go f(c)
	select {
	case <-ctx.Done():
		return nil, errors.New("time out")
	case result := <-c:
		if err, ok := result.(error); ok {
			return nil, err
		} else {
			return result, nil
		}
	}
}
