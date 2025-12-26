package rpc

import "encoding/json"

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#responseMessage
var (
	// Protocol standard errors
	ParseError     = newResponseError(-32700, "RPC: parse error")
	InvalidRequest = newResponseError(-32600, "RPC: invalid request")
	MethodNotFound = newResponseError(-32601, "RPC: method not found")
	InvalidParams  = newResponseError(-32602, "RPC: invalid parameters")
	InternalError  = newResponseError(-32603, "RPC: internal error")

	ServerNotInitialized = newResponseError(-32002, "RPC: Server not initialized")
	UnknownErrorCode     = newResponseError(-32001, "RPC: Unknown error code")

	// A request failed but it was syntactically correct, e.g the
	// method name was known and the parameters were valid. The error
	// message should contain human readable information about why
	// the request failed.
	RequestFailed = newResponseError(-32803, "RPC: Request failed")

	// The server cancelled the request. This error code should
	// only be used for requests that explicitly support being
	// server cancellable.
	ServerCancelled = newResponseError(-32802, "RPC: Server cancelled")

	// The server detected that the content of a document got
	// modified outside normal conditions. A server should
	// NOT send this error code if it detects a content change
	// in its unprocessed messages. The result even computed
	// on an older state might still be useful for the client.
	//
	// If a client decides that a result is not of any use anymore
	// the client should cancel the request.
	ContentModified = newResponseError(-32801, "RPC: Content modified")

	// The client has canceled a request and a server has detected
	// the cancel.
	RequestCancelled = newResponseError(-32800, "RPC: Request cancelled")

	jsonrpcReservedErrorRangeStart = newResponseError(-32099, "RPC: Reserved range start")
	jsonrpcReservedErrorRangeEnd   = newResponseError(-32000, "RPC: Reserved range end")
	lspReservedErrorRangeStart     = newResponseError(-32899, "RPC: Reserved range start: LSP")
	lspReservedErrorRangeEnd       = newResponseError(-32800, "RPC: Reserved range end: LSP")
)

type ResponseError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func newResponseError(code int, msg string) ResponseError {
	return ResponseError{Code: code, Message: msg}
}

func (e ResponseError) Error() string {
	return e.Message
}
