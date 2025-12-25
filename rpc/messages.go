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
	Error   *ResponseError  `json:"error,omitempty"`
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
	Error   *ResponseError  `json:"error,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notificationMessage
type NotificationMessage struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func (m RequestMessage) marshal(ptr *Envelope) {
	m.Version = ptr.Version
	m.ID = ptr.ID
	m.Method = ptr.Method
	m.Params = ptr.Params
}

func (m ResponseMessage) marshal(ptr *Envelope) {
	m.Version = ptr.Version
	m.ID = ptr.ID
	m.Result = ptr.Result
	m.Error = ptr.Error
}

func (m NotificationMessage) marshal(ptr *Envelope) {
	m.Version = ptr.Version
	m.Method = ptr.Method
	m.Params = ptr.Params
}
