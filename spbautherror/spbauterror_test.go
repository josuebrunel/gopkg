package spbautherror

import (
	"testing"

	"github.com/josuebrunel/gopkg/assert"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		errorString string
		want        *Response
		wantErr     bool
	}{
		{
			name:        "valid error string with prefix",
			errorString: `API Error: {"code":400,"error_code":"bad_request","msg":"Invalid request payload"}`,
			want: &Response{
				Code:      400,
				ErrorCode: "bad_request",
				Msg:       "Invalid request payload",
			},
			wantErr: false,
		},
		{
			name:        "string is only json",
			errorString: `{"code":500,"error_code":"internal_error","msg":"An internal error occurred"}`,
			want: &Response{
				Code:      500,
				ErrorCode: "internal_error",
				Msg:       "An internal error occurred",
			},
			wantErr: false,
		},
		{
			name:        "string without json object",
			errorString: "A plain error message",
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "string with malformed json",
			errorString: `Error: {"code":400,"error_code":"bad_request",`,
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "empty string",
			errorString: "",
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.errorString)
			assert.Eq(t, (err != nil), tt.wantErr)
			assert.Eq(t, got, tt.want)
		})
	}
}
