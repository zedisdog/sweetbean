package oauth

import (
	"fmt"
	"strings"
)

func WithScope(scope string) func(*redirectOptions) {
	return func(ro *redirectOptions) {
		ro.scope = scope
	}
}

func WithState(state map[string]string) func(*redirectOptions) {
	return func(ro *redirectOptions) {
		ro.state = state
	}
}

type redirectOptions struct {
	scope string
	state map[string]string
}

func (r redirectOptions) State() string {
	tmp := make([]string, 0, len(r.state))
	for key, value := range r.state {
		tmp = append(tmp, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(tmp, "|")
}
