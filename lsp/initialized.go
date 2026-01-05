package lsp

type TextDocSyncKind int

const (
	None TextDocSyncKind = iota
	Full
	Incremental
)

type InitializeServerResult struct {
	Caps *ServerCapabilities `json:"capabilities,omitempty"`
	Info *ServerInfo         `json:"serverInfo,omitempty"`
}

type ServerCapabilities struct {
	TextDocSync TextDocSync        `json:"textDocumentSync"`
	Completion  *CompletionOptions `json:"completionProvider,omitempty"`
	Hover       bool               `json:"hoverProvider"`
	Definition  bool               `json:"definitionProvider"`
}

type TextDocSync struct {
	OpenClose bool            `json:"openClose"`
	Change    TextDocSyncKind `json:"change"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type CompletionOptions struct {
	TriggerChars []string `json:"triggerCharacters"`
	Resolve      bool     `json:"resolveProvider,omitempty"`
}
