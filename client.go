package httpclient

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// ClientSetting http client configure
type ClientSetting struct {
	UserAgent      string
	Proxy          int
	Proto          string
	ProtoMajor     int
	ProtoMinor     int
	Timeout        time.Duration
	Retries        int // if set to -1 means will retry forever
	proxyType      ProxyType
	proxyURL       string
	ProxyTransport *http.Transport
	err            error
}

// ProxyType 代理类型
type ProxyType int

// 代理枚举类型
const (
	NoProxy      ProxyType = iota // 不使用代理
	DefaultProxy                  // 使用系统代理
	CustomProxy                   // 自定义代理
)

var (
	// DefaultSetting default configure
	setting *ClientSetting
)

// Settings 获取设置项
func Settings() *ClientSetting {

	if setting == nil {
		setting = &ClientSetting{
			UserAgent:  "lets-go-go httpclient",
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			proxyType:  NoProxy,
		}
	}

	return setting
}

// SetUserAgent sets the useragent of the request.
func (c *ClientSetting) SetUserAgent(userAgent string) *ClientSetting {
	c.UserAgent = userAgent
	return c
}

// SetProto sets the useragent of the request.
func (c *ClientSetting) SetProto(proto string) *ClientSetting {
	protoMajor, protoMinor, ok := http.ParseHTTPVersion(c.Proto)
	if ok {
		c.Proto = proto
		c.ProtoMajor, c.ProtoMinor = protoMajor, protoMinor

		c.err = errors.New("invalid PROTOCOL version")
	}

	return c
}

// SetProxy sets the address of the proxy which used by the request.
func (c *ClientSetting) SetProxy(proxyType ProxyType, addr string) *ClientSetting {

	if proxyType == NoProxy {
		return c
	}

	// 系统默认代理
	if proxyType == DefaultProxy {
		c.ProxyTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		return c
	}

	// 自定义代理
	u, err := url.Parse(addr)

	if err != nil {
		c.err = err
		return c
	}

	switch u.Scheme {
	case "http", "https":
		c.ProxyTransport = &http.Transport{
			Proxy: http.ProxyURL(u),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	case "socks5":
		dialer, err := proxy.FromURL(u, proxy.Direct)

		if err != nil {
			c.err = err
			return c
		}

		c.ProxyTransport = &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	}

	return c
}
