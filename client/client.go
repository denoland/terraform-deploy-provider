package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

var (
	baseUrl, _ = url.Parse("https://dash.deno.com")
)

type Client struct {
	HTTPClient *http.Client
	Token      string
}

type PageOptions struct {
	Page  int
	Limit int
}

type PagingInfo struct {
	Page       int `json:"page"`
	Count      int `json:"count"`
	Limit      int `json:"limit"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

func New(token string) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		Token:      token,
	}
}

func (c *Client) request(method, requestPath string, query url.Values, body io.Reader, responseStruct interface{}) error {
	r, err := c.newRequest(method, requestPath, query, body)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("status: %d, body: %v", resp.StatusCode, string(bodyContents))
	}

	if responseStruct == nil {
		return nil
	}

	err = json.Unmarshal(bodyContents, responseStruct)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(method, requestPath string, query url.Values, body io.Reader) (*http.Request, error) {
	url := baseUrl
	url.Path = path.Join(url.Path, requestPath)
	url.RawQuery = query.Encode()
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return req, err
	}

	if c.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	req.Header.Add("Content-Type", "application/json")
	return req, err
}
