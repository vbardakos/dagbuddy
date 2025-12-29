package rpc

import (
	"testing"
)

func TestNewID(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		raw     any
		want    ID
		wantErr bool
	}{
		{
			name: "Float64",
			raw: float64(10.0),
			want: ID{value: int64(10)},
		},
		{
			name: "String",
			raw: "hello world",
			want: ID{value: "hello world"},
		},
		{
			name: "Integer",
			raw: 10,
			want: ID{value: int64(10)},
		},
		{
			name: "Invalid",
			raw: nil,
			want: ID{value: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewID(tt.raw)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				if gotErr.Error() != InternalError.Error() {
						t.Fatal("NewID() succeeded unexpectedly")
				}
				return
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("NewID() = %v, want %v", got, tt.want)
			}
		})
	}
}

