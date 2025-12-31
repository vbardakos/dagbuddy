package codec_test

import (
	"encoding/json"
	"reflect"
	"testing"

	p "github.com/vbardakos/dagbuddy/rpc/protocol"
	c "github.com/vbardakos/dagbuddy/rpc/codec"
)

func msgEquals(m1 p.RPCMessage, m2 p.RPCMessage) bool {
	if r1, ok := m1.(p.RequestMessage); ok {
		var p1 string
		json.Unmarshal(r1.Params, &p1)

		if r2, ok := m2.(p.RequestMessage); ok {
			var p2 string
			json.Unmarshal(r2.Params, &p2)
			return r1.ID.Value() == r2.ID.Value() && r1.Method == r2.Method && p1 == p2
		}

		return false
	}

	if n1, ok := m1.(p.NotificationMessage); ok {
		var p1 string
		json.Unmarshal(n1.Params, &p1)

		if n2, ok := m2.(p.NotificationMessage); ok {
			var p2 string
			json.Unmarshal(n2.Params, &p2)
			return n1.Method == n2.Method && p1 == p2
		}

		return false
	}

	if r1, ok := m1.(p.ResponseMessage); ok {
		var res1 string
		json.Unmarshal(r1.Result, &res1)

		if r2, ok := m2.(p.ResponseMessage); ok {
			var res2 string
			json.Unmarshal(r2.Result, &res2)
			return r1.ID.Value() == r2.ID.Value() && r1.Error == r2.Error && res1 == res2
		}

		return false
	}

	return false
}

func errEquals(err error, msg p.RPCMessage) bool {
	if r, ok := msg.(p.ResponseMessage); ok {
		if rerr, ok := err.(p.ResponseError); ok {
			return rerr.Code == r.Error.Code && rerr.Message == r.Error.Message
		}
	}
	return false
}

func TestEncodeMessage(t *testing.T) {
	id, _ := p.NewID(10)

	tests := []struct {
		name    string
		msg     p.RPCMessage
		want    []byte
		wantErr bool
	}{
		{
			name: "request",
			msg: p.RequestMessage{
				ID:     id,
				Method: "method",
				Params: json.RawMessage("{\"value\": \"hello\"}"),
			},
			want:    []byte("{\"jsonrpc\":\"2.0\",\"id\":10,\"method\":\"method\",\"params\":{\"value\":\"hello\"}}"),
			wantErr: false,
		},
		{
			name: "notification",
			msg: p.NotificationMessage{
				Method: "method",
				Params: json.RawMessage("{\"value\": \"hello\"}"),
			},
			want:    []byte("{\"jsonrpc\":\"2.0\",\"method\":\"method\",\"params\":{\"value\":\"hello\"}}"),
			wantErr: false,
		},
		{
			name: "response result",
			msg: p.ResponseMessage{
				ID:     id,
				Result: json.RawMessage("{\"value\": \"hello\"}"),
			},
			want:    []byte("{\"jsonrpc\":\"2.0\",\"id\":10,\"result\":{\"value\":\"hello\"}}"),
			wantErr: false,
		},
		{
			name: "response result",
			msg: p.ResponseMessage{
				ID: id,
				Error: &p.ResponseError{
					Code:    10,
					Message: "hello",
				},
			},
			want:    []byte("{\"jsonrpc\":\"2.0\",\"id\":10,\"error\":{\"code\":10,\"message\":\"hello\"}}"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := c.EncodeMessage(tt.msg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("EncodeMessage() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("EncodeMessage() succeeded unexpectedly")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeMessage() = %s, want %s", got, tt.want)
				t.Errorf("\n%v\n%v\n", got, tt.want)
			}
		})
	}
}

func TestDecodeMessage(t *testing.T) {
	id, _ := p.NewID(10)

	tests := []struct {
		name    string
		data    []byte
		want    p.RPCMessage
		wantErr bool
	}{
		{
			name: "request",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10, \"method\": \"init\", \"params\": {\"hello\":\"world\"}}"),
			want: p.RequestMessage{
				ID:     id,
				Method: "init",
				Params: json.RawMessage([]byte("{\"hello\":\"world\"}")),
			},
			wantErr: false,
		},
		{
			name: "request no params",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10, \"method\": \"init\"}"),
			want: p.RequestMessage{
				ID:     id,
				Method: "init",
			},
			wantErr: false,
		},
		{
			name: "notification",
			data: []byte("{\"jsonrpc\": \"2.0\", \"method\": \"init\", \"params\": {\"hello\":\"world\"}}"),
			want: p.NotificationMessage{
				Method: "init",
				Params: json.RawMessage([]byte("{\"hello\":\"world\"}")),
			},
			wantErr: false,
		},
		{
			name:    "notification no params",
			data:    []byte("{\"jsonrpc\": \"2.0\", \"method\": \"init\"}"),
			want:    p.NotificationMessage{Method: "init"},
			wantErr: false,
		},
		{
			name:    "response result",
			data:    []byte("{\"jsonrpc\": \"2.0\", \"id\": 10, \"result\": {\"init\": 0}}"),
			want:    p.ResponseMessage{ID: id, Result: json.RawMessage([]byte("{\"init\": 0}"))},
			wantErr: false,
		},
		{
			name: "response error",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10, \"error\": {\"code\": 0, \"message\": \"hello\"}}"),
			want: p.ResponseMessage{
				ID: id,
				Error: &p.ResponseError{
					Code:    0,
					Message: "hello",
				},
			},
			wantErr: true,
		},
		{
			name: "error InternalError",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10}"),
			want: p.ResponseMessage{
				Error: &p.InternalError,
			},
			wantErr: true,
		},
		{
			name: "error ParseError (from ID)",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10.2,\"result\":{}}"),
			want: p.ResponseMessage{
				Error: &p.ParseError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := c.DecodeMessage(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					if !errEquals(gotErr, tt.want) || got != nil {
						t.Errorf("DecodeMessage() failed: %+v", gotErr)
					}
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DecodeMessage() succeeded unexpectedly")
			}

			if !msgEquals(got, tt.want) {
				t.Errorf("DecodeMessage() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
