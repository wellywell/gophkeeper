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

	"github.com/wellywell/gophkeeper/internal/types"
)

const Token = "X-Auth-Token"

type Client struct {
	address string
	client  *http.Client
}

type ItemInfo struct {
	Data []byte
	Type types.ItemType
	View string
}

func NewClient(addr string) (*Client, error) {
	caCert, err := os.ReadFile("ca.key")
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
		address: addr,
		client:  client,
	}, nil

}

func (c *Client) Login(login string, password string) (string, error) {
	return c.getAuthToken(login, password, "login")
}

func (c *Client) Register(login string, password string) (string, error) {
	return c.getAuthToken(login, password, "register")
}

func (c *Client) doRequest(URL string, method string, data []byte, headers map[string]string) (*http.Response, error) {

	fmt.Println(data)

	body := bytes.NewBuffer(data)

    req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request")
	}

	for key, val := range(headers) {
		req.Header.Set(key, val)
	}
	return c.client.Do(req)

}

func (c *Client) CreateBinaryItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("https://%s/api/item/binary", c.address), http.MethodPost, data, headers)
}

func (c *Client) UpdateBinaryItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("https://%s/api/item/binary", c.address), http.MethodPut, data, headers)
}

func (c *Client) CreateLoginPasswordItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("https://%s/api/item/login_password", c.address), http.MethodPost, data, headers)
}

func (c *Client) CreateCreditCardItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("https://%s/api/item/credit_card", c.address), http.MethodPost, data, headers)
}

func (c *Client) CreateTextItem(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("https://%s/api/item/text", c.address), http.MethodPost, data, headers)
}

func (c *Client) GetItem(token string, key string) (data []byte, err error) {

	resp, err := c.doRequest(fmt.Sprintf("https://%s/api/item/%s", c.address, key), http.MethodGet, nil, map[string]string{Token: token})
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

func (c *Client) DeleteItem(token string, key string) error {

	resp, err := c.doRequest(fmt.Sprintf("https://%s/api/item/%s", c.address, key), http.MethodDelete, nil, map[string]string{Token: token})
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

func (c *Client) UpdateLogoPassData(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("https://%s/api/item/login_password", c.address), http.MethodPut, data, headers)
}

func (c *Client) UpdateCreditCardData(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("https://%s/api/item/credit_card", c.address), http.MethodPut, data, headers)
}

func (c *Client) UpdateTextData(data []byte, headers map[string]string) (*http.Response, error) {
	return c.doRequest(fmt.Sprintf("https://%s/api/item/text", c.address), http.MethodPut, data, headers)
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

	resp, err := c.client.Post(fmt.Sprintf("https://%s/api/user/%s", c.address, method), "application/json", bytes.NewBuffer(request))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
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
		Token: token,
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
	fmt.Println([]byte(*newItem.Data))
	_, err = io.Copy(mediaPart, bytes.NewReader(*newItem.Data))
	if err != nil {
		return nil, err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary()),
		"Content-Length": fmt.Sprintf("%d", body.Len()),
		Token: token,
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
