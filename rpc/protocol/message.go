package protocol

import (
	"encoding/json"
)

type RPCKind int

const (
	RequestKind RPCKind = iota
	ResponseKind
	NotificationKind
	VoidKind
)

type RPCMessage interface {
	Type() RPCKind
	Marshal(ptr *envelope)
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#requestMessage
type RequestMessage struct {
	ID     ID              `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#responseMessage
type ResponseMessage struct {
	ID     ID              `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ResponseError  `json:"error,omitempty"` // note :: convert into error???
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notificationMessage
type NotificationMessage struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

type NoResponseMessage struct{}

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

func (RequestMessage) Type() RPCKind {
	return RequestKind
}

func (ResponseMessage) Type() RPCKind {
	return ResponseKind
}

func (NotificationMessage) Type() RPCKind {
	return NotificationKind
}

func (NoResponseMessage) Type() RPCKind {
	return VoidKind
}

func (NoResponseMessage) Marshal(ptr *envelope) {}

func (m RequestMessage) Marshal(ptr *envelope) {
	ptr.ID = m.ID.Value()
	ptr.Method = m.Method
	ptr.Params = m.Params
}

func (m ResponseMessage) Marshal(ptr *envelope) {
	ptr.ID = m.ID.Value()
	ptr.Result = m.Result
	ptr.Error = m.Error
}

func (m NotificationMessage) Marshal(ptr *envelope) {
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
