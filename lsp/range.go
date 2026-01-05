package lsp

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Row int `json:"line"`
	Col int `json:"character"`
}
