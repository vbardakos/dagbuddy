package rpc

import (
	"math"
)

type ID struct {
	str *string
	num *int64
}

func NewID(raw any) (ID, error) {
	if s, ok := raw.(string); ok {
		return ID{str: &s}, nil
	}

	if n, ok := raw.(int); ok {
		int := int64(n)
		return ID{num: &int}, nil
	}

	if n, ok := raw.(float64); ok {
		if math.Trunc(n) == n {
			int := int64(n)
			return ID{num: &int}, nil
		}
	}

	return ID{}, ParseError
}

func (id ID) Value() any {
	if id.str != nil {
		return *id.str
	}
	return *id.num
}
