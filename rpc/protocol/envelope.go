package protocol

import "encoding/json"

const rpcVersion = "2.0"

// generic message container
type envelope struct {
	Version string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

func NewEmptyEnvelope() envelope {
	return envelope{}
}

func NewVersionedEnvelope() envelope {
	return envelope{Version: rpcVersion}
}

func (e envelope) VersionOk() bool {
	return e.Version == rpcVersion
}
