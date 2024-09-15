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
)

type Client struct {
	address string
	client  *http.Client
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

	token := resp.Header.Get("X-Auth-Token")
	if token == "" {
		return "", fmt.Errorf("empty token")
	}
	return token, nil
}
