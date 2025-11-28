package json

import (
	encodingjson "encoding/json"

	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/support/convert"
)

type Json struct {
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func New() foundation.Json {
	return &Json{
		marshal:   encodingjson.Marshal,
		unmarshal: encodingjson.Unmarshal,
	}
}

func (j *Json) Marshal(v any) ([]byte, error) {
	return j.marshal(v)
}

func (j *Json) Unmarshal(data []byte, v any) error {
	return j.unmarshal(data, v)
}

func (j *Json) MarshalString(a any) (string, error) {
	b, err := j.Marshal(a)
	if err != nil {
		return "", err
	}
	return convert.UnsafeString(b), nil
}

func (j *Json) UnmarshalString(s string, a any) error {
	return j.Unmarshal(convert.UnsafeBytes(s), a)
}
