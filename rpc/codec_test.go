package rpc_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/vbardakos/dagbuddy/rpc"
)

func msgEquals(m1 rpc.Message, m2 rpc.Message) bool {
	if r1, ok := m1.(rpc.RequestMessage); ok {
		var p1 string
		json.Unmarshal(r1.Params, &p1)

		if r2, ok := m2.(rpc.RequestMessage); ok {
			var p2 string
			json.Unmarshal(r2.Params, &p2)
			return r1.ID.Value() == r2.ID.Value() && r1.Method == r2.Method && p1 == p2
		}

		return false
	}

	if n1, ok := m1.(rpc.NotificationMessage); ok {
		var p1 string
		json.Unmarshal(n1.Params, &p1)

		if n2, ok := m2.(rpc.NotificationMessage); ok {
			var p2 string
			json.Unmarshal(n2.Params, &p2)
			return n1.Method == n2.Method && p1 == p2
		}

		return false
	}

	if r1, ok := m1.(rpc.ResponseMessage); ok {
		var res1 string
		json.Unmarshal(r1.Result, &res1)

		if r2, ok := m2.(rpc.ResponseMessage); ok {
			var res2 string
			json.Unmarshal(r2.Result, &res2)
			return r1.ID.Value() == r2.ID.Value() && r1.Error == r2.Error && res1 == res2
		}

		return false
	}

	return false
}

func errEquals(err error, msg rpc.Message) bool {
	if r, ok := msg.(rpc.ResponseMessage); ok {
		if rerr, ok := err.(rpc.ResponseError); ok {
			return rerr.Code == r.Error.Code && rerr.Message == r.Error.Message
		}
	}
	return false
}

func TestEncodeMessage(t *testing.T) {
	id, _ := rpc.NewID(10)

	tests := []struct {
		name    string
		msg     rpc.Message
		want    []byte
		wantErr bool
	}{
		{
			name: "Request",
			msg: rpc.RequestMessage{
				ID:     id,
				Method: "method",
				Params: json.RawMessage("{\"value\": \"hello\"}"),
			},
			want:    []byte("{\"jsonrpc\":\"2.0\",\"id\":10,\"method\":\"method\",\"params\":{\"value\":\"hello\"}}"),
			wantErr: false,
		},
		{
			name: "Notification",
			msg: rpc.NotificationMessage{
				Method: "method",
				Params: json.RawMessage("{\"value\": \"hello\"}"),
			},
			want:    []byte("{\"jsonrpc\":\"2.0\",\"method\":\"method\",\"params\":{\"value\":\"hello\"}}"),
			wantErr: false,
		},
		{
			name: "ResponseResult",
			msg: rpc.ResponseMessage{
				ID:     id,
				Result: json.RawMessage("{\"value\": \"hello\"}"),
			},
			want:    []byte("{\"jsonrpc\":\"2.0\",\"id\":10,\"result\":{\"value\":\"hello\"}}"),
			wantErr: false,
		},
		{
			name: "ResponseResult",
			msg: rpc.ResponseMessage{
				ID: id,
				Error: &rpc.ResponseError{
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
			got, gotErr := rpc.EncodeMessage(tt.msg)
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
	id, _ := rpc.NewID(10)

	tests := []struct {
		name    string
		data    []byte
		want    rpc.Message
		wantErr bool
	}{
		{
			name: "Request",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10, \"method\": \"init\", \"params\": {\"hello\":\"world\"}}"),
			want: rpc.RequestMessage{
				ID:     id,
				Method: "init",
				Params: json.RawMessage([]byte("{\"hello\":\"world\"}")),
			},
			wantErr: false,
		},
		{
			name: "Request::NoParams",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10, \"method\": \"init\"}"),
			want: rpc.RequestMessage{
				ID:     id,
				Method: "init",
			},
			wantErr: false,
		},
		{
			name: "Notification",
			data: []byte("{\"jsonrpc\": \"2.0\", \"method\": \"init\", \"params\": {\"hello\":\"world\"}}"),
			want: rpc.NotificationMessage{
				Method: "init",
				Params: json.RawMessage([]byte("{\"hello\":\"world\"}")),
			},
			wantErr: false,
		},
		{
			name:    "Notification::NoParams",
			data:    []byte("{\"jsonrpc\": \"2.0\", \"method\": \"init\"}"),
			want:    rpc.NotificationMessage{Method: "init"},
			wantErr: false,
		},
		{
			name:    "Response::Result",
			data:    []byte("{\"jsonrpc\": \"2.0\", \"id\": 10, \"result\": {\"init\": 0}}"),
			want:    rpc.ResponseMessage{ID: id, Result: json.RawMessage([]byte("{\"init\": 0}"))},
			wantErr: false,
		},
		{
			name: "Response::Error",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10, \"error\": {\"code\": 0, \"message\": \"hello\"}}"),
			want: rpc.ResponseMessage{
				ID: id,
				Error: &rpc.ResponseError{
					Code:    0,
					Message: "hello",
				},
			},
			wantErr: true,
		},
		{
			name: "InternalError",
			data: []byte("{\"jsonrpc\": \"2.0\", \"id\": 10}"),
			want: rpc.ResponseMessage{
				ID:    id,
				Error: &rpc.InternalError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := rpc.DecodeMessage(tt.data)
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
