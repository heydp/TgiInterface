package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const (
	X_TOKEN = "x-token"
)

type Client struct {
	client   *http.Client
	host     string
	schema   string
	basePath string
	port     int64
	token    string
}

func NewClient(client *http.Client, schema string, host string, port int64, basePath string, token string) *Client {
	return &Client{
		client:   client,
		host:     host,
		schema:   schema,
		port:     port,
		basePath: basePath,
		token:    token,
	}
}

func (c *Client) Host() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *Client) Scheme() string {
	return c.schema
}

func (c *Client) BaseUrlWithReqPath(reqPath string, query url.Values) url.URL {
	queryString := ""
	if query != nil {
		queryString = query.Encode()
	}

	reqURL := url.URL{
		Scheme:   c.Scheme(),
		Host:     c.Host(),
		Path:     path.Join(c.basePath, reqPath),
		RawQuery: queryString,
	}

	return reqURL
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add(X_TOKEN, c.token)

	return c.client.Do(req)
}

func (c *Client) ReadData(httpResp *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	if httpResp.StatusCode < 199 || httpResp.StatusCode > 299 {
		return fmt.Errorf("received an error from API: %v", string(body))
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	return nil
}
