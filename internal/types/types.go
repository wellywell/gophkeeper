package types

import "fmt"

type CreditCardData struct {
	Number string `json:"number"`
	Valid  string `json:"valid"`
	Name   string `json:"name"`
	CVC    string `json:"cvc"`
}

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

type CreditCardItem struct {
	Item Item           `json:"item"`
	Data CreditCardData `json:"data"`
}

type LoginPasswordItem struct {
	Item Item          `json:"item"`
	Data LoginPassword `json:"data"`
}

func (l LoginPasswordItem) String() string {
	return fmt.Sprintf("Key: %s\nInfo: %s\nType: %s\nLogin: %s\nPassword: %s\n", l.Item.Key, l.Item.Info, l.Item.Type, l.Data.Login, l.Data.Password)
}

type TextItem struct {
	Item Item   `json:"item"`
	Data string `json:"data"`
}

type BinaryItem struct {
	Item Item   `json:"item"`
	Data []byte `json:"data"`
}

type ItemData interface {
	LoginPassword | CreditCardData | ~string | ~[]byte
}

type GenericItem[T ItemData] struct {
	Item Item `json:"item"`
	Data []T  `json:"data"`
}

type AnyItem struct {
	Item Item        `json:"item"`
	Data interface{} `json:"data"`
}
