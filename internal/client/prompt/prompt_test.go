package prompt

import (
	"reflect"
	"testing"
)

func TestEnterKey(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnterKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnterKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EnterKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnterMetadata(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnterMetadata()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnterMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EnterMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnterLoginPassword(t *testing.T) {
	tests := []struct {
		name    string
		want    *LoginPassword
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnterLoginPassword()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnterLoginPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnterLoginPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnterText(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnterText()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnterText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EnterText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnterFile(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnterFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnterFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EnterFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnterCreditCardData(t *testing.T) {
	tests := []struct {
		name    string
		want    *CreditCardData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnterCreditCardData()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnterCreditCardData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnterCreditCardData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChooseDataType(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChooseDataType()
			if (err != nil) != tt.wantErr {
				t.Errorf("ChooseDataType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChooseDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChooseLoginOrRegister(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChooseLoginOrRegister()
			if (err != nil) != tt.wantErr {
				t.Errorf("ChooseLoginOrRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChooseLoginOrRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	type args struct {
		method func(string, string) (string, error)
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
			got, err := Authenticate(tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}
