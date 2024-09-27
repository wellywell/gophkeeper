// Package types определяет основные типы данных для работы сервера и клиента
package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wellywell/gophkeeper/internal/encrypt"
)

// ItemType определяет возможные типы данных для хранения на сервере
type ItemType string

const (
	TypeCreditCard ItemType = "credit_card"
	TypeText       ItemType = "text"
	TypeBinary     ItemType = "binary"
	TypeLogoPass   ItemType = "logopass"
)

// Item - структура для хранения метаданных о любом объекте, хранимом на сервере
type Item struct {
	Id   int      `db:"id"`
	Key  string   `json:"key" db:"key"`
	Info string   `json:"info" db:"info"`
	Type ItemType `json:"type" db:"item_type"`
}

// String метод для возвращения строкового представления Item
func (i Item) String() string {
	return fmt.Sprintf("\nKey: %s\nInfo: %s\nType: %s\n", i.Key, i.Info, i.Type)
}

// CreditCardData тип для хранения данных крединтных карт
type CreditCardData struct {
	Number     string    `json:"number" db:"number"`
	ValidMonth string    `json:"valid_month" db:"-"`
	ValidYear  string    `json:"valid_year" db:"-"`
	Name       string    `json:"name" db:"owner_name"`
	CVC        string    `json:"cvc" db:"cvc"`
	ValidDate  time.Time `db:"valid_till"`
}

// Encrypt зашифровывает данные кредитной карты перед отправкой на сервер
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

// Decrypt расшифровывает данные кредитной карты для клиента
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

// String строковое представлени данных о кредитной карте
func (c *CreditCardData) String() string {
	return fmt.Sprintf("\nNumber: %s\nValid: %s/%s\nName: %s\nCVC: %s\n", c.Number, c.ValidMonth, c.ValidYear, c.Name, c.CVC)
}

// LoginPassword структура для хранения пароля и логина
type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Encrypt зашифровывает пароль и логин перед передачей на сервер
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

// Decrypt расшифровывает пароль и логин, чтобы показать пользователю
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

// String строкове представление логина и пароля
func (l *LoginPassword) String() string {
	return fmt.Sprintf("\nLogin: %s\nPassword: %s\n", l.Login, l.Password)
}

// TextData тип для хранения простых текстовых данных
type TextData string


// Encrypt зашифровывает данные перед отправкой на сервер
func (t *TextData) Encrypt(key string) error {
	enc, err := encrypt.Encrypt(string(*t), key)

	if err != nil {
		return err
	}

	*t = TextData(enc)
	return nil
}

// Decrypt расшифровывает текстовые данные для передачи клиенту
func (t *TextData) Decrypt(key string) error {
	enc, err := encrypt.Decrypt(string(*t), key)

	if err != nil {
		return err
	}

	*t = TextData(enc)
	return nil
}

// String строковое представление текстовых данных
func (t *TextData) String() string {
	return fmt.Sprintf("%s\n", string(*t))
}

// BinaryData тип для хранения произвольных бинарных данных
type BinaryData []byte

// Encrypt зашифровывает данные перед отправкой на сервер
func (b *BinaryData) Encrypt(key string) error {
	enc, err := encrypt.EncryptBytes(*b, key)

	if err != nil {
		return err
	}

	*b = enc
	return nil
}

// Decrypt расшифровывает данные для показа клиенту
func (b *BinaryData) Decrypt(key string) error {
	dec, err := encrypt.EncryptBytes(*b, key)

	if err != nil {
		return err
	}

	*b = dec
	return nil
}

// String строковое представление, показываемое пользователю
func (b *BinaryData) String() string {
	return "Binary data"
}

// CreditCardItem структура для хранения данных о кредитной карте вместе с метаданными
type CreditCardItem struct {
	Item Item            `json:"item"`
	Data *CreditCardData `json:"data"`
}

// LoginPasswordItem тип для хранения логина пароля и метаданных
type LoginPasswordItem struct {
	Item Item           `json:"item"`
	Data *LoginPassword `json:"data"`
}

// TextItem тип для хранения текстовых данных и метаданных
type TextItem struct {
	Item Item     `json:"item"`
	Data TextData `json:"data"`
}

// BinaryItem тип для хранения бинарных данных и метаданных
type BinaryItem struct {
	Item Item   `json:"item"`
	Data []byte `json:"data"`
}

// AnyItem тип для передачи любого типа данных (из поддерживаемых), без уточнения конкретного типа
type AnyItem struct {
	Item Item        `json:"item"`
	Data interface{} `json:"data"`
}

// ItemData интерфейс, определяющий ограничения для обобщенного типа GenericItem
type ItemData interface {
	*LoginPassword | *CreditCardData | *TextData | *BinaryData
	String() string
	Encrypt(string) error
	Decrypt(string) error
}

// GenericItem обобщенный тип для хранения данных и метаданных
type GenericItem[T ItemData] struct {
	Item Item `json:"item"`
	Data T    `json:"data"`
}

// ParseItem преобразует массив байтов в типизированный GenericItem и расшифровывает зашифрованные данные
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
