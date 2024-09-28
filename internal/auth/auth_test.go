package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	hash1, _ := HashPassword("some")
	hash2, _ := HashPassword("some other")

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"match", args{"some", hash1}, true},
		{"match", args{"some", hash2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordHash(tt.args.password, tt.args.hash); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
