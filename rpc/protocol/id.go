package protocol

import (
	"math"
)

type ID struct {
	value any
}

// note :: jsonrpc expects numeric or string
func NewID(raw any) (ID, error) {
	var id ID

	if s, ok := raw.(string); ok {
		id.value = s
	}

	// note :: numerics expected to be parsed as f64
	if n, ok := raw.(float64); ok {
		if math.Trunc(n) != n {
			return id, ParseError
		}
		id.value = int64(n)
	}

	if n, ok := raw.(int); ok {
		id.value = int64(n)
	}

	if id.value == nil {
		return id, ParseError
	}

	return id, nil
}

func (id ID) Value() any {
	return id.value
}
