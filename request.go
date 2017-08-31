package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Errors used by this package.
var (
	ErrNotPOST        = errors.New("request: method is not POST when using form")
	ErrLackURL        = errors.New("request: request lacks URL")
	ErrLackMethod     = errors.New("request: request lacks method")
	ErrBodyAlreadySet = errors.New("request: request body has already been set")
	ErrStatusNotOk    = errors.New("request: status code is not ok (>= 400)")
)

type maxRedirects int

func (mr maxRedirects) check(req *http.Request, via []*http.Request) error {
	if len(via) >= int(mr) {
		return fmt.Errorf("request: exceed max redirects")
	}
	return nil
}

type basicAuthInfo struct {
	name     string
	password string
}

// Client is a HTTP client which provides usable and chainable methods.
type Client struct {
	cli       *http.Client
	req       *http.Request
	res       *Response
	method    string
	url       *url.URL
	queryVals url.Values
	formVals  url.Values
	mw        *multipart.Writer
	mwBuf     *bytes.Buffer
	body      io.Reader
	basicAuth *basicAuthInfo
	header    http.Header
	cookies   []*http.Cookie
	timeout   time.Duration
	redirects maxRedirects
	err       error
}

// New returns a new instance of Client.
func New() *Client {
	c := &Client{
		cli:      new(http.Client),
		header:   make(http.Header),
		formVals: make(url.Values),
		cookies:  make([]*http.Cookie, 0),
		mwBuf:    bytes.NewBuffer(nil),
	}
	c.mw = multipart.NewWriter(c.mwBuf)

	return c
}

// To defines the method and URL of the request.
func (c *Client) To(method string, URL string) *Client {
	c.method = method
	u, err := url.Parse(URL)

	if err != nil {
		c.err = err
		return c
	}

	c.url = u
	c.queryVals = u.Query()

	return c
}

// SetHeader sets the request header entries associated with key to the single
// element value. It replaces any existing values associated with key.
func (c *Client) SetHeader(key, value string) *Client {
	c.header.Set(key, value)

	return c
}

// AddHeader adds the key, value pair to the request header.It appends to any
// existing values associated with key.
func (c *Client) AddHeader(key, value string) *Client {
	c.header.Add(key, value)

	return c
}

// SetHeaders sets all key, value pairs in h to the request header, it replaces any
// existing values associated with key.
func (c *Client) SetHeaders(h http.Header) *Client {
	c.header = h
	return c
}

var typesMap = map[string]string{
	"html":       "text/html",
	"json":       "application/json",
	"xml":        "application/xml",
	"text":       "text/plain",
	"urlencoded": "application/x-www-form-urlencoded",
	"form":       "application/x-www-form-urlencoded",
	"form-data":  "application/x-www-form-urlencoded",
	"multipart":  "multipart/form-data",
}

// SetContentType sets the "Content-Type" request header to the given value.
// Some shorthands are supported:
//
// "html":       "text/html"
// "json":       "application/json"
// "xml":        "application/xml"
// "text":       "text/plain"
// "urlencoded": "application/x-www-form-urlencoded"
// "form":       "application/x-www-form-urlencoded"
// "form-data":  "application/x-www-form-urlencoded"
// "multipart":  "multipart/form-data"
//
// So you can just call .Type("html") to set the "Content-Type"
// header to "text/html".
func (c *Client) SetContentType(t string) *Client {
	if typ, ok := typesMap[strings.TrimSpace(strings.ToLower(t))]; ok {
		return c.SetHeader("Content-Type", typ)
	}

	return c.SetHeader("Content-Type", t)
}

// Accept sets the "Accept" request header to the given value.
// Some shorthands are supported:
//
// "html":       "text/html"
// "json":       "application/json"
// "xml":        "application/xml"
// "text":       "text/plain"
// "urlencoded": "application/x-www-form-urlencoded"
// "form":       "application/x-www-form-urlencoded"
// "form-data":  "application/x-www-form-urlencoded"
// "multipart":  "multipart/form-data"
//
// So you can just call .Accept("json") to set the "Accept"
// header to "application/json".
func (c *Client) Accept(t string) *Client {
	if typ, ok := typesMap[strings.TrimSpace(strings.ToLower(t))]; ok {
		return c.SetHeader("Accept", typ)
	}

	return c.SetHeader("Accept", t)
}

// Query adds the the given value to request's URL query-string.
func (c *Client) Query(vals url.Values) *Client {
	for k, vs := range vals {
		for _, v := range vs {
			c.queryVals.Add(k, v)
		}
	}

	return c
}

// SendBody sends the body in JSON format, body can be anything which can be
// Marshaled or just Marshaled JSON string.
func (c *Client) SendBody(body interface{}) *Client {
	if c.body != nil || c.mwBuf.Len() != 0 {
		c.err = ErrBodyAlreadySet
		return c
	}

	switch body := body.(type) {
	case string:
		c.body = bytes.NewBufferString(body)
	default:
		j, err := json.Marshal(body)

		if err != nil {
			c.err = err
			return c
		}

		c.body = bytes.NewReader(j)
	}

	c.SetContentType("json")
	return c
}

// Cookies adds get cookie from the response.
func (c *Client) Cookies() []*http.Cookie {
	return c.res.Cookies()
}

// AddCookie adds the cookie to the request.
func (c *Client) AddCookie(cookie *http.Cookie) *Client {
	c.cookies = append(c.cookies, cookie)

	return c
}

// SetCookies adds get cookie from the response.
func (c *Client) SetCookies(cookies []*http.Cookie) *Client {
	for _, cookie := range cookies {
		c.AddCookie(cookie)
	}

	return c
}

// AddCookieJar adds all cookies in the cookie jar to the request.
func (c *Client) AddCookieJar(jar http.CookieJar) *Client {
	for _, cookie := range jar.Cookies(c.url) {
		c.AddCookie(cookie)
	}

	return c
}

// SetTimeout specifies a time limit for the request.
// The timeout includes connection time, any
// redirects, and reading the response body. The timer remains
// running after Get, Head, Post, or End return and will
// interrupt reading of the response body.
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.cli.Timeout = timeout

	return c
}

// Redirects sets the max redirects count for the request.
// If not set, request will use its default policy,
// which is to stop after 10 consecutive requests.
func (c *Client) Redirects(count int) *Client {
	c.cli.CheckRedirect = maxRedirects(count).check

	return c
}

// SetAuth sets the request's Authorization header to use HTTP Basic
// Authentication with the provided username and password.
//
// With HTTP Basic Authentication the provided username and password are not
// encrypted.
func (c *Client) SetAuth(name, password string) *Client {
	c.basicAuth = &basicAuthInfo{name: name, password: password}

	return c
}

// AddFields sets the field values like form fields in HTML. Once it was set,
// the "Content-Type" header of the request will be automatically set to
// "application/x-www-form-urlencoded".
func (c *Client) AddFields(vals url.Values) *Client {
	for k, vs := range vals {
		for _, v := range vs {
			c.formVals.Add(k, v)
		}
	}

	c.SetContentType("application/x-www-form-urlencoded")
	return c
}

// AttachFile adds the attachment file to the form. Once the attachment was
// set, the "Content-Type" will be set to "multipart/form-data; boundary=xxx"
// automatically.
func (c *Client) AttachFile(fieldname, path, filename string) *Client {
	if c.body != nil {
		c.err = ErrBodyAlreadySet
		return c
	}

	file, err := os.Open(path)

	if err != nil {
		c.err = err
		return c
	}

	fw, err := c.mw.CreateFormFile(fieldname, filename)

	if err != nil {
		c.err = err
		return c
	}

	if _, err = io.Copy(fw, file); err != nil {
		c.err = err
		return c
	}

	return c
}

// Execute sends the HTTP request and returns the HTTP reponse.
//
// An error is returned if caused by client policy (such as timeout), or
// failure to speak HTTP (such as a network connectivity problem), or generated
// by former chained methods. A non-2xx status code doesn't cause an error.
func (c *Client) Execute() (*Response, error) {
	if c.url == nil {
		return nil, ErrLackURL
	}

	if c.method == "" {
		return nil, ErrLackMethod
	}

	if c.err != nil || c.res != nil {
		return c.res, c.err
	}

	if err := c.assemble(); err != nil {
		c.err = err
		return nil, err
	}

	response, err := c.cli.Do(c.req)

	if err != nil {
		c.err = err
		return nil, err
	}

	c.res = &Response{Response: response}

	return c.res, nil
}

// Req returns the representing http.Request instance of this request.
// It is often used in wirting tests.
func (c *Client) Req() (*http.Request, error) {
	if c.url == nil {
		return nil, ErrLackURL
	}

	if c.method == "" {
		return nil, ErrLackMethod
	}

	if c.err != nil {
		return nil, c.err
	}

	if err := c.assemble(); err != nil {
		c.err = err
		return nil, err
	}

	return c.req, nil
}

// JSON sends the HTTP request and returns the reponse body with JSON format.
func (c *Client) JSON(v ...interface{}) (interface{}, error) {
	if _, err := c.Execute(); err != nil {
		return nil, err
	}

	return c.res.JSON(v...)
}

// Text sends the HTTP request and returns the reponse body with text format.
func (c *Client) Text() (string, error) {
	if _, err := c.Execute(); err != nil {
		return "", err
	}

	return c.res.Text()
}

// Bytes sends the HTTP request and returns the reponse body with []byte format.
func (c *Client) Bytes() ([]byte, error) {
	if _, err := c.Execute(); err != nil {
		return nil, err
	}

	return c.res.Content()
}

// Dump Dump request
func (c *Client) Dump() error {
	// c.res.Dump()
	return nil
}

// ToFile download file to local
func (c *Client) ToFile(dir, fileName string) error {

	if _, err := c.Execute(); err != nil {
		return err
	}

	if fileName == "" {
		part := multipart.Part{Header: textproto.MIMEHeader(c.res.Header)}
		fileName = part.FileName()

		if fileName == "" {
			var surfix string
			exts, _ := mime.ExtensionsByType(c.res.ContentType())

			if len(exts) > 0 {
				surfix = exts[0]
			} else {
				surfix = GetContentTypeSufix(c.res.ContentType())
			}

			if c.req.URL.Path == "" {
				fileName = time.Now().Format("20060102150405")
			} else {
				fileName = path.Base(c.req.URL.Path)
			}

			if surfix != "" && !strings.Contains(fileName, surfix) {
				fileName = fmt.Sprintf("%s%s", fileName, surfix)
			}
		}
	}

	contentLength := c.res.Header.Get("Content-Length")

	fmt.Printf("header=%+v\n", c.res.Header)
	fmt.Printf("fileName=%+v, contentLength=%v\n", fileName, contentLength)

	filePath := path.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, c.res.Body)
	return err
}

func (c *Client) assemble() error {
	c.url.RawQuery = c.queryVals.Encode()

	if Settings().UserAgent != "" {
		c.SetHeader("User-Agent", Settings().UserAgent)
	}

	if Settings().ProxyTransport != nil {
		c.cli.Transport = Settings().ProxyTransport
	}

	var buf io.Reader

	if c.mwBuf.Len() != 0 {
		if c.formVals != nil {
			for k, vs := range c.formVals {
				for _, v := range vs {
					if err := c.mw.WriteField(k, v); err != nil {
						return err
					}
				}
			}
		}

		buf = c.mwBuf
		c.SetContentType(c.mw.FormDataContentType())
		c.mw.Close()
	} else if c.formVals != nil && c.body == nil {
		buf = strings.NewReader(c.formVals.Encode())
	} else {
		buf = c.body
	}

	req, err := http.NewRequest(c.method, c.url.String(), buf)

	if err != nil {
		return err
	}

	c.req = req
	c.req.Header = c.header

	if c.basicAuth != nil {
		c.req.SetBasicAuth(c.basicAuth.name, c.basicAuth.password)
	}

	for _, cookie := range c.cookies {
		c.req.AddCookie(cookie)
	}

	return nil
}
