package client

import (
	"net/http"
	"reflect"
	"testing"
)

func TestClient_Login(t *testing.T) {
	type fields struct {
		address string
		client  *http.Client
	}
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				address: tt.fields.address,
				client:  tt.fields.client,
			}
			got, err := c.Login(tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Register(t *testing.T) {
	type fields struct {
		address string
		client  *http.Client
	}
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				address: tt.fields.address,
				client:  tt.fields.client,
			}
			got, err := c.Register(tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}
