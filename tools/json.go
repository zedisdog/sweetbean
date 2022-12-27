package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Json string

func (c *Json) UnmarshalJSON(b []byte) (err error) {
	*c = Json(b)
	return nil
}

func (c Json) MarshalJSON() ([]byte, error) {
	return []byte(c), nil
}

func (c Json) Get(name string) (value interface{}, err error) {
	tmp := make(map[string]interface{})
	err = json.Unmarshal([]byte(c), &tmp)
	if err != nil {
		return
	}
	var ok bool
	for _, n := range strings.Split(name, ".") {
		if value == nil {
			value, ok = tmp[n]
		} else {
			value, ok = value.(map[string]interface{})[n]
		}
		if !ok {
			err = fmt.Errorf("value of key <%s> not found", n)
			return
		}
	}

	return
}

func (c Json) GetString(name string) (value string, err error) {
	v, err := c.Get(name)
	if err != nil {
		return
	}
	value = v.(string)
	return
}

func (c Json) As(container interface{}) (err error) {
	err = json.Unmarshal([]byte(c), container)
	return
}
