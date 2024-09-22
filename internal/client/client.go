package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (c *Client) CreateLoginPasswordItem(token string, item types.LoginPasswordItem, encryptKey string) error {
	err := item.Data.Encrypt(encryptKey)

	if err != nil {
		return fmt.Errorf("could not encrypt %w", err)
	}
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("could not serialize data")
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/api/item/login_password", c.address), bytes.NewBuffer(data))

	if err != nil {
		return fmt.Errorf("could not create request")
	}
	req.Header.Set(Token, token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)

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

func (c *Client) CreateCreditCardItem(token string, item types.CreditCardItem, encryptKey string) error {
	err := item.Data.Encrypt(encryptKey)

	if err != nil {
		return fmt.Errorf("could not encrypt %w", err)
	}
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("could not serialize data")
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/api/item/credit_card", c.address), bytes.NewBuffer(data))

	if err != nil {
		return fmt.Errorf("could not create request")
	}
	req.Header.Set(Token, token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)

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

func (c *Client) GetItem(token string, key string) (data []byte, err error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/api/item/%s", c.address, key), nil)

	if err != nil {
		return nil, fmt.Errorf("could not create request")
	}
	req.Header.Set(Token, token)

	resp, err := c.client.Do(req)

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

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("https://%s/api/item/%s", c.address, key), nil)

	if err != nil {
		return fmt.Errorf("could not create request")
	}
	req.Header.Set(Token, token)

	resp, err := c.client.Do(req)
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

func (c *Client) UpdateLogoPassData(token string, data []byte) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://%s/api/item/login_password", c.address), bytes.NewBuffer(data))

	if err != nil {
		return nil, fmt.Errorf("could not create request")
	}
	req.Header.Set(Token, token)

	return c.client.Do(req)
}

func UpdateItem[T types.ItemData](token string, pass string, newItem types.GenericItem[T], method func(string, []byte) (*http.Response, error)) error {
	err := newItem.Data.Encrypt(pass)

	if err != nil {
		return fmt.Errorf("could not encrypt %w", err)
	}
	data, err := json.Marshal(newItem)
	if err != nil {
		return fmt.Errorf("could not serialize data")
	}

	resp, err := method(token, data)

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

func (c *Client) UpdateCreditCardData(token string, data []byte) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://%s/api/item/credit_card", c.address), bytes.NewBuffer(data))

	if err != nil {
		return nil, fmt.Errorf("could not create request")
	}
	req.Header.Set(Token, token)
	return c.client.Do(req)
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
