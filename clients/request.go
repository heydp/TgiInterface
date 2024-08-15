package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Payload struct {
	Text string `json:"text"`
}

type Response struct {
	Result string `json:"result"`
}

func (c *Client) GenerateTextResponse(req string) (*string, error) {
	queryParams := url.Values{}
	var payload = Payload{
		Text: req,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewReader(data)

	reqUrl := c.BaseUrlWithReqPath("", queryParams)

	fmt.Println("the req url is : ", reqUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl.String(), reqBody)
	if err != nil {
		return nil, err
	}

	httpResponse, err := c.Do(httpReq)
	if err != nil {
		return nil, err
	}

	var result Response
	err = c.ReadData(httpResponse, &result)
	if err != nil {
		return nil, err
	}

	return &result.Result, nil
}

func (c *Client) HealthCheck() (*string, error) {
	queryParams := url.Values{}

	reqUrl := c.BaseUrlWithReqPath("/ping", queryParams)

	fmt.Println("the req url is : ", reqUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	curTime := time.Now()
	httpResponse, err := c.Do(httpReq)
	if err != nil {
		return nil, err
	}
	responseTime := time.Since(curTime).String()
	var res = struct {
		Response string `json:"response"`
	}{
		Response: "",
	}

	err = c.ReadData(httpResponse, &res)
	if err != nil {
		return nil, err
	}
	return &responseTime, nil
}
