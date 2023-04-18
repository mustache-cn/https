package https

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MethodType request method
type MethodType string

type ContentType string

const (
	// Get method
	Get    MethodType = "GET"
	Post   MethodType = "POST"
	Put    MethodType = "PUT"
	Patch  MethodType = "PATCH"
	Delete MethodType = "DELETE"

	// JsonType content_type
	JsonType ContentType = "application/json"
	FormType ContentType = "application/x-www-form-urlencoded"

	// timeout const
	timeout = 30 * time.Second // second

)

// Client struct
type Client struct {
	url         string
	headers     map[string]string
	cookies     []*http.Cookie
	timeout     time.Duration
	params      map[string]string
	contentType ContentType
	method      MethodType

	// 'OtherType', you need set value
	body string
}

// NewClient new client instance
func NewClient(url string) *Client {
	return &Client{
		url:         url,
		headers:     map[string]string{},
		params:      map[string]string{},
		timeout:     timeout,
		contentType: JsonType,
	}
}

// SetHeaders set request headers
func (c *Client) SetHeaders(headers map[string]string) *Client {
	c.headers = headers
	return c
}

// AddHeader add request headers
func (c *Client) AddHeader(key, value string) *Client {
	c.headers[key] = value
	return c
}

// SetContentType set content-type
func (c *Client) SetContentType(contentType ContentType) *Client {
	switch contentType {
	case JsonType:
		c.headers["Content-Type"] = "application/json"
	case FormType:
		c.headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	c.contentType = contentType
	return c
}

// AddParam add param
func (c *Client) AddParam(key, value string) *Client {
	c.params[key] = value
	return c
}

// SetCookies set request cookies
func (c *Client) SetCookies(cookies []*http.Cookie) *Client {
	c.cookies = cookies
	return c
}

// SetTimeout set request timeout
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

// SetBody set request body
func (c *Client) SetBody(body string) *Client {
	c.body = body
	return c
}

func (c *Client) setMethodType(method MethodType) *Client {
	c.method = method
	return c
}

// Get request get
func (c *Client) Get() (*Response, error) {
	c.setMethodType(Get)
	return c.do()
}

// Post request post
func (c *Client) Post() (*Response, error) {
	c.setMethodType(Post)
	return c.do()
}

// Put request put
func (c *Client) Put() (*Response, error) {
	c.setMethodType(Put)
	return c.do()
}

// Patch request patch
func (c *Client) Patch() (*Response, error) {
	c.setMethodType(Patch)
	return c.do()
}

// Delete request delete
func (c *Client) Delete() (*Response, error) {
	c.setMethodType(Delete)
	return c.do()
}

// do http request action
func (c *Client) do() (*Response, error) {
	reqUrl, reader, err := c.parse()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	request, err := http.NewRequest(string(c.method), reqUrl, reader)
	if err != nil {
		return nil, err
	}
	for k, v := range c.headers {
		request.Header.Add(k, v)
	}

	for _, cookie := range c.cookies {
		request.AddCookie(cookie)
	}

	httpClient := &http.Client{Timeout: c.timeout}
	res, err := httpClient.Do(request)
	return buildResponse(res, err)
}

func (c *Client) parse() (string, io.Reader, error) {
	switch c.method {
	case Get, Delete:
		res, err := c.parseQuery()
		return res, nil, err
	case Post, Put, Patch:
		return c.parseData()
	}
	return c.url, nil, errors.New("unsupported method")
}

func (c *Client) parseQuery() (string, error) {
	reqUrl := c.url
	paramsQuery := make(url.Values)
	for k, v := range c.params {
		paramsQuery.Add(k, v)
	}
	reqUrlObj, err := url.Parse(reqUrl)
	if err != nil {
		return reqUrl, err
	}
	reqUrlObj.RawQuery = paramsQuery.Encode()
	reqUrl = reqUrlObj.String()
	return reqUrl, nil
}

func (c *Client) parseData() (string, io.Reader, error) {
	if len(c.params) <= 0 {
		return c.url, nil, nil
	}
	switch c.contentType {
	case JsonType:
		b, err := json.Marshal(c.params)
		if err != nil {
			return c.url, nil, err
		}
		return c.url, bytes.NewReader(b), nil
	case FormType:
		postData := url.Values{}
		for k, v := range c.params {
			postData.Add(k, v)
		}
		return c.url, strings.NewReader(postData.Encode()), nil
	}
	return c.url, nil, errors.New("unsupported Content-Type")
}
