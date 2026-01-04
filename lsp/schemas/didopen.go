package schemas

type DidOpenTextDocument struct {
	TextDoc struct {
		URI     string `json:"documentUri"`
		Lang    string `json:"languageId"`
		Version int    `json:"version"`
		Text    string `json:"text"`
	} `json:"textDocument"`
}
