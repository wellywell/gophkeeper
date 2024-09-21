package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func adjustKeyLength(MySecret string) string {
	if len(MySecret) == 16 || len(MySecret) == 24 || len(MySecret) == 32 {
		return MySecret
	}
	if len(MySecret) < 32 {
		pad := make([]string, 32-len(MySecret))
		for i := range len(pad) {
			pad[i] = "0"
		}
		return fmt.Sprintf("%s%s", MySecret, strings.Join(pad, ""))
	} else {
		return string([]rune(MySecret)[:32])
	}
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, MySecret string) (string, error) {

	key := adjustKeyLength(MySecret)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, MySecret string) (string, error) {
	key := adjustKeyLength(MySecret)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
