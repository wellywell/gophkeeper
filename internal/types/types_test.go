// Package types определяет основные типы данных для работы сервера и клиента
package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	textBody       = []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`)
	logopassBody   = []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`)
	logopassItem   = LoginPasswordItem{Item: Item{Key: "111", Type: TypeLogoPass}, Data: &LoginPassword{Login: "112", Password: "222"}}
	creditCardBody = []byte(`{"item": {"type": "credit_card", "key": "111"}, "data": {"cvc": "1", "number":"1", "name":"1", "valid_month": "1", "valid_year": "2000"}}`)
	creditCardItem = CreditCardItem{Item: Item{Key: "111", Type: TypeCreditCard}, Data: &CreditCardData{Number: "1", CVC: "1", Name: "1", ValidMonth: "1", ValidYear: "2000"}}
	textItem       = TextItem{Item: Item{Key: "111", Type: TypeText}, Data: TextData("text")}
)

func TestParseItem(t *testing.T) {
	type args struct {
		data       []byte
		decriptKey string
	}
	tests := []struct {
		name     string
		args     args
		itemType ItemType
	}{
		{"text", args{textBody, "secret"}, TypeText},
		{"credit card", args{creditCardBody, "secret"}, TypeCreditCard},
		{"logopass", args{logopassBody, "secret"}, TypeLogoPass},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			switch tt.itemType {
			case TypeText:
				copyItem := textItem
				_ = copyItem.Data.Encrypt(tt.args.decriptKey)
				body, _ := json.Marshal(copyItem)
				got, err := ParseItem[*TextData](body, tt.args.decriptKey)
				assert.NoError(t, err)
				assert.Equal(t, textItem.Item, got.Item)
				assert.Equal(t, textItem.Data, *got.Data)
			case TypeCreditCard:
				copyItem := CreditCardItem{Data: &CreditCardData{}}
				copyItem.Item = creditCardItem.Item
				*copyItem.Data = *creditCardItem.Data
				_ = copyItem.Data.Encrypt(tt.args.decriptKey)
				body, _ := json.Marshal(copyItem)
				got, err := ParseItem[*CreditCardData](body, tt.args.decriptKey)
				assert.NoError(t, err)
				assert.Equal(t, creditCardItem.Item, got.Item)
				assert.Equal(t, *creditCardItem.Data, *got.Data)
			case TypeLogoPass:
				copyItem := LoginPasswordItem{Data: &LoginPassword{}}
				copyItem.Item = logopassItem.Item
				*copyItem.Data = *logopassItem.Data
				_ = copyItem.Data.Encrypt(tt.args.decriptKey)
				body, _ := json.Marshal(copyItem)
				got, err := ParseItem[*LoginPassword](body, tt.args.decriptKey)
				assert.NoError(t, err)
				assert.Equal(t, logopassItem.Item, got.Item)
				assert.Equal(t, *logopassItem.Data, *got.Data)
			}
		})
	}
}

func TestBinaryData_Encrypt_Decrypt(t *testing.T) {

	text := "some text some text some text some text some text some text some text some text some text some text some text"

	expectEncrypted := []byte{
		125, 249, 121, 247, 12, 220, 4, 198, 224, 186, 123, 29, 124, 254, 76, 113, 138, 167, 125, 67, 84, 18, 124, 206, 229,
		3, 93, 63, 235, 126, 7, 61, 148, 82, 244, 239, 85, 76, 248, 149, 181, 114, 60, 77, 213, 94, 48, 133, 244, 164, 232,
		86, 148, 10, 220, 55, 247, 112, 80, 180, 113, 164, 147, 124, 74, 153, 194, 80, 214, 101, 121, 165, 101, 224, 241, 232,
		181, 5, 43, 23, 88, 199, 221, 189, 246, 211, 89, 156, 242, 118, 93, 182, 216, 94, 48, 142, 195, 190, 55, 178, 93, 218,
		112, 73, 183, 186, 122, 239, 140}

	data := BinaryData([]byte(text))
	err := data.Encrypt("secret")
	assert.NoError(t, err)
	assert.Equal(t, expectEncrypted, []byte(data))
	assert.NotEqual(t, []byte(text), []byte(data))

	err = data.Decrypt("secret")
	assert.NoError(t, err)
	assert.NotEqual(t, expectEncrypted, []byte(data))
	assert.Equal(t, []byte(text), []byte(data))

}

func TestTextData_String(t *testing.T) {
	assert.Equal(t, "text\n", textItem.Data.String())
}

func TestTextData_Encrypt_Decrypt(t *testing.T) {

	text := "some text some text some text some text some text some text some text some text some text some text some text"
	expectEncrypted := "ffl59wzcBMbgunsdfP5McYqnfUNUEnzO5QNdP+t+Bz2UUvTvVUz4lbVyPE3VXjCF9KToVpQK3Df3cFC0caSTfEqZwlDWZXmlZeDx6LUFKxdYx9299tNZnPJ2XbbYXjCOw743sl3acEm3unrvjA=="

	data := TextData(text)
	err := data.Encrypt("secret")
	assert.NoError(t, err)
	assert.Equal(t, expectEncrypted, string(data))
	assert.NotEqual(t, []byte(text), []byte(data))

	err = data.Decrypt("secret")
	assert.NoError(t, err)
	assert.NotEqual(t, expectEncrypted, data)
	assert.Equal(t, []byte(text), []byte(data))

}

func TestLoginPassword_String(t *testing.T) {
	assert.Equal(t, "\nLogin: 112\nPassword: 222\n", logopassItem.Data.String())
}

func TestLoginPassword_Encrypt_Decrypt(t *testing.T) {

	copyItem := LoginPasswordItem{Data: &LoginPassword{}}
	copyItem.Item = logopassItem.Item
	*copyItem.Data = *logopassItem.Data

	err := copyItem.Data.Encrypt("secret")
	assert.NoError(t, err)
	assert.Equal(t, "\nLogin: P6cm\nPassword: PKQm\n", copyItem.Data.String())
	assert.NotEqual(t, logopassItem.Data.String(), copyItem.Data.String())

	err = copyItem.Data.Decrypt("secret")
	assert.NoError(t, err)
	assert.NotEqual(t, "\nLogin: P6cm\nPassword: PKQm\n", copyItem.Data.String())
	assert.Equal(t, logopassItem.Data.String(), copyItem.Data.String())

}

func TestCreditCardData_String(t *testing.T) {
	assert.Equal(t, "\nNumber: 1\nValid: 1/2000\nName: 1\nCVC: 1\n", creditCardItem.Data.String())
}

func TestCreditCardData_Decrypt_Encrypt(t *testing.T) {
	copyItem := CreditCardItem{Data: &CreditCardData{}}
	copyItem.Item = creditCardItem.Item
	*copyItem.Data = *creditCardItem.Data

	expectEncrypted := "\nNumber: Pw==\nValid: 1/2000\nName: Pw==\nCVC: Pw==\n"

	err := copyItem.Data.Encrypt("secret")
	assert.NoError(t, err)
	assert.Equal(t, expectEncrypted, copyItem.Data.String())
	assert.NotEqual(t, creditCardItem.Data.String(), copyItem.Data.String())

	err = copyItem.Data.Decrypt("secret")
	assert.NoError(t, err)
	assert.NotEqual(t, expectEncrypted, copyItem.Data.String())
	assert.Equal(t, creditCardItem.Data.String(), copyItem.Data.String())

}

func TestItem_String(t *testing.T) {
	assert.Equal(t, "text\n", textItem.Data.String())
}
