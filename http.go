package httpclient

import "net/http"
import "time"

// SetGlobalSetting 设置全局
func SetGlobalSetting(config *ClientSetting) {

	if config.Proto != "" {
		major, minor, ok := http.ParseHTTPVersion(config.Proto)
		if ok {
			config.ProtoMajor = major
			config.ProtoMinor = minor
		}
	}
}

// Access access url
func Access(method, url string, timeout time.Duration) *Client {
	return New().SetTimeout(timeout).To(method, url)
}

// Head start a Head request conveniently.
func Head(url string, timeout time.Duration) *Client {
	return Access(http.MethodHead, url, timeout)
}

// Get start a GET request conveniently.
func Get(url string, timeout time.Duration) *Client {
	return Access(http.MethodGet, url, timeout)
}

// Post tart a POST request conveniently.
func Post(url string, timeout time.Duration) *Client {
	return Access(http.MethodPost, url, timeout)
}

// Put estart a PUT request conveniently.
func Put(url string, timeout time.Duration) *Client {
	return Access(http.MethodPut, url, timeout)
}

// Delete  start a DELETE request conveniently.
func Delete(url string, timeout time.Duration) *Client {
	return Access(http.MethodDelete, url, timeout)
}

// Patch  start a PATCH request conveniently.
func Patch(url string, timeout time.Duration) *Client {
	return Access(http.MethodPatch, url, timeout)
}

// Options  start a Options request conveniently.
func Options(url string, timeout time.Duration) *Client {
	return Access(http.MethodOptions, url, timeout)
}
