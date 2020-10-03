package zetka

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_gatewayURI(t *testing.T) {
	type args struct {
		token    string
		version  string
		encoding Encoding
	}
	tests := []struct {
		name     string
		args     args
		want     string
		response []byte
		wantErr  bool
	}{
		{
			name: "failing - invalid encoding",
			args: args{
				token:    "foo",
				version:  "6",
				encoding: "invalid",
			},
			wantErr: true,
		},
		{
			name: "failing - unsupported etf encoding",
			args: args{
				token:    "foo",
				version:  "6",
				encoding: "etf",
			},
			wantErr: true,
		},
		{
			name: "passing - supported json encoding",
			args: args{
				token:    "foo",
				version:  "6",
				encoding: JSON,
			},
			want:     "wss://test.foo?encoding=json&v=6",
			response: []byte(`{"url":"wss://test.foo"}`),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tt.response)
			}))

			got, err := gatewayURI(srv.URL, tt.args.token, tt.args.version, tt.args.encoding)
			if (err != nil) != tt.wantErr {
				t.Errorf("gatewayURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("gatewayURI() got = %v, want %v", got, tt.want)
			}
		})
	}
}
