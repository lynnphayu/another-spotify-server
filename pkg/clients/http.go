package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HTTPClient struct {
	Timeout int
}

func New(timeout int) *HTTPClient {
	return &HTTPClient{Timeout: timeout}
}
func (client *HTTPClient) constructRequest(
	methodType string,
	URL string,
	body map[string]interface{},
	contentType string,
	auth string,
) (*http.Request, error) {
	var bodyByteArr string = ""
	if body != nil {
		if contentType == "application/x-www-form-urlencoded" {
			form := url.Values{}
			for k, v := range body {
				form.Add(k, v.(string))
			}
			bodyByteArr = form.Encode()
		} else {
			b, _ := json.Marshal(body)
			bodyByteArr = string(b)
		}
	}
	request, err := http.NewRequest(methodType, URL, strings.NewReader(bodyByteArr))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", auth)
	request.Header.Add("Content-Type", contentType)
	request.Header.Add("Content-Length", strconv.Itoa(len(bodyByteArr)))
	return request, nil
}

func (client *HTTPClient) getClient(sec int) *http.Client {
	timeout := time.Duration(time.Duration(sec) * time.Second)
	return &http.Client{
		Timeout: timeout,
	}
}

// Request - http request with parameters and return http response from endpoint
func (client *HTTPClient) Request(methodType string, URL string, body map[string]interface{}, contentType string, auth string) (*http.Response, error) {
	fmt.Print(URL)
	request, err := client.constructRequest(methodType, URL, body, contentType, auth)
	if err != nil {
		return nil, err
	}
	resp, err := client.getClient(client.Timeout).Do(request)
	if err != nil {
		return nil, err
	}
	// if resp.StatusCode != http.StatusOK {
	// 	return resp, errors.New("remote server error")
	// }
	return resp, err
}
