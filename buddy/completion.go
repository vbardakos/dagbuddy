package buddy

import (
	"context"
	"encoding/json"

	s "github.com/vbardakos/dagbuddy/lsp/schemas"
	"github.com/vbardakos/dagbuddy/rpc"
)

type CompletionRequest struct{}

func (*CompletionRequest) Handle(ctx context.Context, msg rpc.Message) (rpc.Message, error) {
	rq, ok := msg.(rpc.Request)
	if !ok {
		return nil, rpc.InternalError
	}

	var params map[string]any
	if err := json.Unmarshal(rq.Params, &params); err != nil {
		return nil, err
	}

	sess := GetLspSession(ctx)
	sess.Log.Printf("completion params: %s\n", params)

	return rpc.NewResponse(rq.ID, s.CompletionList{
		// IsIncomplete: true,
		// ItemDefaults: &s.CompletionItemDefaults{
		// 	EditRange: &s.CompletionEditRange{
		// 		Insert:  s.Range{Start: s.Position{Row: 0, Col: 0}, End: s.Position{Row: 99, Col: 99}},
		// 		Replace: s.Range{Start: s.Position{Row: 0, Col: 0}, End: s.Position{Row: 99, Col: 99}},
		// 	},
		// },
		Items: []s.CompletionItem{
			{
				Label: "HelloWorld",
				Kind:  s.CompleteText,
				// InsertText: "HELLO_WORLD",
				// InsertMode: s.AsIs,
				TextEdit: &s.TextEdit{
					Range: s.Range{Start: s.Position{Row: 0, Col: 0}, End: s.Position{Row: 99, Col: 99}},
					Text:  "HELLOWORLD",
				},
			},
			// {
			// 	Label:      "HelloYou",
			// 	Kind:       s.CompleteFunction,
			// 	InsertText: "HELLO_YOU",
			// 	InsertMode: s.AsIs,
			// },
		},
	}, nil)
}
