package lsp

import (
	"encoding/json"
)

type (
	TraceMode  string
	EncKind    string
	SymbolKind int
	SymbolTag  int
)

const (
	FileSymbol SymbolKind = iota + 1
	ModuleSymbol
	NamespaceSymbol
	PackageSymbol
	ClassSymbol
	MethodSymbol
	PropertySymbol
	FieldSymbol
	ConstructorSymbol
	EnumSymbol
	InterfaceSymbol
	FunctionSymbol
	VariableSymbol
	ConstantSymbol
	StringSymbol
	NumberSymbol
	BooleanSymbol
	ArraySymbol
	ObjectSymbol
	KeySymbol
	NullSymbol
	EnumMemberSymbol
	StructSymbol
	EventSymbol
	OperatorSymbol
	TypeParameterSymbol

	Deprecated SymbolTag = 1

	UTF8  EncKind = "utf-8"
	UTF16 EncKind = "utf-16"
	UTF32 EncKind = "utf-32"

	Off      TraceMode = "off"
	Messages TraceMode = "messages"
	Verbose  TraceMode = "verbose"
)

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#initializeParams
type InitializeClientParams struct {
	// WorkDoneTok *struct {
	// 	ProgressTok any `json:"token"` // string | int
	// 	Value       any `json:"value"` // type: T
	// } `json:"workDoneToken,omitempty"`

	PID  int `json:"processId,omitempty"`
	Info *struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"clientInfo,omitempty"`

	Folders []struct {
		URI  string `json:"uri"`
		Name string `json:"name"`
	} `json:"workspaceFolders,omitempty"`

	Locale       string          `json:"locale,omitempty"`
	Root         string          `json:"rootUri,omitempty"`
	InitOpts     json.RawMessage `json:"initializationOptions,omitempty"`
	Capabilities Capabilities    `json:"capabilities"`
	Trace        TraceMode       `json:"trace,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#clientCapabilities
type Capabilities struct {
	Workspace    workspace       `json:"workspace"`
	TextDoc      json.RawMessage `json:"textDocument,omitempty"` // todo : https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocumentClientCapabilities
	NoteDoc      *noteDocCaps    `json:"notebookDocument,omitempty"`
	Window       *window         `json:"window,omitempty"`
	General      *general        `json:"general,omitempty"`
	Experimental json.RawMessage `json:"experimental,omitempty"`
}

// The client supports applying batch edits to the workspace by supporting the request 'workspace/applyEdit'
// Capabilities specific to `WorkspaceEdit`s
// Capabilities specific to the `workspace/didChangeConfiguration` notification
// Capabilities specific to the `workspace/didChangeWatchedFiles` notification.
// Capabilities specific to the `workspace/symbol` request.
// Capabilities specific to the `workspace/executeCommand` request.
// The client has support for workspace folders.
// The client supports `workspace/configuration` requests.
// Capabilities specific to the semantic token requests scoped to the workspace.
// Capabilities specific to the code lens requests scoped to the workspace.
// The client has support for file requests/notifications.
// Client workspace capabilities specific to inline values.
// Client workspace capabilities specific to inlay hints.
// Client workspace capabilities specific to diagnostics.
type workspace struct {
	ApplyEdit     bool            `json:"applyEdit,omitempty"`
	WorkspaceEdit json.RawMessage `json:"workspaceEdit,omitempty"` // todo :: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#workspaceEditClientCapabilities
	DidChangeConf *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"didChangeConfiguration,omitempty"`

	DidChangeWatched *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
		RelativePattern     bool `json:"relativePatternSupport,omitempty"`
	} `json:"didChangeWatchedFiles,omitempty"`

	Symbol *workspaceSymbolCaps `json:"symbol,omitempty"`

	ExecCmd *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	} `json:"executeCommand,omitempty"`
	Folders bool `json:"workspaceFolders,omitempty"`
	Conf    bool `json:"configuration,omitempty"`

	SemanticToks *struct {
		Refresh bool `json:"refreshSupport,omitempty"`
	} `json:"semanticTokens,omitempty"`

	CodeLens *struct {
		Refresh bool `json:"refreshSupport,omitempty"`
	} `json:"codeLens,omitempty"`

	FileOps *struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
		DidCreate           bool `json:"didCreate,omitempty"`
		WillCreate          bool `json:"willCreate,omitempty"`
		DidRename           bool `json:"didRename,omitempty"`
		WillRename          bool `json:"willRename,omitempty"`
		DidDelete           bool `json:"didDelete,omitempty"`
		WillDelete          bool `json:"willDelete,omitempty"`
	} `json:"fileOperations,omitempty"`

	Inline *struct {
		Refresh bool `json:"refreshSupport,omitempty"`
	} `json:"inlineValue,omitempty"`

	Inlay *struct {
		Refresh bool `json:"refreshSupport,omitempty"`
	} `json:"inlayHint,omitempty"`

	Diags *struct {
		Refresh bool `json:"refreshSupport,omitempty"`
	} `json:"diagnostics,omitempty"`
}

type workspaceSymbolCaps struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	Kind *struct {
		Values []SymbolKind `json:"valueSet,omitempty"`
	} `json:"symbolKind,omitempty"`

	Tag *struct {
		Values []SymbolTag `json:"valueSet"`
	} `json:"tagSupport,omitempty"`

	Resolve *struct {
		Props []string `json:"properties"`
	} `json:"resolveSupport,omitempty"`
}

type noteDocCaps struct {
	Sync struct {
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
		ExecSummary         bool `json:"executionSummarySupport,omitempty"`
	} `json:"synchronization"`
}

type window struct {
	ProgressDone bool `json:"workDoneProgress,omitempty"`

	ShowMsg *struct {
		ActionItem *struct {
			Props bool `json:"additionalPropertiesSupport,omitempty"`
		} `json:"messageActionItem,omitempty"`
	} `json:"showMessage,omitempty"`

	ShowDoc *struct {
		Support bool `json:"support"`
	} `json:"showDocument,omitempty"`
}

type general struct {
	StaleRequest *struct {
		Cancel     bool     `json:"cancel"`
		RetryOnMod []string `json:"retryOnContentModified"`
	} `json:"staleRequestSupport,omitempty"`

	Regex *struct {
		Engine  string `json:"engine"`
		Version string `json:"version,omitempty"`
	} `json:"regularExpressions,omitempty"`

	Markdown *struct {
		Parser       string   `json:"parser"`
		Version      string   `json:"version,omitempty"`
		TagWhitelist []string `json:"allowedTags,omitempty"`
	} `json:"markdown,omitempty"`

	PosEncs []EncKind `json:"positionEncodings,omitempty"`
}
