package httpclient

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Response represents the response from a HTTP request.
type Response struct {
	*http.Response

	raw     *bytes.Buffer
	content []byte
}

// Raw returns the raw bytes body of the response.
func (r *Response) Raw() ([]byte, error) {
	if r.raw != nil {
		return r.raw.Bytes(), nil
	}

	b, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		return nil, err
	}

	r.raw = bytes.NewBuffer(b)

	return b, nil
}

// ContentType returns the content-type of the response.
func (r *Response) ContentType() string {

	return r.Header.Get("Content-Type")
}

// Content returns the content of the response body, it will handle
// the compression.
func (r *Response) Content() ([]byte, error) {
	if r.content != nil {
		return r.content, nil
	}

	rawBytes, err := r.Raw()

	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser

	switch r.Header.Get("ContentEncoding") {
	case "gzip":
		if reader, err = gzip.NewReader(bytes.NewBuffer(r.raw.Bytes())); err != nil {
			return nil, err
		}
	case "deflate":
		reader = flate.NewReader(bytes.NewBuffer(r.raw.Bytes()))
	}

	if reader == nil {
		r.content = rawBytes

		return rawBytes, nil
	}

	defer reader.Close()
	b, err := ioutil.ReadAll(reader)

	// If gzip or deflate decoding failed, try zlib decoding instead.
	// The body may be wrapped in the zlib data format.
	if err != nil {
		var zlibReader io.ReadCloser

		if zlibReader, err = zlib.NewReader(bytes.NewBuffer(r.raw.Bytes())); err != nil {
			return nil, err
		}
		defer zlibReader.Close()

		if b, err = ioutil.ReadAll(zlibReader); err != nil {
			return nil, err
		}
	}

	r.content = b

	return b, nil
}

// JSON returns the response body with JSON format.
func (r *Response) JSON(v ...interface{}) (interface{}, error) {
	b, err := r.Content()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(r.Header.Get("ContentType"), "application/json") {
		err := r.Status
		if len(b) > 0 {
			err = string(b)
		}
		return nil, errors.New(err)
	}

	var res interface{}
	if len(v) > 0 {
		res = v[0]
	} else {
		res = new(map[string]interface{})
	}

	if err = json.Unmarshal(b, res); err != nil {
		return nil, err
	}

	if !r.OK() {
		return res, ErrStatusNotOk{statusCode: r.StatusCode}
	}

	return res, nil
}

// Text returns the response body with text format.
func (r *Response) Text() (string, error) {
	b, err := r.Content()

	if err != nil {
		return "", err
	}

	if !r.OK() {
		return string(b), ErrStatusNotOk{statusCode: r.StatusCode}
	}

	return string(b), nil
}

// URL returns url of the final request.
func (r *Response) URL() (*url.URL, error) {
	u := r.Request.URL

	if r.StatusCode == http.StatusMovedPermanently ||
		r.StatusCode == http.StatusFound ||
		r.StatusCode == http.StatusSeeOther ||
		r.StatusCode == http.StatusTemporaryRedirect {
		location, err := r.Location()

		if err != nil {
			return nil, err
		}

		u = u.ResolveReference(location)
	}

	return u, nil
}

// Reason returns the status text of the response status code.
func (r *Response) Reason() string {
	return http.StatusText(r.StatusCode)
}

// OK returns whether the response status code is less than 400.
func (r *Response) OK() bool {
	return r.StatusCode < 400
}
