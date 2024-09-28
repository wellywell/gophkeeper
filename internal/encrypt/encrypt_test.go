// Package encrypt пакет с методами шифрования данных
package encrypt

import (
	"reflect"
	"testing"
)

func TestEncrypt(t *testing.T) {
	type args struct {
		text     string
		MySecret string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"encrypt",args{"текст", "secret"}, "3xTEJ/wSsD9FGA==", false},
		{"empty",args{"", "secret"}, "", false},
		{"empty secret",args{"текст", ""}, "xXbwKfNOD++P0w==", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.text, tt.args.MySecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryptBytes(t *testing.T) {
	type args struct {
		data     []byte
		MySecret string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"encrypt",args{[]byte("a"), "secret"}, []byte{byte(111)}, false},
		{"empty",args{[]byte{}, "secret"}, []byte{}, false},
		{"empty secret",args{[]byte("a"), ""}, []byte{byte(117)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptBytes(tt.args.data, tt.args.MySecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncryptBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		text     string
		MySecret string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"decrypt",args{"3xTEJ/wSsD9FGA==", "secret"}, "текст", false},
		{"empty",args{"", "secret"}, "", false},
		{"empty secret",args{"xXbwKfNOD++P0w==", ""}, "текст", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.text, tt.args.MySecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecryptBytes(t *testing.T) {
	type args struct {
		data     []byte
		MySecret string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"decrypt", args{[]byte{byte(111)}, "secret"}, []byte("a"), false},
		{"empty",args{[]byte{}, "secret"}, []byte{}, false},
		{"empty secret",args{[]byte{byte(117)}, ""}, []byte("a"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptBytes(tt.args.data, tt.args.MySecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
