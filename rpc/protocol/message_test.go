package protocol

import (
	"encoding/json"
	"testing"
)

func equalEnvelopes(e1 envelope, e2 envelope) bool {
	sanity := e1.Version == e2.Version && e1.ID == e2.ID && e1.Method == e2.Method && (e1.Error == nil) == (e2.Error == nil)
	if !sanity {
		return false
	}

	equalBytes := func(r1 json.RawMessage, r2 json.RawMessage) bool {
		if len(r1) != len(r2) {
			return false
		}
		for i := range len(r1) {
			if r1[i] != r2[i] {
				return false
			}
		}
		return true
	}

	rawMessagePairs := map[*json.RawMessage]json.RawMessage{
		&e1.Params: e2.Params,
		&e1.Result: e2.Result,
	}

	for r1, r2 := range rawMessagePairs {
		if !equalBytes(*r1, r2) {
			return false
		}
	}

	if e1.Error != nil {
		sanity := e1.Error.Code == e2.Error.Code && e2.Error.Error() == e2.Error.Error()
		if !sanity {
			return false
		}
		return equalBytes(e1.Error.Data, e2.Error.Data)
	}

	return true
}

func TestMessage_Marshal(t *testing.T) {
	tests := []struct {
		name string
		data RPCMessage
		got  envelope
	}{
		{
			name: "response result",
			data: ResponseMessage{
				ID:     ID{value: 10},
				Result: json.RawMessage([]byte("{\"hello\":\"here\"}")),
			},
			got: envelope{
				Version: "",
				ID:      ID{value: 10}.Value(),
				Method:  "",
				Params:  json.RawMessage{},
				Result:  json.RawMessage([]byte("{\"hello\":\"here\"}")),
				Error:   nil,
			},
		},
		{
			name: "response error",
			data: ResponseMessage{
				ID:    ID{value: 10},
				Error: &InternalError,
			},
			got: envelope{
				Version: "",
				ID:      ID{value: 10}.Value(),
				Method:  "",
				Params:  json.RawMessage{},
				Result:  json.RawMessage{},
				Error:   &InternalError,
			},
		},
		{
			name: "request",
			data: RequestMessage{
				ID:     ID{value: 10},
				Method: "hello/world",
				Params: json.RawMessage([]byte("{\"hello\":\"here\"}")),
			},
			got: envelope{
				Version: "",
				ID:      ID{value: 10}.Value(),
				Method:  "hello/world",
				Params:  json.RawMessage([]byte("{\"hello\":\"here\"}")),
				Result:  json.RawMessage{},
				Error:   nil,
			},
		},
		{
			name: "request no params",
			data: RequestMessage{
				ID:     ID{value: 10},
				Method: "hello/world",
			},
			got: envelope{
				Version: "",
				ID:      ID{value: 10}.Value(),
				Method:  "hello/world",
				Params:  json.RawMessage{},
				Result:  json.RawMessage{},
				Error:   nil,
			},
		},
		{
			name: "request",
			data: NotificationMessage{
				Method: "hello/world",
				Params: json.RawMessage([]byte("{\"hello\":\"here\"}")),
			},
			got: envelope{
				Version: "",
				ID:      nil,
				Method:  "hello/world",
				Params:  json.RawMessage([]byte("{\"hello\":\"here\"}")),
				Result:  json.RawMessage{},
				Error:   nil,
			},
		},
		{
			name: "notification no params",
			data: NotificationMessage{
				Method: "hello/world",
			},
			got: envelope{
				Version: "",
				ID:      nil,
				Method:  "hello/world",
				Params:  json.RawMessage{},
				Result:  json.RawMessage{},
				Error:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := envelope{}
			tt.data.Marshal(&e)
			if !equalEnvelopes(e, tt.got) {
				t.Fatalf("envelopes do not match: %+v; want: %+v", e, tt.got)
			}
		})
	}
}
