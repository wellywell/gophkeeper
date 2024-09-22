package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wellywell/gophkeeper/internal/encrypt"
)

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

type CreditCardData struct {
	Number     string    `json:"number" db:"number"`
	ValidMonth string    `json:"valid_month" db:"-"`
	ValidYear  string    `json:"valid_year" db:"-"`
	Name       string    `json:"name" db:"owner_name"`
	CVC        string    `json:"cvc" db:"cvc"`
	ValidDate  time.Time `db:"valid_till"`
}

func (c *CreditCardData) Encrypt(key string) error {
	num, err := encrypt.Encrypt(c.Number, key)

	if err != nil {
		return err
	}
	name, err := encrypt.Encrypt(c.Name, key)
	if err != nil {
		return err
	}

	cvc, err := encrypt.Encrypt(c.CVC, key)
	if err != nil {
		return err
	}

	c.Number = num
	c.Name = name
	c.CVC = cvc

	return nil
}

func (c *CreditCardData) Decrypt(key string) error {
	num, err := encrypt.Decrypt(c.Number, key)

	if err != nil {
		return err
	}
	name, err := encrypt.Decrypt(c.Name, key)
	if err != nil {
		return err
	}

	cvc, err := encrypt.Decrypt(c.CVC, key)
	if err != nil {
		return err
	}

	c.Number = num
	c.Name = name
	c.CVC = cvc

	return nil
}

func (c *CreditCardData) String() string {
	return fmt.Sprintf("\nNumber: %s\nValid: %s/%s\nName: %s\nCVC: %s\n", c.Number, c.ValidMonth, c.ValidYear, c.Name, c.CVC)
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

type CreditCardItem struct {
	Item Item            `json:"item"`
	Data *CreditCardData `json:"data"`
}
type LoginPasswordItem struct {
	Item Item           `json:"item"`
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
		return nil, err
	}
	return item, nil
}
