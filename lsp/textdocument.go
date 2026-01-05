package lsp

type TextDocumentDidOpenNotificationParams struct {
	TextDoc struct {
		URI     string `json:"documentUri"`
		Lang    string `json:"languageId"`
		Version int32  `json:"version"`
		Text    string `json:"text"`
	} `json:"textDocument"`
}

type TextDocumentDidChangeNotificationParams struct {
	TextDoc struct {
		URI     string `json:"uri"`
		Version int32  `json:"version"`
	} `json:"textDocument"`

	Changes []TextDocumentChanges `json:"contentChanges"`
}

type TextDocumentDidCloseNotificationParams struct {
	TextDoc struct {
		URI     string `json:"uri"`
		Version int32  `json:"version"`
	} `json:"textDocument"`
}

type TextDocumentChanges struct {
	Range  *Range `json:"range,omitempty"`
	Length int    `json:"rangeLength,omitempty"`
	Text   string `json:"text"`
}
