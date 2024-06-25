package rpc_test

import (
	"testing"

	"github.com/tjgurwara99/ruby-lsp/rpc"
)

func TestEncodeMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  any
		want string
	}{
		{
			name: "random struct",
			msg: struct {
				Haha string
			}{
				Haha: "content",
			},
			want: "Content-Length: 18\r\n\r\n{\"Haha\":\"content\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rpc.EncodeMessage(tt.msg); got != tt.want {
				t.Errorf("EncodeMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
