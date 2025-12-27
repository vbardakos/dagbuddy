package rpc

import (
	"encoding/json"
	"fmt"
	"os"
)

const rpcVersion = "2.0"

type Message interface {
	marshal(ptr *Envelope)
}

type Envelope struct {
	Version string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#requestMessage
type RequestMessage struct {
	Version string          `json:"jsonrpc"`
	ID      ID              `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#responseMessage
type ResponseMessage struct {
	Version string          `json:"jsonrpc"`
	ID      ID              `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"` // note :: convert into error???
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notificationMessage
type NotificationMessage struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type ID struct {
	value any
}

func NewRequest(id ID, method string, params any) (*RequestMessage, error) {
	ps, merr := marshalRawMessage(params)
	return &RequestMessage{ID: id, Method: method, Params: ps}, merr
}

func NewResponse(id ID, result any, err *ResponseError) (*ResponseMessage, error) {
	r, merr := marshalRawMessage(result)
	return &ResponseMessage{ID: id, Result: r, Error: err}, merr
}

func NewNotification(m string, params any) (*NotificationMessage, error) {
	ps, err := marshalRawMessage(params)
	return &NotificationMessage{Method: m, Params: ps}, err
}

func newID(raw any) (ID, error) {
	var id ID
	switch raw.(type) {
	case nil, string, int64:
		id = ID{value: raw}
	case float64:
		fmt.Fprintf(os.Stderr, "Caught float64: %d", raw)
		id = ID{value: int64(raw.(float64))}
	default:
		return id, fmt.Errorf("Unknown ID type. Got: %s\n", raw)
	}
	return id, nil
}

func (id ID) Value() any {
	return id.value
}

func (id ID) IsOk() bool {
	switch id.value.(type) {
	case string:
		return true
	case int64:
		code := id.value.(int)
		rpcErrRange := code >= jsonrpcReservedErrorRangeStart.Code && code < jsonrpcReservedErrorRangeEnd.Code
		return rpcErrRange
	default:
		return false
	}
}

func (m RequestMessage) marshal(ptr *Envelope) {
	ptr.Version = m.Version
	ptr.ID = m.ID.value
	ptr.Method = m.Method
	ptr.Params = m.Params
}

func (m ResponseMessage) marshal(ptr *Envelope) {
	ptr.Version = m.Version
	ptr.ID = m.ID.value
	ptr.Result = m.Result
	ptr.Error = m.Error
}

func (m NotificationMessage) marshal(ptr *Envelope) {
	ptr.Version = m.Version
	ptr.Method = m.Method
	ptr.Params = m.Params
}

func marshalRawMessage(value any) (json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(data), nil
}
