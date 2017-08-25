package httpclient

import "net/http"

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

// Head start a Head request conveniently.
func Head(URL string) *Client {
	return New().To(http.MethodHead, URL)
}

// Get start a GET request conveniently.
func Get(URL string) *Client {
	return New().To(http.MethodGet, URL)
}

// Post tart a POST request conveniently.
func Post(URL string) *Client {
	return New().To(http.MethodPost, URL)
}

// Put estart a PUT request conveniently.
func Put(URL string) *Client {
	return New().To(http.MethodPut, URL)
}

// Delete  start a DELETE request conveniently.
func Delete(URL string) *Client {
	return New().To(http.MethodDelete, URL)
}

// Patch  start a PATCH request conveniently.
func Patch(URL string) *Client {
	return New().To(http.MethodPatch, URL)
}
