package rpc

import (
	"encoding/json"
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
	Error   error           `json:"error,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#requestMessage
type RequestMessage struct {
	Version string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#responseMessage
type ResponseMessage struct {
	Version string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   error           `json:"error,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notificationMessage
type NotificationMessage struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func NewResponse(id any, result any, err error) (*ResponseMessage, error) {
	r, _err := marshalRawMessage(result)
	return &ResponseMessage{ID: id, Result: r, Error: err}, _err
}

func NewNotification(m string, params any) (*NotificationMessage, error) {
	ps, err := marshalRawMessage(params)
	return &NotificationMessage{Method: m, Params: ps}, err
}

func (m RequestMessage) marshal(ptr *Envelope) {
	ptr.Version = m.Version
	ptr.ID = m.ID
	ptr.Method = m.Method
	ptr.Params = m.Params
}

func (m ResponseMessage) marshal(ptr *Envelope) {
	ptr.Version = m.Version
	ptr.ID = m.ID
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
