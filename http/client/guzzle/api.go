package guzzle

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// Config is used to configure the creation of a client
type Config struct {
	// Address is the address of the Guzzle core
	Address string

	// Scheme is the URI scheme for the Guzzle core
	Scheme string

	// Transport is the Transport to use for the http client.
	Transport *http.Transport

	// HttpClient is the client to use. Default will be
	// used if not provided.
	HttpClient *http.Client

	// HttpAuth is the auth info to use for http access.
	HttpAuth *HttpBasicAuth

	// WaitTime limits how long a Watch will block. If not provided,
	// the agent default values will be used.
	WaitTime time.Duration

	TLSConfig TLSConfig

	ValidateHost bool
}

// HttpBasicAuth is used to authenticate http client with HTTP Basic Authentication
type HttpBasicAuth struct {
	// Username to use for HTTP Basic Authentication
	Username string

	// Password to use for HTTP Basic Authentication
	Password string
}

// TLSConfig is used to generate a TLSClientConfig that's useful for talking to
// Guzzle using TLS.
type TLSConfig struct {
	// Address is the optional address of the Guzzle core. The port, if any
	// will be removed from here and this will be set to the ServerName of the
	// resulting config.
	Address string

	// CAFile is the optional path to the CA certificate used for Guzzle
	// communication, defaults to the system bundle if not specified.
	CAFile string

	// CAPath is the optional path to a directory of CA certificates to use for
	// Guzzle communication, defaults to the system bundle if not specified.
	CAPath string

	// CertFile is the optional path to the certificate for Guzzle
	// communication. If this is set then you need to also set KeyFile.
	CertFile string

	// KeyFile is the optional path to the private key for Guzzle communication.
	// If this is set then you need to also set CertFile.
	KeyFile string

	// InsecureSkipVerify if set to true will disable TLS host verification.
	InsecureSkipVerify bool
}

// request is used to help build up a request
type request struct {
	config *Config
	method string
	url    *url.URL
	params url.Values
	body   io.Reader
	header http.Header
	obj    interface{}
	ctx    context.Context
}

type Client struct {
	config Config
}

func (c *Client) SetTransport(transportFn func() *http.Transport) *Client {
	c.config.Transport = transportFn()
	return c
}

func (c *Client) SetValidateHost(v bool) *Client {
	c.config.ValidateHost = v
	return c
}

func (c *Client) SetBaseAuth(username, password string) *Client {
	c.config.HttpAuth = &HttpBasicAuth{
		Username: username,
		Password: password,
	}
	return c
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {
	// bootstrap the config
	defConfig := DefaultConfig()

	if len(config.Address) == 0 {
		config.Address = defConfig.Address
	}

	if len(config.Scheme) == 0 {
		config.Scheme = defConfig.Scheme
	}

	if config.Transport == nil {
		config.Transport = defConfig.Transport
	} else {
		//If parameter options are configured, they will override the default options
		if config.Transport.Proxy != nil {
			defConfig.Transport.Proxy = config.Transport.Proxy
		}
		if config.Transport.DialContext != nil {
			defConfig.Transport.DialContext = config.Transport.DialContext
		}
		if config.Transport.MaxIdleConns >= 0 {
			defConfig.Transport.MaxIdleConns = config.Transport.MaxIdleConns
		}
		if config.Transport.MaxIdleConnsPerHost >= 0 {
			defConfig.Transport.MaxIdleConnsPerHost = config.Transport.MaxIdleConnsPerHost
		}
		if config.Transport.IdleConnTimeout >= 0 {
			defConfig.Transport.IdleConnTimeout = config.Transport.IdleConnTimeout
		}
		if config.Transport.TLSHandshakeTimeout >= 0 {
			defConfig.Transport.TLSHandshakeTimeout = config.Transport.TLSHandshakeTimeout
		}
		if config.Transport.ExpectContinueTimeout >= 0 {
			defConfig.Transport.ExpectContinueTimeout = config.Transport.ExpectContinueTimeout
		}
		if config.Transport.ResponseHeaderTimeout >= 0 {
			defConfig.Transport.ResponseHeaderTimeout = config.Transport.ResponseHeaderTimeout
		}
		if config.Transport.TLSClientConfig != nil {
			defConfig.Transport.TLSClientConfig = config.Transport.TLSClientConfig
		}
	}

	if config.TLSConfig.Address == "" {
		config.TLSConfig.Address = defConfig.TLSConfig.Address
	}

	if config.TLSConfig.CAFile == "" {
		config.TLSConfig.CAFile = defConfig.TLSConfig.CAFile
	}

	if config.TLSConfig.CAPath == "" {
		config.TLSConfig.CAPath = defConfig.TLSConfig.CAPath
	}

	if config.TLSConfig.CertFile == "" {
		config.TLSConfig.CertFile = defConfig.TLSConfig.CertFile
	}

	if config.TLSConfig.KeyFile == "" {
		config.TLSConfig.KeyFile = defConfig.TLSConfig.KeyFile
	}

	if !config.TLSConfig.InsecureSkipVerify {
		config.TLSConfig.InsecureSkipVerify = defConfig.TLSConfig.InsecureSkipVerify
	}

	if config.WaitTime == 0 {
		config.WaitTime = defConfig.Transport.IdleConnTimeout
	} else {
		defConfig.Transport.IdleConnTimeout = config.WaitTime
	}

	if config.HttpClient == nil {
		var err error
		config.HttpClient, err = NewHttpClient(config.Transport, config.TLSConfig)
		if err != nil {
			return nil, err
		}
	}

	parts := strings.SplitN(config.Address, "://", 2)
	if len(parts) == 2 {
		switch parts[0] {
		case "http":
			config.Scheme = "http"
		case "https":
			config.Scheme = "https"
		case "unix":
			trans := cleanhttp.DefaultTransport()
			trans.DialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", parts[1])
			}
			config.HttpClient = &http.Client{
				Transport: trans,
			}
		default:
			return nil, fmt.Errorf("unknown protocol scheme: %s", parts[0])
		}
		config.Address = parts[1]
	}

	if config.ValidateHost {
		err := ValidateHostConn(config)
		if err != nil {
			return nil, err
		}
	}

	return &Client{config: *config}, nil
}

func ValidateHostConn(config *Config) error {
	var connErr = errors.New("connecting error")
	var host string
	var port string
	var address = config.Address
	if strings.Contains(address, ":") {
		m := strings.Split(address, ":")
		host = m[0]
		port = m[1]
	} else {
		host = address
		port = "80"
	}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), config.WaitTime)
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Close()
		return nil
	}
	return connErr
}

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) {
	// Encode the query parameters
	r.url.RawQuery = r.params.Encode()

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		b, err := encodeBody(r.obj)
		if err != nil {
			return nil, err
		}
		r.body = b
	}

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	if host := r.header.Get("Host"); host != "" {
		req.Host = host
	} else {
		req.Host = r.url.Host
	}
	req.Header = r.header

	// Setup auth
	if r.config.HttpAuth != nil {
		req.SetBasicAuth(r.config.HttpAuth.Username, r.config.HttpAuth.Password)
	}
	if r.ctx != nil {
		return req.WithContext(r.ctx), nil
	}

	return req, nil
}

// NewDoRequest runs a request instance with http client
func (c *Client) NewDoRequest(r *request) (time.Duration, *http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return 0, nil, err
	}
	start := time.Now()
	//c.config.HttpClient.Timeout = c.config.WaitTime

	resp, err := c.config.HttpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	diff := time.Since(start)

	return diff, resp, err
}

// DoNewRequest new a request is used to create a request struct
func (c *Client) DoNewRequest(method, path string) *request {
	r := &request{
		config: &c.config,
		method: method,
		url: &url.URL{
			Scheme: c.config.Scheme,
			Host:   c.config.Address,
			Path:   path,
		},
		params: make(map[string][]string),
		header: make(http.Header),
	}

	return r
}

// RequireOK  is used to wrap doRequest
func RequireOK(d time.Duration, resp *http.Response, e error) *Response {
	return buildResponse(d, resp, e)
}

// generateUnexpectedResponseCodeError consumes the rest of the body, closes
// the body stream and generates an error indicating the status code was
// unexpected.
func generateUnexpectedResponseCodeError(resp *http.Response) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return fmt.Errorf("unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
}

func (resp *Response) Close() {
	err := resp.RawResponse.Body.Close()
	if err != nil {
		return
	}
}

func buildResponse(d time.Duration, resp *http.Response, e error) *Response {

	defaultResp := &Response{}
	defaultResp.StatusCode = resp.StatusCode
	defaultResp.RawResponse = resp
	defaultResp.Header = resp.Header
	defaultResp.Cookies = resp.Cookies()
	defaultResp.Duration = d
	defaultResp.Err = e
	return defaultResp
}

func requireNotFoundOrOK(d time.Duration, resp *http.Response, e error) (bool, time.Duration, *http.Response, error) {
	if e != nil {
		if resp != nil {
			err := resp.Body.Close()
			if err != nil {
				return false, 0, nil, err
			}
		}
		return false, d, nil, e
	}
	switch resp.StatusCode {
	case 200:
		return true, d, resp, nil
	case 404:
		return false, d, resp, nil
	default:
		return false, d, nil, generateUnexpectedResponseCodeError(resp)
	}
}
