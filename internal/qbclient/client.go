package qbclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// Client communicates with a qBittorrent Web API server.
type Client struct {
	host     string
	username string
	password string
	client   *http.Client
}

// New builds a client with an internal HTTP client and cookie jar.
func New(host, username, password string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		host:     strings.TrimRight(host, "/"),
		username: username,
		password: password,
		client:   &http.Client{Jar: jar},
	}
}

// NewWithClient builds a client using a provided http.Client (for testing).
func NewWithClient(host, username, password string, httpClient *http.Client) *Client {
	return &Client{
		host:     strings.TrimRight(host, "/"),
		username: username,
		password: password,
		client:   httpClient,
	}
}

// Login authenticates with qBittorrent and stores the session cookie.
func (c *Client) Login() error {
	form := url.Values{}
	form.Set("username", c.username)
	form.Set("password", c.password)

	req, err := http.NewRequest("POST", c.host+"/api/v2/auth/login", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

// AddMagnet sends a magnet URL to qBittorrent.
func (c *Client) AddMagnet(magnet string) error {
	if magnet == "" {
		return errors.New("magnet is empty")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	w, err := writer.CreateFormField("urls")
	if err != nil {
		return err
	}
	if _, err := io.WriteString(w, magnet); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.host+"/api/v2/torrents/add", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qBittorrent error: %s", strings.TrimSpace(string(body)))
	}

	return nil
}
