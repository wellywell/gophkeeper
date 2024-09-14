package auth

import (
	"net/http"
	"testing"
)

func TestVerifyUser(t *testing.T) {
	type args struct {
		r      *http.Request
		secret []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyUser(tt.args.r, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildJWTString(t *testing.T) {
	type args struct {
		user   string
		secret []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildJWTString(tt.args.user, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildJWTString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildJWTString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		tokenString string
		secret      []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.tokenString, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
