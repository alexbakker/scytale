package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Response represents the base API response structure.
type Response struct {
	Success bool            `json:"success"`
	Error   string          `json:"error"`
	Content json.RawMessage `json:"content"`
}

// UploadResponse represents the structure of an upload response.
type UploadResponse struct {
	Filename string `json:"filename"`
}

type Client struct {
	key    Key
	client *http.Client
}

func NewClient(key Key) *Client {
	return &Client{
		key:    key,
		client: new(http.Client),
	}
}

func (c *Client) Upload(uri string, extension string, encrypted bool, r io.Reader) (*UploadResponse, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	query := url.Query()
	query.Set("ext", extension)
	query.Set("encrypted", strconv.FormatBool(encrypted))
	url.RawQuery = query.Encode()

	req, err := http.NewRequest("POST", url.String(), r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Key", c.key.String())

	httpRes, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	var rawRes Response
	if err = json.NewDecoder(httpRes.Body).Decode(&rawRes); err != nil {
		return nil, err
	}
	if !rawRes.Success {
		return nil, fmt.Errorf("request error: %s", rawRes.Error)
	}

	var res UploadResponse
	if err = json.Unmarshal(rawRes.Content, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
