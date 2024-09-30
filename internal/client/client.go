// Package client содержит методы для создания http-запросов от клиента на сервер
package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"

	"github.com/wellywell/gophkeeper/internal/config"
	"github.com/wellywell/gophkeeper/internal/types"
)

const Token = "X-Auth-Token"

// Client тип для работы с http-клиетом
type Client struct {
	address string
	client  *http.Client
}

// NewClient инициализирует клиент
func NewClient(conf *config.ClientConfig) (*Client, error) {
	caCert, err := os.ReadFile(conf.SSLKey)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: true,
			},
		},
	}
	return &Client{
		address: conf.ServerAddress,
		client:  client,
	}, nil

}

// Login авторизация пользователя на сервере и получение токена для последующих запросов
func (c *Client) Login(login string, password string) (string, error) {
	return c.getAuthToken(login, password, "login")
}

// Register регистрация пользователя на сервере и получение токена для последующих запросов
func (c *Client) Register(login string, password string) (string, error) {
	return c.getAuthToken(login, password, "register")
}

// CreateBinaryItem сохранение на сервере бинарных данных
func (c *Client) CreateBinaryItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("%s/api/item/binary", c.address), http.MethodPost, data, headers)
}

// UpdateBinaryItem обновление бинарных данных, хранимых на сервере
func (c *Client) UpdateBinaryItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("%s/api/item/binary", c.address), http.MethodPut, data, headers)
}

// CreateLoginPasswordItem сохранение на сервере данных типа "логин и пароль"
func (c *Client) CreateLoginPasswordItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("%s/api/item/login_password", c.address), http.MethodPost, data, headers)
}

// CreateCreditCardItem сохранение на сервере данных кредитной карты
func (c *Client) CreateCreditCardItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("%s/api/item/credit_card", c.address), http.MethodPost, data, headers)
}

// CreateTextItem сохранение на сервере текстовых данных
func (c *Client) CreateTextItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("%s/api/item/text", c.address), http.MethodPost, data, headers)
}

// GetItem получение с сервера данных произвольного типа (из числа поддерживаемых)
func (c *Client) GetItem(token string, key string) (data []byte, err error) {

	resp, err := c.doRequest(fmt.Sprintf("%s/api/item/%s", c.address, key), http.MethodGet, nil, map[string]string{Token: token})
	if err != nil {
		return nil, fmt.Errorf("could not make request %w", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching item %s %s", resp.Status, bodyBytes)
	}
	return bodyBytes, nil
}

// SeeRecords получение списка записей, хранимых на сервере
func (c *Client) SeeRecords(token string, pass string, page int, pageSize int) ([]types.Item, error) {

	resp, err := c.doRequest(fmt.Sprintf("%s/api/item/list?page=%d&limit=%d", c.address, page, pageSize), http.MethodGet, nil, map[string]string{Token: token})
	if err != nil {
		return nil, fmt.Errorf("could not make request %w", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching item %s %s", resp.Status, bodyBytes)
	}

	items := make([]types.Item, pageSize)

	err = json.Unmarshal(bodyBytes, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// DownloadBinaryData cкачивание бинарных данных с сервера
func (c *Client) DownloadBinaryData(token string, pass string, key string) (data []byte, err error) {
	resp, err := c.doRequest(fmt.Sprintf("%s/api/item/binary/%s/download", c.address, key), http.MethodGet, nil, map[string]string{Token: token})
	if err != nil {
		return nil, fmt.Errorf("could not make request %w", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching item %s %s", resp.Status, bodyBytes)
	}

	result := types.BinaryData(bodyBytes)
	err = result.Decrypt(pass)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteItem удаление данных с сервера
func (c *Client) DeleteItem(token string, key string) error {

	resp, err := c.doRequest(fmt.Sprintf("%s/api/item/%s", c.address, key), http.MethodDelete, nil, map[string]string{Token: token})
	if err != nil {
		return fmt.Errorf("could not make request %w", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting item %s %s", resp.Status, bodyBytes)
	}
	return nil
}

// UpdateLogoPassData обновление логина и пароля, хранимых на сервере
func (c *Client) UpdateLogoPassData(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("%s/api/item/login_password", c.address), http.MethodPut, data, headers)
}

// UpdateCreditCardData обновление данных кредитной карты
func (c *Client) UpdateCreditCardData(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("%s/api/item/credit_card", c.address), http.MethodPut, data, headers)
}

// UpdateTextData обновление текстовых данных
func (c *Client) UpdateTextData(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("%s/api/item/text", c.address), http.MethodPut, data, headers)
}

// UpdateItem обобщенный метод для обновления данных типа T
func UpdateItem[T types.ItemData](token string, pass string, newItem types.GenericItem[T], method func([]byte, map[string]string) (*http.Response, error)) error {

	resp, err := saveItem(token, pass, newItem, method)

	if err != nil {
		return fmt.Errorf("could not make request %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		return fmt.Errorf("error updating item %s %s", resp.Status, bodyBytes)
	}
	return nil
}

// CreateItem обобщенный метод для сохранения на сервере данных типа T
func CreateItem[T types.ItemData](token string, pass string, newItem types.GenericItem[T], method func([]byte, map[string]string) (*http.Response, error)) error {
	resp, err := saveItem(token, pass, newItem, method)

	if err != nil {
		return fmt.Errorf("could not make request %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		return fmt.Errorf("error creating item %s %s", resp.Status, bodyBytes)
	}
	return nil
}

func (c *Client) getAuthToken(login string, password string, method string) (string, error) {

	data := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    login,
		Password: password,
	}

	request, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("could not serialize data")
	}

	resp, err := c.client.Post(fmt.Sprintf("%s/api/user/%s", c.address, method), "application/json", bytes.NewBuffer(request))

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("not authenticated")
	}

	token := resp.Header.Get(Token)
	if token == "" {
		return "", fmt.Errorf("empty token")
	}
	return token, nil
}

func saveJSONItem[T types.ItemData](token string, pass string, newItem types.GenericItem[T], method func([]byte, map[string]string) (*http.Response, error)) (*http.Response, error) {
	err := newItem.Data.Encrypt(pass)

	if err != nil {
		return nil, fmt.Errorf("could not encrypt %w", err)
	}
	data, err := json.Marshal(newItem)
	if err != nil {
		return nil, fmt.Errorf("could not serialize data")
	}
	headers := map[string]string{
		Token:          token,
		"Content-Type": "application/json",
	}
	return method(data, headers)
}

func saveBinaryItem[T types.BinaryData](token string, pass string, newItem types.GenericItem[*types.BinaryData], method func([]byte, map[string]string) (*http.Response, error)) (*http.Response, error) {
	err := newItem.Data.Encrypt(pass)

	if err != nil {
		return nil, fmt.Errorf("could not encrypt %w", err)
	}

	metadata, err := json.Marshal(newItem.Item)
	if err != nil {
		return nil, fmt.Errorf("could not convert %w", err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Metadata part
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json")
	// Create new multipart part
	part, err := writer.CreatePart(metadataHeader)
	if err != nil {
		return nil, err
	}
	// Write the part body
	_, err = part.Write(metadata)
	if err != nil {
		return nil, err
	}
	// Media part
	mediaHeader := textproto.MIMEHeader{}
	mediaHeader.Set("Content-Type", "application/octet-stream")

	mediaPart, err := writer.CreatePart(mediaHeader)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(mediaPart, bytes.NewReader(*newItem.Data))
	if err != nil {
		return nil, err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type":   fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary()),
		"Content-Length": fmt.Sprintf("%d", body.Len()),
		Token:            token,
	}

	return method(body.Bytes(), headers)
}

func saveItem[T types.ItemData](token string, pass string, newItem types.GenericItem[T], method func([]byte, map[string]string) (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error

	switch newItem.Item.Type {
	case types.TypeBinary:
		binaryItem, ok := any(newItem).(types.GenericItem[*types.BinaryData])
		if !ok {
			return nil, fmt.Errorf("failed to convert binary data")
		}
		resp, err = saveBinaryItem(token, pass, binaryItem, method)
	default:
		resp, err = saveJSONItem(token, pass, newItem, method)
	}
	return resp, err
}

func (c *Client) doRequest(URL string, method string, data []byte, headers map[string]string) (*http.Response, error) {
	body := bytes.NewBuffer(data)

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request")
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}
	return c.client.Do(req)

}
