package handlers

import (
	"net/http"
	"testing"
)

func TestHandlerSet_HandleLogin(t *testing.T) {
	type fields struct {
		secret   []byte
		database Database
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleLogin(tt.args.w, tt.args.req)
		})
	}
}

func TestHandlerSet_HandleRegisterUser(t *testing.T) {
	type fields struct {
		secret   []byte
		database Database
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleRegisterUser(tt.args.w, tt.args.req)
		})
	}
}
