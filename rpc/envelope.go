package rpc

import "encoding/json"

const rpcVersion = "2.0"

// generic container of a message
type envelope struct {
	Version string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

func versionedEnvelope() envelope {
	return envelope{Version: rpcVersion}
}
