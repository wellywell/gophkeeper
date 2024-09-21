package types

import (
	"encoding/json"
	"fmt"

	"github.com/wellywell/gophkeeper/internal/encrypt"
)

type CreditCardData struct {
	Number string `json:"number"`
	Valid  string `json:"valid"`
	Name   string `json:"name"`
	CVC    string `json:"cvc"`
}

type ItemType string

const (
	TypeCreditCard ItemType = "credit_card"
	TypeText       ItemType = "text"
	TypeBinary     ItemType = "binary"
	TypeLogoPass   ItemType = "logopass"
)

type Item struct {
	Id   int      `db:"id"`
	Key  string   `json:"key" db:"key"`
	Info string   `json:"info" db:"info"`
	Type ItemType `json:"type" db:"item_type"`
}
func (i Item) String() string {
	return fmt.Sprintf("\nKey: %s\nInfo: %s\nType: %s\n", i.Key, i.Info, i.Type)
}

type CreditCardItem struct {
	Item Item           `json:"item"`
	Data CreditCardData `json:"data"`
}

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (l *LoginPassword) Encrypt(key string) error {
	lg, err := encrypt.Encrypt(l.Login, key)

	if err != nil {
		return err
	}
	psswd, err := encrypt.Encrypt(l.Password, key)

	if err != nil {
		return err
	}
	l.Login = lg
	l.Password = psswd

	return nil
}

func (l *LoginPassword) Decrypt(key string) error {
	lg, err := encrypt.Decrypt(l.Login, key)

	if err != nil {
		return err
	}
	psswd, err := encrypt.Decrypt(l.Password, key)

	if err != nil {
		return err
	}
	l.Login = lg
	l.Password = psswd

	return nil
}

func (l *LoginPassword) String() string {
	return fmt.Sprintf("\nLogin: %s\nPassword: %s\n", l.Login, l.Password)
}

type LoginPasswordItem struct {
	Item Item          `json:"item"`
	Data *LoginPassword `json:"data"`
}

type TextItem struct {
	Item Item   `json:"item"`
	Data string `json:"data"`
}

type BinaryItem struct {
	Item Item   `json:"item"`
	Data []byte `json:"data"`
}

type AnyItem struct {
	Item Item        `json:"item"`
	Data interface{} `json:"data"`
}

type ItemData interface {
	*LoginPassword | *CreditCardData | ~string | ~[]byte
	String() string
	Encrypt(string) error
	Decrypt(string) error
}

type GenericItem[T ItemData] struct {
	Item Item `json:"item"`
	Data T    `json:"data"`
}

func ParseItem[T ItemData](data []byte, decriptKey string) (*GenericItem[T], error) {
	var item *GenericItem[T]
	err := json.Unmarshal(data, &item)
	if err != nil {
		return nil, err
	}

	err = item.Data.Decrypt(decriptKey)
	if err != nil {
		return nil,  err
	}
	return item, nil
}
